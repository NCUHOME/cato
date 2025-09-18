package plugins

import (
	"io"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/src/plugins/db"
)

type FieldsPlugger struct {
	fieldValue *protogen.Field
	fields     []*strings.Builder
}

func (fp *FieldsPlugger) BorrowWriter() io.Writer {
	fp.fields = append(fp.fields, &strings.Builder{})
	return fp.fields[len(fp.fields)-1]
}

func (fp *FieldsPlugger) GetContent() string {
	ss := make([]string, len(fp.fields))
	for i, field := range fp.fields {
		ss[i] = field.String()
	}
	return strings.Join(ss, " ")
}

func (fp *FieldsPlugger) Active() (bool, error) {
	descriptor := protodesc.ToFieldDescriptorProto(fp.fieldValue.Desc)
	colExt := new(db.ColumnFieldEx)
	value := proto.GetExtension(descriptor.Options, colExt.FromExtType())
	colExt.Init(value)
	colExt.From(fp.fieldValue)
	colExt.SetWriter(fp.BorrowWriter())
	err := colExt.Register()
	if err != nil {
		return false, err
	}
	return true, nil

}
