package typesutil

import "go/types"

func IsComplex(kind types.BasicKind) bool {
	return kind == types.Complex128 || kind == types.Complex64
}

func IsFloat(kind types.BasicKind) bool {
	return kind == types.Float64 || kind == types.Float32
}

// IsBasic reports wheret `typ` is basic. If `force` is true
// it will dereference Alias, Pointers and Named types.
func IsBasic(typ types.Type, force bool) bool {
	switch t := typ.(type) {
	case *types.Pointer:
		return IsBasic(t.Elem(), force)
	case *types.Alias:
		return IsBasic(t.Underlying(), force)
	case *types.Named:
		return IsBasic(t.Underlying(), force)
	case *types.Basic:
		return true
	default:
		return false
	}
}
