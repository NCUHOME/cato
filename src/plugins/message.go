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

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/db"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type MessagesPlugger struct {
	message *protogen.Message
	context *common.GenContext

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

func (mp *MessagesPlugger) LoadContext(gc *common.GenContext) {
	mp.context = gc
	mp.message = gc.GetNowMessage()
	mp.fields = make(map[string]*FieldsPlugger)
	mp.imports = make([]*strings.Builder, 0)
	mp.methods = make([]*strings.Builder, 0)
	mp.extra = make([]*strings.Builder, 0)
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
		PackageName: utils.GetGoImportName(mp.context.GetNowFile().GoImportPath),
		Imports:     imports,
		ModelName:   mp.message.GoIdent.GoName,
		Fields:      fields,
		Methods:     methods,
	}
}

func (mp *MessagesPlugger) ParseFields() {
	for _, field := range mp.message.Fields {
		fp := new(FieldsPlugger)
		fp.LoadContext(mp.context.WithField(field))
		mp.fields[field.GoName] = fp
	}
}

func (mp *MessagesPlugger) GetTemplateName() string {
	return "models.tmpl"
}

func (mp *MessagesPlugger) Init(template *template.Template) {
	mp.tmpl = template
}

func (mp *MessagesPlugger) GenerateFile() string {
	return fmt.Sprintf("%s.cato.go", mp.outputFileName())
}

func (mp *MessagesPlugger) outputFileName() string {
	patterns := utils.SplitCamelWords(mp.message.GoIdent.GoName)
	mapper := utils.GetStringMapper(generated.FieldMapper_CATO_FIELD_MAPPER_SNAKE_CASE)
	return mapper(patterns)
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
	descriptor := protodesc.ToDescriptorProto(mp.message.Desc)
	butter := db.ChooseButter(mp.message.Desc)
	for index := range butter {
		if !proto.HasExtension(descriptor.Options, butter[index].FromExtType()) {
			continue
		}
		value := proto.GetExtension(descriptor.Options, butter[index].FromExtType())
		butter[index].Init(mp.context, value)
		butter[index].SetWriter(mp.BorrowMethodsWriter(), mp.BorrowExtraWriter())
		err := butter[index].Register()
		if err != nil {
			return false, err
		}
	}
	for name, field := range mp.fields {
		_, err := field.Active()
		if err != nil {
			log.Fatalln("[-] cato generate field error, ", name, err)
		}
	}
	return true, nil
}
