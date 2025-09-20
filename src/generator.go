package src

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/generated"
)

type CatoGenerator interface {
	Generate() *pluginpb.CodeGeneratorResponse
	Accept(element protoreflect.ExtensionType)
}

func GetCatoPackageFromFile(filedesc protoreflect.FileDescriptor) (string, bool) {
	if !proto.HasExtension(filedesc.Options(), generated.E_CatoOpt) {
		return "", false
	}
	catoOptions := proto.GetExtension(filedesc.Options(), generated.E_CatoOpt).(*generated.CatoOptions)
	return catoOptions.CatoPackage, catoOptions.CatoPackage != ""
}

func GetImportPathFromFile(file *protogen.File) []string {
	imports := file.Desc.Imports()
	results := make([]string, 0)
	for index := 0; index < imports.Len(); index++ {
		catoPath, ok := GetCatoPackageFromFile(imports.Get(index))
		if !ok {
			continue
		}
		results = append(results, catoPath)
	}
	return results
}
