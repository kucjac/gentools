package astreflect

import (
	"testing"
)

func TestParsePackages(t *testing.T) {
	// On test purpose try to parse this package.
	pkgs, err := LoadPackages(&LoadConfig{Paths: []string{"."}})
	if err != nil {
		t.Errorf("Parsing packages failed: %v", err)
		return
	}

	// Get package reflection by it's identifier.
	thisPkg, ok := pkgs.GetByIdentifier("astreflect")
	if !ok {
		t.Error("Package 'astreflect' not found by identifier")
		return
	}

	// Get the 'Type' interface.
	typeInterface, ok := thisPkg.GetInterfaceType("Type")
	if !ok {
		t.Error("Interface 'Type' not found'")
		return
	}

	// Get the 'StructType' struct type.
	structType, ok := thisPkg.GetStructType("StructType")
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
	pointer := &PointerType{PointedType: structType}
	if ok := Implements(pointer, typeInterface); !ok {
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
			expectedKind     Kind
			expectedElemKind Kind
		)

		switch i {
		case 0:
			expectedName = "PackagePath"
			expectedType = "PkgPath"
			expectedKind = Wrapper
			expectedElemKind = String
		case 1:
			expectedName = "Comment"
			expectedType = "string"
			expectedKind = String
		case 2:
			expectedName = "TypeName"
			expectedType = "string"
			expectedKind = String
		case 3:
			expectedName = "Fields"
			expectedType = "[]StructField"
			expectedKind = Slice
			expectedElemKind = Struct
		case 4:
			expectedName = "Methods"
			expectedType = "[]FunctionType"
			expectedKind = Slice
			expectedElemKind = Struct
		}
		if sField.Name != expectedName {
			t.Errorf("Expected field name mismatch. Expected: %s, is %s", expectedName, sField.Name)
		}
		if sField.Type.Name(false) != expectedType {
			t.Errorf("Expected field type mismatch. Expected: %s, is %s", expectedType, sField.Type.Name(false))
		}
		if sField.Type.Kind() != expectedKind {
			t.Errorf("Expected field kind mismatch. Expected: %s is %s", expectedKind, sField.Type.Kind())
		}
		if expectedElemKind != Invalid {
			if sField.Type.Elem().Kind() != expectedElemKind {
				t.Errorf("Expected elem Kind mismatch. Expected: %s is %s", expectedElemKind, sField.Type.Elem().Kind())
			}
		}
	}

	tt, ok := thisPkg.GetStructType("testingType")
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
	if embeddedField.Type.Kind() != Struct {
		t.Errorf("'embeddedType' field kind should be 'Struct' is: %s", embeddedField.Type.Kind())
		return
	}
	et, ok := embeddedField.Type.(*StructType)
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
	if etField.Type.Name(true) != "*testing.T" {
		t.Errorf("'embeddedType' field type should be '*testing.T' but is: %v", etField.Type.Name(true))
	}
}
