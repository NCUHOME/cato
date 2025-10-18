package tray

import (
	"io"
	"strings"
)

type FieldTray struct {
	tags      []*strings.Builder
	jsonTrans bool
}

func NewFieldTray() *FieldTray {
	cheese := &FieldTray{}
	cheese.tags = make([]*strings.Builder, 0)
	return cheese
}

func (fp *FieldTray) BorrowTagWriter() io.Writer {
	fp.tags = append(fp.tags, new(strings.Builder))
	return fp.tags[len(fp.tags)-1]
}

func (fp *FieldTray) GetTags() []string {
	data := make([]string, len(fp.tags))
	for i, tag := range fp.tags {
		data[i] = tag.String()
	}
	return data
}

func (fp *FieldTray) SetJsonTrans(b bool) {
	fp.jsonTrans = b
}

func (fp *FieldTray) IsJsonTrans() bool {
	return fp.jsonTrans
}
