package common

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Butter interface {
	FromExtType() protoreflect.ExtensionType
	WorkOn(desc protoreflect.Descriptor) bool
	GetTmplFileName() string
	Init(value interface{})

	AsTmplPack(ctx *GenContext) interface{}
	Register(ctx *GenContext) error
}
