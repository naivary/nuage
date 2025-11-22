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
	tagKey string
	// agnostic options
	name       string
	required   bool
	explode    bool
	deprecated bool
	style      Style
	example    string

	// query param
	queryKeys []string
}

func parseParamTagOpts(field reflect.StructField) (*paramTagOpts, error) {
	opts := paramTagOpts{}
	paramTagKey, paramTagValue := paramTagKeyOf(field)
	if paramTagKey == "" {
		return nil, errTagNotFound
	}
	opts.tagKey = paramTagKey
	if !isIgnoredFromJSONMarshal(field) {
		return nil, fmt.Errorf(
			`parse paramn tag opts: make sure that if you are tagging a struct field to be decoded from parameters it is ignored from json payloads e.g. json:"-"`,
		)
	}

	values := strings.Split(paramTagValue, ",")
	if len(values) == 0 {
		return nil, fmt.Errorf("empty tag (%s in %v): need at least one name value", paramTagKey, field)
	}
	// first element of the tag is always the name
	opts.name = values[0]
	values = values[1:]
	if slices.Contains(values, "deprecated") {
		opts.deprecated = true
	}
	if slices.Contains(values, "required") {
		opts.required = true
	}
	if slices.Contains(values, "explode") {
		opts.explode = true
	}
	exampleTagValue, ok := field.Tag.Lookup(_tagKeyParamExample)
	if ok && exampleTagValue == "" {
		return nil, fmt.Errorf("paramexample cannot be empty: %v", field)
	}
	opts.example = exampleTagValue

	styleTagValue, ok := field.Tag.Lookup(_tagKeyParamStyle)
	opts.style = Style(styleTagValue)
	if ok && !opts.style.IsValid() {
		return nil, fmt.Errorf("invalid param style: %s in %v", paramTagKey, field)
	}

	queryKeys, ok := field.Tag.Lookup(_tagKeyQueryKeys)
	if ok && queryKeys == "" {
		return nil, fmt.Errorf("querykeys tag cannot be empty: %v", field)
	}
	opts.queryKeys = strings.Split(queryKeys, ",")
	return &opts, nil
}

func paramTagKeyOf(field reflect.StructField) (string, string) {
	// find the tag key of the parameter e.g. header, query etc.
	for _, tagKey := range _tagKeys {
		paramTagValue, ok := field.Tag.Lookup(tagKey)
		if !ok {
			continue
		}
		return tagKey, paramTagValue
	}
	return "", ""
}
