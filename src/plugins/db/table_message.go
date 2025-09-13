package db

import (
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins"
	"text/template"
)

type TableMessageEx struct {
	message *plugins.ModelsPlugger

	value *generated.TableOption

	tmpl *template.Template
}

type TableMessageExTmlPack struct {
	MessageName string
	TableName   string
	Comment     string
}

func (t *TableMessageEx) GetTmplFileName() string {
	return "table_name.tmpl"
}

func (t *TableMessageEx) Init(tmpl *template.Template) {
	t.tmpl = tmpl
}

func (t *TableMessageEx) LoadPlugger(message *plugins.ModelsPlugger) {
	t.message = message
}

func (t *TableMessageEx) AsTmplPack() interface{} {
	nameOpt := t.value.GetNameOption()
	if nameOpt == nil || nameOpt.GetLazyName() || nameOpt.GetSimpleName() == "" {
		return nil
	}
	return &TableMessageExTmlPack{
		MessageName: t.message.GetMessageName(),
		TableName:   nameOpt.GetSimpleName(),
		Comment:     t.value.GetComment(),
	}
}

func (t *TableMessageEx) Register() error {
	if t.value == nil || t.message == nil {
		return nil
	}
	pack := &TableMessageExTmlPack{
		MessageName: t.message.GetMessageName(),
		Comment:     t.value.GetNameOption().GetSimpleName(),
	}
	// check if the table name is simple
	if t.value.NameOption.GetSimpleName() != "" {
		pack.TableName = t.value.NameOption.GetSimpleName()
		return t.tmpl.Execute(t.message.BorrowMethodsWriter(), pack)
	}
	// empty table name will impl in extra file
	return t.tmpl.Execute(t.message.BorrowExtraWriter(), pack)
}
