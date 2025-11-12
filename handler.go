package nuage

import (
	"encoding/json"
	"net/http"
)

type HandlerFuncErr[I, O any] func(r *http.Request, w http.ResponseWriter, input *I) (*O, error)

func (h HandlerFuncErr[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input I
	err := Decode(r, &input)
	if err != nil {
	}

	output, err := h(r, w, &input)
	if err == nil {
		return
	}
	_, isHttpErr := err.(*HTTPError)
	if !isHttpErr {
		// Force the usage of the HTTPError type to make every API RFC 9457
		// compatible
	}
	// check if its RFC 9471 error type
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
	}
}
