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

// Type is the interface used by all golang type reflections in package.
type Type interface {
	// Name gets the type name with or without package identifier.
	// An optional packageContext parameter defines the name of the package (full package name) that is expected to be
	// within given context of search. This could be used to get the chain of names with respect to some package.
	// Example:
	//	Developer wants to generate some additional method for the type 'X' within package 'my.com/testing/pkg'.
	//	In order to generate valid names for the imported types the identity needs to be set to 'true'.
	//	But current package context ('my.com/testing/pkg') should not be used be prefixed with the identifier.
	//	Thus an optional 'packageContext' parameter needs to be set to 'my.com/testing/pkg'.
	Name(identified bool, packageContext string) string
	// FullName gets the full name of given type with the full package name and a type.
	FullName() string
	// PkgPath gets the PkgPath for given type.
	PkgPath() PkgPath
	// Kind gets the Kind of given type.
	Kind() Kind
	// Elem gets the wrapped, pointed, base of
	Elem() Type
	// String gets the full name string representation of given type.
	String() string
	// Zero gets zero value string of given type.
	Zero(identified bool, packageContext string) string
	// Equal checks if the types matches exact.
	Equal(another Type) bool
}
