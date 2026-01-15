package flags

import "strings"

const (
	FlagExtOutDir = "ext_out_dir"
	SwaggerPath   = "swagger_path"
	ApiHost       = "api_host"
)

func ParseProtoOptFlag(param string) map[string]string {
	ss := strings.Split(param, ",")
	m := make(map[string]string)
	for _, s := range ss {
		kv := strings.Split(s, "=")
		if len(kv) != 2 {
			continue
		}
		m[kv[0]] = kv[1]
	}
	return m
}
