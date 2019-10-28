package go2typings

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

type Options struct {
}

func New() *StructToTS {
	return &StructToTS{
		seen: map[reflect.Type]*Struct{},
		opts: &Options{},
	}
}

type StructToTS struct {
	structs []*Struct
	seen    map[reflect.Type]*Struct
	opts    *Options
}

func (s *StructToTS) Add(v interface{}) *Struct { return s.AddWithName(v, "") }

func (s *StructToTS) AddWithName(v interface{}, name string) *Struct {
	var t reflect.Type
	switch v := v.(type) {
	case reflect.Type:
		t = v
	case reflect.Value:
		t = v.Type()
	default:
		t = reflect.TypeOf(v)
	}

	return s.addType(t, name, "")
}

type GetTypeName = func(t reflect.Type) string

// Call this func on each step of type processing.
// This func returns type string representation
func typeToString(t reflect.Type, getTypeName GetTypeName) string {
	k := t.Kind()
	switch {
	case k == reflect.Ptr:
		t = indirect(t)
		return fmt.Sprintf("%s | null", typeToString(t, getTypeName))
	case k == reflect.Struct:
		if isDate(t) {
			return "string"
		}
		return getTypeName(t)
	case isNumber(k):
		return "number"
	case k == reflect.String:
		return "string"
	case k == reflect.Bool:
		return "boolean"
	case k == reflect.Slice || k == reflect.Array:
		return fmt.Sprintf("Array<%s> | null", typeToString(t.Elem(), getTypeName))
	case k == reflect.Interface || t == jsonRawMessageType:
		return "any"
	case k == reflect.Map:
		KeyType, ValType := typeToString(t.Key(), getTypeName), typeToString(t.Elem(), getTypeName)
		return fmt.Sprintf("Record<%s, %s>", KeyType, ValType)
	}
	return t.String()
}

func (s *StructToTS) visitType(t reflect.Type, name, namespace string) {
	k := t.Kind()
	switch {
	case k == reflect.Ptr:
		t = indirect(t)
		s.visitType(t, name, namespace)
		break
	case k == reflect.Struct:
		if isDate(t) {
			break
		}
		if t.Name() != "" {
			name = t.Name()
		}
		s.addType(t, name, namespace)
		break
	case isNumber(k):
	case k == reflect.String:
	case k == reflect.Bool:
		break
	case k == reflect.Slice || k == reflect.Array:
		s.visitType(t.Elem(), name, namespace)
		break
	case k == reflect.Map:
		s.visitType(t.Elem(), name, namespace)
		s.visitType(t.Key(), name, namespace)
		break
	}
}

func (s *StructToTS) addType(t reflect.Type, name, namespace string) (out *Struct) {

	t = indirect(t)

	if out = s.seen[t]; out != nil {
		return out
	}

	if name == "" {
		name = t.Name()
	}
	if namespace == "" {
		namespace = path.Base(t.PkgPath())
	}

	fullName := capitalize(name)
	out = &Struct{
		Namespace:     namespace,
		Name:          fullName,
		ReferenceName: namespace + "." + fullName,
		Fields:        make([]*Field, 0, t.NumField()),
		InheritedType: []string{},
		T:             t,
	}

	s.seen[t] = out
	for i := 0; i < t.NumField(); i++ {
		var (
			sf  = t.Field(i)
			sft = sf.Type
			k   = sft.Kind()
		)
		tf := Field{T: sft}
		if tf.setProps(sf, sft) && !sf.Anonymous {
			continue
		}

		fullFieldName := sft.Name()
		if fullFieldName == "" {
			fullFieldName = sf.Name + fullName
		}
		s.visitType(sf.Type, fullFieldName, namespace)

		if sf.Anonymous && k == reflect.Struct {
			extendsType := s.seen[sft].Name
			out.InheritedType = append(out.InheritedType, extendsType)
			continue
		}
		out.Fields = append(out.Fields, &tf)
	}

	s.structs = append(s.structs, out)
	return
}

func (root *StructToTS) getTypeName(t reflect.Type) string {
	return root.seen[t].ReferenceName
}
func (root *StructToTS) setStructTypes(s *Struct) {
	for _, field := range s.Fields {
		if field.TsType == "" {
			field.TsType = typeToString(field.T, root.getTypeName)
		}
	}
}

func inArray(val int, array []*Struct) bool {
	return len(array) > val && val > 0
}
func (s *StructToTS) RenderTo(w io.Writer) (err error) {
	if _, err = fmt.Fprintf(w, "export namespace %s {\n", s.structs[0].Namespace); err != nil {
		return err
	}
	for i, st := range s.structs {
		s.setStructTypes(st)
		if inArray(i-1, s.structs) && s.structs[i-1].Namespace != st.Namespace {
			if _, err = fmt.Fprintf(w, "export namespace %s {\n", st.Namespace); err != nil {
				return err
			}
		}
		if err = st.RenderTo(s.opts, w); err != nil {
			return err
		}
		if inArray(i+1, s.structs) && s.structs[i+1].Namespace != st.Namespace {
			if _, err = fmt.Fprint(w, "}\n\n"); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprint(w, "}\n\n"); err != nil {
		return err
	}
	return
}

func (s *StructToTS) GenerateFile(path string) (err error) {
	interfacesPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	interfacesFile, err := os.Create(interfacesPath)
	if err != nil {
		return
	}

	if _, err = interfacesFile.WriteString("/* tslint:disable */\n"); err != nil {
		return
	}
	if _, err = interfacesFile.WriteString("/* eslint-disable */\n"); err != nil {
		return
	}

	if err := s.RenderTo(interfacesFile); err != nil {
		return err
	}
	f, err := os.Open(interfacesPath)
	if err != nil {
		return
	}
	if err = f.Close(); err != nil {
		return
	}
	return
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

var jsonRawMessageType = reflect.TypeOf((*json.RawMessage)(nil)).Elem()

func capitalize(s string) string {
	return strings.Title(s)
}
