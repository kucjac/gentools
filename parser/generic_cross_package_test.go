package parser

import (
	"testing"

	"github.com/kucjac/gentools/types"
)

func TestLoadImportedGenericInstantiation(t *testing.T) {
	const consumerPackage = "github.com/kucjac/gentools/parser/testcases/genericimport/consumer"
	pkgs, err := LoadPackages(LoadConfig{PkgNames: []string{consumerPackage}})
	if err != nil {
		t.Fatalf("LoadPackages() error = %v", err)
	}
	consumer := pkgs.MustGetByPath(consumerPackage)
	holder, ok := consumer.GetStruct("Holder")
	if !ok || len(holder.Fields) != 1 {
		t.Fatalf("Holder = %#v", holder)
	}
	value, ok := holder.Fields[0].Type.(*types.Instantiation)
	if !ok {
		t.Fatalf("Holder.Value = %T, want *types.Instantiation", holder.Fields[0].Type)
	}
	if value.Origin.FullName() != "github.com/kucjac/gentools/parser/testcases/genericimport/Box" {
		t.Fatalf("origin = %q", value.Origin.FullName())
	}
	if len(value.Arguments) != 1 || value.Arguments[0].Kind() != types.KindString {
		t.Fatalf("arguments = %#v", value.Arguments)
	}
}

func TestUpdatePackagesReusesLoadedGenericDependency(t *testing.T) {
	const (
		dependencyPackage = "github.com/kucjac/gentools/parser/testcases/genericimport"
		consumerPackage   = dependencyPackage + "/consumer"
	)
	packages, err := LoadPackages(LoadConfig{PkgNames: []string{dependencyPackage}})
	if err != nil {
		t.Fatalf("LoadPackages() error = %v", err)
	}
	dependency := packages.MustGetByPath(dependencyPackage)
	box, ok := dependency.GetStruct("Box")
	if !ok {
		t.Fatal("loaded dependency Box not found")
	}

	if err := UpdatePackages(packages, LoadConfig{PkgNames: []string{consumerPackage}}); err != nil {
		t.Fatalf("UpdatePackages() error = %v", err)
	}
	if packages.MustGetByPath(dependencyPackage) != dependency {
		t.Fatal("UpdatePackages() replaced the loaded dependency package")
	}
	consumer := packages.MustGetByPath(consumerPackage)
	holder, ok := consumer.GetStruct("Holder")
	if !ok || len(holder.Fields) != 1 {
		t.Fatalf("Holder = %#v", holder)
	}
	value, ok := holder.Fields[0].Type.(*types.Instantiation)
	if !ok {
		t.Fatalf("Holder.Value = %T, want *types.Instantiation", holder.Fields[0].Type)
	}
	if value.Origin != box {
		t.Fatalf("Holder.Value origin = %#v, want loaded Box %#v", value.Origin, box)
	}
}
