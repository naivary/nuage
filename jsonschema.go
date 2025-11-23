package nuage

import (
	"errors"
	"reflect"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

func jsonSchemaFor[T any](opts *jsonschema.ForOptions) (*jsonschema.Schema, error) {
	if !isStruct[T]() {
		return nil, errors.New("jsonschema: input is not a struct")
	}
	if opts == nil {
		opts = &jsonschema.ForOptions{}
	}
	schema, err := jsonschema.For[T](opts)
	if err != nil {
		return nil, err
	}
	fields, err := fieldsOf[T]()
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		jsonName := jsonNameOf(field)
		propertySchema := schema.Properties[jsonName]
		jsonOpts, err := parseJSONSchemaTagOpts(field)
		if err != nil {
			return nil, err
		}
		err = jsonOpts.applyToSchema(propertySchema)
		if err != nil {
			return nil, err
		}
		schema.Properties[jsonName] = propertySchema
	}
	return schema, nil
}

func jsonNameOf(field reflect.StructField) string {
	jsonTagValue, found := field.Tag.Lookup("json")
	if found && jsonTagValue != "" {
		return strings.Split(jsonTagValue, ",")[0]
	}
	return field.Name
}
