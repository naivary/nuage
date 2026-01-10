package nuage

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"

	"github.com/naivary/nuage/nuagetest"
)

func TestOpenAPIQuerier_SchemaFor(t *testing.T) {
	type testRequest struct {
		Foo    string
		Bar    int
		FooBar []string
	}

	mux := http.NewServeMux()
	mux.Handle("GET /foo/bar/{p1}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	schema, err := jsonschema.For[testRequest](nil)
	if err != nil {
		t.Errorf("OpenAPI Querier Schema For: %v", err)
	}
	q, err := NewOpenAPIQuerier(&OpenAPI{
		Paths: map[string]*PathItem{
			"/foo/bar/{p1}": {
				Get: &Operation{
					RequestBody: &RequestBody{
						Content: map[string]*MediaType{
							ContentTypeJSON: {
								Schema: schema,
							},
						},
					},
				},
			},
		},
	}, mux)
	if err != nil {
		t.Errorf("OpenAPI Querier Schema For: %v", err)
	}
	r := nuagetest.NewRequest(http.MethodGet, "foo/bar/foobar", nil)
	r.Header.Set("Content-Type", ContentTypeJSON)

	got, err := q.RequestSchemaOf(r)
	if err != nil {
		t.Errorf("OpenAPI Querier Schema For: %v", err)
	}

	gotData, err := json.Marshal(got)
	if err != nil {
		t.Errorf("OpenAPI Querier Schema For: %v", err)
	}
	wantData, err := json.Marshal(schema)
	if err != nil {
		t.Errorf("OpenAPI Querier Schema For: %v", err)
	}
	if !bytes.Equal(gotData, wantData) {
		t.Errorf("OpenAPI Querier Schema For: Schemas are not equal. Got: %s; Want: %s", gotData, wantData)
	}
}
