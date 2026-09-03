package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateSelectedAccessors(t *testing.T) {
	output := filepath.Join(t.TempDir(), "account_accessors.go")
	if err := generate("testdata", "Account", []string{"ID", "Email"}, output); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/zz_account_accessors.golden.go")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("generated accessors differ\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestGenerateRejectsMissingOrUnexportedField(t *testing.T) {
	output := filepath.Join(t.TempDir(), "account_accessors.go")
	if err := generate("testdata", "Account", []string{"Missing"}, output); err == nil {
		t.Fatal("expected missing-field error")
	}
	if err := generate("testdata", "Account", []string{"internal"}, output); err == nil {
		t.Fatal("expected unexported-field error")
	}
}
