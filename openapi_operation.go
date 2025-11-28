package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/jsonschema-go/jsonschema"
)

type Operation struct {
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty"`
	Deprecated  bool                  `json:"deprecated,omitempty"`
	Security    []SecurityRequirement `json:"security,omitempty"`

	Parameters  []*Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody         `json:"requestBody,omitempty"`
	Responses   map[string]*Response `json:"responses,omitempty"`

	// nuage specific
	Pattern string `json:"-"`

	ResponseContentType string `json:"-"`
	ResponseDesc        string `json:"-"`
	ResponseStatusCode  int    `json:"-"`

	RequestContentType    string `json:"-"`
	RequestDesc           string `json:"-"`
	IsRequestBodyRequired *bool  `json:"-"`
}

func (o *Operation) IsValid() error {
	if o.Summary == "" && o.Description == "" {
		return errors.New("operation validation: one of summary or description has to be non-empty")
	}
	if o.OperationID == "" {
		return errors.New("operation validation: operation id undefined")
	}
	if o.Pattern == "" {
		return errors.New("operation validation: pattern undefined")
	}
	if o.RequestContentType != ContentTypeJSONPatch && o.RequestContentType != ContentTypeMergePatch && methodOf(o.Pattern) == http.MethodPatch {
		return fmt.Errorf(
			"operation validation: PATCH operations cannot have another Content-Type than %s or %s",
			ContentTypeMergePatch,
			ContentTypeJSONPatch,
		)
	}

	if o.RequestBody != nil && methodOf(o.Pattern) == http.MethodGet {
		return errors.New("operation validation: GET operations cannot have a request body")
	}
	if !isHTTPStatus(o.ResponseStatusCode) {
		return errors.New(
			"operation validation: response status code has to be between 100 and 599. For further information about HTTP status codes see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status",
		)
	}
	if o.ResponseDesc == "" {
		return errors.New("operation validation: missing response description")
	}

	return nil
}

func operationSpecFor[I, O any](op *Operation) error {
	paramSpecs, err := paramSpecsFor[I]()
	if err != nil {
		return err
	}
	op.Parameters = paramSpecs

	requestBody, err := requestBodyFor[I](op)
	if err != nil {
		return err
	}
	op.RequestBody = requestBody
	if op.IsRequestBodyRequired == nil {
		op.IsRequestBodyRequired = Ptr(true)
	}

	response, err := responseFor[O](op)
	if err != nil {
		return err
	}
	if op.Responses == nil {
		op.Responses = make(map[string]*Response, 1)
	}
	op.Responses[strconv.Itoa(op.ResponseStatusCode)] = response
	return nil
}

func requestBodyFor[I any](op *Operation) (*RequestBody, error) {
	method := methodOf(op.Pattern)
	if method == http.MethodGet {
		return nil, nil
	}
	if op.RequestBody != nil {
		return op.RequestBody, nil
	}
	reqBody := &RequestBody{
		Description: op.RequestDesc,
		Required:    op.IsRequestBodyRequired,
	}
	if !isJSONish(op.RequestContentType) {
		return reqBody, nil
	}
	if isEmptyJSON[I]() {
		return reqBody, nil
	}
	schema, err := jsonSchemaFor[I](nil)
	if err != nil {
		return nil, err
	}
	reqBody.Content[op.RequestContentType] = &MediaType{Schema: schema}
	return reqBody, nil
}

func responseFor[O any](op *Operation) (*Response, error) {
	headers, err := headersFor[O]()
	if err != nil {
		return nil, err
	}
	res := &Response{
		Description: op.ResponseDesc,
		Headers:     headers,
	}
	if isEmptyJSON[O]() || !isJSONish(op.ResponseContentType) {
		return res, nil
	}
	if res.Content == nil {
		res.Content = make(map[string]*MediaType, 1)
	}
	_, isUserDefined := res.Content[op.ResponseContentType]
	if isUserDefined {
		return res, nil
	}
	schema, err := jsonSchemaFor[O](nil)
	if err != nil {
		return nil, err
	}
	res.Content[op.ResponseContentType] = &MediaType{
		Schema: schema,
	}
	return res, nil
}

func headersFor[O any]() (map[string]*Parameter, error) {
	fields, err := fieldsOf[O]()
	if err != nil {
		return nil, err
	}
	headers := make(map[string]*Parameter, len(fields))
	for _, field := range fields {
		opts, err := parseParamTagOpts(field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if opts.tagKey != _tagKeyHeader {
			continue
		}
		schema, err := jsonschema.ForType(field.Type, &jsonschema.ForOptions{})
		if err != nil {
			return nil, err
		}
		headers[opts.name] = &Parameter{
			Style:      StyleSimple,
			Explode:    opts.explode,
			Deprecated: opts.deprecated,
			Schema:     schema,
		}
	}
	return headers, nil
}
