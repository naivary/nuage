package nuage

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"

	"github.com/naivary/nuage/openapi"
)

func serializePathParam(v string, typ reflect.Type, style openapi.Style, explode bool) ([]string, error) {
	if v == "" {
		return []string{}, nil
	}
	if style == "" {
		style = openapi.StyleSimple
	}
	kind := deref(typ).Kind()
	switch style {
	case openapi.StyleSimple:
		return serializePathParamStyleSimple(v, kind, explode)
	case openapi.StyleLabel:
		return serializePathParamStyleLabel(v, kind, explode)
	case openapi.StyleMatrix:
		return serializePathParamStyleMatrix(v, kind, explode)
	}
	return nil, fmt.Errorf("serialize path param: invalid style %s", style)
}

func serializePathParamStyleSimple(v string, kind reflect.Kind, explode bool) ([]string, error) {
	switch kind {
	case reflect.Slice:
		return strings.Split(v, ","), nil
	case reflect.Map:
		if explode {
			return pathParamKeyValuePairs(v, ",")
		}
		return strings.Split(v, ","), nil
	}
	return []string{v}, nil
}

func serializePathParamStyleLabel(v string, kind reflect.Kind, explode bool) ([]string, error) {
	switch kind {
	case reflect.Slice:
		sep := ","
		if explode {
			sep = "."
		}
		return strings.Split(v[1:], sep), nil
	case reflect.Map:
		if explode {
			return pathParamKeyValuePairs(v[1:], ".")
		}
		return strings.Split(v[1:], ","), nil
	}
	return []string{v[1:]}, nil
}

func serializePathParamStyleMatrix(v string, kind reflect.Kind, explode bool) ([]string, error) {
	switch kind {
	case reflect.Slice:
		if explode {
			pairs, err := pathParamKeyValuePairs(v[1:], ";")
			if err != nil {
				return nil, err
			}
			values := make([]string, 0, len(pairs))
			for i := 1; i < len(pairs); i += 2 {
				values = append(values, pairs[i])
			}
			return values, nil
		}
		values := strings.Split(v, "=")
		if len(values) != 2 {
			return nil, fmt.Errorf("serialize path param: invalid syntax for matrix style %s", v)
		}
		return strings.Split(values[1], ","), nil
	case reflect.Map:
		if explode {
			return pathParamKeyValuePairs(v[1:], ";")
		}
		values := strings.Split(v, "=")
		if len(values) != 2 {
			return nil, fmt.Errorf("serialize path param: invalid syntax for matrix style %s", v)
		}
		return strings.Split(values[1], ","), nil
	}
	_, value, found := strings.Cut(v[1:], "=")
	if !found {
		return nil, fmt.Errorf("serialize path param: invalid syntax for primitive %s", v)
	}
	return []string{value}, nil
}

// pathParamKeyValuePairs parses a string containing key-value pairs separated by a given separator.
//
// The input string `v` should contain one or more key-value pairs in the form "key=value",
// separated by the string `sep`. For example, with `v = "id=123;name=alice"` and `sep = ";"`,
// the function returns: []string{"id", "123", "name", "alice"}.
//
// If any pair does not contain an '=', the function returns an error indicating invalid syntax.
func pathParamKeyValuePairs(v, sep string) ([]string, error) {
	keyValuePairs := strings.Split(v, sep)
	values := make([]string, 0, len(keyValuePairs)*2)
	for _, pair := range keyValuePairs {
		key, value, found := strings.Cut(pair, "=")
		if !found {
			return nil, fmt.Errorf("serialize path param: invalid syntax %s", v)
		}
		values = append(values, key, value)
	}
	return values, nil
}

func serializeQueryParam(q url.Values, name string, keys []string, typ reflect.Type, style openapi.Style, explode bool) ([]string, error) {
	if style == "" {
		style = openapi.StyleForm
	}

	typ = deref(typ)
	switch style {
	case openapi.StyleForm:
		return serializeQueryParamStyleForm(q, name, keys, typ.Kind(), explode)
	case openapi.StyleDeepObject:
		return serializeQueryParamStyleDeepObject(q, name, keys, typ.Kind(), explode)
	case openapi.StyleSpaceDelim:
		if explode {
			return q[name], nil
		}
		value := q.Get(name)
		return strings.Split(value, " "), nil
	case openapi.StylePipeDelim:
		if explode {
			return q[name], nil
		}
		value := q.Get(name)
		return strings.Split(value, "|"), nil
	}
	return nil, fmt.Errorf("invalid style: %s", style)
}

func serializeQueryParamStyleForm(q url.Values, name string, keys []string, kind reflect.Kind, explode bool) ([]string, error) {
	switch kind {
	case reflect.Map:
		if !explode {
			break
		}
		values := make([]string, 0, len(name))
		for _, key := range keys {
			value := q.Get(key)
			values = append(values, key, value)
		}
		return values, nil
	}
	return q[name], nil
}

func serializeQueryParamStyleDeepObject(q url.Values, name string, keys []string, kind reflect.Kind, explode bool) ([]string, error) {
	if !explode {
		return nil, fmt.Errorf("serialize query param: deep object can only be used with explode=true")
	}
	if kind != reflect.Map {
		return nil, fmt.Errorf("invalid kind: kind has to be map for a parameter type of DeepObject")
	}
	values := make([]string, 0, len(q))
	for key := range q {
		if !strings.HasPrefix(key, name) {
			continue
		}
		keyDeepObj := strings.TrimPrefix(key, name)
		keyDeepObj = keyDeepObj[1 : len(keyDeepObj)-1]
		if !slices.Contains(keys, keyDeepObj) && len(keys) > 0 {
			// unknown keys will be skipped
			continue
		}
		value := q.Get(key)
		values = append(values, keyDeepObj, value)
	}
	return values, nil
}

func serializeHeaderParam(header http.Header, key string, typ reflect.Type, style openapi.Style, explode bool) ([]string, error) {
	if style == "" {
		style = openapi.StyleSimple
	}
	if style != openapi.StyleSimple {
		return nil, fmt.Errorf("invalid style: %s", style)
	}
	typ = deref(typ)
	switch typ.Kind() {
	case reflect.Map:
		value := header.Get(key)
		if explode {
			return pathParamKeyValuePairs(value, ",")
		}
	default:
		value := header.Get(key)
		return strings.Split(value, ","), nil
	}
	return nil, fmt.Errorf("invalid kind: %v", typ.Kind())
}

func serializeCookieParam(cookie *http.Cookie, typ reflect.Type, style openapi.Style, explode bool) ([]string, error) {
	if style != openapi.StyleForm {
		return nil, fmt.Errorf("invalid style: %s", style)
	}
	typ = deref(typ)
	kind := typ.Kind()
	if (kind == reflect.Slice || kind == reflect.Map) && explode {
		return nil, fmt.Errorf("cannot serialize exploded cookie parameter into slice or map")
	}
	return strings.Split(cookie.Value, ","), nil
}
