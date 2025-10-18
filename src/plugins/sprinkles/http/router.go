package http

import (
	"errors"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
)

func init() {
	sprinkles.Register(func() sprinkles.Sprinkle {
		return new(RouterSprinkle)
	})
}

type RouterSprinkle struct {
	value *generated.RouterOptions
}

func (r *RouterSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_RouterOpt
}

func (r *RouterSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.MethodDescriptor)
	return ok
}

func (r *RouterSprinkle) Init(value interface{}) {
	r.value, _ = value.(*generated.RouterOptions)
}

func (r *RouterSprinkle) Register(ctx *common.GenContext) error {
	if r.value == nil {
		return nil
	}
	// check has router message
	relativePath := r.value.GetRouter()
	if relativePath == "" {
		return nil
	}
	// register router into handler service and protocol interface
	return errors.Join(
		r.registerHandler(ctx),
		r.registerMethod(ctx),
		r.RegisterTier(ctx),
	)
}

// registerHandler will register method http handler wrap implement into handler service
func (r *RouterSprinkle) registerHandler(ctx *common.GenContext) error {
	method := ctx.GetNowMethod()
	pack := &packs.RouteRegisterTmplPack{
		HttpMethod:     r.value.Method,
		HttpMethodPath: r.value.Router,
		MethodName:     method.GoName,
	}
	writer := ctx.GetNowServiceContainer().BorrowRegistersWriter()
	return config.GetTemplate(config.HttpHandlerRegisterTmpl).Execute(writer, pack)
}

// registerMethod will register method into service interface
func (r *RouterSprinkle) registerMethod(ctx *common.GenContext) error {
	method := ctx.GetNowMethod()
	pack := &packs.RouterProtocolMethodTmplPack{
		MethodName:   method.GoName,
		RequestType:  common.MapperGoTypeNameFromMessage(ctx, method.Input.Desc),
		ResponseType: common.MapperGoTypeNameFromMessage(ctx, method.Output.Desc),
	}
	writer := ctx.GetNowServiceContainer().BorrowMethodsWriter()
	return config.GetTemplate(config.HttpProtocolMethodTmpl).Execute(writer, pack)
}

func (r *RouterSprinkle) RegisterTier(ctx *common.GenContext) error {
	method := ctx.GetNowMethod()
	pack := &packs.RouterProtocolTierTmplPack{
		MethodName:   method.GoName,
		RequestType:  common.MapperGoTypeNameFromMessage(ctx, method.Input.Desc),
		ResponseType: common.MapperGoTypeNameFromMessage(ctx, method.Output.Desc),
	}
	writer := ctx.GetNowServiceContainer().BorrowTiersWriter()
	return config.GetTemplate(config.HttpProtocolTierTmpl).Execute(writer, pack)
}
