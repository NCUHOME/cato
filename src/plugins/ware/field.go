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

// FieldWare work for parsing message's field option
// field ware is the last node in ware three and has no any sub wares
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
	// field ware has no sub wares
	return []WorkWare{}
}

func (fw *FieldWare) RegisterContext(gc *common.GenContext) *common.GenContext {
	fc := tray.NewFieldTray()
	ctx := gc.WithField(fw.field, fc)
	return ctx
}

func (fw *FieldWare) asTmplPack(fieldType string, tags []string, comments []string) *packs.FieldPack {
	pack := &packs.FieldPack{
		Field: &models.Field{
			Name:   fw.field.GoName,
			GoType: fieldType,
		},
		// comments use "," to join together
		Comments: strings.Join(comments, ", "),
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
	ok, err := CommonWareActive(ctx, fw)
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
		_, err = target.Write([]byte(tagData))
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
	// json trans will change field type when field has json tag
	// origin field will become string type to load json raw value
	// but need better way to check if field has json-trans tag
	if ctx.GetNowFieldContainer().IsJsonTrans() {
		fieldType = "string"
		mc := ctx.GetNowMessageContainer()
		mc.SetScopeColType(fw.field.GoName, fieldType)
	}
	fdc := ctx.GetNowFieldContainer()
	pack := fw.asTmplPack(fieldType, fdc.GetTags(), fdc.GetComments())
	return config.GetTemplate(config.FieldTmpl).Execute(wr, pack)
}

func (fw *FieldWare) StoreExtraFiles(_ []*models.GenerateFileDesc) {}

func (fw *FieldWare) GetExtraFiles(_ *common.GenContext) ([]*models.GenerateFileDesc, error) {
	return []*models.GenerateFileDesc{}, nil
}
