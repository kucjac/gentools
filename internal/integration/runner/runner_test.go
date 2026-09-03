package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("RUNNER_TEST_FAKE_GO") == "1" {
		if content, err := os.ReadFile("main_test.go"); err == nil && strings.Contains(string(content), "mutated") {
			os.Exit(1)
		}
		fmt.Printf("executed %s\n", strings.Join(os.Args[1:], " "))
		os.Exit(0)
	}
	os.Exit(m.Run())
}

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
	goBin := os.Args[0]
	t.Setenv("RUNNER_TEST_FAKE_GO", "1")
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
	goBin := os.Args[0]
	t.Setenv("RUNNER_TEST_FAKE_GO", "1")
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
