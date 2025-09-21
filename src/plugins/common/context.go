package common

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type GenContext struct {
	file    *protogen.File
	message *protogen.Message
	field   *protogen.Field
}

func (gc *GenContext) WithFile(file *protogen.File) *GenContext {
	return &GenContext{
		file: file,
	}
}

func (gc *GenContext) GetNowFile() *protogen.File {
	return gc.file
}

func (gc *GenContext) WithMessage(message *protogen.Message) *GenContext {
	return &GenContext{
		file:    gc.file,
		message: message,
	}
}

func (gc *GenContext) GetNowMessage() *protogen.Message {
	return gc.message
}

func (gc *GenContext) WithField(field *protogen.Field) *GenContext {
	return &GenContext{
		file:    gc.file,
		message: gc.message,
		field:   field,
	}
}

func (gc *GenContext) GetNowField() *protogen.Field {
	return gc.field
}
