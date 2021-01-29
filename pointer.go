package astreflect

// PointerTo creates a pointer of given type.
func PointerTo(pointsTo Type) *PointerType {
	return &PointerType{PointedType: pointsTo}
}

var _ Type = PointerType{}

// PointerType is the type implementation that defines pointer type.
type PointerType struct {
	PointedType Type
}

// Name implements Type interface.
func (p PointerType) Name(identified bool, pkgContext string) string {
	return "*" + p.PointedType.Name(identified, pkgContext)
}

// FullName implements Type interface.
func (p PointerType) FullName() string {
	return "*" + p.PointedType.FullName()
}

// Kind implements Type interface.
func (p PointerType) Kind() Kind {
	return Ptr
}

// Elem implements Type interface.
func (p PointerType) Elem() Type {
	return p.PointedType
}

// String implements Type interface.
func (p PointerType) String() string {
	return p.Name(true, "")
}

// Zero implements Type interface.
func (p PointerType) Zero(_ bool, _ string) string {
	return "nil"
}

// Equal implements Type interface.
func (p PointerType) Equal(another Type) bool {
	var pt PointerType
	switch t := another.(type) {
	case PointerType:
		pt = t
	case *PointerType:
		pt = *t
	default:
		return false
	}
	return p.PointedType.Equal(pt.PointedType)
}
