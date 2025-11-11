package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"slices"

	"github.com/google/jsonschema-go/jsonschema"
)

const (
	_defaultPathParamStyle = StyleSimple
)

func paramsFor[I any]() ([]*Parameter, error) {
	s := reflect.TypeFor[I]()
	params := make([]*Parameter, 0, s.NumField())
	for i := range s.NumField() {
		field := s.Field(i)
		schema, err := jsonschema.ForType(field.Type, nil)
		if err != nil {
			return nil, err
		}
		for _, tagKey := range _tagKeys {
			opts, err := parseParamTagOpts(tagKey, field)
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

			// one field can only be one type of parameter.
			break
		}
	}
	return params, nil
}

func newPathParam(opts *paramTagOpts) (*Parameter, error) {
	if opts.Style == "" {
		opts.Style = StyleSimple
	}
	validStyles := []Style{StyleSimple, StyleLabel, StyleMatrix}
	if !slices.Contains(validStyles, opts.Style) {
		return nil, fmt.Errorf("invalid style: %s. Valid styles are: %v", opts.Style, validStyles)
	}
	return &Parameter{
		ParamIn:    ParamInPath,
		Name:       opts.Name,
		Deprecated: opts.Deprecated,
		Style:      opts.Style,
		Explode:    opts.Explode,
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
		// Headers are always style simple
		Style:    StyleSimple,
		Required: opts.Required,
		Example:  opts.Example,
	}, nil
}

func newQueryParam(opts *paramTagOpts) (*Parameter, error) {
	if opts.Style == "" {
		opts.Style = StyleForm
		opts.Explode = true
	}
	validStyles := []Style{StyleForm, StyleSpaceDelim, StylePipeDelim, StyleDeepObject}
	if !slices.Contains(validStyles, opts.Style) {
		return nil, fmt.Errorf("invalid style: %s. Valid styles are: %v", opts.Style, validStyles)
	}
	if opts.Style == StyleDeepObject {
		opts.Explode = true
	}
	return &Parameter{
		ParamIn:    ParamInQuery,
		Name:       opts.Name,
		Deprecated: opts.Deprecated,
		Style:      opts.Style,
		Required:   opts.Required,
		Example:    opts.Example,
		Explode:    opts.Explode,
	}, nil
}

func newCookieParam(opts *paramTagOpts) (*Parameter, error) {
	return &Parameter{
		ParamIn:    ParamInCookie,
		Name:       opts.Name,
		Deprecated: opts.Deprecated,
		Style:      StyleForm,
		Required:   opts.Required,
		Example:    opts.Example,
	}, nil
}
