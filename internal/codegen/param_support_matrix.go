package codegen

import (
	"fmt"
	"go/types"

	"github.com/naivary/nuage/internal/openapiutil"
	"github.com/naivary/nuage/internal/typesutil"
	"github.com/naivary/nuage/openapi"
)

func isSupportedParamType(opts *openapiutil.ParamOpts, typ types.Type) bool {
	switch opts.In {
	case openapi.ParamInPath:
		return isSupportedPathParamType(typ)
	case openapi.ParamInHeader:
		return isSupportedHeaderParamType(typ)
	case openapi.ParamInQuery:
		return isSupportedQueryParamType(opts, typ)
	case openapi.ParamInCookie:
		return isSupportedCookieParamType(typ)
	}
	return false
}

func isSupportedPathParamType(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Pointer:
		return isSupportedPathParamType(t.Elem())
	case *types.Named:
		return isSupportedPathParamType(t.Underlying())
	case *types.Basic:
		kind := t.Kind()
		return typesutil.IsInt(kind) ||
			typesutil.IsUint(kind) ||
			typesutil.IsString(kind)
	default:
		return false
	}
}

func isSupportedCookieParamType(typ types.Type) bool {
	ptr, isPtr := typ.(*types.Pointer)
	if !isPtr {
		return false
	}
	typ = ptr.Elem()
	named, isNamed := typ.(*types.Named)
	if !isNamed {
		return false
	}
	fqpn := fmt.Sprintf("%s.%s", named.Obj().Pkg().Name(), named.Obj().Name())
	return fqpn == "http.Cookie"
}

func isSupportedHeaderParamType(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Pointer:
		return isSupportedHeaderParamType(t.Elem())
	case *types.Named:
		fqpn := fmt.Sprintf("%s.%s", t.Obj().Pkg(), t.Obj().Name())
		if fqpn == "time.Time" {
			return true
		}
		return isSupportedHeaderParamType(t.Underlying())
	case *types.Basic:
		return isSupportedHeaderParamBasicType(t)
	default:
		return false
	}
}

func isSupportedHeaderParamBasicType(typ types.Type) bool {
	basic, isBasic := typ.(*types.Basic)
	if !isBasic {
		return false
	}
	kind := basic.Kind()
	return typesutil.IsInt(kind) ||
		typesutil.IsUint(kind) ||
		typesutil.IsString(kind) ||
		typesutil.IsBool(kind)
}

func isSupportedQueryParamType(opts *openapiutil.ParamOpts, typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Pointer:
		return isSupportedQueryParamType(opts, t.Elem())
	case *types.Named:
		fqpn := fmt.Sprintf("%s.%s", t.Obj().Pkg(), t.Obj().Name())
		if fqpn == "time.Time" {
			return true
		}
		return isSupportedHeaderParamType(t.Underlying())
	case *types.Basic:
		return isSupportedQueryParamBasicType(t)
	case *types.Slice:
		elem := typesutil.Underlying(t.Elem())
		return isSupportedQueryParamBasicType(elem)
	case *types.Map:
		key := typesutil.Underlying(t.Key())
		val := typesutil.Underlying(t.Elem())
		basicKey, isKeyBasic := key.(*types.Basic)
		basicVal, isValBasic := val.(*types.Basic)
		if !isKeyBasic || !isValBasic {
			return false
		}
		return typesutil.IsString(basicKey.Kind()) &&
			typesutil.IsString(basicVal.Kind())
	}
	return true
}

func isSupportedQueryParamBasicType(typ types.Type) bool {
	basic, isBasic := typ.(*types.Basic)
	if !isBasic {
		return false
	}
	kind := basic.Kind()
	return typesutil.IsInt(kind) ||
		typesutil.IsUint(kind) ||
		typesutil.IsString(kind) ||
		typesutil.IsBool(kind)
}
