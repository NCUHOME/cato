package common

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type MapperType struct {
	TypeRaw  string
	IsSlice  bool
	IsStruct bool
}

func (mp MapperType) GoType() string {
	if mp.IsSlice {
		if mp.IsStruct {
			return fmt.Sprintf("[]*%s", mp.TypeRaw)
		}
		return fmt.Sprintf("[]%s", mp.TypeRaw)
	}
	if mp.IsStruct {
		return fmt.Sprintf("*%s", mp.TypeRaw)
	}
	return fmt.Sprintf("%s", mp.TypeRaw)
}

func (mp MapperType) RawType() string {
	return mp.TypeRaw
}

func MapperGoTypeNameFromField(ctx *GenContext, field protoreflect.FieldDescriptor) MapperType {
	switch field.Kind() {
	case protoreflect.MessageKind:
		if field.IsMap() {
			keyType := MapperGoTypeNameFromField(ctx, field.MapKey())
			valueType := MapperGoTypeNameFromField(ctx, field.MapValue())
			return MapperType{
				TypeRaw: fmt.Sprintf("map[%s][%s]", keyType.GoType(), valueType.GoType()),
			}
		}
		typeName := MapperGoTypeNameFromMessage(ctx, field.Message().(protoreflect.MessageDescriptor))
		return MapperType{TypeRaw: typeName.GoType(), IsSlice: field.IsList()}
	case protoreflect.EnumKind:
		// todo can define if enum map to string or int
		return MapperType{TypeRaw: "int32"}
	default:
		return MapperType{TypeRaw: field.Kind().String(), IsSlice: field.IsList()}
	}
}

func MapperGoTypeNameFromMessage(ctx *GenContext, messageDesc protoreflect.MessageDescriptor) MapperType {
	if messageDesc.IsMapEntry() {
		keyType := MapperGoTypeNameFromField(ctx, messageDesc.Fields().Get(0))
		valueType := MapperGoTypeNameFromField(ctx, messageDesc.Fields().Get(1))
		return MapperType{
			TypeRaw: fmt.Sprintf("map[%s][%s]", keyType.GoType(), valueType.GoType()),
		}
	}
	typeName := string(messageDesc.Name())
	alias := ctx.GetImportPathAlias(messageDesc)
	if alias != "" {
		typeName = fmt.Sprintf("%s.%s", alias, typeName)
	}
	return MapperType{TypeRaw: typeName, IsStruct: true}
}
