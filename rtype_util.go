package nuage

import (
	"fmt"
	"reflect"
	"slices"
)

func isStruct[T any]() bool {
	rtype := reflect.TypeFor[T]()
	return deref(rtype).Kind() == reflect.Struct
}

func deref(rtype reflect.Type) reflect.Type {
	if rtype.Kind() == reflect.Pointer {
		return rtype.Elem()
	}
	return rtype
}

// isAssignable reports whether value a is assignable to value b. It considers
// more options than `CanConvert`. For example CanConvert will report false for
// int and int64 but it might be possible to fit the value of int64 into int.
func isAssignable(a, b reflect.Value) bool {
	// we want to check if a = b is possible.
	if !a.CanSet() {
		// value `a` cannot be changed
		return false
	}
	kind := a.Kind()
	if kind == reflect.Invalid {
		return false
	}
	if isInt(kind) && isInt(b.Kind()) {
		return !a.OverflowInt(b.Int())
	}
	// sane default is false because false positives may result in panics.
	return false
}

func assign(a, b reflect.Value) error {
	if !isAssignable(a, b) {
		return fmt.Errorf("not assignable: %v = %v", a, b)
	}
	kind := a.Kind()
	if isInt(kind) {
		a.SetInt(b.Int())
	}
	return nil
}

func isInt(kind reflect.Kind) bool {
	kinds := []reflect.Kind{
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
	}
	return slices.Contains(kinds, kind)
}
