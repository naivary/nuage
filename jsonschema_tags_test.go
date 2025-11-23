package nuage

import (
	"reflect"
	"testing"
)

func TestParseJSONSchemaTagOpts(t *testing.T) {
	tests := []struct {
		field   reflect.StructField
		isValid func(got *jsonSchemaTagOpts) bool
	}{
		// --- metadata ---
		{
			field: reflect.StructField{Tag: `default:"123"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return string(got.dflt) == `"123"`
			},
		},
		{
			field: reflect.StructField{Tag: `deprecated:"true"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.deprecated
			},
		},
		{
			field: reflect.StructField{Tag: `readOnly:"true"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.readOnly
			},
		},
		{
			field: reflect.StructField{Tag: `writeOnly:"true"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.writeOnly
			},
		},

		// --- enum ---
		{
			field: reflect.StructField{Tag: `enum:"a,b,c"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return len(got.enum) == 3 &&
					got.enum[0] == "a" &&
					got.enum[1] == "b" &&
					got.enum[2] == "c"
			},
		},

		// --- numeric validation ---
		{
			field: reflect.StructField{Tag: `multipleOf:"2.5"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.multipleOf != nil && *got.multipleOf == 2.5
			},
		},
		{
			field: reflect.StructField{Tag: `minimum:"1"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.minimum != nil && *got.minimum == 1
			},
		},
		{
			field: reflect.StructField{Tag: `maximum:"10"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.maximum != nil && *got.maximum == 10
			},
		},
		{
			field: reflect.StructField{Tag: `exclusiveMinimum:"5"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.exclusiveMinimum != nil && *got.exclusiveMinimum == 5
			},
		},
		{
			field: reflect.StructField{Tag: `exclusiveMaximum:"20"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.exclusiveMaximum != nil && *got.exclusiveMaximum == 20
			},
		},
		{
			field: reflect.StructField{Tag: `minLength:"3"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.minLength != nil && *got.minLength == 3
			},
		},
		{
			field: reflect.StructField{Tag: `maxLength:"50"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.maxLength != nil && *got.maxLength == 50
			},
		},
		{
			field: reflect.StructField{Tag: `pattern:"^[a-z]+$"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.pattern == "^[a-z]+$"
			},
		},

		// --- array validation ---
		{
			field: reflect.StructField{Tag: `minItems:"1"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.minItems != nil && *got.minItems == 1
			},
		},
		{
			field: reflect.StructField{Tag: `maxItems:"5"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.maxItems != nil && *got.maxItems == 5
			},
		},
		{
			field: reflect.StructField{Tag: `uniqueItems:"true"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.uniqueItems
			},
		},
		{
			field: reflect.StructField{Tag: `minContains:"2"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.minContains != nil && *got.minContains == 2
			},
		},
		{
			field: reflect.StructField{Tag: `maxContains:"4"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.maxContains != nil && *got.maxContains == 4
			},
		},

		// --- object validation ---
		{
			field: reflect.StructField{Tag: `minProperties:"1"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.minProperties != nil && *got.minProperties == 1
			},
		},
		{
			field: reflect.StructField{Tag: `maxProperties:"12"`},
			isValid: func(got *jsonSchemaTagOpts) bool {
				return got.maxProperties != nil && *got.maxProperties == 12
			},
		},
		{
			field: reflect.StructField{Name: "a", Tag: `dependentRequired:"b,c"`},
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
		t.Run("", func(t *testing.T) {
			got, err := parseJSONSchemaTagOpts(tc.field)
			if err != nil {
				t.Errorf("err: %v", err)
			}
			if !tc.isValid(got) {
				t.Errorf("Unexpected JSON Options for tag %s", tc.field.Tag)
				t.Log(got)
			}
		})
	}
}
