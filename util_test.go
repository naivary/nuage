package nuage

import "testing"

func TestIsPatternMatchingPathVars(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		params  []*Parameter
		want    bool
	}{
		{
			name:    "matching",
			want:    true,
			pattern: "/path/to/{p1}/endpoint/{p2}",
			params: []*Parameter{
				{Name: "p1", ParamIn: ParamInPath},
				{Name: "p2", ParamIn: ParamInPath},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isPatternMatchingDefinedParams(tc.pattern, tc.params)
			if got != tc.want {
				t.Errorf("Want: %t. Got: %t", tc.want, got)
			}
		})
	}
}
