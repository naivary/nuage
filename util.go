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
