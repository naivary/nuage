package nuage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

func Decode[T any](r *http.Request, v *T) error {
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
		for _, tagKey := range _tagKeys {
			opts, err := parseParamTagOpts(tagKey, field)
			if errors.Is(err, errTagNotFound) {
				continue
			}
			if err != nil {
				return err
			}

			fieldValue := rvalue.Field(i)
			var rhs []string
			switch tagKey {
			case _tagKeyPath:
				slug := r.PathValue(opts.Name)
				if slug == "" && opts.Required {
					return fmt.Errorf("decode: missing required path param %v", opts.Name)
				}
				rhs, err = SerializePathParam(slug, field.Type, opts.Style, opts.Explode)
			case _tagKeyQuery:
				rhs, err = SerializeQueryParam(r.URL.Query(), opts.Name, opts.QueryKeys, field.Type, opts.Style, opts.Explode)
			case _tagKeyHeader:
				rhs, err = SerializeHeaderParam(r.Header, opts.Name, field.Type, opts.Style, opts.Explode)
			case _tagKeyCookie:
				rhs, err = SerializeHeaderParam(r.Header, opts.Name, field.Type, opts.Style, opts.Explode)
			}
			if err != nil {
				return err
			}
			err = assign(fieldValue, rhs...)
			if err != nil {
				return err
			}
			// one field can only be one type of parameter.
			break
		}
	}
	// decode payload
	if r.Body == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(v)
}
