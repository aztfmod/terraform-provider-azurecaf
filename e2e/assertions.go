package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

type terraformOutputValue struct {
	Sensitive bool            `json:"sensitive"`
	Type      interface{}     `json:"type"`
	Value     json.RawMessage `json:"value"`
}

var buildProviderOnce sync.Once
var buildProviderErr error

func setupTerraformTest(t *testing.T, testName, tfConfig string) string {
	t.Helper()

	buildProvider(t)

	testDir := newLocalTestDir(t, testName)
	configPath := filepath.Join(testDir, "main.tf")
	if err := os.WriteFile(configPath, []byte(tfConfig), 0o644); err != nil {
		t.Fatalf("failed to write terraform config: %v", err)
	}

	providerPath, err := filepath.Abs("../terraform-provider-azurecaf")
	if err != nil {
		t.Fatalf("failed to resolve provider binary path: %v", err)
	}

	overrideConfig := fmt.Sprintf(`provider_installation {
  dev_overrides {
    "aztfmod/azurecaf" = "%s"
  }
  direct {}
}
`, filepath.Dir(providerPath))

	rcPath := filepath.Join(testDir, "terraform.rc")
	if err := os.WriteFile(rcPath, []byte(overrideConfig), 0o644); err != nil {
		t.Fatalf("failed to write terraform.rc: %v", err)
	}

	runTerraformExpectSuccess(t, testDir, "init", "-input=false")
	return testDir
}

func buildProvider(t *testing.T) {
	t.Helper()

	buildProviderOnce.Do(func() {
		makePath, err := findMakeBinary()
		if err != nil {
			buildProviderErr = fmt.Errorf("failed to find make binary: %w", err)
			return
		}

		cmd := exec.Command(makePath, "build")
		cmd.Dir = ".."
		output, err := cmd.CombinedOutput()
		if err != nil {
			buildProviderErr = fmt.Errorf("failed to build provider: %w\nOutput: %s", err, output)
			return
		}
	})

	if buildProviderErr != nil {
		t.Fatal(buildProviderErr)
	}
}

func newLocalTestDir(t *testing.T, testName string) string {
	t.Helper()

	relativeDir := filepath.Join(".e2e-work", sanitizeTestName(testName))
	absoluteDir, err := filepath.Abs(relativeDir)
	if err != nil {
		t.Fatalf("failed to resolve test directory %s: %v", relativeDir, err)
	}
	if err := os.RemoveAll(absoluteDir); err != nil {
		t.Fatalf("failed to clean test directory %s: %v", absoluteDir, err)
	}
	if err := os.MkdirAll(absoluteDir, 0o750); err != nil {
		t.Fatalf("failed to create test directory %s: %v", absoluteDir, err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(absoluteDir); err != nil {
			t.Logf("warning: failed to cleanup %s: %v", absoluteDir, err)
		}
	})

	return absoluteDir
}

func sanitizeTestName(name string) string {
	replacer := strings.NewReplacer("/", "-", " ", "-", ":", "-", "\\", "-", "\t", "-", "\n", "-")
	cleaned := replacer.Replace(strings.ToLower(name))
	cleaned = strings.Trim(cleaned, "-")
	if cleaned == "" {
		return fmt.Sprintf("e2e-%d", time.Now().UnixNano())
	}
	return cleaned
}

func terraformEnv(dir string) []string {
	rcPath := filepath.Join(dir, "terraform.rc")
	return append(os.Environ(),
		"TF_CLI_CONFIG_FILE="+rcPath,
		"CHECKPOINT_DISABLE=1",
		"TF_IN_AUTOMATION=1",
	)
}

func runTerraform(t *testing.T, dir string, args ...string) (string, int, error) {
	t.Helper()

	terraformPath, err := findTerraformBinary()
	if err != nil {
		t.Fatalf("failed to find terraform binary: %v", err)
	}

	cmd := exec.Command(terraformPath, args...)
	cmd.Dir = dir
	cmd.Env = terraformEnv(dir)
	output, err := cmd.CombinedOutput()
	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}
	return string(output), exitCode, err
}

