package nuage

import (
	"fmt"
	"reflect"
	"strconv"
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

func isPointer(rtype reflect.Type) bool {
	return rtype.Kind() == reflect.Pointer
}

func assign(lhs reflect.Value, rhs string) error {
	// if !lhs.CanSet() {
	// 	return fmt.Errorf("lhs is cannot be set: %s", lhs)
	// }
	if isPointer(lhs.Type()) && !lhs.IsNil() {
		lhs = lhs.Elem()
	}

	switch deref(lhs.Type()).Kind() {
	case reflect.String:
		lhs.SetString(rhs)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		integer, err := strconv.ParseInt(rhs, 10, 64)
		if err != nil {
			return err
		}
		if lhs.OverflowInt(integer) {
			return fmt.Errorf("overflow: %d to %s", integer, lhs.Kind())
		}
		lhs.SetInt(integer)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uinteger, err := strconv.ParseUint(rhs, 10, 64)
		if err != nil {
			return err
		}
		if lhs.OverflowUint(uinteger) {
			return fmt.Errorf("overflow: %d to %s", uinteger, lhs.Kind())
		}
		lhs.SetUint(uinteger)
	case reflect.Float32, reflect.Float64:
		float, err := strconv.ParseFloat(rhs, 64)
		if err != nil {
			return err
		}
		if lhs.OverflowFloat(float) {
			return fmt.Errorf("overflow: %f to %s", float, lhs.Kind())
		}
	case reflect.Complex64, reflect.Complex128:
		c, err := strconv.ParseComplex(rhs, 128)
		if err != nil {
			return err
		}
		if lhs.OverflowComplex(c) {
			return fmt.Errorf("overflow: %f to %s", c, lhs.Kind())
		}
		lhs.SetComplex(c)
	case reflect.Bool:
		boolean, err := strconv.ParseBool(rhs)
		if err != nil {
			return err
		}
		lhs.SetBool(boolean)
	case reflect.Slice:
		elemType := lhs.Type().Elem()
		isPointer := isPointer(elemType)
		if isPointer {
			elemType = elemType.Elem()
		}
		elem := reflect.New(elemType)
		err := assign(elem, rhs)
		if err != nil {
			return err
		}
		if !isPointer {
			reflect.Append(lhs, elem.Elem())
			return nil
		}
		reflect.Append(lhs, elem)
	default:
		return fmt.Errorf("cannot assign: %s to %s", rhs, lhs)
	}
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

func isString(kind reflect.Kind) bool {
	return kind == reflect.String
}

func ptrTo[T any](v T) *T {
	return &v
}
