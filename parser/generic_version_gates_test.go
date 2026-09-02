package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersionGatedFixturesDeclareTheirMinimumToolchain(t *testing.T) {
	fixtures := map[string]string{
		"genericalias/models_go124.go":   "//go:build go1.24",
		"genericmethods/models_go127.go": "//go:build go1.27",
	}
	for relativePath, buildConstraint := range fixtures {
		contents, err := os.ReadFile(filepath.Join("testcases", relativePath))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(contents), buildConstraint) {
			t.Errorf("%s does not declare %q", relativePath, buildConstraint)
		}
	}
}
