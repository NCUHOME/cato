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
	return sw.completeContent(ctx)
}

func (sw *ServiceWare) completeContent(ctx *common.GenContext) error {
	locate := ctx.GetNowFileContainer().GetCatoPackage()
	sc := ctx.GetNowServiceContainer()
	service := ctx.GetNowService()
	tmpl := config.GetTemplate(config.HttpProtocolTmpl)
	pack := &packs.HttpProtocolTmplPack{
		HttpHandlerPackage: utils.GetGoPackageName(locate.GetPath()),
		ServiceName:        service.GoName,
		Methods:            sc.GetMethods(),
		TierMethods:        sc.GetTiers(),
	}
	fc := ctx.GetNowFileContainer()
	err := tmpl.Execute(fc.BorrowContentWriter(), pack)
	if err != nil {
		return err
	}
	for _, mi := range sc.GetExtraImport() {
		_, err = fc.BorrowImportWriter().Write([]byte(mi))
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
