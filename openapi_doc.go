package nuage

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/google/jsonschema-go/jsonschema"
)

func operationSpecFor[I, O any](op *Operation, schemaReg map[reflect.Type]*jsonschema.Schema) error {
	paramSpecs, err := paramSpecsFor[I]()
	if err != nil {
		return err
	}
	op.Parameters = paramSpecs

	requestSchema, err := jsonschema.For[I](nil)
	if err != nil {
		return err
	}
	schemaReg[reflect.TypeFor[I]()] = requestSchema
	_, customSchemaIsExisting := op.RequestBody.Content[op.ContentType]
	if !customSchemaIsExisting {
		op.RequestBody.Content[op.ContentType] = &MediaType{Schema: requestSchema}
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
