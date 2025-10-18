package sprinkles

import (
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/src/plugins/common"
)

type Sprinkle interface {
	FromExtType() protoreflect.ExtensionType
	WorkOn(desc protoreflect.Descriptor) bool
	Init(value interface{})
	Register(ctx *common.GenContext) error
}

var (
	factory     []func() Sprinkle
	factoryOnce = new(sync.Once)
)

func Register(builder func() Sprinkle) {
	factoryOnce.Do(func() {
		factory = make([]func() Sprinkle, 0)
	})
	factory = append(factory, builder)
}

func ChooseSprinkle(desc protoreflect.Descriptor) []Sprinkle {
	chosen := make([]Sprinkle, 0)
	for index := range factory {
		b := factory[index]()
		if b.WorkOn(desc) {
			chosen = append(chosen, b)
		}
	}
	return chosen
}
