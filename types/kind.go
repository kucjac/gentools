package types

// A Kind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type Kind uint

// KindString implements fmt.Stringer interface.
func (k Kind) String() string {
	name, ok := kindNameMap[k]
	if !ok {
		return "Invalid"
	}
	return name
}

// IsNumber checks if given kind is of number (integers, floats)
func (k Kind) IsNumber() bool {
	return k >= KindInt && k <= KindFloat64
}

// IsBuiltin checks if given kind is
func (k Kind) IsBuiltin() bool {
	return k >= KindBool && k <= KindString
}

// BuiltInName gets the name of the builtin kind.
func (k Kind) BuiltInName() string {
	if k == Invalid {
		return "_INVALID_KIND_"
	}
	if !k.IsBuiltin() {
		return "_NOT_BUILT_IN_"
	}
	return builtInNames[k-1]
}

// Enumerated kind representations.
const (
	Invalid Kind = iota
	KindBool
	KindInt
	KindInt8
	KindInt16
	KindInt32
	KindInt64
	KindUint
	KindUint8
	KindUint16
	KindUint32
	KindUint64
	KindUintptr
	KindFloat32
	KindFloat64
	KindComplex64
	KindComplex128
	KindString
	KindArray
	KindChan
	KindFunc
	KindInterface
	KindMap
	KindPtr
	KindSlice
	KindStruct
	KindUnsafePointer
	KindWrapper
)

var stdKindMap = map[string]Kind{"int": KindInt, "int8": KindInt8, "int16": KindInt16, "int32": KindInt32, "int64": KindInt64, "uint": KindUint, "uint8": KindUint8, "uint16": KindUint16, "uint32": KindUint32, "uint64": KindUint64, "float32": KindFloat32, "float64": KindFloat64, "string": KindString, "bool": KindBool, "uintptr": KindUintptr, "complex64": KindComplex64, "complex128": KindComplex128}

var builtInNames = [KindString]string{"bool", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64", "complex64", "complex128", "string"}

var kindNameMap = map[Kind]string{Invalid: "Invalid", KindBool: "Bool", KindInt: "Int", KindInt8: "Int8", KindInt16: "Int16", KindInt32: "Int32", KindInt64: "Int64", KindUint: "Uint", KindUint8: "Uint8", KindUint16: "Uint16", KindUint32: "Uint32", KindUint64: "Uint64", KindUintptr: "Uintptr", KindFloat32: "Float32", KindFloat64: "Float64", KindComplex64: "Complex64", KindComplex128: "Complex128", KindArray: "Array", KindChan: "Chan", KindFunc: "Func", KindInterface: "Interface", KindMap: "Map", KindPtr: "Ptr", KindSlice: "Slice", KindString: "String", KindStruct: "Struct", KindUnsafePointer: "UnsafePointer", KindWrapper: "Wrapper"}
