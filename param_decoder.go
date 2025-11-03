package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

var errInvalidKind = errors.New("invalid kind")

func DecodePath[T any](r *http.Request, v *T) error {
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
		opts, err := parseTagOpts(_tagKeyPath, field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		pathValue := r.PathValue(opts.Name)
		if pathValue == "" && !opts.Optional {
			return fmt.Errorf("empty path value: %s", opts.Name)
		}

		// decode it into the value provided
		switch field.Type.Kind() {
		case reflect.String:
			rvalue.Field(i).SetString(pathValue)
		case reflect.Int:
			integer, err := strconv.ParseInt(pathValue, 10, 64)
			if err != nil {
				return err
			}
			rvalue.Field(i).SetInt(integer)
		}
	}
	return nil
}

func DecodeHeader[T any](r *http.Request, v *T) error {
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
		opts, err := parseTagOpts(_tagKeyHeader, field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		headerValue := r.Header.Get(opts.Name)
		if headerValue == "" && opts.Required {
			return fmt.Errorf("required header value not set: %s", opts.Name)
		}
		v, err := parseValue(field.Type.Kind(), headerValue)
		if err != nil {
			return err
		}
		if isAssignable(rvalue.Field(i), v) {
			err = assign(rvalue.Field(i), v)
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func DecodeQuery[T any](r *http.Request, v *T) error {
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
		opts, err := parseTagOpts(_tagKeyQuery, field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		queryValue := r.URL.Query().Get(opts.Name)
		if queryValue == "" && opts.Required {
			return fmt.Errorf("required query value not set: %s", opts.Name)
		}
		if queryValue == "" {
			continue
		}
	}
	return nil
}

func DecodeCookie[T any](r *http.Request, v *T) error {
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
		opts, err := parseTagOpts(_tagKeyCookie, field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		cookie, err := r.Cookie(opts.Name)
		if errors.Is(err, http.ErrNoCookie) && opts.Required {
			return fmt.Errorf("required cookie not set: %s", opts.Name)
		}
		if errors.Is(err, http.ErrNoCookie) && !opts.Required {
			continue
		}
		if err != nil {
			return err
		}
		_ = cookie
	}
	return nil
}

func parseValue(kind reflect.Kind, s string) (reflect.Value, error) {
	rvzero := reflect.Value{}
	switch kind {
	case reflect.String:
		return reflect.ValueOf(s), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		integer, err := strconv.ParseInt(s, 10, 64)
		return reflect.ValueOf(integer), err
	case reflect.Float32, reflect.Float64:
		float, err := strconv.ParseFloat(s, 64)
		return reflect.ValueOf(float), err
	case reflect.Bool:
		boolean, err := strconv.ParseBool(s)
		return reflect.ValueOf(boolean), err
	default:
		return rvzero, errInvalidKind
	}
}
