package db

import (
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/src/plugins/common"
)

var (
	butterFactory []func() common.Butter
	factoryOnce   = new(sync.Once) // todo register as tool
)

func register(builder func() common.Butter) {
	factoryOnce.Do(func() {
		butterFactory = make([]func() common.Butter, 0)
	})
	butterFactory = append(butterFactory, builder)
}

func ChooseButter(desc protoreflect.Descriptor) []common.Butter {
	chosen := make([]common.Butter, 0)
	for index := range butterFactory {
		butter := butterFactory[index]()
		if butter.WorkOn(desc) {
			chosen = append(chosen, butter)
		}
	}
	return chosen
}
