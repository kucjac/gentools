package generated

import (
	"testing"

	_ "github.com/kucjac/gentools/internal/integration/protobuf"
	"github.com/kucjac/gentools/parser"
	"github.com/kucjac/gentools/types"
)

func TestGeneratedProtobufConsumerContract(t *testing.T) {
	pkgs, err := parser.LoadPackages(parser.LoadConfig{Paths: []string{"../../protobuf"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, pkg := range pkgs {
		if value, ok := pkg.GetType("TestingMessage"); ok {
			message, ok := value.(*types.Struct)
			if !ok {
				t.Fatalf("TestingMessage = %T", value)
			}
			found := map[string]bool{}
			for _, field := range message.Fields {
				found[field.Name] = true
			}
			for _, name := range []string{"Any", "Duration", "Timestamp", "File"} {
				if !found[name] {
					t.Fatalf("missing generated field %s", name)
				}
			}
			return
		}
	}
	t.Fatal("TestingMessage was not found")
}
