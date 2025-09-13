package db

import (
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins"
	"github.com/ncuhome/cato/src/plugins/common"
	"text/template"
)

type ColumnFieldEx struct {
	field  *plugins.FieldsPlugger
	parent *plugins.ModelsPlugger

	value *generated.ColumnOption

	tmpl *template.Template
	tags map[string]*common.Kv
}

type ColumnFieldExTmlPack struct {
	Name   string
	GoType string
	Tags   []common.Kv
}

func (c *ColumnFieldEx) GetTmplFileName() string {
	return "column_field.tmpl"
}

func (c *ColumnFieldEx) Init(tmpl *template.Template) {
	c.tmpl = tmpl
}

func (c *ColumnFieldEx) LoadPlugger(field *plugins.FieldsPlugger, message *plugins.ModelsPlugger) {
	c.field = field
	c.parent = message
}

func (c *ColumnFieldEx) AsTmplPack() interface{} {
	tags, index := make([]common.Kv, len(c.tags)), 0
	for k, v := range c.tags {
		tags[index] = common.Kv{
			Key:   k,
			Value: v.Value,
		}
		index++
	}
	return &ColumnFieldExTmlPack{
		Name:   c.field.GetName(),
		GoType: c.field.GetGoType(),
		Tags:   tags,
	}
}

func (c *ColumnFieldEx) Register() error {
	return nil
}
