package astreflect

var _ Type = (*InterfaceType)(nil)

// InterfaceType is the interface type model definition.
type InterfaceType struct {
	PackagePath   PkgPath
	Comment       string
	InterfaceName string
	Methods       []FunctionType
}

// Name implements Type interface.
func (i *InterfaceType) Name(identified bool, packageContext string) string {
	if identified && packageContext != i.PackagePath.FullName() {
		if identifier := i.PackagePath.Identifier(); identifier != "" {
			return identifier + "." + i.InterfaceName
		}
	}
	return i.InterfaceName
}

// FullName implements Type interface.
func (i *InterfaceType) FullName() string {
	return string(i.PackagePath) + "/" + i.InterfaceName
}

// PkgPath implements Type interface
func (i *InterfaceType) PkgPath() PkgPath {
	return i.PackagePath
}

// Kind implements Type interface.
func (i *InterfaceType) Kind() Kind {
	return Interface
}

// Elem implements Type interface.
func (i *InterfaceType) Elem() Type {
	return nil
}

// String implements Type interface.
func (i *InterfaceType) String() string {
	return i.Name(true, "")
}
