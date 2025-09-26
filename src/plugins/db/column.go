package db

import (
	"fmt"
	"log"
	"text/template"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
)

func init() {
	register(func() common.Butter {
		return new(ColFieldButter)
	})
}

type ColFieldButter struct {
	value *generated.ColumnOption
	tmpl  *template.Template
	tags  map[string]*common.Kv
}

type ColFieldButterTmplPack struct {
	*common.FieldPack
	Tags []common.Kv
}

func (c *ColFieldButter) Init(value interface{}) {
	exValue, ok := value.(*generated.ColumnOption)
	if !ok {
		log.Fatalln("[-] cato ColFieldButter except ColumnOption")
	}
	c.value = exValue
	c.tmpl = config.GetTemplate(c.GetTmplFileName())
	c.tags = make(map[string]*common.Kv)
}

func (c *ColFieldButter) GetTmplFileName() string {
	return config.CommonTagTmpl
}

func (c *ColFieldButter) LoadPlugger() {

}

func (c *ColFieldButter) AsTmplPack(ctx *common.GenContext) interface{} {
	tags, index := make([]common.Kv, len(c.tags)), 0
	for k, v := range c.tags {
		tags[index] = common.Kv{
			Key:   k,
			Value: v.Value,
		}
		index++
	}
	return tags
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
	selfTags := c.value.GetTags()
	for _, tag := range selfTags {
		c.tags[tag.TagName] = &common.Kv{
			Key:   tag.TagName,
			Value: tag.TagValue,
		}
	}
	writers := ctx.GetWriters()
	// check if the value has a json-trans option
	hasJsonTrans := c.value.GetJsonTrans() != nil
	if hasJsonTrans {
		// register str type raw field
		fieldWriter := writers.FieldWriter()
		extraFieldData := &common.FieldPack{
			Name:   fmt.Sprintf("%sRaw", ctx.GetNowField().GoName),
			GoType: "string",
		}
		err := config.GetTemplate(config.CommonFieldTmpl).Execute(fieldWriter, extraFieldData)
		if err != nil {
			return err
		}
		// register json trans
	}
	packData := c.AsTmplPack(ctx)
	return c.tmpl.Execute(writers.TagWriter(), packData)
}
