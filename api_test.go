package nuage

import (
	"encoding/json"
	"net/http"
	"testing"
)

type testRequest struct {
	P1 string `path:"p1" json:"-"`

	H1 string `header:"H1" json:"-"`

	Q1 string `query:"q1" json:"-"`

	C1 string `cookie:"c1" json:"-"`
}

func TestHandle(t *testing.T) {
	doc := NewOpenAPI(&Info{
		Version: "1.0",
		Title:   "nuage API Test",
	})
	api, err := NewAPI(doc, nil)
	if err != nil {
		t.Errorf("new api: %v", err)
	}
	handler := HandlerFuncErr[testRequest, struct{}](func(r *http.Request, input *testRequest) (*struct{}, error) {
		return nil, nil
	})
	err = Handle(api, &Operation{
		Description:        "something",
		Pattern:            "GET /path/to/handler",
		OperationID:        "test-operation-id",
		ResponseStatusCode: http.StatusOK,
		ResponseDesc:       "something",
	}, handler)
	if err != nil {
		t.Errorf("handle: %v", err)
	}
	data, err := json.Marshal(api.doc)
	if err != nil {
		t.Errorf("json: %v", err)
	}
	t.Logf("OpenAPI Doc: %s", data)
}
