package parser

import (
	"testing"

	"github.com/kucjac/gentools/types"
)

func TestParsePackages(t *testing.T) {
	// On test purpose try to parse this package.
	pkgs, err := LoadPackages(LoadConfig{Paths: []string{"."}, Verbose: true})
	if err != nil {
		t.Errorf("Parsing packages failed: %v", err)
		return
	}

	// Get package reflection by it's identifier.
	typesPkg, ok := pkgs.PackageByPath("github.com/kucjac/gentools/types")
	if !ok {
		t.Error("Package 'types' not found by path")
		return
	}

	// Get the 'Type' interface.
	typeInterface, ok := typesPkg.GetInterfaceType("Type")
	if !ok {
		t.Error("Interface 'Type' not found'")
		return
	}

	// Get the 'StructType' struct type.
	structType, ok := typesPkg.GetStruct("StructType")
	if !ok {
		t.Error("Struct 'StructType' not found")
		return
	}

	// Check if 'StructType' implements 'Type' interface.
	// In fact it shouldn't implement it - only *StructType implements it.
	if ok = structType.Implements(typeInterface, false); ok {
		t.Error("'StructType' implements 'Type' interface but it shouldn't.")
		return
	}

	if ok = structType.Implements(typeInterface, true); !ok {
		t.Error("'*StructType' doesn't implement 'Type' interface, but it should.")
		return
	}

	// The pointer to the 'StructType' should in fact implement the 'Type' interface. Let's check it.
	pointer := &types.Pointer{PointedType: structType}
	if ok := types.Implements(pointer, typeInterface); !ok {
		t.Error("'*StructType' doesn't implement 'Type' interface")
		return
	}

	// The API allows to check the fields for given struct type.
	if len(structType.Fields) != 5 {
		t.Errorf("'StructType' should have 5 fields but have: %d", len(structType.Fields))
		return
	}
	for i, sField := range structType.Fields {
		var (
			expectedName     string
			expectedType     string
			expectedKind     types.Kind
			expectedElemKind types.Kind
		)

		switch i {
		case 0:
			expectedName = "Pkg"
			expectedType = "*Package"
			expectedKind = types.KindPtr
			expectedElemKind = types.KindStruct
		case 1:
			expectedName = "Comment"
			expectedType = "string"
			expectedKind = types.KindString
		case 2:
			expectedName = "TypeName"
			expectedType = "string"
			expectedKind = types.KindString
		case 3:
			expectedName = "Fields"
			expectedType = "[]StructField"
			expectedKind = types.KindSlice
			expectedElemKind = types.KindStruct
		case 4:
			expectedName = "Methods"
			expectedType = "[]FunctionType"
			expectedKind = types.KindSlice
			expectedElemKind = types.KindStruct
		}
		if sField.Name != expectedName {
			t.Errorf("Expected field name mismatch. Expected: %s, is %s", expectedName, sField.Name)
		}
		if sField.Type.Name(false, "") != expectedType {
			t.Errorf("Expected field type mismatch. Expected: %s, is %s", expectedType, sField.Type.Name(false, ""))
		}
		if sField.Type.Kind() != expectedKind {
			t.Errorf("Expected field kind mismatch. Expected: %s is %s", expectedKind, sField.Type.Kind())
		}
		if expectedElemKind != types.Invalid {
			if sField.Type.Elem().Kind() != expectedElemKind {
				t.Errorf("Expected elem Kind mismatch. Expected: %s is %s", expectedElemKind, sField.Type.Elem().Kind())
			}
		}
	}

	thisPkg, ok := pkgs.PackageByPath("github.com/kucjac/gentools/parser")
	if !ok {
		t.Error("Package parser not found")
		return
	}

	tt, ok := thisPkg.GetStruct("testingType")
	if !ok {
		t.Error("TestingType not found")
		return
	}

	// Testing type should have two fields:
	// 1. embeddedType
	// 2. Integer int
	if len(tt.Fields) != 2 {
		t.Errorf("testingType should have 2 fields")
		return
	}
	embeddedField := tt.Fields[0]
	if !embeddedField.Embedded {
		t.Errorf("'embeddedType' field should be embedded")
		return
	}
	if embeddedField.Type.Kind() != types.KindStruct {
		t.Errorf("'embeddedType' field kind should be 'Struct' is: %s", embeddedField.Type.Kind())
		return
	}
	et, ok := embeddedField.Type.(*types.Struct)
	if !ok {
		t.Errorf("'embeddedType' field should be of 'StructType' field. Is: %T", embeddedField.Type)
		return
	}

	// It should have one imported field with a tag.
	if len(et.Fields) != 1 {
		t.Errorf("'embeddedType' should have exactly one field, but have: %d", len(et.Fields))
		return
	}

	etField := et.Fields[0]
	if etField.Name != "Imported" {
		t.Errorf("'embeddedType' field name should be 'Imported' but is: %s", etField.Name)
		return
	}
	if etField.Type.Name(true, "") != "*testing.T" {
		t.Errorf("'embeddedType' field type should be '*testing.T' but is: %v", etField.Type.Name(true, ""))
	}
	if etField.Type.Name(true, "testing") != "*T" {
		t.Errorf("'embeddedType' field type with 'testing' package context should be '*T' but is: %v", etField.Type.Name(true, ""))
	}
}
