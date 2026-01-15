package src

import (
	"errors"
	"go/format"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/imports"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/flags"
	"github.com/ncuhome/cato/src/plugins/ware"
)

type CatoGenerator struct {
	req    *pluginpb.CodeGeneratorRequest
	params map[string]string
}

func NewCatoGenerator(req *pluginpb.CodeGeneratorRequest) *CatoGenerator {
	g := &CatoGenerator{req: req}
	g.params = make(map[string]string)
	if g.req.GetParameter() != "" {
		g.params = flags.ParseProtoOptFlag(g.req.GetParameter())
	}
	return g
}

func (g *CatoGenerator) Generate() []*pluginpb.CodeGeneratorResponse_File {
	genOption, err := protogen.Options{}.New(g.req)
	outdir := g.params[flags.FlagExtOutDir]
	if err != nil {
		log.Fatalln(err)
	}
	respFiles := make([]*pluginpb.CodeGeneratorResponse_File, 0)
	root := new(common.GenContext).Init(g.params)
	for _, file := range genOption.Files {
		fc := ware.NewFileWare(file)
		ctx := fc.RegisterContext(root)
		_, err = ware.CommonWareActive(ctx, fc)
		if err != nil {
			log.Fatalln(err)
		}
		err = fc.Complete(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		fs, err := fc.GetExtraFiles(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		for _, gf := range fs {
			if gf.CheckExists {
				// check if file exists
				_, err := os.Stat(filepath.Join(outdir, gf.Name))
				if errors.Is(err, os.ErrNotExist) {
					respFiles = append(respFiles, g.outputContent(gf.Name, gf.Content))
				} else if err != nil {
					log.Fatalln(err)
				}
			} else {
				respFiles = append(respFiles, g.outputContent(gf.Name, gf.Content))
			}
		}
		err = fc.Complete(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if root.NeedDoc() {
		err = root.GenerateSwagger()
		if err != nil {
			log.Fatalln(err)
		}
	}
	return respFiles
}

func (g *CatoGenerator) outputContent(filename, content string) *pluginpb.CodeGeneratorResponse_File {
	options := &imports.Options{
		TabWidth:  4,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}
	c, err := imports.Process("", []byte(content), options)
	if err != nil {
		log.Fatalf("[-] cato import %s file content %s error\n", filename, content)
	}
	// 使用go/format进行标准格式化
	c, err = format.Source(c)
	if err != nil {
		log.Fatalf("[-] cato formatted %s file content %s error\n", filename, content)
	}
	formatted := string(c)
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &filename,
		Content: &formatted,
	}
}
