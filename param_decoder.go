package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

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

		// decode it into the value provided
		switch field.Type.Kind() {
		case reflect.String:
			rvalue.Field(i).SetString(headerValue)
		case reflect.Int:
			integer, err := strconv.ParseInt(headerValue, 10, 64)
			if err != nil {
				return err
			}
			rvalue.Field(i).SetInt(integer)
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
		// decode it into the value provided
		switch field.Type.Kind() {
		case reflect.String:
			rvalue.Field(i).SetString(queryValue)
		case reflect.Int:
			integer, err := strconv.ParseInt(queryValue, 10, 64)
			if err != nil {
				return err
			}
			rvalue.Field(i).SetInt(integer)
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
		fmt.Println(err)
		if errors.Is(err, http.ErrNoCookie) && opts.Required {
			return fmt.Errorf("required cookie not set: %s", opts.Name)
		}
		if err != nil {
			return err
		}
		cookieValue := cookie.Value
		// decode it into the value provided
		switch field.Type.Kind() {
		case reflect.String:
			rvalue.Field(i).SetString(cookieValue)
		case reflect.Int:
			integer, err := strconv.ParseInt(cookieValue, 10, 64)
			if err != nil {
				return err
			}
			rvalue.Field(i).SetInt(integer)
		}
	}
	return nil
}
