package http

import (
	"bytes"
	"log"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
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
		return new(ParamSprinkle)
	})
}

type ParamSprinkle struct {
	option *generated.HttpParamOption
}

func (m *ParamSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_HttpParamOpt
}

func (m *ParamSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.MessageDescriptor)
	return ok
}

func (m *ParamSprinkle) Init(value interface{}) {
	m.option = value.(*generated.HttpParamOption)
}

func (m *ParamSprinkle) Register(ctx *common.GenContext) error {
	if m.option == nil {
		return nil
	}
	if !m.option.GetImpl() {
		return nil
	}
	mc := ctx.GetNowMessageContainer()
	// need impl param interface for message
	tmpl := config.GetTemplate(config.HttpParamImplTmpl)
	if tmpl == nil {
		return nil
	}
	pack := &packs.HttpParamPack{ModelType: ctx.GetNowMessageTypeName()}
	w := new(bytes.Buffer)
	err := tmpl.Execute(w, pack)
	if err != nil {
		return err
	}
	swaggerMessage := m.registerSwagger(ctx, ctx.GetNowMessage())
	ctx.AddDocMessage(swaggerMessage.Identify, swaggerMessage)
	_, err = mc.BorrowMethodsWriter().Write(w.Bytes())
	return err
}

func (m *ParamSprinkle) registerSwagger(ctx *common.GenContext, message *protogen.Message) *models.SwaggerMessage {
	identifyName := string(message.Desc.FullName())
	properties, required := make([]*models.SwaggerMessageField, 0), make([]string, 0)
	for _, field := range message.Fields {
		swaggerField := &models.SwaggerMessageField{
			Name:        string(field.Desc.Name()),
			Description: field.Comments.Leading.String(),
			Type:        field.Desc.Kind().String(),
			Identify:    string(field.Desc.FullName().Name()),
		}
		if field.Enum != nil && len(field.Enum.Values) > 0 {
			swaggerField.Enum = m.transEnums(field.Enum.Values)
			swaggerField.Type = "string"
		}
		if field.Desc.Kind() == protoreflect.MessageKind && field.Message.Desc.FullName() != message.Desc.FullName() {
			swaggerField.Refer = m.registerSwagger(ctx, field.Message)
			swaggerField.Type = "object"
		}
		if proto.HasExtension(field.Desc.Options(), generated.E_HttpPfOpt) {
			value := proto.GetExtension(field.Desc.Options(), generated.E_HttpPfOpt)
			opt, ok := value.(*generated.HttpParamFieldOption)
			if !ok {
				log.Printf("[-] param field %#v as pf error\n", value)
				properties = append(properties, swaggerField)
				continue
			}
			if opt.Skip {
				continue
			}
			swaggerField.Example = opt.Example
			swaggerField.Description += opt.ExtraDesc
			swaggerField.Format = opt.Format
			swaggerField.Required = opt.Must
			if swaggerField.Required {
				required = append(required, swaggerField.Name)
			}
		}
		properties = append(properties, swaggerField)
	}

	swaggerMessage := &models.SwaggerMessage{
		Name:       string(message.Desc.Name()),
		Type:       "object",
		Required:   required,
		Properties: properties,
		Identify:   identifyName,
	}
	ctx.AddDocMessage(identifyName, swaggerMessage)
	return swaggerMessage
}

func (m *ParamSprinkle) transEnums(enums []*protogen.EnumValue) []string {
	ss := make([]string, len(enums))
	for i, enum := range enums {
		ss[i] = enum.GoIdent.GoName
	}
	return ss
}
