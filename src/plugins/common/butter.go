package common

import (
	"io"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type Butter interface {
	FromExtType() protoreflect.ExtensionType
	WorkOn(desc protoreflect.Descriptor) bool
	GetTmplFileName() string
	
	Init(gc *GenContext, value interface{})
	SetWriter(writers ...io.Writer)
	AsTmplPack() interface{}
	Register() error
}
