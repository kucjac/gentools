package parser

import (
	"errors"
	"fmt"
	"go/ast"
	gotypes "go/types"
	"io/ioutil"
	"log"
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
	// WithComments
	WithComments bool
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
	if err := p.parsePackages(&cfg, pkgs...); err != nil {
		return nil, err
	}
	return p.pkgMap, nil
}

// UpdatePackages updates the packages in the given PackageMap.
// The function would get and parse only packages that doesn't currently exists in given map.
func UpdatePackages(p types.PackageMap, cfg LoadConfig) error {
	pkgNames, err := getPackageNames(&cfg)
	if err != nil {
		return err
	}

	// Parse into a copy so callers never observe a partially initialized map
	// when loading or conversion fails.
	staged := make(types.PackageMap, len(p))
	for path, pkg := range p {
		staged[path] = pkg
	}
	pm := &packageMap{pkgMap: staged}
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
		if err := pm.parsePackages(&cfg, pkgs...); err != nil {
			return err
		}
		for path, pkg := range staged {
			if _, exists := p[path]; !exists {
				p[path] = pkg
			}
		}
	}
	return nil
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

func (p *packageMap) loadPackages(cfg *LoadConfig, pkgNames ...string) ([]*packages.Package, error) {
	now := time.Now()
	mode := packages.NeedName | packages.NeedDeps | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo
	pkgCfg := &packages.Config{
		Mode:       mode,
		BuildFlags: cfg.BuildFlags,
	}

	pkgs, err := packages.Load(pkgCfg, pkgNames...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		var diagnostics []string
		for _, pkg := range pkgs {
			for _, diagnostic := range pkg.Errors {
				diagnostics = append(diagnostics, diagnostic.Error())
			}
		}
		return nil, fmt.Errorf("package load failed: %s", strings.Join(diagnostics, "; "))
	}
	if cfg.Verbose {
		log.Printf("AST packages loaded in: %s\n", time.Since(now))
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
	pkgNames = append(pkgNames, cfg.PkgNames...)
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

func (p *packageMap) parsePackages(cfg *LoadConfig, newPkgs ...*packages.Package) error {
	now := time.Now()

	var pkgs []*packages.Package
	for _, pkg := range newPkgs {
		// Check if the package is not already scanned.
		if _, ok := p.read(pkg.PkgPath); !ok {
			pkgs = append(pkgs, pkg)
		}
	}
	if len(pkgs) == 0 {
		return nil
	}

	initWg, finishGroup := &sync.WaitGroup{}, &sync.WaitGroup{}
	packageMap := map[string]*importedPackage{}
	for _, pkg := range pkgs {
		getAllImports(pkg, packageMap)
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

	rootPkgs := map[*gotypes.Package]*rootPackage{}
	for _, importedPkg := range pkgList {
		existing, _ := p.read(importedPkg.typesPkg.Path())
		rootPkg := &rootPackage{
			rootPackages:    rootPkgs,
			pkgPkg:          importedPkg.pkgPkg,
			typesPkg:        importedPkg.typesPkg,
			refPkg:          existing,
			pkgMap:          p,
			loadConfig:      cfg,
			typesInProgress: map[string]types.Type{},
			typeParameters:  map[gotypes.Type]*types.TypeParameter{},
			mappedAliases:   map[string]struct{}{},
			namedAliases:    map[string]*gotypes.Named{},
		}
		rootPkgs[importedPkg.typesPkg] = rootPkg
	}

	errs := make(chan error, len(rootPkgs))
	for _, rootPkg := range rootPkgs {
		go rootPkg.parseTypePkg(initWg, finishGroup, errs)
	}

	finishGroup.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			return err
		}
	}

	if cfg.Verbose {
		log.Printf("gentools packages parsed in %s\n", time.Since(now))
	}
	return nil
}

type importedPackage struct {
	pkgPkg   *packages.Package
	typesPkg *gotypes.Package
	importNo int
}

func getAllImports(pkg *packages.Package, imports map[string]*importedPackage) {
	typesPkg := pkg.Types
	imports[typesPkg.Path()] = &importedPackage{pkgPkg: pkg, typesPkg: typesPkg, importNo: len(typesPkg.Imports())}
	for path, sub := range pkg.Imports {
		if _, ok := imports[path]; ok {
			continue
		}
		getAllImports(sub, imports)
	}
}

type rootPackage struct {
	pkgPkg          *packages.Package
	typesPkg        *gotypes.Package
	refPkg          *types.Package
	pkgMap          *packageMap
	rootPackages    map[*gotypes.Package]*rootPackage
	mappedAliases   map[string]struct{}
	namedAliases    map[string]*gotypes.Named
	loadConfig      *LoadConfig
	declNames       []string
	typesInProgress map[string]types.Type
	typeParameters  map[gotypes.Type]*types.TypeParameter
}

func (r *rootPackage) setTypeInProgress(name string, tp types.Type) {
	r.typesInProgress[name] = tp
	r.refPkg.SetNamedType(name, tp)
}

func (r *rootPackage) parseTypePkg(initWg, fg *sync.WaitGroup, errs chan<- error) {
	defer fg.Done()
	if r.refPkg != nil {
		initWg.Done()
		initWg.Wait()
		return
	}

	p := r.pkgMap.newPackage(r.typesPkg.Path(), r.typesPkg.Name())

	r.refPkg = p
	r.scaffoldPackageObjects()
	initWg.Done()
	initWg.Wait()

	s := r.typesPkg.Scope()
	r.resolveInProgressTypes(s, p)
	if err := r.resolveIdentAliases(); err != nil {
		errs <- fmt.Errorf("resolve aliases in package %s: %w", r.typesPkg.Path(), err)
		return
	}
	r.defineDeclarations(s, p)
	if r.loadConfig.WithComments {
		r.parseComments(p)
	}
}

// resolveIdentAliases resolves compatibility aliases exactly once.  The former
// retry queue appended an unresolved TypeSpec while using the growing slice as
// its loop bound, so a permanently unresolved type grew memory without limit.
// Type information is already available from go/packages and provides the
// semantic RHS directly, including generic instantiations.
func (r *rootPackage) resolveIdentAliases() error {
	for _, file := range r.pkgPkg.Syntax {
		for _, decl := range file.Decls {
			dt, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range dt.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if _, isAlias := r.mappedAliases[ts.Name.Name]; !isAlias {
					continue
				}
				named, ok := r.namedAliases[ts.Name.Name]
				if !ok {
					continue
				}

				ip, ok := r.typesInProgress[ts.Name.Name]
				if !ok {
					continue
				}
				alias, ok := ip.(*types.Alias)
				if !ok {
					return fmt.Errorf("%s: expected compatibility alias, got %T", ts.Name.Name, ip)
				}

				rhs := r.pkgPkg.TypesInfo.TypeOf(ts.Type)
				if rhs == nil {
					return fmt.Errorf("%s: missing type information for alias target", ts.Name.Name)
				}
				underlying, ok := r.dereferenceType(r.refPkg, rhs)
				if !ok {
					return fmt.Errorf("%s: cannot resolve alias target %s", ts.Name.Name, rhs.String())
				}
				if !r.finishNamedAliasType(named, alias, underlying) {
					return fmt.Errorf("%s: cannot finish compatibility alias", ts.Name.Name)
				}
			}
		}
	}
	return nil
}

