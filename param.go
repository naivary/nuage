package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"slices"

	"github.com/google/jsonschema-go/jsonschema"
)

// ParamSpecsFor generates a list of OpenAPI parameter specifications for a given
// Go struct type `I`. It inspects the struct fields tags and derives parameter definitions
// based on recognized struct tags (such as `path`, `query`, `header`, or `cookie`).
//
// Each struct field can define at most one parameter type through its tag. The
// function automatically infers the JSON Schema for the field’s Go type and
// attaches it to the resulting parameter definition.
//
// Example usage:
//
//	type GetUserInput struct {
//	    UserID string `path:"user_id" json:"-"`
//	    Filter string `query:"filter,omitempty" json:"-"`
//	}
//
//	params, err := ParamSpecsFor[GetUserInput]()
//	if err != nil {
//	    log.Fatal(err)
//	}
func ParamSpecsFor[I any]() ([]*Parameter, error) {
	s := reflect.TypeFor[I]()
	params := make([]*Parameter, 0, s.NumField())
	for i := range s.NumField() {
		field := s.Field(i)
		schema, err := jsonschema.ForType(field.Type, &jsonschema.ForOptions{})
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
