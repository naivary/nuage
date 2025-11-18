package nuage

import (
	"net/http"
	"strconv"

	"github.com/google/jsonschema-go/jsonschema"

	"github.com/naivary/nuage/openapi"
)

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
	var output O
	responser, isResponser := any(output).(Responser)
	responseDesc := "Successfull Request!"
	responseStatusCode := http.StatusOK
	if isResponser {
		responseDesc = responser.Description()
		responseStatusCode = responser.StatusCode()
	}
	responseHeaders, err := responseHeaderSpecs[O]()
	if err != nil {
		return err
	}
	op.Responses[strconv.Itoa(responseStatusCode)] = &openapi.Response{
		Description: responseDesc,
		Headers: map[string]*openapi.Parameter{

		},
		Content: map[string]*openapi.MediaType{
			ContentTypeJSON: {
				Schema: responseSchema,
			},
		},
	}
	return nil
}

func responseHeaderSpecs[O any]() ([]*openapi.Parameter, error) {
	fields, err := fieldsOf[O]()
	if err != nil {
		return nil, err
	}
	headers := make([]*openapi.Parameter, 0, len(fields))
	for _, field := range fields {
		opts, err := parseParamTagOpts(field)
		if err != nil {
			return nil, err
		}
		if opts.tagKey != _tagKeyHeader {
			continue
		}
		headers = append(headers, &openapi.Parameter{
			Name:       opts.name,
			ParamIn:    openapi.ParamInHeader,
			Style:      openapi.StyleSimple,
			Explode:    opts.explode,
			Deprecated: opts.deprecated,
		})
	}
	return headers, nil
}
