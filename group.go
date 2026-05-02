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

func (g *Group) Get(group string, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.Get(group, path, handlerFunc, mw...)
}

func (g *Group) Post(group string, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.Post(group, path, handlerFunc, mw...)
}

func (g *Group) Put(group string, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.Put(group, path, handlerFunc, mw...)
}

func (g *Group) Patch(group string, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.Patch(group, path, handlerFunc, mw...)
}

func (g *Group) Delete(group string, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.Delete(group, path, handlerFunc, mw...)
}

func (g *Group) HandleFunc(method string, path string, handlerFunc HandlerFunc, mw ...MidFunc) {
	g.app.HandleFunc(method, g.prefix, path, handlerFunc, mw...)
}
