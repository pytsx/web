# web

### mini framework web

## Exemplo de uso
```
func main() {
  webapi := web.New(logger)

  root := webapi.Router("api")

  root.Get("/home", func(ctx context.Context, r *http.Request) web.Encoder {
    return encoder.NewJSON(map[string]string{"hello": "world"})
  })

  swcfg := swagger.Config{
    Title:       "API",
    Version:     "1.0.0",
    Description: "API documentation",
    BaseURL:     "/",
  }

  sw := swagger.NewSwagger(swcfg, "/docs")
  sw.Serve(webapi)

  srv := &web.Server{
    Addr:    ":8080",
    Handler: webapi,
  }

  if err := srv.ListenAndServe(context.Background()); err != nil {
    panic(err)
  }
}
```