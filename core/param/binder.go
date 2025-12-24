package param

import (
	"net/http"
)

type ResponseBinder interface {
	Marshal(v interface{}) ([]byte, error)
}

type RequestBinder interface {
	Bind(req *http.Request, v interface{}) error
}
