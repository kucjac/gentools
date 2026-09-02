package types

import "testing"

func TestGenericInfoIsAdditive(t *testing.T) {
	pkg := NewPackage("example.test/models", "models")
	nonGeneric := &Struct{Pkg: pkg, TypeName: "Plain"}
	if info := nonGeneric.GenericInfo(); len(info.TypeParameters) != 0 {
		t.Fatalf("non-generic struct unexpectedly has metadata: %#v", info)
	}

	generic := &Struct{Pkg: pkg, TypeName: "Box"}
	SetGenericInfo(generic, GenericInfo{TypeParameters: []TypeParameter{{
		Identifier: "T",
		Owner:      generic.FullName(),
		Position:   0,
		Constraint: &Constraint{Type: &Interface{}, Expression: "any"},
	}}})
	info := generic.GenericInfo()
	if len(info.TypeParameters) != 1 || info.TypeParameters[0].Identifier != "T" {
		t.Fatalf("generic metadata = %#v, want one T parameter", info)
	}
	info.TypeParameters[0].Identifier = "changed"
	if generic.GenericInfo().TypeParameters[0].Identifier != "T" {
		t.Fatal("GenericInfo returned mutable stored metadata")
	}
}

func TestInstantiationIdentityAndRendering(t *testing.T) {
	pkg := NewPackage("example.test/models", "models")
	box := &Struct{Pkg: pkg, TypeName: "Box"}
	intBox := &Instantiation{Origin: box, Arguments: []Type{Int}}
	stringBox := &Instantiation{Origin: box, Arguments: []Type{String}}
	if got, want := intBox.String(), "models.Box[int]"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
	if !intBox.Equal(&Instantiation{Origin: box, Arguments: []Type{Int}}) {
		t.Fatal("equivalent instantiations are not equal")
	}
	if intBox.Equal(stringBox) {
		t.Fatal("different type arguments compare equal")
	}
}
