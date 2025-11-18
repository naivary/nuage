package nuage

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/naivary/nuage/openapi"
)

type requestTypeTest struct{}

type responseTypeTest struct{}

func TestHandle(t *testing.T) {
	api, err := NewAPI(&APIConfig{
		LoggerOpts: &slog.HandlerOptions{},
		Doc:        openapi.New("1.0.0", nil),
	})
	if err != nil {
		t.Fatalf("new api: %v", err)
	}
	hl := HandlerFuncErr[requestTypeTest, responseTypeTest](func(r *http.Request, input *requestTypeTest) (*responseTypeTest, error) {
		return nil, nil
	})
	err = Handle(api, &openapi.Operation{
		Pattern:     "GET /path/to/endpoint",
		OperationID: "CreateUser",
	}, hl)
	if err != nil {
		t.Fatalf("handle: %v", err)
	}
	t.Log(json.NewEncoder(os.Stdout).Encode(api.doc.Paths))
}
