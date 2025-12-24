package structs

import (
	"bytes"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
)

func init() {
	sprinkles.Register(func() sprinkles.Sprinkle {
		return new(MessageParamSprinkle)
	})
}

type MessageParamSprinkle struct {
	option *generated.HttpParamOption
}

func (m *MessageParamSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_HttpParamOpt
}

func (m *MessageParamSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.MessageDescriptor)
	return ok
}

func (m *MessageParamSprinkle) Init(value interface{}) {
	m.option = value.(*generated.HttpParamOption)
}

func (m *MessageParamSprinkle) Register(ctx *common.GenContext) error {
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
	_, err = mc.BorrowMethodsWriter().Write(w.Bytes())
	return err
}
