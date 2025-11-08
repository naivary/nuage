package nuage

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func assign(lhs reflect.Value, rhs ...string) error {
	if len(rhs) == 0 {
		return errors.New("rhs is empty")
	}
	if isPointer(lhs.Type()) && !lhs.IsNil() {
		lhs = lhs.Elem()
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
		elemType := lhs.Type().Elem()
		isPointer := isPointer(elemType)
		if isPointer {
			elemType = elemType.Elem()
		}
		elems := make([]reflect.Value, 0, len(rhs))
		for _, rh := range rhs {
			elem := reflect.New(elemType)
			err := assign(elem, rh)
			if err != nil {
				return err
			}
			if !isPointer {
				elems = append(elems, elem.Elem())
				continue
			}
			elems = append(elems, elem)
		}
		reflect.Append(lhs, elems...)
	case reflect.Map:
		// TODO(naivary): finding and setting the value is pretty repetitive
		// with the slice type. should find a better solution or outsource to
		// another function.
		if len(rhs) < 2 {
			return fmt.Errorf("invalid rhs: map expects at least two rhs values. Got: %d", len(rhs))
		}
		keyType := lhs.Type().Key()
		isKeyTypePtr := isPointer(keyType)
		if isKeyTypePtr {
			keyType = keyType.Elem()
		}
		valueType := lhs.Type().Elem()
		isValueTypePtr := isPointer(valueType)
		if isValueTypePtr {
			valueType = valueType.Elem()
		}
		key := reflect.New(keyType)
		value := reflect.New(valueType)
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
			if !isKeyTypePtr {
				k = key.Elem()
			}
			if !isValueTypePtr {
				v = value.Elem()
			}
			lhs.SetMapIndex(k, v)
		}
	default:
		return fmt.Errorf("cannot assign: %s to %s", rhs, lhs)
	}
	return nil
}

func newVar(typ reflect.Type) (reflect.Value, bool) {
	isPtr := isPointer(typ)
	if isPtr {
		typ = typ.Elem()
	}
	return reflect.New(typ), isPtr
}

func ptrTo[T any](v T) *T {
	return &v
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
