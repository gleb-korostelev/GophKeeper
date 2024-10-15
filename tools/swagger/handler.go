package swagger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/gleb-korostelev/GophKeeper/config"
)

// Handler struct
type Handler struct {
	ResponseBody     interface{}
	RequestBody      interface{}
	HandlerFunc      http.HandlerFunc
	Path             string
	Method           string
	Description      string
	ResponseMimeType mimeType
	Tag              string
	Opts             []Option
	IsResponseFile   bool
}

// Parameter struct
type Parameter struct {
	Name        string
	Type        swaggerFieldType
	ItemsType   swaggerFieldType
	Description string
	Required    bool
}

func (h *Handler) AppendRequestBody(schema interface{}) {
	h.Opts = append(h.Opts, bodyOpts{
		BodySchema: schema,
	})
}

func (h *Handler) AddMimeTypeProduce() string {
	if h.ResponseMimeType == "" {
		h.ResponseMimeType = MimeJson
	}

	mimes := []string{
		string(MimeJson),
	}

	if MimeJson != h.ResponseMimeType {
		mimes = append(mimes, string(h.ResponseMimeType))
	}

	if len(mimes) == 2 {
		return fmt.Sprintf(responseProduce, fmt.Sprintf(`"%s", "%s"`, mimes[0], mimes[1]))
	}
	return fmt.Sprintf(responseProduce, fmt.Sprintf(`"%s"`, mimes[0]))
}

func GenerateDoc(nameAPI string, handlers []Handler) (string, error) {
	var (
		handlersByPath = make(map[string][]Handler, len(handlers))
		paths          = make([]string, 0, len(handlers))
		definitions    = make([]interface{}, 0, len(handlers))
	)

	for _, h := range handlers {
		if _, ok := handlersByPath[h.Path]; !ok {
			handlersByPath[h.Path] = []Handler{h}
		} else {
			handlersByPath[h.Path] = append(handlersByPath[h.Path], h)
		}
	}
	// Default err response
	definitions = append(definitions, RpcStatus{})

	for pathHandler, hdls := range handlersByPath {
		handlersPath := make([]string, 0, len(hdls))
		for _, h := range hdls {
			var (
				commonParams    = make([]string, 0)
				responseSchema  = EmptyObject
				handlerResponse string
			)

			handlerResponse = HandlerRaw1

			if h.ResponseBody != nil {
				responseSchema = GetFullSchema(h.ResponseBody)
				handlerResponse = HandlerRawJSON
			}

			if h.IsResponseFile {
				responseSchema = GetFileResponse()
			}

			if h.RequestBody != nil {
				reqBodyType := reflect.TypeOf(h.RequestBody)
				if reqBodyType.Kind() == reflect.Ptr {
					reqBodyType = reqBodyType.Elem()
				}

				if reqBodyType.Kind() == reflect.Struct {
					if reqBodyType.NumField() > 0 {
						var (
							paramStr string
							paramArr []string
						)
						paramStr += `{
											"name": "body",
											"in": "path",
											"required": true,
											"schema": {
												"type": "object",
												"properties": {`
						for i := 0; i < reqBodyType.NumField(); i++ {
							field := reqBodyType.Field(i)
							name := field.Tag.Get("json")
							if name == "" {
								name = field.Name
							}
							nameType := mapKindToSwaggerType(field.Type.Kind())
							paramArr = append(paramArr, fmt.Sprintf(parameterBody, name, nameType))
						}
						paramStr += strings.Join(paramArr, ",") + `}}}`
						commonParams = append(commonParams, paramStr)
					}
				}

			}

			if len(h.Opts) > 0 {

				for _, opt := range h.Opts {
					commonParams = append(commonParams, opt.GetParameter())
				}
			}

			tag := nameAPI
			if len(h.Tag) != 0 {
				tag = h.Tag
			}

			p := fmt.Sprintf(
				handlerResponse,
				strings.ToLower(h.Method),
				h.Description,
				`"parameters": [`+strings.Join(commonParams, ",")+`],`,
				responseSchema,
				GetFullSchema(RpcStatus{}),
				"",
				tag,
			)
			handlersPath = append(handlersPath, p)
		}
		paths = append(paths, fmt.Sprintf(HandlerBrakesRaw, pathHandler, strings.Join(handlersPath, ",")))
	}

	sort.Slice(paths, func(i, j int) bool {
		return paths[i] < paths[j]
	})

	host := config.GetConfigString(config.HttpsHost)

	return fmt.Sprintf(SwaggerJSON, nameAPI, host, strings.Join(paths, ",")), nil
}

func GetFullSchema(item interface{}) string {
	_, realTypeDef := getRealType(reflect.TypeOf(item))
	definition := map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}

	appendDefsToResult(realTypeDef, definition["properties"].(map[string]interface{}), nil, nil)
	schema, err := json.Marshal(definition)
	if err != nil {
		return "{}"
	}
	return string(schema)
}
