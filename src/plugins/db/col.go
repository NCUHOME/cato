package db

import (
	"io"
	"log"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/utils"
)

func init() {
	register(func() common.Butter {
		return new(ColumnFieldEx)
	})
}

type ColumnFieldEx struct {
	fromField *protogen.Field
	value     *generated.ColumnOption
	tmpl      *template.Template
	tags      map[string]*common.Kv
	tagWriter io.Writer
}

type ColumnFieldExTmlPack struct {
	*common.FieldPack
	Tags []common.Kv
}

func (c *ColumnFieldEx) Init(gc *common.GenContext, value interface{}) {
	exValue, ok := value.(*generated.ColumnOption)
	if !ok {
		log.Fatalln("[-] cato ColumnFieldEx except ColumnOption")
	}
	c.value = exValue
	c.tmpl = config.GetTemplate(c.GetTmplFileName())
	c.fromField = gc.GetNowField()
}

func (c *ColumnFieldEx) SetWriter(writers ...io.Writer) {
	if len(writers) == 0 {
		log.Fatalln("[-] cato ColumnFieldEx except at least one writer")
	}
	c.tagWriter = writers[0]
}

func (c *ColumnFieldEx) GetTmplFileName() string {
	return "column_field.tmpl"
}

func (c *ColumnFieldEx) LoadPlugger() {

}

func (c *ColumnFieldEx) AsTmplPack() interface{} {
	tags, index := make([]common.Kv, len(c.tags)), 0
	for k, v := range c.tags {
		tags[index] = common.Kv{
			Key:   k,
			Value: v.Value,
		}
		index++
	}
	return &ColumnFieldExTmlPack{
		FieldPack: &common.FieldPack{
			Name:   c.fromField.GoName,
			GoType: utils.MapperGoTypeName(c.fromField.Desc),
		},
		Tags: tags,
	}
}

func (c *ColumnFieldEx) FromExtType() protoreflect.ExtensionType {
	return generated.E_ColumnOpt
}

func (c *ColumnFieldEx) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FieldDescriptor)
	return ok
}

func (c *ColumnFieldEx) Register() error {
	// self-tags has the highest priority
	selfTags := c.value.GetTags()
	for _, tag := range selfTags {
		c.tags[tag.TagName] = &common.Kv{
			Key:   tag.TagName,
			Value: tag.TagValue,
		}
	}
	// check if value has a json-trans option
	hasJsonTrans := c.value.GetJsonTrans()
	fieldKind := c.fromField.Desc.Kind()
	if hasJsonTrans && fieldKind == protoreflect.MessageKind {
		// need register into message-field

	}
	packData := c.AsTmplPack()
	return c.tmpl.Execute(c.tagWriter, packData)
}
