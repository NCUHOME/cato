package db

import (
	"log"
	"text/template"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
)

func init() {
	register(func() common.Butter {
		return new(TableMessageButter)
	})
}

type TableMessageButter struct {
	value *generated.TableOption
	tmpl  *template.Template
}

type TableMessageButterTmplPack struct {
	MessageTypeName string
	TableName       string
	Comment         string
}

func (t *TableMessageButter) GetTmplFileName() string {
	return "table_name.tmpl"
}

func (t *TableMessageButter) Init(value interface{}) {
	exValue, ok := value.(*generated.TableOption)
	if !ok {
		log.Fatalf("[-] can not convert %#v to TableOption", value)
	}

	t.value = exValue
	t.tmpl = config.GetTemplate(t.GetTmplFileName())
}

func (t *TableMessageButter) AsTmplPack(ctx *common.GenContext) interface{} {
	nameOpt := t.value.GetNameOption()
	if nameOpt == nil || nameOpt.GetLazyName() || nameOpt.GetSimpleName() == "" {
		return nil
	}
	return &TableMessageButterTmplPack{
		MessageTypeName: ctx.GetNowMessage().GoIdent.GoName,
		TableName:       nameOpt.GetSimpleName(),
		Comment:         t.value.GetComment(),
	}
}

func (t *TableMessageButter) FromExtType() protoreflect.ExtensionType {
	return generated.E_TableOpt
}

func (t *TableMessageButter) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.MessageDescriptor)
	return ok
}

func (t *TableMessageButter) Register(ctx *common.GenContext) error {
	if t.value == nil {
		return nil
	}
	pack := &TableMessageButterTmplPack{
		MessageTypeName: ctx.GetNowMessage().GoIdent.GoName,
		Comment:         t.value.GetComment(),
	}
	// check if the table name is simple
	writers := ctx.GetWriters()
	if t.value.NameOption.GetSimpleName() != "" {
		pack.TableName = t.value.NameOption.GetSimpleName()
		return t.tmpl.Execute(writers.MethodWriter(), pack)
	}
	// table name will impl in an extra file
	return t.tmpl.Execute(writers.ExtraWriter(), pack)
}
