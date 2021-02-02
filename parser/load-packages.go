package parser

import (
	"errors"
	"fmt"
	gotypes "go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"

	"github.com/kucjac/gentools/types"
)

// LoadConfig contains configuration used while loading packages.
type LoadConfig struct {
	// Paths could be absolute or relative path to given directory.
	Paths []string
	// PkgNames should be full pkg name i.e.: 'golang.org/x/mod/modfile'
	PkgNames []string
	// BuildFlags are the flags used by the ast.
	BuildFlags []string
	// Verbose sets the loader in verbose mode.
	Verbose bool
}

// LoadPackages parses Golang packages using AST.
func LoadPackages(cfg LoadConfig) (types.PackageMap, error) {
	pkgNames, err := getPackageNames(&cfg)
	if err != nil {
		return nil, err
	}

	p := &packageMap{pkgMap: types.PackageMap{}}
	pkgs, err := p.loadPackages(&cfg, pkgNames...)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, errors.New("no packages found")
	}
	p.parsePackages(&cfg, pkgs...)
	return p.pkgMap, nil
}

// UpdatePackages updates the packages in the given PackageMap.
// The function would get and parse only packages that doesn't currently exists in given map.
func UpdatePackages(p types.PackageMap, cfg LoadConfig) error {
	pkgNames, err := getPackageNames(&cfg)
	if err != nil {
		return err
	}
	pm := &packageMap{pkgMap: p}
	pkgNames = pm.resolveLoadedPackages(pkgNames)
	switch len(pkgNames) {
	case 0:
		if cfg.Verbose {
			fmt.Println("all packages from the input already loaded")
		}
	default:
		pkgs, err := pm.loadPackages(&cfg, pkgNames...)
		if err != nil {
			return err
		}
		if len(pkgs) == 0 {
			return errors.New("no packages found")
		}
		pm.parsePackages(&cfg, pkgs...)
	}
	return nil
}

func (p *packageMap) loadPackages(cfg *LoadConfig, pkgNames ...string) ([]*packages.Package, error) {
	now := time.Now()
	pkgCfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedDeps | packages.NeedImports | packages.NeedTypes,
		BuildFlags: cfg.BuildFlags,
	}

	pkgs, err := packages.Load(pkgCfg, pkgNames...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 1 {
		return nil, errors.New("error while loading import packages")
	}
	if cfg.Verbose {
		fmt.Printf("AST packages loaded in: %s\n", time.Since(now))
	}
	return pkgs, nil
}

func getPackageNames(cfg *LoadConfig) (pkgNames []string, err error) {
	for _, pkgPath := range cfg.Paths {
		if !filepath.IsAbs(pkgPath) {
			pkgPath, err = filepath.Abs(pkgPath)
			if err != nil {
				return nil, err
			}
		}
		pkgName, err := PackageNameOfDir(pkgPath)
		if err != nil {
			return nil, err
		}
		pkgNames = append(pkgNames, pkgName)
	}
	for _, imports := range cfg.PkgNames {
		pkgNames = append(pkgNames, imports)
	}
	return pkgNames, nil
}

func (p *packageMap) resolveLoadedPackages(pkgNames []string) (result []string) {
	for _, pkgName := range pkgNames {
		_, ok := p.read(pkgName)
		if !ok {
			result = append(result, pkgName)
		}
	}
	return result
}

func (p *packageMap) parsePackages(cfg *LoadConfig, newPkgs ...*packages.Package) {
	now := time.Now()

	var pkgs []*packages.Package
	for _, pkg := range newPkgs {
		// Check if the package is not already scanned.
		if _, ok := p.read(pkg.PkgPath); !ok {
			pkgs = append(pkgs, pkg)
		}
	}
	if len(pkgs) == 0 {
		return
	}

	initWg, finishGroup := &sync.WaitGroup{}, &sync.WaitGroup{}
	packageMap := map[string]*importedPackage{}
	for _, pkg := range pkgs {
		getAllImports(pkg.Types, packageMap)
	}
	pkgList := make([]*importedPackage, len(packageMap))
	var i int
	for _, v := range packageMap {
		pkgList[i] = v
		i++
	}
	sort.Slice(pkgList, func(i, j int) bool { return pkgList[i].importNo < pkgList[j].importNo })
	initWg.Add(len(packageMap))
	finishGroup.Add(len(packageMap))

	var rootPkgs []*rootPackage
	for _, importedPkg := range pkgList {
		rootPkg := &rootPackage{typesPkg: importedPkg.pkg, pkgMap: p, loadConfig: cfg, typesInProgress: map[string]types.Type{}}
		rootPkgs = append(rootPkgs, rootPkg)
		go rootPkg.parseTypePkg(initWg, finishGroup)
	}
	finishGroup.Wait()

	if cfg.Verbose {
		fmt.Printf("astreflect packages parsed in %s\n", time.Since(now))
	}

}

