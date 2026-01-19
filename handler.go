package nuage

import (
	"encoding/json"
	"net/http"
)

var _ http.Handler = (HandlerFuncErr[struct{}, struct{}])(nil)

// HandlerFuncErr represents the primary request handler function signature
// used by the framework to implement REST API endpoints.
type HandlerFuncErr[RequestModel, ResponseModel any] func(ctx *Context, r RequestModel) (ResponseModel, error)

func (hl HandlerFuncErr[RequestModel, ResponseModel]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req RequestModel
	ctx := NewCtx()
	res, err := hl(ctx, req)
	if err != nil {
		// handle error
		return
	}
	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		// handle err 
		return
	}
}

func Handle[RequestModel, ResponseModel any](
	n *Nuage,
	hl HandlerFuncErr[RequestModel, RequestModel],
	op *Operation,
) error {
	return nil
}
