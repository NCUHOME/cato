package packs

type JsonTransTmplPack struct {
	MessageTypeName string
	FieldName       string
	FieldType       string
	FieldTypeRaw    string
	IsSlice         bool
	LazyLoad        bool
}
