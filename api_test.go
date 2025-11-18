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
	hl := HandlerFuncErr[requestTypeTest, responseTypeTest](func(r *http.Request, input *requestTypeTest) (Response[responseTypeTest], error) {
		var output responseTypeTest
		return Return(output, http.StatusOK), nil
	})
	err = Handle(api, &Operation{
		Pattern:     "GET /path/to/{p1}",
		OperationID: "CreateUser",
	}, hl)
	if err != nil {
		t.Fatalf("handle: %v", err)
	}
	t.Log(json.NewEncoder(os.Stdout).Encode(api.Doc.Paths))
}
