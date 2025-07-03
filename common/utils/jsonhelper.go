package utils

import (
	"encoding/json"
)

// NonNilSlice 是一個泛型 wrapper，用來保證 JSON 序列化時 nil slice → []
type NonNilSlice[T any] []T

func (s NonNilSlice[T]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]T(s))
}
