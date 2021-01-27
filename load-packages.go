package astreflect

import (
	"errors"
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
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
func LoadPackages(cfg *LoadConfig) (PackageMap, error) {
	pkgNames, err := getPackageNames(cfg)
	if err != nil {
		return nil, err
	}
	pkgNames = resolveLoadedPackages(pkgNames)
	switch len(pkgNames) {
	case 0:
		if cfg.Verbose {
			fmt.Println("all packages from the input already loaded")
		}
	default:
		pkgs, err := loadPackages(cfg, pkgNames...)
		if err != nil {
			return nil, err
		}
		if len(pkgs) == 0 {
			return nil, errors.New("no packages found")
		}
		parsePackages(cfg, pkgs...)
	}
	return getPackageMap(), nil
}

// GetPackages gets the package map.
func GetPackages() PackageMap {
	return getPackageMap()
}

func getPackageMap() PackageMap {
	return pkgMap.pkgMap
}

func loadPackages(cfg *LoadConfig, pkgNames ...string) ([]*packages.Package, error) {
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

func resolveLoadedPackages(pkgNames []string) (result []string) {
	for _, pkgName := range pkgNames {
		_, ok := pkgMap.read(pkgName)
		if !ok {
			result = append(result, pkgName)
		}
	}
	return result
}

func parsePackages(cfg *LoadConfig, newPkgs ...*packages.Package) {
	now := time.Now()

	var pkgs []*packages.Package
	for _, pkg := range newPkgs {
		if _, ok := pkgMap.read(pkg.PkgPath); !ok {
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
		rootPkg := &rootPackage{typesPkg: importedPkg.pkg}
		rootPkgs = append(rootPkgs, rootPkg)
		go rootPkg.parseTypePkg(initWg, finishGroup)
	}
	finishGroup.Wait()

	if cfg.Verbose {
		fmt.Printf("astreflect packages parsed in %s\n", time.Since(now))
	}
}

type importedPackage struct {
	pkg      *types.Package
	importNo int
}

func getAllImports(pkg *types.Package, imports map[string]*importedPackage) {
	imports[pkg.Path()] = &importedPackage{pkg: pkg, importNo: len(pkg.Imports())}
	for _, sub := range pkg.Imports() {
		if _, ok := imports[sub.Path()]; ok {
			continue
		}
		getAllImports(sub, imports)
	}
}

type rootPackage struct {
	typesPkg *types.Package
	refPkg   *Package
}

func (r *rootPackage) parseTypePkg(initWg, fg *sync.WaitGroup) {
	p := newPackage(r.typesPkg.Path(), r.typesPkg.Name())
	r.refPkg = p
	r.scaffoldPackageObjects()
	initWg.Done()
	initWg.Wait()

	typesScope := r.typesPkg.Scope()
	type tuple struct {
		name string
		tp   Type
	}

	inProgress := make([]tuple, len(p.typesInProgress))
	var i int
	for name, tp := range p.typesInProgress {
		inProgress[i] = tuple{name, tp}
		i++
	}
	for _, tpl := range inProgress {
		name, tt := tpl.name, tpl.tp
		tp := typesScope.Lookup(name)
		switch t := tp.Type().(type) {
		case *types.Named:
			if ok := r.finishNamedType(t, tt); !ok {
				continue
			}
		case *types.Signature:
			ok := r.finishNamedFunc(name, t, tt)
			if !ok {
				continue
			}
		default:
			wt, ok := tt.(*WrappedType)
			if !ok {
				continue
			}
			r.finishWrappedType(t.Underlying(), name, wt)
		}
	}
	fg.Done()
}

func (r *rootPackage) scaffoldPackageObjects() {
	s := r.typesPkg.Scope()
	for _, name := range s.Names() {
		obj := s.Lookup(name)
		switch obj.(type) {
		case *types.TypeName, *types.Func:
		default:
			continue
		}

		switch ot := obj.Type().(type) {
		case *types.Named:
			switch t := ot.Underlying().(type) {
			case *types.Interface:
				it := &InterfaceType{
					PackagePath:   r.refPkg.Path,
					InterfaceName: name,
					Methods:       make([]FunctionType, t.NumMethods()),
				}
				r.refPkg.setInProgressType(name, it)
			case *types.Struct:
				st := &StructType{
					PackagePath: r.refPkg.Path,
					TypeName:    name,
					Fields:      make([]StructField, t.NumFields()),
				}
				r.refPkg.setInProgressType(name, st)
			default:
				wt := &WrappedType{
					PackagePath: r.refPkg.Path,
					WrapperName: name,
				}
				r.refPkg.setInProgressType(name, wt)
			}
		case *types.Signature:
			if ot.Recv() != nil {
				// This is a method which should not be extracted as function.
				continue
			}
			fT := &FunctionType{
				PackagePath: r.refPkg.Path,
				FuncName:    name,
			}
			r.refPkg.setInProgressType(name, fT)
		}
	}
}

func (r *rootPackage) finishNamedType(named *types.Named, t Type) bool {
	switch ot := t.(type) {
	case *StructType:
		return r.finishNamedStructType(ot, named)
	case *InterfaceType:
		return r.finishNamedInterfaceType(named, ot)
	case *WrappedType:
		return r.finishWrappedType(named.Underlying(), named.Obj().Name(), ot)
	default:
		panic("invalid type")
	}
}

func (r *rootPackage) finishWrappedType(underlying types.Type, name string, wt *WrappedType) bool {
	tp, ok := r.dereferenceType(r.refPkg, underlying)
	if !ok {
		fmt.Printf("Finishing WrappedType: %s failed: %s\n", name, underlying)
		return false
	}
	wt.Type = tp
	r.refPkg.markTypeDone(name)
	return true
}

func getNamedType(et *types.Named) (Type, bool) {
	if et.Obj().Pkg() == nil {
		t, ok := GetBuiltInType(et.Obj().Name())
		return t, ok
	}
	p, ok := GetPackage(et.Obj().Pkg().Path())
	if !ok {
		return nil, ok
	}
	return p.GetType(et.Obj().Name())
}

func (r *rootPackage) dereferenceType(p *Package, tp types.Type) (Type, bool) {
	switch et := tp.(type) {
	case *types.Named:
		return getNamedType(et)
	case *types.Struct:
		return r.parseStructType(p, et)
	case *types.Interface:
		return r.parseInterfaceType(p, et)
	case *types.Basic:
		return GetBuiltInType(et.Name())
	case *types.Slice:
		st, ok := r.dereferenceType(p, et.Elem())
		if !ok {
			return nil, false
		}
		return &ArrayType{
			ArrayKind: Slice,
			Type:      st,
		}, true
	case *types.Pointer:
		st, ok := r.dereferenceType(p, et.Elem())
		if !ok {
			return nil, ok
		}
		return &PointerType{PointedType: st}, true
	case *types.Map:
		kt, ok := r.dereferenceType(p, et.Key())
		if !ok {
			return nil, false
		}
		vt, ok := r.dereferenceType(p, et.Elem())
		if !ok {
			return nil, false
		}
		return &MapType{
			Key:   kt,
			Value: vt,
		}, true
	case *types.Array:
		st, ok := r.dereferenceType(p, et.Elem())
		if !ok {
			return nil, ok
		}
		return &ArrayType{
			ArrayKind: Array,
			Type:      st,
			ArraySize: int(et.Len()),
		}, true
	case *types.Chan:
		st, ok := r.dereferenceType(p, et.Elem())
		if !ok {
			return nil, ok
		}
		return &ChanType{Type: st, Dir: ChanDir(et.Dir())}, true
	case *types.Signature:
		ft := &FunctionType{PackagePath: p.Path}
		if !r.parseSignatureType(p, et, ft, true) {
			return nil, false
		}
		return ft, true
	default:
		fmt.Printf("type not found for dereferencing: %s, %T\n", et.String(), et)
		return nil, false
	}
}

func (r *rootPackage) finishNamedInterfaceType(named *types.Named, intf *InterfaceType) bool {
	p := r.refPkg
	it, ok := named.Underlying().(*types.Interface)
	if !ok {
		return false
	}
	r.parseInterfaceMethods(p, it, intf)
	p.markTypeDone(named.Obj().Name())
	return true
}

func (r *rootPackage) parseInterfaceType(p *Package, it *types.Interface) (*InterfaceType, bool) {
	intf := &InterfaceType{PackagePath: p.Path}
	if it.NumMethods() != 0 {
		intf.Methods = make([]FunctionType, it.NumMethods())
		if ok := r.parseInterfaceMethods(p, it, intf); !ok {
			return nil, false
		}
	}
	return intf, true
}

func (r *rootPackage) parseInterfaceMethods(p *Package, it *types.Interface, intf *InterfaceType) bool {
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

func (r *rootPackage) finishNamedStructType(t *StructType, named *types.Named) bool {
	p := r.refPkg
	r.parseStructFields(p, named.Underlying().(*types.Struct), t)

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

	p.markTypeDone(named.Obj().Name())
	return true
}

func (r *rootPackage) parseStructType(p *Package, ot *types.Struct) (*StructType, bool) {
	t := &StructType{PackagePath: p.Path}

	if ot.NumFields() != 0 {
		t.Fields = make([]StructField, ot.NumFields())
		if ok := r.parseStructFields(p, ot, t); !ok {
			return nil, false
		}
	}
	return t, true
}

func (r *rootPackage) parseStructFields(p *Package, ot *types.Struct, t *StructType) bool {
	for i := 0; i < ot.NumFields(); i++ {
		f := ot.Field(i)
		ft, ok := r.dereferenceType(p, f.Type())
		if !ok {
			return ok
		}
		sField := StructField{
			Name:      f.Name(),
			Tag:       StructTag(ot.Tag(i)),
			Type:      ft,
			Index:     []int{i},
			Embedded:  f.Embedded(),
			Anonymous: f.Anonymous(),
		}
		t.Fields[i] = sField
	}
	return false
}

type methoder interface {
	Method(int) *types.Func
}

func (r *rootPackage) parseMethod(p *Package, named methoder, i int, needReceiver bool) (FunctionType, bool) {
	m := named.Method(i)

	s, ok := m.Type().(*types.Signature)
	if !ok {
		fmt.Printf("method type is not a signature: %+v\n", m.Type())
		return FunctionType{}, false
	}
	ft := FunctionType{FuncName: m.Name(), PackagePath: p.Path}
	if ok = r.parseSignatureType(p, s, &ft, needReceiver); !ok {
		return FunctionType{}, false
	}
	return ft, true
}

func (r *rootPackage) parseSignatureType(p *Package, s *types.Signature, xm *FunctionType, needReceiver bool) bool {
	xm.Variadic = s.Variadic()
	if needReceiver && s.Recv() != nil {
		xm.Receiver = &Receiver{Name: s.Recv().Name()}
		_, xm.Receiver.Pointer = s.Recv().Type().(*types.Pointer)
	}
	if params := s.Params(); params != nil {
		xm.In = make([]IOParam, params.Len())
		for j := 0; j < params.Len(); j++ {
			pm := params.At(j)
			in := IOParam{Name: pm.Name()}
			pt, ok := r.dereferenceType(p, pm.Type())
			if !ok {
				return false
			}
			in.Type = pt
			xm.In[j] = in
		}
	}

	if results := s.Results(); results != nil {
		xm.Out = make([]IOParam, results.Len())
		for j := 0; j < results.Len(); j++ {
			pm := results.At(j)
			out := IOParam{Name: pm.Name()}
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

func (r *rootPackage) finishNamedFunc(name string, st *types.Signature, t Type) bool {
	ft, ok := t.(*FunctionType)
	if !ok {
		return false
	}
	if !r.parseSignatureType(r.refPkg, st, ft, false) {
		return false
	}
	r.refPkg.markTypeDone(name)
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
