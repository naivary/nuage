package nuage

import (
	"net/http"
)

type HandlerFuncErr[I, O any] func(r *http.Request, input I) (O, error)

type handler[I, O any] struct {
	fn HandlerFuncErr[I, O]
}

func (h handler[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
