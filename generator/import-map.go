package genutils

import (
	"github.com/kucjac/gentools/types"
)

func newImportMap(pkg *types.Package) *ImportMap {
	return &ImportMap{Package: pkg, Map: map[string]*types.Package{}}
}

// ImportMap is the map of imports for given type.
type ImportMap struct {
	Package *types.Package
	Map     map[string]*types.Package
}

// Sorted gets sorted imports from provided map.
func (im *ImportMap) Sorted() *Imports {
	pkgs := make([]*types.Package, len(im.Map))
	var i int
	for _, pkg := range im.Map {
		pkgs[i] = pkg
		i++
	}
	return SortImports(im.Package.Path, pkgs)
}

// // AddFunctionImports adds the imports from given function type.
// func (im *ImportMap) AddFunctionImports(ctxPackage string, f *types.FunctionType) {
// 	for i := range f.In {
// 		im.addImportPackages(ctxPackage, getTypePackages(f.In[i].Type)...)
// 	}
// 	for i := range f.Out {
// 		im.addImportPackages(ctxPackage, getTypePackages(f.Out[i].Type)...)
// 	}
// 	if f.Receiver != nil {
// 		im.addImportPackages(ctxPackage, getTypePackages(f.Receiver.Type)...)
// 	}
// }

func (im *ImportMap) addImportPackages(packages ...*types.Package) {
	for _, pkg := range packages {
		if pkg.Path == im.Package.Path {
			continue
		}
		im.Map[pkg.Path] = pkg
	}
}

func getTypePackages(tp types.Type) []*types.Package {
	for tp != nil {
		switch t := tp.(type) {
		case types.Packager:
			return []*types.Package{t.Package()}
		case *types.Map:
			k, v := getTypePackages(t.Key), getTypePackages(t.Value)
			return append(k, v...)
		}
		tp = tp.Elem()
	}
	return nil
}
