package utils

import (
	"strings"
	"unicode"

	"github.com/ncuhome/cato/generated"
)

func GetStringMapper(mapper generated.FieldMapper) func(s []string) string {
	switch mapper {
	case generated.FieldMapper_CATO_FIELD_MAPPER_CAMEL:
		return buildCamelWords
	case generated.FieldMapper_CATO_FIELD_MAPPER_SNAKE_CASE:
		return buildSnakeCaseWords
	default:
		return func(s []string) string { return strings.Join(s, "") }
	}
}

func buildCamelWords(ss []string) string {
	for index, s := range ss {
		sb := new(strings.Builder)
		sb.WriteRune(unicode.ToUpper(rune(s[0])))
		sb.WriteString(s[1:])
		ss[index] = sb.String()
	}
	return strings.Join(ss, "")
}

func buildSnakeCaseWords(ss []string) string {
	for index, s := range ss {
		ss[index] = strings.ToLower(s)
	}
	return strings.Join(ss, "_")
}

func SplitCamelWords(s string) []string {
	currentWord := new(strings.Builder)
	words := make([]string, 0)
	for _, r := range s {
		if unicode.IsUpper(r) && currentWord.Len() > 0 {
			words = append(words, currentWord.String())
			currentWord.Reset()
		}
		currentWord.WriteRune(unicode.ToLower(r))
	}
	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}
	return words
}

func SpliSnakeCaseWords(s string) []string {
	return strings.Split(s, "_")
}
