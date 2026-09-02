// Package consumer instantiates an imported generic type.
package consumer

import "github.com/kucjac/gentools/parser/testcases/genericimport"

// Holder contains an imported generic instantiation.
type Holder struct {
	Value genericimport.Box[string]
}
