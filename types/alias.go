package types

// AliasOf creates an Alias of given aliasType at given package with given name.
func AliasOf(pkg *Package, name string, aliasType Type) (*Alias, error) {
	if err := pkg.hasName(name); err != nil {
		return nil, err
	}
	a := aliasOf(pkg, name, aliasType)
	return a, nil
}

var _ Type = (*Alias)(nil)

// Alias is the type that represents wrapped and named another type.
// I.e.: 'type Custom int' would be an Alias over BuiltIn(int) type.
type Alias struct {
	Comment   string
	Pkg       *Package
	AliasName string
	Type      Type
	Methods   []Function
}

// Name implements Type interface.
func (a *Alias) Name(identified bool, packageContext string) string {
	if identified && packageContext != a.Pkg.Path {
		if i := a.Pkg.Identifier; i != "" {
			return i + "." + a.AliasName
		}
	}
	return a.AliasName
}

// FullName implements Type interface.
func (a *Alias) FullName() string {
	return a.Pkg.Path + "/" + a.AliasName
}

// Package implements Type interface.
func (a *Alias) Package() *Package {
	return a.Pkg
}

// Kind implements Type interface.
func (a *Alias) Kind() Kind {
	return a.Type.Kind()
}

// Elem implements Type interface.
func (a *Alias) Elem() Type {
	return a.Type
}

// KindString implements Type interface.
func (a Alias) String() string {
	return a.Name(true, "")
}

// Zero implements Type interface.
func (a *Alias) Zero(identified bool, packageContext string) string {
	t := a.Type

	return a.Name(identified, packageContext) + "(" + t.Zero(identified, packageContext) + ")"
}

// Equal implements Type interface.
func (a *Alias) Equal(another Type) bool {
	wt, ok := another.(*Alias)
	if !ok {
		return false
	}
	return a.Pkg == wt.Pkg && wt.AliasName == a.AliasName
}

// Implements checks if the alias types implements provided interface.
// The argument isPointer states if given the pointer to alias or an alias by itself implements given interface.
func (a *Alias) Implements(interfaceType *Interface, isPointer bool) bool {
	return implements(interfaceType, a, isPointer)
}

func (a *Alias) getMethods() []Function {
	return a.Methods
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
