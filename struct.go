package go2typings

import (
	"fmt"
	"io"
	"path"
	"reflect"
	"strconv"
	"strings"
)

type Kind int

const (
	RegularType = iota
	Enum
)

type Struct struct {
	Type          Kind
	ReferenceName string
	Namespace     string
	Name          string
	Fields        []*Field
	InheritedType []string
	Values        []reflect.Value
	T             reflect.Type
}

const template = `  //%s.%s
  export interface %s %s{%s}
`

const enumTemplate = `  //%s.%s
  export type %s = %s
`

func MakeStruct(t reflect.Type, name, namespace string) *Struct {
	if name == "" {
		name = t.Name()
	}
	if namespace == "" {
		namespace = path.Base(t.PkgPath())
	}

	fullName := capitalize(name)
	out := &Struct{
		Namespace:     namespace,
		Name:          fullName,
		ReferenceName: namespace + "." + fullName,
		InheritedType: []string{},
		T:             t,
	}
	return out

}

func (s *Struct) RenderTo(opts *Options, w io.Writer) (err error) {
	extendsType := ""
	if len(s.InheritedType) != 0 {
		extendsType = fmt.Sprintf("extends %s ", strings.Join(s.InheritedType, ", "))
	}

	fields := ""
	for n, field := range s.Fields {
		name, t := Type(field)
		fields += fmt.Sprintf("\n    %s: %s;", name, t)
		if n == len(s.Fields)-1 {
			fields += "\n  "
		}
	}
	_, err = fmt.Fprintf(w, template, s.T.PkgPath(), s.T.Name(), s.Name, extendsType, fields)
	return
}

func (s *Struct) RenderEnum(opts *Options, w io.Writer) (err error) {
	union := ""
	for i, v := range s.Values {
		k := v.Type().Kind()
		switch k {
		case reflect.String:
			union += strconv.Quote(v.String())
		case reflect.Int:
			_, hasToString := v.Type().MethodByName("String")
			if hasToString {
				union += strconv.Quote(fmt.Sprintf("%v", v))
			} else {
				union += fmt.Sprintf("%d", v.Int())
			}
		}
		if i != len(s.Values)-1 {
			union += " | "
		}
	}
	_, err = fmt.Fprintf(w, enumTemplate, s.T.PkgPath(), s.T.Name(), s.Name, union)
	return
}