func runTerraformExpectSuccess(t *testing.T, dir string, args ...string) string {
	t.Helper()

	output, _, err := runTerraform(t, dir, args...)
	if err != nil {
		t.Fatalf("terraform %s failed: %v\nOutput:\n%s", strings.Join(args, " "), err, output)
	}
	return output
}

func runTerraformPlan(t *testing.T, dir string) (string, int) {
	t.Helper()

	output, exitCode, err := runTerraform(t, dir, "plan", "-input=false", "-lock=false", "-detailed-exitcode")
	if err != nil && exitCode != 2 {
		t.Fatalf("terraform plan failed: %v\nOutput:\n%s", err, output)
	}
	return output, exitCode
}

func runTerraformPlanExpectError(t *testing.T, dir string) string {
	t.Helper()

	output, _, err := runTerraform(t, dir, "plan", "-input=false", "-lock=false")
	if err == nil {
		t.Fatalf("expected terraform plan to fail for %s", dir)
	}
	return output
}

func applyTerraform(t *testing.T, dir string) string {
	t.Helper()
	return runTerraformExpectSuccess(t, dir, "apply", "-auto-approve", "-input=false", "-lock=false")
}

func runTerraformApplyExpectError(t *testing.T, dir string) string {
	t.Helper()

	output, _, err := runTerraform(t, dir, "apply", "-auto-approve", "-input=false", "-lock=false")
	if err == nil {
		t.Fatalf("expected terraform apply to fail for %s", dir)
	}
	return output
}

func assertOutputMatches(t *testing.T, dir, outputName, expectedPattern string) {
	t.Helper()

	value := getOutputString(t, dir, outputName)
	matched, err := regexp.MatchString(expectedPattern, value)
	if err != nil {
		t.Fatalf("invalid output pattern %q: %v", expectedPattern, err)
	}
	if !matched {
		t.Fatalf("output %s value %q does not match %q", outputName, value, expectedPattern)
	}
}

func assertOutputContains(t *testing.T, dir, outputName, substring string) {
	t.Helper()

	value := getOutputString(t, dir, outputName)
	if !strings.Contains(value, substring) {
		t.Fatalf("output %s value %q does not contain %q", outputName, value, substring)
	}
}

func assertOutputEquals(t *testing.T, dir, outputName, expected string) {
	t.Helper()

	value := getOutputString(t, dir, outputName)
	if value != expected {
		t.Fatalf("output %s mismatch: expected %q, got %q", outputName, expected, value)
	}
}

func assertPlanNoChanges(t *testing.T, dir string) {
	t.Helper()

	output, exitCode := runTerraformPlan(t, dir)
	if exitCode != 0 {
		t.Fatalf("expected no terraform changes, got exit code %d\nOutput:\n%s", exitCode, output)
	}
	if !strings.Contains(output, "No changes") && !strings.Contains(output, "0 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected terraform plan to report no changes\nOutput:\n%s", output)
	}
}

func terraformOutputs(t *testing.T, dir string) map[string]terraformOutputValue {
	t.Helper()

	output := runTerraformExpectSuccess(t, dir, "output", "-json")
	outputs := make(map[string]terraformOutputValue)
	if err := json.Unmarshal([]byte(output), &outputs); err != nil {
		t.Fatalf("failed to parse terraform output JSON: %v\nOutput:\n%s", err, output)
	}
	return outputs
}

func getOutputString(t *testing.T, dir, outputName string) string {
	t.Helper()

	outputs := terraformOutputs(t, dir)
	output, exists := outputs[outputName]
	if !exists {
		t.Fatalf("terraform output %q not found", outputName)
	}

	var value string
	if err := json.Unmarshal(output.Value, &value); err != nil {
		t.Fatalf("failed to parse output %s as string: %v", outputName, err)
	}
	return value
}

func getOutputStringMap(t *testing.T, dir, outputName string) map[string]string {
	t.Helper()

	outputs := terraformOutputs(t, dir)
	output, exists := outputs[outputName]
	if !exists {
		t.Fatalf("terraform output %q not found", outputName)
	}

	value := make(map[string]string)
	if err := json.Unmarshal(output.Value, &value); err != nil {
		t.Fatalf("failed to parse output %s as string map: %v", outputName, err)
	}
	return value
}
