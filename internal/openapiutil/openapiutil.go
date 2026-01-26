package openapiutil

import (
	"reflect"

	"github.com/naivary/nuage/openapi"
)

func ParamLocation(tag reflect.StructTag) openapi.ParamIn {
	if _, ok := tag.Lookup(openapi.TagKeyPathParam); ok {
		return openapi.ParamInPath
	}
	if _, ok := tag.Lookup(openapi.TagKeyQueryParam); ok {
		return openapi.ParamInQuery
	}
	if _, ok := tag.Lookup(openapi.TagKeyHeaderParam); ok {
		return openapi.ParamInHeader
	}
	if _, ok := tag.Lookup(openapi.TagKeyCookieParam); ok {
		return openapi.ParamInCookie
	}
	return ""
}
