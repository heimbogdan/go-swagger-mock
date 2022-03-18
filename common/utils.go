package common

import "encoding/json"

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
