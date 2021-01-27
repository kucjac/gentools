package astreflect

var _ Type = &InterfaceType{}

// InterfaceType is the interface type model definition.
type InterfaceType struct {
	PackagePath   PkgPath
	Comment       string
	InterfaceName string
	Methods       []FunctionType
}

func (i *InterfaceType) Name(identified bool) string {
	if identified {
		if identifier := i.PackagePath.Identifier(); identifier != "" {
			return identifier + "." + i.InterfaceName
		}
	}
	return i.InterfaceName
}

func (i *InterfaceType) FullName() string {
	return string(i.PackagePath) + "/" + i.InterfaceName
}

func (i *InterfaceType) PkgPath() PkgPath {
	return i.PackagePath
}

func (i *InterfaceType) Kind() Kind {
	return Interface
}

func (i *InterfaceType) Elem() Type {
	return nil
}

// Receiver is the function (method) receiver with the name and a pointer flag.
type Receiver struct {
	Name    string
	Pointer bool
}

// IOParam is the input/output parameter of functions and methods.
type IOParam struct {
	Name string
	Type Type
}
