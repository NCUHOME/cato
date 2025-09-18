package src

import (
	"log"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/config"
	db2 "github.com/ncuhome/cato/src/plugins"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type DbGenerator struct {
	req  *pluginpb.CodeGeneratorRequest
	resp *pluginpb.CodeGeneratorResponse
}

func NewDBGenerator(req *pluginpb.CodeGeneratorRequest) *DbGenerator {
	return &DbGenerator{req: req}
}

func (g *DbGenerator) Generate(resp *pluginpb.CodeGeneratorResponse) *pluginpb.CodeGeneratorResponse {
	genOption, err := protogen.Options{}.New(g.req)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range genOption.Files {
		goPackageName := utils.GetGoPackageName(file.GoImportPath)
		if goPackageName == "" {
			continue
		}
		for _, message := range file.Messages {
			// test for single plugger
			mp := new(db2.MessagesPlugger)
			mp.Init(config.GetTemplate(mp.GetTemplateName()))
			mp.LoadContext(message, file)
			ok, err := mp.Active()
			if err != nil {
				log.Fatalln(err)
			}
			if !ok {
				continue
			}
			fileName := filepath.Join(utils.GetGoFilePath(file.GoImportPath), mp.GenerateFile())
			content := mp.GenerateContent()
			resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    &fileName,
				Content: &content,
			})
		}
	}
	return nil
}
