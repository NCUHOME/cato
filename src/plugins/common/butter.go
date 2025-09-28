package common

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Butter interface {
	FromExtType() protoreflect.ExtensionType
	WorkOn(desc protoreflect.Descriptor) bool
	Init(value interface{})
	Register(ctx *GenContext) error
}
