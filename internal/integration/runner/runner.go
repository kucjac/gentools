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
		if s.Name == "" || s.SourceInput == "" || s.ModuleBoundary == "" || s.ExpectedResult == "" || s.FailureExpectation == "" || s.MutationCheck == "" {
			return nil, fmt.Errorf("incomplete scenario %q", s.Name)
		}
	}
	return m.Scenarios, nil
}

// Run executes a manifest-selected consumer scenario from the integration
// module directory. It deliberately labels this as candidate-source evidence:
// publication is verified separately by the tag-triggered smoke workflow.
func Run(manifestPath, selection, goBin string, output func(string)) error {
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
		cmd := exec.Command(goBin, "test", "./...")
		cmd.Dir = filepath.Join(filepath.Dir(manifestPath), scenario.Name)
		cmd.Env = append(os.Environ(), "GODEBUG=gotypesalias=1")
		cmd.Env = removeEnv(cmd.Env, "GOROOT", "GOTOOLCHAIN")
		result, err := cmd.CombinedOutput()
		if len(result) > 0 {
			output(strings.TrimSpace(string(result)))
		}
		if err != nil {
			return fmt.Errorf("scenario %s failed: %w", scenario.Name, err)
		}
	}
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
	flag.Parse()
	if err := Run(*manifestPath, *selection, *goBin, func(line string) { fmt.Println(line) }); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
