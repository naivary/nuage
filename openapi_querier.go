package nuage

import (
	"encoding/json"
	"net/http"

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
