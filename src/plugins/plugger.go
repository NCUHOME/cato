package plugins

import "google.golang.org/protobuf/reflect/protoreflect"

type Plugger interface {
	GetExtensionType() protoreflect.ExtensionType
	Name() string
}
