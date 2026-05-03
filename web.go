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
	routes  []*Route
}

func New(log Logger, mw ...MidFunc) *App {
	return &App{
		log:    log,
		mux:    http.NewServeMux(),
		mw:     mw,
		routes: []*Route{},
	}
}

func (a *App) Origins() []string {
	return a.origins
}

func (a *App) EnableCORS(origins []string) {
	a.origins = origins
}

func (a *App) Use(mw ...MidFunc) {
	a.mw = append(a.mw, mw...)
}

func (a *App) Routes() []Route {

	out := make([]Route, 0, len(a.routes))

	for _, r := range a.routes {
		if r == nil {
			continue
		}

		out = append(out, *r)
	}

	return out
}

func (a *App) Router(prefix string) *Router {
	return newRouter(a, prefix)
}

func (a *App) Get(path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	return a.HandleFunc(http.MethodGet, path, handlerFunc, mw...)
}

func (a *App) Post(path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	return a.HandleFunc(http.MethodPost, path, handlerFunc, mw...)
}

func (a *App) Put(path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	return a.HandleFunc(http.MethodPut, path, handlerFunc, mw...)
}

func (a *App) Patch(path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	return a.HandleFunc(http.MethodPatch, path, handlerFunc, mw...)
}

func (a *App) Delete(path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	return a.HandleFunc(http.MethodDelete, path, handlerFunc, mw...)
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

func (a *App) HandleFunc(method string, path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	route := newRoute(method, path, handlerFunc, mw...)
	a.routes = append(a.routes, route)

	handlerFunc = applyMiddleware(mw, handlerFunc)
	handlerFunc = applyMiddleware(a.mw, handlerFunc)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		resp := handlerFunc(ctx, r)

		if err := Respond(ctx, w, resp); err != nil {
			a.log(ctx, "web-respond", "ERROR", err)
			return
		}
	})

	finalPath := fmt.Sprintf("%s %s", method, path)

	a.mux.HandleFunc(finalPath, h)

	return route
}
