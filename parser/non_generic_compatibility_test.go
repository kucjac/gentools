package parser

import (
	"testing"

	"github.com/kucjac/gentools/types"
)

func TestNonGenericFixtureRemainsConcreteAndUnchanged(t *testing.T) {
	const fixture = "github.com/kucjac/gentools/parser/testcases"
	pkgs, err := LoadPackages(LoadConfig{PkgNames: []string{fixture}})
	if err != nil {
		t.Fatalf("LoadPackages() error = %v", err)
	}
	pkg := pkgs.MustGetByPath(fixture)
	foo, ok := pkg.GetStruct("Foo")
	if !ok {
		t.Fatal("Foo no longer loads as *types.Struct")
	}
	if foo.FullName() != fixture+"/Foo" || foo.Kind() != types.KindStruct || len(foo.Fields) != 9 {
		t.Fatalf("Foo compatibility changed: %#v", foo)
	}
	if len(foo.GenericInfo().TypeParameters) != 0 {
		t.Fatalf("non-generic Foo has generic metadata: %#v", foo.GenericInfo())
	}
	alias, ok := pkg.GetAlias("FooAlias")
	if !ok || !alias.Equal(&types.Alias{Pkg: pkg, AliasName: "FooAlias"}) {
		t.Fatalf("FooAlias compatibility changed: %#v", alias)
	}
}
