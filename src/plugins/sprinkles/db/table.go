package db

import (
	"log"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/sprinkles"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
)

func init() {
	sprinkles.Register(func() sprinkles.Sprinkle {
		return new(TableBasicSprinkle)
	})
}

type TableBasicSprinkle struct {
	value *generated.TableOption
}

func (t *TableBasicSprinkle) extendTmplName() string {
	return "table_extend.tmpl"
}

func (t *TableBasicSprinkle) Init(value interface{}) {
	exValue, ok := value.(*generated.TableOption)
	if !ok {
		log.Fatalf("[-] can not convert %#v to TableOption", value)
	}

	t.value = exValue
}

func (t *TableBasicSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_TableOpt
}

func (t *TableBasicSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.MessageDescriptor)
	return ok
}

func (t *TableBasicSprinkle) Register(ctx *common.GenContext) error {
	if t.value == nil {
		return nil
	}
	pack := &packs.TableBasicTmplPack{
		MessageTypeName: ctx.GetNowMessage().GoIdent.GoName,
		Comment:         t.value.GetComment(),
	}
	// set extension
	mc := ctx.GetNowMessageContainer()
	_, err := mc.BorrowFieldWriter().Write([]byte("ext *extension"))
	// data object auto set ext package
	mc.SetNeedExtraFile(true)
	if err != nil {
		return err
	}
	// need set ext file
	tmpl := config.GetTemplate(config.TableNameTmpl)
	// check if the table name is simple
	if t.value.NameOption.GetSimpleName() != "" {
		pack.TableName = t.value.NameOption.GetSimpleName()
		return tmpl.Execute(mc.BorrowMethodsWriter(), pack)
	}
	// table name will impl in an extra file
	return tmpl.Execute(mc.BorrowExtraWriter(), pack)
}
