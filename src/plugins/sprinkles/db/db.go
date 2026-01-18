package db

import (
	"github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
)

func init() {
	sprinkles.Register(func() sprinkles.Sprinkle {
		return new(DataBaseOptSprinkle)
	})
}

type DataBaseOptSprinkle struct {
	value *generated.DbOption
}

func (d *DataBaseOptSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_DbOpt
}

func (d *DataBaseOptSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.MessageDescriptor)
	return ok
}

func (d *DataBaseOptSprinkle) Init(value interface{}) {
	v, ok := value.(*generated.DbOption)
	if !ok {
		return
	}
	d.value = v
}

func (d *DataBaseOptSprinkle) loadSqlPlaceholder() squirrel.PlaceholderFormat {
	switch d.value.GetDbType() {
	case generated.DbTypeEnum_CATO_DB_TYPE_POSTGRESQL:
		return squirrel.Dollar
	case generated.DbTypeEnum_CATO_DB_TYPE_SQLLITE:
		return squirrel.AtP
	default:
		return squirrel.Question
	}
}

func (d *DataBaseOptSprinkle) Register(ctx *common.GenContext) error {
	var placeHolder string
	switch d.value.GetDbType() {
	case generated.DbTypeEnum_CATO_DB_TYPE_POSTGRESQL:
		placeHolder = "squirrel.Dollar"
	case generated.DbTypeEnum_CATO_DB_TYPE_SQLLITE:
		placeHolder = "squirrel.AtP"
	default:
		placeHolder = "squirrel.Question"
	}
	// register squirrel package
	mc := ctx.GetNowMessageContainer()
	sqImport := "\"github.com/Masterminds/squirrel\""
	_, err := mc.BorrowImportWriter().Write([]byte(sqImport))
	if err != nil {
		return err
	}
	pack := &packs.TableStatPack{
		PlaceHolder:     placeHolder,
		MessageTypeName: ctx.GetNowMessageTypeName(),
	}
	return config.GetTemplate(config.TableStatTmpl).Execute(mc.BorrowMethodsWriter(), pack)
}
