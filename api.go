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

const ContentTypeJSON = "application/json"

type APIConfig struct {
	LoggerOpts *slog.HandlerOptions

	// Supported formats of the REST API. If the format cannot be found an error
	// will be returned and not format negotiation can be succesfully completed.
	Formats map[string]Formater

	DefaultFormat string
}

func DefaultAPIConfig() *APIConfig {
	return &APIConfig{
		LoggerOpts: &slog.HandlerOptions{
			AddSource: true,
		},
		Formats: map[string]Formater{
			ContentTypeJSON: &JSONFormater{},
		},
		DefaultFormat: ContentTypeJSON,
	}
}

type api struct {
	doc        *openapi.OpenAPI
	mux        *http.ServeMux
	logger     *slog.Logger
	operations map[string]struct{}
	formats    map[string]Formater
}

func NewAPI(doc *openapi.OpenAPI, cfg *APIConfig) (*api, error) {
	if cfg == nil {
		cfg = DefaultAPIConfig()
	}
	if doc == nil {
		return nil, errors.New("new api: `doc` cannot be nil")
	}
	return &api{
		doc:        doc,
		mux:        http.NewServeMux(),
		operations: make(map[string]struct{}),
		logger:     slog.New(slog.NewJSONHandler(os.Stdout, cfg.LoggerOpts)),
		formats:    cfg.Formats,
	}, nil
}

func Handle[I, O any](api *api, op *openapi.Operation, handler HandlerFuncErr[I, O]) error {
	if !isStruct[I]() || !isStruct[O]() {
		return errors.New("handle: both input and output data types have to be of kind struct")
	}
	method, uri, isValidPatternSyntax := strings.Cut(op.Pattern, " ")
	if !isValidPatternSyntax {
		return fmt.Errorf("handle: invalid pattern syntax `%s`. Make sure to use the standard library syntax of [METHOD ][HOST]/[PATH]", op.Pattern)
	}
	if err := isValidOperation(api, op); err != nil {
		return err
	}
	api.operations[op.OperationID] = struct{}{}
	if err := buildOperationSpec[I, O](op); err != nil {
		return err
	}
	e := &endpoint[I, O]{
		handler: handler,
		doc:     op,
		logger:  api.logger,
		formats: api.formats,
	}
	pathItem := api.doc.Paths[op.Pattern]
	if pathItem == nil {
		pathItem = &openapi.PathItem{}
	}
	if err := pathItem.AddOperation(method, op); err != nil {
		return err
	}
	api.doc.Paths[uri] = pathItem
	api.mux.Handle(op.Pattern, e)
	return nil
}

func isValidOperation(api *api, op *openapi.Operation) error {
	if op.OperationID == "" {
		return fmt.Errorf("operation validation: missing operation id")
	}
	if _, isIDUnique := api.operations[op.OperationID]; isIDUnique {
		return fmt.Errorf("operation validation: repeated operation id `%s`", op.OperationID)
	}
	if op.Pattern == "" {
		return fmt.Errorf("operation validation: missing pattern for operation with id `%s`", op.OperationID)
	}
	return nil
}
