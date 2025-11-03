package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

const (
	_tagKeyHeader = "header"
	_tagKeyCookie = "cookie"
	_tagKeyPath   = "path"
	_tagKeyQuery  = "query"
)

var errTagNotFound = errors.New("tag not found")

type paramTagOpts struct {
	Required   bool
	Name       string
	Deprecated bool
	Style      Style
	Optional   bool
	Example    any
}

func parseTagOpts(tagKey string, field reflect.StructField) (*paramTagOpts, error) {
	// TODO(naivary): exampels hould be defined in the tag e.g. example or
	// example=. This make order irrelevant
	p := paramTagOpts{}
	tagValue, ok := field.Tag.Lookup(tagKey)
	if !ok {
		return nil, errTagNotFound
	}
	values := strings.Split(tagValue, ",")
	if len(values) == 0 {
		return nil, fmt.Errorf("empty tag(%s) for %v", tagKey, field)
	}
	// first element of the tag is always the name
	p.Name = values[0]
	if slices.Contains(values, "deprecated") {
		p.Deprecated = true
	}
	if slices.Contains(values, "optional") {
		p.Optional = true
	}
	if slices.Contains(values, "required") {
		p.Required = true
	}
	styles := []Style{
		StyleMatrix,
		StyleLabel,
		StyleSimple,
		StyleForm,
		StyleSpaceDelim,
		StylePipeDelim,
		StyleDeepObject,
		StyleCookie,
	}
	for _, style := range styles {
		if slices.Contains(values, style.String()) {
			p.Style = style
			break
		}
	}
	// example
	for _, value := range values {
		k, v, found := strings.Cut(value, "=")
		if !found {
			continue
		}
		switch k {
		case "example":
			p.Example = v
		}
	}
	return &p, nil
}

func paramsFor[I any]() ([]*Parameter, error) {
	path, err := pathParams[I]()
	if err != nil {
		return nil, err
	}
	header, err := headerParams[I]()
	if err != nil {
		return nil, err
	}
	query, err := queryParams[I]()
	if err != nil {
		return nil, err
	}
	cookie, err := cookieParams[I]()
	if err != nil {
		return nil, err
	}
	return slices.Concat(path, header, query, cookie), nil
}

func pathParams[I any]() ([]*Parameter, error) {
	s := reflect.TypeFor[I]()
	params := make([]*Parameter, 0, s.NumField())
	for i := range s.NumField() {
		field := s.Field(i)
		schema, err := jsonschema.ForType(field.Type, &jsonschema.ForOptions{})
		if err != nil {
			return nil, err
		}
		opts, err := parseTagOpts(_tagKeyPath, field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

		param := Parameter{
			ParamIn:    ParamInPath,
			Schema:     schema,
			Name:       opts.Name,
			Deprecated: opts.Deprecated,
			Style:      opts.Style,
			// Path Parameters are always required.
			Required: true,
			Example:  opts.Example,
		}
		params = append(params, &param)
	}
	return params, nil
}

func headerParams[I any]() ([]*Parameter, error) {
	s := reflect.TypeFor[I]()
	params := make([]*Parameter, 0, s.NumField())
	for i := range s.NumField() {
		field := s.Field(i)
		schema, err := jsonschema.ForType(field.Type, &jsonschema.ForOptions{})
		if err != nil {
			return nil, err
		}
		opts, err := parseTagOpts(_tagKeyHeader, field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

		// Header key must be canonical
		opts.Name = http.CanonicalHeaderKey(opts.Name)
		param := Parameter{
			ParamIn:    ParamInHeader,
			Schema:     schema,
			Name:       opts.Name,
			Deprecated: opts.Deprecated,
			Style:      opts.Style,
			Required:   opts.Required,
			Example:    opts.Example,
		}
		params = append(params, &param)
	}
	return params, nil
}

func queryParams[I any]() ([]*Parameter, error) {
	s := reflect.TypeFor[I]()
	params := make([]*Parameter, 0, s.NumField())
	for i := range s.NumField() {
		field := s.Field(i)
		schema, err := jsonschema.ForType(field.Type, &jsonschema.ForOptions{})
		if err != nil {
			return nil, err
		}
		opts, err := parseTagOpts(_tagKeyQuery, field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

		param := Parameter{
			ParamIn:    ParamInQuery,
			Schema:     schema,
			Name:       opts.Name,
			Deprecated: opts.Deprecated,
			Style:      opts.Style,
			Required:   opts.Required,
			Example:    opts.Example,
		}
		params = append(params, &param)
	}
	return params, nil
}

func cookieParams[I any]() ([]*Parameter, error) {
	s := reflect.TypeFor[I]()
	params := make([]*Parameter, 0, s.NumField())
	for i := range s.NumField() {
		field := s.Field(i)
		schema, err := jsonschema.ForType(field.Type, &jsonschema.ForOptions{})
		if err != nil {
			return nil, err
		}
		opts, err := parseTagOpts(_tagKeyCookie, field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

		param := Parameter{
			ParamIn:    ParamInCookie,
			Schema:     schema,
			Name:       opts.Name,
			Deprecated: opts.Deprecated,
			Style:      opts.Style,
			Required:   opts.Required,
			Example:    opts.Example,
		}
		params = append(params, &param)
	}
	return params, nil
}
