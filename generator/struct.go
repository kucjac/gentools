package genutils

import (
	"fmt"
	"strings"

	"github.com/kucjac/gentools/types"
)

// Struct is the structure definition generator.
type Struct struct {
	tp      *types.Struct
	methods []*FuncDef
	defFunc func(d *StructDef)
}

func (s *Struct) Method(name string, methodCreator func(f MethodCreator)) error {
	return nil
}

// Type get the structure type.
func (s *Struct) Type() *types.Struct {
	return s.tp
}

// Methods returns structure method definitions.
func (s *Struct) Methods() []*FuncDef {
	return s.methods
}

// MethodCreator is an interface used for generating type method.
type MethodCreator interface {
	FuncCreator
	Receiver(name string, pointer bool)
}

// StructDef is the structure generator.
type StructDef struct {
	s          *Struct
	err        error
	methods    []*FuncDef
	fieldNames map[string]struct{}
}

// Field adds the field definition to the given structure.
func (s *StructDef) Field(name string, tp types.Type, tags ...FieldTagger) types.StructField {
	if s.err != nil {
		return types.StructField{}
	}
	if _, ok := s.fieldNames[name]; ok {
		s.err = fmt.Errorf("struct: '%s' field: '%s' already defined", s.s.tp.TypeName, name)
	}
	index := []int{len(s.s.tp.Fields)}
	sField := types.StructField{
		Name:  name,
		Type:  tp,
		Index: index,
	}
	for _, tagger := range tags {
		tagger(&sField)
	}
	s.s.tp.Fields = append(s.s.tp.Fields, sField)
	return sField
}

func (s *StructDef) Method(name string, cf func(fc MethodCreator)) {
}

func (s *StructDef) MethodComment(name string, comment string) {
}

// Content implements Block interface.
func (s *Struct) Content() ([]byte, error) {
	return nil, nil
}

type fieldValue struct {
	Name, Value string
}

type structValuer struct {
	fields []fieldValue
	st     *types.Struct
	err    error
	ptr    bool
}

func (s *structValuer) Field(name string, value interface{}) {
	if s.err != nil {
		return
	}

}

func (s *structValuer) Ptr() {
	s.ptr = true
}

// FieldTagger is the function type that sets the tag for given struct field.
type FieldTagger func(*types.StructField)

// NameFunc is the function used for converting input name into some common naming convention i.e.: CamelCase.
type NameFunc func(name string) string

// JSONTagger is the function that sets up the JSON field tag with provided name func for given struct field.
func JSONTagger(nameFunc NameFunc, omitempty bool) FieldTagger {
	return func(sField *types.StructField) {
		var sb strings.Builder
		sb.WriteString(nameFunc(sField.Name))
		if omitempty {
			sb.WriteString(",omitempty")
		}

		var found bool
		tuples := sField.Tag.Split()
		for i := range tuples {
			if tuples[i].Key == "json" {
				// Replace the tag value.
				tuples[i].Value = sb.String()
				found = true
				break
			}
		}
		if !found {
			tuples = append(tuples, types.StructTagTuple{Key: "json", Value: sb.String()})
		}
		sField.Tag = tuples.Join()
	}
}