func (r *rootPackage) resolveInProgressTypes(s *gotypes.Scope, p *types.Package) {
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
			if _, isLegacyAlias := r.mappedAliases[name]; isLegacyAlias {
				r.namedAliases[name] = t
				continue
			}
			if ok := r.finishNamedType(t, tt); !ok {
				if r.loadConfig.Verbose {
					log.Printf("package: %s, type: %s not mapped\n", p.Path, t.Obj().Name())
				}
				continue
			}
		case *gotypes.Signature:
			if !r.finishNamedFunc(t, tt) {
				continue
			}
		default:
			alias, ok := tt.(*types.Alias)
			typeName, isTypeName := tp.(*gotypes.TypeName)
			if !ok || !isTypeName || !r.finishAliasType(p, typeName, alias) {
				if r.loadConfig.Verbose {
					log.Printf("package: %s, alias: %s not mapped\n", p.Path, name)
				}
			}
		}
	}
}

func (r *rootPackage) defineDeclarations(s *gotypes.Scope, p *types.Package) {
	for _, name := range r.declNames {
		obj := s.Lookup(name)
		switch ot := obj.(type) {
		case *gotypes.Const:
			declType, ok := r.dereferenceType(p, ot.Type())
			if !ok {
				continue
			}
			if err := p.NewConstant(ot.Name(), declType, ot.Val()); err != nil {
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
}

func (r *rootPackage) parseComments(p *types.Package) {
	for _, file := range r.pkgPkg.Syntax {
	declLoop:
		for _, decl := range file.Decls {
			switch dt := decl.(type) {
			case *ast.GenDecl:
			specLoop:
				for _, spec := range dt.Specs {
					switch st := spec.(type) {
					case *ast.TypeSpec:
						tp, ok := p.Types[st.Name.Name]
						if !ok {
							if r.loadConfig.Verbose {
								log.Printf("type: '%s' not found in the pkg declaration", st.Name.Name)
							}
							continue
						}
						var comment string
						if st.Doc != nil {
							comment = st.Doc.Text()
						}
						if comment == "" && dt.Doc != nil {
							comment = dt.Doc.Text()
						}

					ptrLoop:
						for {
							switch tt := tp.(type) {
							case *types.Pointer:
								// Dereference the type.
								tp = tt.Elem()
								continue ptrLoop
							case *types.Struct:
								tt.Comment = comment
								structType, ok := r.extractStructExpr(file, st.Type)
								if !ok {
									if r.loadConfig.Verbose {
										log.Printf("getting ast struct type failed. package: '%s', name: '%s'\n", r.refPkg.Path, st.Name.Name)
									}
									continue specLoop
								}

								for j, field := range structType.Fields.List {
									if j >= len(tt.Fields) {
										break
									}
									var fc string
									if field.Doc != nil {
										fc = field.Doc.Text()
									}
									if len(tt.Fields) != structType.Fields.NumFields() {
										fmt.Println("Different!")
									}
									tt.Fields[j].Comment = fc
								}
							case *types.Interface:
								tt.Comment = comment

								interfaceType, ok := r.extractInterfaceExpr(file, st.Type)
								if !ok {
									if r.loadConfig.Verbose {
										log.Printf("Getting ast interface type failed.  Package: '%s', name: '%s'", r.pkgPkg.Name, st.Name.Name)
									}
									continue specLoop
								}

								for j, method := range interfaceType.Methods.List {
									if j >= len(tt.Methods) {
										break
									}
									var fc string
									if method.Doc != nil {
										fc = method.Doc.Text()
									}
									tt.Methods[j].Comment = fc
								}
							case *types.Alias:
								tt.Comment = comment
								tp = tt.Type
								continue ptrLoop
							case *types.Function:
								tt.Comment = comment
							}
							continue specLoop
						}
					case *ast.ValueSpec:
						if st.Doc == nil {
							continue
						}

						for _, name := range st.Names {
							decl, ok := p.Declarations[name.Name]
							if !ok {
								continue
							}
							decl.Comment = st.Doc.Text()
							p.Declarations[name.Name] = decl
						}
					}
				}
			case *ast.FuncDecl:
				if dt.Doc == nil {
					continue
				}

				// Check if it is a method or a func.
				var funType *types.Function
				if dt.Recv != nil {
					for _, rcv := range dt.Recv.List {
						var (
							tp types.Type
							ok bool
						)
						astExpr := rcv.Type
					infLoop:
						for {
							switch rtp := astExpr.(type) {
							case *ast.StarExpr:
								astExpr = rtp.X
							case *ast.Ident:
								tp, ok = p.Types[rtp.Name]
								if !ok {
									continue declLoop
								}
								break infLoop
							default:
								continue declLoop
							}
						}
						switch tpType := tp.(type) {
						case *types.Alias:
							for i, method := range tpType.Methods {
								if method.FuncName == dt.Name.Name {
									funType = &tpType.Methods[i]
								}
							}
						case *types.Struct:
							for i, method := range tpType.Methods {
								if method.FuncName == dt.Name.Name {
									funType = &tpType.Methods[i]
									break
								}
							}
						}
						break
					}
					if funType == nil {
						if r.loadConfig.Verbose {
							log.Printf("method: '%s' not found in the pkg: %s declaration\n", dt.Name.Name, p.Path)
						}
						continue
					}
				} else {
					var ok bool
					funType, ok = p.GetFunction(dt.Name.Name)
					if !ok {
						if r.loadConfig.Verbose {
							log.Printf("function: '%s' not found in the pkg: %s declaration\n", dt.Name.Name, p.Path)
						}
						continue
					}
				}
				funType.Comment = dt.Doc.Text()
			}
		}
	}
}

func (r *rootPackage) scaffoldPackageObjects() {
	s := r.typesPkg.Scope()

	for _, file := range r.pkgPkg.Syntax {
		for _, decl := range file.Decls {
			dt, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range dt.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				tsType := ts.Type

				// This is a workaround only for the invalid type wrapper.
			outLoop:
				for {
					switch tt := tsType.(type) {
					case *ast.Ident, *ast.SelectorExpr:
						r.mappedAliases[ts.Name.Name] = struct{}{}
						break outLoop
					case *ast.StarExpr:
						tsType = tt.X
					default:
						break outLoop
					}
				}
			}
		}
	}

	for _, name := range s.Names() {
		obj := s.Lookup(name)
		switch ot := obj.(type) {
		case *gotypes.TypeName:
			_, isLegacyAlias := r.mappedAliases[name]
			if ot.IsAlias() || isLegacyAlias {
				r.mappedAliases[name] = struct{}{}
				wt := &types.Alias{
					Pkg:       r.refPkg,
					AliasName: name,
				}
				r.setTypeInProgress(name, wt)
				continue
			}

			switch tp := ot.Type().(type) {
			case *gotypes.Named:
				r.scaffoldNamedObject(tp, name)
			default:
				wt := &types.Alias{
					Pkg:       r.refPkg,
					AliasName: name,
				}
				r.setTypeInProgress(name, wt)
			}
		case *gotypes.Func:
			sig, ok := ot.Type().(*gotypes.Signature)
			if !ok {
				continue
			}
			if sig.Recv() != nil {
				// This is a method which should not be extracted as function.
				continue
			}
			fT := &types.Function{
				Pkg:      r.refPkg,
				FuncName: name,
			}
			r.setTypeInProgress(name, fT)
		case *gotypes.Const, *gotypes.Var:
			r.declNames = append(r.declNames, name)
			continue
		default:
			continue
		}
	}
}

func (r *rootPackage) scaffoldNamedObject(named *gotypes.Named, name string) {
	switch t := named.Underlying().(type) {
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
}

func (r *rootPackage) finishNamedType(named *gotypes.Named, t types.Type) bool {
	r.attachNamedGenericInfo(t, named)
	switch ot := t.(type) {
	case *types.Struct:
		return r.finishNamedStructType(ot, named)
	case *types.Interface:
		return r.finishNamedInterfaceType(named, ot)
	case *types.Alias:
		underlying, ok := r.dereferenceType(r.refPkg, named.Underlying())
		if !ok {
			log.Printf("no underlying alias type found: %v", named.Underlying())
			return false
		}
		return r.finishNamedAliasType(named, ot, underlying)
	default:
		log.Printf("invalid named Gentools type %T", t)
		return false
	}
}

func (r *rootPackage) finishNamedAliasType(named *gotypes.Named, alias *types.Alias, underlying types.Type) bool {
	p := r.refPkg
	alias.Type = underlying

	// Map methods.
	for i := 0; i < named.NumMethods(); i++ {
		xm, ok := r.parseMethod(p, named, i, true)
		if !ok {
			return ok
		}
		alias.Methods = append(alias.Methods, xm)
	}
	sort.Slice(alias.Methods, func(i, j int) bool { return alias.Methods[i].FuncName < alias.Methods[j].FuncName })

	return true
}

func (r *rootPackage) parseNamedType(et *gotypes.Named) (types.Type, bool) {
	if et.Obj().Pkg() == nil {
		t, ok := types.GetBuiltInType(et.Obj().Name())
		return t, ok
	}
	p, ok := r.pkgMap.read(et.Obj().Pkg().Path())
	if !ok {
		return nil, ok
	}
	origin, ok := p.GetType(et.Obj().Name())
	if !ok {
		return nil, false
	}
	return r.instantiateNamed(origin, et)
}

func (r *rootPackage) dereferenceType(p *types.Package, tp gotypes.Type) (types.Type, bool) {
	switch et := tp.(type) {
	case *gotypes.TypeParam:
		return r.typeParameter(p, p.Path, et), true
	case *gotypes.Named:
		return r.parseNamedType(et)
	case *gotypes.Struct:
		return r.parseStructType(p, et)
	case *gotypes.Interface:
		return r.parseInterfaceType(p, et)
	case *gotypes.Basic:
		if et.Kind() == gotypes.UnsafePointer {
			return types.UnsafePointer, true
		}
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
		if parsed, ok := r.parsePost119Type(p, et); ok {
			return parsed, true
		}
		log.Printf("type not found for dereferencing: %s, %T\n", et.String(), et)
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
	if !r.parseStructFields(p, named.Underlying().(*gotypes.Struct), t) {
		return false
	}

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
	return true
}

type astMethoder interface {
	Method(int) *gotypes.Func
}

func (r *rootPackage) parseMethod(p *types.Package, named astMethoder, i int, needReceiver bool) (types.Function, bool) {
	m := named.Method(i)

	s, ok := m.Type().(*gotypes.Signature)
	if !ok {
		log.Printf("method type is not a signature: %+v\n", m.Type())
		return types.Function{}, false
	}
	ft := types.Function{FuncName: m.Name(), Pkg: p}
	if ok = r.parseSignatureType(p, s, &ft, needReceiver); !ok {
		return types.Function{}, false
	}
	r.attachSignatureGenericInfo(&ft, s)
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
	r.attachSignatureGenericInfo(ft, st)
	return true
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