type importedPackage struct {
	pkg      *gotypes.Package
	importNo int
}

func getAllImports(pkg *gotypes.Package, imports map[string]*importedPackage) {
	imports[pkg.Path()] = &importedPackage{pkg: pkg, importNo: len(pkg.Imports())}
	for _, sub := range pkg.Imports() {
		if _, ok := imports[sub.Path()]; ok {
			continue
		}
		getAllImports(sub, imports)
	}
}

type rootPackage struct {
	typesPkg        *gotypes.Package
	refPkg          *types.Package
	pkgMap          *packageMap
	loadConfig      *LoadConfig
	declNames       []string
	typesInProgress map[string]types.Type
}

func (r *rootPackage) setTypeInProgress(name string, tp types.Type) {
	r.typesInProgress[name] = tp
	r.refPkg.SetNamedType(name, tp)
}

func (r *rootPackage) parseTypePkg(initWg, fg *sync.WaitGroup) {
	p := r.pkgMap.newPackage(r.typesPkg.Path(), r.typesPkg.Name())
	r.refPkg = p
	r.scaffoldPackageObjects()
	initWg.Done()
	initWg.Wait()

	s := r.typesPkg.Scope()
	typesScope := s
	type tuple struct {
		name string
		tp   types.Type
	}

	inProgress := make([]tuple, len(r.typesInProgress))
	var i int
	for name, tp := range r.typesInProgress {
		inProgress[i] = tuple{name, tp}
		i++
	}
	for _, tpl := range inProgress {
		name, tt := tpl.name, tpl.tp
		tp := typesScope.Lookup(name)
		switch t := tp.Type().(type) {
		case *gotypes.Named:
			if ok := r.finishNamedType(t, tt); !ok {
				if r.loadConfig.Verbose {
					fmt.Printf("package: %s, type: %s not mapped\n", p.Path, t.Obj().Name())
				}
				continue
			}
		case *gotypes.Signature:
			ok := r.finishNamedFunc(t, tt)
			if !ok {
				continue
			}
		default:
			wt, ok := tt.(*types.Alias)
			if !ok {
				continue
			}
			r.finishWrappedType(t.Underlying(), name, wt)
		}
	}
	for _, name := range r.declNames {
		obj := s.Lookup(name)
		switch ot := obj.(type) {
		case *gotypes.Const:
			declType, ok := r.dereferenceType(p, ot.Type())
			if !ok {
				continue
			}
			if err := p.NewConstant(ot.Name(), declType); err != nil {
				if r.loadConfig.Verbose {
					fmt.Println(err)
				}
			}
		case *gotypes.Var:
			declType, ok := r.dereferenceType(p, ot.Type())
			if !ok {
				continue
			}
			if err := p.NewVariable(ot.Name(), declType); err != nil {
				if r.loadConfig.Verbose {
					fmt.Println(err)
				}
			}
		default:
			continue
		}
	}
	fg.Done()
}

func (r *rootPackage) scaffoldPackageObjects() {
	s := r.typesPkg.Scope()
	for _, name := range s.Names() {
		obj := s.Lookup(name)
		switch obj.(type) {
		case *gotypes.TypeName, *gotypes.Func:
		case *gotypes.Const, *gotypes.Var:
			r.declNames = append(r.declNames, name)
			continue
		default:
			continue
		}

		switch ot := obj.Type().(type) {
		case *gotypes.Named:
			switch t := ot.Underlying().(type) {
			case *gotypes.Interface:
				it := &types.Interface{
					Pkg:           r.refPkg,
					InterfaceName: name,
					Methods:       make([]types.Function, t.NumMethods()),
				}
				r.setTypeInProgress(name, it)
			case *gotypes.Struct:
				st := &types.Struct{
					Pkg:      r.refPkg,
					TypeName: name,
					Fields:   make([]types.StructField, t.NumFields()),
				}
				r.setTypeInProgress(name, st)
			default:
				wt := &types.Alias{
					Pkg:       r.refPkg,
					AliasName: name,
				}
				r.setTypeInProgress(name, wt)
			}
		case *gotypes.Signature:
			if ot.Recv() != nil {
				// This is a method which should not be extracted as function.
				continue
			}
			fT := &types.Function{
				Pkg:      r.refPkg,
				FuncName: name,
			}
			r.setTypeInProgress(name, fT)
		}
	}
}

func (r *rootPackage) finishNamedType(named *gotypes.Named, t types.Type) bool {
	switch ot := t.(type) {
	case *types.Struct:
		return r.finishNamedStructType(ot, named)
	case *types.Interface:
		return r.finishNamedInterfaceType(named, ot)
	case *types.Alias:
		return r.finishWrappedType(named.Underlying(), named.Obj().Name(), ot)
	default:
		panic("invalid type")
	}
}

