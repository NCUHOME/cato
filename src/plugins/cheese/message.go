package cheese

import (
	"io"
	"strings"

	"github.com/ncuhome/cato/src/plugins/models"
)

type MessageCheese struct {
	methods []*strings.Builder
	fields  []*strings.Builder
	extra   []*strings.Builder
	repo    []*strings.Builder
	imports []*strings.Builder

	scopeTags    map[string]*models.Tag
	scopeCols    map[string]*models.Col
	scopeImports map[string]*models.Import
}

func NewMessageCheese() *MessageCheese {
	return &MessageCheese{
		methods:      make([]*strings.Builder, 0),
		fields:       make([]*strings.Builder, 0),
		extra:        make([]*strings.Builder, 0),
		repo:         make([]*strings.Builder, 0),
		imports:      make([]*strings.Builder, 0),
		scopeTags:    make(map[string]*models.Tag),
		scopeCols:    make(map[string]*models.Col),
		scopeImports: make(map[string]*models.Import),
	}
}

func (mc *MessageCheese) BorrowImportWriter() io.Writer {
	mc.imports = append(mc.imports, &strings.Builder{})
	return mc.imports[len(mc.imports)-1]
}

func (mc *MessageCheese) BorrowMethodsWriter() io.Writer {
	mc.methods = append(mc.methods, new(strings.Builder))
	return mc.methods[len(mc.methods)-1]
}

func (mc *MessageCheese) BorrowExtraWriter() io.Writer {
	mc.extra = append(mc.extra, new(strings.Builder))
	return mc.extra[len(mc.extra)-1]
}

func (mc *MessageCheese) BorrowRepoWriter() io.Writer {
	mc.repo = append(mc.repo, new(strings.Builder))
	return mc.repo[len(mc.repo)-1]
}

func (mc *MessageCheese) BorrowFieldWriter() io.Writer {
	mc.fields = append(mc.fields, new(strings.Builder))
	return mc.fields[len(mc.fields)-1]
}

func (mc *MessageCheese) AddScopeCol(col *models.Col) {
	if col == nil {
		return
	}
	colName := col.ColName
	mc.scopeCols[colName] = col
}

func (mc *MessageCheese) AddScopeTag(tag *models.Tag) {
	if tag == nil || tag.KV == nil {
		return
	}
	mc.scopeTags[tag.KV.Key] = tag
}

func (mc *MessageCheese) GetScopeTags() []*models.Tag {
	tags, index := make([]*models.Tag, len(mc.scopeTags)), 0
	for _, tag := range mc.scopeTags {
		tags[index] = tag
		index++
	}
	return tags
}

func (mc *MessageCheese) AsBasicTmpl() []*models.Col {}
