package astreflect

var _ Type = (*InterfaceType)(nil)

// InterfaceType is the interface type model definition.
type InterfaceType struct {
	Pkg           *Package
	Comment       string
	InterfaceName string
	Methods       []FunctionType
}

// Name implements Type interface.
func (i InterfaceType) Name(identified bool, packageContext string) string {
	if identified && packageContext != i.Pkg.Path {
		if identifier := i.Pkg.Identifier; identifier != "" {
			return identifier + "." + i.InterfaceName
		}
	}
	return i.InterfaceName
}

// FullName implements Type interface.
func (i InterfaceType) FullName() string {
	return i.Pkg.Path + "/" + i.InterfaceName
}

// Package implements Type interface
func (i InterfaceType) Package() *Package {
	return i.Pkg
}

// Kind implements Type interface.
func (i InterfaceType) Kind() Kind {
	return Interface
}

// Elem implements Type interface.
func (i InterfaceType) Elem() Type {
	return nil
}

// String implements Type interface.
func (i InterfaceType) String() string {
	return i.Name(true, "")
}

// Zero implements Type interface.
func (i InterfaceType) Zero(_ bool, _ string) string {
	return "nil"
}

// Implements checks if given interface implements another interface.
func (i InterfaceType) Implements(another *InterfaceType) bool {
	return Implements(i, another)
}

// Equal implements Type interface.
func (i InterfaceType) Equal(another Type) bool {
	var it InterfaceType
	switch t := another.(type) {
	case *InterfaceType:
		it = *t
	case InterfaceType:
		it = t
	default:
		return false
	}
	return it.Pkg == i.Pkg && it.InterfaceName == i.InterfaceName
}

// Implements checks if the type t implements interface 'interfaceType'.
func Implements(t Type, interfaceType *InterfaceType) bool {
	var (
		s         *StructType
		isPointer bool
	)
	for s == nil {
		switch tt := t.(type) {
		case *PointerType:
			isPointer = true
			t = tt.PointedType
		case *WrappedType:
			t = tt.Type
		case *StructType:
			s = tt
		default:
			return false
		}
	}
	return s.Implements(interfaceType, isPointer)
}

func implements(interfaceToImplement *InterfaceType, implementer methoder, pointer bool) bool {
	implMethods := implementer.getMethods()
	if len(interfaceToImplement.Methods) > len(implMethods) {
		return false
	}
	for _, iMethod := range interfaceToImplement.Methods {
		var found bool
		for _, sMethod := range implMethods {
			if sMethod.FuncName == iMethod.FuncName {
				if len(iMethod.In) != len(sMethod.In) {
					return false
				}
				if len(iMethod.Out) != len(sMethod.Out) {
					return false
				}
				if sMethod.Variadic != iMethod.Variadic {
					return false
				}
				if sMethod.Receiver.IsPointer() && !pointer {
					return false
				}

				for i := 0; i < len(iMethod.In); i++ {
					if iMethod.In[i].Type.FullName() != sMethod.In[i].Type.FullName() {
						return false
					}
				}
				for i := 0; i < len(iMethod.Out); i++ {
					if iMethod.Out[i].Type.FullName() != sMethod.Out[i].Type.FullName() {
						return false
					}
				}
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type methoder interface {
	getMethods() []FunctionType
}
