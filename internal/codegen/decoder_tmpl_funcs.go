package codegen

import "strconv"

func BitSize(typ string) int {
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
