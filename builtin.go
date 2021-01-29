package astreflect

import (
	"fmt"
)

const (
	builtInPkgName = "builtin"
	builtInPkgPath = PkgPath(builtInPkgName)
)

var builtIn *builtInPackage

func init() {
	builtIn = &builtInPackage{
		Package: &Package{
			Path:            builtInPkgPath,
			Identifier:      "",
			Types:           map[string]Type{},
			typesInProgress: map[string]Type{},
		},
	}

	for name, kind := range stdKindMap {
		st := BuiltInType{
			TypeName: name,
			StdKind:  kind,
		}
		builtIn.BuiltInTypes = append(builtIn.BuiltInTypes, st)
		builtIn.Types[name] = &st
		switch kind {
		case Uint8:
			wt := &WrappedType{
				PackagePath: builtInPkgPath,
				WrapperName: "byte",
				Type:        &st,
			}
			builtIn.WrappedTypes = append(builtIn.WrappedTypes, wt)
			builtIn.Types["byte"] = wt
		case Int32:
			wt := &WrappedType{
				PackagePath: builtInPkgPath,
				WrapperName: "rune",
				Type:        &st,
			}
			builtIn.WrappedTypes = append(builtIn.WrappedTypes, wt)
			builtIn.Types["rune"] = wt
		}
	}
	stringType, _ := GetBuiltInType("string")
	er := &InterfaceType{
		PackagePath:   builtInPkgPath,
		InterfaceName: "error",
		Methods: []FunctionType{{
			FuncName: "Error",
			Out:      []FuncParam{{Type: stringType}},
		}},
	}
	builtIn.Interfaces = append(builtIn.Interfaces, er)
	builtIn.Types["error"] = er
}

// A Kind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type Kind uint

// String implements fmt.Stringer interface.
func (k Kind) String() string {
	name, ok := kindNameMap[k]
	if !ok {
		return "Invalid"
	}
	return name
}

// IsNumber checks if given kind is of number (integers, floats)
func (k Kind) IsNumber() bool {
	return k >= Int && k <= Float64
}

// IsBuiltin checks if given kind is
func (k Kind) IsBuiltin() bool {
	return k >= Bool && k <= String
}

// Enumerated kind representations.
const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	Struct
	UnsafePointer
	Wrapper
)

var stdKindMap = map[string]Kind{
	"int":        Int,
	"int8":       Int8,
	"int16":      Int16,
	"int32":      Int32,
	"int64":      Int64,
	"uint":       Uint,
	"uint8":      Uint8,
	"uint16":     Uint16,
	"uint32":     Uint32,
	"uint64":     Uint64,
	"float32":    Float32,
	"float64":    Float64,
	"string":     String,
	"bool":       Bool,
	"uintptr":    Uintptr,
	"complex64":  Complex64,
	"complex128": Complex128,
}

var kindNameMap = map[Kind]string{
	Invalid:       "Invalid",
	Bool:          "Bool",
	Int:           "Int",
	Int8:          "Int8",
	Int16:         "Int16",
	Int32:         "Int32",
	Int64:         "Int64",
	Uint:          "Uint",
	Uint8:         "Uint8",
	Uint16:        "Uint16",
	Uint32:        "Uint32",
	Uint64:        "Uint64",
	Uintptr:       "Uintptr",
	Float32:       "Float32",
	Float64:       "Float64",
	Complex64:     "Complex64",
	Complex128:    "Complex128",
	Array:         "Array",
	Chan:          "Chan",
	Func:          "Func",
	Interface:     "Interface",
	Map:           "Map",
	Ptr:           "Ptr",
	Slice:         "Slice",
	String:        "String",
	Struct:        "Struct",
	UnsafePointer: "UnsafePointer",
	Wrapper:       "Wrapper",
}

// IsKindBuiltIn checks if given name is a built in type.
func IsKindBuiltIn(kindName string) (Kind, bool) {
	k, ok := stdKindMap[kindName]
	return k, ok
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
	TypeName string
	StdKind  Kind
}

// Name implements Type interface.
func (b BuiltInType) Name(_ bool, _ string) string {
	return b.TypeName
}

// FullName implements Type interface.
func (b BuiltInType) FullName() string {
	return b.TypeName
}

// PkgPath implements Type interface.
func (b BuiltInType) PkgPath() PkgPath {
	return builtInPkgPath
}

// Kind implements Type interface.
func (b BuiltInType) Kind() Kind {
	return b.StdKind
}

// Elem implements Type interface.
func (b BuiltInType) Elem() Type {
	return nil
}

// String implements Type interface.
func (b BuiltInType) String() string {
	return b.TypeName
}

func (b BuiltInType) Equal(another Type) bool {
	var bt BuiltInType
	switch t := another.(type) {
	case *BuiltInType:
		bt = *t
	case BuiltInType:
		bt = t
	default:
		return false
	}
	return b.StdKind == bt.StdKind
}

// Zero implements Type interface.
func (b BuiltInType) Zero(_ bool, _ string) string {
	if b.Kind().IsNumber() {
		return "0"
	}
	switch b.Kind() {
	case String:
		return "\"\""
	case Complex64, Complex128:
		return "complex(0,0)"
	default:
		return "nil"
	}
}

type builtInPackage struct {
	*Package
	BuiltInTypes []BuiltInType
}
