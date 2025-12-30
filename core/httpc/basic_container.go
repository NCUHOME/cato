package httpc

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	ErrInvalidHandlerMapKey = errors.New("InvalidHandlerMapKey")
)

func NewBasicMapContainer() Container {
	c := new(BasicMapContainer)
	c.Init()
	return c
}

type BasicMapContainer struct {
	c    *sync.Mutex
	data map[string]http.HandlerFunc
}

func (b *BasicMapContainer) Init() {
	b.c = &sync.Mutex{}
	b.data = make(map[string]http.HandlerFunc)
}

func (b *BasicMapContainer) EncodeKey(method string, path string) string {
	// need base64
	pathB64 := base64.StdEncoding.EncodeToString([]byte(path))
	return fmt.Sprintf("%s_%s", method, pathB64)
}

func (b *BasicMapContainer) DecodeKey(key string) (method string, path string, err error) {
	pattern := strings.Split(key, "_")
	if len(pattern) != 2 {
		return "", "", fmt.Errorf("err is %w, Invalid key legth: %s", ErrInvalidHandlerMapKey, key)
	}
	pathEncoded := pattern[1]
	pathRaw, err := base64.StdEncoding.DecodeString(pathEncoded)
	if err != nil {
		return "", "", fmt.Errorf("err is %w, Invalid key path encode: %s", err, key)
	}
	return pattern[0], string(pathRaw), nil
}

func (b *BasicMapContainer) Set(key string, runner func(w http.ResponseWriter, r *http.Request)) {
	b.c.Lock()
	defer b.c.Unlock()
	b.data[key] = runner
}

func (b *BasicMapContainer) Get(key string) (func(w http.ResponseWriter, r *http.Request), bool) {
	b.c.Lock()
	defer b.c.Unlock()
	runner, ok := b.data[key]
	return runner, ok
}

func (b *BasicMapContainer) ToMap() map[string]http.HandlerFunc {
	copied := make(map[string]http.HandlerFunc)
	b.c.Lock()
	defer b.c.Unlock()
	for k, v := range b.data {
		copied[k] = v
	}
	return copied
}
