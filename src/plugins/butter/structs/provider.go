package structs

import (
	"sync"

	"github.com/ncuhome/cato/src/plugins/butter"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	butterFactory []func() butter.Butter
	factoryOnce   = new(sync.Once) // todo register as tool
)

func register(builder func() butter.Butter) {
	factoryOnce.Do(func() {
		butterFactory = make([]func() butter.Butter, 0)
	})
	butterFactory = append(butterFactory, builder)
}

func ChooseButter(desc protoreflect.Descriptor) []butter.Butter {
	chosen := make([]butter.Butter, 0)
	for index := range butterFactory {
		b := butterFactory[index]()
		if b.WorkOn(desc) {
			chosen = append(chosen, b)
		}
	}
	return chosen
}
