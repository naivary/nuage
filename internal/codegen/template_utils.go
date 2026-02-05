package codegen

import (
	"slices"
	"strconv"
	"text/template"

	"github.com/naivary/nuage/openapi"
)

var FuncsMap = template.FuncMap{
	"BitSize":             bitSize,
	"Dict":                dict,
	"IsString":            isString,
	"IsBasic":             isBasic,
	"IsInteger":           isInteger,
	"ElemType":            elemType,
	"IsQueryParamDefined": isQueryParamDefined,
}

func bitSize(typ string) int {
	switch typ {
	case "int8", "uint8":
		return 8
	case "int16", "uint16":
		return 16
	case "int32", "uint32", "float32":
		return 32
	case "int64", "uint64", "float64":
		return 64
	default:
		return strconv.IntSize
	}
}

func dict(pairs ...any) map[string]any {
	d := make(map[string]any, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		d[pairs[i].(string)] = pairs[i+1]
	}
	return d
}

func isString(info *typeInfo) bool {
	switch info.Kind {
	case "string":
		return true
	case kindPtr:
		return slices.ContainsFunc(info.Children, isString)
	case kindNamed:
		return slices.ContainsFunc(info.Children, isString)
	default:
		return false
	}
}

func isBasic(info *typeInfo) bool {
	switch info.Kind {
	case kindPtr, kindMap, kindSlice, kindStruct, kindNamed:
		return false
	default:
		return true
	}
}

func isInteger(kind string) bool {
	switch kind {
	case "int", "int8", "int16", "int32", "int64":
		return true
	case "uint", "uin8", "uin16", "uin32", "uin64":
		return true
	default:
		return false
	}
}

func elemType(info *typeInfo) string {
	t := ""
	if info.Kind == kindPtr {
		t += "*"
		info = info.Children[0]
	}
	if info.Kind == kindNamed {
		t += info.Ident
	}
	if isBasic(info) {
		t += info.Kind
	}
	return t
}

func isQueryParamDefined(params []*parameter) bool {
	return slices.ContainsFunc(params, func(p *parameter) bool {
		return p.In == openapi.ParamInQuery
	})
}
