package nuage

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

type APIConfig struct {
	LoggerOpts *slog.HandlerOptions

	// Supported formats of the REST API. If the format cannot be found an error
	// will be returned and not format negotiation can be succesfully completed.
	Formats map[string]Formater

	DefaultContentType string
}

func DefaultAPIConfig() *APIConfig {
	return &APIConfig{
		LoggerOpts: &slog.HandlerOptions{
			AddSource: true,
		},
		Formats:            map[string]Formater{},
		DefaultContentType: ContentTypeJSON,
	}
}

type api struct {
	doc        *OpenAPI
	mux        *http.ServeMux
	logger     *slog.Logger
	operations map[string]struct{}
	schemaReg  map[reflect.Type]*jsonschema.Schema

	cfg *APIConfig
}

func NewAPI(doc *OpenAPI, cfg *APIConfig) (*api, error) {
	if cfg == nil {
		cfg = DefaultAPIConfig()
	}
	if doc == nil {
		return nil, errors.New("new api: `doc` cannot be nil")
	}
	a := &api{
		doc:        doc,
		mux:        http.NewServeMux(),
		operations: make(map[string]struct{}, 1),
		logger:     slog.New(slog.NewJSONHandler(os.Stdout, cfg.LoggerOpts)),
		cfg:        cfg,
	}

	// add defaul formaters
	maps.Copy(cfg.Formats, map[string]Formater{
		ContentTypeJSON: &jsonFormater{},
	})
	return a, nil
}

// TODO: check if request input struct has path parameters which are defined in
// the path also in the pattern.
// TODO: automatically set contentype of patch and put to RFC standard type. And provide utilities to check if the handler is idempotent
func Handle[I, O any](api *api, op *Operation, handler HandlerFuncErr[I, O]) error {
	if !isStruct[I]() || !isStruct[O]() {
		return errors.New("handle: both input and output data types have to be of kind struct")
	}
	method, uri, isValidPatternSyntax := strings.Cut(op.Pattern, " ")
	if !isValidPatternSyntax {
		return fmt.Errorf("handle: invalid pattern syntax `%s`. Make sure to use the standard library syntax of [METHOD ][HOST]/[PATH]", op.Pattern)
	}
	if op.ContentType == "" {
		op.ContentType = api.cfg.DefaultContentType
	}
	if method == http.MethodPatch && isPatchContentTypeRFCCompatible(op.ContentType) {
		return fmt.Errorf("handle: patch operation is not conveying to RFC 7386 or 6902 as content type")
	}
	if err := isValidOperation(api, op); err != nil {
		return err
	}
	if err := operationSpecFor[I, O](op, api.schemaReg); err != nil {
		return err
	}
	e := &endpoint[I, O]{
		handler: handler,
		doc:     op,
		logger:  api.logger,
		formats: api.cfg.Formats,
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
	api.mux.Handle(op.Pattern, e)
	return nil
}

func isValidOperation(api *api, op *Operation) error {
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
