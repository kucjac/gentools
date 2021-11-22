package types

import (
	"fmt"
	"go/constant"
	"strings"
	"sync"
)

// Package is the golang package reflection container. It contains all interfaces, structs, functions
// and type wrappers that are located inside of it.
type Package struct {
	Path         string
	Identifier   string
	Interfaces   []*Interface
	Structs      []*Struct
	Functions    []*Function
	Aliases      []*Alias
	Types        map[string]Type
	Declarations map[string]Declaration
	sync.Mutex
}

// NewPackage creates new package definition.
func NewPackage(pkgPath, identifier string) *Package {
	return &Package{Path: pkgPath, Identifier: identifier, Types: map[string]Type{}, Declarations: map[string]Declaration{}}
}

// NewVariable adds new package variable Declaration.
func (p *Package) NewVariable(name string, tp Type) error {
	if err := p.hasName(name); err != nil {
		return err
	}
	p.Declarations[name] = Declaration{
		Name:     name,
		Type:     tp,
		Constant: false,
		Package:  p,
	}
	return nil
}

// NewConstant adds new package constant Declaration.
func (p *Package) NewConstant(name string, tp Type, value constant.Value) error {
	if err := p.hasName(name); err != nil {
		return err
	}
	p.Declarations[name] = Declaration{
		Name:     name,
		Type:     tp,
		Constant: true,
		Val:      value,
		Package:  p,
	}
	return nil
}

// NewNamedType adds new named type to the package definition.
func (p *Package) NewNamedType(name string, namedType Type) error {
	if err := p.hasName(name); err != nil {
		return err
	}
	switch nt := namedType.(type) {
	case *Alias:
		p.SetNamedType(name, nt)
	case *Interface:
		p.SetNamedType(name, nt)
	case *Struct:
		p.SetNamedType(name, nt)
	case *Function:
		p.SetNamedType(name, nt)
	default:
		return fmt.Errorf("invalid named type definition: %T", namedType)
	}
	return nil
}

func (p *Package) hasName(name string) error {
	if _, ok := p.Declarations[name]; ok {
		return fmt.Errorf("package: %s already contains definition with given name: '%s'", p.Path, name)
	}
	if _, ok := p.Types[name]; ok {
		return fmt.Errorf("package: '%s' already contains a type with given name: '%s'", p.Path, name)
	}
	return nil
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
func (p *Package) GetInterfaceType(name string) (*Interface, bool) {
	p.Lock()
	defer p.Unlock()
	t, ok := p.Types[name]
	if !ok {
		return nil, false
	}
	i, ok := t.(*Interface)
	return i, ok
}

// GetStruct gets the struct type by it's name.
func (p *Package) GetStruct(name string) (*Struct, bool) {
	p.Lock()
	defer p.Unlock()
	t, ok := p.Types[name]
	if !ok {
		return nil, false
	}
	s, ok := t.(*Struct)
	return s, ok
}

// MustStruct gets the struct declaration with the 'name'.
// If the struct is not found it panics.
func (p *Package) MustStruct(name string) *Struct {
	st, ok := p.GetStruct(name)
	if !ok {
		panic(fmt.Sprintf("struct: '%s' not found in the package: '%s'", name, p.Path))
	}
	return st
}

// MustFunction gets the function declaration with the 'name'.
// If the function is not found it panics.
func (p *Package) MustFunction(name string) *Function {
	f, ok := p.GetFunction(name)
	if !ok {
		panic(fmt.Sprintf("function %s not found in the package: '%s'", name, p.Path))
	}
	return f
}

// GetFunction gets the function type by it's name.
func (p *Package) GetFunction(name string) (*Function, bool) {
	p.Lock()
	defer p.Unlock()
	t, ok := p.Types[name]
	if !ok {
		return nil, false
	}
	s, ok := t.(*Function)
	return s, ok
}

// GetAlias gets the wrapped type by it's name.
func (p *Package) GetAlias(name string) (*Alias, bool) {
	p.Lock()
	defer p.Unlock()
	t, ok := p.Types[name]
	if !ok {
		return nil, false
	}
	s, ok := t.(*Alias)
	return s, ok
}

// IsStandard checks if given package is a standard package.
func (p *Package) IsStandard() bool {
	return p == builtIn
}

// SetNamedType sets the type with given name to be stored withing given package.
func (p *Package) SetNamedType(name string, tp Type) {
	p.Lock()
	defer p.Unlock()
	p.Types[name] = tp
	switch t := tp.(type) {
	case *Function:
		p.Functions = append(p.Functions, t)
	case *Struct:
		p.Structs = append(p.Structs, t)
	case *Alias:
		p.Aliases = append(p.Aliases, t)
	case *Interface:
		p.Interfaces = append(p.Interfaces, t)
	}
}

// SetIdentifier sets the package identifier.
func (p *Package) SetIdentifier(s string) {
	p.Identifier = s
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
