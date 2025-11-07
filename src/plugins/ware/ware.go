package ware

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ncuhome/cato/src/plugins/common"
	"github.com/ncuhome/cato/src/plugins/models"
	"github.com/ncuhome/cato/src/plugins/sprinkles"
)

type WorkWare interface {
	GetDescriptor() protoreflect.Descriptor
	RegisterContext(gc *common.GenContext) *common.GenContext
	Active(ctx *common.GenContext) (bool, error)
	Complete(ctx *common.GenContext) error
	GetSubWares() []WorkWare
	GetExtraFiles(ctx *common.GenContext) ([]*models.GenerateFileDesc, error)
	StoreExtraFiles(files []*models.GenerateFileDesc)
}

func CommonWareActive(ctx *common.GenContext, ware WorkWare) (bool, error) {
	descriptor := ware.GetDescriptor()
	spks := sprinkles.ChooseSprinkle(ware.GetDescriptor())
	for index := range spks {
		if !proto.HasExtension(descriptor.Options(), spks[index].FromExtType()) {
			continue
		}
		value := proto.GetExtension(descriptor.Options(), spks[index].FromExtType())
		spks[index].Init(value)
		err := spks[index].Register(ctx)
		if err != nil {
			return false, err
		}
	}
	// for sub ware generate
	subs := ware.GetSubWares()
	for _, sub := range subs {
		if sub == nil {
			continue
		}
		subCtx := sub.RegisterContext(ctx)
		_, err := sub.Active(subCtx)
		if err != nil {
			return false, err
		}
		err = sub.Complete(subCtx)
		if err != nil {
			return false, err
		}
		extraFiles, err := sub.GetExtraFiles(subCtx)
		if err != nil {
			return false, err
		}
		ware.StoreExtraFiles(extraFiles)
	}
	return true, nil
}
