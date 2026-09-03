package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateMatchesGoldenAndReusesModel(t *testing.T) {
	output := filepath.Join(t.TempDir(), "openapi.json")
	if err := generate("testdata", output); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/openapi.golden.json")
	if err != nil {
		t.Fatal(err)
	}
	if normalizeNewlines(strings.TrimSpace(string(got))) != normalizeNewlines(strings.TrimSpace(string(want))) {
		t.Fatalf("generated contract differs\nwant:\n%s\ngot:\n%s", want, got)
	}
	var document document
	if err := json.Unmarshal(got, &document); err != nil {
		t.Fatal(err)
	}
	if len(document.Paths) != 2 || len(document.Components.Schemas) != 2 {
		t.Fatalf("unexpected contract shape: %#v", document)
	}
	if got := document.Components.Schemas["Pet"].Properties["owner"].Ref; got != "#/components/schemas/shared_Profile" {
		t.Fatalf("owner reference = %q", got)
	}
}

func normalizeNewlines(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(value, "\r\n", "\n"), "\r", "\n")
}

func TestGenerateRejectsInvalidInputWithoutReplacingOutput(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  string
	}{
		{"malformed", "testdata/invalid/malformed", "operation Broken"},
		{"duplicate", "testdata/invalid/duplicate", "duplicate"},
		{"unresolved", "testdata/invalid/unresolved", "unresolved model"},
		{"unsupported", "testdata/invalid/unsupported", "Values"},
		{"unknown directive", "testdata/invalid/unknown-directive", "unknown annotation key"},
		{"missing route", "testdata/invalid/missing-route", "requires method, route"},
		{"missing method", "testdata/invalid/missing-method", "requires method, route"},
		{"missing response", "testdata/invalid/missing-response", "requires method, route"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "contract.json")
			if err := os.WriteFile(output, []byte("existing"), 0o644); err != nil {
				t.Fatal(err)
			}
			err := generate(tc.input, output)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
			got, readErr := os.ReadFile(output)
			if readErr != nil || string(got) != "existing" {
				t.Fatalf("output was changed: %q, %v", got, readErr)
			}
		})
	}

	t.Run("absent output", func(t *testing.T) {
		output := filepath.Join(t.TempDir(), "contract.json")
		if err := generate("testdata/invalid/unknown-directive", output); err == nil {
			t.Fatal("expected generation error")
		}
		if _, err := os.Stat(output); !os.IsNotExist(err) {
			t.Fatalf("invalid generation created output: %v", err)
		}
	})
}

func TestGenerateReplacesValidExistingOutput(t *testing.T) {
	output := filepath.Join(t.TempDir(), "contract.json")
	if err := os.WriteFile(output, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := generate("testdata", output); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "stale") || !strings.Contains(string(data), `"openapi": "3.0.3"`) {
		t.Fatalf("existing output was not replaced: %s", data)
	}
}
