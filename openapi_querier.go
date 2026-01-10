package nuage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

func (o *openAPIQuerier) RequestSchemaOf(r *http.Request) (*jsonschema.Schema, error) {
	_, uri := o.mux.Handler(r)
	_, pattern, _ := strings.Cut(uri, " ")
	item, ok := o.doc.Paths[pattern]
	if !ok {
		return nil, fmt.Errorf("request schema for: path item not found for %s", pattern)
	}
	contentType := r.Header.Get("Content-Type")
	// TODO: Get Operation dynamicall based on the method from strings.Cut
	mediaType, ok := item.Get.RequestBody.Content[contentType]
	if !ok {
		return nil, fmt.Errorf("requestSchemaFor: no request schema found for %s and content type %s", pattern, contentType)
	}
	return mediaType.Schema, nil
}
