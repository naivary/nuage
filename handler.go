package nuage

import (
	"encoding/json"
	"net/http"
)

type Decoder interface {
	Decode(r *http.Request) error
}

var _ Decoder = (*request)(nil)

// used to satisfy Decoder and to be used for compile-time
// check that the Handler is implementing http.Handler.
type request struct{}

func (r request) Decode(req *http.Request) error {
	return nil
}

var _ http.Handler = (HandlerFuncErr[request, struct{}])(nil)

// HandlerFuncErr represents the primary request handler function signature
// used by the framework to implement REST API endpoints.
type HandlerFuncErr[Request Decoder, Response any] func(ctx *Context, r Request) (Response, error)

func (hl HandlerFuncErr[RequestModel, ResponseModel]) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
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

func Handle[RequestModel Decoder, ResponseModel any](
	n *Nuage,
	hl HandlerFuncErr[RequestModel, RequestModel],
	op *Operation,
) error {
	return nil
}
