package plugins

import (
	"google.golang.org/protobuf/compiler/protogen"
	"io"
	"strings"
)

type FieldsPlugger struct {
	fieldValue *protogen.Field
	fields     []*strings.Builder
}

func (fp *FieldsPlugger) BorrowWriter() io.Writer {
	fp.fields = append(fp.fields, &strings.Builder{})
	return fp.fields[len(fp.fields)-1]
}

func (fp *FieldsPlugger) GetName() string {
	return fp.fieldValue.GoName
}

func (fp *FieldsPlugger) GetGoType() string {
	// todo: type map
	return ""
}
