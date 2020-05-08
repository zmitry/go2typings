package go2typings

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
)

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
	case isNumber(k) && isEnum(t):
		return getTypeName(t)
	case isNumber(k):
		return "number"
	case k == reflect.String && isEnum(t):
		return getTypeName(t)
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
		if st.Type == Enum {
			err := st.RenderEnum(s.opts, w)
			if err != nil {
				return err
			}
			continue
		} else {
			if err = st.RenderTo(s.opts, w); err != nil {
				return err
			}
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

func (root *StructToTS) setStructTypes(s *Struct) {
	for _, field := range s.Fields {
		if field.TsType == "" {
			field.TsType = typeToString(field.T, root.GetTypeName)
		}
	}
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
