package v2

import (
	"reflect"
	"strings"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func indirect(t reflect.Type) reflect.Type {
	k := t.Kind()
	for k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	return t
}

func isNumber(k reflect.Kind) bool {
	switch k {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func isDate(t reflect.Type) bool {
	return t.Name() == "Time" && t.PkgPath() == "time"
}

func isEnum(t reflect.Type) bool {
	return t.PkgPath() != ""
}

type PropertyState int

const (
	Auto PropertyState = iota
	Ignored
	Optional
	Null
	NotNull
)

type ParseResult struct {
	FieldName string
	FieldType string
	Ignore    bool
	State     PropertyState
}

func saveGet(arr []string, i int) string {
	if i >= 0 && i < len(arr) {
		return arr[i]
	}
	return ""
}
func ParseStructTag(structTag reflect.StructTag) (*ParseResult, error) {
	result := &ParseResult{}
	var (
		jsonTag = strings.Split(structTag.Get("json"), ",")
		tsTag   = strings.Split(structTag.Get("ts"), ",")
	)
	jsonTagVal := saveGet(jsonTag, 0)
	tsTagVal := saveGet(tsTag, 0)
	jsonTagOption := saveGet(jsonTag, 1)
	tsTagOptions := saveGet(tsTag, 1)

	if jsonTagVal == "-" || tsTagVal == "-" {
		result.State = Ignored
	}

	if result.State != Ignored {
		result.FieldName = jsonTagVal
		result.FieldType = tsTagVal
		switch tsTagOptions {
		case "no-null":
			result.State = NotNull
		case "null":
			result.State = Null
		case "optional":
			result.State = Optional
		}
		if jsonTagOption == "omitempty" {
			result.State = Optional
		}
	}
	return result, nil
}
