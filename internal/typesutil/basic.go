package typesutil

import "go/types"

func IsComplex(kind types.BasicKind) bool {
	return kind == types.Complex128 || kind == types.Complex64
}
