package swagger

import (
	"context"
	"net/http"

	"github.com/pytsx/web"
	"github.com/pytsx/web/encoder"
)

type Config struct {
	Title       string
	Version     string
	Description string
	BaseURL     string
}

type Swagger struct {
	title       string
	version     string
	description string

	jsonPath string
	uiPath   string
	baseURL  string
}

type SwaggerPage struct {
	html string
}

func (s SwaggerPage) Encode() ([]byte, string, error) {
	return []byte(s.html), "text/html", nil
}

func NewSwagger(cfg Config, path string) Swagger {
	jsonPath := "openapi.json"
	uiPath := path

	if path != "" {
		jsonPath = uiPath + "/" + jsonPath
	}

	return Swagger{
		title:       cfg.Title,
		version:     cfg.Version,
		description: cfg.Description,

		baseURL:  cfg.BaseURL,
		jsonPath: jsonPath,
		uiPath:   uiPath,
	}
}

func (s *Swagger) Serve(app *web.App, mw ...web.MidFunc) {
	router := app.Router(s.baseURL)

	openapi := router.Get(s.jsonPath, func(ctx context.Context, r *http.Request) web.Encoder {
		return encoder.NewJSON(BuildOpenAPI(s, app.Routes()))
	}, mw...)

	router.Get(s.uiPath, func(ctx context.Context, r *http.Request) web.Encoder {
		return SwaggerPage{mountHTML(openapi.GetPath())}
	}, mw...)
}

func mountHTML(openapi string) string {
	return `
		<!doctype html>
		<html>
			<head>
				<meta charset="utf-8">
				<title>API Docs</title>
				<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
			</head>
			<body>
				<div id="swagger-ui"></div>
				<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
				<script>
					window.onload = function() {
						window.ui = SwaggerUIBundle({
							url: "` + openapi + `",
							dom_id: "#swagger-ui"
						});
					};
				</script>
			</body>
		</html>
	`
}
