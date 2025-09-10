package src

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
)

type CatoGenerator interface {
	Generate() *pluginpb.CodeGeneratorResponse
	Accept(element protoreflect.ExtensionType)
}
