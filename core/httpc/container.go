package httpc

import (
	"net/http"
)

type Container interface {
	Set(key string, runner func(w http.ResponseWriter, r *http.Request))
	Get(key string) (func(w http.ResponseWriter, r *http.Request), bool)
	ToMap() map[string]http.HandlerFunc
	EncodeKey(method, path string) string
	DecodeKey(key string) (string, string, error)
}
