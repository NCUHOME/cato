package common

type Kv struct {
	Key   string
	Value string
}

type Tag struct {
	KV     *Kv
	Mapper func(s string) string
}
