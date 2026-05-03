package swagger

import (
	"reflect"
	"strings"

	"github.com/pytsx/web"
)

func BuildOpenAPI(sw *Swagger, routes []web.Route) OpenAPI {
	spec := OpenAPI{
		OpenAPI: "3.0.3",
		Info: OpenAPIInfo{
			Title:       sw.title,
			Version:     sw.version,
			Description: sw.description,
		},
		Paths: make(map[string]PathItem),
	}

	if sw.baseURL != "" {
		spec.Servers = []OpenAPIServer{{URL: sw.baseURL}}
	}

	for _, route := range routes {
		path := normalizeOpenAPIPath(route.GetPath())
		method := strings.ToLower(route.GetMethod())

		if _, ok := spec.Paths[path]; !ok {
			spec.Paths[path] = PathItem{}
		}

		spec.Paths[path][method] = toOperation(route)
	}

	return spec
}

func normalizeOpenAPIPath(path string) string {
	if path == "" {
		return "/"
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}

func schemaFrom(v any) map[string]interface{} {
	if v == nil {
		return map[string]interface{}{"type": "object"}
	}

	t := reflect.TypeOf(v)

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return schemaFromType(t)
}

func schemaFromType(t reflect.Type) map[string]interface{} {
	switch t.Kind() {
	case reflect.String:
		return map[string]interface{}{"type": "string"}

	case reflect.Bool:
		return map[string]interface{}{"type": "boolean"}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]interface{}{"type": "integer"}

	case reflect.Float32, reflect.Float64:
		return map[string]interface{}{"type": "number"}

	case reflect.Slice, reflect.Array:
		return map[string]interface{}{
			"type":  "array",
			"items": schemaFromType(t.Elem()),
		}

	case reflect.Map:
		return map[string]interface{}{
			"type":                 "object",
			"additionalProperties": true,
		}

	case reflect.Struct:
		props := map[string]interface{}{}
		required := []string{}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			if field.PkgPath != "" {
				continue
			}

			name := fieldName(field)
			if name == "-" {
				continue
			}

			props[name] = schemaFromType(field.Type)

			if !strings.Contains(field.Tag.Get("json"), "omitempty") {
				required = append(required, name)
			}
		}

		out := map[string]interface{}{
			"type":       "object",
			"properties": props,
		}

		if len(required) > 0 {
			out["required"] = required
		}

		return out

	default:
		return map[string]interface{}{"type": "object"}
	}
}

func fieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return lowerFirst(field.Name)
	}

	parts := strings.Split(tag, ",")
	if parts[0] == "" {
		return lowerFirst(field.Name)
	}

	return parts[0]
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}

	return strings.ToLower(s[:1]) + s[1:]
}

func nonEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}

	negative := v < 0
	if negative {
		v = -v
	}

	var buf [20]byte
	i := len(buf)

	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}

	if negative {
		i--
		buf[i] = '-'
	}

	return string(buf[i:])
}