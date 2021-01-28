package astreflect

import (
	"strings"
)

var _ Type = (*FunctionType)(nil)

// FunctionType is the function type used for getting.
type FunctionType struct {
	Comment     string
	PackagePath PkgPath
	Receiver    *Receiver
	FuncName    string
	In          []FuncParam
	Out         []FuncParam
	Variadic    bool
}

// Name implements Type interface.
func (f *FunctionType) Name(identified bool, packageContext string) string {
	if identified && packageContext != f.PackagePath.FullName() {
		if i := f.PackagePath.Identifier(); i != "" {
			return i + "." + f.FuncName
		}
	}
	return f.FuncName
}

// FullName implements Type interface.
func (f *FunctionType) FullName() string {
	return string(f.PackagePath) + "/" + f.FuncName
}

// PkgPath implements Type interface.
func (f *FunctionType) PkgPath() PkgPath {
	return f.PackagePath
}

// Kind implements Type interface.
func (f *FunctionType) Kind() Kind {
	return Func
}

// Elem implements Type interface.
func (f *FunctionType) Elem() Type {
	return nil
}

// String implements Type interface.
func (f *FunctionType) String() string {
	sb := strings.Builder{}
	sb.WriteString("func ")
	if f.Receiver != nil {
		sb.WriteString(f.Receiver.String())
		sb.WriteRune(' ')
	}
	sb.WriteString(f.FuncName)
	sb.WriteRune('(')
	for i := range f.In {
		sb.WriteString(f.In[i].String())
		if i != len(f.In)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(") ")
	for i := range f.Out {
		sb.WriteString(f.Out[i].Type.String())
		if i != len(f.Out)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

// FuncParam is the input/output parameter of functions and methods.
type FuncParam struct {
	Name string
	Type Type
}

// String implements fmt.Stringer interface.
func (f FuncParam) String() string {
	return f.Name + " " + f.Type.String()
}

// Receiver is the function (method) receiver with the name and a pointer flag.
type Receiver struct {
	Name string
	Type Type
}

func (r *Receiver) String() string {
	return "(" + r.Name + " " + r.Type.Name(false, "") + ")"
}

// IsPointer checks if the receiver type is a pointer.
func (r *Receiver) IsPointer() bool {
	_, ok := r.Type.(*PointerType)
	return ok
}
