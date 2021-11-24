package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"log"
	"strconv"
	"strings"

	"github.com/kucjac/gentools/types"
	"golang.org/x/tools/go/packages"
)

var errIdentNotFound = errors.New("ident not found")

func (r *rootPackage) extractAliasExpr(f *ast.File, expr ast.Expr) (types.Type, error) {
	switch x := expr.(type) {
	case *ast.StarExpr:
		of, err := r.extractAliasExpr(f, x.X)
		if err != nil {
			return nil, err
		}
		return types.PointerTo(of), nil
	case *ast.SelectorExpr:
		i, ok := x.X.(*ast.Ident)
		if !ok {
			panic("selector expression is not an ident")
		}
		for _, im := range f.Imports {
			pkgPath := strings.Trim(im.Path.Value, "\"")
			ident := pkgPath
			if im.Name == nil {
				if li := strings.LastIndex(ident, "/"); li != -1 {
					ident = ident[li+1:]
				}
			} else {
				ident = im.Name.Name
			}
			if ident != i.Name {
				continue
			}

			if pkgPath == "unsafe" && x.Sel.Name == "Pointer" {
				return types.UnsafePointer, nil
			}

			pkg, ok := r.pkgPkg.Imports[pkgPath]
			if !ok {
				return nil, fmt.Errorf("imported package not found %s", im.Path.Value)
			}

			root, ok := r.rootPackages[pkg.Types]
			if !ok {
				return nil, fmt.Errorf("root package: %s not found", pkg.PkgPath)
			}
			// Find matching file for the selector definition.
			file, ok := r.findFileMatchingIdent(pkg, x.Sel)
			if !ok {
				return nil, fmt.Errorf("file matching given ident not found: %s", x.Sel.Name)
			}

			return root.extractAliasExpr(file, x.Sel)
		}

		// Need to search for the file specific import alias.
		return nil, fmt.Errorf("no matching imported package to given selector: '%s'", x.Sel.Name)
	case *ast.Ident:
		if tp, ok := types.GetBuiltInType(x.Name); ok {
			return tp, nil
		}
		if r.refPkg.Path == "unsafe" && x.Name == "Pointer" {
			return types.UnsafePointer, nil
		}
		tp, ok := r.refPkg.GetType(x.Name)
		if !ok {
			if r.loadConfig.Verbose {
				log.Printf("ident: '%s' not found in the package: '%s'", x.Name, r.refPkg.Path)
			}
			return nil, errIdentNotFound
		}
		return tp, nil
	case *ast.MapType:
		k, err := r.extractAliasExpr(f, x.Key)
		if err != nil {
			return nil, err
		}
		v, err := r.extractAliasExpr(f, x.Value)
		if err != nil {
			return nil, err
		}
		return types.MapOf(k, v), nil
	case *ast.ChanType:
		tp, err := r.extractAliasExpr(f, x.Value)
		if err != nil {
			return nil, err
		}
		dir := types.SendRecv
		if x.Dir&ast.RECV == 0 {
			dir = types.SendOnly
		}
		if x.Dir&ast.SEND == 0 {
			dir = types.RecvOnly
		}
		return types.ChanOf(dir, tp), nil
	case *ast.ArrayType:
		tp, err := r.extractAliasExpr(f, x.Elt)
		if err != nil {
			return nil, err
		}
		var size int
		switch l := x.Len.(type) {
		case *ast.Ellipsis:
		case *ast.Ident:
			for _, decl := range r.refPkg.Declarations {
				if l.Name != decl.Name {
					continue
				}
				if !decl.Constant {
					return nil, fmt.Errorf("non constant alias array size: %v", l.Name)
				}
				constVal := decl.ConstValue()
				if ui, ok := constVal.(uint); ok {
					size = int(ui)
				} else if i, ok := constVal.(int); ok {
					size = i
				}
				break
			}
		case *ast.BasicLit:
			i, err := strconv.Atoi(l.Value)
			if err != nil {
				return nil, fmt.Errorf("array type basic lit value is not an integer: %s -  %w", l.Value, err)
			}
			size = i
		}
		if size == 0 {
			return types.SliceOf(tp), nil
		}
		return types.ArrayOf(tp, size), nil
	default:
		return nil, nil
	}
}

func (r *rootPackage) extractStructExpr(expr ast.Expr) (*ast.StructType, bool) {
	switch x := expr.(type) {
	case *ast.Ident:
		for _, f := range r.pkgPkg.Syntax {
			for _, decl := range f.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				for _, spec := range gen.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					if ts.Name.Name != x.Name {
						continue
					}
					return r.extractStructExpr(ts.Type)
				}
			}
		}
		return nil, false
	case *ast.SelectorExpr:
		i, ok := x.X.(*ast.Ident)
		if !ok {
			panic("selector expression is not an ident")
		}
		for _, im := range r.typesPkg.Imports() {
			if im.Name() != i.Name {
				continue
			}

			root, ok := r.rootPackages[im]
			if !ok {
				if r.loadConfig.Verbose {
					log.Printf("root package: %s not found", im.Path())
				}
				return nil, false
			}
			return root.extractStructExpr(x.Sel)
		}
	case *ast.StarExpr:
		return r.extractStructExpr(x.X)
	case *ast.StructType:
		return x, true
	}
	return nil, false
}

func (r *rootPackage) extractInterfaceExpr(expr ast.Expr) (*ast.InterfaceType, bool) {
	switch x := expr.(type) {
	case *ast.Ident:
		for _, f := range r.pkgPkg.Syntax {
			for _, decl := range f.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				for _, spec := range gen.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					if ts.Name.Name != x.Name {
						continue
					}
					return r.extractInterfaceExpr(ts.Type)
				}
			}
		}
		return nil, false
	case *ast.SelectorExpr:
		i, ok := x.X.(*ast.Ident)
		if !ok {
			panic("selector expression is not an ident")
		}
		for _, im := range r.typesPkg.Imports() {
			if im.Name() != i.Name {
				continue
			}

			root, ok := r.rootPackages[im]
			if !ok {
				if r.loadConfig.Verbose {
					log.Printf("root package: %s not found", im.Path())
				}
				return nil, false
			}
			return root.extractInterfaceExpr(x.Sel)
		}
	case *ast.StarExpr:
		return r.extractInterfaceExpr(x.X)
	case *ast.InterfaceType:
		return x, true
	}
	return nil, false
}

func (r *rootPackage) findFileMatchingIdent(pkg *packages.Package, ident *ast.Ident) (*ast.File, bool) {
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gd.Specs {
				switch t := spec.(type) {
				case *ast.TypeSpec:
					if t.Name.Name == ident.Name {
						return file, true
					}
				case *ast.ValueSpec:
					for _, name := range t.Names {
						if name.Name == ident.Name {
							return file, true
						}
					}
				}
			}
		}
	}
	return nil, false
}
