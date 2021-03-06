package types

var _ Type = (*Interface)(nil)

// Interface is the interface type model definition.
type Interface struct {
	Pkg           *Package
	Comment       string
	InterfaceName string
	Methods       []Function
}

// Name implements Type interface.
func (i Interface) Name(identified bool, packageContext string) string {
	if identified && packageContext != i.Pkg.Path {
		if identifier := i.Pkg.Identifier; identifier != "" {
			return identifier + "." + i.InterfaceName
		}
	}
	return i.InterfaceName
}

// FullName implements Type interface.
func (i Interface) FullName() string {
	return i.Pkg.Path + "/" + i.InterfaceName
}

// Package implements Type interface
func (i *Interface) Package() *Package {
	return i.Pkg
}

// Kind implements Type interface.
func (i *Interface) Kind() Kind {
	return KindInterface
}

// Elem implements Type interface.
func (i *Interface) Elem() Type {
	return nil
}

// KindString implements Type interface.
func (i Interface) String() string {
	return i.Name(true, "")
}

// Zero implements Type interface.
func (i Interface) Zero(_ bool, _ string) string {
	return "nil"
}

// Implements checks if given interface implements another interface.
func (i *Interface) Implements(another *Interface) bool {
	return Implements(i, another)
}

// Equal implements Type interface.
func (i *Interface) Equal(another Type) bool {
	it, ok := another.(*Interface)
	if !ok {
		return false
	}
	return it.Pkg == i.Pkg && it.InterfaceName == i.InterfaceName
}

// Implements checks if the type t implements interface 'interfaceType'.
func Implements(t Type, interfaceType *Interface) bool {
	var (
		s         *Struct
		isPointer bool
	)
	for s == nil {
		switch tt := t.(type) {
		case *Pointer:
			isPointer = true
			t = tt.PointedType
		case *Alias:
			t = tt.Type
		case *Struct:
			s = tt
		default:
			return false
		}
	}
	return s.Implements(interfaceType, isPointer)
}

func implements(interfaceToImplement *Interface, implementer methoder, pointer bool) bool {
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
	getMethods() []Function
}
