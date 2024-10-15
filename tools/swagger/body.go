package swagger

import (
	"fmt"
	"reflect"
	"strings"
)

type constantMapItem struct {
	Example string
	Values  []string
}

var constantMap = map[string]*constantMapItem{}

func GetDefSchema(item interface{}) string {
	name, _ := getRealType(reflect.TypeOf(item))
	return fmt.Sprintf(bodySchema, name)
}

func GetFileResponse() string {
	return fileResponse
}

func prepareDefs(items ...interface{}) map[string]reflect.Type {
	toGenerate := map[string]reflect.Type{}

	for _, inner := range items {
		if inner == nil {
			continue
		}

		name, realTypeDef := getRealType(reflect.TypeOf(inner))
		if isSimple(realTypeDef.Kind()) {
			continue
		}

		if _, ok := toGenerate[name]; !ok {
			toGenerate[name] = realTypeDef
		}
	}

	hasNewRecords := false
	for {
		for _, item := range toGenerate {
			for i := 0; i < item.NumField(); i++ {
				fieldData := item.Field(i)
				name, realTypeDef := getRealType(skipPointers(fieldData.Type))

				if isSimple(realTypeDef.Kind()) {
					continue
				}

				if _, ok := toGenerate[name]; !ok {
					toGenerate[name] = realTypeDef
					hasNewRecords = true
				}
			}
		}

		if hasNewRecords {
			hasNewRecords = false
			continue
		}

		break
	}

	return toGenerate
}

func appendDefsToResult(def reflect.Type, innerProps map[string]interface{}, allDefs map[string]reflect.Type, required *[]string) {
	if def.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < def.NumField(); i++ {
		swaggerFieldDefinition := map[string]interface{}{}

		field := def.Field(i)

		if field.Anonymous {
			embeddedTypeName := getTypeName(field.Type)

			if v, ok := allDefs[embeddedTypeName]; ok {
				appendDefsToResult(v, innerProps, allDefs, required)
				continue
			}
		}

		fieldName := field.Name

		if js := field.Tag.Get("json"); len(js) > 0 {
			if js == "-" {
				continue
			}
			if strings.Contains(js, ",omitempty") {
				js = strings.TrimSuffix(js, ",omitempty")
			}
			fieldName = js
		}

		if swag := field.Tag.Get("swag"); len(swag) > 0 {
			for _, group := range strings.Split(swag, ";") {
				if len(group) == 0 {
					continue
				}

				parsedGroup := strings.Split(group, ":")

				if parsedGroup[0] == "required" {
					*required = append(*required, fieldName)
				}
			}
		}

		topType := skipPointers(field.Type)
		swType := mapKindToSwaggerType(topType.Kind())

		swaggerFieldDefinition["type"] = swType

		if swType == "map" {
			swaggerFieldDefinition["type"] = "object"
		}

		name, innerType := getRealType(field.Type)

		simple := isSimple(innerType.Kind())

		//nolint:exhaustive // exhaustive enum check is not required here
		switch swType {
		case Array:
			if simple {
				innerSw := mapKindToSwaggerType(innerType.Kind())

				itemsMap := map[string]interface{}{
					"type": innerSw,
				}

				if con, ok1 := constantMap[name]; ok1 {
					itemsMap["enum"] = con.Values
					itemsMap["example"] = con.Example
				}

				swaggerFieldDefinition["items"] = itemsMap
			} else {
				swaggerFieldDefinition["items"] = GetFullSchemaMap(innerType)
			}
		case Map:
			if simple {
				innerSw := mapKindToSwaggerType(innerType.Kind())

				itemsMap := map[string]interface{}{
					"type": innerSw,
				}

				if con, ok1 := constantMap[name]; ok1 {
					itemsMap["enum"] = con.Values
					itemsMap["example"] = con.Example
				}

				swaggerFieldDefinition["additionalProperties"] = itemsMap
			} else {
				swaggerFieldDefinition["additionalProperties"] = GetFullSchemaMap(innerType)
			}

		default:

			if simple {
				innerSw := mapKindToSwaggerType(innerType.Kind())

				if con, ok1 := constantMap[name]; ok1 {
					swaggerFieldDefinition["enum"] = con.Values
					swaggerFieldDefinition["example"] = con.Example
				}

				swaggerFieldDefinition["type"] = innerSw
			} else {
				swaggerFieldDefinition = GetFullSchemaMap(innerType)
			}
		}

		innerProps[fieldName] = swaggerFieldDefinition
	}
}

func BuildDefs(items ...interface{}) map[string]interface{} {
	allDefs := prepareDefs(items...)
	result := make(map[string]interface{})

	for key, def := range allDefs {
		topProps := map[string]interface{}{
			"type": "object",
		}

		innerProps := map[string]interface{}{}
		topProps["properties"] = innerProps

		result[key] = topProps

		var required []string

		appendDefsToResult(def, innerProps, allDefs, &required)

		if len(required) > 0 {
			topProps["required"] = required
		}
	}

	return result
}

func isSimple(kind reflect.Kind) bool {
	if kind != reflect.Ptr && kind != reflect.Slice && kind != reflect.Map && kind != reflect.Struct &&
		kind != reflect.Invalid {
		return true
	}

	return false
}

func getRealType(expectedToBeStruct reflect.Type) (string, reflect.Type) {
	val := expectedToBeStruct

	if val.Kind() == reflect.Struct {
		return getTypeName(val), val
	}

	if val.Kind() == reflect.Slice {
		n, t := getRealType(expectedToBeStruct.Elem())
		return n, t
	}
	if val.Kind() == reflect.Ptr {
		n, t := getRealType(expectedToBeStruct.Elem())
		return n, t
	}
	if val.Kind() == reflect.Map {
		n, t := getRealType(expectedToBeStruct.Elem())
		return n, t
	}

	return getTypeName(val), val
}

var replacementArr = []string{"/", "-", " ", "{}", "[", "]"}

func getTypeName(val reflect.Type) string {
	result := fmt.Sprintf("%v_%v", val.PkgPath(), val.Name())
	for _, i := range replacementArr {
		result = strings.ReplaceAll(result, i, ".")
	}

	return result
}

func skipPointers(expectedToBeStruct reflect.Type) reflect.Type {
	if expectedToBeStruct == nil {
		return nil
	}

	if expectedToBeStruct.Kind() == reflect.Ptr {
		n := skipPointers(expectedToBeStruct.Elem())
		return n
	}

	return expectedToBeStruct
}

func mapKindToSwaggerType(kind reflect.Kind) swaggerFieldType {
	swType := swaggerFieldType("")

	//nolint:exhaustive // Exhaustive enum check is not required here
	switch kind {
	case reflect.Map:
		swType = Map
	case reflect.Struct:
		swType = Object
	case reflect.Interface:
		swType = Object
	case reflect.Slice:
		swType = Array
	case reflect.Bool:
		swType = Boolean
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Int64:
		swType = Integer
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		swType = Float64
	case reflect.String:
		swType = String
	default:
		swType = String
	}

	return swType
}

func GetFullSchemaMap(item reflect.Type) map[string]interface{} {
	definition := map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}

	appendDefsToResult(item, definition["properties"].(map[string]interface{}), nil, nil)
	return definition
}
