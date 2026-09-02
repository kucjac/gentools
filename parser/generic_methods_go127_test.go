//go:build go1.27

package parser

import "testing"

func TestLoadGenericMethods(t *testing.T) {
	const methodsPackage = "github.com/kucjac/gentools/parser/testcases/genericmethods"
	pkgs, err := LoadPackages(LoadConfig{PkgNames: []string{methodsPackage}})
	if err != nil {
		t.Fatalf("LoadPackages() error = %v", err)
	}
	pkg := pkgs.MustGetByPath(methodsPackage)
	box, ok := pkg.GetStruct("Box")
	if !ok || len(box.Methods) != 1 {
		t.Fatalf("Box = %#v", box)
	}
	method := box.Methods[0]
	if method.FuncName != "Convert" || !method.Variadic || len(method.In) != 1 || len(method.Out) != 2 {
		t.Fatalf("Convert = %#v", method)
	}
	info := method.GenericInfo()
	if len(info.ReceiverTypeParameters) != 1 || info.ReceiverTypeParameters[0].Identifier != "T" {
		t.Fatalf("receiver metadata = %#v", info.ReceiverTypeParameters)
	}
	if len(info.TypeParameters) != 0 {
		t.Fatalf("method-local metadata = %#v, want none", info.TypeParameters)
	}
}
