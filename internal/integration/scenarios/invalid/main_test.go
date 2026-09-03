package invalid

import (
	"github.com/kucjac/gentools/parser"
	"strings"
	"testing"
)

func TestInvalidConsumerContract(t *testing.T) {
	_, err := parser.LoadPackages(parser.LoadConfig{Paths: []string{"testdata"}})
	if err == nil {
		t.Fatal("expected invalid source error")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Fatalf("error is not actionable: %v", err)
	}
}
