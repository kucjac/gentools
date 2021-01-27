package astreflect

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
			Out:      []IOParam{{Type: stringType}},
		}},
	}
	builtIn.Interfaces = append(builtIn.Interfaces, er)
	builtIn.Types["error"] = er
}

// A Kind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type Kind uint

func (k Kind) String() string {
	name, ok := kindNameMap[k]
	if !ok {
		return "Invalid"
	}
	return name
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
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
	Wrapper
)

var stdKindMap = map[string]Kind{
	"int":     Int,
	"int8":    Int8,
	"int16":   Int16,
	"int32":   Int32,
	"int64":   Int64,
	"uint":    Uint,
	"uint8":   Uint8,
	"uint16":  Uint16,
	"uint32":  Uint32,
	"uint64":  Uint64,
	"float32": Float32,
	"float64": Float64,
	"string":  String,
	"bool":    Bool,
	"uintptr": Uintptr,
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

var _ Type = &BuiltInType{}

// BuiltInType is the built in type definition.
type BuiltInType struct {
	TypeName string
	StdKind  Kind
}

// FullName implements Type interface.
func (s *BuiltInType) FullName() string {
	return s.TypeName
}

// Name implements Type interface.
func (s *BuiltInType) Name(_ bool) string {
	return s.TypeName
}

// PkgPath implements Type interface.
func (s *BuiltInType) PkgPath() PkgPath {
	return builtInPkgPath
}

// Kind implements Type interface.
func (s *BuiltInType) Kind() Kind {
	return s.StdKind
}

func (s *BuiltInType) Elem() Type {
	return nil
}

type builtInPackage struct {
	*Package
	BuiltInTypes []BuiltInType
}
