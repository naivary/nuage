package nuage

import (
	"errors"
	"net/http"
)

type Operation struct {
	// Tags associated with this operation
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty"`
	Deprecated  bool                  `json:"deprecated,omitempty"`
	Security    []SecurityRequirement `json:"security,omitempty"`

	Parameters  []*Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody         `json:"requestBody,omitempty"`
	Responses   map[string]*Response `json:"responses,omitempty"`

	// nuage specific members

	Pattern string `json:"-"`

	ResponseContentType string `json:"-"`
	ResponseDesc        string `json:"-"`
	ResponseStatusCode  int    `json:"-"`

	RequestContentType    string `json:"-"`
	RequestDesc           string `json:"-"`
	IsRequestBodyRequired bool   `json:"-"`
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
	if o.RequestBody != nil && methodOf(o.Pattern) == http.MethodGet {
		return errors.New("operation validation: GET operations cannot have a request body")
	}
	if !isStatusCodeInRange(o.ResponseStatusCode) {
		return errors.New(
			"operation validation: response status code has to be between 100 and 599. For further information about HTTP status codes see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status",
		)
	}
	if o.ResponseDesc == "" {
		return errors.New("operation validation: missing response description")
	}
	if !isStatusCodeInRange(o.ResponseStatusCode) {
		return errors.New(
			"operation validation: response status code has to be between 100 and 599. For further information about HTTP status codes see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status",
		)
	}

	return nil
}

// TODO: get operations are not allowed to have request body
func operationSpecFor[I, O any](op *Operation) error {
	paramSpecs, err := paramSpecsFor[I]()
	if err != nil {
		return err
	}
	op.Parameters = paramSpecs
	return nil
}

func responseHeadersFor[O any]() (map[string]*Parameter, error) {
	fields, err := fieldsOf[O]()
	if err != nil {
		return nil, err
	}
	headers := make(map[string]*Parameter, len(fields))
	for _, field := range fields {
		opts, err := parseParamTagOpts(field)
		if err != nil {
			return nil, err
		}
		if opts.tagKey != _tagKeyHeader {
			continue
		}
		headers[opts.name] = &Parameter{
			ParamIn:    ParamInHeader,
			Style:      StyleSimple,
			Explode:    opts.explode,
			Deprecated: opts.deprecated,
		}
	}
	return headers, nil
}
