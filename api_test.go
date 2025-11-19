package nuage

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naivary/nuage/openapi"
)

type requestTypeTest struct {
	Limit int `header:"Limit" json:"-"`

	PlayerName string `json:"playerName"`
}

type responseTypeTest struct{}

func TestAPIHandle(t *testing.T) {
	api, err := NewAPI(openapi.New("1.0.0", nil), DefaultAPIConfig())
	if err != nil {
		t.Fatalf("new api: %v", err)
	}
	hl := HandlerFuncErr[requestTypeTest, responseTypeTest](func(r *http.Request, input *requestTypeTest) (*responseTypeTest, error) {
		return nil, nil
	})
	err = Handle(api, &openapi.Operation{
		Pattern:     "POST /path/to/endpoint",
		OperationID: "CreateUser",
	}, hl)
	if err != nil {
		t.Fatalf("handle: %v", err)
	}

	w := httptest.NewRecorder()
	data, err := json.Marshal(&requestTypeTest{PlayerName: "nuage"})
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}
	r, err := http.NewRequest(http.MethodPost, "/path/to/endpoint", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	r.Header.Add(_headerKeyContentType, ContentTypeJSON)
	r.Header.Add("Limit", "3")
	api.mux.ServeHTTP(w, r)
}
