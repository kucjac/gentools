package astreflect

var _ Type = WrappedType{}

// WrappedType is the type that represents wrapped and named another type.
// I.e.: 'type Custom int' would be a WrappedType over BuiltIn(int) type.
type WrappedType struct {
	Comment     string
	PackagePath PkgPath
	WrapperName string
	Type        Type
}

// Name implements Type interface.
func (w WrappedType) Name(identified bool, packageContext string) string {
	if identified && packageContext != w.PackagePath.FullName() {
		if i := w.PackagePath.Identifier(); i != "" {
			return i + "." + w.WrapperName
		}
	}
	return w.WrapperName
}

// FullName implements Type interface.
func (w WrappedType) FullName() string {
	return string(w.PackagePath) + "/" + w.WrapperName
}

// PkgPath implements Type interface.
func (w WrappedType) PkgPath() PkgPath {
	return w.PackagePath
}

// Kind implements Type interface.
func (w WrappedType) Kind() Kind {
	return Wrapper
}

// Elem implements Type interface.
func (w WrappedType) Elem() Type {
	return w.Type
}

// String implements Type interface.
func (w WrappedType) String() string {
	return w.Name(true, "")
}

// Zero implements Type interface.
func (w WrappedType) Zero(identified bool, packageContext string) string {
	t := w.Type
	for t.Kind() == Wrapper {
		t = t.Elem()
	}

	if t.Kind().IsBuiltin() {
		return w.Name(identified, packageContext) + "(" + t.Zero(identified, packageContext) + ")"
	}

	switch t.Kind() {
	case Struct, Array:
		return w.Name(identified, packageContext) + "{}"
	case Slice, Interface, Chan, Map, Func, Ptr:
		return "nil"
	default:
		return "nil"
	}
}

// Equal implements Type interface.
func (w WrappedType) Equal(another Type) bool {
	var wt WrappedType
	switch t := another.(type) {
	case WrappedType:
		wt = t
	case *WrappedType:
		wt = *t
	default:
		return false
	}
	return w.PackagePath == wt.PackagePath && wt.WrapperName == w.WrapperName
}
