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

	CatoFileTmpl = "cato_file.tmpl"

	ColArrivalTmpl  = "cols_arrival.tmpl"
	TableNameTmpl   = "table_name.tmpl"
	TableExtendTmpl = "table_extend.tmpl"
	TimeFormatTmpl  = "time_format.tmpl"
	JsonTransTmpl   = "json_trans.tmpl"
	TagTmpl         = "tag.tmpl"
	FieldTmpl       = "field.tmpl"
	ModelTmpl       = "model.tmpl"
	FileTmpl        = "file.tmpl"
	TableColTmpl    = "table_col.tmpl"
	ColsGroupTmpl   = "cols_group.tmpl"
	TableStatTmpl   = "table_stat.tmpl"

	RdbTmpl       = "rdb.tmpl"
	RdbDeleteTmpl = "rdb_delete.tmpl"
	RdbFetchTmpl  = "rdb_fetch.tmpl"
	RdbUpdateTmpl = "rdb_update.tmpl"
	RdbInsertTmpl = "rdb_insert.tmpl"

	RepoTmpl       = "repo.tmpl"
	RepoExtTmpl    = "repo_ext.tmpl"
	RepoDeleteTmpl = "repo_delete.tmpl"
	RepoFetchTmpl  = "repo_fetch.tmpl"
	RepoUpdateTmpl = "repo_update.tmpl"
	RepoInsertTmpl = "repo_insert.tmpl"

	HttpHandlerTmpl         = "http_handler.tmpl"
	HttpHandlerRegisterTmpl = "http_handler_register.tmpl"
	HttpProtocolTmpl        = "http_protocol.tmpl"
	HttpProtocolMethodTmpl  = "http_protocol_method.tmpl"
	HttpProtocolTierTmpl    = "http_protocol_tier.tmpl"
	HttpParamImplTmpl       = "http_param_impl.tmpl"
)

var (
	//go:embed templates/*
	templatesFs embed.FS

	templates = make(map[string]*template.Template)
)

var templateNames = []string{
	CatoFileTmpl,

	ColArrivalTmpl,
	TableNameTmpl,
	TableExtendTmpl,
	TimeFormatTmpl,
	JsonTransTmpl,
	FieldTmpl,
	TagTmpl,
	ModelTmpl,
	FileTmpl,
	TableColTmpl,
	ColsGroupTmpl,
	TableStatTmpl,

	RepoTmpl,
	RepoExtTmpl,
	RepoDeleteTmpl,
	RepoFetchTmpl,
	RepoUpdateTmpl,
	RepoInsertTmpl,

	RdbTmpl,
	RdbDeleteTmpl,
	RdbFetchTmpl,
	RdbUpdateTmpl,
	RdbInsertTmpl,

	HttpHandlerTmpl,
	HttpHandlerRegisterTmpl,
	HttpProtocolTmpl,
	HttpProtocolMethodTmpl,
	HttpProtocolTierTmpl,
	HttpParamImplTmpl,
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
