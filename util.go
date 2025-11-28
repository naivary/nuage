package nuage

import (
	"regexp"
	"slices"
	"strings"
)

// _regExpPathParam captures all path parameters in the standard syntax definitoin of the http.ServeMux
var _regExpPathParam = regexp.MustCompile(`(?m){([a-z0-9]+)}`)

func methodOf(pattern string) string {
	method, _, _ := strings.Cut(pattern, " ")
	return method
}

func isHTTPStatus(statusCode int) bool {
	return statusCode >= 100 && statusCode <= 599
}

func isJSONish(contentType string) bool {
	if contentType == ContentTypeJSON {
		return true
	}
	_, after, found := strings.Cut(contentType, "+")
	if !found {
		return false
	}
	return after == "json"
}

func isEmptyJSON[T any]() bool {
	fields, err := fieldsOf[T]()
	if err != nil {
		return false
	}
	for _, field := range fields {
		tagValue, ok := field.Tag.Lookup("json")
		if tagValue != "-" && tagValue != "" && ok {
			return false
		}
	}
	return true
}

func isPatternAndPathParamsConsistent(pattern string, params []*Parameter) bool {
	pathParams := filter(params, func(el *Parameter) bool {
		return el.ParamIn == ParamInPath
	})
	pathParamNames := mapper(pathParams, func(param *Parameter) string {
		return param.Name
	})
	matchInfos := _regExpPathParam.FindAllStringSubmatch(pattern, -1)
	if len(matchInfos) != len(pathParamNames) {
		return false
	}
	for _, matchInfo := range matchInfos {
		if len(matchInfo) != 2 {
			return false
		}
		slug := matchInfo[1]
		if !slices.Contains(pathParamNames, slug) {
			return false
		}
	}
	return true
}

func filter[T any](s []T, fn func(el T) bool) []T {
	result := make([]T, 0, len(s))
	for _, el := range s {
		take := fn(el)
		if take {
			result = append(result, el)
		}
	}
	return result
}

func mapper[T, I any](s []T, fn func(el T) I) []I {
	result := make([]I, 0, len(s))
	for _, el := range s {
		result = append(result, fn(el))
	}
	return result
}

// Ptr returns a pointer to the provided value. It will be deprecated in the future
// with the introduction of the new(expression) function accepting expression.
func Ptr[T any](v T) *T {
	return &v
}
