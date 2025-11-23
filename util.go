package nuage

import (
	"strings"
)

func methodOf(pattern string) string {
	method, _, _ := strings.Cut(pattern, " ")
	return method
}

func isStatusCodeInRange(statusCode int) bool {
	return statusCode >= 100 && statusCode <= 599
}

func isJSONishContentType(contentType string) bool {
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
