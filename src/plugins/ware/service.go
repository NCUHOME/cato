package ware

import (
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
	service *protogen.Service
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
	return active(ctx, sw)
}

func (sw *ServiceWare) Complete(_ *common.GenContext) error { return nil }

func (sw *ServiceWare) GetFiles(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	// check has package file entry
	fc := ctx.GetNowFileContainer()
	locate := fc.GetHttpHandlerPackage()
	if locate.IsEmpty() {
		return []*models.GenerateFileDesc{}, nil
	}
	// generate http service register handler and generate
	generator := []fileGenerator{
		sw.generateProtocol,
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
	return files, nil
}

func (sw *ServiceWare) generateHandler(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	locate := ctx.GetNowFileContainer().GetHttpHandlerPackage()
	sc := ctx.GetNowServiceContainer()
	service := ctx.GetNowService()
	tmpl := config.GetTemplate(config.HttpHandlerTmpl)
	pack := &packs.HttpHandlerTmplPack{
		HttpHandlerPackage:     utils.GetGoPackageName(locate.GetPath()),
		ServiceName:            service.GoName,
		RegisterServiceMethods: sc.GetRegisters(),
		RouterBasePath:         sc.GetRouterBasePath(),
	}
	sb := new(strings.Builder)
	err := tmpl.Execute(sb, pack)
	if err != nil {
		return nil, err
	}
	filename := filepath.Join(locate.ImportPath, "handlers.cato.go")
	return []*models.GenerateFileDesc{{
		Name:        filename,
		Content:     sb.String(),
		CheckExists: false,
	}}, nil
}

func (sw *ServiceWare) generateProtocol(ctx *common.GenContext) ([]*models.GenerateFileDesc, error) {
	locate := ctx.GetNowFileContainer().GetHttpHandlerPackage()
	sc := ctx.GetNowServiceContainer()
	service := ctx.GetNowService()
	tmpl := config.GetTemplate(config.HttpProtocolTmpl)
	imports := append([]string{}, ctx.GetNowFileContainer().GetImports()...)
	pack := &packs.HttpProtocolTmplPack{
		HttpHandlerPackage:    utils.GetGoPackageName(locate.GetPath()),
		ProtocolParamPackages: append(imports, sc.GetExtraImport()...),
		ServiceName:           service.GoName,
		Methods:               sc.GetMethods(),
		TierMethods:           sc.GetTiers(),
	}
	sb := new(strings.Builder)
	err := tmpl.Execute(sb, pack)
	if err != nil {
		return nil, err
	}
	filename := filepath.Join(locate.ImportPath, "protocol.cato.go")
	return []*models.GenerateFileDesc{{
		Name:        filename,
		Content:     sb.String(),
		CheckExists: false,
	}}, nil
}
