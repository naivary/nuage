package nuage

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/naivary/nuage/openapi"
)

type Responser interface {
	StatusCode() int
	Description() string
}

type HandlerFuncErr[I, O any] func(r *http.Request, input *I) (*O, error)

type endpoint[I, O any] struct {
	handler HandlerFuncErr[I, O]
	doc     *openapi.Operation
	logger  *slog.Logger
	formats map[string]Formater
}

func (e endpoint[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	format := r.Header.Get(_headerKeyContentType)
	formater, isSupportedFormat := e.formats[format]
	if !isSupportedFormat {
		// unssuported format
		e.logger.Error("format not suypported", "format", format)
		return
	}
	err := e.validateParams(r)
	if err != nil {
		e.logger.Error(err.Error())
		return
	}
	var input I
	err = decodeParams(r, &input)
	if err != nil {
		e.logger.Error(err.Error())
		return
	}
	if err := formater.Decode(r.Body, &input); err != nil {
		// bad request internal error of decoding format
		e.logger.Error(err.Error())
		return
	}
	err = e.validateRequestBody(input)
	if err != nil {
		return
	}
	fmt.Println(input)
	validator, canValidate := any(&input).(Validator)
	if canValidate {
		errs := validator.Validate(ctx)
		if len(errs) > 0 {
			// send error message
			return
		}
	}
	res, err := e.handler(r, &input)
	if err == nil {
		responser := any(res).(Responser)
		err = encode(w, responser.StatusCode(), res)
		if err != nil {
			e.logger.Error(err.Error())
		}
		return
	}
	// error was returning by the handler and it should be a httpError to
	// convey RFC9457
	_, isHTTPErr := err.(httpError)
	if !isHTTPErr {
		// non-rfc9457 errors will only be logged and not retunred to the client
		// because of security risks of exposing internal functionalities
		e.logger.Error(err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(err)
	if err != nil {
		e.logger.Error(err.Error())
	}
}

func (e endpoint[I, O]) validateParams(r *http.Request) error {
	for _, param := range e.doc.Parameters {
		var value any
		switch param.ParamIn {
		case openapi.ParamInHeader:
			value = r.Header.Get(param.Name)
		case openapi.ParamInQuery:
			value = 3
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
