package plugins

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/butter"
	"github.com/ncuhome/cato/src/plugins/cheese"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/utils"
)

const (
	modelImportAlias     = "model"
	repoImportAlias      = "repo"
	repoFetchOneFuncName = "FetchOne"
	repoFetchAllFuncName = "FetchAll"
)

type MessageWorker struct {
	message *protogen.Message
}

func NewMessageWorker(msg *protogen.Message) *MessageWorker {
	mp := new(MessageWorker)
	mp.message = msg
	return mp
}

// RegisterContext because generate a file from a message, so a file-level writer for a message generates progress
func (mw *MessageWorker) RegisterContext(gc *common.GenContext) *common.GenContext {
	mc := cheese.NewMessageCheese()
	ctx := gc.WithMessage(mw.message, mc)
	return ctx
}

func (mw *MessageWorker) filename() string {
	patterns := utils.SplitCamelWords(mw.message.GoIdent.GoName)
	mapper := utils.GetStringsMapper(generated.FieldMapper_CATO_FIELD_MAPPER_SNAKE_CASE)
	return mapper(patterns)
}

func (mw *MessageWorker) GetFiles(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	files := make([]*models.GenerateFileDesc, 0)
	gen := []func(ctx *common.GenContext) ([]*models.GenerateFileDesc, error){
		mw.generateModelFile, mw.generateModelExtendFile, mw.generateModelRepoFiles, mw.generateModelRdbFiles,
	}
	var err error
	for _, f := range gen {
		gfs, ex := f(ctx)
		if ex != nil {
			err = errors.Join(err, ex)
		}
		files = append(files, gfs...)
	}
	return files, err
}

func (mw *MessageWorker) generateModelFile(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	sw := new(strings.Builder)
	tmpl := config.GetTemplate(config.ModelTmpl)
	mc := ctx.GetNowMessageContainer()
	fc := ctx.GetNowFileContainer()
	modelPack := fc.GetCatoPackage()
	if modelPack == nil || modelPack.IsEmpty() {
		return []*models.GenerateFileDesc{}, nil
	}
	pack := &packs.ModelContentTmplPack{
		PackageName: utils.GetGoPackageName(modelPack.GetPath()),
		ModelName:   mw.message.GoIdent.GoName,
		Imports:     append(fc.GetImports(), mc.GetImports()...),
		Fields:      mc.GetField(),
		Methods:     mc.GetMethods(),
	}
	err := tmpl.Execute(sw, pack)
	if err != nil {
		log.Fatalln("[-] models plugger exec tmpl error, ", err)
		return nil, err
	}
	filename := filepath.Join(modelPack.ImportPath, fmt.Sprintf("%s.cato.go", mw.filename()))
	return []*models.GenerateFileDesc{{
		Name:        filename,
		Content:     sw.String(),
		CheckExists: false,
	}}, nil
}

func (mw *MessageWorker) generateModelExtendFile(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	sw := new(strings.Builder)
	tmpl := config.GetTemplate(config.TableExtendTmpl)
	fc := ctx.GetNowFileContainer()
	mc := ctx.GetNowMessageContainer()
	modelExtPack := fc.GetCatoExtPackage()
	if modelExtPack == nil || modelExtPack.IsEmpty() {
		return []*models.GenerateFileDesc{}, nil
	}
	pack := &packs.TableExtendTmplPack{
		PackageName:   utils.GetGoPackageName(modelExtPack.ImportPath),
		ExtendMethods: mc.GetExtra(),
	}
	err := tmpl.Execute(sw, pack)
	if err != nil {
		log.Fatalln("[-] plugger model exec extend tmpl error, ", err)
		return nil, err
	}
	filename := filepath.Join(modelExtPack.ImportPath, fmt.Sprintf("%s_extend.go", mw.filename()))
	return []*models.GenerateFileDesc{{
		Name:        filename,
		Content:     sw.String(),
		CheckExists: true,
	}}, nil
}

