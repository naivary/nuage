package codegen

import (
	"strconv"
	"strings"
	"text/template"
)

var FuncsMap = template.FuncMap{
	"BitSize":    bitSize,
	"Capitalize": capitalize,
	"Dict":       dict,
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

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	first := strings.ToUpper(string(s[0]))
	return first + s[1:]
}

func dict(pairs ...any) map[string]any {
	d := make(map[string]any, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		d[pairs[i].(string)] = pairs[i+1]
	}
	return d
}
