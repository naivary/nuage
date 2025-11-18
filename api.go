package nuage

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/naivary/nuage/openapi"
)

type APIConfig struct {
	LoggerOpts *slog.HandlerOptions

	Doc *openapi.OpenAPI
}

func DefaultAPIConfig() *APIConfig {
	return &APIConfig{
		LoggerOpts: &slog.HandlerOptions{
			AddSource: true,
		},
	}
}

type api struct {
	Doc *openapi.OpenAPI

	Mux *http.ServeMux

	logger *slog.Logger

	operations map[string]struct{}
}

func NewAPI(cfg *APIConfig) (*api, error) {
	if cfg == nil {
		cfg = DefaultAPIConfig()
	}
	if cfg.Doc == nil {
		return nil, errors.New("new api: missing openapi documentation")
	}
	return &api{
		Doc:        cfg.Doc,
		Mux:        http.NewServeMux(),
		operations: make(map[string]struct{}),
		logger:     slog.New(slog.NewJSONHandler(os.Stdout, cfg.LoggerOpts)),
	}, nil
}

func Handle[I, O any](api *api, pattern string, op *openapi.Operation, handler HandlerFuncErr[I, O]) error {
	if !isStruct[I]() {
		return errors.New("non struct input type")
	}
	if !isStruct[O]() {
		return errors.New("non struct output type")
	}
	method, uri, isValidPatternSyntax := strings.Cut(pattern, " ")
	if !isValidPatternSyntax {
		return fmt.Errorf("invalid pattern syntax: %s", pattern)
	}
	if op.OperationID == "" {
		return fmt.Errorf("handle: operation id missing")
	}
	if _, isIDExisting := api.operations[op.OperationID]; isIDExisting {
		return fmt.Errorf("handle: operation id repeated `%s`", op.OperationID)
	}
	api.operations[op.OperationID] = struct{}{}
	if err := buildOperationSpec[I, O](op); err != nil {
		return err
	}
	e := &endpoint[I, O]{
		handler: handler,
		doc:     op,
		logger:  api.logger,
	}
	pathItem := api.Doc.Paths[pattern]
	if pathItem == nil {
		pathItem = &openapi.PathItem{}
	}
	if err := pathItem.AddOperation(method, op); err != nil {
		return err
	}
	api.Doc.Paths[uri] = pathItem
	api.Mux.Handle(pattern, e)
	return nil
}
