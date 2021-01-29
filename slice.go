package astreflect

import "fmt"

// SliceOf creates the Slice of given inner types.
func SliceOf(inner Type) *ArrayType {
	return &ArrayType{ArrayKind: Slice, Type: inner}
}

// ArrayOf creates the Array of given size with inner type.
func ArrayOf(inner Type, size int) *ArrayType {
	return &ArrayType{ArrayKind: Array, Type: inner, ArraySize: size}
}

var _ Type = ArrayType{}

// ArrayType is the array or slice representing type.
type ArrayType struct {
	ArrayKind Kind // Could be either Slice or Array
	Type      Type
	ArraySize int
}

// Name implements Type interface.
func (a ArrayType) Name(identifier bool, packageContext string) string {
	if a.ArrayKind == Slice {
		return "[]" + a.Type.Name(identifier, packageContext)
	}
	return fmt.Sprintf("[%d]%s", a.ArraySize, a.Type.Name(identifier, packageContext))
}

// FullName implements Type interface.
func (a ArrayType) FullName() string {
	if a.ArrayKind == Slice {
		return "[]" + a.Type.FullName()
	}
	return fmt.Sprintf("[%d]%s", a.ArraySize, a.Type.FullName())
}

// Kind implements Type interface.
func (a ArrayType) Kind() Kind {
	return a.ArrayKind
}

// Elem implements Type interface.
func (a ArrayType) Elem() Type {
	return a.Type
}

// String implements Type interface.
func (a ArrayType) String() string {
	return a.Name(true, "")
}

// Zero implements Type interface.
func (a ArrayType) Zero(identified bool, packageContext string) string {
	if a.Kind() == Slice {
		return "nil"
	}
	return a.Name(identified, packageContext) + "{}"
}

// Equal implements Type interface.
func (a ArrayType) Equal(another Type) bool {
	var at ArrayType
	switch t := another.(type) {
	case *ArrayType:
		at = *t
	case ArrayType:
		at = t
	default:
		return false
	}
	if a.ArrayKind != at.ArrayKind {
		return false
	}
	if a.ArrayKind == Array && a.ArraySize != at.ArraySize {
		return false
	}
	return a.Type.Equal(at.Type)
}
