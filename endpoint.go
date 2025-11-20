package nuage

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/naivary/nuage/openapi"
)

const _headerKeyContentType = "Content-Type"

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

	paramDocs map[string]*openapi.Parameter
}

// use transformer model to add $schema to the response struct
func (e endpoint[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	format := r.Header.Get(_headerKeyContentType)
	formater, isSupportedFormat := e.formats[format]
	if !isSupportedFormat {
		// unssuported format
		e.logger.Error("format not supported", "format", format)
		return
	}
	var input I
	err := decodeParams(r, e.paramDocsMap(), &input)
	if err != nil {
		e.logger.Error(err.Error())
		return
	}
	if err := formater.Decode(r.Body, &input); err != nil {
		// bad request internal error of decoding format
		e.logger.Error(err.Error())
		return
	}
	resolver, err := e.doc.RequestBody.Content[format].Schema.Resolve(nil)
	if err != nil {
		// internal error cannot resolve schema
		e.logger.Error(err.Error())
	}
	m, err := structToMap(&input)
	if err != nil {
		// internal error: struct to map conversion failed
		e.logger.Error(err.Error())
	}
	err = resolver.Validate(m)
	if err != nil {
		e.logger.Error(err.Error())
	}
	_, err = e.handler(r, &input)
	if err != nil {
		// error from the handler ahs to be HTTPError as by RFC 9457
		e.logger.Error(err.Error())
	}
}

func (e *endpoint[I, O]) paramDocsMap() map[string]*openapi.Parameter {
	if e.paramDocs != nil {
		return e.paramDocs
	}
	m := make(map[string]*openapi.Parameter, len(e.doc.Parameters))
	for _, param := range e.doc.Parameters {
		m[param.Name] = param
	}
	e.paramDocs = m
	return m
}

func structToMap[S any](v *S) (map[string]any, error) {
	if !isStruct[S]() {
		return nil, fmt.Errorf("struct to map: type is not struct")
	}
	m := make(map[string]any)
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return m, json.Unmarshal(data, &m)
}
