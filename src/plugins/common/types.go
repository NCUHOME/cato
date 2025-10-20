package common

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func MapperGoTypeNameFromField(ctx *GenContext, field protoreflect.FieldDescriptor) string {
	switch field.Kind() {
	case protoreflect.MessageKind:
		if field.IsMap() {
			keyType := MapperGoTypeNameFromField(ctx, field.MapKey())
			valueType := MapperGoTypeNameFromField(ctx, field.MapValue())
			return fmt.Sprintf("map[%s][%s]", keyType, valueType)
		}
		typeName := MapperGoTypeNameFromMessage(ctx, field.Message().(protoreflect.MessageDescriptor))
		if field.IsList() {
			return fmt.Sprintf("[]%s", typeName)
		}
		return fmt.Sprintf("%s", typeName)
	case protoreflect.EnumKind:
		// todo can define if enum map to string or int
		return "int32"
	default:
		return field.Kind().String()
	}
}

func MapperGoTypeNameFromMessage(ctx *GenContext, messageDesc protoreflect.MessageDescriptor) string {
	if messageDesc.IsMapEntry() {
		keyType := MapperGoTypeNameFromField(ctx, messageDesc.Fields().Get(0))
		valueType := MapperGoTypeNameFromField(ctx, messageDesc.Fields().Get(1))
		return fmt.Sprintf("map[%s]%s", keyType, valueType)
	}
	typeName := string(messageDesc.Name())
	alias := ctx.GetImportPathAlias(messageDesc)
	if alias != "" {
		typeName = fmt.Sprintf("%s.%s", alias, typeName)
	}
	return fmt.Sprintf("*%s", typeName)
}

func UnwrapPointType(typeRaw string) string {
	firstRune := typeRaw[0]
	if firstRune != '[' && firstRune != '*' {
		return typeRaw
	}
	return strings.ReplaceAll(typeRaw, "*", "")
}
