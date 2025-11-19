package nuage

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func assign(lhs reflect.Value, rhs ...string) error {
	if len(rhs) == 0 {
		return errors.New("assign: rhs is empty")
	}
	if canCallIsNil(lhs.Type().Kind()) && lhs.IsNil() {
		lhs.Set(reflect.New(lhs.Type().Elem()))
	}
	if isPointer(lhs.Type()) {
		lhs = lhs.Elem()
	}
	if !lhs.IsValid() {
		return fmt.Errorf("assign: lhs value is invalid %v", lhs)
	}

	switch deref(lhs.Type()).Kind() {
	case reflect.String:
		lhs.SetString(rhs[0])
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		integer, err := strconv.ParseInt(rhs[0], 10, 64)
		if err != nil {
			return err
		}
		if lhs.OverflowInt(integer) {
			return fmt.Errorf("overflow: %d to %s", integer, lhs.Kind())
		}
		lhs.SetInt(integer)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uinteger, err := strconv.ParseUint(rhs[0], 10, 64)
		if err != nil {
			return err
		}
		if lhs.OverflowUint(uinteger) {
			return fmt.Errorf("overflow: %d to %s", uinteger, lhs.Kind())
		}
		lhs.SetUint(uinteger)
	case reflect.Float32, reflect.Float64:
		float, err := strconv.ParseFloat(rhs[0], 64)
		if err != nil {
			return err
		}
		if lhs.OverflowFloat(float) {
			return fmt.Errorf("overflow: %f to %s", float, lhs.Kind())
		}
	case reflect.Complex64, reflect.Complex128:
		c, err := strconv.ParseComplex(rhs[0], 128)
		if err != nil {
			return err
		}
		if lhs.OverflowComplex(c) {
			return fmt.Errorf("overflow: %f to %s", c, lhs.Kind())
		}
		lhs.SetComplex(c)
	case reflect.Bool:
		boolean, err := strconv.ParseBool(rhs[0])
		if err != nil {
			return err
		}
		lhs.SetBool(boolean)
	case reflect.Slice:
		elems := make([]reflect.Value, 0, len(rhs))
		for _, rh := range rhs {
			elem, isPtr := newVar(lhs.Type().Elem())
			err := assign(elem, rh)
			if err != nil {
				return err
			}
			if !isPtr {
				elems = append(elems, elem.Elem())
				continue
			}
			elems = append(elems, elem)
		}
		s := reflect.Append(lhs, elems...)
		lhs.Set(s)
	case reflect.Map:
		if len(rhs)%2 != 0 {
			return fmt.Errorf("assign: map expects an even number of rhs values. Got: %d", len(rhs))
		}
		key, isKeyPtr := newVar(lhs.Type().Key())
		value, isValuePtr := newVar(lhs.Type().Elem())
		for i := 0; i < len(rhs); i += 2 {
			err := assign(key, rhs[i])
			if err != nil {
				return err
			}
			err = assign(value, rhs[i+1])
			if err != nil {
				return err
			}
			var (
				k = key
				v = value
			)
			if !isKeyPtr {
				k = key.Elem()
			}
			if !isValuePtr {
				v = value.Elem()
			}
			lhs.SetMapIndex(k, v)
		}
	default:
		return fmt.Errorf("cannot assign: %s to %s", rhs, lhs)
	}
	return nil
}

// newVar returns a new reflect.Value based on the given reflect.Type and
// assures that it is not a double pointer. If the original type was a pointer
// it will be indicated by the second return value.
func newVar(typ reflect.Type) (reflect.Value, bool) {
	isPtr := isPointer(typ)
	if isPtr {
		typ = typ.Elem()
	}
	return reflect.New(typ), isPtr
}

func deref(rtype reflect.Type) reflect.Type {
	if rtype.Kind() == reflect.Pointer {
		return rtype.Elem()
	}
	return rtype
}

func derefValue(rvalue reflect.Value) reflect.Value {
	if isPointer(rvalue.Type()) {
		rvalue = rvalue.Elem()
	}
	return rvalue
}

func isPointer(rtype reflect.Type) bool {
	return rtype.Kind() == reflect.Pointer
}

func isStruct[T any]() bool {
	rtype := reflect.TypeFor[T]()
	return deref(rtype).Kind() == reflect.Struct
}

// canCallIsNil is reporting if it is save to call the IsNil method of
// reflect.Value without it panicing.
func canCallIsNil(kind reflect.Kind) bool {
	switch kind {
	case reflect.Pointer, reflect.Interface, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice:
		return true
	default:
		return false
	}
}

func fieldsOf[S any]() ([]reflect.StructField, error) {
	if !isStruct[S]() {
		return nil, errors.New("fields of: is not struct")
	}
	rtype := reflect.TypeFor[S]()
	fields := make([]reflect.StructField, 0, rtype.NumField())
	for i := range rtype.NumField() {
		fields = append(fields, rtype.Field(i))
	}
	return fields, nil
}
