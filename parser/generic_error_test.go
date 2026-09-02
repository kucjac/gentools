package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kucjac/gentools/types"
)

func TestLoadPackagesReturnsCompilerDiagnosticsForInvalidGenericSource(t *testing.T) {
	path := filepath.Join(t.TempDir(), "broken.go")
	if err := os.WriteFile(path, []byte("package broken\nfunc Broken[T any]( {\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadPackages(LoadConfig{PkgNames: []string{"file=" + path}})
	if err == nil {
		t.Fatal("LoadPackages() succeeded for invalid generic source")
	}
	if !strings.Contains(err.Error(), "expected") {
		t.Fatalf("LoadPackages() error = %q, want compiler diagnostic", err)
	}
}

func TestUpdatePackagesLeavesMapUnchangedOnFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "broken.go")
	if err := os.WriteFile(path, []byte("package broken\nfunc Broken[T any]( {\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	existing := types.NewPackage("example.test/existing", "existing")
	packages := types.PackageMap{existing.Path: existing}
	if err := UpdatePackages(packages, LoadConfig{PkgNames: []string{"file=" + path}}); err == nil {
		t.Fatal("UpdatePackages() succeeded for invalid generic source")
	}
	if len(packages) != 1 || packages[existing.Path] != existing {
		t.Fatalf("UpdatePackages() mutated package map on failure: %#v", packages)
	}
}
