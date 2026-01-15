package http

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
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
	if !ctx.NeedDoc() {
		return nil
	}
	name := string(ctx.GetNowService().Desc.Name())
	tag := &models.SwaggerTag{Name: name, Description: name}
	sv := &models.SwaggerApis{
		Tags:       []*models.SwaggerTag{tag},
		BasePath:   a.value.GroupPrefix,
		Containers: make(map[string]map[string]*models.SwaggerRoute),
	}
	ctx.AddDocService(name, sv)
	return nil
}
