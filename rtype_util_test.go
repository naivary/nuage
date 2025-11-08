package nuage

import (
	"reflect"
	"testing"
)

func TestIsAssignable(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
		lhs   reflect.Value
		rhs   []string
	}{
		{
			name:  "ptr string",
			valid: true,
			lhs:   reflect.ValueOf(ptrTo("")),
			rhs:   []string{"test"},
		},
		{
			name:  "int",
			valid: true,
			lhs:   reflect.ValueOf(ptrTo(0)),
			rhs:   []string{"12937812"},
		},
		{
			name:  "float",
			valid: true,
			lhs:   reflect.ValueOf(ptrTo(0.0)),
			rhs:   []string{"12937812.3123"},
		},
		{
			name:  "map",
			valid: true,
			lhs:   reflect.ValueOf(map[string]int{}),
			rhs:   []string{"t1", "1", "t2", "2"},
		},
		{
			name:  "map key ptr",
			valid: true,
			lhs:   reflect.ValueOf(map[*string]int{}),
			rhs:   []string{"t1", "1", "t2", "2"},
		},
		{
			name:  "map value ptr",
			valid: true,
			lhs:   reflect.ValueOf(map[string]*int{}),
			rhs:   []string{"t1", "1", "t2", "2"},
		},
		{
			name:  "slice",
			valid: true,
			lhs:   reflect.ValueOf(ptrTo([]string{})).Elem(),
			rhs:   []string{"e1", "e2"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := assign(tc.lhs, tc.rhs...)
			if err != nil && tc.valid {
				t.Errorf("expected to assign %v = %v", tc.lhs.Kind(), tc.rhs)
			}
			t.Logf("lhs: %v", tc.lhs)
		})
	}
}
