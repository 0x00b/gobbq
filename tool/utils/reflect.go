package utils

import (
	"reflect"
	"strings"
)

func HasJsonName(field reflect.StructField) bool {
	tag := field.Tag.Get("json")
	if tag != "" {
		splitTags := strings.Split(tag, ",")
		if splitTags[0] != "" && splitTags[0] != "-" {
			return true
		}
	}
	return false
}

func GetFiledTag(field reflect.StructField, tag string) []string {
	str := field.Tag.Get(tag)
	var splitTags []string
	if str != "" {
		splitTags = strings.Split(str, ",")
	}
	return splitTags
}

func GetJsonName(field reflect.StructField) string {
	tagName := field.Name
	splitTags := GetFiledTag(field, "json")
	if len(splitTags) > 0 && splitTags[0] != "" {
		tagName = splitTags[0]
	}
	return tagName
}

func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
