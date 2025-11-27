package nuage

import (
	"net/http"

	"github.com/theory/jsonpath"
)

type openAPIDocQueryRequest struct {
	JSONPath string `query:"jsonpath"`
	URI      string `query:"uri"`
}

type openAPIDocQueryResponse struct {
	Result jsonpath.NodeList `json:"result"`
}

func queryOpenAPIDoc(q *openAPIQuerier) HandlerFuncErr[openAPIDocQueryRequest, *openAPIDocQueryResponse] {
	return HandlerFuncErr[openAPIDocQueryRequest, *openAPIDocQueryResponse](
		func(r *http.Request, input openAPIDocQueryRequest) (*openAPIDocQueryResponse, error) {
			nodes, err := q.Select(input.JSONPath)
			if err != nil {
				// bad request status code
				return nil, err
			}
			return &openAPIDocQueryResponse{Result: nodes}, nil
		},
	)
}