func (mw *MessageWorker) generateModelRepoFiles(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
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
	}
	files := make([]*models.GenerateFileDesc, 0)
	sw := new(strings.Builder)
	err := config.GetTemplate(config.RepoTmpl).Execute(sw, pack)
	if err != nil {
		log.Fatalln("[-] plugger repo exec tmpl error, ", err)
		return nil, err
	}
	filename := filepath.Join(repoPack.ImportPath, fmt.Sprintf("%s_repo.cato.go", mw.filename()))
	files = append(files, &models.GenerateFileDesc{
		Name:        filename,
		Content:     sw.String(),
		CheckExists: false,
	})
	extraSw := new(strings.Builder)
	err = config.GetTemplate(config.RepoRepoTmpl).Execute(extraSw, pack)
	if err != nil {
		log.Fatalln("[-] plugger repo repo tmpl error, ", err)
		return nil, err
	}
	filename = filepath.Join(repoPack.ImportPath, fmt.Sprintf("%s_repo.go", mw.filename()))
	files = append(files, &models.GenerateFileDesc{
		Name:        filename,
		Content:     extraSw.String(),
		CheckExists: true,
	})
	return files, nil
}

func (mw *MessageWorker) generateModelRdbFiles(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
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
	// todo rdb engine funcs can move into messageCheese
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
		log.Fatalln("[-] plugger rdb exec tmpl error, ", err)
		return nil, err
	}
	filename := filepath.Join(rdbPack.ImportPath, fmt.Sprintf("%s_rdb.cato.go", mw.filename()))
	files = append(files, &models.GenerateFileDesc{
		Name:        filename,
		Content:     sw.String(),
		CheckExists: false,
	})
	extraSw := new(strings.Builder)
	err = config.GetTemplate(config.RdbRepoTmpl).Execute(extraSw, pack)
	if err != nil {
		log.Fatalln("[-] plugger rdb repo tmpl error, ", err)
		return nil, err
	}
	filename = filepath.Join(rdbPack.ImportPath, fmt.Sprintf("%s_rdb.go", mw.filename()))
	files = append(files, &models.GenerateFileDesc{
		Name:        filename,
		Content:     extraSw.String(),
		CheckExists: true,
	})
	return files, nil
}

