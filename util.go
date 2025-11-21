package nuage

import (
	"encoding/json"
	"fmt"
	"strings"
)

func structToMap[S any](v *S) (map[string]any, error) {
	if !isStruct[S]() {
		return nil, fmt.Errorf("struct to map: type is not struct")
	}
	m := make(map[string]any)
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return m, json.Unmarshal(data, &m)
}

func isContentTypeJSON(contentType string) bool {
	if contentType == ContentTypeJSON {
		return true
	}
	_, after, found := strings.Cut(contentType, "+")
	if !found {
		return false
	}
	return after == "json"
}
