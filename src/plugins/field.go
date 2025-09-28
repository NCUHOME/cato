package plugins

import (
	"io"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/db"
)

type FieldCheese struct {
	field       *protogen.Field
	tags        []*strings.Builder
	DefaultTags []*common.Kv
}

type FieldCheesePack struct {
	*common.FieldPack
	Tags []string
}

func NewFieldCheese(field *protogen.Field) *FieldCheese {
	return &FieldCheese{
		field: field,
		tags:  make([]*strings.Builder, 0),
	}
}

func (fp *FieldCheese) RegisterContext(gc *common.GenContext) *common.GenContext {
	ctx := gc.WithField(fp.field)
	writers := ctx.GetWriters()
	writers.TagWriter = fp.borrowTagWriter
	return ctx
}

func (fp *FieldCheese) borrowTagWriter() io.Writer {
	fp.tags = append(fp.tags, new(strings.Builder))
	return fp.tags[len(fp.tags)-1]
}

func (fp *FieldCheese) AsTmplPack(ctx *common.GenContext) interface{} {
	commonType := common.MapperGoTypeName(ctx, fp.field.Desc)
	if fp.willAsJsonType() {
		commonType = "string"
	}
	pack := &FieldCheesePack{
		FieldPack: &common.FieldPack{
			Name:   fp.field.GoName,
			GoType: commonType,
		},
		Tags: make([]string, len(fp.tags)),
	}
	for index := range fp.tags {
		pack.Tags[index] = fp.tags[index].String()
	}
	return pack
}

func (fp *FieldCheese) Active(ctx *common.GenContext) (bool, error) {
	butter := db.ChooseButter(fp.field.Desc)
	descriptor := protodesc.ToFieldDescriptorProto(fp.field.Desc)
	for index := range butter {
		if !proto.HasExtension(descriptor.Options, butter[index].FromExtType()) {
			continue
		}
		value := proto.GetExtension(descriptor.Options, butter[index].FromExtType())
		butter[index].Init(value)
		err := butter[index].Register(ctx)
		if err != nil {
			return false, err
		}
	}
	wr := ctx.GetWriters().FieldWriter()
	// register into field writer
	pack := fp.AsTmplPack(ctx)
	err := config.GetTemplate(config.CommonFieldTmpl).Execute(wr, pack)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fp *FieldCheese) willAsJsonType() bool {
	descriptor := protodesc.ToFieldDescriptorProto(fp.field.Desc)
	if !proto.HasExtension(descriptor.Options, generated.E_ColumnOpt) {
		return false
	}
	colOpt := proto.GetExtension(descriptor.Options, generated.E_ColumnOpt).(*generated.ColumnOption)
	jsonTransOpt := colOpt.GetJsonTrans()
	if jsonTransOpt == nil {
		return false
	}
	return true
}
