package types

// PointerTo creates a pointer of given type.
func PointerTo(pointsTo Type) *Pointer {
	return &Pointer{PointedType: pointsTo}
}

var _ Type = (*Pointer)(nil)

// Pointer is the type implementation that defines pointer type.
type Pointer struct {
	PointedType Type
}

// Name implements Type interface.
func (p *Pointer) Name(identified bool, pkgContext string) string {
	return "*" + p.PointedType.Name(identified, pkgContext)
}

// FullName implements Type interface.
func (p *Pointer) FullName() string {
	return "*" + p.PointedType.FullName()
}

// Kind implements Type interface.
func (p *Pointer) Kind() Kind {
	return KindPtr
}

// Elem implements Type interface.
func (p *Pointer) Elem() Type {
	return p.PointedType
}

// KindString implements Type interface.
func (p Pointer) String() string {
	return p.Name(true, "")
}

// Zero implements Type interface.
func (p *Pointer) Zero(_ bool, _ string) string {
	return "nil"
}

// Equal implements Type interface.
func (p *Pointer) Equal(another Type) bool {
	pt, ok := another.(*Pointer)
	if !ok {
		return false
	}
	return p.PointedType.Equal(pt.PointedType)
}
