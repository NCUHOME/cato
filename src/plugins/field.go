package plugins

import (
	"fmt"
	"io"
	"strings"

	"github.com/ncuhome/cato/src/plugins/butter/db"
	"github.com/ncuhome/cato/src/plugins/butter/structs"
	"github.com/ncuhome/cato/src/plugins/cheese"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/utils"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
)

type FieldWorker struct {
	field       *protogen.Field
	tags        []*strings.Builder
	DefaultTags []*models.Kv
}

func NewFieldCheese(field *protogen.Field) *FieldWorker {
	return &FieldWorker{
		field: field,
		tags:  make([]*strings.Builder, 0),
	}
}

func (fp *FieldWorker) RegisterContext(gc *common.GenContext) *common.GenContext {
	fc := cheese.NewFieldCheese()
	ctx := gc.WithField(fp.field, fc)
	return ctx
}

func (fp *FieldWorker) borrowTagWriter() io.Writer {
	fp.tags = append(fp.tags, new(strings.Builder))
	return fp.tags[len(fp.tags)-1]
}

func (fp *FieldWorker) AsTmplPack(ctx *common.GenContext) interface{} {
	commonType := common.MapperGoTypeName(ctx, fp.field.Desc)
	if fp.willAsJsonType() {
		commonType = "string"
	}
	pack := &packs.FieldPack{
		Field: &models.Field{
			Name:   fp.field.GoName,
			GoType: commonType,
		},
	}
	tags := make([]string, len(fp.tags))
	tagMap := make(map[string]struct{})
	for index := range fp.tags {
		raw := fp.tags[index].String()
		tagKey := utils.GetTagKey(raw)
		_, hasTag := tagMap[tagKey]
		if tagKey == "" || hasTag {
			continue
		}
		tags[index] = fp.tags[index].String()
		tagMap[tagKey] = struct{}{}
	}
	pack.Tags = strings.Join(tags, " ")
	return pack
}

func (fp *FieldWorker) Active(ctx *common.GenContext) (bool, error) {
	butter := db.ChooseButter(fp.field.Desc)
	butter = append(butter, structs.ChooseButter(fp.field.Desc)...)

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
	// need register tags in ctx
	for _, scopeTag := range ctx.GetNowMessageContainer().GetScopeTags() {
		if scopeTag.KV == nil {
			continue
		}
		target := fp.borrowTagWriter()
		tagData := fmt.Sprintf("%s:\"%s\"", scopeTag.KV.Key, scopeTag.GetTagValue(fp.field.GoName))
		_, err := target.Write([]byte(tagData))
		if err != nil {
			return false, err
		}
	}
	wr := ctx.GetNowMessageContainer().BorrowFieldWriter()
	// register into field writer
	pack := fp.AsTmplPack(ctx)
	err := config.GetTemplate(config.FieldTmpl).Execute(wr, pack)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fp *FieldWorker) willAsJsonType() bool {
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
