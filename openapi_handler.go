package nuage

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/theory/jsonpath"
)

type openAPIDocQueryRequest struct {
	JSONPath string `query:"jsonpath"`
	URI      string `query:"uri"`
}

type openAPIDocQueryResponse struct {
	Result jsonpath.NodeList `json:"result"`
}

func queryOpenAPIDoc(doc *OpenAPI) HandlerFuncErr[openAPIDocQueryRequest, *openAPIDocQueryResponse] {
	var (
		init           sync.Once
		openAPIJSON    any
		openAPIJSONErr error
	)
	// TODO: if uri is not empty get the path item of that uri
	return HandlerFuncErr[openAPIDocQueryRequest, *openAPIDocQueryResponse](
		func(r *http.Request, input openAPIDocQueryRequest) (*openAPIDocQueryResponse, error) {
			init.Do(func() {
				openAPIJSON, openAPIJSONErr = json.Marshal(doc)
			})
			if openAPIJSONErr != nil {
				return nil, openAPIJSONErr
			}
			p, err := jsonpath.Parse(input.JSONPath)
			if err != nil {
				return nil, err
			}
			nodes := p.Select(openAPIJSON)
			return &openAPIDocQueryResponse{Result: nodes}, nil
		},
	)
}
