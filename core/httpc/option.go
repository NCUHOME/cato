package httpc

import (
	"net/http"
)

type HttpOption interface {
	Next(t http.HandlerFunc) http.HandlerFunc
}

func WrapHandler(f http.HandlerFunc, opts ...HttpOption) http.HandlerFunc {
	for index := 0; index < len(opts); index++ {
		f = opts[index].Next(f)
	}
	return f
}
