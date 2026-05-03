package web

import (
	"net/http"
	"strings"
)

type Router struct {
	prefix string
	app    *App
}

func newRouter(app *App, prefix string) *Router {
	return &Router{
		prefix: prefix,
		app:    app,
	}
}

func (g *Router) Get(path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	p := joinPath(g.prefix, path)
	return g.app.HandleFunc(http.MethodGet, p, handlerFunc, mw...)
}

func (g *Router) Post(path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	p := joinPath(g.prefix, path)
	return g.app.HandleFunc(http.MethodPost, p, handlerFunc, mw...)
}

func (g *Router) Put(path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	p := joinPath(g.prefix, path)
	return g.app.HandleFunc(http.MethodPut, p, handlerFunc, mw...)
}

func (g *Router) Patch(path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	p := joinPath(g.prefix, path)
	return g.app.HandleFunc(http.MethodPatch, p, handlerFunc, mw...)
}

func (g *Router) Delete(path string, handlerFunc HandlerFunc, mw ...MidFunc) *Route {
	p := joinPath(g.prefix, path)
	return g.app.HandleFunc(http.MethodDelete, p, handlerFunc, mw...)
}

func (g *Router) Subrouter(prefix string) *Router {
	if prefix == "" {
		return g
	}

	return newRouter(g.app, joinPath(g.prefix, prefix))
}

func joinPath(route string, path string) string {
	route = strings.Trim(route, "/")
	path = strings.TrimLeft(path, "/")

	if route == "" && path == "" {
		return "/"
	}

	if route == "" {
		return "/" + path
	}

	if path == "" {
		return "/" + route
	}

	return "/" + route + "/" + path
}
