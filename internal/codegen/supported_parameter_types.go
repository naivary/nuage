package codegen

import (
	"go/types"

	"github.com/naivary/nuage/internal/openapiutil"
	"github.com/naivary/nuage/internal/typesutil"
	"github.com/naivary/nuage/openapi"
)

const (
	_timeTypeName   = "time.Time"
	_cookieTypeName = "http.Cookie"
)

func isSupportedParamType(opts *openapiutil.ParamOpts, typ types.Type) bool {
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

func isSupportedPathParamType(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Pointer:
		return isSupportedPathParamType(t.Elem())
	case *types.Named:
		return isSupportedPathParamType(t.Underlying())
	case *types.Basic:
		kind := t.Kind()
		return typesutil.IsInt(kind) || typesutil.IsUint(kind) || typesutil.IsString(kind)
	default:
		return false
	}
}
