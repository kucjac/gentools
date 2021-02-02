package types

import (
	"strings"
)

var _ Type = (*Function)(nil)

// Function is the function type used for getting.
type Function struct {
	Comment  string
	Pkg      *Package
	Receiver *Receiver
	FuncName string
	In       []FuncParam
	Out      []FuncParam
	Variadic bool
}

// Name implements Type interface.
func (f Function) Name(identified bool, packageContext string) string {
	if identified && packageContext != f.Pkg.Path {
		if i := f.Pkg.Identifier; i != "" {
			return i + "." + f.FuncName
		}
	}
	return f.FuncName
}

// FullName implements Type interface.
func (f *Function) FullName() string {
	return f.Pkg.Path + "/" + f.FuncName
}

// Package implements Type interface.
func (f *Function) Package() *Package {
	return f.Pkg
}

// Kind implements Type interface.
func (f *Function) Kind() Kind {
	return KindFunc
}

// Elem implements Type interface.
func (f *Function) Elem() Type {
	return nil
}

// KindString implements Type interface.
func (f Function) String() string {
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

// Zero implements Type interface.
func (f *Function) Zero(_ bool, _ string) string {
	return "nil"
}

// Equal implements Type interface.
func (f *Function) Equal(another Type) bool {
	ft, ok := another.(*Function)
	if !ok {
		return false
	}
	if (f.Receiver == nil && ft.Receiver != nil) || (ft.Receiver == nil && f.Receiver != nil) {
		return false
	}
	return f.Pkg == ft.Pkg && f.FuncName == ft.FuncName
}

// FuncParam is the input/output parameter of functions and methods.
type FuncParam struct {
	Name string
	Type Type
}

// KindString implements fmt.Stringer interface.
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
	_, ok := r.Type.(*Pointer)
	return ok
}
