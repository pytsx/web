package encoder

import "encoding/json"

var jsonContentType = "application/json"

type JSON[T any] struct {
	data T
}

func NewJSON[T any](data T) JSON[T] {
	return JSON[T]{
		data: data,
	}
}

func (j JSON[T]) Encode() ([]byte, string, error) {
	data, err := json.Marshal(j.data)

	return data, jsonContentType, err
}
