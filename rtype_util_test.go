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
		rhs   string
	}{
		{
			name:  "ptr string",
			valid: true,
			lhs:   reflect.ValueOf(ptrTo("")),
			rhs:   "test",
		},
		{
			name:  "int",
			valid: true,
			lhs:   reflect.ValueOf(ptrTo(0)),
			rhs:   "12937812",
		},
		{
			name:  "float",
			valid: true,
			lhs:   reflect.ValueOf(ptrTo(0.0)),
			rhs:   "12937812.3123",
		},
		{
			name:  "slice with ptr type",
			valid: true,
			lhs:   reflect.ValueOf([]*string{}),
			rhs:   "test",
		},
		{
			name:  "slice",
			valid: true,
			lhs:   reflect.ValueOf([]string{}),
			rhs:   "test",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := assign(tc.lhs, tc.rhs)
			if err != nil && tc.valid {
				t.Errorf("expected to assign %v = %v", tc.lhs.Kind(), tc.rhs)
			}
		})
	}
}
