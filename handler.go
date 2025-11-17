package nuage

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/naivary/nuage/openapi"
)

type endpoint[I, O any] struct {
	handler HandlerFuncErr[I, O]
	doc     *openapi.Operation
}

func (e endpoint[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e.validateParams(r)
	if err != nil {
		return
	}
}

func (e endpoint[I, O]) validateParams(r *http.Request) error {
	for _, param := range e.doc.Parameters {
		var value string
		switch param.ParamIn {
		case openapi.ParamInHeader:
			value = r.Header.Get(param.Name)
		}
		resolver, err := param.Schema.Resolve(nil)
		if err != nil {
			return err
		}
		err = resolver.Validate(value)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: remove http.ResponseWriter from this
type HandlerFuncErr[I, O any] func(r *http.Request, w http.ResponseWriter, input *I) (*O, error)

func (h HandlerFuncErr[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input I
	ctx := r.Context()
	err := decode(r, &input)
	if err != nil {
		slog.Error(err.Error())
		return
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
		err = encode(w, http.StatusOK, output)
		if err != nil {
			slog.Error(err.Error())
		}
		return
	}
	// error was returning by the handler and it should be a httpError to
	// convey RFC9457
	_, isHTTPErr := err.(httpError)
	if !isHTTPErr {
		// non-rfc9457 errors will only be logged and not retunred to the client
		// because of security risks of exposing internal functionalities
		slog.Error(err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(err)
	if err != nil {
		slog.Error(err.Error())
	}
}
