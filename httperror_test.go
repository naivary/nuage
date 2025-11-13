package nuage

import (
	"testing"
)

func TestHTTPError_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		err     *HTTPError[any, any]
		isValid bool
	}{
		{
			name: "extra members",
			err: &HTTPError[any, any]{
				Type: "type",
				AdditionalMembers: struct {
					Name string `json:"name"`
				}{
					Name: "test_name",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.err.MarshalJSON()
			if err != nil {
				t.Fatalf("marshal JSON: %v", err)
			}
			t.Log(string(data))
		})
	}
}
