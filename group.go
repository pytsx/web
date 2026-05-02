package web

type Group struct {
	prefix string
	app    *App
}

func NewGroup(app *App, prefix string) *Group {
	return &Group{
		prefix: prefix,
		app:    app,
	}
}

func (g *Group) Get(path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.Get(g.prefix, path, handlerFunc, mw...)
}

func (g *Group) Post(path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.Post(g.prefix, path, handlerFunc, mw...)
}

func (g *Group) Put(path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.Put(g.prefix, path, handlerFunc, mw...)
}

func (g *Group) Patch(path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.Patch(g.prefix, path, handlerFunc, mw...)
}

func (g *Group) Delete(path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.Delete(g.prefix, path, handlerFunc, mw...)
}

func (g *Group) HandleFunc(method string, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.HandleFunc(method, g.prefix, path, handlerFunc, mw...)
}
