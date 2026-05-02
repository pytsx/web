package encoder

import "encoding/json"

var jsonProblemContentType = "application/problem+json"

type JSONProblem[T any] struct {
	data T
}

func NewJSONProblem[T any](data T) JSONProblem[T] {
	return JSONProblem[T]{
		data: data,
	}
}

func (j JSONProblem[T]) Encode() ([]byte, string, error) {
	data, err := json.Marshal(j.data)
	return data, jsonProblemContentType, err
}
