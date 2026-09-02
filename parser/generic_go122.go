package parser

import (
	gotypes "go/types"

	"github.com/kucjac/gentools/types"
)

func (r *rootPackage) parsePost119Type(p *types.Package, value gotypes.Type) (types.Type, bool) {
	alias, ok := value.(*gotypes.Alias)
	if !ok {
		return nil, false
	}
	origin, ok := r.aliasOrigin(alias)
	if !ok {
		origin, ok = r.dereferenceType(p, alias.Rhs())
	}
	if !ok {
		return nil, false
	}
	arguments := alias.TypeArgs()
	if arguments == nil || arguments.Len() == 0 {
		return origin, true
	}
	result := &types.Instantiation{Origin: origin, Arguments: make([]types.Type, arguments.Len())}
	for index := 0; index < arguments.Len(); index++ {
		argument, ok := r.dereferenceType(p, arguments.At(index))
		if !ok {
			return nil, false
		}
		result.Arguments[index] = argument
	}
	return result, true
}

func (r *rootPackage) finishAliasType(p *types.Package, object *gotypes.TypeName, target *types.Alias) bool {
	alias, ok := object.Type().(*gotypes.Alias)
	if !ok {
		underlying, ok := r.dereferenceType(p, object.Type())
		if !ok {
			return false
		}
		target.Type = underlying
		return true
	}
	underlying, ok := r.dereferenceType(p, alias.Rhs())
	if !ok {
		return false
	}
	target.Type = underlying
	params := r.parameters(p, object.Pkg().Path()+"/"+object.Name(), alias.TypeParams())
	if len(params) != 0 {
		types.SetGenericInfo(target, types.GenericInfo{TypeParameters: params})
	}
	return true
}

func (r *rootPackage) aliasOrigin(alias *gotypes.Alias) (types.Type, bool) {
	object := alias.Origin().Obj()
	if object.Pkg() == nil {
		return nil, false
	}
	pkg, ok := r.pkgMap.read(object.Pkg().Path())
	if !ok {
		return nil, false
	}
	return pkg.GetType(object.Name())
}
