package types

import "testing"

func TestStructTagSplitLookupEscapingAndAbsentKey(t *testing.T) {
	tag := StructTag(`json:"name,omitempty" note:"say \"hello\"" empty:""`)
	tuples := tag.Split()
	if len(tuples) != 3 || tuples[1].Value != `say "hello"` {
		t.Fatalf("unexpected split: %#v", tuples)
	}
	if got, ok := tag.Lookup("empty"); !ok || got != "" {
		t.Fatalf("empty value lookup = %q, %v", got, ok)
	}
	if _, ok := tag.Lookup("missing"); ok || tag.Get("missing") != "" {
		t.Fatal("absent tag key was found")
	}
	if tuples.Join() != tag {
		t.Fatalf("split/join did not round trip: %q", tuples.Join())
	}
}
