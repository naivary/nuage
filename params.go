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
	opts := paramTagOpts{}
	tagValue, ok := field.Tag.Lookup(tagKey)
	if !ok {
		return nil, errTagNotFound
	}
	values := strings.Split(tagValue, ",")
	if len(values) == 0 {
		return nil, fmt.Errorf("empty tag(%s) for %v", tagKey, field)
	}
	// first element of the tag is always the name
	opts.Name = values[0]
	if slices.Contains(values, "deprecated") {
		opts.Deprecated = true
	}
	if slices.Contains(values, "optional") {
		opts.Optional = true
	}
	if slices.Contains(values, "required") {
		opts.Required = true
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
			opts.Style = style
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
			opts.Example = v
		}
	}
	return &opts, nil
}

func paramsFor[I any]() ([]*Parameter, error) {
	s := reflect.TypeFor[I]()
	params := make([]*Parameter, 0, s.NumField())
	for i := range s.NumField() {
		field := s.Field(i)
		schema, err := jsonschema.ForType(field.Type, &jsonschema.ForOptions{})
		if err != nil {
			return nil, err
		}
		tagKeys := []string{_tagKeyPath, _tagKeyHeader, _tagKeyQuery, _tagKeyCookie}
		for _, tagKey := range tagKeys {
			opts, err := parseTagOpts(tagKey, field)
			if errors.Is(err, errTagNotFound) {
				continue
			}
			if err != nil {
				return nil, err
			}

			var (
				param       *Parameter
				newParamErr error
			)
			switch tagKey {
			case _tagKeyPath:
				param, newParamErr = newPathParam(opts)
			case _tagKeyHeader:
				param, newParamErr = newHeaderParam(opts)
			case _tagKeyQuery:
				param, newParamErr = newQueryParam(opts)
			case _tagKeyCookie:
				param, newParamErr = newCookieParam(opts)
			}
			if newParamErr != nil {
				return nil, newParamErr
			}
			param.Schema = schema
			params = append(params, param)

			// only one parameter type is allowed per field.
			break
		}
	}
	return params, nil
}

func newPathParam(opts *paramTagOpts) (*Parameter, error) {
	return &Parameter{
		ParamIn:    ParamInPath,
		Name:       opts.Name,
		Deprecated: opts.Deprecated,
		Style:      opts.Style,
		// Path Parameters are always required.
		Required: true,
		Example:  opts.Example,
	}, nil
}

func newHeaderParam(opts *paramTagOpts) (*Parameter, error) {
	// Header key must be canonical
	canonicalName := http.CanonicalHeaderKey(opts.Name)
	if canonicalName != opts.Name {
		return nil, fmt.Errorf("header name is not canonical: %s. Change it to: %s", opts.Name, canonicalName)
	}
	return &Parameter{
		ParamIn:    ParamInHeader,
		Name:       canonicalName,
		Deprecated: opts.Deprecated,
		Style:      opts.Style,
		Required:   opts.Required,
		Example:    opts.Example,
	}, nil
}

func newQueryParam(opts *paramTagOpts) (*Parameter, error) {
	return &Parameter{
		ParamIn:    ParamInQuery,
		Name:       opts.Name,
		Deprecated: opts.Deprecated,
		Style:      opts.Style,
		Required:   opts.Required,
		Example:    opts.Example,
	}, nil
}

func newCookieParam(opts *paramTagOpts) (*Parameter, error) {
	return &Parameter{
		ParamIn:    ParamInCookie,
		Name:       opts.Name,
		Deprecated: opts.Deprecated,
		Style:      opts.Style,
		Required:   opts.Required,
		Example:    opts.Example,
	}, nil
}
