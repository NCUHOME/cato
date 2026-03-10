package packs

type ModelContentTmplPack struct {
	PackageName string
	Imports     []string
	ModelName   string
	Fields      []string
	Methods     []string
	NeedEmpty   bool
}
