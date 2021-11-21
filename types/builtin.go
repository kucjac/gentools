package types

import (
	"fmt"
)

const (
	builtInPkgName = "builtin"
	builtInPkgPath = PkgPath(builtInPkgName)
)

var builtIn *Package

// Builtin types definitions.
var (
	Error      Type
	Byte       Type
	Rune       Type
	Bool       Type
	Int        Type
	Int8       Type
	Int16      Type
	Int32      Type
	Int64      Type
	Uint       Type
	Uint8      Type
	Uint16     Type
	Uint32     Type
	Uint64     Type
	Uintptr    Type
	Float32    Type
	Float64    Type
	Complex64  Type
	Complex128 Type
	String     Type
)

func init() {
	builtIn = NewPackage(builtInPkgName, "")
	for name, kind := range stdKindMap {
		st := &BuiltInType{BuiltInKind: kind}
		builtIn.Types[name] = st
		switch kind {
		case KindBool:
			Bool = st
		case KindInt:
			Int = st
		case KindInt8:
			Int8 = st
		case KindInt16:
			Int16 = st
		case KindInt32:
			Int32 = st
			wt := &Alias{
				Pkg:       builtIn,
				AliasName: "rune",
				Type:      st,
			}
			builtIn.Aliases = append(builtIn.Aliases, wt)
			builtIn.Types["rune"] = wt
			Rune = wt
		case KindInt64:
			Int64 = st
		case KindUint:
			Uint = st
		case KindUint8:
			Uint8 = st
			wt := &Alias{
				Pkg:       builtIn,
				AliasName: "byte",
				Type:      st,
			}
			builtIn.Aliases = append(builtIn.Aliases, wt)
			builtIn.Types["byte"] = wt
			Byte = wt
		case KindUint16:
			Uint16 = st
		case KindUint32:
			Uint32 = st
		case KindUint64:
			Uint64 = st
		case KindUintptr:
			Uintptr = st
		case KindFloat32:
			Float32 = st
		case KindFloat64:
			Float64 = st
		case KindComplex64:
			Complex64 = st
		case KindComplex128:
			Complex128 = st
		case KindString:
			String = st
		}
	}

	er := &Interface{
		Pkg:           builtIn,
		InterfaceName: "error",
		Methods: []Function{{
			FuncName: "Error",
			Out:      []FuncParam{{Type: String}},
		}},
	}

	builtIn.Interfaces = append(builtIn.Interfaces, er)
	builtIn.Types["error"] = er
	Error = er
}

// IsBuiltIn checks if given name is a built in type.
func IsBuiltIn(kindName string) (Kind, bool) {
	switch kindName {
	case "byte":
		return KindUint8, true
	case "rune":
		return KindInt32, true
	default:
		k, ok := stdKindMap[kindName]
		return k, ok
	}
}

// BuiltInOf gets the built in of specific kind.
func BuiltInOf(x interface{}) Type {
	switch xt := x.(type) {
	case Kind:
		if !xt.IsBuiltin() {
			panic(fmt.Sprintf("provided kind is not a builtin kind: %s", xt))
		}
		bt, _ := builtIn.GetType(xt.BuiltInName())
		return bt
	case string:
		if xt == "" {
			bt, _ := builtIn.GetType(KindString.BuiltInName())
			return bt
		}
		return MustGetBuiltInType(xt)
	case bool:
		bt, _ := builtIn.GetType(KindBool.BuiltInName())
		return bt
	case int:
		bt, _ := builtIn.GetType(KindInt.BuiltInName())
		return bt
	case int8:
		bt, _ := builtIn.GetType(KindInt8.BuiltInName())
		return bt
	case int16:
		bt, _ := builtIn.GetType(KindInt16.BuiltInName())
		return bt
	case int32:
		bt, _ := builtIn.GetType(KindInt32.BuiltInName())
		return bt
	case int64:
		bt, _ := builtIn.GetType(KindInt64.BuiltInName())
		return bt
	case uint:
		bt, _ := builtIn.GetType(KindUint.BuiltInName())
		return bt
	case uint8:
		bt, _ := builtIn.GetType(KindUint8.BuiltInName())
		return bt
	case uint16:
		bt, _ := builtIn.GetType(KindUint16.BuiltInName())
		return bt
	case uint32:
		bt, _ := builtIn.GetType(KindUint32.BuiltInName())
		return bt
	case uint64:
		bt, _ := builtIn.GetType(KindUint64.BuiltInName())
		return bt
	case complex64:
		bt, _ := builtIn.GetType(KindComplex64.BuiltInName())
		return bt
	case complex128:
		bt, _ := builtIn.GetType(KindComplex128.BuiltInName())
		return bt
	case Type:
		return xt
	default:
		panic(fmt.Sprintf("unknown builtin type: %T", xt))
	}
}

// GetBuiltInType gets the built in type with given name.
func GetBuiltInType(name string) (Type, bool) {
	return builtIn.GetType(name)
}

// MustGetBuiltInType gets the built in type with given name. If the type is not found the function panics.
func MustGetBuiltInType(name string) Type {
	t, ok := builtIn.GetType(name)
	if !ok {
		panic(fmt.Sprintf("builtin type: '%s' not found", name))
	}
	return t
}

var _ Type = (*BuiltInType)(nil)

// BuiltInType is the built in type definition.
type BuiltInType struct {
	BuiltInKind Kind
}

// Name implements Type interface.
func (b *BuiltInType) Name(_ bool, _ string) string {
	return b.BuiltInKind.BuiltInName()
}

// FullName implements Type interface.
func (b *BuiltInType) FullName() string {
	return b.BuiltInKind.BuiltInName()
}

// Kind implements Type interface.
func (b *BuiltInType) Kind() Kind {
	return b.BuiltInKind
}

// Elem implements Type interface.
func (b *BuiltInType) Elem() Type {
	return nil
}

// KindString implements Type interface.
func (b BuiltInType) String() string {
	return b.BuiltInKind.String()
}

// Equal checks if given built in type is equal to another Type.
func (b *BuiltInType) Equal(another Type) bool {
	bt, ok := another.(*BuiltInType)
	if !ok {
		return false
	}
	return b.BuiltInKind == bt.BuiltInKind
}

// Zero implements Type interface.
func (b BuiltInType) Zero(_ bool, _ string) string {
	if b.Kind().IsNumber() {
		return "0"
	}
	switch b.Kind() {
	case KindString:
		return "\"\""
	case KindComplex64, KindComplex128:
		return "complex(0,0)"
	default:
		return "nil"
	}
}
