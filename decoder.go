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
	for i := range rtype.NumField() {
		field := rtype.Field(i)
		if !field.IsExported() {
			continue
		}
		for _, tagKey := range _tagKeys {
			opts, err := parseTagOpts(tagKey, field)
			if errors.Is(err, errTagNotFound) {
				continue
			}
			if err != nil {
				return err
			}

			var value string
			switch tagKey {
			case _tagKeyPath:
				value = r.PathValue(opts.Name)
				if value == "" && opts.Required {
					return fmt.Errorf("path parameter required: %s", opts.Name)
				}
				seriliazed, err := serializePathParam(value, rvalue.Field(i).Type(), opts.Style, opts.Explode)
				if err != nil {
					return err
				}
				err = assign(rvalue.Field(i), seriliazed...)
				if err != nil {
					return err
				}
			case _tagKeyHeader:
				value = r.Header.Get(opts.Name)
				if value == "" && opts.Required {
					return fmt.Errorf("header required: %s", opts.Name)
				}
			case _tagKeyQuery:
				value = r.URL.Query().Get(opts.Name)
				if value == "" && opts.Required {
					return fmt.Errorf("query required: %s", opts.Name)
				}
			case _tagKeyCookie:
				c, err := r.Cookie(opts.Name)
				if errors.Is(err, http.ErrNoCookie) && opts.Required {
					return fmt.Errorf("cookie required: %s", opts.Name)
				}
				value = c.Value
			}
		}
	}
	if r.Body != nil {
		return json.NewDecoder(r.Body).Decode(v)
	}
	return nil
}
