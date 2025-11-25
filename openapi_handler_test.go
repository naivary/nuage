package nuage

import (
	"net/http"
	"testing"
)

func TestQueryOpenAPIDoc(t *testing.T) {
	tests := []struct {
		name string
		r    *http.Request
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
		})
	}
}
