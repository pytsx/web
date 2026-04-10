package web

import (
	"context"
	"fmt"
	"net/http"
)

type Encoder interface {
	Encode() (data []byte, contentType string, err error)
}

type HandlerFunc func(ctx context.Context, r *http.Request) Encoder

type Logger func(ctx context.Context, msg string, args ...any)

type App struct {
	log     Logger
	mux     *http.ServeMux
	mw      []MidFunc
	origins []string
}

func New(log Logger, mw ...MidFunc) *App {
	mux := http.NewServeMux()
	return &App{
		log: log,
		mux: mux,
		mw:  mw,
	}
}

func (a *App) Origins() []string {
	return a.origins
}

func (a *App) EnableCORS(origins []string) {
	a.origins = origins
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.origins != nil {

		requestOrigin := r.Header.Get("Origin")
		for _, origin := range a.origins {
			if origin == "*" || origin == requestOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

	a.mux.ServeHTTP(w, r)
}

func (a *App) HandleFunc(method string, group string, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	handlerFunc = applyMiddleware(mw, handlerFunc)
	handlerFunc = applyMiddleware(a.mw, handlerFunc)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := setWriter(r.Context(), w)

		resp := handlerFunc(ctx, r)

		if err := Respond(ctx, w, resp); err != nil {
			a.log(ctx, "web-respond", "ERROR", err)
			return
		}
	})

	finalPath := path
	if group != "" {
		finalPath = fmt.Sprintf("/%s%s", group, path)
	}
	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	a.mux.HandleFunc(finalPath, h)
}
