package typesutil

import "go/types"

func IsComplex(kind types.BasicKind) bool {
	return kind == types.Complex128 || kind == types.Complex64
}

func IsFloat(kind types.BasicKind) bool {
	return kind == types.Float64 || kind == types.Float32
}

// IsBasic reports wheret `typ` is basic. If `deref` is true
// it will dereference Alias, Pointers and Named types.
func IsBasic(typ types.Type, deref bool) bool {
	switch t := typ.(type) {
	case *types.Pointer:
		if !deref {
			return false
		}
		return IsBasic(t.Elem(), deref)
	case *types.Alias:
		if !deref {
			return false
		}
		return IsBasic(t.Underlying(), deref)
	case *types.Named:
		if !deref {
			return false
		}
		return IsBasic(t.Underlying(), deref)
	case *types.Basic:
		return true
	default:
		return false
	}
}

func IsSlice(typ types.Type, deref bool) bool {
	switch t := typ.(type) {
	case *types.Pointer:
		if !deref {
			return false
		}
		return IsSlice(t.Elem(), deref)
	case *types.Alias:
		if !deref {
			return false
		}
		return IsSlice(t.Underlying(), deref)
	case *types.Named:
		if !deref {
			return false
		}
		return IsSlice(t.Underlying(), deref)
	case *types.Slice:
		return true
	default:
		return false
	}
}

func IsStruct(typ types.Type, deref bool) bool {
	switch t := typ.(type) {
	case *types.Pointer:
		if !deref {
			return false
		}
		return IsStruct(t.Elem(), deref)
	case *types.Alias:
		if !deref {
			return false
		}
		return IsStruct(t.Underlying(), deref)
	case *types.Named:
		if !deref {
			return false
		}
		return IsStruct(t.Underlying(), deref)
	case *types.Struct:
		return true
	default:
		return false
	}
}

func IsPointer(typ types.Type) bool {
	_, isPtr := typ.(*types.Pointer)
	return isPtr
}
