package typesutil

import (
	"go/types"
)

func IsComplex(kind types.BasicKind) bool {
	return kind == types.Complex128 || kind == types.Complex64
}

func IsFloat(kind types.BasicKind) bool {
	return kind == types.Float64 || kind == types.Float32
}

func IsInt(kind types.BasicKind) bool {
	switch kind {
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
		return true
	default:
		return false
	}
}

func IsUint(kind types.BasicKind) bool {
	switch kind {
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		return true
	default:
		return false
	}
}

func IsString(kind types.BasicKind) bool {
	return kind == types.String
}

func IsBool(kind types.BasicKind) bool {
	return kind == types.Bool
}

// IsBasic reports wheret `typ` is basic. If `deref` is true
// it will dereference Alias, Pointers and Named types.
func IsBasic(typ types.Type, deref bool) bool {
	if deref {
		typ = Underlying(typ)
	}
	_, isBasic := typ.(*types.Basic)
	return isBasic
}

func IsSlice(typ types.Type, deref bool) bool {
	if deref {
		typ = Underlying(typ)
	}
	_, isSlice := typ.(*types.Slice)
	return isSlice
}

func IsStruct(typ types.Type, deref bool) bool {
	if deref {
		typ = Underlying(typ)
	}
	_, isStruct := typ.(*types.Struct)
	return isStruct
}

func IsMap(typ types.Type, deref bool) bool {
	if deref {
		typ = Underlying(typ)
	}
	_, isMap := typ.(*types.Map)
	return isMap
}

func IsPointer(typ types.Type) bool {
	_, isPtr := typ.(*types.Pointer)
	return isPtr
}

func Deref(typ types.Type) types.Type {
	if IsPointer(typ) {
		return typ.(*types.Pointer).Elem()
	}
	return typ
}

func Underlying(typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Pointer:
		return Underlying(t.Elem())
	case *types.Alias:
		return Underlying(t.Rhs())
	case *types.Named:
		return Underlying(t.Underlying())
	default:
		return t
	}
}
