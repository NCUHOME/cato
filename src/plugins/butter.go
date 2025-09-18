package plugins

import (
	"io"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type Butter interface {
	GetTmplFileName() string
	Init(value interface{})
	SetWriter(writers ...io.Writer)
	AsTmplPack() interface{}
	Register() error
	FromExtType() protoreflect.ExtensionType
}