func (r *rootPackage) finishWrappedType(underlying gotypes.Type, name string, wt *types.Alias) bool {
	tp, ok := r.dereferenceType(r.refPkg, underlying)
	if !ok {
		fmt.Printf("Finishing WrappedType: %s failed: %s\n", name, underlying)
		return false
	}
	wt.Type = tp
	return true
}

func (r *rootPackage) getNamedType(et *gotypes.Named) (types.Type, bool) {
	if et.Obj().Pkg() == nil {
		t, ok := types.GetBuiltInType(et.Obj().Name())
		return t, ok
	}
	p, ok := r.pkgMap.read(et.Obj().Pkg().Path())
	if !ok {
		return nil, ok
	}
	return p.GetType(et.Obj().Name())
}

func (r *rootPackage) dereferenceType(p *types.Package, tp gotypes.Type) (types.Type, bool) {
	switch et := tp.(type) {
	case *gotypes.Named:
		return r.getNamedType(et)
	case *gotypes.Struct:
		return r.parseStructType(p, et)
	case *gotypes.Interface:
		return r.parseInterfaceType(p, et)
	case *gotypes.Basic:
		return types.GetBuiltInType(et.Name())
	case *gotypes.Slice:
		st, ok := r.dereferenceType(p, et.Elem())
		if !ok {
			return nil, false
		}
		return &types.Array{
			ArrayKind: types.KindSlice,
			Type:      st,
		}, true
	case *gotypes.Pointer:
		st, ok := r.dereferenceType(p, et.Elem())
		if !ok {
			return nil, ok
		}
		return &types.Pointer{PointedType: st}, true
	case *gotypes.Map:
		kt, ok := r.dereferenceType(p, et.Key())
		if !ok {
			return nil, false
		}
		vt, ok := r.dereferenceType(p, et.Elem())
		if !ok {
			return nil, false
		}
		return &types.Map{
			Key:   kt,
			Value: vt,
		}, true
	case *gotypes.Array:
		st, ok := r.dereferenceType(p, et.Elem())
		if !ok {
			return nil, ok
		}
		return &types.Array{
			ArrayKind: types.KindArray,
			Type:      st,
			ArraySize: int(et.Len()),
		}, true
	case *gotypes.Chan:
		st, ok := r.dereferenceType(p, et.Elem())
		if !ok {
			return nil, ok
		}
		return &types.Chan{Type: st, Dir: types.ChanDir(et.Dir())}, true
	case *gotypes.Signature:
		ft := &types.Function{Pkg: p}
		if !r.parseSignatureType(p, et, ft, true) {
			return nil, false
		}
		return ft, true
	default:
		fmt.Printf("type not found for dereferencing: %s, %T\n", et.String(), et)
		return nil, false
	}
}

func (r *rootPackage) finishNamedInterfaceType(named *gotypes.Named, intf *types.Interface) bool {
	p := r.refPkg
	it, ok := named.Underlying().(*gotypes.Interface)
	if !ok {
		return false
	}
	r.parseInterfaceMethods(p, it, intf)
	return true
}

func (r *rootPackage) parseInterfaceType(p *types.Package, it *gotypes.Interface) (*types.Interface, bool) {
	intf := &types.Interface{Pkg: p}
	if it.NumMethods() != 0 {
		intf.Methods = make([]types.Function, it.NumMethods())
		if ok := r.parseInterfaceMethods(p, it, intf); !ok {
			return nil, false
		}
	}
	return intf, true
}

func (r *rootPackage) parseInterfaceMethods(p *types.Package, it *gotypes.Interface, intf *types.Interface) bool {
	for i := 0; i < it.NumMethods(); i++ {
		xm, ok := r.parseMethod(p, it, i, false)
		if !ok {
			return ok
		}
		intf.Methods[i] = xm
	}
	// Sort methods by their names.
	sort.Slice(intf.Methods, func(i, j int) bool { return intf.Methods[i].FuncName < intf.Methods[j].FuncName })
	return true
}

func (r *rootPackage) finishNamedStructType(t *types.Struct, named *gotypes.Named) bool {
	p := r.refPkg
	r.parseStructFields(p, named.Underlying().(*gotypes.Struct), t)

	t.TypeName = named.Obj().Name()
	// Map methods.
	for i := 0; i < named.NumMethods(); i++ {
		xm, ok := r.parseMethod(p, named, i, true)
		if !ok {
			return ok
		}
		t.Methods = append(t.Methods, xm)
	}
	sort.Slice(t.Methods, func(i, j int) bool { return t.Methods[i].FuncName < t.Methods[j].FuncName })

	return true
}

