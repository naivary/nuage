package codegen

import (
	"errors"
	"go/types"

	"github.com/naivary/nuage/internal/typesutil"
	"github.com/naivary/nuage/openapi"
)

const (
	_timeTypeName   = "time.Time"
	_cookieTypeName = "http.Cookie"
)

func isSupportedParamType(in openapi.ParamIn, typ types.Type) error {
	switch in {
	case openapi.ParamInPath:
		return isSupportedPathParamType(typ)
	case openapi.ParamInHeader:
		return isSupportedHeaderParamType(typ)
	case openapi.ParamInQuery:
		return isSupportedQueryType(typ)
	case openapi.ParamInCookie:
		return isSupportedCookieType(typ)
	default:
		return nil
	}
}

func isSupportedPathParamType(typ types.Type) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedPathParamType(t.Underlying())
	case *types.Pointer:
		return isSupportedPathParamType(t.Elem())
	case *types.Basic:
		kind := t.Kind()
		if kind == types.Bool {
			return errors.New("path parameters cannot be of type boolean")
		}
		if typesutil.IsComplex(kind) {
			return errors.New("path parameters cannot be of type complex")
		}
		if typesutil.IsFloat(kind) {
			return errors.New("path parameters cannot be of type float")
		}
	case *types.Named:
		return isSupportedPathParamType(t.Underlying())
	case *types.Slice:
		if !typesutil.IsBasic(t.Elem(), true) {
			return errors.New("slices can only be of basic type")
		}
	default:
		return errors.New("type is not supported for path parameter")
	}
	return nil
}

func isSupportedHeaderParamType(typ types.Type) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedHeaderParamType(t.Underlying())
	case *types.Pointer:
		return isSupportedHeaderParamType(t.Elem())
	case *types.Basic:
		kind := t.Kind()
		if kind == types.Bool {
			return errors.New("header parameters cannot be of type boolean")
		}
		if typesutil.IsComplex(kind) {
			return errors.New("header parameters cannot be of type complex")
		}
		if typesutil.IsFloat(kind) {
			return errors.New("header parameters cannot be of type float")
		}
	case *types.Named:
		name := t.String()
		if name == _timeTypeName {
			return nil
		}
		return isSupportedHeaderParamType(t.Underlying())
	case *types.Slice:
		if typesutil.IsSlice(t.Elem(), true) {
			return errors.New("header parameters cannot be nested slices")
		}
		return isSupportedHeaderParamType(t.Elem())
	default:
		return errors.New("type is not supported for header parameter")
	}
	return nil
}

func isSupportedCookieType(typ types.Type) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedCookieType(t.Underlying())
	case *types.Pointer:
		return isSupportedCookieType(t.Elem())
	case *types.Named:
		name := t.String()
		if name == _cookieTypeName {
			return nil
		}
	default:
		return errors.New("type is not supported for cookie parameter. Use http.Cookie")
	}
	return nil
}

func isSupportedQueryType(typ types.Type) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedQueryType(t.Underlying())
	case *types.Pointer:
		return isSupportedQueryType(t.Elem())
	case *types.Basic:
		kind := t.Kind()
		if typesutil.IsFloat(kind) {
			return errors.New("query parameters cannot be of type float")
		}
		if typesutil.IsComplex(kind) {
			return errors.New("query parameters cannot be of type complex")
		}
	case *types.Named:
		name := t.String()
		if name == _timeTypeName {
			return nil
		}
		return isSupportedQueryType(t.Underlying())
	case *types.Slice:
		if !typesutil.IsBasic(t.Elem(), true) {
			return errors.New("slices can only be of basic type")
		}
	case *types.Map:
		isKeyTypeBasic := typesutil.IsBasic(t.Key(), true)
		isValTypeBasic := typesutil.IsBasic(t.Elem(), true)
		if !isKeyTypeBasic || !isValTypeBasic {
			return errors.New("map types for query parameters can only be of type map[string]string")
		}
		key := t.Key().(*types.Basic)
		if key.Kind() != types.String {
			return errors.New("the map key type of a query parameters has to be string")
		}
		val := t.Elem().(*types.Basic)
		if val.Kind() != types.String {
			return errors.New("the map value type of a query parameters has to be string")
		}
	case *types.Struct:
		for field := range t.Fields() {
			if !typesutil.IsBasic(field.Type(), true) {
				return errors.New("when using a struct as a query parameter only primitive types are allowed to use")
			}
		}
	default:
		return errors.New("type is not supported for query parameter")
	}
	return nil
}
