package codegen

import (
	"errors"
	"fmt"
	"go/types"

	"github.com/naivary/nuage/internal/typesutil"
	"github.com/naivary/nuage/openapi"
)

func isSupportedParamType(in openapi.ParamIn, field *types.Var) error {
	typ := field.Type()
	ptr, isPtr := typ.(*types.Pointer)
	if isPtr {
		typ = ptr.Elem()
	}

	// the following types are not serializable and thus not supported
	// by any of the parameter types.
	switch t := typ.(type) {
	case *types.Signature, *types.Chan:
		return fmt.Errorf("functions or channels are not supported as parameter types because they are not serializable: %s", field.Name())
	case *types.Basic:
		kind := t.Kind()
		if kind == types.Uintptr || kind == types.UnsafePointer {
			return errors.New("uintptr and unsafe pointer are not supported as parameters of any kind")
		}
	case *types.Array:
		return errors.New("arrays are a valid type parameter in general but are hard to support for nuage. For this version it is not supported. Use slices instead")
	}

	switch in {
	case openapi.ParamInPath:
		return isSupportedPathParamType(field, field.Type(), false)
	case openapi.ParamInHeader:
		return isSupportedHeaderParamType(field, field.Type(), false)
	case openapi.ParamInQuery:
		return isSupportedQueryType(field, field.Type(), false)
	case openapi.ParamInCookie:
		return isSupportedCookieType(field, field.Type())
	}
	return nil
}

func isSupportedPathParamType(field *types.Var, typ types.Type, isSlice bool) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedPathParamType(field, t.Underlying(), false)
	case *types.Pointer:
		return isSupportedPathParamType(field, t.Elem(), false)
	case *types.Basic:
		kind := t.Kind()
		if kind == types.Bool {
			return fmt.Errorf("path parameters cannot be of type boolean: %s", field.Name())
		}
		if typesutil.IsComplex(kind) {
			return fmt.Errorf("path parameters cannot be of type complex: %s", field.Name())
		}
		if typesutil.IsFloat(kind) {
			return fmt.Errorf("path parameters cannot be of type float: %s", field.Name())
		}
	case *types.Named:
		return isSupportedPathParamType(field, t.Underlying(), false)
	case *types.Slice:
		if isSlice {
			return fmt.Errorf("path parameters cannot be nested slices: %s", field.Name())
		}
		return isSupportedPathParamType(field, t.Elem(), true)
	default:
		return fmt.Errorf("type `%s` is not supported for path parameters", typ.String())
	}
	return nil
}

func isSupportedHeaderParamType(field *types.Var, typ types.Type, isSlice bool) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedHeaderParamType(field, t.Underlying(), false)
	case *types.Pointer:
		return isSupportedHeaderParamType(field, t.Elem(), false)
	case *types.Basic:
		kind := t.Kind()
		if kind == types.Bool {
			return fmt.Errorf("header parameters cannot be of type boolean: %s", field.Name())
		}
		if typesutil.IsComplex(kind) {
			return fmt.Errorf("header parameters cannot be of type complex: %s", field.Name())
		}
		if typesutil.IsFloat(kind) {
			return fmt.Errorf("header parameters cannot be of type float: %s", field.Name())
		}
	case *types.Named:
		name := t.String()
		if name == _timeTypeName {
			return nil
		}
		return isSupportedHeaderParamType(field, t.Underlying(), false)
	case *types.Slice:
		if isSlice {
			return fmt.Errorf("header parameters cannot be nested slices: %s", field.Name())
		}
		return isSupportedHeaderParamType(field, t.Elem(), true)
	default:
		return fmt.Errorf("type `%s` is not supported for header parameters", typ.String())
	}
	return nil
}

func isSupportedCookieType(field *types.Var, typ types.Type) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedCookieType(field, t.Underlying())
	case *types.Pointer:
		return isSupportedCookieType(field, t.Elem())
	case *types.Named:
		name := t.String()
		if name == _cookieTypeName {
			return nil
		}
	default:
		return fmt.Errorf("type `%s` is not supported for cookie parameters. Use http.Cookie", typ.String())
	}
	return nil
}

func isSupportedQueryType(field *types.Var, typ types.Type, isSlice bool) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedQueryType(field, t.Underlying(), false)
	case *types.Pointer:
		return isSupportedQueryType(field, t.Elem(), false)
	case *types.Basic:
		kind := t.Kind()
		if typesutil.IsFloat(kind) {
			return fmt.Errorf("query parameters cannot be of type float: %s", field.Name())
		}
		if typesutil.IsComplex(kind) {
			return fmt.Errorf("query parameters cannot be of type complex: %s", field.Name())
		}
	case *types.Named:
		name := t.String()
		if name == _timeTypeName {
			return nil
		}
		return isSupportedQueryType(field, t.Underlying(), false)
	case *types.Slice:
		if isSlice {
			return fmt.Errorf("query parameters cannot be nested slices: %s", field.Name())
		}
		return isSupportedQueryType(field, t.Elem(), true)
	case *types.Map:
		isKeyTypeBasic := typesutil.IsBasic(t.Key(), true)
		isValTypeBasic := typesutil.IsBasic(t.Elem(), true)
		if !isKeyTypeBasic || !isValTypeBasic {
			return fmt.Errorf("map types for query parameters can only be of form map[string]string")
		}
		key := t.Key().(*types.Basic)
		if key.Kind() != types.String {
			return fmt.Errorf("the map key type of a query parameters has to be string")
		}
		val := t.Elem().(*types.Basic)
		if val.Kind() != types.String {
			return fmt.Errorf("the map value type of a query parameters has to be string")
		}
	case *types.Struct:
		for field := range t.Fields() {
			err := isSupportedQueryType(field, field.Type(), false)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("type `%s` is not supported for query parameters", typ.String())
	}
	return nil
}
