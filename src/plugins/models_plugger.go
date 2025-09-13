package plugins

import (
	"github.com/ncuhome/cato/generated"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"io"
	"strings"
)

type ModelsPlugger struct {
	message *protogen.Message

	fields  map[string]*FieldsPlugger
	imports []*strings.Builder
	methods []*strings.Builder
	extra   []*strings.Builder
}

func (mp *ModelsPlugger) init(message *protogen.Message) {
	mp.message = message
	mp.fields = make(map[string]*FieldsPlugger)
	mp.imports = make([]*strings.Builder, 0)
	mp.methods = make([]*strings.Builder, 0)
	mp.extra = make([]*strings.Builder, 0)
}

func (mp *ModelsPlugger) findField(name string) (*protogen.Field, bool) {
	for _, field := range mp.message.Fields {
		if string(field.Desc.Name()) == name {
			return field, true
		}
	}
	return nil, false
}

func (mp *ModelsPlugger) BorrowFieldsWriter(name string) (io.Writer, bool) {
	_, ok := mp.fields[name]
	if !ok {
		fieldDesc, ok := mp.findField(name)
		if !ok {
			return nil, false
		}
		mp.fields[name] = &FieldsPlugger{fieldDesc, make([]*strings.Builder, 0)}
	}
	return mp.fields[name].BorrowWriter(), true
}

func (mp *ModelsPlugger) BorrowMethodsWriter() io.Writer {
	mp.methods = append(mp.methods, new(strings.Builder))
	return mp.methods[len(mp.methods)-1]
}

func (mp *ModelsPlugger) BorrowImportsWriter() io.Writer {
	mp.imports = append(mp.imports, new(strings.Builder))
	return mp.imports[len(mp.imports)-1]
}

func (mp *ModelsPlugger) BorrowExtraWriter() io.Writer {
	mp.extra = append(mp.extra, new(strings.Builder))
	return mp.extra[len(mp.extra)-1]
}

func (mp *ModelsPlugger) GetExtensionType() protoreflect.ExtensionType {
	return generated.E_DbOpt
}

func (mp *ModelsPlugger) GetMessageName() string {
	return mp.message.GoIdent.GoName
}
