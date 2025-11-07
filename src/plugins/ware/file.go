package ware

import (
	"fmt"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/tray"
	"github.com/ncuhome/cato/src/plugins/utils"
)

// FileWare works for parsing file options
// also a file is the root of a ware tree because of cato_package option
type FileWare struct {
	file       *protogen.File
	context    *common.GenContext
	extraFiles []*models.GenerateFileDesc
}

func (fw *FileWare) GetDescriptor() protoreflect.Descriptor {
	return fw.file.Desc
}

func NewFileWare(file *protogen.File) *FileWare {
	fc := new(FileWare)
	fc.file = file
	fc.extraFiles = make([]*models.GenerateFileDesc, 0)
	return fc
}

func (fw *FileWare) RegisterContext(gc *common.GenContext) *common.GenContext {
	f := tray.NewFileTray(fw.file)
	ctx := gc.WithFile(fw.file, f)
	return ctx
}

func (fw *FileWare) Active(ctx *common.GenContext) (bool, error) {
	return CommonWareActive(ctx, fw)
}

func (fw *FileWare) GetSubWares() []WorkWare {
	subs := make([]WorkWare, 0)
	// load messages descriptor in a file
	// for a file ware, message and service ware is the direct sub ware
	msgs := fw.loadAllMessages(fw.file.Messages)
	for _, msg := range msgs {
		mw := NewMessageWare(msg)
		subs = append(subs, mw)
	}
	// load services ware
	for _, svr := range fw.file.Services {
		subs = append(subs, NewServiceWare(svr))
	}
	return subs
}

func (fw *FileWare) Complete(ctx *common.GenContext) error {
	fc := ctx.GetNowFileContainer()
	catoPackage := fc.GetCatoPackage()
	// cato plugin use 'cato_package option' to determine if this proto file will be parsed by cato
	// if no cato package specified, cato plugin will not generate file content
	if catoPackage == nil || catoPackage.IsEmpty() {
		return nil
	}
	allImports := append(fc.GetImports(), fc.GetOutImports()...)
	allContent := fc.GetOutContent()
	pack := &packs.CatoFileTmplPack{
		Imports:     allImports,
		ContentList: allContent,
		PackageName: utils.GetGoPackageName(catoPackage.ImportPath),
	}
	sb := new(strings.Builder)
	err := config.GetTemplate(config.CatoFileTmpl).Execute(sb, pack)
	if err != nil {
		return err
	}
	filename := filepath.Join(catoPackage.ImportPath, fmt.Sprintf("%s.cato.go", fw.filename()))
	fw.extraFiles = append(fw.extraFiles, &models.GenerateFileDesc{
		Name:        filename,
		Content:     sb.String(),
		CheckExists: false,
	})
	return nil
}

func (fw *FileWare) filename() string {
	patterns := strings.Split(fw.file.GeneratedFilenamePrefix, "/")
	if len(patterns) == 0 {
		return "cato_generated"
	}
	return patterns[len(patterns)-1]
}

// loadAllMessages will load all messages include nested from parent messages
func (fw *FileWare) loadAllMessages(parents []*protogen.Message) []*protogen.Message {
	if len(parents) == 0 {
		return make([]*protogen.Message, 0)
	}
	results := make([]*protogen.Message, 0)
	// first load self
	results = append(results, parents...)
	for _, message := range parents {
		// load message children
		results = append(results, fw.loadAllMessages(message.Messages)...)
	}
	return results
}

func (fw *FileWare) StoreExtraFiles(files []*models.GenerateFileDesc) {
	fw.extraFiles = append(fw.extraFiles, files...)
}

func (fw *FileWare) GetExtraFiles(_ *common.GenContext) ([]*models.GenerateFileDesc, error) {
	return fw.extraFiles, nil
}
