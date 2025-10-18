package db

import (
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
)

func init() {
	sprinkles.Register(func() sprinkles.Sprinkle {
		return new(TimeOptionSprinkle)
	})
}

type TimeOptionSprinkle struct {
	value      *generated.TimeOption
	timeFormat string
}

func (t *TimeOptionSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_ColumnOpt
}

func (t *TimeOptionSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FieldDescriptor)
	return ok
}

func (t *TimeOptionSprinkle) tmplName() string {
	return "time_format.tmpl"
}

func (t *TimeOptionSprinkle) Init(value interface{}) {
	colOpt, ok := value.(*generated.ColumnOption)
	if !ok {
		return
	}
	t.value = colOpt.TimeOption
	t.timeFormat = time.RFC3339
}

func (t *TimeOptionSprinkle) Register(ctx *common.GenContext) error {
	if t.value == nil {
		return nil
	}
	timeOpt := t.value
	if timeOpt.GetTimeFormat() != "" {
		t.timeFormat = timeOpt.GetTimeFormat()
	}
	mWriter := ctx.GetNowMessageContainer().BorrowMethodsWriter()
	pack := &packs.TimeOptionTmplPack{
		MessageTypeName: ctx.GetNowMessageTypeName(),
		FieldName:       ctx.GetNowField().GoName,
		Format:          t.timeFormat,
	}
	tmpl := config.GetTemplate(config.TimeFormatTmpl)
	err := tmpl.Execute(mWriter, pack)
	if err != nil {
		return err
	}
	_, err = ctx.GetNowMessageContainer().BorrowImportWriter().Write([]byte("\"time\""))
	return err
}
