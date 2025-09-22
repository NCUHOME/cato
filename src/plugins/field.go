package plugins

import (
	"io"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/db"
)

type FieldsPlugger struct {
	context  *common.GenContext
	field    *protogen.Field
	writters []*strings.Builder
}

func (fp *FieldsPlugger) LoadContext(gc *common.GenContext) {
	fp.context = gc
	fp.field = fp.context.GetNowField()
	fp.writters = make([]*strings.Builder, 0)
}

func (fp *FieldsPlugger) BorrowWriter() io.Writer {
	fp.writters = append(fp.writters, &strings.Builder{})
	return fp.writters[len(fp.writters)-1]
}

func (fp *FieldsPlugger) GetContent() string {
	ss := make([]string, len(fp.writters))
	for i, field := range fp.writters {
		ss[i] = field.String()
	}
	return strings.Join(ss, " ")
}

func (fp *FieldsPlugger) Active() (bool, error) {
	butter := db.ChooseButter(fp.field.Desc)
	descriptor := protodesc.ToFieldDescriptorProto(fp.context.GetNowField().Desc)
	for index := range butter {
		if !proto.HasExtension(descriptor.Options, butter[index].FromExtType()) {
			continue
		}
		value := proto.GetExtension(descriptor.Options, butter[index].FromExtType())
		butter[index].Init(fp.context, value)
		butter[index].SetWriter(fp.BorrowWriter())
		err := butter[index].Register()
		if err != nil {
			return false, err
		}
	}
	if len(fp.writters) == 0 {
		err := config.GetTemplate(config.CommonFieldTmpl).
			Execute(fp.BorrowWriter(), &common.FieldPack{
				Name:   fp.field.GoName,
				GoType: common.MapperGoTypeName(fp.context, fp.field.Desc),
			})
		if err != nil {
			return false, err
		}
	}
	return true, nil

}
