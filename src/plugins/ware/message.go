package ware

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/tray"
	"github.com/ncuhome/cato/src/plugins/utils"
)

const (
	modelImportAlias     = "model"
	repoImportAlias      = "repo"
	repoFetchOneFuncName = "FetchOne"
	repoFetchAllFuncName = "FetchAll"
)

type MessageWare struct {
	message       *protogen.Message
	fieldGenFiles []*models.GenerateFileDesc
}

func NewMessageWare(msg *protogen.Message) *MessageWare {
	mp := new(MessageWare)
	mp.message = msg
	return mp
}

// RegisterContext because generate a file from a message, so a file-level writer for a message generates progress
func (mw *MessageWare) RegisterContext(gc *common.GenContext) *common.GenContext {
	mc := tray.NewMessageTray()
	ctx := gc.WithMessage(mw.message, mc)
	return ctx
}

func (mw *MessageWare) GetSubWares() []WorkWare {
	fields := make([]WorkWare, len(mw.message.Fields))
	for i, field := range mw.message.Fields {
		fields[i] = NewFieldWare(field)
	}
	return fields
}

func (mw *MessageWare) GetDescriptor() protoreflect.Descriptor {
	return mw.message.Desc
}

func (mw *MessageWare) Active(ctx *common.GenContext) (bool, error) {
	return CommonWareActive(ctx, mw)
}

func (mw *MessageWare) filename() string {
	patterns := utils.SplitCamelWords(mw.message.GoIdent.GoName)
	mapper := utils.GetStringsMapper(generated.FieldMapper_CATO_FIELD_MAPPER_SNAKE_CASE)
	return mapper(patterns)
}

func (mw *MessageWare) GetExtraFiles(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	files := make([]*models.GenerateFileDesc, 0)
	gen := []fileGenerator{
		mw.generateModelExtendFiles,
		mw.generateModelRepoFiles,
		mw.generateModelRdbFiles,
	}
	var err error
	for _, f := range gen {
		gfs, ex := f(ctx)
		if ex != nil {
			err = errors.Join(err, ex)
		}
		files = append(files, gfs...)
	}
	return append(mw.fieldGenFiles, files...), err
}

func (mw *MessageWare) generateModelExtendFiles(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	sw := new(strings.Builder)
	tmpl := config.GetTemplate(config.TableExtendTmpl)
	fc := ctx.GetNowFileContainer()
	mc := ctx.GetNowMessageContainer()
	modelPack := fc.GetCatoPackage()
	if modelPack == nil || !mc.IsNeedExtraFile() {
		return []*models.GenerateFileDesc{}, nil
	}
	pack := &packs.TableExtendTmplPack{
		PackageName:   utils.GetGoPackageName(modelPack.ImportPath),
		ExtendMethods: mc.GetExtra(),
	}
	err := tmpl.Execute(sw, pack)
	if err != nil {
		return nil, err
	}
	filename := filepath.Join(modelPack.ImportPath, fmt.Sprintf("%s_extend.go", mw.filename()))
	return []*models.GenerateFileDesc{{
		Name:        filename,
		Content:     sw.String(),
		CheckExists: true,
	}}, nil
}

func (mw *MessageWare) Complete(ctx *common.GenContext) error {
	err := mw.completeCols(ctx)
	if err != nil {
		return err
	}
	keyParams := mw.loadKeyTmplPacks(ctx)
	err = mw.completeRepo(ctx, keyParams)
	if err != nil {
		return err
	}
	return mw.completeMessageContent(ctx)
}

func (mw *MessageWare) completeCols(ctx *common.GenContext) error {
	mc := ctx.GetNowMessageContainer()
	cols := mc.GetScopeCols()
	if len(cols) == 0 {
		return nil
	}
	tmpl := config.GetTemplate(config.TableColTmpl)
	pack := &packs.TableColTmplPack{
		MessageTypeName: ctx.GetNowMessageTypeName(),
		Cols:            cols,
	}
	err := tmpl.Execute(mc.BorrowMethodsWriter(), pack)
	if err != nil {
		return err
	}
	// check if it has col groups
	colGroups := mc.GetColGroups()
	for groupName, cls := range colGroups {
		groupPack := &packs.ColsGroupPack{
			MessageTypeName: ctx.GetNowMessageTypeName(),
			GroupName:       groupName,
			Cols:            cls,
		}
		colGroupTmpl := config.GetTemplate(config.ColsGroupTmpl)

		err = errors.Join(colGroupTmpl.Execute(mc.BorrowMethodsWriter(), groupPack))
	}
	return err
}

type repoCompleteParam struct {
	path       *models.Import
	modelType  string
	isRepoSame bool
	tmpls      []string
	ukTmpls    []string
	sgTmpls    []string
	writer     func() io.Writer
}

