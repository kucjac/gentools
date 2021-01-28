package astreflect

import (
	"sync"
)

// Package is the golang package reflection container. It contains all interfaces, structs, functions
// and type wrappers that are located inside of it.
type Package struct {
	Path            PkgPath
	Identifier      string
	Interfaces      []*InterfaceType
	Structs         []*StructType
	Functions       []*FunctionType
	WrappedTypes    []*WrappedType
	Types           map[string]Type
	typesInProgress map[string]Type
	sync.Mutex
}

// PackageMap is a slice wrapper over Package type.
type PackageMap map[string]*Package

// GetByIdentifier gets the package by provided identifier. If there is more than one package with given identifier
// The function would return the first matching package.
func (p PackageMap) GetByIdentifier(identifier string) (*Package, bool) {
	for _, pkg := range p {
		if pkg.Identifier == identifier {
			return pkg, true
		}
	}
	return nil, false
}

// GetByPath gets the package by provided path.
func (p PackageMap) GetByPath(path string) (*Package, bool) {
	for _, pkg := range p {
		if string(pkg.Path) == path {
			return pkg, true
		}
	}
	return nil, false
}

// GetType gets concurrently package type.
func (p *Package) GetType(name string) (Type, bool) {
	p.Lock()
	defer p.Unlock()
	t, ok := p.Types[name]
	return t, ok
}

// GetInterfaceType gets the interface by it's name.
func (p *Package) GetInterfaceType(name string) (*InterfaceType, bool) {
	p.Lock()
	defer p.Unlock()
	t, ok := p.Types[name]
	if !ok {
		return nil, false
	}
	i, ok := t.(*InterfaceType)
	return i, ok
}

// GetStructType gets the struct type by it's name.
func (p *Package) GetStructType(name string) (*StructType, bool) {
	p.Lock()
	defer p.Unlock()
	t, ok := p.Types[name]
	if !ok {
		return nil, false
	}
	s, ok := t.(*StructType)
	return s, ok
}

// GetFunction gets the function type by it's name.
func (p *Package) GetFunction(name string) (*FunctionType, bool) {
	p.Lock()
	defer p.Unlock()
	t, ok := p.Types[name]
	if !ok {
		return nil, false
	}
	s, ok := t.(*FunctionType)
	return s, ok
}

// GetWrappedType gets the wrapped type by it's name.
func (p *Package) GetWrappedType(name string) (*WrappedType, bool) {
	p.Lock()
	defer p.Unlock()
	t, ok := p.Types[name]
	if !ok {
		return nil, false
	}
	s, ok := t.(*WrappedType)
	return s, ok
}

// IsStandard checks if given package is a standard package.
func (p *Package) IsStandard() bool {
	return p == builtIn.Package
}

// setInProgressType sets the in-progress type with given name.
func (p *Package) setInProgressType(name string, tp Type) {
	p.Lock()
	defer p.Unlock()
	p.Types[name] = tp
	p.typesInProgress[name] = tp
	switch t := tp.(type) {
	case *FunctionType:
		p.Functions = append(p.Functions, t)
	case *StructType:
		p.Structs = append(p.Structs, t)
	case *WrappedType:
		p.WrappedTypes = append(p.WrappedTypes, t)
	case *InterfaceType:
		p.Interfaces = append(p.Interfaces, t)
	}
}

// markTypeDone marks type that was in progress as done.
func (p *Package) markTypeDone(name string) {
	p.Lock()
	defer p.Unlock()
	delete(p.typesInProgress, name)
}

// newPackage creates new package for given pkgPath and identifier.
func newPackage(pkgPath, identifier string) *Package {
	pkg := &Package{Path: PkgPath(pkgPath), Identifier: identifier, Types: map[string]Type{}, typesInProgress: map[string]Type{}}
	pkgMap.write(pkgPath, pkg)
	return pkg
}

var pkgMap = &packageMap{pkgMap: map[string]*Package{}}

type packageMap struct {
	sync.Mutex
	pkgMap map[string]*Package
}

func (r *packageMap) read(key string) (*Package, bool) {
	r.Lock()
	defer r.Unlock()
	v, ok := r.pkgMap[key]
	return v, ok
}

func (r *packageMap) write(key string, value *Package) {
	r.Lock()
	defer r.Unlock()
	r.pkgMap[key] = value
}

// GetPackage gets the package path for given string value.
func GetPackage(pkgPath string) (*Package, bool) {
	p, ok := pkgMap.read(pkgPath)
	return p, ok
}

// PkgPath is the string package that contains full package name.
type PkgPath string

// Identifier gets package identifier.
func (p PkgPath) Identifier() string {
	v, ok := pkgMap.read(string(p))
	if ok {
		return v.Identifier
	}
	return ""
}

// FullName gets the full name of given PkgPath in a string type.
func (p PkgPath) FullName() string {
	return string(p)
}

// IsStandard checks if the package is standard.
func (p PkgPath) IsStandard() bool {
	return p == builtInPkgPath
}
