package nuage

import (
	"net/http"
	"testing"

	"github.com/naivary/nuage/nuagetest"
)

func TestOpenAPIQuerier_DocOf(t *testing.T) {
	doc := &OpenAPI{
		Info: &Info{
			Version: "v1.0.0",
			Title:   "test-title",
		},
	}
	mux := &http.ServeMux{}
	mux.Handle("/path/to/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	q, err := NewOpenAPIQuerier(doc, mux)
	if err != nil {
		t.Errorf("OpenAPI Querier: %v", err)
	}

	r := nuagetest.NewRequest(http.MethodGet, "/path/to/endpoint", nil)
	q.DocOf(r)
}
