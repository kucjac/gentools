package types

// AliasOf creates an Alias of given aliasType at given package with given name.
func AliasOf(pkg *Package, name string, aliasType Type) (*Alias, error) {
	if err := pkg.hasName(name); err != nil {
		return nil, err
	}
	a := aliasOf(pkg, name, aliasType)
	return a, nil
}

func aliasOf(pkg *Package, name string, aType Type) *Alias {
	a := &Alias{
		Pkg:       pkg,
		AliasName: name,
		Type:      aType,
	}
	pkg.Types[name] = a
	pkg.Aliases = append(pkg.Aliases, a)
	return a
}

var _ Type = (*Alias)(nil)

// Alias is the type that represents wrapped and named another type.
// I.e.: 'type Custom int' would be a Alias over BuiltIn(int) type.
type Alias struct {
	Comment   string
	Pkg       *Package
	AliasName string
	Type      Type
}

// Name implements Type interface.
func (w *Alias) Name(identified bool, packageContext string) string {
	if identified && packageContext != w.Pkg.Path {
		if i := w.Pkg.Identifier; i != "" {
			return i + "." + w.AliasName
		}
	}
	return w.AliasName
}

// FullName implements Type interface.
func (w *Alias) FullName() string {
	return w.Pkg.Path + "/" + w.AliasName
}

// PkgPath implements Type interface.
func (w *Alias) Package() *Package {
	return w.Pkg
}

// Kind implements Type interface.
func (w *Alias) Kind() Kind {
	return KindWrapper
}

// Elem implements Type interface.
func (w *Alias) Elem() Type {
	return w.Type
}

// KindString implements Type interface.
func (w Alias) String() string {
	return w.Name(true, "")
}

// Zero implements Type interface.
func (w *Alias) Zero(identified bool, packageContext string) string {
	t := w.Type
	for t.Kind() == KindWrapper {
		t = t.Elem()
	}

	if t.Kind().IsBuiltin() {
		return w.Name(identified, packageContext) + "(" + t.Zero(identified, packageContext) + ")"
	}

	switch t.Kind() {
	case KindStruct, KindArray:
		return w.Name(identified, packageContext) + "{}"
	case KindSlice, KindInterface, KindChan, KindMap, KindFunc, KindPtr:
		return "nil"
	default:
		return "nil"
	}
}

// Equal implements Type interface.
func (w *Alias) Equal(another Type) bool {
	wt, ok := another.(*Alias)
	if !ok {
		return false
	}
	return w.Pkg == wt.Pkg && wt.AliasName == w.AliasName
}
