package openapi

import "reflect"

const (
	_tagKeyPathParam   = "path"
	_tagKeyQueryParam  = "query"
	_tagKeyHeaderParam = "header"
	_tagKeyCookieParam = "cookie"
)

func ParamLocation(tag reflect.StructTag) ParamIn {
	if _, ok := tag.Lookup(_tagKeyPathParam); ok {
		return ParamInPath
	}
	if _, ok := tag.Lookup(_tagKeyQueryParam); ok {
		return ParamInQuery
	}
	if _, ok := tag.Lookup(_tagKeyHeaderParam); ok {
		return ParamInHeader
	}
	if _, ok := tag.Lookup(_tagKeyCookieParam); ok {
		return ParamInCookie
	}
	return ""
}
