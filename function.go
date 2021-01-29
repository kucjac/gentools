package astreflect

import (
	"strings"
)

var _ Type = FunctionType{}

// FunctionType is the function type used for getting.
type FunctionType struct {
	Comment  string
	Pkg      *Package
	Receiver *Receiver
	FuncName string
	In       []FuncParam
	Out      []FuncParam
	Variadic bool
}

// Name implements Type interface.
func (f FunctionType) Name(identified bool, packageContext string) string {
	if identified && packageContext != f.Pkg.Path {
		if i := f.Pkg.Identifier; i != "" {
			return i + "." + f.FuncName
		}
	}
	return f.FuncName
}

// FullName implements Type interface.
func (f FunctionType) FullName() string {
	return f.Pkg.Path + "/" + f.FuncName
}

// Package implements Type interface.
func (f FunctionType) Package() *Package {
	return f.Pkg
}

// Kind implements Type interface.
func (f FunctionType) Kind() Kind {
	return Func
}

// Elem implements Type interface.
func (f FunctionType) Elem() Type {
	return nil
}

// String implements Type interface.
func (f FunctionType) String() string {
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
func (f FunctionType) Zero(_ bool, _ string) string {
	return "nil"
}

// Equal implements Type interface.
func (f FunctionType) Equal(another Type) bool {
	var ft FunctionType
	switch t := another.(type) {
	case *FunctionType:
		ft = *t
	case FunctionType:
		ft = t
	default:
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
