package parser

import "testing"

func TestPackageNameOfDirAndInvalidConfiguration(t *testing.T) {
	if _, err := PackageNameOfDir(t.TempDir()); err == nil {
		t.Fatal("expected no-source error")
	}
	if _, err := PackageNameOfDir("missing-directory"); err == nil {
		t.Fatal("expected missing-directory error")
	}
}
