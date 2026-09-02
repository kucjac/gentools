// Package genericimport supplies generic declarations used by another fixture package.
package genericimport

// Box is imported and instantiated by consumer.
type Box[T any] struct {
	Value T
}
