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

// Implements checks if the type t implements interface 'interfaceType'.
func Implements(t Type, interfaceType *InterfaceType) bool {
	var (
		s         *StructType
		isPointer bool
	)
	for s == nil {
		switch tt := t.(type) {
		case *PointerType:
			isPointer = true
			t = tt.PointedType
		case *WrappedType:
			t = tt.Type
		case *StructType:
			s = tt
		default:
			return false
		}
	}
	return s.Implements(interfaceType, isPointer)
}

// Implements checks if given structure implements provided interface.
func (s *StructType) Implements(interfaceType *InterfaceType, pointer bool) bool {
	if len(interfaceType.Methods) > len(s.Methods) {
		return false
	}
	for _, iMethod := range interfaceType.Methods {
		var found bool
		for _, sMethod := range s.Methods {
			if sMethod.FuncName == iMethod.FuncName {
				if len(iMethod.In) != len(sMethod.In) {
					return false
				}
				if len(iMethod.Out) != len(sMethod.Out) {
					return false
				}
				if sMethod.Variadic != iMethod.Variadic {
					return false
				}
				if sMethod.Receiver.Pointer && !pointer {
					return false
				}

				for _, iIn := range iMethod.In {
					for _, sIn := range sMethod.In {
						if iIn.Type.FullName() != sIn.Type.FullName() {
							return false
						}
					}
				}

				for _, iOut := range iMethod.Out {
					for _, sOut := range sMethod.Out {
						if iOut.Type.FullName() != sOut.Type.FullName() {
							return false
						}
					}
				}
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (s *StructType) Name(identifier bool) string {
	if identifier {
		if i := s.PackagePath.Identifier(); i != "" {
			return i + "." + s.TypeName
		}
	}
	return s.TypeName
}

func (s *StructType) FullName() string {
	return string(s.PackagePath) + "/" + s.TypeName
}

func (s *StructType) PkgPath() PkgPath {
	return s.PackagePath
}

func (s *StructType) Kind() Kind {
	return Struct
}

func (s *StructType) Elem() Type {
	return nil
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
