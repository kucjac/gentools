package astreflect

var _ Type = (*WrappedType)(nil)

// WrappedType is the type that represents wrapped and named another type.
// I.e.: 'type Custom int' would be a WrappedType over BuiltIn(int) type.
type WrappedType struct {
	Comment     string
	PackagePath PkgPath
	WrapperName string
	Type        Type
}

// Name implements Type interface.
func (w *WrappedType) Name(identified bool) string {
	if identified {
		if i := w.PackagePath.Identifier(); i != "" {
			return i + "." + w.WrapperName
		}
	}
	return w.WrapperName
}

// FullName implements Type interface.
func (w *WrappedType) FullName() string {
	return string(w.PackagePath) + "/" + w.WrapperName
}

// PkgPath implements Type interface.
func (w *WrappedType) PkgPath() PkgPath {
	return w.PackagePath
}

// Kind implements Type interface.
func (w *WrappedType) Kind() Kind {
	return Wrapper
}

// Elem implements Type interface.
func (w *WrappedType) Elem() Type {
	return w.Type
}
