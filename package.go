package astreflect

import (
	"fmt"
	"strings"
	"sync"
)

// Package is the golang package reflection container. It contains all interfaces, structs, functions
// and type wrappers that are located inside of it.
type Package struct {
	Path            string
	Identifier      string
	Interfaces      []*InterfaceType
	Structs         []*StructType
	Functions       []*FunctionType
	WrappedTypes    []*WrappedType
	Types           map[string]Type
	typesInProgress map[string]Type
	sync.Mutex
}

// GetPkgPath gets the PkgPath for given package.
func (p *Package) GetPkgPath() PkgPath {
	if p.IsStandard() {
		return ""
	}
	return PkgPath(p.Identifier + " " + p.Path)
}

// MustGetType get the type with given 'name' from given package. If the type is not found the function panics.
func (p *Package) MustGetType(name string) Type {
	p.Lock()
	defer p.Unlock()
	t, ok := p.Types[name]
	if !ok {
		panic(fmt.Sprintf("Type: '%s' not found in the package: '%s'", name, p.Path))
	}
	return t
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
	return p == builtIn
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

// PkgPath is the string package that contains full package name.
type PkgPath string

// Identifier gets package identifier.
func (p PkgPath) Identifier() string {
	if i := strings.IndexRune(string(p), ' '); i != -1 {
		return string(p)[:i]
	}
	return ""
}

// FullName gets the full name of given PkgPath in a string type.
func (p PkgPath) FullName() string {
	if i := strings.IndexRune(string(p), ' '); i != -1 && len(p)-1 != i {
		return string(p)[i+1:]
	}
	return ""
}

// IsStandard checks if the package is standard.
func (p PkgPath) IsStandard() bool {
	return p == builtInPkgPath
}

func trimZeroRuneSpace(typeOf string) string {
	if len(typeOf) > 1 && typeOf[0] == ' ' {
		typeOf = typeOf[1:]
	}
	return typeOf
}
