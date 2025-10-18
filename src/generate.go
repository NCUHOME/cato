package src

import (
	"errors"
	"go/format"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/flags"
	"github.com/ncuhome/cato/src/plugins/models"
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
	outdir := flags.ParseProtoOptFlag(g.req.GetParameter())[flags.FlagExtOutDir]
	if err != nil {
		log.Fatalln(err)
	}
	respFiles := make([]*pluginpb.CodeGeneratorResponse_File, 0)
	for _, file := range genOption.Files {
		fc := ware.NewFileWare(file)
		ctx := fc.RegisterContext(new(common.GenContext))
		_, err = fc.Active(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		files := make([]*models.GenerateFileDesc, 0)
		allMessages := g.loadAllMessages(file.Messages)
		for _, message := range allMessages {
			// init message scope tray
			mc := ware.NewMessageWare(message)
			mctx := mc.RegisterContext(ctx)
			ok, err := mc.Active(mctx)
			if err != nil || !ok {
				log.Fatalf("[-] cato could not activate message %s: %v\n", mctx.GetNowMessageTypeName(), err)
			}
			err = mc.Complete(mctx)
			if err != nil {
				log.Fatalf("[-] cato could not complete message %s: %v\n", mctx.GetNowMessageTypeName(), err)
			}
			genFiles, err := mc.GetFiles(mctx)
			if err != nil {
				log.Fatalln(err)
			}
			files = append(files, genFiles...)
		}
		// todo as complete func
		for _, gf := range files {
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
	return respFiles
}

func (g *CatoGenerator) outputContent(filename, content string) *pluginpb.CodeGeneratorResponse_File {
	c, err := format.Source([]byte(content))
	if err != nil {
		log.Fatalf("[-] cato formatted %s file content %s error\n", filename, content)
	}
	formatted := string(c)
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &filename,
		Content: &formatted,
	}
}

// loadAllMessages will load all messages include nested from parent messages
func (g *CatoGenerator) loadAllMessages(parents []*protogen.Message) []*protogen.Message {
	if len(parents) == 0 {
		return make([]*protogen.Message, 0)
	}
	results := make([]*protogen.Message, 0)
	// first load self
	results = append(results, parents...)
	for _, message := range parents {
		// load message children
		results = append(results, g.loadAllMessages(message.Messages)...)
	}
	return results
}
