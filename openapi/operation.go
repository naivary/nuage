package openapi

import "github.com/google/jsonschema-go/jsonschema"

type Operation struct {
	// Tags associated with this operation
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty"`
	Deprecated  bool                  `json:"deprecated,omitempty"`
	Security    []SecurityRequirement `json:"security,omitempty"`
	Pattern     string                `json:"-"`

	Parameters  []*Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody         `json:"requestBody,omitempty"`
	Responses   map[string]*Response `json:"responses"`
}

func (o *Operation) GetParamSchema(name string, in ParamIn) *jsonschema.Schema {
	for _, param := range o.Parameters {
		if param.Name == name && param.ParamIn == in {
			return param.Schema
		}
	}
	return nil
}
