package nuage

import "net/http"

type HandlerFuncErr[I any, O any] func(r *http.Request, data *I) (*O, error)

func (h HandlerFuncErr[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data I
	res, err := h(r, &data)
	if err != nil {
		panic(err)
	}
}
