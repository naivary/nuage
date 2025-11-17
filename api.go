package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/naivary/nuage/openapi"
)

type api struct {
	openAPI *openapi.OpenAPI

	mux *http.ServeMux
}

func NewAPI(root *openapi.OpenAPI) *api {
	return &api{
		openAPI: root,
		mux:     http.NewServeMux(),
	}
}

func Handle[I, O any](api *api, operation *openapi.Operation, handler HandlerFuncErr[I, O]) error {
	if !isStruct[I]() {
		return errors.New("non struct input type")
	}
	if !isStruct[O]() {
		return errors.New("non struct output type")
	}
	_, pattern, found := strings.Cut(operation.Pattern, " ")
	if !found {
		return fmt.Errorf("invalid pattern syntax: %s", operation.Pattern)
	}
	params, err := paramSpecsFor[I]()
	if err != nil {
		return err
	}
	operation.Parameters = params

	api.openAPI.Paths[pattern] = &openapi.PathItem{}
	api.mux.Handle(operation.Pattern, handler)
	return nil
}
