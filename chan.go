package astreflect

import (
	"strings"
)

var _ Type = (*ChanType)(nil)

// ChanType is the type representing channel.
type ChanType struct {
	Type Type
	Dir  ChanDir
}

// Name implements Type interface.
func (c ChanType) Name(identified bool, packageContext string) string {
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
func (c ChanType) FullName() string {
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

// PkgPath implements Type interface.
func (c ChanType) PkgPath() PkgPath {
	return builtInPkgPath
}

// Kind gets the kind of the type.
func (c ChanType) Kind() Kind {
	return Chan
}

// Elem gets the channel element type.
func (c ChanType) Elem() Type {
	return c.Type
}

// String implements fmt.Stringer interface.
func (c ChanType) String() string {
	return c.Name(true, "")
}

// Zero implements Type interface.
func (c ChanType) Zero(_ bool, _ string) string {
	return "nil"
}

// Equal implements Type interface.
func (c ChanType) Equal(another Type) bool {
	var ct ChanType
	switch t := another.(type) {
	case *ChanType:
		ct = *t
	case ChanType:
		ct = t
	default:
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
