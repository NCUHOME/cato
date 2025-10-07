package db

import (
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type FieldKeysButter struct {
	values []*generated.DBKey
}

func (f *FieldKeysButter) FromExtType() protoreflect.ExtensionType {
	return generated.E_ColumnOpt
}

func (f *FieldKeysButter) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FieldDescriptor)
	return ok
}

func (f *FieldKeysButter) Init(value interface{}) {
	v, ok := value.(*generated.ColumnOption)
	if !ok {
		return
	}
	f.values = v.GetKeys()
}

func (f *FieldKeysButter) Register(ctx *common.GenContext) error {
	nowField := ctx.GetNowField()
	fieldName := nowField.GoName
	fieldType := common.MapperGoTypeName(ctx, nowField.Desc)
	field := &models.Field{
		Name:   fieldName,
		GoType: fieldType,
	}
	for index := range f.values {
		key := &models.Key{
			Fields:  []*models.Field{field},
			KeyType: f.values[index].KeyType,
		}
		ctx.GetScope().AddScopeKey(key)
	}
	return nil
}
