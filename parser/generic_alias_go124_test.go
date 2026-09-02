//go:build go1.24

package parser

import (
	"os"
	"testing"

	"github.com/kucjac/gentools/types"
)

func TestLoadGenericAliases(t *testing.T) {
	if debug := os.Getenv("GODEBUG"); debug == "" {
		t.Setenv("GODEBUG", "gotypesalias=1")
	} else {
		t.Setenv("GODEBUG", debug+",gotypesalias=1")
	}
	const aliasPackage = "github.com/kucjac/gentools/parser/testcases/genericalias"
	pkgs, err := LoadPackages(LoadConfig{PkgNames: []string{aliasPackage}})
	if err != nil {
		t.Fatalf("LoadPackages() error = %v", err)
	}
	pkg := pkgs.MustGetByPath(aliasPackage)
	set, ok := pkg.GetAlias("Set")
	if !ok {
		t.Fatal("Set not loaded as an alias")
	}
	if params := set.GenericInfo().TypeParameters; len(params) != 1 || params[0].Identifier != "K" {
		t.Fatalf("Set generic metadata = %#v", params)
	}
	uses, ok := pkg.GetStruct("Uses")
	if !ok || len(uses.Fields) != 1 {
		t.Fatalf("Uses = %#v", uses)
	}
	value, ok := uses.Fields[0].Type.(*types.Instantiation)
	if !ok {
		t.Fatalf("Uses.Value = %T, want *types.Instantiation", uses.Fields[0].Type)
	}
	if value.Origin != set || len(value.Arguments) != 1 || value.Arguments[0].Kind() != types.KindString {
		t.Fatalf("generic alias instantiation = %#v", value)
	}
}
