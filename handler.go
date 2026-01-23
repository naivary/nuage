package nuage

import (
	"encoding/json"
	"net/http"

	"github.com/naivary/nuage/openapi"
)

type Decoder interface {
	Decode(r *http.Request) error
}

var _ Decoder = (*request)(nil)

// used to satisfy Decoder and to be used for compile-time
// check that the Handler is implementing http.Handler.
type request struct{}

func (r request) Decode(req *http.Request) error { return nil }

var _ http.Handler = (HandlerFuncErr[request, struct{}])(nil)

// HandlerFuncErr is the primary request handler function signature used by the
// framework to implement REST API endpoints.
//
// # Request Model
//
// The Request type defines the structure of incoming requests. It may include
// parameters from the following sources:
//   - Path
//   - Query
//   - Header
//   - Cookie
//
// You define your own request model using supported Go types, as outlined
// below.
//
// # Supported Types by Parameter Source
//
// Path Parameters
//   - string
//   - int8, int16, int32, int64
//   - uint8, uint16, uint32, uint64
//   - float32, float64
//
// Query Parameters
//
//   - string
//
//   - int8, int16, int32, int64
//
//   - uint8, uint16, uint32, uint64
//
//   - float32, float64
//
//   - bool
//
//   - time.Duration
//
//   - time.Time (RFC3339 format only)
//
//     Additional support:
//
//   - []T (slice of supported types), except time.Time and time.Duration
//
//   - map[string]string when using query style `deepObject`
//
// Header Parameters
//
//   - string
//
//   - int8, int16, int32, int64
//
//   - uint8, uint16, uint32, uint64
//
//   - float32, float64
//
//   - bool
//
//   - time.Duration
//
//   - time.Time (RFC3339 format only)
//
//     Additional support:
//
//   - []T (slice of supported types), except time.Time and time.Duration
//
// Cookie Parameters
//   - *http.Cookie
//
// Notes
//   - Slice types ([]T) are supported for Path, Query, and Header parameters
//     unless explicitly stated otherwise.
//   - time.Time and time.Duration do not support slice variants.
//   - All time.Time values must be provided in RFC3339 format.
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
	op *openapi.Operation,
) error {
	return nil
}
