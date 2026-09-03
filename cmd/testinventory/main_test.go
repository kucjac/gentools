package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAssessRejectsMissingEvidence(t *testing.T) {
	dir := t.TempDir()
	inv := filepath.Join(dir, "inventory.json")
	if err := os.WriteFile(inv, []byte(`{"entries":[{"id":"x","symbol":"parser.*","owner":"parser","risk":"low","test_level":"unit","status":"verified","evidence":"missing"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := assess(inv, filepath.Join(dir, "report.md")); err == nil {
		t.Fatal("expected missing evidence failure")
	}
}
