package httpc

import (
	"net/http"
)

type HttpService interface {
	BaseUrl() string
	Urls() map[string][]string
	Handlers(method, uri string) (http.HandlerFunc, bool)
}
