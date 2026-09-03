package types

import "testing"

func TestPackageDeclarationsCollisionsLookupsAndPath(t *testing.T) {
	pkg := NewPackage("example.test/widget", "widget")
	widget := &Struct{Pkg: pkg, TypeName: "Widget"}
	if err := pkg.NewNamedType("Widget", widget); err != nil {
		t.Fatal(err)
	}
	if _, ok := pkg.GetStruct("Widget"); !ok || pkg.MustStruct("Widget") != widget {
		t.Fatal("struct lookup failed")
	}
	if err := pkg.NewNamedType("Widget", widget); err == nil {
		t.Fatal("duplicate type was accepted")
	}
	if got := pkg.GetPkgPath(); got.Identifier() != "widget" || got.FullName() != "example.test/widget" {
		t.Fatalf("unexpected package path: %q", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("MustFunction did not panic for a missing function")
		}
	}()
	pkg.MustFunction("missing")
}
