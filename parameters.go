package nuage

import (
	"fmt"
	"net/http"
)

type paramTagOpts struct {
	in           ParamIn
	name         string
	style        ParamStyle
	required     bool
	explode      bool
	isDeprecated bool
	queryKeys    []string
}

func newPathParam(opts *paramTagOpts) (*Parameter, error) {
	switch opts.style {
	case ParamStyleSimple, ParamStyleLabel, ParamStyleMatrix:
	default:
		return nil, fmt.Errorf("path parameter: invalid style `%s`", opts.style)
	}
	return &Parameter{
		ParamIn:    ParamInPath,
		Name:       opts.name,
		Deprecated: opts.isDeprecated,
		Style:      opts.style,
		Explode:    opts.explode,
		// Path Parameters are always required
		Required: true,
	}, nil
}

func newHeaderParam(opts *paramTagOpts) (*Parameter, error) {
	// Header key must be canonical
	canonicalName := http.CanonicalHeaderKey(opts.name)
	if canonicalName != opts.name {
		return nil, fmt.Errorf("header parameter: name is not canonical. Change it to: %s", canonicalName)
	}
	return &Parameter{
		ParamIn:    ParamInHeader,
		Name:       canonicalName,
		Deprecated: opts.isDeprecated,
		// Headers are always style simple
		Style:    ParamStyleSimple,
		Required: opts.required,
	}, nil
}

func newQueryParam(opts *paramTagOpts) (*Parameter, error) {
	switch opts.style {
	case ParamStyleForm, ParamStyleSpaceDelim, ParamStylePipeDelim, ParamStyleDeepObject:
	default:
		return nil, fmt.Errorf("query param: invalid style `%s`", opts.style)
	}
	if opts.style == ParamStyleDeepObject {
		opts.explode = true
	}
	return &Parameter{
		ParamIn:    ParamInQuery,
		Name:       opts.name,
		Deprecated: opts.isDeprecated,
		Style:      opts.style,
		Required:   opts.required,
		Explode:    opts.explode,
	}, nil
}

func newCookieParam(opts *paramTagOpts) (*Parameter, error) {
	return &Parameter{
		ParamIn:    ParamInCookie,
		Name:       opts.name,
		Deprecated: opts.isDeprecated,
		Style:      ParamStyleForm,
		Required:   opts.required,
	}, nil
}
