package openapiutil

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/naivary/nuage/openapi"
)

var ErrParamStyleNotSupported = errors.New(`
This parameter style is valid in OpenAPI but intentionally unsupported in this framework version.

nuage only supports a restricted set of parameter styles to guarantee clarity, interoperability, and predictable request parsing.
Currently the following parameter styles are supported:
	1. Path: simple
	2. Query: form
	3. Header: simple
	4. Cookie: form

If you have a concrete use case that cannot be expressed with the supported styles, please open a GitHub issue and describe the problem you are trying to solve`)

func ParamLocation(tag reflect.StructTag) openapi.ParamIn {
	if _, ok := tag.Lookup(openapi.ParamInPath.String()); ok {
		return openapi.ParamInPath
	}
	if _, ok := tag.Lookup(openapi.ParamInQuery.String()); ok {
		return openapi.ParamInQuery
	}
	if _, ok := tag.Lookup(openapi.ParamInHeader.String()); ok {
		return openapi.ParamInHeader
	}
	if _, ok := tag.Lookup(openapi.ParamInCookie.String()); ok {
		return openapi.ParamInCookie
	}
	return ""
}

func defaultParamOpts(name string, in openapi.ParamIn) *ParamOpts {
	switch in {
	case openapi.ParamInPath:
		return &ParamOpts{
			In:       in,
			Name:     name,
			Style:    openapi.ParamStyleSimple,
			Explode:  false,
			Required: true,
		}
	case openapi.ParamInQuery:
		return &ParamOpts{
			In:      in,
			Name:    name,
			Style:   openapi.ParamStyleForm,
			Explode: true,
		}
	case openapi.ParamInHeader:
		return &ParamOpts{
			In:    in,
			Name:  name,
			Style: openapi.ParamStyleSimple,
		}
	case openapi.ParamInCookie:
		return &ParamOpts{
			In:      in,
			Name:    name,
			Style:   openapi.ParamStyleForm,
			Explode: true,
		}
	default:
		return nil
	}
}

type ParamOpts struct {
	In           openapi.ParamIn
	Name         string
	Style        openapi.ParamStyle
	Required     bool
	Explode      bool
	IsDeprecated bool
	Default      any
}

func ParseParamOpts(tag reflect.StructTag) (*ParamOpts, error) {
	in := ParamLocation(tag)
	if in == "" {
		// tag is not a parameter or at an invalid location
		return nil, nil
	}
	tagValue := tag.Get(in.String())
	if len(tagValue) == 0 {
		return nil, errors.New("parameter tag value is empty")
	}
	definedOpts := strings.Split(tagValue, ",")
	opts := defaultParamOpts(definedOpts[0], in)
	// per default all parameters are required and become
	// optional when a default is set.
	for _, opt := range definedOpts[1:] {
		if opt == "deprecated" {
			opts.IsDeprecated = true
		}
		if opt == "required" {
			opts.Required = true
		}
		if strings.HasPrefix(opt, "explode") {
			_, value, _ := strings.Cut(opt, "=")
			if value == "" {
				return nil, fmt.Errorf("rhs of `%s` is empty", opt)
			}
			e, err := strconv.ParseBool(value)
			if err != nil {
				return nil, err
			}
			opts.Explode = e
		}
		if strings.HasPrefix(opt, "style") {
			_, value, _ := strings.Cut(opt, "=")
			opts.Style = openapi.ParamStyle(value)
		}
		if strings.HasPrefix(opt, "default") && in != openapi.ParamInPath {
			_, value, _ := strings.Cut(opt, "=")
			opts.Default = any(value)
			opts.Required = false
		}
	}
	return opts, nil
}

func NewPathParam(opts *ParamOpts) (*openapi.Parameter, error) {
	switch opts.Style {
	case openapi.ParamStyleSimple:
	default:
		return nil, ErrParamStyleNotSupported
	}
	return &openapi.Parameter{
		ParamIn:    openapi.ParamInPath,
		Name:       opts.Name,
		Deprecated: opts.IsDeprecated,
		Style:      opts.Style,
		Explode:    opts.Explode,
		// Path Parameters are always required
		Required: true,
	}, nil
}

func NewHeaderParam(opts *ParamOpts) (*openapi.Parameter, error) {
	switch opts.Style {
	case openapi.ParamStyleSimple:
	default:
		return nil, ErrParamStyleNotSupported
	}
	canonicalName := http.CanonicalHeaderKey(opts.Name)
	if canonicalName != opts.Name {
		return nil, fmt.Errorf(
			"header parameter: name is not canonical. Change it to: %s",
			canonicalName,
		)
	}
	return &openapi.Parameter{
		ParamIn:    openapi.ParamInHeader,
		Name:       canonicalName,
		Deprecated: opts.IsDeprecated,
		Style:      opts.Style,
		Required:   opts.Required,
	}, nil
}

func NewCookieParam(opts *ParamOpts) (*openapi.Parameter, error) {
	switch opts.Style {
	case openapi.ParamStyleForm:
	default:
		return nil, ErrParamStyleNotSupported
	}
	return &openapi.Parameter{
		ParamIn:    openapi.ParamInCookie,
		Name:       opts.Name,
		Deprecated: opts.IsDeprecated,
		Style:      opts.Style,
		Required:   opts.Required,
	}, nil
}

func NewQueryParam(opts *ParamOpts) (*openapi.Parameter, error) {
	switch opts.Style {
	case openapi.ParamStyleForm:
	default:
		return nil, ErrParamStyleNotSupported
	}
	return &openapi.Parameter{
		ParamIn:    openapi.ParamInQuery,
		Name:       opts.Name,
		Deprecated: opts.IsDeprecated,
		Style:      opts.Style,
		Required:   opts.Required,
		Explode:    opts.Explode,
	}, nil
}
