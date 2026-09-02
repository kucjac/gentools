package parser

import (
	gotypes "go/types"

	"github.com/kucjac/gentools/types"
)

func (r *rootPackage) typeParameter(p *types.Package, owner string, parameter *gotypes.TypeParam) *types.TypeParameter {
	if existing, ok := r.typeParameters[parameter]; ok {
		if owner != "" && existing.Owner != owner {
			existing.Owner = owner
		}
		return existing
	}
	result := &types.TypeParameter{
		Identifier: parameter.Obj().Name(),
		Owner:      owner,
		Position:   parameter.Index(),
	}
	r.typeParameters[parameter] = result
	if constraint, ok := r.dereferenceType(p, parameter.Constraint()); ok {
		result.Constraint = &types.Constraint{Type: constraint, Expression: parameter.Constraint().Underlying().String()}
	}
	return result
}

func (r *rootPackage) parameters(p *types.Package, owner string, list *gotypes.TypeParamList) []types.TypeParameter {
	if list == nil || list.Len() == 0 {
		return nil
	}
	result := make([]types.TypeParameter, list.Len())
	for index := 0; index < list.Len(); index++ {
		parameter := r.typeParameter(p, owner, list.At(index))
		result[index] = *parameter
	}
	return result
}

func (r *rootPackage) attachNamedGenericInfo(owner types.Type, named *gotypes.Named) {
	params := r.parameters(r.refPkg, named.Obj().Pkg().Path()+"/"+named.Obj().Name(), named.TypeParams())
	if len(params) != 0 {
		types.SetGenericInfo(owner, types.GenericInfo{TypeParameters: params})
	}
}

func (r *rootPackage) attachSignatureGenericInfo(owner types.Type, signature *gotypes.Signature) {
	ownerName := ""
	if function, ok := owner.(*types.Function); ok {
		ownerName = function.FullName()
	}
	info := types.GenericInfo{
		TypeParameters:         r.parameters(r.refPkg, ownerName, signature.TypeParams()),
		ReceiverTypeParameters: r.parameters(r.refPkg, ownerName, signature.RecvTypeParams()),
	}
	if len(info.TypeParameters) != 0 || len(info.ReceiverTypeParameters) != 0 {
		types.SetGenericInfo(owner, info)
	}
}

func (r *rootPackage) instantiateNamed(origin types.Type, named *gotypes.Named) (types.Type, bool) {
	arguments := named.TypeArgs()
	if arguments == nil || arguments.Len() == 0 {
		return origin, true
	}
	result := &types.Instantiation{Origin: origin, Arguments: make([]types.Type, arguments.Len())}
	for index := 0; index < arguments.Len(); index++ {
		argument, ok := r.dereferenceType(r.refPkg, arguments.At(index))
		if !ok {
			return nil, false
		}
		result.Arguments[index] = argument
	}
	return result, true
}
