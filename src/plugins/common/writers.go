package common

import (
	"io"
)

type WriterProvider func() io.Writer

type ContextWriter struct {
	ImportWriter WriterProvider
	ExtraWriter  WriterProvider

	MethodWriter WriterProvider
	FieldWriter  WriterProvider
	TagWriter    WriterProvider
}
