package nuage

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

var errTagNotFound = errors.New("tag not found")

var _tagKeys = []string{_tagKeyPath, _tagKeyHeader, _tagKeyQuery, _tagKeyCookie}

const (
	// query params
	_tagKeyQuery     = "query"
	_tagKeyQueryKeys = "querykeys"

	// header
	_tagKeyHeader = "header"

	// cookie
	_tagKeyCookie = "cookie"

	// path
	_tagKeyPath = "path"

	// agnostic
	_tagKeyParamStyle = "paramstyle"
)

type paramTagOpts struct {
	Required   bool
	Explode    bool
	Name       string
	Deprecated bool
	Style      Style
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
	if slices.Contains(values, "required") {
		opts.Required = true
	}
	if slices.Contains(values, "explode") {
		opts.Explode = true
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
	// example option
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
