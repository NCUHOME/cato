package plugins

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/ncuhome/cato/src/plugins/common"
)

type FileCheese struct {
	file    *protogen.File
	context *common.GenContext
}

func NewFileCheese(file *protogen.File) *FileCheese {
	fc := new(FileCheese)
	fc.file = file
	return fc
}

func (fc *FileCheese) RegisterContext(gc *common.GenContext) *common.GenContext {
	ctx := gc.WithFile(fc.file)
	return ctx
}
