package inspect

import (
	"github.com/kucjac/gentools/parser"
	"github.com/kucjac/gentools/types"
	"testing"
)

func TestInspectConsumerContract(t *testing.T) {
	pkgs, err := parser.LoadPackages(parser.LoadConfig{Paths: []string{"testdata"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, pkg := range pkgs {
		if widget, ok := pkg.GetType("Widget"); ok {
			if _, ok := widget.(*types.Struct); !ok {
				t.Fatalf("Widget = %T", widget)
			}
			return
		}
	}
	t.Fatal("Widget was not found")
}
