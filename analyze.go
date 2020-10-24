package go2typings

import (
	"encoding/json"
	"fmt"
	"go/constant"
	"go/types"
	"path"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
)

type (
	Options struct {
	}
	TypescriptEnumMember struct {
		Name    string
		Value   string
		Comment string
	}
	StructToTS struct {
		structs []*Struct
		seen    map[reflect.Type]*Struct
		opts    *Options
	}
)

func New() *StructToTS {
	return &StructToTS{
		seen: map[reflect.Type]*Struct{},
		opts: &Options{},
	}
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

var cache = map[string][]*packages.Package{}

func getPackages(pkgName string) ([]*packages.Package, error) {
	if len(cache[pkgName]) > 0 {
		return cache[pkgName], nil
	}
	res, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes,
	}, pkgName)

	if err != nil {
		return nil, err
	}
	fmt.Println(pkgName)
	cache[pkgName] = res
	return res, nil
}
func getEnumValues(pkgName, typename string) ([]constant.Value, error) {
	res, err := getPackages(pkgName)
	if err != nil {
		return nil, err
	}
	enums := []constant.Value{}
	pkg := res[0].Types.Scope()
	for _, name := range pkg.Names() {
		v := pkg.Lookup(name)
		// It has format similar to this "type.T".
		baseTypename := path.Base(v.Type().String())
		if v != nil && baseTypename == typename {
			switch t := v.(type) {
			case *types.Const:
				{
					enums = append(enums, t.Val())
				}
			}
		}
	}
	return enums, nil
}

func (s *StructToTS) visitType(t reflect.Type, name, namespace string) {
	k := t.Kind()
	switch {
	case k == reflect.Ptr:
		t = indirect(t)
		s.visitType(t, name, namespace)
	case k == reflect.Struct:
		if isDate(t) {
			break
		}
		if t.Name() != "" {
			name = t.Name()
		}
		s.addType(t, name, namespace)
	case k == reflect.Slice || k == reflect.Array:
		s.visitType(t.Elem(), name, namespace)
	case k == reflect.Map:
		s.visitType(t.Elem(), name, namespace)
		s.visitType(t.Key(), name, namespace)
	case (isNumber(k) || k == reflect.String) && isEnum(t):
		{
			s.addTypeEnum(t, "", "")
		}
	}
}

func isEnum(t reflect.Type) bool {
	return t.PkgPath() != ""
}

func getTypedEnumValues(t reflect.Type) []reflect.Value {
	pkg := t.PkgPath()
	values, err := getEnumValues(pkg, t.String())
	if err != nil {
		panic(err)
	}
	enumStrValues := []reflect.Value{}
	for _, v := range values {
		reflectValue := reflect.New(t).Elem()
		newVal := constant.Val(v)
		switch t.Kind() {
		case reflect.String:
			reflectValue.SetString(constant.StringVal(v))
		case reflect.Int:
			value, ok := constant.Int64Val(v)
			if !ok {
				panic("failed to convert")
			}
			reflectValue.SetInt(value)
		default:
			fmt.Println(reflect.TypeOf(newVal), newVal, reflectValue, v.Kind(), t)
			panic("unknown type")
		}

		enumStrValues = append(enumStrValues, reflectValue)
	}
	return enumStrValues
}

func (s *StructToTS) addTypeEnum(t reflect.Type, name, namespace string) (out *Struct) {
	t = indirect(t)
	if out = s.seen[t]; out != nil {
		return out
	}
	out = MakeStruct(t, name, namespace)
	out.Values = getTypedEnumValues(t)
	out.Type = Enum
	s.seen[t] = out
	s.structs = append(s.structs, out)
	return
}

func (s *StructToTS) addType(t reflect.Type, name, namespace string) (out *Struct) {

	t = indirect(t)

	if out = s.seen[t]; out != nil {
		return out
	}
	out = MakeStruct(t, name, namespace)
	fullName := out.Name
	out.Type = RegularType
	out.Fields = make([]*Field, 0, t.NumField())
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
			extendsType := s.seen[sft].ReferenceName
			out.InheritedType = append(out.InheritedType, extendsType)
			continue
		}
		out.Fields = append(out.Fields, &tf)
	}

	s.structs = append(s.structs, out)
	return
}

func (root *StructToTS) GetTypeName(t reflect.Type) string {
	if root.seen[t] != nil {
		return root.seen[t].ReferenceName
	}
	return root.seen[indirect(t)].ReferenceName
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
