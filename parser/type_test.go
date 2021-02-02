package parser

import (
	"testing"

	"github.com/kucjac/gentools/types"
)

func TestZero(t *testing.T) {
	pkgs, err := LoadPackages(LoadConfig{
		Paths:      []string{"."},
		PkgNames:   nil,
		BuildFlags: nil,
		Verbose:    false,
	})
	if err != nil {
		t.Errorf("Loading packages failed: %v", err)
		return
	}

	thisPkg, ok := pkgs.PackageByPath("github.com/kucjac/gentools/types")
	if !ok {
		t.Error("This package not found by path")
		return
	}
	thisPkg.SetIdentifier("gentypes")

	_, ok := pkgs.TypeOf("types.Type", nil)
	if !ok {
		t.Error("TypeOf find types.Type failed")
	}
	_, ok = pkgs.TypeOf("types.StructType", nil)
	if !ok {
		t.Error("TypeOf find 'types.StructType' failed")
	}

	for _, expected := range []struct {
		PkgContext *types.Package
		Name       string
		Type       types.Type
	}{
		{
			Name:       "Type",
			Type:       thisPkg.MustGetType("Type"),
			PkgContext: thisPkg,
		},
		{
			Name: "gentypes.Type",
			Type: thisPkg.MustGetType("Type"),
		},
		{
			Name: "[]*gentypes.StructType",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("StructType")}}
			}(),
		},
		{
			Name: "[]*StructType",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("StructType")}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[]*StructType{}",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("StructType")}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[3][]*InterfaceType",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindArray, ArraySize: 3, Type: &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("InterfaceType")}}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[3][]*gentypes.InterfaceType",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindArray, ArraySize: 3, Type: &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("InterfaceType")}}}
			}(),
		},
		{
			Name: "[3][]<- chan InterfaceType",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindArray, ArraySize: 3, Type: &types.Array{ArrayKind: types.KindSlice, Type: &types.Chan{Dir: types.RecvOnly, Type: thisPkg.MustGetType("InterfaceType")}}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[3][]chan <- *gentypes.InterfaceType",
			Type: func() types.Type {
				return &types.Array{
					ArrayKind: types.KindArray, ArraySize: 3, Type: &types.Array{
						ArrayKind: types.KindSlice, Type: &types.Chan{
							Dir: types.SendOnly, Type: &types.Pointer{PointedType: thisPkg.MustGetType("InterfaceType")}},
					},
				}
			}(),
		},
		{
			Name:       "map[string][]*StructType",
			Type:       &types.Map{Key: &types.BuiltInType{BuiltInKind: types.KindString}, Value: &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("StructType")}}},
			PkgContext: thisPkg,
		},
		{
			Name: "[]gentypes.Kind",
			Type: &types.Array{ArrayKind: types.KindSlice, Type: thisPkg.MustGetType("Kind")},
		},
		{
			Name: "[]byte",
			Type: &types.Array{ArrayKind: types.KindSlice, Type: types.MustGetBuiltInType("byte")},
		},
	} {
		tp, ok := pkgs.TypeOf(expected.Name, expected.PkgContext)
		if !ok {
			t.Errorf("Package: TypeOf('%s') failed", expected.Name)
			continue
		}
		if !tp.Equal(expected.Type) {
			t.Errorf("TypeOf resulting type: '%s' is not equal to expected: '%s'", tp, expected.Type)
		}
	}
}
