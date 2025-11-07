package nuage

import (
	"fmt"
	"reflect"
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
	if deref(b.Type()).ConvertibleTo(deref(a.Type())) {
		return true
	}
	varKind := deref(a.Type()).Kind()
	valueKind := deref(b.Type()).Kind()
	if varKind == reflect.Invalid {
		return false
	}
	if varKind == reflect.String && valueKind == reflect.String {
		return true
	}
	if varKind == reflect.Bool && valueKind == reflect.Bool {
		return true
	}
	if isInt(varKind) && isInt(valueKind) {
		return !a.OverflowInt(b.Int())
	}
	if isFloat(varKind) && isFloat(valueKind) {
		return !a.OverflowFloat(b.Float())
	}
	if isUint(varKind) && isUint(valueKind) {
		return !a.OverflowUint(b.Uint())
	}
	if isComplex(varKind) && isComplex(valueKind) {
		return !a.OverflowComplex(b.Complex())
	}
	// sane default is false because false-positives are more costly resulting
	// in panics in production.
	return false
}

func assign[T any](a reflect.Value, b T) error {
	bvalue := reflect.ValueOf(b)
	if !isAssignable(a, bvalue) {
		return fmt.Errorf("not assignable: %v = %v", a, b)
	}
	a.Set(bvalue)
	return nil
}

func isInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return true
	default:
		return false

	}
}

func isFloat(kind reflect.Kind) bool {
	switch kind {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func isUint(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return true
	default:
		return false
	}
}

func isComplex(kind reflect.Kind) bool {
	switch kind {
	case reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}
