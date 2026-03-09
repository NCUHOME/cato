package http

import (
	"errors"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
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
		r.registerRouter(ctx),
		r.registerApiMethod(ctx),
		r.registerHttpMethod(ctx),
		r.registerSwagger(ctx),
	)
}

func (r *RouterSprinkle) registerRouter(ctx *common.GenContext) error {
	method := ctx.GetNowMethod()
	pack := &packs.HttpRouterTmplPack{
		Method: r.value.Method,
		Uri:    r.value.Router,
		Func:   method.GoName,
	}
	writer := ctx.GetNowServiceContainer().BorrowRouterssWriter()
	return config.GetTemplate(config.HandlersRouterTmpl).Execute(writer, pack)
}

func (r *RouterSprinkle) registerApiMethod(ctx *common.GenContext) error {
	method := ctx.GetNowMethod()
	pack := &packs.RouterProtocolMethodTmplPack{
		MethodName:   method.GoName,
		RequestType:  common.MapperGoTypeNameFromMessage(ctx, method.Input.Desc).GoType(),
		ResponseType: common.MapperGoTypeNameFromMessage(ctx, method.Output.Desc).GoType(),
	}
	sc := ctx.GetNowServiceContainer()
	return config.GetTemplate(config.ApiMethodTmpl).Execute(sc.BorrowApisWriter(), pack)

}

func (r *RouterSprinkle) registerHttpMethod(ctx *common.GenContext) error {
	method := ctx.GetNowMethod()
	sc := ctx.GetNowServiceContainer()
	_, err := sc.BorrowMethodsWriter().Write([]byte(method.GoName))
	return err
}

func (r *RouterSprinkle) registerSwagger(ctx *common.GenContext) error {
	if !ctx.NeedDoc() {
		return nil
	}
	method := r.value.GetMethod()
	comment := ctx.GetNowMethod().Comments.Leading.String()
	route := &models.SwaggerRoute{
		Description: comment,
		Produces:    []string{"application/json"},
		Requests: &models.SwaggerMessageRef{
			FullName: string(ctx.GetNowMethodContainer().Request().Desc.FullName()),
		},
		Responses: &models.SwaggerMessageRef{
			FullName: string(ctx.GetNowMethodContainer().Response().Desc.FullName()),
		},
	}
	switch method {
	case "GET":
		route.Requests.Location = "query"
	default:
		route.Requests.Location = "body"
		route.Consumes = []string{"application/json"}
	}
	ctx.AddDocApi(r.value.Router, method, route)
	return nil
}
