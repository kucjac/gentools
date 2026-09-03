// Command runner executes the checked-in consumer scenario contract.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Scenario struct {
	Name               string `json:"name"`
	SourceInput        string `json:"source_input"`
	ModuleBoundary     string `json:"module_boundary"`
	ExpectedResult     string `json:"expected_result"`
	FailureExpectation string `json:"failure_expectation"`
	MutationCheck      string `json:"mutation_check"`
	MutationFile       string `json:"mutation_file"`
	MutationOld        string `json:"mutation_old"`
	MutationNew        string `json:"mutation_new"`
}
type manifest struct {
	Scenarios []Scenario `json:"scenarios"`
}

func Load(path string) ([]Scenario, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if len(m.Scenarios) != 4 {
		return nil, fmt.Errorf("expected four scenarios, got %d", len(m.Scenarios))
	}
	for _, s := range m.Scenarios {
		if s.Name == "" || s.SourceInput == "" || s.ModuleBoundary == "" || s.ExpectedResult == "" || s.FailureExpectation == "" || s.MutationCheck == "" || s.MutationFile == "" || s.MutationOld == "" || s.MutationNew == "" {
			return nil, fmt.Errorf("incomplete scenario %q", s.Name)
		}
	}
	return m.Scenarios, nil
}

// Run executes a manifest-selected consumer scenario from the integration
// module directory. It deliberately labels this as candidate-source evidence:
// publication is verified separately by the tag-triggered smoke workflow.
func Run(manifestPath, selection, goBin string, output func(string)) error {
	return run(manifestPath, selection, goBin, false, output)
}

// RunWithMutationVerification proves each selected scenario's assertions fail
// when its declared expectation is deliberately changed.
func RunWithMutationVerification(manifestPath, selection, goBin string, output func(string)) error {
	return run(manifestPath, selection, goBin, true, output)
}

func run(manifestPath, selection, goBin string, verifyMutations bool, output func(string)) error {
	scenarios, err := Load(manifestPath)
	if err != nil {
		return err
	}
	if goBin == "" {
		goBin = "go"
	}
	var selected []Scenario
	for _, scenario := range scenarios {
		if selection == "all" || selection == scenario.Name {
			selected = append(selected, scenario)
		}
	}
	if len(selected) == 0 {
		return fmt.Errorf("unknown scenario: %s", selection)
	}
	for _, scenario := range selected {
		output(fmt.Sprintf("scenario=%s boundary=candidate-source expected=%s", scenario.Name, scenario.ExpectedResult))
		dir := filepath.Join(filepath.Dir(manifestPath), scenario.Name)
		result, err := runScenario(goBin, dir)
		outputResult(output, result)
		if err != nil {
			return fmt.Errorf("scenario %s failed: %w", scenario.Name, err)
		}
		if verifyMutations {
			if err := verifyMutation(goBin, dir, scenario, output); err != nil {
				return err
			}
		}
	}
	return nil
}

func runScenario(goBin, dir string) ([]byte, error) {
	cmd := exec.Command(goBin, "test", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GODEBUG=gotypesalias=1")
	cmd.Env = removeEnv(cmd.Env, "GOROOT", "GOTOOLCHAIN")
	return cmd.CombinedOutput()
}

func outputResult(output func(string), result []byte) {
	if len(result) > 0 {
		output(strings.TrimSpace(string(result)))
	}
}

func verifyMutation(goBin, dir string, scenario Scenario, output func(string)) error {
	path := filepath.Join(dir, scenario.MutationFile)
	original, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("scenario %s mutation input: %w", scenario.Name, err)
	}
	if !strings.Contains(string(original), scenario.MutationOld) {
		return fmt.Errorf("scenario %s mutation expectation not found", scenario.Name)
	}
	mutated := strings.Replace(string(original), scenario.MutationOld, scenario.MutationNew, 1)
	if err := os.WriteFile(path, []byte(mutated), 0o644); err != nil {
		return fmt.Errorf("scenario %s write mutation: %w", scenario.Name, err)
	}
	defer func() { _ = os.WriteFile(path, original, 0o644) }()
	result, err := runScenario(goBin, dir)
	outputResult(output, result)
	if err == nil {
		return fmt.Errorf("scenario %s mutation unexpectedly passed", scenario.Name)
	}
	output(fmt.Sprintf("scenario=%s mutation=verified", scenario.Name))
	return nil
}

func removeEnv(values []string, names ...string) []string {
	result := values[:0]
	for _, value := range values {
		keep := true
		for _, name := range names {
			if strings.HasPrefix(value, name+"=") {
				keep = false
				break
			}
		}
		if keep {
			result = append(result, value)
		}
	}
	return result
}

func main() {
	manifestPath := flag.String("manifest", "scenarios/scenarios.json", "scenario manifest")
	selection := flag.String("scenario", "all", "scenario name or all")
	goBin := flag.String("go", "go", "Go command")
	verifyMutations := flag.Bool("verify-mutations", false, "prove selected scenario assertions fail when mutated")
	flag.Parse()
	runFn := Run
	if *verifyMutations {
		runFn = RunWithMutationVerification
	}
	if err := runFn(*manifestPath, *selection, *goBin, func(line string) { fmt.Println(line) }); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
