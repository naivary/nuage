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
	_tagKeyParamStyle   = "paramstyle"
	_tagKeyParamExample = "paramexample"
)

type paramTagOpts struct {
	// agnostic options
	Name       string
	Required   bool
	Explode    bool
	Deprecated bool
	Style      Style
	Example    string

	// query param
	QueryKeys []string
}

func parseParamTagOpts(paramTagKey string, field reflect.StructField) (*paramTagOpts, error) {
	opts := paramTagOpts{}
	paramTagValue, ok := field.Tag.Lookup(paramTagKey)
	if !ok {
		return nil, errTagNotFound
	}
	jsonTagValue, ok := field.Tag.Lookup("json")
	if !ok {
		return nil, fmt.Errorf(
			`parse paramn tag opts: make sure that if you are tagging a struct field to be decoded from parameters it is also excluded from json payloads e.g. json:"-"`,
		)
	}
	if jsonTagValue != "-" {
		return nil, fmt.Errorf("parse param tag opts: only tag option for json is '-' for a parameter tagged struct field")
	}

	values := strings.Split(paramTagValue, ",")
	if len(values) == 0 {
		return nil, fmt.Errorf("empty tag (%s in %v): need at least one name value", paramTagKey, field)
	}
	// first element of the tag is always the name
	opts.Name = values[0]
	values = values[1:]
	if slices.Contains(values, "deprecated") {
		opts.Deprecated = true
	}
	if slices.Contains(values, "required") {
		opts.Required = true
	}
	if slices.Contains(values, "explode") {
		opts.Explode = true
	}
	exampleTagValue, ok := field.Tag.Lookup(_tagKeyParamExample)
	if ok && exampleTagValue == "" {
		return nil, fmt.Errorf("paramexample cannot be empty: %v", field)
	}
	opts.Example = exampleTagValue

	styleTagValue, ok := field.Tag.Lookup(_tagKeyParamStyle)
	if ok && !Style(styleTagValue).IsValid() {
		return nil, fmt.Errorf("invalid param style: %s in %v", paramTagKey, field)
	}
	opts.Style = Style(styleTagValue)

	queryKeys, ok := field.Tag.Lookup(_tagKeyQueryKeys)
	if ok && queryKeys == "" {
		return nil, fmt.Errorf("querykeys tag cannot be empty: %v", field)
	}
	opts.QueryKeys = strings.Split(queryKeys, ",")
	return &opts, nil
}
