//go:generate go tool go-enum --marshal --nocomments
package nuage

import (
	"encoding/json"

	"github.com/theory/jsonpath"
)

// ENUM(request, response)
type SchemaType string

type openAPIQuerier struct {
	doc *OpenAPI

	json any
}

func NewOpenAPIQuerier(doc *OpenAPI) (*openAPIQuerier, error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	q := openAPIQuerier{
		doc: doc,
	}
	if err := json.Unmarshal(data, &q.json); err != nil {
		return nil, err
	}
	return &q, nil
}

func (o *openAPIQuerier) Select(jsonPath string) (jsonpath.NodeList, error) {
	p, err := jsonpath.Parse(jsonPath)
	if err != nil {
		return nil, err
	}
	nodes := p.Select(o.json)
	return nodes, nil
}
