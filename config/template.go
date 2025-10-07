package config

import (
	"embed"
	"io"
	"log"
	"path/filepath"
	"text/template"
)

const (
	tmplPath = "templates"

	ColArrivalTmpl  = "col_arrival.tmpl"
	TableNameTmpl   = "table_name.tmpl"
	TableExtendTmpl = "table_extend.tmpl"
	TimeFormatTmpl  = "time_format.tmpl"
	JsonTransTmpl   = "json_trans.tmpl"
	TagTmpl         = "tag.tmpl"
	FieldTmpl       = "field.tmpl"
	ModelTmpl       = "model.tmpl"
	FileTmpl        = "file.tmpl"
)

var (
	//go:embed templates/*
	templatesFs embed.FS

	templates = make(map[string]*template.Template)
)

var templateNames = []string{

	// todo: need treat this name as const

	ColArrivalTmpl,
	TableNameTmpl,
	TableExtendTmpl,
	TimeFormatTmpl,
	JsonTransTmpl,
	FieldTmpl,
	TagTmpl,
	ModelTmpl,
	FileTmpl,
}

func init() {
	for _, name := range templateNames {
		path := filepath.Join(tmplPath, name)
		fs, err := templatesFs.Open(path)
		if err != nil {
			log.Fatalln(err)
		}
		data, err := io.ReadAll(fs)
		if err != nil {
			log.Fatalln(err)
		}
		tmpl, err := template.New(name).Parse(string(data))
		if err != nil {
			log.Fatalln(err)
		}
		templates[name] = tmpl
	}
}

func GetTemplate(name string) *template.Template {
	tmpl, ok := templates[name]
	if !ok {
		log.Fatalf("template %q not found", name)
	}
	return tmpl
}
