package nuage

import (
	"net/http"
	"testing"
)

type request struct {
	verbose int   `query:"verbose"`
	id      int64 `path:"id"`
}

type response struct{}

func TestHandlerFuncErr(t *testing.T) {
	hl := HandlerFuncErr[request, response](func(r *http.Request, data *request) (*response, error) {
		return nil, nil
	})

	mux := http.NewServeMux()
	mux.Handle("GET /test/endpoint", hl)
}
