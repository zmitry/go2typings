package go2typings

import (
	"reflect"
	"strings"
)

type Field struct {
	Name       string `json:"name"`
	TsType     string `json:"type"`
	KeyType    string `json:"keyType,omitempty"`
	ValType    string `json:"valType,omitempty"`
	CanBeNull  bool   `json:"canBeNull"`
	IsOptional bool   `json:"isOptional"`
	IsDate     bool   `json:"isDate"`
	T          reflect.Type
}

func Type(field *Field) (string, string) {
	fieldType := field.TsType
	name := field.Name
	if field.IsDate {
		fieldType = "string"
	}
	if field.CanBeNull {
		fieldType += " | null"
	}
	if field.IsOptional {
		name += "?"
	}
	return name, fieldType
}

func (f *Field) setProps(sf reflect.StructField, sft reflect.Type) (ignore bool) {
	var (
		jsonTag = strings.Split(sf.Tag.Get("json"), ",")
		tsTag   = strings.Split(sf.Tag.Get("ts"), ",")
	)
	if jsonTag[0] == "" {
		return true
	}
	if ignore = len(tsTag) > 0 && tsTag[0] == "-" || len(jsonTag) > 0 && jsonTag[0] == "-"; ignore {
		return true
	}
	if f.Name = sf.Name; len(jsonTag) > 0 && jsonTag[0] != "" {
		f.Name = jsonTag[0]
	}
	f.IsDate = isDate(sft) || len(tsTag) > 0 && tsTag[0] == "date" || sft.Kind() == reflect.Int64 && strings.HasSuffix(f.Name, "TS")
	if len(tsTag) > 1 {
		switch tsTag[1] {
		case "no-null":
			f.CanBeNull = false
		case "null":
			f.CanBeNull = true
		case "optional":
			f.IsOptional = true
		}
	}
	f.IsOptional = f.IsOptional || len(jsonTag) > 1 && jsonTag[1] == "omitempty"
	if tsTag[0] == "string" {
		f.TsType = "string"
	}
	return
}

func isDate(t reflect.Type) bool {
	return t.Name() == "Time" && t.PkgPath() == "time"
}
