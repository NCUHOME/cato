package utils

import (
	"strings"
	"unicode"

	"github.com/ncuhome/cato/generated"
)

func GetStringsMapper(mapper generated.FieldMapper) func(s []string) string {
	switch mapper {
	case generated.FieldMapper_CATO_FIELD_MAPPER_CAMEL:
		return buildCamelWords
	case generated.FieldMapper_CATO_FIELD_MAPPER_SNAKE_CASE:
		return buildSnakeCaseWords
	default:
		return func(s []string) string { return strings.Join(s, "") }
	}
}

func GetWordMapper(mapper generated.FieldMapper) func(s string) string {
	runner := GetStringsMapper(mapper)
	return func(s string) string {
		ss := GetSplitor(s)(s)
		return runner(ss)
	}
}

func GetSplitor(s string) func(s string) []string {
	if strings.Contains(s, "_") {
		return SplitSnakeCaseWords
	}
	return SplitCamelWords
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
	runes := []rune(s)
	if len(runes) == 0 {
		return nil
	}
	words := make([]string, 0)
	var b strings.Builder
	for i, r := range runes {
		if i == 0 {
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		prev := runes[i-1]
		nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
		// Check if we're at the start of a new word (lowercase to uppercase, or uppercase sequence ending)
		isNewWord := unicode.IsLower(prev) || (unicode.IsUpper(prev) && nextIsLower)
		// Split when current char is uppercase and marks a new word boundary
		shouldSplit := unicode.IsUpper(r) && isNewWord
		if shouldSplit && b.Len() > 0 {
			words = append(words, b.String())
			b.Reset()
		}
		b.WriteRune(unicode.ToLower(r))
	}
	if b.Len() > 0 {
		words = append(words, b.String())
	}
	return words
}

func SplitSnakeCaseWords(s string) []string {
	return strings.Split(s, "_")
}
