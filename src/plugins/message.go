package plugins

import (
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/db"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type MessagesPlugger struct {
	message *protogen.Message
	root    *protogen.File

	fields  map[string]*FieldsPlugger
	imports []*strings.Builder
	methods []*strings.Builder
	extra   []*strings.Builder
	tmpl    *template.Template
}

type ModelsPluggerPack struct {
	PackageName string
	Imports     []string
	ModelName   string
	Fields      []string
	Methods     []string
}

func (mp *MessagesPlugger) LoadContext(message *protogen.Message, file *protogen.File) {
	mp.message = message
	mp.root = file
	mp.fields = make(map[string]*FieldsPlugger)
	mp.imports = make([]*strings.Builder, 0)
	mp.methods = make([]*strings.Builder, 0)
	mp.extra = make([]*strings.Builder, 0)
}

func (mp *MessagesPlugger) findField(name string) (*protogen.Field, bool) {
	for _, field := range mp.message.Fields {
		if string(field.Desc.Name()) == name {
			return field, true
		}
	}
	return nil, false
}

func (mp *MessagesPlugger) BorrowFieldsWriter(name string) (io.Writer, bool) {
	_, ok := mp.fields[name]
	if !ok {
		fieldDesc, ok := mp.findField(name)
		if !ok {
			return nil, false
		}
		mp.fields[name] = &FieldsPlugger{fieldDesc, make([]*strings.Builder, 0)}
	}
	return mp.fields[name].BorrowWriter(), true
}

func (mp *MessagesPlugger) BorrowMethodsWriter() io.Writer {
	mp.methods = append(mp.methods, new(strings.Builder))
	return mp.methods[len(mp.methods)-1]
}

func (mp *MessagesPlugger) BorrowImportsWriter() io.Writer {
	mp.imports = append(mp.imports, new(strings.Builder))
	return mp.imports[len(mp.imports)-1]
}

func (mp *MessagesPlugger) BorrowExtraWriter() io.Writer {
	mp.extra = append(mp.extra, new(strings.Builder))
	return mp.extra[len(mp.extra)-1]
}

func (mp *MessagesPlugger) GetExtensionType() protoreflect.ExtensionType {
	return generated.E_DbOpt
}

func (mp *MessagesPlugger) GetMessageName() string {
	return mp.message.GoIdent.GoName
}

func (mp *MessagesPlugger) AsTmplPack() *ModelsPluggerPack {
	imports := make([]string, len(mp.imports))
	for i, imp := range mp.imports {
		imports[i] = imp.String()
	}
	fields := make([]string, len(mp.fields))
	fieldsIndex := 0
	for _, field := range mp.fields {
		value := field.GetContent()
		fields[fieldsIndex] = value
		fieldsIndex++
	}
	methods := make([]string, len(mp.methods))
	for index, method := range mp.methods {
		methods[index] = method.String()
	}
	return &ModelsPluggerPack{
		PackageName: utils.GetGoPackageName(mp.root.GoImportPath),
		Imports:     imports,
		ModelName:   mp.message.GoIdent.GoName,
		Fields:      fields,
		Methods:     methods,
	}
}

func (mp *MessagesPlugger) ParseFields() {
	for _, field := range mp.message.Fields {
		fieldPlugger := &FieldsPlugger{
			fieldValue: field,
			fields:     make([]*strings.Builder, 0),
		}
		mp.fields[field.GoName] = fieldPlugger
	}
}

func (mp *MessagesPlugger) GetTemplateName() string {
	return "models.tmpl"
}

func (mp *MessagesPlugger) Init(template *template.Template) {
	mp.tmpl = template
}

func (mp *MessagesPlugger) GenerateFile() string {
	return fmt.Sprintf("%s.cato.go", strings.ToLower(mp.message.GoIdent.GoName))
}

func (mp *MessagesPlugger) GenerateContent() string {
	sw := new(strings.Builder)
	err := mp.tmpl.Execute(sw, mp.AsTmplPack())
	if err != nil {
		log.Fatalln("[-] models plugger exec tmpl error, ", err)
	}
	return sw.String()
}

func (mp *MessagesPlugger) Active() (bool, error) {
	mp.ParseFields()
	for name, field := range mp.fields {
		_, err := field.Active()
		if err != nil {
			log.Fatalln("[-] cato generte field error, ", name, err)
		}
	}

	tableExt := new(db.TableMessageEx)
	descriptor := protodesc.ToDescriptorProto(mp.message.Desc)
	if !proto.HasExtension(descriptor.Options, tableExt.FromExtType()) {
		return false, nil
	}
	value := proto.GetExtension(descriptor.Options, tableExt.FromExtType())
	tableExt.Init(value)
	tableExt.From(mp.message)
	tableExt.SetWriter(mp.BorrowMethodsWriter(), mp.BorrowExtraWriter())
	err := tableExt.Register()
	if err != nil {
		return false, err
	}
	return true, nil
}
