package nuage

import (
	"context"
	"encoding/json"
	"net/http"
)

type HandlerFuncErr[I, O any] func(r *http.Request, w http.ResponseWriter, input *I) (*O, error)

func (h HandlerFuncErr[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input I
	err := Decode(r, &input)
	if err != nil {
	}
	validator, ok := any(&input).(Validator)
	if ok {
		errs := validator.Validate(ctx)
		if len(errs) > 0 {
			// validation of input failed
			return
		}
	}

	output, err := h(r, w, &input)
	if err == nil {
		return
	}
	_, isHTTPErr := err.(httpError)
	if isHTTPErr {
		err = json.NewEncoder(w).Encode(err)
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
	}
}

func handleInput[I any](ctx context.Context, input *I) error {
	return nil
}
