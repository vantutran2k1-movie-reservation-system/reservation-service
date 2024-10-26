package utils

import (
	"reflect"
	"strings"
)

func StructToMap(input any) map[string]any {
	output := make(map[string]any)
	val := reflect.ValueOf(input)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return output
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		tag := field.Tag.Get("json")
		if tag == "-" || (isNilValue(value) && strings.Contains(tag, "omitempty")) || field.PkgPath != "" {
			continue
		}

		tagName := parseJSONTag(tag)
		if tagName == "" {
			tagName = field.Name
		}

		output[tagName] = value.Interface()
	}

	return output
}

func SliceToMaps(slice any) []map[string]any {
	output := make([]map[string]any, 0)
	val := reflect.ValueOf(slice)

	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			structMap := StructToMap(val.Index(i).Interface())
			output = append(output, structMap)
		}
	}

	return output
}

func parseJSONTag(tag string) string {
	if tag == "" {
		return ""
	}
	tagParts := strings.Split(tag, ",")
	return tagParts[0]
}

func isNilValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface, reflect.Chan:
		return v.IsNil()
	default:
		return false
	}
}
