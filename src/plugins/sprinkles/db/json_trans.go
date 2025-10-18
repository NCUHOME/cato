package db

import (
	"fmt"

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
		return new(JsonTransSprinkle)
	})
}

type JsonTransSprinkle struct {
	value *generated.ColumnOption
}

func (j *JsonTransSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_ColumnOpt
}

func (j *JsonTransSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FieldDescriptor)
	return ok
}

func (j *JsonTransSprinkle) Init(value interface{}) {
	data, ok := value.(*generated.ColumnOption)
	if !ok {
		return
	}
	j.value = data
}

func (j *JsonTransSprinkle) Register(ctx *common.GenContext) error {
	if j.value == nil || j.value.GetJsonTrans() == nil {
		return nil
	}
	transOpt := j.value.GetJsonTrans()
	// set json trans flag
	ctx.GetNowFieldContainer().SetJsonTrans(true)

	nowField := ctx.GetNowField()
	fieldType := common.MapperGoTypeNameFromField(ctx, nowField.Desc)

	if transOpt.LazyLoad {
		// need to register extra inner field into message-fields map
		extraField := &packs.FieldPack{
			Field: &models.Field{
				Name:   fmt.Sprintf("inner%s", nowField.GoName),
				GoType: fieldType,
			},
		}
		writer := ctx.GetNowMessageContainer().BorrowFieldWriter()
		err := config.GetTemplate(config.FieldTmpl).Execute(writer, extraField)
		if err != nil {
			return err
		}
	}
	_, err := ctx.GetNowMessageContainer().BorrowImportWriter().Write([]byte("\"encoding/json\""))
	if err != nil {
		return err
	}
	mWriter := ctx.GetNowMessageContainer().BorrowMethodsWriter()
	pack := &packs.JsonTransTmplPack{
		MessageTypeName: ctx.GetNowMessageTypeName(),
		FieldName:       nowField.GoName,
		FieldType:       fieldType,
		FieldTypeRaw:    common.UnwrapPointType(fieldType),
		LazyLoad:        j.value.JsonTrans.LazyLoad,
	}
	return config.GetTemplate(config.JsonTransTmpl).Execute(mWriter, pack)
}
