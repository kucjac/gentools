//go:build go1.27

package genericmethods

// Box is the receiver used by the generic-receiver-method fixture.
type Box[T any] struct {
	Value T
}

// Convert exercises receiver parameters, variadics, inputs, and results.
// Go does not permit methods to declare their own type parameters.
func (box Box[T]) Convert(values ...T) (T, []T) {
	return box.Value, values
}
