package common

import (
	"github.com/ncuhome/cato/src/plugins/cheese"
	"github.com/ncuhome/cato/src/plugins/utils"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type GenContext struct {
	catoPackage   string
	file          *protogen.File
	fileContainer *cheese.FileCheese

	message          *protogen.Message
	messageContainer *cheese.MessageCheese

	field          *protogen.Field
	fieldContainer *cheese.FieldCheese
}

func (gc *GenContext) WithFile(file *protogen.File, container *cheese.FileCheese) *GenContext {
	ctx := &GenContext{
		file:          file,
		fileContainer: container,
	}
	catoPackage, ok := utils.GetCatoPackageFromFile(file.Desc)
	if !ok {
		return ctx
	}
	ctx.catoPackage = catoPackage
	return ctx
}

func (gc *GenContext) GetCatoPackage() string {
	return gc.catoPackage
}

func (gc *GenContext) GetFilePackage() string {
	return utils.GetGoImportName(gc.GetNowFile().GoImportPath)
}

func (gc *GenContext) CatoPackage() string {
	return gc.catoPackage
}

func (gc *GenContext) GetImportPathAlias(desc protoreflect.MessageDescriptor) string {
	parent := desc.FullName().Parent()
	current := gc.file.Desc.Package()
	if parent == current {
		return ""
	}
	return gc.fileContainer.GetImportPathAlias(string(parent))
}

func (gc *GenContext) GetNowFile() *protogen.File {
	return gc.file
}

func (gc *GenContext) GetNowFileContainer() *cheese.FileCheese {
	return gc.fileContainer
}

func (gc *GenContext) WithMessage(message *protogen.Message, container *cheese.MessageCheese) *GenContext {
	return &GenContext{
		catoPackage:      gc.catoPackage,
		file:             gc.file,
		fileContainer:    gc.fileContainer,
		message:          message,
		messageContainer: container,
	}
}

func (gc *GenContext) GetNowMessage() *protogen.Message {
	return gc.message
}

func (gc *GenContext) GetNowMessageContainer() *cheese.MessageCheese {
	return gc.messageContainer
}

func (gc *GenContext) GetNowMessageTypeName() string {
	return gc.GetNowMessage().GoIdent.GoName
}

func (gc *GenContext) WithField(field *protogen.Field, container *cheese.FieldCheese) *GenContext {
	return &GenContext{
		catoPackage:      gc.catoPackage,
		file:             gc.file,
		fileContainer:    gc.fileContainer,
		message:          gc.message,
		messageContainer: gc.messageContainer,
		field:            field,
		fieldContainer:   container,
	}
}

func (gc *GenContext) GetNowField() *protogen.Field {
	return gc.field
}

func (gc *GenContext) GetNowFieldContainer() *cheese.FieldCheese {
	return gc.fieldContainer
}
