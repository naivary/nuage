package nuage

import (
	"net/http"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"

	"github.com/naivary/nuage/nuagetest"
	"github.com/naivary/nuage/openapi"
)

type decodeParamsRequestTypeTest struct {
	P1 int `json:"-" path:"p1"`
}

func jsonSchemaFor[T any]() *jsonschema.Schema {
	schema, err := jsonschema.For[T](nil)
	if err != nil {
		panic(err)
	}
	return schema
}

func TestDecodeParams(t *testing.T) {
	tests := []struct {
		name      string
		r         *http.Request
		isValid   func(input *decodeParamsRequestTypeTest) bool
		isInvalid bool
	}{
		{
			name: "integer path parameter",
			r: func() *http.Request {
				r := nuagetest.NewRequest(http.MethodGet, "", nil)
				r.SetPathValue("p1", "3")
				return r
			}(),
			isValid: func(input *decodeParamsRequestTypeTest) bool {
				return input.P1 == 3
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var input decodeParamsRequestTypeTest
			err := decodeParams(tc.r, map[string]*openapi.Parameter{
				"p1": {Schema: jsonSchemaFor[int]()},
			}, &input)
			if err != nil && !tc.isInvalid {
				t.Errorf("unexpected err: %v", err)
			}
			if !tc.isValid(&input) {
				t.Errorf("input is in unexpected form: %v", input)
			}
			t.Log(input)
		})
	}
}
