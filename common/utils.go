package common

import (
	"encoding/json"
	"regexp"
)

func Pointer[T any](i T) *T {
	return &i
}

func ToString[T any](i T) (string, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Regex
const (
	bracketRegex    = `{(.*?)}`
	colonReplacer   = `:$1`
	colonRegex      = `:(\w+[^\/])`
	bracketRepalcer = `{$1}`
)

func ToGinPath(path string) string {
	matcher := regexp.MustCompile(bracketRegex)
	return string(matcher.ReplaceAll([]byte(path), []byte(colonReplacer)))
}

func ToSwaggerPath(path string) string {
	matcher := regexp.MustCompile(colonRegex)
	return string(matcher.ReplaceAll([]byte(path), []byte(bracketRepalcer)))
}
