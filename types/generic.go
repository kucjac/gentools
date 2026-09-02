package types

import (
	"fmt"
	"strings"
	"sync"
)

// TypeParameter describes a declared generic type parameter. Owner is the full
// name of the declaration that owns the parameter and Position is its zero-based
// position in that declaration's parameter list.
type TypeParameter struct {
	Identifier string
	Constraint *Constraint
	Owner      string
	Position   int
}

// GenericInfo describes type parameters declared by a named type or function.
// ReceiverTypeParameters are present for methods whose receiver is an
// instantiation of a generic named type.
type GenericInfo struct {
	TypeParameters         []TypeParameter
	ReceiverTypeParameters []TypeParameter
}

// Constraint is the resolved representation of a type-parameter constraint.
// Type retains the existing Gentools model where possible. Expression preserves
// the go/types representation of interface type sets (such as ~int | ~string),
// which the historical Interface model cannot otherwise represent structurally.
type Constraint struct {
	Type       Type
	Expression string
}

// Instantiation is a named generic type supplied with concrete type arguments.
type Instantiation struct {
	Origin    Type
	Arguments []Type
}

var _ Type = (*TypeParameter)(nil)
var _ Type = (*Constraint)(nil)
var _ Type = (*Instantiation)(nil)

func (p *TypeParameter) Name(bool, string) string { return p.Identifier }
func (p *TypeParameter) FullName() string         { return p.Owner + "/" + p.Identifier }
func (p *TypeParameter) Kind() Kind               { return KindTypeParameter }
func (p *TypeParameter) Elem() Type               { return nil }
func (p *TypeParameter) String() string           { return p.Identifier }
func (p *TypeParameter) Zero(bool, string) string { return "*new(" + p.Identifier + ")" }
func (p *TypeParameter) Equal(other Type) bool {
	q, ok := other.(*TypeParameter)
	return ok && p.Owner == q.Owner && p.Position == q.Position
}

// Name implements Type for Constraint.
func (c *Constraint) Name(identified bool, packageContext string) string {
	if c == nil || c.Type == nil {
		return c.Expression
	}
	return c.Type.Name(identified, packageContext)
}

// FullName implements Type for Constraint.
func (c *Constraint) FullName() string {
	if c == nil || c.Type == nil {
		return c.Expression
	}
	return c.Type.FullName()
}

// Kind implements Type for Constraint.
func (c *Constraint) Kind() Kind {
	if c == nil || c.Type == nil {
		return Invalid
	}
	return c.Type.Kind()
}

// Elem implements Type for Constraint.
func (c *Constraint) Elem() Type {
	if c == nil || c.Type == nil {
		return nil
	}
	return c.Type.Elem()
}

// String implements fmt.Stringer for Constraint.
func (c *Constraint) String() string { return c.Name(true, "") }

// Zero implements Type for Constraint.
func (c *Constraint) Zero(identified bool, packageContext string) string {
	if c == nil || c.Type == nil {
		return ""
	}
	return c.Type.Zero(identified, packageContext)
}

// Equal implements Type for Constraint.
func (c *Constraint) Equal(other Type) bool {
	otherConstraint, ok := other.(*Constraint)
	if !ok || c == nil || otherConstraint == nil {
		return false
	}
	if c.Type == nil || otherConstraint.Type == nil {
		return c.Type == nil && otherConstraint.Type == nil && c.Expression == otherConstraint.Expression
	}
	return c.Type.Equal(otherConstraint.Type) && c.Expression == otherConstraint.Expression
}

func (i *Instantiation) Name(identified bool, context string) string {
	parts := make([]string, len(i.Arguments))
	for n, arg := range i.Arguments {
		parts[n] = arg.Name(identified, context)
	}
	return i.Origin.Name(identified, context) + "[" + strings.Join(parts, ", ") + "]"
}
func (i *Instantiation) FullName() string { return i.Name(true, "") }
func (i *Instantiation) Kind() Kind       { return i.Origin.Kind() }
func (i *Instantiation) Elem() Type       { return i.Origin.Elem() }
func (i *Instantiation) String() string   { return i.Name(true, "") }
func (i *Instantiation) Zero(identified bool, context string) string {
	return fmt.Sprintf("%s{}", i.Name(identified, context))
}
func (i *Instantiation) Equal(other Type) bool {
	j, ok := other.(*Instantiation)
	if !ok || !i.Origin.Equal(j.Origin) || len(i.Arguments) != len(j.Arguments) {
		return false
	}
	for n := range i.Arguments {
		if !i.Arguments[n].Equal(j.Arguments[n]) {
			return false
		}
	}
	return true
}

var genericInfo sync.Map // map[string]GenericInfo

// SetGenericInfo associates additive generic metadata with a named model value.
func SetGenericInfo(owner Type, info GenericInfo) {
	if key, ok := genericInfoKey(owner); ok {
		genericInfo.Store(key, cloneGenericInfo(info))
	}
}

// Generics returns generic metadata for owner, or an empty value for non-generic types.
func Generics(owner Type) GenericInfo {
	if key, ok := genericInfoKey(owner); ok {
		if value, found := genericInfo.Load(key); found {
			return cloneGenericInfo(value.(GenericInfo))
		}
	}
	return GenericInfo{}
}

func genericInfoKey(owner Type) (string, bool) {
	switch value := owner.(type) {
	case *Struct:
		return "struct:" + value.FullName(), true
	case *Interface:
		return "interface:" + value.FullName(), true
	case *Alias:
		return "alias:" + value.FullName(), true
	case *Function:
		if value.Pkg == nil {
			return "", false
		}
		receiver := ""
		if value.Receiver != nil && value.Receiver.Type != nil {
			receiver = value.Receiver.Type.FullName() + "."
		}
		return "function:" + value.Pkg.Path + "/" + receiver + value.FuncName, true
	default:
		return "", false
	}
}

func cloneGenericInfo(info GenericInfo) GenericInfo {
	info.TypeParameters = append([]TypeParameter(nil), info.TypeParameters...)
	info.ReceiverTypeParameters = append([]TypeParameter(nil), info.ReceiverTypeParameters...)
	return info
}
