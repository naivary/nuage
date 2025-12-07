package nuage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/theory/jsonpath"
)

type openAPIQuerier struct {
	doc  *OpenAPI
	json any
	mux  *http.ServeMux
}

func NewOpenAPIQuerier(doc *OpenAPI, mux *http.ServeMux) (*openAPIQuerier, error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	q := openAPIQuerier{
		doc: doc,
		mux: mux,
	}
	if err := json.Unmarshal(data, &q.json); err != nil {
		return nil, err
	}
	return &q, nil
}

func (o *openAPIQuerier) Select(jsonPath string, input any) (jsonpath.NodeList, error) {
	p, err := jsonpath.Parse(jsonPath)
	if err != nil {
		return nil, err
	}
	if input == nil {
		input = o.json
	}
	return p.Select(input), nil
}

func (o *openAPIQuerier) operationFor(r *http.Request) (*Operation, error) {
	_, pattern := o.mux.Handler(r)
	if pattern == "" {
		return nil, fmt.Errorf("openapi querier: request schema not found for %s", r.URL.RawPath)
	}
	// pathItem MUST exist
	return o.doc.Paths[pattern].OperationFor(r.Method), nil
}

func (o *openAPIQuerier) RequestBodySchemaFor(r *http.Request) (*jsonschema.Schema, error) {
	op, err := o.operationFor(r)
	if err != nil {
		return nil, err
	}
	contentType := r.Header.Get("Content-Type")
	mediaType, found := op.RequestBody.Content[contentType]
	if !found {
		return nil, fmt.Errorf("openapi querier: media type `%s` does not exist", contentType)
	}
	return mediaType.Schema, nil
}

func (o *openAPIQuerier) ResponseSchemaFor(r *http.Request, code int) (*jsonschema.Schema, error) {
	op, err := o.operationFor(r)
	if err != nil {
		return nil, err
	}
	codeAsText := strconv.Itoa(code)
	res, found := op.Responses[codeAsText]
	if !found {
		return nil, fmt.Errorf("openapi querier: no response for `%s` status code", codeAsText)
	}
	contentType := r.Header.Get("Content-Type")
	mediaType, found := res.Content[contentType]
	if !found {
		return nil, fmt.Errorf("openapi querier: media type `%s` does not exist", contentType)
	}
	return mediaType.Schema, nil
}
