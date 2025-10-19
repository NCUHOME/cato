package tray

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type FileTray struct {
	imports map[string]*models.Import
	// todo optimize as repo map
	catoPackage        *models.Import
	catoExtPackage     *models.Import
	repoPackage        *models.Import
	rdbRepoPackage     *models.Import
	httpHandlerPackage *models.Import
}

func NewFileTray(file *protogen.File) *FileTray {
	cheese := new(FileTray)
	cheese.imports = make(map[string]*models.Import)

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

func (fc *FileTray) GetImportPathAlias(path string) string {
	v, ok := fc.imports[path]
	if !ok {
		return ""
	}
	return v.Alias
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

func (fc *FileTray) SetCatoExtPackage(packagePath string) {
	i := new(models.Import).Init(packagePath)
	fc.catoExtPackage = i
}

func (fc *FileTray) GetCatoExtPackage() *models.Import {
	return fc.catoExtPackage
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

func (fc *FileTray) SetHttpHandlerPackage(packagePath string) {
	i := new(models.Import).Init(packagePath)
	fc.httpHandlerPackage = i
}

func (fc *FileTray) GetHttpHandlerPackage() *models.Import {
	return fc.httpHandlerPackage
}
