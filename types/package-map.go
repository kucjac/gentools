package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// PackageMap is a slice wrapper over Package type.
type PackageMap map[string]*Package

// NewPackage creates new package definition in given package map.
func (p PackageMap) NewPackage(name, identifier string) (*Package, error) {
	if name == "" {
		return nil, errors.New("empty package name")
	}
	if _, ok := p[name]; ok {
		return nil, errors.New("package with given name already defined")
	}
	pkg := NewPackage(name, identifier)
	p[name] = pkg
	return pkg, nil
}

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
			return &Array{Type: tp, ArrayKind: KindSlice}, true
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
		return &Array{Type: tp, ArrayKind: KindArray, ArraySize: size}, true
	case '*':
		if len(typeOf) == 1 {
			return nil, false
		}
		tp, ok := p.decomposeStringType(typeOf[1:], ctxPkg)
		if !ok {
			return nil, false
		}
		return &Pointer{PointedType: tp}, true
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
		return &Chan{Type: t, Dir: dir}, true
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
		return &Map{Key: kt, Value: vt}, true
	}

	var pkgPath string
	indexDot := strings.LastIndex(typeOf, ".")
	if indexDot != -1 {
		pkgPath = typeOf[:indexDot]
		typeOf = typeOf[indexDot+1:]
	}

	var isRoundBracket bool
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
			isRoundBracket = runes[i] == '('
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
	if pkgPath != "" {
		var ok bool
		if len(typeOf) == 0 {
			return nil, false
		}
		if thisPkg == nil || (thisPkg != nil && pkgPath != thisPkg.Path && pkgPath != thisPkg.Identifier) {
			if strings.ContainsRune(pkgPath, '/') {
				thisPkg, ok = p.PackageByPath(pkgPath)
			} else {
				thisPkg, ok = p.PackageByIdentifier(pkgPath)
			}
			if !ok {
				return nil, false
			}
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
	tp, ok := thisPkg.GetType(typeOf)
	if !ok {
		return nil, false
	}
	if next == "" || next == ")" || next == "}" {
		return tp, true
	}

	if next != "" {
		alias, isAlias := tp.(*Alias)
		if !isAlias {
			return nil, false
		}
		lastRune := runes[len(runes)-1]
		if isRoundBracket && lastRune != ')' {
			return nil, false
		}

		if !isRoundBracket && lastRune != '}' {
			return nil, false
		}
		zero := alias.Type.Zero(false, ctxPkg.Path)
		next = next[:len(next)-1]
		if next != zero {
			return nil, false
		}
	}
	return tp, true
}
