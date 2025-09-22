package common

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func MapperGoTypeName(ctx *GenContext, field protoreflect.FieldDescriptor) string {
	switch field.Kind() {
	case protoreflect.MessageKind:
		if field.IsMap() {
			keyType := MapperGoTypeName(ctx, field.MapKey())
			valueType := MapperGoTypeName(ctx, field.MapValue())
			return fmt.Sprintf("map[%s][%s]", keyType, valueType)
		}
		messageDesc := field.Message().(protoreflect.MessageDescriptor)
		typeName := string(messageDesc.Name())
		alias := ctx.GetImportPathAlias(messageDesc)
		if alias != "" {
			typeName = fmt.Sprintf("%s.%s", alias, typeName)
		}
		if field.IsList() {
			return fmt.Sprintf("*%s[]", typeName)
		}
		return fmt.Sprintf("*%s", typeName)
	case protoreflect.EnumKind:
		// todo can define if enum map to string or int
		return "int32"
	default:
		return field.Kind().String()
	}
}
