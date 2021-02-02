package types

import (
	"fmt"
)

// SliceOf creates the KindSlice of given inner types.
func SliceOf(inner Type) *Array {
	return &Array{ArrayKind: KindSlice, Type: inner}
}

// ArrayOf creates the KindArray of given size with inner type.
func ArrayOf(inner Type, size int) *Array {
	return &Array{ArrayKind: KindArray, Type: inner, ArraySize: size}
}

var _ Type = (*Array)(nil)

// Array is the array or slice representing type.
type Array struct {
	ArrayKind Kind // Could be either KindSlice or KindArray
	Type      Type
	ArraySize int
}

// Name implements Type interface.
func (a *Array) Name(identifier bool, packageContext string) string {
	if a.ArrayKind == KindSlice {
		return "[]" + a.Type.Name(identifier, packageContext)
	}
	return fmt.Sprintf("[%d]%s", a.ArraySize, a.Type.Name(identifier, packageContext))
}

// FullName implements Type interface.
func (a *Array) FullName() string {
	if a.ArrayKind == KindSlice {
		return "[]" + a.Type.FullName()
	}
	return fmt.Sprintf("[%d]%s", a.ArraySize, a.Type.FullName())
}

// Kind implements Type interface.
func (a *Array) Kind() Kind {
	return a.ArrayKind
}

// Elem implements Type interface.
func (a *Array) Elem() Type {
	return a.Type
}

// KindString implements Type interface.
func (a Array) String() string {
	return a.Name(true, "")
}

// Zero implements Type interface.
func (a *Array) Zero(identified bool, packageContext string) string {
	if a.Kind() == KindSlice {
		return "nil"
	}
	return a.Name(identified, packageContext) + "{}"
}

// Equal implements Type interface.
func (a *Array) Equal(another Type) bool {
	at, ok := another.(*Array)
	if !ok {
		return false
	}
	if a.ArrayKind != at.ArrayKind {
		return false
	}
	if a.ArrayKind == KindArray && a.ArraySize != at.ArraySize {
		return false
	}
	return a.Type.Equal(at.Type)
}
