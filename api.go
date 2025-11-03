package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type api struct {
	openAPI *OpenAPI

	mux *http.ServeMux
}

func NewAPI(root *OpenAPI) *api {
	return &api{
		openAPI: root,
		mux:     http.NewServeMux(),
	}
}

func Handle[I, O any](api *api, operation *Operation, handler HandlerFuncErr[I, O]) error {
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
	params, err := paramsFor[I]()
	if err != nil {
		return err
	}
	operation.Parameters = params

	api.openAPI.Paths[pattern] = &pathItem{}
	api.mux.Handle(operation.Pattern, handler)
	return nil
}

func paramsFor[I any]() ([]*Parameter, error) {
	path, err := pathParams[I]()
	if err != nil {
		return nil, err
	}
	header, err := headerParams[I]()
	if err != nil {
		return nil, err
	}
	query, err := queryParams[I]()
	if err != nil {
		return nil, err
	}
	cookie, err := cookieParams[I]()
	if err != nil {
		return nil, err
	}
	return slices.Concat(path, header, query, cookie), nil
}
