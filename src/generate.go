package src

import (
	"go/format"
	"log"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/src/plugins"
	"github.com/ncuhome/cato/src/plugins/common"
)

type CatoGenerator struct {
	req  *pluginpb.CodeGeneratorRequest
	resp *pluginpb.CodeGeneratorResponse
}

func NewCatoGenerator(req *pluginpb.CodeGeneratorRequest) *CatoGenerator {
	return &CatoGenerator{req: req}
}

func (g *CatoGenerator) Generate(resp *pluginpb.CodeGeneratorResponse) *pluginpb.CodeGeneratorResponse {
	genOption, err := protogen.Options{}.New(g.req)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range genOption.Files {
		fc := plugins.NewFileCheese(file)
		ctx := fc.RegisterContext(new(common.GenContext))
		for _, message := range file.Messages {
			mc := plugins.NewMessageCheese(message)
			mctx := mc.RegisterContext(ctx)
			mc.Init(config.GetTemplate(mc.GetTemplateName()))
			ok, err := mc.Active(mctx)
			if err != nil {
				log.Fatalln(err)
			}
			if !ok {
				continue
			}
			fileName := filepath.Join(mctx.CatoPackage(), mc.GenerateFile())
			content := mc.GenerateContent(mctx)
			formattedContent, err := format.Source([]byte(content))
			if err != nil {
				log.Fatalf("[-] cato formatted %s file content error\n", fileName)
			}
			ft := string(formattedContent)
			resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    &fileName,
				Content: &ft,
			})
		}
	}
	return nil
}
