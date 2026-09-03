package types

import "testing"

func TestInterfaceCompatibilityAndNegativeMethodContracts(t *testing.T) {
	pkg := NewPackage("example.test/contracts", "contracts")
	required := &Interface{Pkg: pkg, InterfaceName: "Reader", Methods: []Function{{FuncName: "Read", Pkg: pkg, In: []FuncParam{{Type: String}}, Out: []FuncParam{{Type: Int}}}}}
	valueMethod := Function{FuncName: "Read", Pkg: pkg, In: []FuncParam{{Type: String}}, Out: []FuncParam{{Type: Int}}}
	value := &Struct{Pkg: pkg, TypeName: "Value", Methods: []Function{valueMethod}}
	if !Implements(value, required) {
		t.Fatal("matching value receiver did not implement interface")
	}
	pointerMethod := valueMethod
	pointerMethod.Receiver = &Receiver{Name: "v", Type: PointerTo(value)}
	pointer := &Struct{Pkg: pkg, TypeName: "Pointer", Methods: []Function{pointerMethod}}
	if Implements(pointer, required) || !Implements(PointerTo(pointer), required) {
		t.Fatal("pointer receiver compatibility changed")
	}
	wrong := &Struct{Pkg: pkg, TypeName: "Wrong", Methods: []Function{{FuncName: "Read", Pkg: pkg, Variadic: true, In: []FuncParam{{Type: String}}, Out: []FuncParam{{Type: Int}}}}}
	if Implements(wrong, required) {
		t.Fatal("variadic mismatch implemented interface")
	}
}
