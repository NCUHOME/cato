package ware

import (
	"errors"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/tray"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type ServiceWare struct {
	service        *protogen.Service
	methodGenFiles []*models.GenerateFileDesc
}

func NewServiceWare(service *protogen.Service) *ServiceWare {
	return &ServiceWare{service: service}
}

func (sw *ServiceWare) GetDescriptor() protoreflect.Descriptor {
	return sw.service.Desc
}

func (sw *ServiceWare) GetSubWares() []WorkWare {
	subs := make([]WorkWare, len(sw.service.Methods))
	for i, m := range sw.service.Methods {
		subs[i] = NewMethodWare(m)
	}
	return subs
}

func (sw *ServiceWare) RegisterContext(gc *common.GenContext) *common.GenContext {
	sc := tray.NewServiceTray()
	ctx := gc.WithService(sw.service, sc)
	return ctx
}

func (sw *ServiceWare) Active(ctx *common.GenContext) (bool, error) {
	return CommonWareActive(ctx, sw)
}

func (sw *ServiceWare) Complete(ctx *common.GenContext) error {
	sc := ctx.GetNowServiceContainer()
	fc := ctx.GetNowFileContainer()
	for _, mi := range sc.GetExtraImport() {
		_, err := fc.BorrowImportWriter().Write([]byte(mi))
		if err != nil {
			return err
		}
	}
	return nil
}

func (sw *ServiceWare) StoreExtraFiles(files []*models.GenerateFileDesc) {
	sw.methodGenFiles = append(sw.methodGenFiles, files...)
}

func (sw *ServiceWare) GetExtraFiles(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	// check has package file entry
	fc := ctx.GetNowFileContainer()
	locate := fc.GetCatoPackage()
	if locate.IsEmpty() {
		return []*models.GenerateFileDesc{}, nil
	}
	// generate http service register handler and generate
	generator := []fileGenerator{
		sw.generateHandler,
	}
	files := make([]*models.GenerateFileDesc, 0)
	for _, f := range generator {
		fs, err := f(ctx)
		if err != nil {
			return nil, err
		}
		files = append(files, fs...)
	}
	return append(sw.methodGenFiles, files...), nil
}

func (sw *ServiceWare) generateHandler(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	locate := ctx.GetNowFileContainer().GetCatoPackage()
	sc := ctx.GetNowServiceContainer()
	service := ctx.GetNowService()
	pack := &packs.HttpProtocolTmplPack{
		ServicePackage:   utils.GetGoPackageName(locate.GetPath()),
		ImportPackages:   append(sc.GetExtraImport(), ctx.GetNowFileContainer().GetImports()...),
		ServiceName:      service.GoName,
		ServiceNameInner: utils.FirstLower(service.GoName),
		Methods:          sc.GetApis(),
		RouterBasePath:   sc.GetRouterBasePath(),
		HttpMethods:      sc.GetMethods(),
		HttpRouters:      sc.GetRouters(),
	}
	handlers, handlersCustom := new(strings.Builder), new(strings.Builder)
	apis, apiCustom := new(strings.Builder), new(strings.Builder)
	// handlers file
	err := errors.Join(
		config.GetTemplate(config.HandlersTmpl).Execute(handlers, pack),
		config.GetTemplate(config.HandlersCustomTmpl).Execute(handlersCustom, pack),
		config.GetTemplate(config.ApiCustomTmpl).Execute(apiCustom, pack),
		config.GetTemplate(config.ApiTmpl).Execute(apis, pack),
	)
	if err != nil {
		return nil, err
	}
	files := []*models.GenerateFileDesc{
		{
			Name:        filepath.Join(locate.ImportPath, "handlers_custom.go"),
			Content:     handlersCustom.String(),
			CheckExists: true,
		},
		{
			Name:        filepath.Join(locate.ImportPath, "handlers.cato.go"),
			Content:     handlers.String(),
			CheckExists: false,
		},
		{
			Name:        filepath.Join(locate.ImportPath, "api_custom.go"),
			Content:     apiCustom.String(),
			CheckExists: true,
		},
		{
			Name:        filepath.Join(locate.ImportPath, "api.cato.go"),
			Content:     apis.String(),
			CheckExists: false,
		},
	}
	return files, nil
}
