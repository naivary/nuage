package nuage

import (
	"reflect"
	"testing"
)

func TestParseJSONSchemaTagOpts(t *testing.T) {
	tests := []struct {
		name    string
		field   reflect.StructField
		isValid func(got *jsonSchemaTagOpts) bool
	}{
		// metadata
		{
			name:  "default",
			field: reflect.StructField{Type: reflect.TypeFor[string](), Tag: `default:"123"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return string(got.dflt) == `"123"`
			},
		},
		{
			name:  "deprecated",
			field: reflect.StructField{Type: reflect.TypeFor[string](), Tag: `deprecated:"true"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.deprecated
			},
		},
		{
			name:  "readOnly",
			field: reflect.StructField{Type: reflect.TypeFor[string](), Tag: `readOnly:"true"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.readOnly
			},
		},
		{
			name:  "writeOnly",
			field: reflect.StructField{Type: reflect.TypeFor[string](), Tag: `writeOnly:"true"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.writeOnly
			},
		},

		// enum
		{
			name:  "enum",
			field: reflect.StructField{Type: reflect.TypeFor[string](), Tag: `enum:"a,b,c"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				if len(got.enum) != 3 {
					return false
				}
				return got.enum[0] == "a" && got.enum[1] == "b" && got.enum[2] == "c"
			},
		},

		// numeric validation
		{
			name:  "multipleOf",
			field: reflect.StructField{Type: reflect.TypeFor[float64](), Tag: `multipleOf:"2.5"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.multipleOf != nil && *got.multipleOf == 2.5
			},
		},
		{
			name:  "minimum",
			field: reflect.StructField{Type: reflect.TypeFor[float64](), Tag: `minimum:"1"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.minimum != nil && *got.minimum == 1
			},
		},
		{
			name:  "maximum",
			field: reflect.StructField{Type: reflect.TypeFor[float64](), Tag: `maximum:"10"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.maximum != nil && *got.maximum == 10
			},
		},
		{
			name:  "exclusiveMinimum",
			field: reflect.StructField{Type: reflect.TypeFor[float64](), Tag: `exclusiveMinimum:"5"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.exclusiveMinimum != nil && *got.exclusiveMinimum == 5
			},
		},
		{
			name:  "exclusiveMaximum",
			field: reflect.StructField{Type: reflect.TypeFor[float64](), Tag: `exclusiveMaximum:"20"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.exclusiveMaximum != nil && *got.exclusiveMaximum == 20
			},
		},
		// string
		{
			name:  "minLength",
			field: reflect.StructField{Type: reflect.TypeFor[string](), Tag: `minLength:"3"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.minLength != nil && *got.minLength == 3
			},
		},
		{
			name:  "maxLength",
			field: reflect.StructField{Type: reflect.TypeFor[string](), Tag: `maxLength:"50"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.maxLength != nil && *got.maxLength == 50
			},
		},
		{
			name:  "pattern",
			field: reflect.StructField{Type: reflect.TypeFor[string](), Tag: `pattern:"^[a-z]+$"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.pattern == "^[a-z]+$"
			},
		},

		// array validation
		{
			name:  "minItems",
			field: reflect.StructField{Type: reflect.TypeFor[[]string](), Tag: `minItems:"1"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.minItems != nil && *got.minItems == 1
			},
		},
		{
			name:  "maxItems",
			field: reflect.StructField{Type: reflect.TypeFor[[]string](), Tag: `maxItems:"5"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.maxItems != nil && *got.maxItems == 5
			},
		},
		{
			name:  "uniqueItems",
			field: reflect.StructField{Type: reflect.TypeFor[[]string](), Tag: `uniqueItems:"true"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.uniqueItems
			},
		},
		{
			name:  "minContains",
			field: reflect.StructField{Type: reflect.TypeFor[[]string](), Tag: `minContains:"2"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.minContains != nil && *got.minContains == 2
			},
		},
		{
			name:  "maxContains",
			field: reflect.StructField{Type: reflect.TypeFor[[]string](), Tag: `maxContains:"4"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.maxContains != nil && *got.maxContains == 4
			},
		},

		// object validation
		{
			name:  "minProperties",
			field: reflect.StructField{Type: reflect.TypeFor[map[string]string](), Tag: `minProperties:"1"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.minProperties != nil && *got.minProperties == 1
			},
		},
		{
			name:  "maxProperties",
			field: reflect.StructField{Type: reflect.TypeFor[map[string]string](), Tag: `maxProperties:"12"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.maxProperties != nil && *got.maxProperties == 12
			},
		},
		{
			name:  "dependentRequired",
			field: reflect.StructField{Type: reflect.TypeFor[map[string]string](), Name: "a", Tag: `dependentRequired:"b,c"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				if len(got.dependentRequired) != 1 {
					return false
				}
				properties := got.dependentRequired["a"]
				if len(properties) != 2 {
					return false
				}
				return properties[0] == "b" && properties[1] == "c"
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseJSONSchemaTagOpts(tc.field)
			if err != nil {
				t.Errorf("err: %v", err)
			}
			if !tc.isValid(got) {
				t.Errorf("invalid %s", tc.field.Tag)
			}
		})
	}
}
