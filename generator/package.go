package genutils

import (
	"fmt"

	"github.com/kucjac/gentools/types"
)

// Package is the package generator definition.
type Package struct {
	tp *types.Package
	pm types.PackageMap
}

// NewPackage creates a new package.
func NewPackage(pm types.PackageMap, name, identifier string) (*Package, error) {
	pkg, err := pm.NewPackage(name, identifier)
	if err != nil {
		return nil, err
	}
	return &Package{tp: pkg, pm: pm}, nil
}

// Func creates new function definition generator.
func (p *Package) Func(name string, cf func(c FuncCreator)) (*FuncDef, error) {
	if _, ok := p.tp.Types[name]; ok {
		return nil, fmt.Errorf("package: '%s' has already a type with a name: '%s'", p.tp.Path, name)
	}
	if _, ok := p.tp.Declarations[name]; ok {
		return nil, fmt.Errorf("package: '%s' has already a type with a name: '%s'", p.tp.Path, name)
	}
	return &FuncDef{tp: &types.Function{Pkg: p.tp, FuncName: name}, contentFunc: cf}, nil
}

// Struct creates a struct definition.
func (p *Package) Struct(name string, defFunc func(d *StructDef)) (*Struct, error) {
	st := &types.Struct{Pkg: p.tp, TypeName: name}
	if err := p.tp.NewNamedType(name, st); err != nil {
		return nil, err
	}
	return &Struct{tp: st, defFunc: defFunc}, nil
}
