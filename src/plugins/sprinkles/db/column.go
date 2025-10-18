package db

import (
	"log"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
	"github.com/ncuhome/cato/src/plugins/utils"
)

func init() {
	sprinkles.Register(func() sprinkles.Sprinkle {
		return new(ColFieldSprinkle)
	})
}

type ColFieldSprinkle struct {
	value *generated.ColumnOption
	tags  map[string]string
}

func (c *ColFieldSprinkle) Init(value interface{}) {
	exValue, ok := value.(*generated.ColumnOption)
	if !ok {
		log.Fatalln("[-] cato ColFieldSprinkle except ColumnOption")
	}
	c.value = exValue
	c.tags = make(map[string]string)
}

func (c *ColFieldSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_ColumnOpt
}

func (c *ColFieldSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FieldDescriptor)
	return ok
}

func (c *ColFieldSprinkle) Register(ctx *common.GenContext) error {
	// self-tags has the highest priority
	err := c.addTagInfo(ctx)
	if err != nil {
		return err
	}
	// add col func
	err = c.addColInfo(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *ColFieldSprinkle) addColInfo(ctx *common.GenContext) error {
	mc := ctx.GetNowMessageContainer()
	field := ctx.GetNowField()
	fieldType := common.MapperGoTypeNameFromField(ctx, field.Desc)
	// todo: need a better way to check if set string type
	if field.Desc.Kind() == protoreflect.MessageKind {
		fieldType = "string"
	}
	fieldPack := &models.Field{Name: field.GoName, GoType: fieldType}
	colDesc := c.value.GetColDesc()
	if colDesc == nil {
		return nil
	}
	colName := colDesc.FieldName
	if colName == "" {
		colName = string(field.Desc.Name())
	}
	col := &models.Col{ColName: colName, Field: fieldPack}
	mc.AddScopeCol(col)
	colArrivalPack := &packs.ColArrivalTmplPack{
		MessageTypeName: ctx.GetNowMessageTypeName(),
		FieldName:       field.GoName,
		ColName:         colName,
		FieldType:       fieldType,
	}
	tmpl := config.GetTemplate(config.ColArrivalTmpl)
	return tmpl.Execute(mc.BorrowMethodsWriter(), colArrivalPack)
}

func (c *ColFieldSprinkle) addTagInfo(ctx *common.GenContext) error {
	selfTags := c.value.GetTags()
	if len(selfTags) == 0 {
		return nil
	}
	for _, tag := range selfTags {
		t := &models.Tag{
			KV:     &models.Kv{Key: tag.TagName, Value: tag.TagValue},
			Mapper: utils.GetWordMapper(tag.Mapper),
		}
		c.tags[t.KV.Key] = t.GetTagValue(ctx.GetNowField().GoName)
	}
	fc := ctx.GetNowFieldContainer()
	// as tag tmpl pack
	tags, index := make([]models.Kv, len(c.tags)), 0
	for k, v := range c.tags {
		tags[index] = models.Kv{
			Key:   k,
			Value: v,
		}
		index++
	}
	tmpl := config.GetTemplate(config.TagTmpl)
	return tmpl.Execute(fc.BorrowTagWriter(), tags)
}
