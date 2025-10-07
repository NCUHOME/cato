package models

import (
	"fmt"
	"strings"

	"github.com/ncuhome/cato/generated"
)

type Import struct {
	Alias      string
	ImportPath string
}

func (cip *Import) GetPath() string {
	value := fmt.Sprintf("\"%s\"", cip.ImportPath)
	if cip.Alias != "" {
		value = fmt.Sprintf("%s \"%s\"", cip.Alias, cip.ImportPath)
	}
	return value
}

func (cip *Import) Init(importPath string) *Import {
	cip.ImportPath = importPath
	pattern := strings.Split(importPath, "/")
	length := len(pattern)
	if length <= 2 {
		return cip
	}
	cip.Alias = strings.Join([]string{pattern[length-2], pattern[length-1]}, "")
	return cip
}

type Kv struct {
	Key   string
	Value string
}

type Tag struct {
	KV     *Kv
	Mapper func(s string) string
}

func (t *Tag) GetTagValue(by string) string {
	if t.KV == nil {
		return by
	}
	if t.KV.Value != "" {
		return t.KV.Value
	}
	return t.Mapper(by)
}

type Field struct {
	Name   string
	GoType string
}

type Key struct {
	// this represents from field and type
	KeyName string
	KeyType generated.DBKeyType
	Fields  []*Field
}

type Col struct {
	ColName string
	Field   *Field
}
