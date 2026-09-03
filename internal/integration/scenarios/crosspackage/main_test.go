package crosspackage

import (
	"github.com/kucjac/gentools/parser"
	"github.com/kucjac/gentools/types"
	"testing"
)

func TestCrossPackageGenericConsumerContract(t *testing.T) {
	pkgs, err := parser.LoadPackages(parser.LoadConfig{Paths: []string{"testdata/producer", "testdata/consumer"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, pkg := range pkgs {
		if holder, ok := pkg.GetType("Holder"); ok {
			s, ok := holder.(*types.Struct)
			if !ok || len(s.Fields) != 1 {
				t.Fatalf("Holder = %#v", holder)
			}
			if _, ok := s.Fields[0].Type.(*types.Instantiation); !ok {
				t.Fatalf("field = %T", s.Fields[0].Type)
			}
			return
		}
	}
	t.Fatal("Holder was not found")
}
