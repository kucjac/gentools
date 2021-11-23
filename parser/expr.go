package parser

import (
	"fmt"
	"go/ast"
	"strconv"

	"github.com/kucjac/gentools/types"
)

func (r *rootPackage) extractAliasExpr(expr ast.Expr) (types.Type, error) {
	switch x := expr.(type) {
	case *ast.StarExpr:
		of, err := r.extractAliasExpr(x.X)
		if err != nil {
			return nil, err
		}
		return types.PointerTo(of), nil
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
				return nil, fmt.Errorf("root package: %s not found", im.Path())
			}

			return root.extractAliasExpr(x.Sel)
		}
		return nil, fmt.Errorf("no matching imported package to given selector: '%s'", x.Sel.Name)
	case *ast.Ident:
		if tp, ok := types.GetBuiltInType(x.Name); ok {
			return tp, nil
		}
		tp, ok := r.refPkg.GetType(x.Name)
		if !ok {
			return nil, fmt.Errorf("ident: '%s' not found in the package: '%s'", x.Name, r.refPkg.Path)
		}
		return tp, nil
	case *ast.MapType:
		k, err := r.extractAliasExpr(x.Key)
		if err != nil {
			return nil, err
		}
		v, err := r.extractAliasExpr(x.Value)
		if err != nil {
			return nil, err
		}
		return types.MapOf(k, v), nil
	case *ast.ChanType:
		tp, err := r.extractAliasExpr(x.Value)
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
		tp, err := r.extractAliasExpr(x.Elt)
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
		if size != 0 {
			return types.SliceOf(tp), nil
		}
		return types.ArrayOf(tp, size), nil
	default:
		return nil, nil
	}
}
