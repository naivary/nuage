package nuage

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

type pathParamOpts struct {
	Name       string
	Deprecated bool
}

func newPathParamOpts(field reflect.StructField, tag string) (*pathParamOpts, error) {
	p := pathParamOpts{}
	values := strings.Split(tag, ",")
	if len(values) == 0 {
		return nil, fmt.Errorf("path tag cannot be empty: %s", field.Name)
	}
	if slices.Contains(values, "deprecated") {
		p.Deprecated = true
	}
	return &p, nil
}

func pathParams[I any]() ([]*Parameter, error) {
	s := reflect.TypeFor[I]()
	params := make([]*Parameter, 0, s.NumField())
	for i := range s.NumField() {
		field := s.Field(i)
		tag, ok := field.Tag.Lookup("path")
		if !ok {
			continue
		}
		schema, err := jsonschema.For[I](nil)
		if err != nil {
			return nil, err
		}
		opts, err := newPathParamOpts(field, tag)
		if err != nil {
			return nil, err
		}
		param := Parameter{
			Name:       opts.Name,
			ParamIn:    ParamInPath,
			Schema:     schema,
			Required:   true,
			Deprecated: opts.Deprecated,
		}
		params = append(params, &param)
	}
	return params, nil
}
