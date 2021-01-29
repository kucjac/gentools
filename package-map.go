package astreflect

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// PackageMap is a slice wrapper over Package type.
type PackageMap map[string]*Package

func (p PackageMap) MustGetByPath(path string) *Package {
	pkg, ok := p[path]
	if !ok {
		panic(fmt.Sprintf("Package: '%s' not found", path))
	}
	return pkg
}

// PackageByIdentifier gets the package by provided identifier. If there is more than one package with given identifier
// The function would return the first matching package.
func (p PackageMap) PackageByIdentifier(identifier string) (*Package, bool) {
	for _, pkg := range p {
		if pkg.Identifier == identifier {
			return pkg, true
		}
	}
	return nil, false
}

// PackageByPath gets the package by provided path.
func (p PackageMap) PackageByPath(path string) (*Package, bool) {
	pkg, ok := p[path]
	return pkg, ok
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
		thisPkg, ok = p.PackageByIdentifier(typeOf[:indexDot])
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

type packageMap struct {
	sync.Mutex
	pkgMap map[string]*Package
}

func (p *packageMap) read(key string) (*Package, bool) {
	p.Lock()
	defer p.Unlock()
	v, ok := p.pkgMap[key]
	return v, ok
}

func (p *packageMap) write(key string, value *Package) {
	p.Lock()
	defer p.Unlock()
	p.pkgMap[key] = value
}

// newPackage creates new package for given pkgPath and identifier.
func (p *packageMap) newPackage(pkgPath, identifier string) *Package {
	pkg := &Package{Path: pkgPath, Identifier: identifier, Types: map[string]Type{}, typesInProgress: map[string]Type{}}
	p.write(pkgPath, pkg)
	return pkg
}
