package ware

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/tray"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type FieldWare struct {
	field       *protogen.Field
	DefaultTags []*models.Kv
}

func NewFieldWare(field *protogen.Field) *FieldWare {
	return &FieldWare{
		field: field,
	}
}

func (fw *FieldWare) GetDescriptor() protoreflect.Descriptor {
	return fw.field.Desc
}

func (fw *FieldWare) GetSubWares() []WorkWare {
	return []WorkWare{}
}

func (fw *FieldWare) RegisterContext(gc *common.GenContext) *common.GenContext {
	fc := tray.NewFieldTray()
	ctx := gc.WithField(fw.field, fc)
	return ctx
}

func (fw *FieldWare) asTmplPack(fieldType string, tags []string) *packs.FieldPack {
	pack := &packs.FieldPack{
		Field: &models.Field{
			Name:   fw.field.GoName,
			GoType: fieldType,
		},
	}

	filterTags := make([]string, len(tags))
	tagMap := make(map[string]struct{})
	for index := range tags {
		raw := tags[index]
		tagKey := utils.GetTagKey(raw)
		_, hasTag := tagMap[tagKey]
		if tagKey == "" || hasTag {
			continue
		}
		filterTags[index] = raw
		tagMap[tagKey] = struct{}{}
	}
	pack.Tags = strings.Join(filterTags, " ")
	return pack
}

func (fw *FieldWare) Active(ctx *common.GenContext) (bool, error) {
	ok, err := active(ctx, fw)
	if !ok || err != nil {
		return false, err
	}
	fdc := ctx.GetNowFieldContainer()
	// need register tags in ctx
	for _, scopeTag := range ctx.GetNowMessageContainer().GetScopeTags() {
		if scopeTag.KV == nil {
			continue
		}
		target := fdc.BorrowTagWriter()
		tagData := fmt.Sprintf("%s:\"%s\"", scopeTag.KV.Key, scopeTag.GetTagValue(fw.field.GoName))
		_, err := target.Write([]byte(tagData))
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (fw *FieldWare) Complete(ctx *common.GenContext) error {
	wr := ctx.GetNowMessageContainer().BorrowFieldWriter()
	// register into field writer
	fieldType := common.MapperGoTypeNameFromField(ctx, fw.field.Desc)
	if ctx.GetNowFieldContainer().IsJsonTrans() {
		fieldType = "string"
		mc := ctx.GetNowMessageContainer()
		mc.SetScopeColType(fw.field.GoName, fieldType)
	}
	fdc := ctx.GetNowFieldContainer()
	pack := fw.asTmplPack(fieldType, fdc.GetTags())
	return config.GetTemplate(config.FieldTmpl).Execute(wr, pack)
}
