package nuage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

func decode[T any](r *http.Request, v *T) error {
	rvalue := reflect.ValueOf(v).Elem()
	rtype := rvalue.Type()
	if !isStruct[T]() {
		return fmt.Errorf("type must be a struct: %v", rtype)
	}
	// decode parameters
	for i := range rtype.NumField() {
		field := rtype.Field(i)
		if !field.IsExported() {
			// unexported fields will be ignored for decoding
			continue
		}
		opts, err := parseParamTagOpts(field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return err
		}

		fieldValue := rvalue.Field(i)
		var rhs []string
		switch opts.tagKey {
		case _tagKeyPath:
			slug := r.PathValue(opts.name)
			if slug == "" && opts.required {
				return fmt.Errorf("decode: missing required path param %v", opts.name)
			}
			rhs, err = serializePathParam(slug, field.Type, opts.style, opts.explode)
		case _tagKeyQuery:
			rhs, err = serializeQueryParam(r.URL.Query(), opts.name, opts.queryKeys, field.Type, opts.style, opts.explode)
		case _tagKeyHeader:
			rhs, err = serializeHeaderParam(r.Header, opts.name, field.Type, opts.style, opts.explode)
		case _tagKeyCookie:
			rhs, err = serializeHeaderParam(r.Header, opts.name, field.Type, opts.style, opts.explode)
		}
		if err != nil {
			return err
		}
		err = assign(fieldValue, rhs...)
		if err != nil {
			return err
		}
	}
	// decode payload
	if r.Body == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(v)
}