func (mw *MessageWorker) Active(ctx *common.GenContext) (bool, error) {
	descriptor := protodesc.ToDescriptorProto(mw.message.Desc)
	butters := butter.ChooseButter(mw.message.Desc)

	for index := range butters {
		if !proto.HasExtension(descriptor.Options, butters[index].FromExtType()) {
			continue
		}
		value := proto.GetExtension(descriptor.Options, butters[index].FromExtType())
		butters[index].Init(value)
		err := butters[index].Register(ctx)
		if err != nil {
			return false, err
		}
	}
	// for fields
	for _, field := range mw.message.Fields {
		fp := NewFieldCheese(field)
		fieldCtx := fp.RegisterContext(ctx)
		_, err := fp.Active(fieldCtx)
		if err != nil {
			return false, err
		}
		err = fp.Complete(fieldCtx)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (mw *MessageWorker) Complete(ctx *common.GenContext) error {
	err := mw.completeCols(ctx)
	if err != nil {
		return err
	}
	keyParams := mw.loadKeyTmplPacks(ctx)
	err = mw.completeRepo(ctx, keyParams)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MessageWorker) completeCols(ctx *common.GenContext) error {
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
	return tmpl.Execute(mc.BorrowMethodsWriter(), pack)
}

type repoCompleteParam struct {
	basePath *models.Import
	path     *models.Import
	tmpls    []string
	ukTmpls  []string
	sgTmpls  []string
	writer   func() io.Writer
}

func (mw *MessageWorker) completeRepo(ctx *common.GenContext, params []*packs.RepoFuncTmplPack) error {
	fc := ctx.GetNowFileContainer()
	mc := ctx.GetNowMessageContainer()
	runParams := make([]*repoCompleteParam, 0)
	repoPackage := fc.GetRepoPackage()
	if repoPackage != nil && !repoPackage.IsEmpty() {
		repoParam := &repoCompleteParam{
			basePath: fc.GetCatoPackage(),
			path:     repoPackage,
			tmpls:    []string{config.RepoFetchTmpl},
			ukTmpls:  []string{config.RepoUpdateTmpl, config.RepoDeleteTmpl},
			sgTmpls:  []string{config.RepoInsertTmpl},
			writer:   mc.BorrowRepoWriter,
		}
		runParams = append(runParams, repoParam)
	}
	rdbPackage := fc.GetRdbRepoPackage()
	if rdbPackage != nil && !rdbPackage.IsEmpty() {
		rdbParam := &repoCompleteParam{
			basePath: fc.GetCatoPackage(),
			path:     rdbPackage,
			tmpls:    []string{config.RdbFetchTmpl},
			ukTmpls:  []string{config.RdbUpdateTmpl, config.RdbDeleteTmpl},
			sgTmpls:  []string{config.RdbInsertTmpl},
			writer:   mc.BorrowRdbWriter,
		}
		runParams = append(runParams, rdbParam)
	}
	var err error
	for _, rp := range runParams {
		err = errors.Join(err, mw.repoInsRunner(rp, params))
	}
	return err
}

func (mw *MessageWorker) repoInsRunner(rcp *repoCompleteParam, params []*packs.RepoFuncTmplPack) error {
	isRepoSame := rcp.basePath.IsSame(rcp.path)
	var (
		err     error
		sgParam *packs.RepoFuncTmplPack
	)

	for _, param := range params {
		cparam := param.Copy()
		cparam.IsModelAnotherPackage = isRepoSame
		if !cparam.IsModelAnotherPackage {
			cparam.ModelType = fmt.Sprintf("%s.%s", cparam.ModelPackageAlias, param.ModelType)
		}
		cparam.Tmpls = append(cparam.Tmpls, rcp.tmpls...)
		if cparam.IsUniqueKey {
			cparam.Tmpls = append(cparam.Tmpls, rcp.ukTmpls...)
			cparam.FetchReturnType = fmt.Sprintf("*%s", cparam.ModelType)
		} else {
			cparam.FetchReturnType = fmt.Sprintf("[]*%s", cparam.ModelType)
		}
		if sgParam == nil {
			sgParam = cparam.Copy()
		}
		for _, tmpl := range cparam.Tmpls {
			err = errors.Join(err, config.GetTemplate(tmpl).Execute(rcp.writer(), cparam))
		}
	}
	// todo need a better way to deal once tmpl
	for _, tmpl := range rcp.sgTmpls {
		err = errors.Join(err, config.GetTemplate(tmpl).Execute(rcp.writer(), sgParam))
	}
	return err
}

func (mw *MessageWorker) loadKeyTmplPacks(ctx *common.GenContext) []*packs.RepoFuncTmplPack {
	keys := ctx.GetNowMessageContainer().GetScopeKeys()
	if len(keys) == 0 {
		return nil
	}
	keysTmplPack := make([]*packs.RepoFuncTmplPack, 0)
	for _, key := range keys {
		keyType := key.KeyType
		pack := &packs.RepoFuncTmplPack{
			KeyNameCombine:    key.GetFieldNameCombine(),
			ModelType:         ctx.GetNowMessageTypeName(),
			ModelPackage:      ctx.GetCatoPackage(),
			ModelPackageAlias: modelImportAlias,
			Tmpls:             make([]string, 0),
		}
		packParams := make([]*packs.RepoFuncTmplPackParam, len(key.Fields))
		for index := range key.Fields {
			packParams[index] = &packs.RepoFuncTmplPackParam{
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
		case generated.DBKeyType_CATO_DB_KEY_TYPE_COMBINE:
			pack.FetchFuncName = repoFetchAllFuncName
		}
		keysTmplPack = append(keysTmplPack, pack)
	}
	return keysTmplPack
}
