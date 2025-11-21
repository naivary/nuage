package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const ContentTypeJSON = "application/json"

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
		return nil, errors.New(`new api: the root OpenAPI documentation passed to NewAPI cannot be nil. It has to provide at least the info field`)
	}
	a := &api{
		doc:        doc,
		mux:        http.NewServeMux(),
		operations: make(map[string]struct{}, 1),
		cfg:        cfg,
	}
	return a, nil
}

// TODO: check if request input struct has path parameters which are defined in
// the path also in the pattern.
func Handle[I, O any](api *api, op *Operation, handler HandlerFuncErr[I, O]) error {
	if !isStruct[I]() || !isStruct[O]() {
		return errors.New("handle: both input and output data types have to be of kind struct")
	}
	if err := op.IsValid(); err != nil {
		return err
	}
	method, uri, isStdSyntax := strings.Cut(op.Pattern, " ")
	if !isStdSyntax {
		return fmt.Errorf("handle: invalid pattern syntax `%s`. Make sure to use the standard library syntax of [METHOD ][HOST]/[PATH]", op.Pattern)
	}
	if _, ok := api.operations[op.OperationID]; ok {
		return fmt.Errorf("handle: operation id repeated `%s`", op.OperationID)
	}
	if err := operationSpecFor[I, O](op); err != nil {
		return err
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
	return nil
}
