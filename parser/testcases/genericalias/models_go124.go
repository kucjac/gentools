//go:build go1.24

package genericalias

// Set is a generic alias with a type-set constraint.
type Set[K comparable] = map[K]bool

// Strings exercises a generic alias instantiation.
type Strings = Set[string]

// Uses retains a generic alias instantiation in a declaration.
type Uses struct {
	Value Set[string]
}
