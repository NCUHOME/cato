package common

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/src/plugins/tray"
)

type GenContext struct {
	file          *protogen.File
	fileContainer *tray.FileTray

	message          *protogen.Message
	messageContainer *tray.MessageTray

	field          *protogen.Field
	fieldContainer *tray.FieldTray

	service          *protogen.Service
	serviceContainer *tray.ServiceTray

	method          *protogen.Method
	methodContainer *tray.MethodTray
}

func (gc *GenContext) WithFile(file *protogen.File, container *tray.FileTray) *GenContext {
	ctx := &GenContext{
		file:          file,
		fileContainer: container,
	}
	return ctx
}

func (gc *GenContext) GetCatoPackage() string {
	return gc.fileContainer.GetCatoPackage().GetPath()
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

func (gc *GenContext) GetNowFileContainer() *tray.FileTray {
	return gc.fileContainer
}

func (gc *GenContext) WithMessage(message *protogen.Message, container *tray.MessageTray) *GenContext {
	return &GenContext{
		file:             gc.file,
		fileContainer:    gc.fileContainer,
		message:          message,
		messageContainer: container,
	}
}

func (gc *GenContext) GetNowMessage() *protogen.Message {
	return gc.message
}

func (gc *GenContext) GetNowMessageContainer() *tray.MessageTray {
	return gc.messageContainer
}

func (gc *GenContext) GetNowMessageTypeName() string {
	return gc.GetNowMessage().GoIdent.GoName
}

func (gc *GenContext) WithField(field *protogen.Field, container *tray.FieldTray) *GenContext {
	return &GenContext{
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

func (gc *GenContext) GetNowFieldContainer() *tray.FieldTray {
	return gc.fieldContainer
}

func (gc *GenContext) WithService(service *protogen.Service, container *tray.ServiceTray) *GenContext {
	return &GenContext{
		file:             gc.file,
		fileContainer:    gc.fileContainer,
		service:          service,
		serviceContainer: container,
	}
}

func (gc *GenContext) GetNowService() *protogen.Service {
	return gc.service
}

func (gc *GenContext) GetNowServiceContainer() *tray.ServiceTray {
	return gc.serviceContainer
}

func (gc *GenContext) WithMethod(method *protogen.Method, container *tray.MethodTray) *GenContext {
	return &GenContext{
		file:             gc.file,
		fileContainer:    gc.fileContainer,
		service:          gc.service,
		serviceContainer: gc.serviceContainer,
		method:           method,
		methodContainer:  container,
	}
}

func (gc *GenContext) GetNowMethod() *protogen.Method {
	return gc.method
}

func (gc *GenContext) GetNowMethodContainer() *tray.MethodTray {
	return gc.methodContainer
}
