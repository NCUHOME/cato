package config

import (
	"embed"
	"io"
	"log"
	"path/filepath"
	"text/template"
)

const (
	TmplPath        = "templates"
	CommonFieldTmpl = "common_field_tmpl"
	CommonTagTmpl   = "common_tag_tmpl"
)

var (
	//go:embed templates/*
	templatesFs embed.FS

	templates = make(map[string]*template.Template)
)

var templateNames = []string{
	"column_field.tmpl",
	"models.tmpl",
	"table_name.tmpl",
}

func init() {
	templates[CommonFieldTmpl], _ = template.New(CommonFieldTmpl).
		Parse(`{{ .Name }} {{ .GoType }} `)
	templates[CommonTagTmpl], _ = template.New(CommonTagTmpl).
		Parse(`{{range .}}{{.Key}}:"{{.Value}}" {{end}}"`)

	for _, name := range templateNames {
		path := filepath.Join(TmplPath, name)
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
