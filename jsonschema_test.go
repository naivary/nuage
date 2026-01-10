package nuage

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
)

func baseSchema(typ reflect.Type, opts *jsonschema.ForOptions) *jsonschema.Schema {
	schema, err := jsonschema.ForType(typ, opts)
	if err != nil {
		panic(err)
	}
	return schema
}

func TestJSONSchemaFor(t *testing.T) {
	tests := []struct {
		name string
		typ  any
		want func(typ reflect.Type, opts *jsonschema.ForOptions) *jsonschema.Schema
	}{
		// string options
		{
			name: "string",
			typ: struct {
				A string `json:"a" minLength:"1" maxLength:"2" pattern:"[a-z]+" deprecated:"true" readOnly:"true" writeOnly:"true" default:"hello" enum:"x,y,z"`
			}{},
			want: func(typ reflect.Type, opts *jsonschema.ForOptions) *jsonschema.Schema {
				base := baseSchema(typ, opts)
				p := base.Properties["a"]

				p.MinLength = jsonschema.Ptr(1)
				p.MaxLength = jsonschema.Ptr(2)
				p.Pattern = "[a-z]+"

				p.Deprecated = true
				p.ReadOnly = true
				p.WriteOnly = true

				p.Default = json.RawMessage(`"hello"`)
				p.Enum = []any{"x", "y", "z"}
				return base
			},
		},

		// int options
		{
			name: "integer",
			typ: struct {
				A int `json:"a" minimum:"1" maximum:"10" exclusiveMinimum:"2" exclusiveMaximum:"9" multipleOf:"2" default:"5" enum:"1,2,3"`
			}{},
			want: func(typ reflect.Type, opts *jsonschema.ForOptions) *jsonschema.Schema {
				base := baseSchema(typ, opts)
				p := base.Properties["a"]

				p.Minimum = jsonschema.Ptr(1.0)
				p.Maximum = jsonschema.Ptr(10.0)
				p.ExclusiveMinimum = jsonschema.Ptr(2.0)
				p.ExclusiveMaximum = jsonschema.Ptr(9.0)
				p.MultipleOf = jsonschema.Ptr(2.0)

				p.Default = json.RawMessage(`5`)
				p.Enum = []any{1, 2, 3}
				return base
			},
		},

		// float64 options
		{
			name: "number",
			typ: struct {
				A float64 `json:"a" minimum:"0.5" maximum:"99.9" exclusiveMinimum:"0.7" exclusiveMaximum:"88.8" multipleOf:"0.1" default:"1.23" enum:"0.1,0.2,0.3"`
			}{},
			want: func(typ reflect.Type, opts *jsonschema.ForOptions) *jsonschema.Schema {
				base := baseSchema(typ, opts)
				p := base.Properties["a"]

				p.Minimum = jsonschema.Ptr(0.5)
				p.Maximum = jsonschema.Ptr(99.9)
				p.ExclusiveMinimum = jsonschema.Ptr(0.7)
				p.ExclusiveMaximum = jsonschema.Ptr(88.8)
				p.MultipleOf = jsonschema.Ptr(0.1)

				p.Default = json.RawMessage(`1.23`)
				p.Enum = []any{0.1, 0.2, 0.3}
				return base
			},
		},

		// slice (array) options
		{
			name: "array",
			typ: struct {
				A []string `json:"a" minItems:"1" maxItems:"5" uniqueItems:"true" minContains:"1" maxContains:"2"`
			}{},
			want: func(typ reflect.Type, opts *jsonschema.ForOptions) *jsonschema.Schema {
				base := baseSchema(typ, opts)
				p := base.Properties["a"]

				p.MinItems = jsonschema.Ptr(1)
				p.MaxItems = jsonschema.Ptr(5)
				p.UniqueItems = true
				p.MinContains = jsonschema.Ptr(1)
				p.MaxContains = jsonschema.Ptr(2)
				return base
			},
		},

		// map/object options
		{
			name: "object",
			typ: struct {
				A map[string]int `json:"a,omitzero" minProperties:"1" maxProperties:"3" dependentRequired:"b,c" enum:"1,2,3"`
			}{},
			want: func(typ reflect.Type, opts *jsonschema.ForOptions) *jsonschema.Schema {
				base := baseSchema(typ, opts)
				p := base.Properties["a"]

				p.MinProperties = jsonschema.Ptr(1)
				p.MaxProperties = jsonschema.Ptr(3)
				p.AdditionalProperties.Enum = []any{1, 2, 3}

				base.DependentRequired = map[string][]string{
					"a": {"b", "c"},
				}
				return base
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			typ := reflect.TypeOf(tc.typ)
			got, err := jsonSchemaForType(typ, nil)
			if err != nil {
				t.Errorf("err: %v", err)
			}
			want := tc.want(typ, nil)

			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Errorf("json marshal: %v", err)
			}
			wantJSON, err := json.Marshal(want)
			if err != nil {
				t.Errorf("json marshal: %v", err)
			}
			if !bytes.Equal(gotJSON, wantJSON) {
				t.Errorf("Got: %s\n Want: %s", gotJSON, wantJSON)
			}
		})
	}
}
