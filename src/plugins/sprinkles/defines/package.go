package defines

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/generated"
	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
)

func init() {
	sprinkles.Register(func() sprinkles.Sprinkle {
		return new(PackageSprinkle)
	})
}

type PackageSprinkle struct {
	value *generated.CatoOptions
}

func (p *PackageSprinkle) FromExtType() protoreflect.ExtensionType {
	return generated.E_CatoOpt
}

func (p *PackageSprinkle) WorkOn(desc protoreflect.Descriptor) bool {
	_, ok := desc.(protoreflect.FileDescriptor)
	return ok
}

func (p *PackageSprinkle) Init(value interface{}) {
	_, ok := value.(*generated.CatoOptions)
	if !ok {
		return
	}
	p.value = value.(*generated.CatoOptions)
}

func (p *PackageSprinkle) Register(ctx *common.GenContext) error {
	fc := ctx.GetNowFileContainer()
	if p.value.GetCatoPackage() != "" {
		fc.SetCatoPackage(p.value.GetCatoPackage())
	}
	if p.value.GetRepoPackage() != "" {
		fc.SetRepoPackage(p.value.GetRepoPackage())
	}
	if p.value.GetRdbRepoPackage() != "" {
		fc.SetRdbRepoPackage(p.value.GetRdbRepoPackage())
	}
	return nil
}