func (mw *MessageWare) completeRepo(ctx *common.GenContext, params []*packs.RepoKeyFuncTmplPack) error {
	fc := ctx.GetNowFileContainer()
	mc := ctx.GetNowMessageContainer()
	runParams := make([]*repoCompleteParam, 0)
	repoPackage := fc.GetRepoPackage()
	if repoPackage != nil && !repoPackage.IsEmpty() {
		repoParam := &repoCompleteParam{
			modelType:  ctx.GetNowMessageTypeName(),
			path:       repoPackage,
			isRepoSame: fc.GetCatoPackage().IsSame(repoPackage),
			tmpls:      []string{config.RepoFetchTmpl},
			ukTmpls:    []string{config.RepoUpdateTmpl, config.RepoDeleteTmpl},
			sgTmpls:    []string{config.RepoInsertTmpl},
			writer:     mc.BorrowRepoWriter,
		}
		runParams = append(runParams, repoParam)
	}
	rdbPackage := fc.GetRdbRepoPackage()
	if rdbPackage != nil && !rdbPackage.IsEmpty() {
		rdbParam := &repoCompleteParam{
			modelType:  ctx.GetNowMessageTypeName(),
			path:       rdbPackage,
			isRepoSame: fc.GetCatoPackage().IsSame(rdbPackage),
			tmpls:      []string{config.RdbFetchTmpl},
			ukTmpls:    []string{config.RdbUpdateTmpl, config.RdbDeleteTmpl},
			sgTmpls:    []string{config.RdbInsertTmpl},
			writer:     mc.BorrowRdbWriter,
		}
		runParams = append(runParams, rdbParam)
	}
	return mw.repoImplRunner(runParams, params)
}

func (mw *MessageWare) repoImplRunner(runParams []*repoCompleteParam, params []*packs.RepoKeyFuncTmplPack) error {
	var err error
	for _, rp := range runParams {
		modelType := rp.modelType
		if !rp.isRepoSame {
			modelType = fmt.Sprintf("%s.%s", modelImportAlias, rp.modelType)
		}
		noKeyPack := &packs.NoKeyFuncTmplPack{ModelType: modelType}
		for _, param := range params {
			cparam := param.Copy()
			cparam.IsModelAnotherPackage = rp.isRepoSame
			cparam.ModelType = modelType
			cparam.Tmpls = append(cparam.Tmpls, rp.tmpls...)
			if cparam.IsUniqueKey {
				cparam.Tmpls = append(cparam.Tmpls, rp.ukTmpls...)
				cparam.FetchReturnType = fmt.Sprintf("*%s", cparam.ModelType)
			} else {
				cparam.FetchReturnType = fmt.Sprintf("[]*%s", cparam.ModelType)
			}
			for _, tmpl := range cparam.Tmpls {
				err = errors.Join(err, config.GetTemplate(tmpl).Execute(rp.writer(), cparam))
			}
		}
		for _, tmpl := range rp.sgTmpls {
			err = errors.Join(err, config.GetTemplate(tmpl).Execute(rp.writer(), noKeyPack))
		}
	}
	return err
}

func (mw *MessageWare) loadKeyTmplPacks(ctx *common.GenContext) []*packs.RepoKeyFuncTmplPack {
	keys := ctx.GetNowMessageContainer().GetScopeKeys()
	if len(keys) == 0 {
		return nil
	}
	keysTmplPack := make([]*packs.RepoKeyFuncTmplPack, 0)
	for _, key := range keys {
		keyType := key.KeyType
		pack := &packs.RepoKeyFuncTmplPack{
			KeyNameCombine:    key.GetFieldNameCombine(),
			ModelPackage:      ctx.GetCatoPackage(),
			ModelPackageAlias: modelImportAlias,
			Tmpls:             make([]string, 0),
		}
		packParams := make([]*packs.RepoKeyFuncTmplPackParam, len(key.Fields))
		for index := range key.Fields {
			packParams[index] = &packs.RepoKeyFuncTmplPackParam{
				FieldName: key.Fields[index].Name,
				ParamName: key.Fields[index].AsParamName(),
			}
		}
		pack.Params = packParams
		switch keyType {
		// unique and primary key will have FetchOne, UpdateBy, DeleteByMethod
		case generated.DBKeyType_CATO_DB_KEY_TYPE_PRIMARY, generated.DBKeyType_CATO_DB_KEY_TYPE_UNIQUE:
			pack.FetchFuncName = repoFetchOneFuncName
			pack.IsUniqueKey = true
		case generated.DBKeyType_CATO_DB_KEY_TYPE_COMBINE, generated.DBKeyType_CATO_DB_KEY_TYPE_INDEX:
			pack.FetchFuncName = repoFetchAllFuncName
		}
		keysTmplPack = append(keysTmplPack, pack)
	}
	return keysTmplPack
}

