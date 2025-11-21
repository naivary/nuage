package nuage

import (
	"net/http"
	"strconv"

	"github.com/google/jsonschema-go/jsonschema"
)

func operationSpecFor[I, O any](op *Operation) error {
	paramSpecs, err := paramSpecsFor[I]()
	if err != nil {
		return err
	}
	op.Parameters = paramSpecs

	requestSchema, err := jsonschema.For[I](nil)
	if err != nil {
		return err
	}
	op.RequestBody = &RequestBody{
		Description: "Successfull Request!",
		Required:    true,
		Content: map[string]*MediaType{
			ContentTypeJSON: {Schema: requestSchema},
		},
	}
	responseSchema, err := jsonschema.For[O](nil)
	if err != nil {
		return err
	}

	var output O
	responser, isResponser := any(output).(Responser)
	responseDesc := "Successfull Request!"
	responseStatusCode := http.StatusOK
	if isResponser {
		responseDesc = responser.Description()
		responseStatusCode = responser.StatusCode()
	}
	responseHeaders, err := responseHeadersFor[O]()
	if err != nil {
		return err
	}
	if op.Responses == nil {
		op.Responses = make(map[string]*Response)
	}
	op.Responses[strconv.Itoa(responseStatusCode)] = &Response{
		Description: responseDesc,
		Headers:     responseHeaders,
		Content: map[string]*MediaType{
			ContentTypeJSON: {
				Schema: responseSchema,
			},
		},
	}
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
