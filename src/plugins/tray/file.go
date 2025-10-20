package tray

import (
	"io"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type FileTray struct {
	imports map[string]*models.Import
	// todo optimize as repo map
	catoPackage    *models.Import
	repoPackage    *models.Import
	rdbRepoPackage *models.Import

	outImports []*strings.Builder
	outContent []*strings.Builder
}

func NewFileTray(file *protogen.File) *FileTray {
	tray := new(FileTray)
	tray.imports = make(map[string]*models.Import)
	tray.outImports = make([]*strings.Builder, 0)
	tray.outContent = make([]*strings.Builder, 0)

	desc := file.Desc
	for index := 0; index < desc.Imports().Len(); index++ {
		importFile := desc.Imports().Get(index)
		importPackage := string(importFile.FileDescriptor.Package())
		importCatoPath, ok := utils.GetCatoPackageFromFile(importFile.FileDescriptor)
		if !ok {
			continue
		}
		tray.imports[importPackage] = new(models.Import).Init(importCatoPath)
	}
	return tray
}

func (fc *FileTray) GetImportPathAlias(path string) string {
	v, ok := fc.imports[path]
	if !ok {
		return ""
	}
	return v.Alias
}

func (fc *FileTray) GetOutImports() []string {
	ss := make([]string, len(fc.outImports))
	for i, v := range fc.outImports {
		ss[i] = v.String()
	}
	return ss
}

func (fc *FileTray) GetOutContent() []string {
	ss := make([]string, len(fc.outContent))
	for i, v := range fc.outContent {
		ss[i] = v.String()
	}
	return ss
}

func (fc *FileTray) GetImports() []string {
	imports, index := make([]string, len(fc.imports)), 0
	for _, v := range fc.imports {
		imports[index] = v.GetPath()
		index++
	}
	return imports
}

func (fc *FileTray) SetCatoPackage(packagePath string) {
	i := new(models.Import).Init(packagePath)
	fc.catoPackage = i
}

func (fc *FileTray) GetCatoPackage() *models.Import {
	return fc.catoPackage
}

func (fc *FileTray) SetRepoPackage(packagePath string) {
	i := new(models.Import).Init(packagePath)
	fc.repoPackage = i
}

func (fc *FileTray) GetRepoPackage() *models.Import {
	return fc.repoPackage
}

func (fc *FileTray) SetRdbRepoPackage(packagePath string) {
	i := new(models.Import).Init(packagePath)
	fc.rdbRepoPackage = i
}

func (fc *FileTray) GetRdbRepoPackage() *models.Import {
	return fc.rdbRepoPackage
}

func (fc *FileTray) BorrowImportWriter() io.Writer {
	fc.outImports = append(fc.outImports, new(strings.Builder))
	return fc.outImports[len(fc.outImports)-1]
}

func (fc *FileTray) BorrowContentWriter() io.Writer {
	fc.outContent = append(fc.outContent, new(strings.Builder))
	return fc.outContent[len(fc.outContent)-1]
}
