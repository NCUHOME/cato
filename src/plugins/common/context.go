package common

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/generated"
)

func GetCatoPackageFromFile(filedesc protoreflect.FileDescriptor) (string, bool) {
	if !proto.HasExtension(filedesc.Options(), generated.E_CatoOpt) {
		return "", false
	}
	catoOptions := proto.GetExtension(filedesc.Options(), generated.E_CatoOpt).(*generated.CatoOptions)
	return catoOptions.CatoPackage, catoOptions.CatoPackage != ""
}

type CatoImportPath struct {
	Alias      string
	ImportPath string
}

func (cip *CatoImportPath) Init(importPath string) *CatoImportPath {
	cip.ImportPath = importPath
	pattern := strings.Split(importPath, "/")
	length := len(pattern)
	if length <= 2 {
		return cip
	}
	cip.Alias = strings.Join([]string{pattern[length-2], pattern[length-1]}, "")
	return cip
}

type GenContext struct {
	file        *protogen.File
	message     *protogen.Message
	field       *protogen.Field
	catoPackage string
	namespaces  map[string]*CatoImportPath
}

func (gc *GenContext) WithFile(file *protogen.File) *GenContext {
	ctx := &GenContext{
		file:       file,
		namespaces: make(map[string]*CatoImportPath),
	}
	desc := file.Desc
	catoPackage, ok := GetCatoPackageFromFile(desc)
	if !ok {
		return ctx
	}
	ctx.catoPackage = catoPackage
	for index := 0; index < desc.Imports().Len(); index++ {
		importFile := desc.Imports().Get(index)
		importPackage := string(importFile.FileDescriptor.Package())
		importCatoPath, ok := GetCatoPackageFromFile(importFile.FileDescriptor)
		if !ok {
			continue
		}
		ctx.namespaces[importPackage] = new(CatoImportPath).Init(importCatoPath)
	}
	return ctx
}

func (gc *GenContext) CatoPackage() string {
	return gc.catoPackage
}

func (gc *GenContext) EmptyNamespace() bool {
	return len(gc.namespaces) == 0
}

func (gc *GenContext) GetImports() []string {
	imports, index := make([]string, len(gc.namespaces)), 0
	for _, v := range gc.namespaces {
		value := fmt.Sprintf("\"%s\"", v.ImportPath)
		if v.Alias != "" {
			value = fmt.Sprintf("%s \"%s\"", v.Alias, v.ImportPath)
		}
		imports[index] = value
	}
	return imports
}

func (gc *GenContext) GetImportPathAlias(desc protoreflect.MessageDescriptor) string {
	parent := desc.FullName().Parent()
	current := gc.file.Desc.Package()
	if parent == current {
		return ""
	}
	v, ok := gc.namespaces[string(parent)]
	if !ok {
		return ""
	}
	return v.Alias
}

func (gc *GenContext) GetNowFile() *protogen.File {
	return gc.file
}

func (gc *GenContext) WithMessage(message *protogen.Message) *GenContext {
	return &GenContext{
		file:       gc.file,
		namespaces: gc.namespaces,
		message:    message,
	}
}

func (gc *GenContext) GetNowMessage() *protogen.Message {
	return gc.message
}

func (gc *GenContext) WithField(field *protogen.Field) *GenContext {
	return &GenContext{
		file:       gc.file,
		namespaces: gc.namespaces,
		message:    gc.message,
		field:      field,
	}
}

func (gc *GenContext) GetNowField() *protogen.Field {
	return gc.field
}
