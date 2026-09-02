package parser

import (
	"strings"
	"testing"

	"github.com/kucjac/gentools/types"
)

const genericFixturePackage = "github.com/kucjac/gentools/parser/testcases/generic"

func TestLoadGenericDeclarationsAndInstantiations(t *testing.T) {
	pkgs, err := LoadPackages(LoadConfig{PkgNames: []string{genericFixturePackage}})
	if err != nil {
		t.Fatalf("LoadPackages() error = %v", err)
	}
	pkg := pkgs.MustGetByPath(genericFixturePackage)

	pair, ok := pkg.GetStruct("Pair")
	if !ok {
		t.Fatal("Pair not loaded as a struct")
	}
	info := pair.GenericInfo()
	if len(info.TypeParameters) != 2 {
		t.Fatalf("Pair type parameters = %#v, want two", info.TypeParameters)
	}
	if info.TypeParameters[0].Identifier != "T" || info.TypeParameters[0].Position != 0 {
		t.Fatalf("first Pair parameter = %#v", info.TypeParameters[0])
	}
	if info.TypeParameters[0].Constraint == nil || !strings.Contains(info.TypeParameters[0].Constraint.Expression, "~int") {
		t.Fatalf("Pair constraint = %#v, want preserved type set", info.TypeParameters[0].Constraint)
	}
	if info.TypeParameters[1].Identifier != "U" || info.TypeParameters[1].Position != 1 {
		t.Fatalf("second Pair parameter = %#v", info.TypeParameters[1])
	}

	identity, ok := pkg.GetFunction("Identity")
	if !ok {
		t.Fatal("Identity not loaded")
	}
	if got := identity.GenericInfo().TypeParameters; len(got) != 1 || got[0].Identifier != "T" {
		t.Fatalf("Identity generic metadata = %#v", got)
	}
	collect, ok := pkg.GetFunction("Collect")
	if !ok || !collect.Variadic || len(collect.In) != 1 {
		t.Fatalf("Collect = %#v, want variadic generic function", collect)
	}

	uses, ok := pkg.GetStruct("Uses")
	if !ok || len(uses.Fields) != 2 {
		t.Fatalf("Uses = %#v", uses)
	}
	instantiation, ok := uses.Fields[0].Type.(*types.Instantiation)
	if !ok {
		t.Fatalf("Uses.Direct = %T, want *types.Instantiation", uses.Fields[0].Type)
	}
	if instantiation.Origin != pair || len(instantiation.Arguments) != 2 {
		t.Fatalf("Uses.Direct instantiation = %#v", instantiation)
	}
	if instantiation.Arguments[0].Kind() != types.KindInt || instantiation.Arguments[1].Kind() != types.KindString {
		t.Fatalf("Uses.Direct arguments = %#v", instantiation.Arguments)
	}

	if len(pair.Methods) != 1 || pair.Methods[0].FuncName != "Copy" {
		t.Fatalf("Pair methods = %#v", pair.Methods)
	}
	if got := pair.Methods[0].GenericInfo().ReceiverTypeParameters; len(got) != 2 {
		t.Fatalf("Copy receiver metadata = %#v", got)
	}
}
