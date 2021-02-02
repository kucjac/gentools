package types

import (
	"strconv"
	"strings"
)

// Compile type check if *Struct implements Type interface.
var _ Type = (*Struct)(nil)

// Struct is the struct type reflection.
type Struct struct {
	Pkg      *Package
	Comment  string
	TypeName string
	Fields   []StructField
	Methods  []Function
}

// Implements checks if given structure implements provided interface.
func (s *Struct) Implements(interfaceType *Interface, pointer bool) bool {
	return implements(interfaceType, s, pointer)
}

// Name implements Type interface.
func (s *Struct) Name(identifier bool, packageContext string) string {
	if identifier && packageContext != s.Pkg.Path {
		if i := s.Pkg.Identifier; i != "" {
			return i + "." + s.TypeName
		}
	}
	return s.TypeName
}

// FullName implements Type interface.
func (s *Struct) FullName() string {
	return s.Pkg.Path + "/" + s.TypeName
}

// PkgPath implements Type interface.
func (s *Struct) Package() *Package {
	return s.Pkg
}

// Kind implements Type interface.
func (s *Struct) Kind() Kind {
	return KindStruct
}

// Elem implements Type interface.
func (s *Struct) Elem() Type {
	return nil
}

// KindString implements Type interface.
func (s *Struct) String() string {
	return s.Name(true, "")
}

// Zero implements Type interface.
func (s *Struct) Zero(identified bool, packageContext string) string {
	return s.Name(identified, packageContext) + "{}"
}

// Equal implements Type interface.
func (s *Struct) Equal(another Type) bool {
	st, ok := another.(*Struct)
	if !ok {
		return false
	}
	return st.Pkg == s.Pkg && st.TypeName == s.TypeName
}

func (s *Struct) getMethods() []Function {
	return s.Methods
}

// StructField is a structure field model.
type StructField struct {
	Name      string
	Comment   string
	Type      Type
	Tag       StructTag
	Index     []int
	Embedded  bool
	Anonymous bool
}

// KindString implements fmt.Stringer interface.
func (s StructField) String() string {
	return s.Name + "\t" + s.Type.String() + "\t`" + string(s.Tag) + "`"
}

// A StructTag is the tag string in a struct field.
//
// By convention, tag strings are a concatenation of
// optionally space-separated key:"value" pairs.
// Each key is a non-empty string consisting of non-control
// characters other than space (U+0020 ' '), quote (U+0022 '"'),
// and colon (U+003A ':').  Each value is quoted using U+0022 '"'
// characters and Go string literal syntax.
type StructTag string

// StructTagTuple is a tuple key,value for the struct tag.
type StructTagTuple struct {
	Key, Value string
}

// StructTagTuples is the slice alias over the the structTag key, value tuples.
// It is used to recreate the StructTag.
type StructTagTuples []StructTagTuple

// Join joins the structTag tuples and creates a single StructTag.
func (s StructTagTuples) Join() StructTag {
	var sb strings.Builder
	for i := range s {
		sb.WriteString(s[i].Key)
		sb.WriteRune(':')
		sb.WriteRune('"')
		sb.WriteString(s[i].Value)
		sb.WriteRune('"')
		if i != len(s)-1 {
			sb.WriteRune(' ')
		}
	}
	return StructTag(sb.String())
}

// Split splits up the struct tag into key, value tuples.
func (tag StructTag) Split() (tuples StructTagTuples) {
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		quotedValue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(quotedValue)
		if err != nil {
			break
		}
		tuples = append(tuples, StructTagTuple{Key: name, Value: value})
	}
	return tuples
}

// Get returns the value associated with key in the tag string.
// If there is no such key in the tag, Get returns the empty string.
// If the tag does not have the conventional format, the value
// returned by Get is unspecified. To determine whether a tag is
// explicitly set to the empty string, use Lookup.
func (tag StructTag) Get(key string) string {
	v, _ := tag.Lookup(key)
	return v
}

// Lookup returns the value associated with key in the tag string.
// If the key is present in the tag the value (which may be empty)
// is returned. Otherwise the returned value will be the empty string.
// The ok return value reports whether the value was explicitly set in
// the tag string. If the tag does not have the conventional format,
// the value returned by Lookup is unspecified.
func (tag StructTag) Lookup(key string) (value string, ok bool) {
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}
