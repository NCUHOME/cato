package src

import (
	"io"
	"log"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/src/plugins"
	"github.com/ncuhome/cato/src/plugins/common"
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
	context := new(common.GenContext)
	for _, file := range genOption.Files {
		catoPackage, ok := GetCatoPackageFromFile(file.Desc)
		if !ok {
			continue
		}
		fc := context.WithFile(file)
		imports := GetImportPathFromFile(file)
		for _, message := range file.Messages {
			// test for single plugger
			mc := fc.WithMessage(message)
			mp := new(plugins.MessagesPlugger)
			mp.Init(config.GetTemplate(mp.GetTemplateName()))
			mp.LoadContext(mc)
			ok, err := mp.Active()
			if err != nil {
				log.Fatalln(err)
			}
			if !ok {
				continue
			}
			fileName := filepath.Join(catoPackage, mp.GenerateFile())
			for _, importName := range imports {
				_, err = io.WriteString(mp.BorrowImportsWriter(), importName)
				if err != nil {
					log.Fatalln(err)
				}
			}
			content := mp.GenerateContent()
			resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    &fileName,
				Content: &content,
			})
		}
	}
	return nil
}
