package db

import (
	"github.com/ncuhome/cato/generated"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Plugin struct {
}

func (p *Plugin) GetExtensionType() protoreflect.ExtensionType {
	return generated.E_DbOpt
}

func (p *Plugin) Name() string {
	return string(p.GetExtensionType().TypeDescriptor().Name())
}
