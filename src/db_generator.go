package src

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type dbGenerator struct {
	req *pluginpb.CodeGeneratorRequest
}

func NewDBGenerator(req *pluginpb.CodeGeneratorRequest) CatoGenerator {
	return &dbGenerator{req: req}
}

func (g *dbGenerator) Generate() *pluginpb.CodeGeneratorResponse {
	for _, file := range g.req.GetProtoFile() {
		if !g.shouldGen(file) {
			continue
		}
		for _, message := range file.GetMessageType() {
			for _, ext := range message.GetExtension() {
				ext.ProtoReflect().Interface()
			}
		}
	}
	return nil
}

func (g *dbGenerator) Accept(element protoreflect.ExtensionType) {

}

func (g *dbGenerator) shouldGen(file *descriptorpb.FileDescriptorProto) bool {
	return file != nil && file.GetOptions().GetGoPackage() != ""
}

func (g *dbGenerator) GenerateAll(req *pluginpb.CodeGeneratorRequest) *pluginpb.CodeGeneratorResponse {
	return g.Generate()
}
