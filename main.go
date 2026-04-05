package main

import (
	"context"
	"net/http"
)

type ctxKey int

const (
	writerKey ctxKey = iota + 1
)

func setWriter(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, writerKey, w)
}

func GetWriter(ctx context.Context) http.ResponseWriter {
	if v, ok := ctx.Value(writerKey).(http.ResponseWriter); ok {
		return v
	}

	return nil
}
