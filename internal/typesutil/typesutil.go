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
	if deref {
		typ = underlying(typ)
	}
	_, isBasic := typ.(*types.Basic)
	return isBasic
}

func IsSlice(typ types.Type, deref bool) bool {
	if deref {
		typ = underlying(typ)
	}
	_, isSlice := typ.(*types.Slice)
	return isSlice
}

func IsStruct(typ types.Type, deref bool) bool {
	if deref {
		typ = underlying(typ)
	}
	_, isStruct := typ.(*types.Struct)
	return isStruct
}

func IsMap(typ types.Type, deref bool) bool {
	if deref {
		typ = underlying(typ)
	}
	_, isMap := typ.(*types.Map)
	return isMap
}

func IsPointer(typ types.Type) bool {
	_, isPtr := typ.(*types.Pointer)
	return isPtr
}

func underlying(typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Pointer:
		return underlying(t.Elem())
	case *types.Alias:
		return underlying(t.Rhs())
	case *types.Named:
		return underlying(t.Underlying())
	default:
		return t
	}
}
