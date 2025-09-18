package db

import (
	"io"
	"log"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
)

type TableMessageEx struct {
	from *protogen.Message

	methodWriter io.Writer
	extraWriter  io.Writer

	value *generated.TableOption
	tmpl  *template.Template
}

type TableMessageExTmplPack struct {
	MessageTypeName string
	TableName       string
	Comment         string
}

func (t *TableMessageEx) GetTmplFileName() string {
	return "table_name.tmpl"
}

func (t *TableMessageEx) FromMessageName() string {
	return t.from.GoIdent.GoName
}

func (t *TableMessageEx) From(message *protogen.Message) {
	t.from = message
}

func (t *TableMessageEx) Init(value interface{}) {
	exValue, ok := value.(*generated.TableOption)
	if !ok {
		log.Fatalf("[-] can not convert %#v to TableOption", value)
	}

	t.value = exValue
	t.tmpl = config.GetTemplate(t.GetTmplFileName())
}

func (t *TableMessageEx) SetWriter(writers ...io.Writer) {
	if len(writers) < 2 {
		log.Fatalln("[-] TableMessageEx need at least two writer")
	}
	t.setWriter(writers[0], writers[1])
}

func (t *TableMessageEx) setWriter(methodWriter, extraWriter io.Writer) {
	t.methodWriter = methodWriter
	t.extraWriter = extraWriter
}

func (t *TableMessageEx) AsTmplPack() interface{} {
	nameOpt := t.value.GetNameOption()
	if nameOpt == nil || nameOpt.GetLazyName() || nameOpt.GetSimpleName() == "" {
		return nil
	}
	return &TableMessageExTmplPack{
		MessageTypeName: t.FromMessageName(),
		TableName:       nameOpt.GetSimpleName(),
		Comment:         t.value.GetComment(),
	}
}

func (t *TableMessageEx) FromExtType() protoreflect.ExtensionType {
	return generated.E_TableOpt
}

func (t *TableMessageEx) Register() error {
	if t.value == nil {
		return nil
	}
	pack := &TableMessageExTmplPack{
		MessageTypeName: t.FromMessageName(),
		Comment:         t.value.GetComment(),
	}
	// check if the table name is simple
	if t.value.NameOption.GetSimpleName() != "" {
		pack.TableName = t.value.NameOption.GetSimpleName()
		return t.tmpl.Execute(t.methodWriter, pack)
	}
	// empty table name will impl in an extra file
	return t.tmpl.Execute(t.extraWriter, pack)
}
