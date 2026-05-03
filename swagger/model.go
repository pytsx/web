package swagger

import (
	"github.com/pytsx/web"
)

type OpenAPI struct {
	OpenAPI string                 `json:"openapi"`
	Info    OpenAPIInfo            `json:"info"`
	Servers []OpenAPIServer        `json:"servers,omitempty"`
	Paths   map[string]PathItem    `json:"paths"`
	Comp    map[string]interface{} `json:"components,omitempty"`
}

type OpenAPIInfo struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

type OpenAPIServer struct {
	URL string `json:"url"`
}

type PathItem map[string]Operation

type Operation struct {
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	Tags        []string            `json:"tags,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses"`
}

type Parameter struct {
	Name        string                 `json:"name"`
	In          string                 `json:"in"`
	Required    bool                   `json:"required"`
	Description string                 `json:"description,omitempty"`
	Schema      map[string]interface{} `json:"schema"`
}

type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Content     map[string]MediaType `json:"content"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema map[string]interface{} `json:"schema"`
}

func toOperation(route web.Route) Operation {
	op := Operation{
		Summary:     route.GetSummary(),
		Description: route.GetDescription(),
		Tags:        route.GetTags(),
		Responses:   make(map[string]Response),
	}

	for _, p := range route.GetParams() {
		op.Parameters = append(op.Parameters, Parameter{
			Name:        p.Name,
			In:          string(p.In),
			Required:    p.Required,
			Description: p.Description,
			Schema: map[string]interface{}{
				"type": string(p.Type),
			},
		})
	}

	body := route.GetBody()
	if body != nil {
		op.RequestBody = &RequestBody{
			Description: body.Description,
			Required:    body.Required,
			Content: map[string]MediaType{
				"application/json": {
					Schema: schemaFrom(body.Schema),
				},
			},
		}
	}

	responses := route.GetResponses()
	if len(responses) == 0 {
		op.Responses["200"] = Response{
			Description: "OK",
		}
	}

	for status, resp := range responses {
		op.Responses[itoa(status)] = Response{
			Description: nonEmpty(resp.Description, "Response"),
			Content: map[string]MediaType{
				"application/json": {
					Schema: schemaFrom(resp.Schema),
				},
			},
		}
	}


	return op
}