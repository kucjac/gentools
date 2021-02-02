package types

import (
	"strings"
)

// ChanOf creates the channel of given type with given direction.
func ChanOf(dir ChanDir, chanType Type) *Chan {
	return &Chan{Type: chanType, Dir: dir}
}

var _ Type = (*Chan)(nil)

// Chan is the type representing channel.
type Chan struct {
	Type Type
	Dir  ChanDir
}

// Name implements Type interface.
func (c *Chan) Name(identified bool, packageContext string) string {
	sb := strings.Builder{}
	if c.Dir == SendOnly {
		sb.WriteString("<-")
	}
	sb.WriteString("chan")
	if c.Dir == RecvOnly {
		sb.WriteString("<-")
	}
	sb.WriteRune(' ')
	sb.WriteString(c.Type.Name(identified, packageContext))
	return sb.String()
}

// FullName implements Type interface.
func (c *Chan) FullName() string {
	sb := strings.Builder{}
	if c.Dir == SendOnly {
		sb.WriteString("<-")
	}
	sb.WriteString("chan")
	if c.Dir == RecvOnly {
		sb.WriteString("<-")
	}
	sb.WriteRune(' ')
	sb.WriteString(c.Type.FullName())
	return sb.String()
}

// Kind gets the kind of the type.
func (c *Chan) Kind() Kind {
	return KindChan
}

// Elem gets the channel element type.
func (c *Chan) Elem() Type {
	return c.Type
}

// KindString implements fmt.Stringer interface.
func (c Chan) String() string {
	return c.Name(true, "")
}

// Zero implements Type interface.
func (c *Chan) Zero(_ bool, _ string) string {
	return "nil"
}

// Equal implements Type interface.
func (c *Chan) Equal(another Type) bool {
	ct, ok := another.(*Chan)
	if !ok {
		return false
	}
	if c.Dir != ct.Dir {
		return false
	}
	return c.Type.Equal(ct.Type)
}

// A ChanDir value indicates a channel direction.
type ChanDir int

// The direction of a channel is indicated by one of these constants.
const (
	SendRecv ChanDir = iota
	SendOnly
	RecvOnly
)
