package nuage

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func serializePathParam(v string, typ reflect.Type, style Style, explode bool) ([]string, error) {
	if v == "" {
		return []string{}, nil
	}
	if style == "" {
		style = _defaultPathParamStyle
	}
	kind := deref(typ).Kind()
	switch style {
	case StyleSimple:
		return serializePathParamStyleSimple(v, kind, explode)
	case StyleLabel:
		return serializePathParamStyleLabel(v, kind, explode)
	case StyleMatrix:
		return serializePathParamStyleMatrix(v, kind, explode)
	}
	return nil, fmt.Errorf("serialized path param: invalid style %s", style)
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

func serializeQueryParam(q url.Values, name []string, typ reflect.Type, style Style, explode bool) ([]string, error) {
	typ = deref(typ)
	switch style {
	case StyleForm:
		switch typ.Kind() {
		case reflect.Map:
			if explode {
				values := make([]string, 0, len(name))
				for _, n := range name {
					values = append(values, n, q.Get(n))
				}
				return values, nil
			}
		default:
			return q[name[0]], nil
		}
	case StyleSpaceDelim:
		if explode {
			return q[name[0]], nil
		}
		value := q.Get(name[0])
		return strings.Split(value, " "), nil
	case StylePipeDelim:
		if explode {
			return q[name[0]], nil
		}
		value := q.Get(name[0])
		return strings.Split(value, "|"), nil
	case StyleDeepObject:
		if !explode {
			return nil, fmt.Errorf("serialize query param: deep object can only be used with explode=true")
		}
		values := make([]string, 0, len(q))
		prefix := name[0]
		for key := range q {
			if !strings.HasPrefix(key, prefix) {
				continue
			}
			keyDeepObj := strings.TrimPrefix(key, prefix)
			keyDeepObj = keyDeepObj[1 : len(keyDeepObj)-1]
			value := q.Get(key)
			values = append(values, keyDeepObj, value)
		}
	}
	return nil, nil
}
