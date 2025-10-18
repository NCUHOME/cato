package db

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
)

func init() {
	sprinkles.Register(func() sprinkles.Sprinkle {
		return new(FieldKeysSprinkle)
	})
}

type FieldKeysSprinkle struct {
	values []*generated.DBKey
}

func (f *FieldKeysSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_ColumnOpt
}

func (f *FieldKeysSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FieldDescriptor)
	return ok
}

func (f *FieldKeysSprinkle) Init(value interface{}) {
	v, ok := value.(*generated.ColumnOption)
	if !ok {
		return
	}
	f.values = v.GetKeys()
}

func (f *FieldKeysSprinkle) Register(ctx *common.GenContext) error {
	nowField := ctx.GetNowField()
	fieldName := nowField.GoName
	fieldType := common.MapperGoTypeNameFromField(ctx, nowField.Desc)
	field := &models.Field{
		Name:   fieldName,
		GoType: fieldType,
	}
	mc := ctx.GetNowMessageContainer()
	for index := range f.values {
		key := &models.Key{
			Fields:  []*models.Field{field},
			KeyType: f.values[index].KeyType,
			KeyName: f.values[index].KeyName,
		}
		mc.AddScopeKey(key)
	}
	return nil
}