func (r *rootPackage) parseStructType(p *types.Package, ot *gotypes.Struct) (*types.Struct, bool) {
	t := &types.Struct{Pkg: p}

	if ot.NumFields() != 0 {
		t.Fields = make([]types.StructField, ot.NumFields())
		if ok := r.parseStructFields(p, ot, t); !ok {
			return nil, false
		}
	}
	return t, true
}

func (r *rootPackage) parseStructFields(p *types.Package, ot *gotypes.Struct, t *types.Struct) bool {
	for i := 0; i < ot.NumFields(); i++ {
		f := ot.Field(i)
		ft, ok := r.dereferenceType(p, f.Type())
		if !ok {
			return ok
		}
		sField := types.StructField{
			Name:      f.Name(),
			Tag:       types.StructTag(ot.Tag(i)),
			Type:      ft,
			Index:     []int{i},
			Embedded:  f.Embedded(),
			Anonymous: f.Anonymous(),
		}
		t.Fields[i] = sField
	}
	return false
}

type astMethoder interface {
	Method(int) *gotypes.Func
}

func (r *rootPackage) parseMethod(p *types.Package, named astMethoder, i int, needReceiver bool) (types.Function, bool) {
	m := named.Method(i)

	s, ok := m.Type().(*gotypes.Signature)
	if !ok {
		fmt.Printf("method type is not a signature: %+v\n", m.Type())
		return types.Function{}, false
	}
	ft := types.Function{FuncName: m.Name(), Pkg: p}
	if ok = r.parseSignatureType(p, s, &ft, needReceiver); !ok {
		return types.Function{}, false
	}
	return ft, true
}

func (r *rootPackage) parseSignatureType(p *types.Package, s *gotypes.Signature, xm *types.Function, needReceiver bool) bool {
	xm.Variadic = s.Variadic()
	if needReceiver && s.Recv() != nil {
		xm.Receiver = &types.Receiver{Name: s.Recv().Name()}
		xm.Receiver.Type, _ = r.dereferenceType(p, s.Recv().Type())
	}
	if params := s.Params(); params != nil {
		xm.In = make([]types.FuncParam, params.Len())
		for j := 0; j < params.Len(); j++ {
			pm := params.At(j)
			in := types.FuncParam{Name: pm.Name()}
			pt, ok := r.dereferenceType(p, pm.Type())
			if !ok {
				return false
			}
			in.Type = pt
			xm.In[j] = in
		}
	}

	if results := s.Results(); results != nil {
		xm.Out = make([]types.FuncParam, results.Len())
		for j := 0; j < results.Len(); j++ {
			pm := results.At(j)
			out := types.FuncParam{Name: pm.Name()}
			pt, ok := r.dereferenceType(p, pm.Type())
			if !ok {
				return false
			}
			out.Type = pt
			xm.Out[j] = out
		}
	}
	return true
}

func (r *rootPackage) finishNamedFunc(st *gotypes.Signature, t types.Type) bool {
	ft, ok := t.(*types.Function)
	if !ok {
		return false
	}
	if !r.parseSignatureType(r.refPkg, st, ft, false) {
		return false
	}
	return true
}

// PackageNameOfDir get package import path via dir
func PackageNameOfDir(srcDir string) (string, error) {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return "", err
	}

	var goFilePath string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			goFilePath = file.Name()
			break
		}
	}
	if goFilePath == "" {
		return "", fmt.Errorf("go source file not found %s", srcDir)
	}

	packageImport, err := parsePackageImport(srcDir)
	if err != nil {
		return "", err
	}
	return packageImport, nil
}

var errOutsideGoPath = errors.New("source directory is outside GOPATH")

func parsePackageImport(srcDir string) (string, error) {
	moduleMode := os.Getenv("GO111MODULE")
	// trying to find the module
	if moduleMode != "off" {
		currentDir := srcDir
		for {
			dat, err := ioutil.ReadFile(filepath.Join(currentDir, "go.mod"))
			if os.IsNotExist(err) {
				if currentDir == filepath.Dir(currentDir) {
					// at the root
					break
				}
				currentDir = filepath.Dir(currentDir)
				continue
			} else if err != nil {
				return "", err
			}
			modulePath := modfile.ModulePath(dat)
			return filepath.ToSlash(filepath.Join(modulePath, strings.TrimPrefix(srcDir, currentDir))), nil
		}
	}
	// fall back to GOPATH mode
	goPaths := os.Getenv("GOPATH")
	if goPaths == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}
	goPathList := strings.Split(goPaths, string(os.PathListSeparator))
	for _, goPath := range goPathList {
		sourceRoot := filepath.Join(goPath, "src") + string(os.PathSeparator)
		if strings.HasPrefix(srcDir, sourceRoot) {
			return filepath.ToSlash(strings.TrimPrefix(srcDir, sourceRoot)), nil
		}
	}
	return "", errOutsideGoPath
}
