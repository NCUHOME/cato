package plugins

import (
	"fmt"
	"log"
	"strings"

	"github.com/ncuhome/cato/src/plugins/butter/db"
	"github.com/ncuhome/cato/src/plugins/butter/structs"
	"github.com/ncuhome/cato/src/plugins/cheese"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type MessageWorker struct {
	message *protogen.Message
}

type MessageCheesePack struct {
	PackageName string
	ModelName   string

	Fields  []string
	Methods []string
	Imports []string
}

func NewMessageCheese(msg *protogen.Message) *MessageWorker {
	mp := new(MessageWorker)
	mp.message = msg
	return mp
}

// RegisterContext because generate a file from a message, so a file-level writer for a message generates progress
func (mp *MessageWorker) RegisterContext(gc *common.GenContext) *common.GenContext {
	mc := cheese.NewMessageCheese()
	ctx := gc.WithMessage(mp.message, mc)
	return ctx
}

func (mp *MessageWorker) asBasicTmpl() *MessageCheesePack {
	imports := make([]string, len(mp.imports))
	for i, imp := range mp.imports {
		imports[i] = imp.String()
	}
	methods := make([]string, len(mp.methods))
	for index, method := range mp.methods {
		methods[index] = method.String()
	}
	fields := make([]string, len(mp.fields))
	for index, field := range mp.fields {
		fields[index] = field.String()
	}
	return &MessageCheesePack{
		ModelName: mp.message.GoIdent.GoName,
		Fields:    fields,
		Methods:   methods,
	}
}

func (mp *MessageWorker) GenerateFile() string {
	return fmt.Sprintf("%s.cato.go", mp.outputFileName())
}

func (mp *MessageWorker) outputFileName() string {
	patterns := utils.SplitCamelWords(mp.message.GoIdent.GoName)
	mapper := utils.GetStringsMapper(generated.FieldMapper_CATO_FIELD_MAPPER_SNAKE_CASE)
	return mapper(patterns)
}

func (mp *MessageWorker) GenerateContent(ctx *common.GenContext) string {
	sw := new(strings.Builder)
	err := mp.tmpl.Execute(sw, mp.AsTmplPack(ctx))
	if err != nil {
		log.Fatalln("[-] models plugger exec tmpl error, ", err)
	}
	return sw.String()
}

func (mp *MessageWorker) GenerateExtraContent(_ *common.GenContext) string {
	return mp.extra.String()
}

func (mp *MessageWorker) Active(ctx *common.GenContext) (bool, error) {
	descriptor := protodesc.ToDescriptorProto(mp.message.Desc)
	butter := db.ChooseButter(mp.message.Desc)

	butter = append(butter, structs.ChooseButter(mp.message.Desc)...)
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
	// for fields
	for _, field := range mp.message.Fields {
		fp := NewFieldCheese(field)
		fieldCtx := fp.RegisterContext(ctx)
		_, err := fp.Active(fieldCtx)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
