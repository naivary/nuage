package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/google/jsonschema-go/jsonschema"
)

// paramSpecsFor generates a list of OpenAPI parameter specifications for a given
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
//	params, err := paramSpecsFor[GetUserInput]()
//	if err != nil {
//	    log.Fatal(err)
//	}
func paramSpecsFor[I any]() ([]*Parameter, error) {
	fields, err := fieldsOf[I]()
	params := make([]*Parameter, 0, len(fields))
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		schema, err := jsonschema.ForType(field.Type, &jsonschema.ForOptions{})
		if err != nil {
			return nil, err
		}
		opts, err := parseParamTagOpts(field)
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
		switch opts.tagKey {
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

	}
	return params, nil
}

func newPathParam(opts *paramTagOpts) (*Parameter, error) {
	if opts.style == "" {
		opts.style = StyleSimple
	}
	validStyles := []Style{StyleSimple, StyleLabel, StyleMatrix}
	if !slices.Contains(validStyles, opts.style) {
		return nil, fmt.Errorf("invalid style: %s. Valid styles are: %v", opts.style, validStyles)
	}
	return &Parameter{
		ParamIn:    ParamInPath,
		Name:       opts.name,
		Deprecated: opts.deprecated,
		Style:      opts.style,
		Explode:    opts.explode,
		// Path Parameters are always required.
		Required: true,
		Example:  opts.example,
	}, nil
}

func newHeaderParam(opts *paramTagOpts) (*Parameter, error) {
	// Header key must be canonical
	canonicalName := http.CanonicalHeaderKey(opts.name)
	if canonicalName != opts.name {
		return nil, fmt.Errorf("header name is not canonical: %s. Change it to: %s", opts.name, canonicalName)
	}
	return &Parameter{
		ParamIn:    ParamInHeader,
		Name:       canonicalName,
		Deprecated: opts.deprecated,
		// Headers are always style simple
		Style:    StyleSimple,
		Required: opts.required,
		Example:  opts.example,
	}, nil
}

func newQueryParam(opts *paramTagOpts) (*Parameter, error) {
	if opts.style == "" {
		opts.style = StyleForm
		opts.explode = true
	}
	validStyles := []Style{StyleForm, StyleSpaceDelim, StylePipeDelim, StyleDeepObject}
	if !slices.Contains(validStyles, opts.style) {
		return nil, fmt.Errorf("invalid style: %s. Valid styles are: %v", opts.style, validStyles)
	}
	if opts.style == StyleDeepObject {
		opts.explode = true
	}
	return &Parameter{
		ParamIn:    ParamInQuery,
		Name:       opts.name,
		Deprecated: opts.deprecated,
		Style:      opts.style,
		Required:   opts.required,
		Example:    opts.example,
		Explode:    opts.explode,
	}, nil
}

func newCookieParam(opts *paramTagOpts) (*Parameter, error) {
	return &Parameter{
		ParamIn:    ParamInCookie,
		Name:       opts.name,
		Deprecated: opts.deprecated,
		Style:      StyleForm,
		Required:   opts.required,
		Example:    opts.example,
	}, nil
}
