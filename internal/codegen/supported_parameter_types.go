package codegen

import (
	"errors"
	"fmt"
	"go/types"

	"github.com/naivary/nuage/internal/openapiutil"
	"github.com/naivary/nuage/internal/typesutil"
	"github.com/naivary/nuage/openapi"
)

const (
	_timeTypeName   = "time.Time"
	_cookieTypeName = "http.Cookie"
)

func isSupportedParamType(opts *openapiutil.ParamOpts, typ types.Type) error {
	switch opts.In {
	case openapi.ParamInPath:
		return isSupportedPathParamType(typ)
	case openapi.ParamInHeader:
		return isSupportedHeaderParamType(typ)
	case openapi.ParamInQuery:
		return isSupportedQueryType(opts, typ)
	case openapi.ParamInCookie:
		return isSupportedCookieType(typ)
	default:
		return nil
	}
}

func isSupportedPathParamType(typ types.Type) error {
	switch t := typ.(type) {
	case *types.Pointer:
		return isSupportedPathParamType(t.Elem())
	case *types.Named:
		return isSupportedPathParamType(t.Underlying())
	case *types.Alias:
		return isSupportedPathParamType(t.Rhs())
	case *types.Basic:
		return isSupportedPathParamBasicType(t)
	case *types.Slice:
		return isSupportedPathParamBasicType(t.Elem())
	default:
		return fmt.Errorf("path parameter type not supported: %s", typ.String())
	}
}

func isSupportedPathParamBasicType(typ types.Type) error {
	isUnsupported := typesutil.IsBasicKind(
		typ,
		true,
		types.Bool, types.Float64, types.Float32, types.Complex128, types.Complex64,
	)
	if isUnsupported {
		return errors.New("path parameter type not supported")
	}
	return nil
}

func isSupportedHeaderParamType(typ types.Type) error {
	switch t := typ.(type) {
	case *types.Alias, *types.Pointer:
		return isSupportedHeaderParamType(typ.Underlying())
	case *types.Basic:
		return isSupportedHeaderParamBasicType(typ)
	case *types.Named:
		name := t.String()
		if name == _timeTypeName {
			return nil
		}
		return isSupportedHeaderParamType(t.Underlying())
	case *types.Slice:
		return isSupportedHeaderParamBasicType(t.Elem())
	default:
		return errors.New("type is not supported for header parameter")
	}
}

func isSupportedHeaderParamBasicType(typ types.Type) error {
	isUnsupported := typesutil.IsBasicKind(
		typ,
		true,
		types.Bool, types.Float64, types.Float32, types.Complex128, types.Complex64,
	)
	if isUnsupported {
		return errors.New("header parameter type not supported")
	}
	return nil
}

func isSupportedCookieType(typ types.Type) error {
	switch t := typ.(type) {
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

func isSupportedQueryType(opts *openapiutil.ParamOpts, typ types.Type) error {
	switch t := typ.(type) {
	case *types.Pointer:
		return isSupportedQueryType(opts, t.Elem())
	case *types.Alias:
		return isSupportedQueryType(opts, t.Rhs())
	case *types.Basic:
		return isSupportedQueryParamBasicType(t)
	case *types.Named:
		return isSupportedQueryParamNamedType(opts, t)
	case *types.Slice:
		return isSupportedQueryParamBasicType(t.Elem())
	case *types.Map:
		if opts.Style != openapi.ParamStyleForm {
			return errors.New("query parameter type map can only be used with style form")
		}
		isKeyTypeBasic := typesutil.IsBasic(t.Key(), true)
		isValTypeBasic := typesutil.IsBasic(t.Elem(), true)
		if !isKeyTypeBasic || !isValTypeBasic {
			return errors.New("map types for query parameters can only be of type map[~string]~string")
		}
		if !typesutil.IsBasicKind(t.Key(), true, types.String) || !typesutil.IsBasicKind(t.Elem(), true, types.String) {
			return errors.New("maps can only have ~string types as key and value")
		}
	case *types.Struct:
		if opts.Style != openapi.ParamStyleDeepObject {
			return errors.New("query parameter type struct can only be used with style deepObject")
		}
		for field := range t.Fields() {
			fieldType := field.Type()
			if typesutil.IsMap(fieldType, true) {
				return errors.New("query parameter types not supported")
			}
			err := isSupportedQueryParamBasicType(fieldType)
			if err == nil {
				continue
			}
			err = isSupportedQueryParamNamedType(opts, fieldType)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("type is not supported for query parameter")
	}
	return nil
}

// TODO: this is not working if you pass in a slice of basic kiknd. It will return nil because it is
// considered supproted by this function.
func isSupportedQueryParamBasicType(typ types.Type) error {
	isBasic := typesutil.IsBasic(typ, true)
	if !isBasic {
		return errors.New("not a basic type")
	}
	isUnsupported := typesutil.IsBasicKind(
		typ,
		true,
		types.Float64, types.Float32, types.Complex128, types.Complex64,
	)
	if isUnsupported {
		return fmt.Errorf("query parameter type not supported: %s", typ.String())
	}
	return nil
}

func isSupportedQueryParamNamedType(opts *openapiutil.ParamOpts, typ types.Type) error {
	typ = typesutil.Deref(typ)
	named, isNamed := typ.(*types.Named)
	if !isNamed {
		return errors.New("not a named type")
	}
	if named.String() == _timeTypeName {
		return nil
	}
	return isSupportedQueryType(opts, named.Underlying())
}
