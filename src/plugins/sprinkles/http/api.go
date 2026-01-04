package http

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
)

func init() {
	sprinkles.Register(func() sprinkles.Sprinkle {
		return new(ApiSprinkle)
	})
}

type ApiSprinkle struct {
	value *generated.HttpOptions
}

func (a *ApiSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_HttpOpt
}

func (a *ApiSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.ServiceDescriptor)
	return ok
}

func (a *ApiSprinkle) Init(value interface{}) {
	a.value, _ = value.(*generated.HttpOptions)
}

func (a *ApiSprinkle) Register(ctx *common.GenContext) error {
	if a.value == nil {
		return nil
	}
	fc := ctx.GetNowServiceContainer()
	fc.SetRouterBasePath(a.value.GroupPrefix)
	fc.SetRegisterHttpApi(a.value.AsHttpService)
	return nil
}
