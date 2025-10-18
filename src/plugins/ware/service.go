package ware

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/tray"
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

func (sw *ServiceWare) Complete(ctx *common.GenContext) error { return nil }
