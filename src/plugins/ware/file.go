package ware

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
	"github.com/ncuhome/cato/src/plugins/tray"
)

type FileWare struct {
	file    *protogen.File
	context *common.GenContext
}

func NewFileWare(file *protogen.File) *FileWare {
	fc := new(FileWare)
	fc.file = file
	return fc
}

func (fc *FileWare) RegisterContext(gc *common.GenContext) *common.GenContext {
	f := tray.NewFileTray(fc.file)
	ctx := gc.WithFile(fc.file, f)
	return ctx
}

func (fc *FileWare) Active(ctx *common.GenContext) (bool, error) {
	descriptor := protodesc.ToFileDescriptorProto(fc.file.Desc)
	butters := sprinkles.ChooseSprinkle(fc.file.Desc)
	for index := range butters {
		if !proto.HasExtension(descriptor.Options, butters[index].FromExtType()) {
			continue
		}
		value := proto.GetExtension(descriptor.Options, butters[index].FromExtType())
		butters[index].Init(value)
		err := butters[index].Register(ctx)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (fc *FileWare) Complete(_ *common.GenContext) error {
	return nil
}
