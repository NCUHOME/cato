package db

import (
	"fmt"
	"log"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/models/packs"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
	"github.com/ncuhome/cato/src/plugins/utils"

	"github.com/ncuhome/cato/config"
	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
)

func init() {
	sprinkles.Register(func() sprinkles.Sprinkle {
		return new(TableBasicSprinkle)
	})
}

type TableBasicSprinkle struct {
	value *generated.TableOption
}

func (t *TableBasicSprinkle) extendTmplName() string {
	return "table_extend.tmpl"
}

func (t *TableBasicSprinkle) Init(value interface{}) {
	exValue, ok := value.(*generated.TableOption)
	if !ok {
		log.Fatalf("[-] can not convert %#v to TableOption", value)
	}

	t.value = exValue
}

func (t *TableBasicSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_TableOpt
}

func (t *TableBasicSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.MessageDescriptor)
	return ok
}

func (t *TableBasicSprinkle) Register(ctx *common.GenContext) error {
	if t.value == nil {
		return nil
	}
	pack := &packs.TableBasicTmplPack{
		MessageTypeName: ctx.GetNowMessage().GoIdent.GoName,
	}
	// set extension
	mc := ctx.GetNowMessageContainer()
	_, err := mc.BorrowFieldWriter().Write([]byte("ext *extension"))
	// data object auto set ext package
	mc.SetNeedExtraFile(true)
	mc.NeedEmpty = true
	if err != nil {
		return err
	}
	if t.value.ImplBaseModel {
		err = t.implBasic(ctx)
		if err != nil {
			return err
		}
	}
	// need set ext file
	tmpl := config.GetTemplate(config.TableNameTmpl)
	// check if the table name is simple
	if t.value.NameOption.GetSimpleName() != "" {
		pack.TableName = t.value.NameOption.GetSimpleName()
		return tmpl.Execute(mc.BorrowMethodsWriter(), pack)
	}
	// table name will impl in an extra file
	return tmpl.Execute(mc.BorrowExtraWriter(), pack)
}

func (t *TableBasicSprinkle) implBasic(ctx *common.GenContext) error {
	idField := &models.Field{Name: "Id", GoType: "int64"}
	fields := []*packs.FieldPack{
		{Field: idField, Comments: "primary key, auto increment"},
		{Field: &models.Field{Name: "CreateTime", GoType: "time.Time"}, Comments: "create time, auto gen when create"},
		{Field: &models.Field{Name: "UpdateTime", GoType: "time.Time"}, Comments: "update time, auto gen when update"},
		{Field: &models.Field{Name: "DeleteTime", GoType: "time.Time"}, Comments: "delete time, auto gen when delete"},
	}
	mc := ctx.GetNowMessageContainer()
	_, _ = mc.BorrowFieldWriter().Write([]byte("\n\t // inject for base entity model"))
	colMapper := utils.GetWordMapper(t.value.FieldMapper)
	for _, field := range fields {
		tags := make([]string, 0)
		for _, tag := range mc.GetScopeTags() {
			ss := fmt.Sprintf("%s:\"%s\"", tag.KV.Key, tag.GetTagValue(field.Name))
			tags = append(tags, ss)
		}
		field.Tags = strings.Join(tags, " ")
		err := config.GetTemplate(config.FieldTmpl).Execute(mc.BorrowFieldWriter(), field)
		if err != nil {
			return err
		}
		colName := colMapper(field.Name)
		mc.AddScopeCol(&models.Col{ColName: colName, Field: field.Field})
		colArrivalPack := &packs.ColArrivalTmplPack{
			MessageTypeName: ctx.GetNowMessageTypeName(),
			FieldName:       field.Name,
			ColName:         colName,
			FieldType:       field.GoType,
		}
		tmpl := config.GetTemplate(config.ColArrivalTmpl)
		err = tmpl.Execute(mc.BorrowMethodsWriter(), colArrivalPack)
		if err != nil {
			return err
		}
	}
	mc.AddScopeKey(&models.Key{
		KeyType: generated.DBKeyType_CATO_DB_KEY_TYPE_PRIMARY,
		KeyName: fmt.Sprintf("cato_pk_idx_%s", colMapper(ctx.GetNowMessage().GoIdent.GoName)),
		Fields:  []*models.Field{idField},
	})
	_, err := mc.BorrowFieldWriter().Write([]byte{'\t', '\n'})
	return err
}
