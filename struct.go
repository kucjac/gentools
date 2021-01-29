package astreflect

import "strconv"

// Compile type check if *StructType implements Type interface.
var _ Type = (*StructType)(nil)

// StructType is the struct type reflection.
type StructType struct {
	PackagePath PkgPath
	Comment     string
	TypeName    string
	Fields      []StructField
	Methods     []FunctionType
}

// Implements checks if given structure implements provided interface.
func (s *StructType) Implements(interfaceType *InterfaceType, pointer bool) bool {
	return implements(interfaceType, s, pointer)
}

// Name implements Type interface.
func (s *StructType) Name(identifier bool, packageContext string) string {
	if identifier && packageContext != s.PackagePath.FullName() {
		if i := s.PackagePath.Identifier(); i != "" {
			return i + "." + s.TypeName
		}
	}
	return s.TypeName
}

// FullName implements Type interface.
func (s *StructType) FullName() string {
	return string(s.PackagePath) + "/" + s.TypeName
}

// PkgPath implements Type interface.
func (s *StructType) PkgPath() PkgPath {
	return s.PackagePath
}

// Kind implements Type interface.
func (s *StructType) Kind() Kind {
	return Struct
}

// Elem implements Type interface.
func (s *StructType) Elem() Type {
	return nil
}

// String implements Type interface.
func (s *StructType) String() string {
	return s.Name(true, "")
}

// Zero implements Type interface.
func (s *StructType) Zero(identified bool, packageContext string) string {
	return s.Name(identified, packageContext) + "{}"
}

// Equal implements Type interface.
func (s *StructType) Equal(another Type) bool {
	st, ok := another.(*StructType)
	if !ok {
		return false
	}
	return st.PackagePath == s.PackagePath && st.TypeName == s.TypeName
}

func (s *StructType) getMethods() []FunctionType {
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

// String implements fmt.Stringer interface.
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
	// When modifying this code, also update the validateStructTag code
	// in cmd/vet/structtag.go.

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
