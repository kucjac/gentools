package types

import "testing"

func TestBuiltinLookupConversionEqualityAndPanic(t *testing.T) {
	if got, ok := GetBuiltInType("string"); !ok || !got.Equal(String) {
		t.Fatal("string lookup did not return String")
	}
	if !BuiltInOf(1).Equal(Int) || !BuiltInOf("").Equal(String) || !BuiltInOf(true).Equal(Bool) {
		t.Fatal("Go value conversion did not select the matching builtin")
	}
	if String.Equal(Int) {
		t.Fatal("distinct builtin types compare equal")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("MustGetBuiltInType did not panic for an unknown name")
		}
	}()
	MustGetBuiltInType("not-a-builtin")
}
