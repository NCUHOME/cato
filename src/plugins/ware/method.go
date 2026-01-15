package ware

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/tray"
)

type MethodWare struct {
	method *protogen.Method
}

func NewMethodWare(m *protogen.Method) *MethodWare {
	return &MethodWare{method: m}
}

func (mw *MethodWare) GetExtraFiles(_ *common.GenContext) ([]*models.GenerateFileDesc, error) {
	return []*models.GenerateFileDesc{}, nil
}

func (mw *MethodWare) StoreExtraFiles(_ []*models.GenerateFileDesc) {}

func (mw *MethodWare) GetDescriptor() protoreflect.Descriptor {
	return mw.method.Desc
}

func (mw *MethodWare) GetSubWares() []WorkWare {
	return []WorkWare{}
}

func (mw *MethodWare) RegisterContext(gc *common.GenContext) *common.GenContext {
	mc := tray.NewMethodTray()
	mc.SetRequest(mw.method.Input)
	mc.SetResponse(mw.method.Output)
	ctx := gc.WithMethod(mw.method, mc)
	return ctx
}

func (mw *MethodWare) Active(ctx *common.GenContext) (bool, error) {
	return CommonWareActive(ctx, mw)
}

func (mw *MethodWare) Complete(_ *common.GenContext) error {
	return nil
}
