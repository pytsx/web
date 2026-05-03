package web

import (
	"net/http"
)

type ParamLocation string

const (
	PathParam   ParamLocation = "path"
	QueryParam  ParamLocation = "query"
	HeaderParam ParamLocation = "header"
)

type DataType string

const (
	String  DataType = "string"
	Integer DataType = "integer"
	Number  DataType = "number"
	Boolean DataType = "boolean"
	Object  DataType = "object"
	Array   DataType = "array"
)

type Route struct {
	method      string
	path        string
	handler     HandlerFunc
	middlewares []MidFunc

	summary     string
	description string
	tags        []string
	params      []ParamDoc
	requestBody *BodyDoc
	responses   map[int]ResponseDoc
}

type ParamDoc struct {
	Name        string
	In          ParamLocation
	Type        DataType
	Required    bool
	Description string
}

type BodyDoc struct {
	Schema      any
	Description string
	Required    bool
}

type ResponseDoc struct {
	Status      int
	Schema      any
	Description string
}

func newRoute(method, path string, handler HandlerFunc, mw ...MidFunc) *Route {
	return &Route{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: mw,
		responses:   make(map[int]ResponseDoc),
	}
}

func (r *Route) GetPath() string {
	return r.path
}

func (r *Route) GetParams() []ParamDoc {
	return r.params
}

func (r *Route) GetDescription() string {
	return r.description
}

func (r *Route) GetSummary() string {
	return r.summary
}

func (r *Route) GetMethod() string {
	return r.method
}

func (r *Route) GetTags() []string {
	return r.tags
}

func (r *Route) GetBody() *BodyDoc {
	return r.requestBody
}
func (r *Route) GetResponses() map[int]ResponseDoc {
	return r.responses
}

func (r *Route) Summary(value string) *Route {
	r.summary = value
	return r
}

func (r *Route) Description(value string) *Route {
	r.description = value
	return r
}

func (r *Route) Tag(tags ...string) *Route {
	r.tags = append(r.tags, tags...)
	return r
}

func (r *Route) Body(schema any, description string, required bool) *Route {
	r.requestBody = &BodyDoc{
		Schema:      schema,
		Description: description,
		Required:    required,
	}
	return r
}

func (r *Route) Param(name string, in ParamLocation, typ DataType, required bool, description string) *Route {
	r.params = append(r.params, ParamDoc{
		Name:        name,
		In:          in,
		Type:        typ,
		Required:    required,
		Description: description,
	})
	return r
}

func (r *Route) PathParam(name string, typ DataType, description string) *Route {
	return r.Param(name, PathParam, typ, true, description)
}

func (r *Route) QueryParam(name string, typ DataType, required bool, description string) *Route {
	return r.Param(name, QueryParam, typ, required, description)
}

func (r *Route) HeaderParam(name string, typ DataType, required bool, description string) *Route {
	return r.Param(name, HeaderParam, typ, required, description)
}

func (r *Route) Response(status int, schema any, description string) *Route {
	r.responses[status] = ResponseDoc{
		Status:      status,
		Schema:      schema,
		Description: description,
	}
	return r
}

func (r *Route) OK(schema any) *Route {
	return r.Response(http.StatusOK, schema, "OK")
}

func (r *Route) Created(schema any) *Route {
	return r.Response(http.StatusCreated, schema, "Created")
}

func (r *Route) BadRequest(schema any) *Route {
	return r.Response(http.StatusBadRequest, schema, "Bad Request")
}

func (r *Route) NotFound(schema any) *Route {
	return r.Response(http.StatusNotFound, schema, "Not Found")
}

func (r *Route) InternalError(schema any) *Route {
	return r.Response(http.StatusInternalServerError, schema, "Internal Server Error")
}
