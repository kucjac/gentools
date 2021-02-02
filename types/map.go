package types

import (
	"fmt"
)

// MapOf creates a map of given key, value types.
func MapOf(key, value Type) *Map {
	return &Map{Key: key, Value: value}
}

var _ Type = (*Map)(nil)

// Map is the type wrapper for the standar key value map type.
type Map struct {
	Key   Type
	Value Type
}

// Name implements Type interface.
func (m *Map) Name(identified bool, packageContext string) string {
	return fmt.Sprintf("map[%s]%s", m.Key.Name(identified, packageContext), m.Value.Name(identified, packageContext))
}

// FullName implements Type interface.
func (m *Map) FullName() string {
	return fmt.Sprintf("map[%s]%s", m.Key.FullName(), m.Value.FullName())
}

// Kind implements Type interface.
func (m *Map) Kind() Kind {
	return KindMap
}

// Elem as the map has both the key and value it needs to be dereferenced manually.
func (m *Map) Elem() Type {
	return nil
}

// KindString implements Type interface.
func (m Map) String() string {
	return m.Name(true, "")
}

// Zero implements Type interface.
func (m *Map) Zero(_ bool, _ string) string {
	return "nil"
}

// Equal implements Type interface.
func (m *Map) Equal(another Type) bool {
	mp, ok := another.(*Map)
	if !ok {
		return false
	}
	if !mp.Key.Equal(m.Key) {
		return false
	}
	return mp.Value.Equal(m.Value)
}
