package go2typings

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/spec"
)

func (s *StructToTS) RenderToSwagger() spec.Swagger {
	def := spec.Definitions{}
	swag := spec.Swagger{SwaggerProps: spec.SwaggerProps{
		Definitions: def,
		Swagger:     "2.0",
	}}
	for _, st := range s.structs {
		schema := s.GenerateOpenApi(st)
		def[st.ReferenceName] = schema
	}
	return swag
}

func typeToSwagger(t reflect.Type, swaggerType spec.Schema, getTypeName GetTypeName) spec.Schema {
	k := t.Kind()
	switch {
	case k == reflect.Ptr:
		t = indirect(t)
		swaggerType.AsNullable()
		return typeToSwagger(t, swaggerType, getTypeName)
	case k == reflect.Struct:
		if isDate(t) {
			swaggerType.AddType("string", "date-time")
			return swaggerType
		}
		ref, err := spec.NewRef("#/definitions/" + getTypeName(t))
		if err != nil {
			fmt.Println(err)
		}
		swaggerType.Ref = ref
		return swaggerType
	case isNumber(k) && isEnum(t):
		ref, err := spec.NewRef("#/definitions/" + getTypeName(t))
		if err != nil {
			fmt.Println(err)
		}
		swaggerType.Ref = ref
		return swaggerType
	case isNumber(k):
		return *swaggerType.Typed("number", "")
	case k == reflect.String && isEnum(t):
		ref, err := spec.NewRef("#/definitions/" + getTypeName(t))
		if err != nil {
			fmt.Println(err)
		}
		swaggerType.Ref = ref
		return swaggerType
	case k == reflect.String:
		return *swaggerType.Typed("string", "")
	case k == reflect.Bool:
		return *swaggerType.Typed("boolean", "")
	case k == reflect.Slice || k == reflect.Array:
		swaggerType.Type = spec.StringOrArray{"array"}
		props := spec.Schema{}
		swaggerType.Nullable = true
		swaggerType.CollectionOf(typeToSwagger(t.Elem(), props, getTypeName))
		return swaggerType
	case k == reflect.Map:
		props := spec.Schema{}
		item := typeToSwagger(t.Elem(), props, getTypeName)
		swaggerType.SchemaProps = spec.MapProperty(&item).SchemaProps
		return swaggerType
	}
	return swaggerType
}

func (root *StructToTS) GenerateOpenApi(s *Struct) spec.Schema {
	t := spec.Schema{SchemaProps: spec.SchemaProps{
		Properties: map[string]spec.Schema{},
	}}
	propertiesTypes := &t
	if s.Type == Enum {
		t.Typed(s.T.Kind().String(), "")
		convertedValues := make([]interface{}, len(s.Values))
		for i, v := range s.Values {
			convertedValues[i] = v.Interface()
		}
		t.WithEnum(convertedValues...)
		return t
	}
	if len(s.InheritedType) != 0 {
		types := make([]spec.Schema, len(s.InheritedType))
		for i, v := range s.InheritedType {
			ref := spec.RefProperty("#/definitions/" + v)
			types[i] = *ref
		}
		propertiesTypes = &spec.Schema{SchemaProps: spec.SchemaProps{
			Properties: map[string]spec.Schema{},
		}}
		t.WithAllOf(types...)
		t.AddToAllOf(*propertiesTypes)
	}
	propertiesTypes.Typed("object", "")
	for _, field := range s.Fields {
		scm := spec.Schema{}
		propertiesTypes.SchemaProps.Properties[field.Name] = typeToSwagger(field.T, scm, root.GetTypeName)
	}
	return t
}
