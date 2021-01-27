package astreflect

import "fmt"

var _ Type = (*ArrayType)(nil)

// ArrayType is the array or slice representing type.
type ArrayType struct {
	ArrayKind Kind // Could be either Slice or Array
	Type      Type
	ArraySize int
}

// Name implements Type interface.
func (a *ArrayType) Name(identifier bool) string {
	if a.ArrayKind == Slice {
		return "[]" + a.Type.Name(identifier)
	}
	return fmt.Sprintf("[%d]%s", a.ArraySize, a.Type.Name(identifier))
}

// FullName implements Type interface.
func (a *ArrayType) FullName() string {
	if a.ArrayKind == Slice {
		return "[]" + a.Type.FullName()
	}
	return fmt.Sprintf("[%d]%s", a.ArraySize, a.Type.FullName())
}

// PkgPath implements Type interface.
func (a *ArrayType) PkgPath() PkgPath {
	return builtInPkgPath
}

// Kind implements Type interface.
func (a *ArrayType) Kind() Kind {
	return a.ArrayKind
}

// Elem implements Type interface.
func (a *ArrayType) Elem() Type {
	return a.Type
}