func (mw *MessageWare) completeMessageContent(ctx *common.GenContext) error {
	tmpl := config.GetTemplate(config.ModelTmpl)
	mc := ctx.GetNowMessageContainer()
	fc := ctx.GetNowFileContainer()
	modelPack := fc.GetCatoPackage()
	if modelPack == nil || modelPack.IsEmpty() {
		return nil
	}
	pack := &packs.ModelContentTmplPack{
		PackageName: utils.GetGoPackageName(modelPack.GetPath()),
		ModelName:   mw.message.GoIdent.GoName,
		Fields:      mc.GetField(),
		Methods:     mc.GetMethods(),
	}
	messageName := mw.message.GoIdent.GoName
	err := tmpl.Execute(fc.BorrowContentWriter(), pack)
	if err != nil {
		return fmt.Errorf("model %s content exec tmpl error, %#v", messageName, err)
	}
	for _, iv := range mc.GetImports() {
		_, err = fc.BorrowImportWriter().Write([]byte(iv))
		if err != nil {
			return fmt.Errorf("model %s import exec tmpl error, %#v", messageName, err)
		}
	}
	return nil
}

func (mw *MessageWare) generateModelRepoFiles(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	fc := ctx.GetNowFileContainer()
	repoPack := fc.GetRepoPackage()
	if repoPack == nil || repoPack.IsEmpty() {
		return []*models.GenerateFileDesc{}, nil
	}
	modelPack := fc.GetCatoPackage()
	mc := ctx.GetNowMessageContainer()
	pack := &packs.RepoTmplPack{
		RepoPackageName:       utils.GetGoPackageName(repoPack.ImportPath),
		IsModelAnotherPackage: modelPack.IsSame(repoPack),
		ModelPackageAlias:     modelImportAlias,
		ModelPackage:          modelPack.ImportPath,
		RepoFuncs:             mc.GetRepo(),
		RdbPackage:            fc.GetRdbRepoPackage().ImportPath,
		ModelType:             ctx.GetNowMessageTypeName(),
	}
	files := make([]*models.GenerateFileDesc, 0)
	sw := new(strings.Builder)
	err := config.GetTemplate(config.RepoTmpl).Execute(sw, pack)
	if err != nil {
		return nil, err
	}
	filename := filepath.Join(repoPack.ImportPath, fmt.Sprintf("%s_repo.cato.go", mw.filename()))
	files = append(files, &models.GenerateFileDesc{
		Name:        filename,
		Content:     sw.String(),
		CheckExists: false,
	})
	extraSw := new(strings.Builder)
	err = config.GetTemplate(config.RepoExtTmpl).Execute(extraSw, pack)
	if err != nil {
		return nil, err
	}
	filename = filepath.Join(repoPack.ImportPath, fmt.Sprintf("extension.go"))
	files = append(files, &models.GenerateFileDesc{
		Name:        filename,
		Content:     extraSw.String(),
		CheckExists: true,
	})
	return files, nil
}

func (mw *MessageWare) generateModelRdbFiles(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	fc := ctx.GetNowFileContainer()
	repoPack := fc.GetRepoPackage()
	modelPack := fc.GetCatoPackage()
	rdbPack := fc.GetRdbRepoPackage()
	if rdbPack == nil || rdbPack.IsEmpty() {
		return []*models.GenerateFileDesc{}, nil
	}
	mc := ctx.GetNowMessageContainer()
	pack := &packs.RdbTmplPack{
		RdbRepoPackage:        utils.GetGoPackageName(rdbPack.ImportPath),
		IsModelAnotherPackage: modelPack.IsSame(rdbPack),
		ModelPackageAlias:     modelImportAlias,
		ModelPackage:          modelPack.ImportPath,
		RdbRepoFuncs:          mc.GetRdb(),
		IsRepoAnotherPackage:  repoPack.IsSame(rdbPack),
		RepoPackageAlias:      repoImportAlias,
		RepoPackage:           repoPack.ImportPath,
		ModelType:             ctx.GetNowMessageTypeName(),
	}
	pack.FetchOneReturnType = fmt.Sprintf("*%s", pack.ModelType)
	if !modelPack.IsSame(rdbPack) {
		pack.FetchOneReturnType = fmt.Sprintf("*%s.%s", modelImportAlias, pack.ModelType)
	}
	pack.FetchAllReturnType = fmt.Sprintf("[]*%s", rdbPack.ImportPath)
	if !modelPack.IsSame(rdbPack) {
		pack.FetchAllReturnType = fmt.Sprintf("[]*%s.%s", modelImportAlias, pack.ModelType)
	}
	files := make([]*models.GenerateFileDesc, 0)
	sw := new(strings.Builder)
	err := config.GetTemplate(config.RdbTmpl).Execute(sw, pack)
	if err != nil {
		return nil, err
	}
	filename := filepath.Join(rdbPack.ImportPath, fmt.Sprintf("%s_rdb.cato.go", mw.filename()))
	files = append(files, &models.GenerateFileDesc{
		Name:        filename,
		Content:     sw.String(),
		CheckExists: false,
	})
	return files, nil
}

func (mw *MessageWare) StoreExtraFiles(files []*models.GenerateFileDesc) {
	mw.fieldGenFiles = append(mw.fieldGenFiles, files...)
}
