package astreflect

// Dereference is getting Type dereferenced basic value.
// If the value is basic returns nil.
func Dereference(t Type) Type {
	var e Type
	for {
		if e = t.Elem(); e == nil {
			break
		}
		t = e
	}
	return e
}

var _ Type = (*PointerType)(nil)

// PointerType is the type implementation that defines pointer type.
type PointerType struct {
	PointedType Type
}

// Name implements Type interface.
func (p *PointerType) Name(identified bool) string {
	return "*" + p.PointedType.Name(identified)
}

// FullName implements Type interface.
func (p *PointerType) FullName() string {
	return "*" + p.PointedType.FullName()
}

// PkgPath implements Type interface.
func (p *PointerType) PkgPath() PkgPath {
	return p.PointedType.PkgPath()
}

// Kind implements Type interface.
func (p *PointerType) Kind() Kind {
	return Ptr
}

// Elem implements Type interface.
func (p *PointerType) Elem() Type {
	return p.PointedType
}

// Type is the interface used by all type reflections in xreflect package.
type Type interface {
	// Name gets the type name with or without package identifier.
	Name(identified bool) string
	// FullName gets the full name of given type with the full package name and a type.
	FullName() string
	// PkgPath gets the PkgPath for given type.
	PkgPath() PkgPath
	// Kind gets the Kind of given type.
	Kind() Kind
	// Elem gets the wrapped, pointed, base of
	Elem() Type
}
