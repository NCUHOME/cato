package plugins

import (
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/db"
	"github.com/ncuhome/cato/src/plugins/utils"
)

type MessageCheese struct {
	message *protogen.Message

	imports []*strings.Builder
	methods []*strings.Builder
	extra   []*strings.Builder
	fields  []*strings.Builder
	tmpl    *template.Template
}

type MessageCheesePack struct {
	PackageName string
	ModelName   string

	Fields  []string
	Methods []string
	Imports []string
}

func NewMessageCheese(msg *protogen.Message) *MessageCheese {
	mp := new(MessageCheese)
	mp.message = msg
	mp.fields = make([]*strings.Builder, 0)
	mp.imports = make([]*strings.Builder, 0)
	mp.methods = make([]*strings.Builder, 0)
	mp.extra = make([]*strings.Builder, 0)
	return mp
}

// RegisterContext because generate a file from a message, so a file-level writer for a message generates progress
func (mp *MessageCheese) RegisterContext(gc *common.GenContext) *common.GenContext {
	ctx := gc.WithMessage(mp.message)
	writers := new(common.ContextWriter)
	writers.ImportWriter = mp.borrowImportsWriter
	writers.MethodWriter = mp.borrowMethodsWriter
	writers.ExtraWriter = mp.borrowExtraWriter

	writers.FieldWriter = mp.borrowFieldWriter
	ctx.SetWriters(writers)
	return ctx
}

func (mp *MessageCheese) borrowMethodsWriter() io.Writer {
	mp.methods = append(mp.methods, new(strings.Builder))
	return mp.methods[len(mp.methods)-1]
}

func (mp *MessageCheese) borrowImportsWriter() io.Writer {
	mp.imports = append(mp.imports, new(strings.Builder))
	return mp.imports[len(mp.imports)-1]
}

func (mp *MessageCheese) borrowExtraWriter() io.Writer {
	mp.extra = append(mp.extra, new(strings.Builder))
	return mp.extra[len(mp.extra)-1]
}

func (mp *MessageCheese) borrowFieldWriter() io.Writer {
	mp.fields = append(mp.fields, new(strings.Builder))
	return mp.fields[len(mp.fields)-1]
}

func (mp *MessageCheese) AsTmplPack(ctx *common.GenContext) *MessageCheesePack {
	imports := make([]string, len(mp.imports))
	for i, imp := range mp.imports {
		imports[i] = imp.String()
	}
	// 组织namespace的imports
	for _, importPath := range ctx.GetImports() {
		imports = append(imports, importPath)
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
		PackageName: utils.GetGoImportName(ctx.GetNowFile().GoImportPath),
		Imports:     imports,
		ModelName:   mp.message.GoIdent.GoName,
		Fields:      fields,
		Methods:     methods,
	}
}

func (mp *MessageCheese) GetTemplateName() string {
	return "models.tmpl"
}

func (mp *MessageCheese) Init(template *template.Template) {
	mp.tmpl = template
}

func (mp *MessageCheese) GenerateFile() string {
	return fmt.Sprintf("%s.cato.go", mp.outputFileName())
}

func (mp *MessageCheese) outputFileName() string {
	patterns := utils.SplitCamelWords(mp.message.GoIdent.GoName)
	mapper := utils.GetStringMapper(generated.FieldMapper_CATO_FIELD_MAPPER_SNAKE_CASE)
	return mapper(patterns)
}

func (mp *MessageCheese) GenerateContent(ctx *common.GenContext) string {
	sw := new(strings.Builder)
	err := mp.tmpl.Execute(sw, mp.AsTmplPack(ctx))
	if err != nil {
		log.Fatalln("[-] models plugger exec tmpl error, ", err)
	}
	return sw.String()
}

func (mp *MessageCheese) Active(ctx *common.GenContext) (bool, error) {
	descriptor := protodesc.ToDescriptorProto(mp.message.Desc)
	butter := db.ChooseButter(mp.message.Desc)
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
