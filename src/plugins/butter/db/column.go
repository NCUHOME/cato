package db

import (
	"log"

	"github.com/ncuhome/cato/src/plugins/butter"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/utils"
)

func init() {
	register(func() butter.Butter {
		return new(ColFieldButter)
	})
}

type ColFieldButter struct {
	value *generated.ColumnOption
	tags  map[string]string
}

func (c *ColFieldButter) Init(value interface{}) {
	exValue, ok := value.(*generated.ColumnOption)
	if !ok {
		log.Fatalln("[-] cato ColFieldButter except ColumnOption")
	}
	c.value = exValue
	c.tags = make(map[string]string)
}

func (c *ColFieldButter) FromExtType() protoreflect.ExtensionType {
	return generated.E_ColumnOpt
}

func (c *ColFieldButter) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FieldDescriptor)
	return ok
}

func (c *ColFieldButter) Register(ctx *common.GenContext) error {
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

func (c *ColFieldButter) addColInfo(ctx *common.GenContext) error {
	mc := ctx.GetNowMessageContainer()
	field := ctx.GetNowField()
	fieldType := common.MapperGoTypeName(ctx, field.Desc)
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

func (c *ColFieldButter) addTagInfo(ctx *common.GenContext) error {
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
