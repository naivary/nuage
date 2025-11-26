package nuage

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/naivary/nuage/nuagetest"
)

func TestQueryOpenAPIDoc(t *testing.T) {
	tests := []struct {
		name  string
		input openAPIDocQueryRequest
	}{
		{
			name: "get title",
			input: openAPIDocQueryRequest{
				JSONPath: "$.info.title",
			},
		},
	}

	doc := &OpenAPI{
		Info: &Info{
			Version: "v1.0.0",
			Title:   "test-title",
		},
	}

	q, err := NewOpenAPIQuerier(doc)
	if err != nil {
		t.Errorf("OpenAPI Querier: %v", err)
	}
	hl := queryOpenAPIDoc(q)
	r := nuagetest.NewRequest(http.MethodGet, "", nil)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := hl(r, tc.input)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			data, err := json.Marshal(res)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			t.Logf("%s", string(data))
		})
	}
}
