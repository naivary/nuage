package nuage

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func assign(lhs reflect.Value, rhs ...string) error {
	if len(rhs) == 0 {
		return errors.New("assign: rhs is empty")
	}
	if canCallIsNil(lhs.Type().Kind()) && lhs.IsNil() {
		switch lhs.Kind() {
		case reflect.Slice:
			lhs.Set(reflect.MakeSlice(lhs.Type(), 0, 1))
		default:
			lhs.Set(reflect.New(lhs.Type().Elem()))
		}
	}
	if isPointer(lhs.Type()) {
		lhs = lhs.Elem()
	}
	if !lhs.IsValid() {
		return fmt.Errorf("assign: lhs value is invalid %v", lhs)
	}

	switch lhs.Interface().(type) {
	case string:
		lhs.SetString(rhs[0])
		return nil
	case int, int8, int16, int32, int64:
		integer, err := strconv.ParseInt(rhs[0], 10, 64)
		if err != nil {
			return err
		}
		if lhs.OverflowInt(integer) {
			return fmt.Errorf("overflow: %d to %s", integer, lhs.Kind())
		}
		lhs.SetInt(integer)
		return nil
	case uint, uint8, uint16, uint32, uint64:
		uinteger, err := strconv.ParseUint(rhs[0], 10, 64)
		if err != nil {
			return err
		}
		if lhs.OverflowUint(uinteger) {
			return fmt.Errorf("overflow: %d to %s", uinteger, lhs.Kind())
		}
		lhs.SetUint(uinteger)
		return nil
	case float32, float64:
		float, err := strconv.ParseFloat(rhs[0], 64)
		if err != nil {
			return err
		}
		if lhs.OverflowFloat(float) {
			return fmt.Errorf("overflow: %f to %s", float, lhs.Kind())
		}
		return nil
	case complex64, complex128:
		c, err := strconv.ParseComplex(rhs[0], 128)
		if err != nil {
			return err
		}
		if lhs.OverflowComplex(c) {
			return fmt.Errorf("overflow: %f to %s", c, lhs.Kind())
		}
		lhs.SetComplex(c)
		return nil
	case bool:
		boolean, err := strconv.ParseBool(rhs[0])
		if err != nil {
			return err
		}
		lhs.SetBool(boolean)
		return nil
	case time.Time:
		t, err := time.Parse(time.RFC3339, rhs[0])
		if err != nil {
			return err
		}
		rvalue := reflect.ValueOf(t)
		lhs.Set(rvalue)
		return nil
	// nil is needed because the type of any is nil
	case any, nil:
		if lhs.Kind() == reflect.Slice || lhs.Kind() == reflect.Map {
			break
		}
		lhs.Set(reflect.ValueOf(rhs[0]))
		return nil
	}

	switch lhs.Kind() {
	case reflect.Slice, reflect.Array:
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
		return nil
	case reflect.Map:
		if len(rhs)%2 != 0 {
			return fmt.Errorf("assign: lhs kind of map expects an even number of rhs values. Got: %d", len(rhs))
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
	}
	return fmt.Errorf("cannot assign: %s to %v", rhs, lhs)
}

// newVar allocates a new value of the given reflect.Type and returns it as a
// reflect.Value. If the provided type is a pointer type, newVar returns a
// pointer to a newly allocated zero value of the element type and the second
// return value is true. If the provided type is not a pointer, the returned
// reflect.Value still contains a pointer (as reflect.New always returns a
// pointer to the value), but the second return value is false.
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

// fieldsOf returns all visible fields of S if S is a struct.
func fieldsOf[S any]() ([]reflect.StructField, error) {
	if !isStruct[S]() {
		return nil, errors.New("fields of: is not struct")
	}
	rtype := reflect.TypeFor[S]()
	return reflect.VisibleFields(rtype), nil
}

func isIgnoredFromJSONMarshal(field reflect.StructField) bool {
	jsonTagValue := field.Tag.Get("json")
	return jsonTagValue == "-"
}
