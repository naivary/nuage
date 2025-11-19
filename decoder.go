package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

func decodeParams[T any](r *http.Request, v *T) error {
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
				return fmt.Errorf("decode: missing required path param %s", opts.name)
			}
			rhs, err = serializePathParam(slug, field.Type, opts.style, opts.explode)
		case _tagKeyQuery:
			if !r.URL.Query().Has(opts.name) && opts.required {
				return fmt.Errorf("decode: missing required query param %s", opts.name)
			}
			rhs, err = serializeQueryParam(r.URL.Query(), opts.name, opts.queryKeys, field.Type, opts.style, opts.explode)
		case _tagKeyHeader:
			if value := r.Header.Get(opts.name); value == "" && opts.required {
				return fmt.Errorf("decode: missing required header param %s", opts.name)
			}
			rhs, err = serializeHeaderParam(r.Header, opts.name, field.Type, opts.style, opts.explode)
		case _tagKeyCookie:
			cookie, err := r.Cookie(opts.name)
			if errors.Is(err, http.ErrNoCookie) && opts.required {
				return fmt.Errorf("decode: missing required cookie param %s", opts.name)
			}
			rhs, err = serializeCookieParam(cookie, field.Type, opts.style, opts.explode)
		}
		if err != nil {
			return err
		}
		err = assign(fieldValue, rhs...)
		if err != nil {
			return err
		}
	}
	return nil
}
