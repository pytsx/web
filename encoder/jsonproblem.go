package encoder

import "encoding/json"

var jsonProblemContentType = "application/problem+json"

type jsonProblemEncoder[T any] struct {
	data T
}

func NewJSONProblem[T any](data T) jsonProblemEncoder[T] {
	return jsonProblemEncoder[T]{
		data: data,
	}
}

func (j jsonProblemEncoder[T]) Encode() ([]byte, string, error) {
	data, err := json.Marshal(j.data)
	return data, jsonProblemContentType, err
}
