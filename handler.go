package nuage

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/naivary/nuage/openapi"
)

// TODO: remove http.ResponseWriter from this
type HandlerFuncErr[I, O any] func(r *http.Request, w http.ResponseWriter, input *I) (*O, error)

type endpoint[I, O any] struct {
	handler HandlerFuncErr[I, O]
	doc     *openapi.Operation
	logger  *slog.Logger
}

func (e endpoint[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := e.validateParams(r)
	if err != nil {
		return
	}
	var input I
	err = decode(r, &input)
	if err != nil {
		return
	}
	err = e.validateRequestBody(input)
	if err != nil {
		return
	}
	validator, canValidate := any(&input).(Validator)
	if canValidate {
		errs := validator.Validate(ctx)
		if len(errs) > 0 {
			// send error message
			return
		}
	}
	output, err := e.handler(r, w, &input)
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

func (e endpoint[I, O]) validateParams(r *http.Request) error {
	for _, param := range e.doc.Parameters {
		var value string
		switch param.ParamIn {
		case openapi.ParamInHeader:
			value = r.Header.Get(param.Name)
		}
		if value == "" && param.Required {
			return fmt.Errorf("parameter validatin: missing required parameter `%s`", param.Name)
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

func (e endpoint[I, O]) validateRequestBody(input I) error {
	schema := e.doc.RequestBody.Content[ContentTypeJSON].Schema
	resolver, err := schema.Resolve(nil)
	if err != nil {
		return err
	}
	return resolver.Validate(&input)
}
