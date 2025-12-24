package tray

import (
	"io"
	"sort"
	"strings"

	"github.com/ncuhome/cato/src/plugins/models"
)

type MessageTray struct {
	methods []*strings.Builder
	fields  []*strings.Builder
	extra   []*strings.Builder
	repo    []*strings.Builder
	rdb     []*strings.Builder
	imports []*strings.Builder

	scopeTags      map[string]*models.Tag
	scopeCols      map[string]*models.Col
	scopeKeys      map[string]*models.Key
	scopeGroupCols map[string][]string

	needExtraFile bool
	asParam       bool
}

func NewMessageTray() *MessageTray {
	return &MessageTray{
		methods:        make([]*strings.Builder, 0),
		fields:         make([]*strings.Builder, 0),
		extra:          make([]*strings.Builder, 0),
		repo:           make([]*strings.Builder, 0),
		rdb:            make([]*strings.Builder, 0),
		imports:        make([]*strings.Builder, 0),
		scopeTags:      make(map[string]*models.Tag),
		scopeCols:      make(map[string]*models.Col),
		scopeKeys:      make(map[string]*models.Key),
		scopeGroupCols: make(map[string][]string),
	}
}

func (mc *MessageTray) BorrowImportWriter() io.Writer {
	mc.imports = append(mc.imports, &strings.Builder{})
	return mc.imports[len(mc.imports)-1]
}

func (mc *MessageTray) GetImports() []string {
	ss := make([]string, len(mc.imports))
	for index := range mc.imports {
		ss[index] = mc.imports[index].String()
	}
	return ss
}

func (mc *MessageTray) BorrowMethodsWriter() io.Writer {
	mc.methods = append(mc.methods, new(strings.Builder))
	return mc.methods[len(mc.methods)-1]
}

func (mc *MessageTray) GetMethods() []string {
	ss := make([]string, len(mc.methods))
	for index := range mc.methods {
		ss[index] = mc.methods[index].String()
	}
	return ss
}

func (mc *MessageTray) BorrowExtraWriter() io.Writer {
	mc.extra = append(mc.extra, new(strings.Builder))
	return mc.extra[len(mc.extra)-1]
}

func (mc *MessageTray) GetExtra() []string {
	ss := make([]string, len(mc.extra))
	for index := range mc.extra {
		ss[index] = mc.extra[index].String()
	}
	return ss
}

func (mc *MessageTray) BorrowRepoWriter() io.Writer {
	mc.repo = append(mc.repo, new(strings.Builder))
	return mc.repo[len(mc.repo)-1]
}

func (mc *MessageTray) BorrowRdbWriter() io.Writer {
	mc.rdb = append(mc.rdb, new(strings.Builder))
	return mc.rdb[len(mc.rdb)-1]
}

func (mc *MessageTray) GetRepo() []string {
	ss := make([]string, len(mc.repo))
	for index := range mc.repo {
		ss[index] = mc.repo[index].String()
	}
	return ss
}

func (mc *MessageTray) GetRdb() []string {
	ss := make([]string, len(mc.rdb))
	for index := range mc.rdb {
		ss[index] = mc.rdb[index].String()
	}
	return ss
}

func (mc *MessageTray) BorrowFieldWriter() io.Writer {
	mc.fields = append(mc.fields, new(strings.Builder))
	return mc.fields[len(mc.fields)-1]
}

func (mc *MessageTray) GetField() []string {
	ss := make([]string, len(mc.fields))
	for index := range mc.fields {
		ss[index] = mc.fields[index].String()
	}
	return ss
}

func (mc *MessageTray) AddScopeCol(col *models.Col) {
	if col == nil {
		return
	}
	colName := col.ColName
	mc.scopeCols[colName] = col
}

func (mc *MessageTray) SetScopeColType(fieldName string, colType string) {
	for _, col := range mc.scopeCols {
		if col.Name == fieldName {
			col.GoType = colType
		}
	}
}

func (mc *MessageTray) AddScopeTag(tag *models.Tag) {
	if tag == nil || tag.KV == nil {
		return
	}
	mc.scopeTags[tag.KV.Key] = tag
}

func (mc *MessageTray) AddScopeKey(key *models.Key) {
	if key == nil || len(key.Fields) == 0 {
		return
	}
	_, ok := mc.scopeKeys[key.KeyName]
	if !ok {
		mc.scopeKeys[key.KeyName] = key
	} else {
		mc.scopeKeys[key.KeyName].Fields = append(mc.scopeKeys[key.KeyName].Fields, key.Fields...)
	}
}

func (mc *MessageTray) GetScopeTags() []*models.Tag {
	tags, index := make([]*models.Tag, len(mc.scopeTags)), 0
	for _, tag := range mc.scopeTags {
		tags[index] = tag
		index++
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].KV.Key < tags[j].KV.Key
	})
	return tags
}

func (mc *MessageTray) GetScopeCols() []*models.Col {
	cols, index := make([]*models.Col, len(mc.scopeCols)), 0
	for _, col := range mc.scopeCols {
		cols[index] = col
		index++
	}
	sort.Slice(cols, func(i, j int) bool {
		return cols[i].Name < cols[j].Name
	})
	return cols
}

func (mc *MessageTray) GetScopeKeys() []*models.Key {
	keys, index := make([]*models.Key, len(mc.scopeKeys)), 0
	for _, key := range mc.scopeKeys {
		keys[index] = key
		index++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].KeyName < keys[j].KeyName
	})
	return keys
}

func (mc *MessageTray) GetColGroups() map[string][]string {
	return mc.scopeGroupCols
}

func (mc *MessageTray) AddColIntoGroup(groupName string, colName string) {
	_, ok := mc.scopeGroupCols[groupName]
	if !ok {
		mc.scopeGroupCols[groupName] = make([]string, 0)
	}
	mc.scopeGroupCols[groupName] = append(mc.scopeGroupCols[groupName], colName)
}

func (mc *MessageTray) SetNeedExtraFile(need bool) {
	mc.needExtraFile = need
}

func (mc *MessageTray) IsNeedExtraFile() bool {
	return mc.needExtraFile
}

func (mc *MessageTray) IsAsParam() bool {
	return mc.asParam
}

func (mc *MessageTray) SetAsParam(b bool) {
	mc.asParam = b
}
