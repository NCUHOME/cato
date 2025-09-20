package utils

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GetGoImportName(importPath protogen.GoImportPath) string {
	patterns := strings.Split(GetGoFilePath(importPath), "/")
	if len(patterns) == 0 || patterns[0] == "." {
		return ""
	}
	return patterns[len(patterns)-1]
}

func GetGoFilePath(importPath protogen.GoImportPath) string {
	return strings.Trim(importPath.String(), "\"")
}

func MapperGoTypeName(field protoreflect.FieldDescriptor) string {
	switch field.Kind() {
	case protoreflect.MessageKind:
		if field.IsMap() {
			keyType := MapperGoTypeName(field.MapKey())
			valueType := MapperGoTypeName(field.MapValue())
			return fmt.Sprintf("map[%s][%s]", keyType, valueType)
		}
		messageDesc := field.Message().(protoreflect.MessageDescriptor)
		if field.IsList() {
			return fmt.Sprintf("*%s[]", messageDesc.Name())
		}
		return fmt.Sprintf("*%s", messageDesc.Name())
	case protoreflect.EnumKind:
		// todo can define if enum map to string or int
		return "int32"
	default:
		return field.Kind().String()
	}
}
