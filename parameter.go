package nuage

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrParamStyleNotSupported = errors.New(`
This parameter style is valid in OpenAPI but intentionally unsupported in this framework version.

nuage only supports a restricted set of parameter styles to guarantee clarity, interoperability, and predictable request parsing.
Currently the following parameter styles are supported:
	1. Path: simple
	2. Query: form, deepObject
	3. Header: simple
	4. Cookie: form

If you have a concrete use case that cannot be expressed with the supported styles, please open a GitHub issue and describe the problem you are trying to solve`)

type paramTagOpts struct {
	in           ParamIn
	name         string
	style        ParamStyle
	required     bool
	explode      bool
	isDeprecated bool
}

func newPathParam(opts *paramTagOpts) (*Parameter, error) {
	switch opts.style {
	case ParamStyleSimple:
	default:
		return nil, ErrParamStyleNotSupported
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
	switch opts.style {
	case ParamStyleSimple:
	default:
		return nil, ErrParamStyleNotSupported
	}

	canonicalName := http.CanonicalHeaderKey(opts.name)
	if canonicalName != opts.name {
		return nil, fmt.Errorf(
			"header parameter: name is not canonical. Change it to: %s",
			canonicalName,
		)
	}
	return &Parameter{
		ParamIn:    ParamInHeader,
		Name:       canonicalName,
		Deprecated: opts.isDeprecated,
		// Headers are always style simple
		Style:    opts.style,
		Required: opts.required,
	}, nil
}

func newCookieParam(opts *paramTagOpts) (*Parameter, error) {
	switch opts.style {
	case ParamStyleForm:
	default:
		return nil, ErrParamStyleNotSupported
	}

	return &Parameter{
		ParamIn:    ParamInCookie,
		Name:       opts.name,
		Deprecated: opts.isDeprecated,
		Style:      opts.style,
		Required:   opts.required,
	}, nil
}

func newQueryParam(opts *paramTagOpts) (*Parameter, error) {
	switch opts.style {
	case ParamStyleForm, ParamStyleDeepObject:
	default:
		return nil, ErrParamStyleNotSupported
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
