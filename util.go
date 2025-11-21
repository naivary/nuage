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
