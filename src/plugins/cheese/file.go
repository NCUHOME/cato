package cheese

import (
	"io"
	"strings"

	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/utils"
	"google.golang.org/protobuf/compiler/protogen"
)

type FileCheese struct {
	imports       map[string]*models.Import
	appendImports []*strings.Builder
}

func NewFileCheese(file *protogen.File) *FileCheese {
	cheese := new(FileCheese)
	cheese.imports = make(map[string]*models.Import)
	cheese.appendImports = make([]*strings.Builder, 0)

	desc := file.Desc
	for index := 0; index < desc.Imports().Len(); index++ {
		importFile := desc.Imports().Get(index)
		importPackage := string(importFile.FileDescriptor.Package())
		importCatoPath, ok := utils.GetCatoPackageFromFile(importFile.FileDescriptor)
		if !ok {
			continue
		}
		cheese.imports[importPackage] = new(models.Import).Init(importCatoPath)
	}
	return cheese
}

func (fc *FileCheese) GetImportPathAlias(path string) string {
	v, ok := fc.imports[path]
	if !ok {
		return ""
	}
	return v.Alias
}

func (fc *FileCheese) GetImports() []string {
	imports, index := make([]string, len(fc.imports)), 0
	for _, v := range fc.imports {
		imports[index] = v.GetPath()
		index++
	}
	return imports
}

func (fc *FileCheese) BorrowImportsWriter() io.Writer {
	fc.appendImports = append(fc.appendImports, new(strings.Builder))
	return fc.appendImports[len(fc.appendImports)-1]
}
