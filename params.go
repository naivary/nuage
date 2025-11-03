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

var errNoTag = errors.New("tag not found")

type paramTagOpts struct {
	Required   bool
	Name       string
	Deprecated bool
	Style      Style
	Optional   bool
	Example    any
}

func newParamTagOpts(tagKey string, field reflect.StructField) (*paramTagOpts, error) {
	p := paramTagOpts{}
	tag, ok := field.Tag.Lookup(tagKey)
	if !ok {
		return nil, errNoTag
	}
	values := strings.Split(tag, ",")
	if len(values) == 0 {
		return nil, fmt.Errorf("path tag cannot be empty: %s", field.Name)
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("the value has to be at least the length of two. First element is the name second is an example value")
	}
	// first element of the tag is always the name
	p.Name = values[0]
	// last element is always an example
	p.Example = values[len(values)-1]
	values = values[1 : len(values)-1]
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
		opts, err := newParamTagOpts(_tagKeyPath, field)
		if errors.Is(err, errNoTag) {
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
		opts, err := newParamTagOpts(_tagKeyHeader, field)
		if errors.Is(err, errNoTag) {
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
		opts, err := newParamTagOpts(_tagKeyQuery, field)
		if errors.Is(err, errNoTag) {
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
		opts, err := newParamTagOpts(_tagKeyCookie, field)
		if errors.Is(err, errNoTag) {
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
