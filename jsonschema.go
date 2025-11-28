package nuage

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

func jsonSchemaForType(typ reflect.Type, opts *jsonschema.ForOptions) (*jsonschema.Schema, error) {
	typ = deref(typ)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("josnschema: type is not struct %s", typ)
	}
	if opts == nil {
		opts = &jsonschema.ForOptions{}
	}
	rootSchema, err := jsonschema.ForType(typ, opts)
	if err != nil {
		return nil, err
	}
	fields := reflect.VisibleFields(typ)
	for _, field := range fields {
		jsonName := jsonNameOf(field)
		propertySchema := rootSchema.Properties[jsonName]
		jsonOpts, err := parseJSONSchemaTagOpts(field)
		if err != nil {
			return nil, err
		}
		err = jsonOpts.applyToSchema(propertySchema, false)
		if err != nil {
			return nil, err
		}
		// ignore min and max properties for the root schema because additionalProperties will be false
		// and the properties are pre-defined with required/dependentRequired
		jsonOpts.minProperties = nil
		jsonOpts.maxProperties = nil
		err = jsonOpts.applyToSchema(rootSchema, true)
		if err != nil {
			return nil, err
		}
		rootSchema.Properties[jsonName] = propertySchema
	}
	return rootSchema, nil
}

func jsonSchemaFor[T any](opts *jsonschema.ForOptions) (*jsonschema.Schema, error) {
	typ := reflect.TypeFor[T]()
	return jsonSchemaForType(typ, opts)
}

func jsonNameOf(field reflect.StructField) string {
	jsonTagValue, found := field.Tag.Lookup("json")
	if found && jsonTagValue != "" {
		return strings.Split(jsonTagValue, ",")[0]
	}
	return field.Name
}
