package db

import (
	"errors"
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
	runner := []func(ctx *common.GenContext) error{
		c.addTagInfo,
		c.addColInfo,
		c.appendColComment,
	}
	var err error
	for _, fn := range runner {
		err = errors.Join(fn(ctx))
	}
	return err
}

func (c *ColFieldSprinkle) addColInfo(ctx *common.GenContext) error {
	mc := ctx.GetNowMessageContainer()
	field := ctx.GetNowField()
	fieldType := common.MapperGoTypeNameFromField(ctx, field.Desc)
	fieldTargetType := fieldType.GoType()
	if fieldType.IsSlice || fieldType.IsStruct {
		fieldTargetType = "string"
	}
	fieldPack := &models.Field{Name: field.GoName, GoType: fieldTargetType}
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
		FieldType:       fieldTargetType,
	}
	tmpl := config.GetTemplate(config.ColArrivalTmpl)
	err := tmpl.Execute(mc.BorrowMethodsWriter(), colArrivalPack)
	if err != nil {
		return err
	}
	// if it has col group, need add col group into group
	if colDesc.GetColGroup() != "" {
		mc.AddColIntoGroup(colDesc.GetColGroup(), colName)
	}
	return nil
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

// appendColComment will append table column comment to struct field
func (c *ColFieldSprinkle) appendColComment(ctx *common.GenContext) error {
	col := c.value.GetColDesc()
	if col == nil || col.Comment == "" {
		return nil
	}
	w := ctx.GetNowFieldContainer().BorrowCommentsWriter()
	_, err := w.Write([]byte(col.Comment))
	return err
}
