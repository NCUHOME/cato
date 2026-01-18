package common

import (
	"encoding/json"
	"log"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/src/plugins/flags"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/tray"
)

type swaggerScope struct {
	messages map[string]*models.SwaggerMessage
	doc      *models.SwaggerDoc
}

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

	needDoc bool
	swagger *swaggerScope
}

func (gc *GenContext) Init(params map[string]string) *GenContext {
	swagegrPath := params[flags.SwaggerPath]
	if swagegrPath == "" {
		return gc
	}
	gc.needDoc = true
	gc.swagger = &swaggerScope{
		messages: make(map[string]*models.SwaggerMessage),
		doc: &models.SwaggerDoc{
			Tags: make([]*models.SwaggerTag, 0),
			Auth: make([]*models.SwaggerAuth, 0),
			Host: params[flags.ApiHost],
			Apis: make(map[string]*models.SwaggerApis),
			Loc:  swagegrPath,
		},
	}
	return gc
}

func (gc *GenContext) WithFile(file *protogen.File, container *tray.FileTray) *GenContext {
	ctx := &GenContext{
		file:          file,
		fileContainer: container,
		needDoc:       gc.needDoc,
		swagger:       gc.swagger,
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
		needDoc:          gc.needDoc,
		swagger:          gc.swagger,
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
		needDoc:          gc.needDoc,
		swagger:          gc.swagger,
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
		needDoc:          gc.needDoc,
		swagger:          gc.swagger,
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
		needDoc:          gc.needDoc,
		swagger:          gc.swagger,
	}
}

func (gc *GenContext) GetNowMethod() *protogen.Method {
	return gc.method
}

func (gc *GenContext) GetNowMethodContainer() *tray.MethodTray {
	return gc.methodContainer
}

func (gc *GenContext) NeedDoc() bool {
	return gc.needDoc
}

func (gc *GenContext) AddDocMessage(id string, param *models.SwaggerMessage) {
	if !gc.needDoc || gc.swagger == nil {
		return
	}
	if gc.swagger.messages == nil {
		gc.swagger.messages = make(map[string]*models.SwaggerMessage)
	}
	gc.swagger.messages[id] = param
}

func (gc *GenContext) AddDocService(name string, service *models.SwaggerApis) {
	if !gc.needDoc || gc.swagger == nil {
		return
	}
	gc.swagger.doc.AddService(name, service)
}

func (gc *GenContext) AddDocApi(path, method string, api *models.SwaggerRoute) {
	if !gc.needDoc || gc.swagger == nil {
		return
	}
	serviceName := string(gc.GetNowService().Desc.Name())
	sv, ok := gc.swagger.doc.Apis[serviceName]
	if !ok {
		log.Printf("WARNING: service %s not found in swagger apis", serviceName)
		return
	}
	sv.AddRoute(path, method, api)
}

func (gc *GenContext) GenerateSwagger() error {
	f, err := os.OpenFile(gc.swagger.doc.Loc, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	doc := gc.swagger.doc.AsDoc(gc.swagger.messages)
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}
