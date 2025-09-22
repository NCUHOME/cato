package utils

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func GetGoImportName(importPath protogen.GoImportPath) string {
	patterns := strings.Split(GetGoFilePath(importPath), "/")
	if len(patterns) == 0 || patterns[0] == "." {
		return ""
	}
	return patterns[len(patterns)-1]
}

func GetGoFilePath(importPath protogen.GoImportPath) string {
	return strings.Trim(importPath.String(), "\"")
}
