package types

import "testing"

func TestCoreValueContracts(t *testing.T) {
	if got, ok := GetBuiltInType("string"); !ok || !got.Equal(String) {
		t.Fatal("string builtin lookup failed")
	}
	if !PointerTo(String).Equal(PointerTo(String)) || !SliceOf(String).Equal(SliceOf(String)) || !MapOf(String, Int).Equal(MapOf(String, Int)) {
		t.Fatal("compound equality failed")
	}
	pkg := NewPackage("example.test/types", "types")
	alias, err := AliasOf(pkg, "Text", String)
	if err != nil {
		t.Fatal(err)
	}
	if alias.Zero(false, "") != `Text("")` {
		t.Fatalf("alias zero = %q", alias.Zero(false, ""))
	}
	if got := PointerTo(String).Elem(); !got.Equal(String) {
		t.Fatalf("pointer element = %s", got)
	}
}

func TestStructTagLookup(t *testing.T) {
	tag := StructTag(`json:"name" xml:""`)
	if tag.Get("json") != "name" {
		t.Fatal("json lookup failed")
	}
	if value, ok := tag.Lookup("xml"); !ok || value != "" {
		t.Fatal("empty lookup failed")
	}
	if _, ok := tag.Lookup("missing"); ok {
		t.Fatal("missing lookup succeeded")
	}
}
