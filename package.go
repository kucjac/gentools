package astreflect

import (
	"fmt"
	"strconv"
	"strings"
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

func (p PackageMap) MustGetByPath(path string) *Package {
	pkg, ok := p[path]
	if !ok {
		panic(fmt.Sprintf("Package: '%s' not found", path))
	}
	return pkg
}

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
	pkg, ok := p[path]
	return pkg, ok
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

// TypeOf gets the resulting type for provided 'typeOf'.
// If the packageContext is defined the values without identifier would be found within given package also.
func (p *PackageMap) TypeOf(typeOf string, packageContext *Package) (Type, bool) {
	if typeOf == "" {
		return nil, false
	}
	return p.decomposeStringType(typeOf, packageContext)
}

func (p *PackageMap) decomposeStringType(typeOf string, ctxPkg *Package) (Type, bool) {
	if typeOf == "" {
		return nil, false
	}
	switch typeOf[0] {
	case '[':
		closing := strings.IndexRune(typeOf, ']')
		if closing == -1 {
			return nil, false
		}

		if closing == 1 {
			if len(typeOf) == 2 {
				return nil, false
			}
			tp, ok := p.decomposeStringType(typeOf[2:], ctxPkg)
			if !ok {
				return nil, false
			}
			return ArrayType{Type: tp, ArrayKind: Slice}, true
		}

		size, err := strconv.Atoi(typeOf[1:closing])
		if err != nil {
			// TODO: support constant base size of array.
			return nil, false
		}
		tp, ok := p.decomposeStringType(typeOf[closing+1:], ctxPkg)
		if !ok {
			return nil, false
		}
		return ArrayType{Type: tp, ArrayKind: Array, ArraySize: size}, true
	case '*':
		if len(typeOf) == 1 {
			return nil, false
		}
		tp, ok := p.decomposeStringType(typeOf[1:], ctxPkg)
		if !ok {
			return nil, false
		}
		return PointerType{PointedType: tp}, true
	}

	var dir ChanDir
	if strings.HasPrefix(typeOf, "<-") {
		dir = RecvOnly
		typeOf = strings.TrimPrefix(typeOf, "<-")
		typeOf = trimZeroRuneSpace(typeOf)
	}
	if strings.HasPrefix(typeOf, "chan") {
		typeOf = strings.TrimPrefix(typeOf, "chan")
		typeOf = trimZeroRuneSpace(typeOf)
		if strings.HasPrefix(typeOf, "<-") {
			if dir != 0 {
				return nil, false
			}
			dir = SendOnly
			typeOf = strings.TrimPrefix(typeOf, "<-")
			typeOf = trimZeroRuneSpace(typeOf)
		}
		t, ok := p.decomposeStringType(typeOf, ctxPkg)
		if !ok {
			return nil, false
		}
		return ChanType{Type: t, Dir: dir}, true
	} else if dir != 0 {
		return nil, false
	}

	if strings.HasPrefix(typeOf, "map") {
		typeOf = strings.TrimPrefix(typeOf, "map")
		if len(typeOf) == 0 {
			return nil, false
		}
		if typeOf[0] != '[' {
			return nil, false
		}
		typeOf = typeOf[1:]
		closing := strings.IndexRune(typeOf, ']')
		if closing == -1 {
			return nil, false
		}
		if len(typeOf)-1 == closing {
			return nil, false
		}
		key := typeOf[:closing]
		value := typeOf[closing+1:]
		kt, ok := p.decomposeStringType(key, ctxPkg)
		if !ok {
			return nil, false
		}
		vt, ok := p.decomposeStringType(value, ctxPkg)
		if !ok {
			return nil, false
		}
		return MapType{Key: kt, Value: vt}, true
	}

	indexNext, indexDot, indexBracket := -1, -1, -1
	runes := []rune(typeOf)
	for i := 0; i < len(typeOf); i++ {
		switch runes[i] {
		case '.':
			if indexDot != -1 {
				return nil, false
			}
			indexDot = i
		case '(', '{':
			if indexBracket != -1 {
				return nil, false
			}
			indexBracket = i
		case '*', '[', '<':
			indexNext = i
		case 'c':
			if strings.HasPrefix(typeOf[i:], "chan") {
				indexNext = i
			}
		}
		if indexNext != -1 {
			break
		}
	}

	thisPkg := ctxPkg
	if indexDot != -1 && (indexNext > indexDot || indexNext == -1) {
		var ok bool
		if len(typeOf)-1 == indexDot {
			return nil, false
		}
		thisPkg, ok = p.GetByIdentifier(typeOf[:indexDot])
		if !ok {
			return nil, false
		}
		typeOf = typeOf[indexDot+1:]
		if indexBracket != -1 {
			indexBracket -= indexDot
		}
	}
	var next string
	if indexBracket != -1 {
		if len(typeOf)-1 == indexBracket {
			return nil, false
		}
		next = typeOf[indexBracket+1:]
		typeOf = typeOf[:indexBracket]
	}
	bt, ok := GetBuiltInType(typeOf)
	if ok {
		return bt, true
	}
	if thisPkg == nil {
		return nil, false
	}
	tp, ok := thisPkg.GetType(typeOf)
	if !ok {
		return nil, false
	}
	if next == "" || next == "()" || next == "{}" {
		return tp, true
	}
	if next != "" && next[0] == '0' && tp.Kind() != Wrapper {
		return nil, false
	}
	return tp, true
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

func trimZeroRuneSpace(typeOf string) string {
	if len(typeOf) > 1 && typeOf[0] == ' ' {
		typeOf = typeOf[1:]
	}
	return typeOf
}
