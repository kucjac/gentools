package astreflect

import (
	"testing"
)

func TestZero(t *testing.T) {
	pkgs, err := LoadPackages(&LoadConfig{
		Paths:      []string{"."},
		PkgNames:   nil,
		BuildFlags: nil,
		Verbose:    false,
	})
	if err != nil {
		t.Errorf("Loading packages failed: %v", err)
		return
	}
	_, ok := pkgs.TypeOf("astreflect.Type", nil)
	if !ok {
		t.Error("TypeOf find astreflect.Type failed")
	}
	_, ok = pkgs.TypeOf("astreflect.StructType", nil)
	if !ok {
		t.Error("TypeOf find 'astreflect.StructType' failed")
	}

	thisPkg, ok := pkgs.GetByPath("github.com/kucjac/astreflect")
	if !ok {
		t.Error("This package not found by path")
		return
	}

	for _, expected := range []struct {
		PkgContext *Package
		Name       string
		Type       Type
	}{
		{
			Name:       "Type",
			Type:       thisPkg.MustGetType("Type"),
			PkgContext: thisPkg,
		},
		{
			Name: "astreflect.Type",
			Type: thisPkg.MustGetType("Type"),
		},
		{
			Name: "[]*astreflect.StructType",
			Type: func() Type {
				return ArrayType{ArrayKind: Slice, Type: PointerType{PointedType: thisPkg.MustGetType("StructType")}}
			}(),
		},
		{
			Name: "[]*StructType",
			Type: func() Type {
				return ArrayType{ArrayKind: Slice, Type: PointerType{PointedType: thisPkg.MustGetType("StructType")}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[]*StructType{}",
			Type: func() Type {
				return ArrayType{ArrayKind: Slice, Type: PointerType{PointedType: thisPkg.MustGetType("StructType")}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[3][]*InterfaceType",
			Type: func() Type {
				return ArrayType{ArrayKind: Array, ArraySize: 3, Type: ArrayType{ArrayKind: Slice, Type: PointerType{PointedType: thisPkg.MustGetType("InterfaceType")}}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[3][]*astreflect.InterfaceType",
			Type: func() Type {
				return ArrayType{ArrayKind: Array, ArraySize: 3, Type: ArrayType{ArrayKind: Slice, Type: PointerType{PointedType: thisPkg.MustGetType("InterfaceType")}}}
			}(),
		},
		{
			Name: "[3][]<- chan InterfaceType",
			Type: func() Type {
				return ArrayType{ArrayKind: Array, ArraySize: 3, Type: ArrayType{ArrayKind: Slice, Type: ChanType{Dir: RecvOnly, Type: thisPkg.MustGetType("InterfaceType")}}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[3][]chan <- *astreflect.InterfaceType",
			Type: func() Type {
				return ArrayType{
					ArrayKind: Array, ArraySize: 3, Type: ArrayType{
						ArrayKind: Slice, Type: ChanType{
							Dir: SendOnly, Type: PointerType{PointedType: thisPkg.MustGetType("InterfaceType")}},
					},
				}
			}(),
		},
		{
			Name:       "map[string][]*StructType",
			Type:       MapType{Key: BuiltInType{StdKind: String}, Value: ArrayType{ArrayKind: Slice, Type: PointerType{PointedType: thisPkg.MustGetType("StructType")}}},
			PkgContext: thisPkg,
		},
		{
			Name: "[]astreflect.Kind",
			Type: ArrayType{ArrayKind: Slice, Type: thisPkg.MustGetType("Kind")},
		},
		{
			Name: "[]byte",
			Type: ArrayType{ArrayKind: Slice, Type: MustGetBuiltInType("byte")},
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
