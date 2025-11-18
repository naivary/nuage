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

type api struct {
	OpenAPI *openapi.OpenAPI

	Mux *http.ServeMux

	logger *slog.Logger

	operations map[string]struct{}
}

func NewAPI(cfg *APIConfig) *api {
	return &api{
		OpenAPI: cfg.Doc,
		Mux:     http.NewServeMux(),
		logger:  slog.New(slog.NewJSONHandler(os.Stdout, cfg.LoggerOpts)),
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

	api.OpenAPI.Paths[pattern] = &openapi.PathItem{}
	e := endpoint[I, O]{handler: handler, doc: operation}
	_, isExisting := api.operations[operation.OperationID]
	if isExisting {
		return fmt.Errorf("handle: non-unique operation id `%s`", operation.OperationID)
	}
	api.operations[operation.OperationID] = struct{}{}
	api.Mux.Handle(operation.Pattern, &e)
	return nil
}
