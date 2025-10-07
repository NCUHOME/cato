package butter

import (
	"github.com/ncuhome/cato/src/plugins/common"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Butter interface {
	FromExtType() protoreflect.ExtensionType
	WorkOn(desc protoreflect.Descriptor) bool
	Init(value interface{})
	Register(ctx *common.GenContext) error
}
