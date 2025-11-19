package nuage

import (
	"reflect"
	"testing"
)

func varFor[T any]() reflect.Value {
	return reflect.New(reflect.TypeFor[T]()).Elem()
}

func isEqualString(lhs reflect.Value, rhs []string) bool {
	if isPointer(lhs.Type()) {
		lhs = lhs.Elem()
	}
	return lhs.Interface().(string) == rhs[0]
}

func TestAssign(t *testing.T) {
	tests := []struct {
		name      string
		lhs       reflect.Value
		rhs       []string
		isEqual   func(lhs reflect.Value, rhs []string) bool
		isInvalid bool
	}{
		{
			name:    "string to string",
			lhs:     varFor[string](),
			rhs:     []string{"t1"},
			isEqual: isEqualString,
		},
		{
			name:    "*string to string",
			lhs:     varFor[*string](),
			rhs:     []string{"t1"},
			isEqual: isEqualString,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := assign(tc.lhs, tc.rhs...)
			if err == nil && tc.isInvalid {
				t.Errorf("expected an error: %v", err)
			}
			if !tc.isEqual(tc.lhs, tc.rhs) {
				t.Errorf("values not equal: lhs: %s; rhs: %v", tc.lhs, tc.rhs)
			}
		})
	}
}
