package types

import "testing"

func TestValueTypeContracts(t *testing.T) {
	array := ArrayOf(String, 2)
	if array.Name(true, "") != "[2]string" || array.Zero(true, "") != "[2]string{}" || !array.Elem().Equal(String) {
		t.Fatalf("unexpected array contract: %s / %s", array.Name(true, ""), array.Zero(true, ""))
	}
	if !SliceOf(String).Equal(SliceOf(String)) || SliceOf(String).Zero(true, "") != "nil" {
		t.Fatal("slice contract changed")
	}
	if !MapOf(String, Int).Equal(MapOf(String, Int)) || MapOf(String, Int).Zero(true, "") != "nil" {
		t.Fatal("map contract changed")
	}
	pointer := PointerTo(String)
	if pointer.Name(true, "") != "*string" || pointer.Zero(true, "") != "nil" || !pointer.Elem().Equal(String) {
		t.Fatal("pointer contract changed")
	}
	channel := ChanOf(RecvOnly, String)
	if channel.Name(true, "") != "chan<- string" || !channel.Elem().Equal(String) || !channel.Equal(ChanOf(RecvOnly, String)) {
		t.Fatal("channel contract changed")
	}
}
