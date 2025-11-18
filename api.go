package nuage

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"

	"github.com/naivary/nuage/openapi"
)

type Operation struct {
	*openapi.Operation

	Pattern string
}

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

func Handle[I, O any](api *api, op *Operation, handler HandlerFuncErr[I, O]) error {
	if !isStruct[I]() {
		return errors.New("non struct input type")
	}
	if !isStruct[O]() {
		return errors.New("non struct output type")
	}
	method, pattern, isValidPatternSyntax := strings.Cut(op.Pattern, " ")
	if !isValidPatternSyntax {
		return fmt.Errorf("invalid pattern syntax: %s", op.Pattern)
	}
	if op.OperationID == "" {
		return fmt.Errorf("handle: operation id missing")
	}
	if _, isIDExisting := api.operations[op.OperationID]; isIDExisting {
		return fmt.Errorf("handle: operation id repeated `%s`", op.OperationID)
	}
	api.operations[op.OperationID] = struct{}{}
	if err := buildOperationSpec[I, O](op.Operation); err != nil {
		return err
	}
	e := &endpoint[I, O]{
		handler: handler,
		doc:     op.Operation,
		logger:  api.logger,
	}
	pathItem := api.Doc.Paths[pattern]
	if pathItem == nil {
		pathItem = &openapi.PathItem{}
	}
	if err := pathItem.AddOperation(method, op.Operation); err != nil {
		return err
	}
	api.Doc.Paths[pattern] = pathItem
	api.Mux.Handle(op.Pattern, e)
	return nil
}

func buildOperationSpec[I, O any](op *openapi.Operation) error {
	paramSpecs, err := paramSpecsFor[I]()
	if err != nil {
		return err
	}
	op.Parameters = paramSpecs

	requestSchema, err := jsonschema.For[I](nil)
	if err != nil {
		return err
	}
	op.RequestBody = &openapi.RequestBody{
		Description: "Successfull Request!",
		Required:    true,
		Content: map[string]*openapi.MediaType{
			ContentTypeJSON: {Schema: requestSchema},
		},
	}
	responseSchema, err := jsonschema.For[O](nil)
	if err != nil {
		return err
	}
	return nil
}
