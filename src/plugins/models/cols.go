package models

import (
	"strings"

	"github.com/ncuhome/cato/generated"
)

type Field struct {
	Name   string
	GoType string
}

func (f *Field) AsParamName() string {
	return strings.ToLower(f.Name)
}

type Key struct {
	// this represents from field and type
	KeyName string
	KeyType generated.DBKeyType
	Fields  []*Field
}

func (k *Key) GetFieldNameCombine() string {
	fieldNames := make([]string, len(k.Fields))
	for i, f := range k.Fields {
		fieldNames[i] = f.Name
	}
	return strings.Join(fieldNames, "And")
}

type Col struct {
	ColName string
	*Field
}
