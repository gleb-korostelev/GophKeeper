package swagger

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Option interface
type Option interface {
	GetParameter() string
}

// QueryOpt param
type QueryOpt Parameter

// GetParameter func
func (q QueryOpt) GetParameter() string {
	queryArray := ""
	if q.Type == Array {
		queryArray = fmt.Sprintf(parameterArray, q.ItemsType)
	}
	return fmt.Sprintf(parameter, q.Name, query, q.Required, q.Description, q.Type, queryArray)
}

// PathOpt param
type PathOpt Parameter

// GetParameter func
func (q PathOpt) GetParameter() string {
	typeF := q.Type
	if typeF == "" {
		typeF = String
	}
	return fmt.Sprintf(parameter, q.Name, path, q.Required, q.Description, typeF, "")
}

type bodyOpts struct {
	BodySchema interface{}
}

// GetParameter func
func (q bodyOpts) GetParameter() string {
	return ""
}

// FormDataOpt param
type FormDataOpt Parameter

// GetParameter func
func (q FormDataOpt) GetParameter() string {
	queryArray := ""
	if q.Type == Array {
		queryArray = fmt.Sprintf(parameterArray, q.ItemsType)
	}
	return fmt.Sprintf(parameterFormData, q.Name, q.Required, q.Description, q.Type, queryArray)
}

// HeaderOpt param
type HeaderOpt Parameter

// GetParameter func
func (q HeaderOpt) GetParameter() string {
	return fmt.Sprintf(parameter, q.Name, header, q.Required, q.Description, q.Type, "")
}

// MergeOptionsJSON func
func MergeOptionsJSON(opts ...Option) ([]string, []map[string]interface{}, string) {
	if len(opts) == 0 {
		return []string{}, nil, ""
	}

	var consumes []string
	params := make([]map[string]interface{}, 0, len(opts))
	var requestBodyStr string

	for _, o := range opts {
		switch opt := o.(type) {
		case FormDataOpt:
			consumes = append(consumes, fmt.Sprintf(requestConsumes, MimeMultipart))
		case bodyOpts:
			if schemaStr, ok := opt.BodySchema.(string); ok {
				var schema map[string]interface{}
				err := json.Unmarshal([]byte(schemaStr), &schema)
				if err != nil {
					return nil, nil, ""
				}
				requestBodyStr = fmt.Sprintf(`{"required": true, "content": {"application/json": {"schema": %s}}}`, toJson(schema))
			} else {
				_, realTypeDef := getRealType(reflect.TypeOf(opt.BodySchema))
				requestBody := map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				}
				appendDefsToResult(realTypeDef, requestBody["properties"].(map[string]interface{}), nil, nil)
				schema, err := json.Marshal(requestBody)
				if err != nil {
					return nil, nil, ""
				}
				requestBodyStr = fmt.Sprintf(`{"required": true, "content": {"application/json": {"schema": %s}}}`, string(schema))
			}

		default:
			param := map[string]interface{}{}
			params = append(params, param)
		}
	}

	return consumes, params, requestBodyStr
}

func toJson(data interface{}) string {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(result)
}
