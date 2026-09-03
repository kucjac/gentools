package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsIncompleteScenario(t *testing.T) {
	p := filepath.Join(t.TempDir(), "scenarios.json")
	if err := os.WriteFile(p, []byte(`{"scenarios":[{}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(p); err == nil {
		t.Fatal("expected incomplete scenario failure")
	}
}

func TestRunSelectsScenarioAndLabelsBoundary(t *testing.T) {
	dir := t.TempDir()
	manifest := filepath.Join(dir, "scenarios.json")
	if err := os.MkdirAll(filepath.Join(dir, "inspect"), 0o755); err != nil {
		t.Fatal(err)
	}
	content := `{"scenarios":[
{"name":"inspect","source_input":"inspect","module_boundary":"snapshot","expected_result":"found","failure_expectation":"failure","mutation_check":"change","mutation_file":"main_test.go","mutation_old":"expected","mutation_new":"mutated"},
{"name":"generated","source_input":"inspect","module_boundary":"snapshot","expected_result":"found","failure_expectation":"failure","mutation_check":"change","mutation_file":"main_test.go","mutation_old":"expected","mutation_new":"mutated"},
{"name":"crosspackage","source_input":"inspect","module_boundary":"snapshot","expected_result":"found","failure_expectation":"failure","mutation_check":"change","mutation_file":"main_test.go","mutation_old":"expected","mutation_new":"mutated"},
{"name":"invalid","source_input":"inspect","module_boundary":"snapshot","expected_result":"found","failure_expectation":"failure","mutation_check":"change","mutation_file":"main_test.go","mutation_old":"expected","mutation_new":"mutated"}]}`
	if err := os.WriteFile(manifest, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "inspect", "main_test.go"), []byte("expected"), 0o644); err != nil {
		t.Fatal(err)
	}
	goBin := filepath.Join(dir, "go")
	if err := os.WriteFile(goBin, []byte("#!/bin/sh\nprintf 'executed %s\\n' \"$*\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	var lines []string
	if err := Run(manifest, "inspect", goBin, func(line string) { lines = append(lines, line) }); err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "scenario=inspect boundary=candidate-source") || !strings.Contains(joined, "executed test ./...") {
		t.Fatalf("unexpected runner output: %s", joined)
	}
}

func TestRunWithMutationVerificationRequiresFailure(t *testing.T) {
	dir := t.TempDir()
	manifest := filepath.Join(dir, "scenarios.json")
	if err := os.MkdirAll(filepath.Join(dir, "inspect"), 0o755); err != nil {
		t.Fatal(err)
	}
	content := `{"scenarios":[
{"name":"inspect","source_input":"inspect","module_boundary":"snapshot","expected_result":"found","failure_expectation":"failure","mutation_check":"change","mutation_file":"main_test.go","mutation_old":"expected","mutation_new":"mutated"},
{"name":"generated","source_input":"inspect","module_boundary":"snapshot","expected_result":"found","failure_expectation":"failure","mutation_check":"change","mutation_file":"main_test.go","mutation_old":"expected","mutation_new":"mutated"},
{"name":"crosspackage","source_input":"inspect","module_boundary":"snapshot","expected_result":"found","failure_expectation":"failure","mutation_check":"change","mutation_file":"main_test.go","mutation_old":"expected","mutation_new":"mutated"},
{"name":"invalid","source_input":"inspect","module_boundary":"snapshot","expected_result":"found","failure_expectation":"failure","mutation_check":"change","mutation_file":"main_test.go","mutation_old":"expected","mutation_new":"mutated"}]}`
	if err := os.WriteFile(manifest, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	testFile := filepath.Join(dir, "inspect", "main_test.go")
	if err := os.WriteFile(testFile, []byte("expected"), 0o644); err != nil {
		t.Fatal(err)
	}
	goBin := filepath.Join(dir, "go")
	script := "#!/bin/sh\nif grep -q mutated main_test.go; then exit 1; fi\nexit 0\n"
	if err := os.WriteFile(goBin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	var lines []string
	if err := RunWithMutationVerification(manifest, "inspect", goBin, func(line string) { lines = append(lines, line) }); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(lines, "\n"), "scenario=inspect mutation=verified") {
		t.Fatalf("mutation verification was not reported: %v", lines)
	}
	contentAfter, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(contentAfter) != "expected" {
		t.Fatalf("mutation was not restored: %q", contentAfter)
	}
}

func TestRunRejectsUnknownScenario(t *testing.T) {
	if err := Run("testdata/does-not-exist", "unknown", "go", func(string) {}); err == nil {
		t.Fatal("expected manifest failure")
	}
}
