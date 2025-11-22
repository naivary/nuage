package nuage

import (
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
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
		fmt.Println(field)
		tagValue, ok := field.Tag.Lookup("json")
		if tagValue != "-" && tagValue != "" && ok {
			return false
		}
	}
	return true
}

func jsonSchemaFor[T any](opts *jsonschema.ForOptions) (*jsonschema.Schema, error) {
	if opts == nil {
		opts = &jsonschema.ForOptions{}
	}
	return nil, nil
}
