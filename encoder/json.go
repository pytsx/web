package encoder

import "encoding/json"

var jsonContentType = "application/json"

type jsonEncoder[T any] struct {
	data T
}

func NewJSON[T any](data T) jsonEncoder[T] {
	return jsonEncoder[T]{
		data: data,
	}
}

func (j jsonEncoder[T]) Encode() ([]byte, string, error) {
	data, err := json.Marshal(j.data)

	return data, jsonContentType, err
}
