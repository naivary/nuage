package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	ContentTypeJSON       = "application/json"
	ContentTypeMergePatch = "application/merge-patch+json"
	ContentTypeJSONPatch  = "application/json-patch+json"
)

type APIConfig struct {
	DefaultContentType string
}

func DefaultAPIConfig() *APIConfig {
	return &APIConfig{
		DefaultContentType: ContentTypeJSON,
	}
}

type api struct {
	doc        *OpenAPI
	mux        *http.ServeMux
	operations map[string]struct{}

	cfg *APIConfig
}

func NewAPI(doc *OpenAPI, cfg *APIConfig) (*api, error) {
	if cfg == nil {
		cfg = DefaultAPIConfig()
	}
	if doc == nil || doc.Info == nil {
		return nil, errors.New(`api: the root OpenAPI documentation passed to NewAPI cannot be nil. It has to provide at least the info field`)
	}
	if doc.Info.Title == "" || doc.Info.Version == "" {
		return nil, errors.New("api: OpenAPI.Info.Title and OpenAPI.Info.Version are both required")
	}
	a := &api{
		doc:        doc,
		mux:        http.NewServeMux(),
		operations: make(map[string]struct{}, 1),
		cfg:        cfg,
	}
	return a, nil
}

func (a *api) Doc() *OpenAPI {
	return a.doc
}

func Handle[I, O any](api *api, op *Operation, fn HandlerFuncErr[I, O]) error {
	if !isStruct[I]() || !isStruct[O]() {
		return errors.New("handle: input and output type parameters have to be of kind struct")
	}
	if err := op.IsValid(); err != nil {
		return err
	}
	method, uri, isStdSyntax := strings.Cut(op.Pattern, " ")
	if !isStdSyntax {
		return fmt.Errorf("handle: invalid pattern syntax `%s`. Make sure to use the standard library syntax of [METHOD ][HOST]/[PATH]", op.Pattern)
	}
	if _, isIDUnique := api.operations[op.OperationID]; isIDUnique {
		return fmt.Errorf("handle: operation id repeated `%s`", op.OperationID)
	}
	if err := operationSpecFor[I, O](op); err != nil {
		return err
	}
	if !isPatternAndPathParamsConsistent(op.Pattern, op.Parameters) {
		return errors.New(
			`every path parameter in the route pattern must have a matching field in your input struct (using path:"name"), and every path:"name" field must correspond to a parameter in the route`,
		)
	}

	pathItem := api.doc.Paths[op.Pattern]
	if pathItem == nil {
		pathItem = &PathItem{}
	}
	if err := pathItem.AddOperation(method, op); err != nil {
		return err
	}
	api.operations[op.OperationID] = struct{}{}
	api.doc.Paths[uri] = pathItem
	handler := handler[I, O]{
		fn: fn,
	}
	api.mux.Handle(op.Pattern, &handler)
	return nil
}
