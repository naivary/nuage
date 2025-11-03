package nuage

import "net/http"

type HandlerFuncErr[I, O any] func(r *http.Request, w http.ResponseWriter, input *I) (*O, error)

func (h HandlerFuncErr[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
