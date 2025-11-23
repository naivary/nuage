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

type testResponse struct {
	Scores       []int             `json:"scores"              maximum:"20" default:"musti"`
	PlayerName   string            `json:"playerName,omitzero"              minLength:"10" dependentRequired:"scores" deprecated:"true"`
	JerseyOwners map[string]string `json:"jerseyOwner"                                                                                minProperties:"1" maxProperties:"10"`
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
	handler := HandlerFuncErr[testRequest, testResponse](func(r *http.Request, input *testRequest) (*testResponse, error) {
		return nil, nil
	})
	err = Handle(api, &Operation{
		Description:         "something",
		Pattern:             "GET /path/to/handler",
		OperationID:         "test-operation-id",
		ResponseStatusCode:  http.StatusOK,
		ResponseDesc:        "something",
		ResponseContentType: ContentTypeJSON,
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
