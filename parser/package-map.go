package parser

import (
	"sync"

	"github.com/kucjac/gentools/types"
)

type packageMap struct {
	sync.Mutex
	pkgMap             map[string]*types.Package
	pkgTypesinProgress map[*types.Package]map[string]types.Type
}

func (p *packageMap) read(key string) (*types.Package, bool) {
	p.Lock()
	defer p.Unlock()
	v, ok := p.pkgMap[key]
	return v, ok
}

func (p *packageMap) write(key string, value *types.Package) {
	p.Lock()
	defer p.Unlock()
	p.pkgMap[key] = value
}

// newPackage creates new package for given pkgPath and identifier.
func (p *packageMap) newPackage(pkgPath, identifier string) *types.Package {
	pkg := types.NewPackage(pkgPath, identifier)
	p.write(pkgPath, pkg)
	return pkg
}
