package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

const lifecycleProviderSource = "aztfmod.com/local/azurecaf"

type lifecycleHarness struct {
	terraformPath string
	cliConfigPath string
	rootDir       string
}

type lifecyclePlan struct {
	PlannedValues struct {
		Outputs map[string]struct {
			Value json.RawMessage `json:"value"`
		} `json:"outputs"`
	} `json:"planned_values"`
	ResourceChanges []struct {
		Address string `json:"address"`
		Change  struct {
			Actions      []string                   `json:"actions"`
			After        map[string]json.RawMessage `json:"after"`
			AfterUnknown map[string]json.RawMessage `json:"after_unknown"`
		} `json:"change"`
	} `json:"resource_changes"`
}

type lifecyclePlanResult struct {
	exitCode int
	path     string
	output   string
	document lifecyclePlan
}

type lifecycleOutput struct {
	Value json.RawMessage `json:"value"`
}

func TestE2ELifecycle(t *testing.T) {
	harness := newLifecycleHarness(t)

	t.Run("seeded resource plan apply replace refresh destroy", func(t *testing.T) {
		dir := harness.newWorkspace(t, "seeded", seededLifecycleConfig(42))

		firstPlan := harness.plan(t, dir, "create")
		if firstPlan.exitCode != 2 {
			t.Fatalf("expected initial plan to contain changes, got exit code %d", firstPlan.exitCode)
		}
		const initialResult = "prod-rg-checkout-hrukp-sea"
		assertPlannedString(t, firstPlan.document, "result", initialResult)
		assertResourceActions(t, firstPlan.document, "azurecaf_name.seeded", []string{"create"})

		harness.applyPlan(t, dir, firstPlan.path)
		if got := harness.outputString(t, dir, "result"); got != initialResult {
			t.Fatalf("expected applied result %q, got %q", initialResult, got)
		}
		harness.assertNoChanges(t, dir, "post-apply")

		refreshPlan := harness.plan(t, dir, "refresh-only", "-refresh-only")
		if refreshPlan.exitCode != 0 {
			t.Fatalf("expected refresh-only plan to be empty, got exit code %d\n%s", refreshPlan.exitCode, refreshPlan.output)
		}

		harness.writeConfig(t, dir, seededLifecycleConfig(43))
		replacementPlan := harness.plan(t, dir, "replacement")
		if replacementPlan.exitCode != 2 {
			t.Fatalf("expected seed change to require replacement, got exit code %d", replacementPlan.exitCode)
		}
		assertReplacementActions(t, replacementPlan.document, "azurecaf_name.seeded")
		replacementResult := plannedString(t, replacementPlan.document, "result")
		if replacementResult == initialResult {
			t.Fatalf("expected seed replacement to change result, both plans produced %q", initialResult)
		}

		harness.applyPlan(t, dir, replacementPlan.path)
		if got := harness.outputString(t, dir, "result"); got != replacementResult {
			t.Fatalf("expected replacement result %q, got %q", replacementResult, got)
		}
		harness.assertNoChanges(t, dir, "post-replacement")
		harness.destroyAndAssertEmpty(t, dir)
	})

	t.Run("unseeded and zero seed remain stable in state", func(t *testing.T) {
		dir := harness.newWorkspace(t, "unseeded", unseededLifecycleConfig())

		firstPlan := harness.plan(t, dir, "create")
		if firstPlan.exitCode != 2 {
			t.Fatalf("expected initial plan to contain changes, got exit code %d", firstPlan.exitCode)
		}
		assertPlannedUnknown(t, firstPlan.document, "omitted_result")
		assertPlannedUnknown(t, firstPlan.document, "zero_result")
		assertResourceResultUnknown(t, firstPlan.document, "azurecaf_name.omitted")
		assertResourceResultUnknown(t, firstPlan.document, "azurecaf_name.zero")

		harness.applyPlan(t, dir, firstPlan.path)
		omittedResult := harness.outputString(t, dir, "omitted_result")
		zeroResult := harness.outputString(t, dir, "zero_result")
		assertRandomResourceName(t, omittedResult)
		assertRandomResourceName(t, zeroResult)

		harness.assertNoChanges(t, dir, "post-apply")
		refreshPlan := harness.plan(t, dir, "refresh-only", "-refresh-only")
		if refreshPlan.exitCode != 0 {
			t.Fatalf("expected refresh-only plan to preserve random state, got exit code %d\n%s", refreshPlan.exitCode, refreshPlan.output)
		}
		harness.mustRun(t, dir, "apply", "-input=false", "-no-color", "-auto-approve", "-refresh-only")
		if got := harness.outputString(t, dir, "omitted_result"); got != omittedResult {
			t.Fatalf("omitted-seed result changed during refresh: before %q, after %q", omittedResult, got)
		}
		if got := harness.outputString(t, dir, "zero_result"); got != zeroResult {
			t.Fatalf("zero-seed result changed during refresh: before %q, after %q", zeroResult, got)
		}
		harness.assertNoChanges(t, dir, "post-refresh")
		harness.destroyAndAssertEmpty(t, dir)
	})

	t.Run("unknown list input defers then converges", func(t *testing.T) {
		dir := harness.newWorkspace(t, "unknown-input", unknownInputLifecycleConfig())

		firstPlan := harness.plan(t, dir, "create")
		if firstPlan.exitCode != 2 {
			t.Fatalf("expected initial plan to contain changes, got exit code %d", firstPlan.exitCode)
		}
		assertPlannedUnknown(t, firstPlan.document, "result")
		assertResourceResultUnknown(t, firstPlan.document, "azurecaf_name.deferred")

		harness.applyPlan(t, dir, firstPlan.path)
		const expected = "dev-rg-app-hrukp"
		if got := harness.outputString(t, dir, "result"); got != expected {
			t.Fatalf("expected deferred result %q after apply, got %q", expected, got)
		}
		harness.assertNoChanges(t, dir, "post-apply")
		harness.destroyAndAssertEmpty(t, dir)
	})

	t.Run("multiple resource types match plan and state", func(t *testing.T) {
		dir := harness.newWorkspace(t, "multiple-types", multipleTypesLifecycleConfig())

		firstPlan := harness.plan(t, dir, "create")
		if firstPlan.exitCode != 2 {
			t.Fatalf("expected initial plan to contain changes, got exit code %d", firstPlan.exitCode)
		}
		const primaryResult = "prod-rg-checkout-hruk"
		expectedResults := map[string]string{
			"azurerm_key_vault":       "prod-kv-checkout-hruk",
			"azurerm_storage_account": "prodstcheckouthruk",
		}
		assertPlannedString(t, firstPlan.document, "result", primaryResult)
		plannedResults := plannedStringMap(t, firstPlan.document, "results")
		if !reflect.DeepEqual(plannedResults, expectedResults) {
			t.Fatalf("unexpected planned results: expected %#v, got %#v", expectedResults, plannedResults)
		}

		harness.applyPlan(t, dir, firstPlan.path)
		if got := harness.outputString(t, dir, "result"); got != primaryResult {
			t.Fatalf("expected primary result %q, got %q", primaryResult, got)
		}
		if got := harness.outputStringMap(t, dir, "results"); !reflect.DeepEqual(got, expectedResults) {
			t.Fatalf("plan/apply results mismatch: expected %#v, got %#v", expectedResults, got)
		}
		harness.assertNoChanges(t, dir, "post-apply")
		harness.destroyAndAssertEmpty(t, dir)
	})

	t.Run("data source seed semantics persist across plans", func(t *testing.T) {
		dir := harness.newWorkspace(t, "data-source-seeds", dataSourceLifecycleConfig())

		firstPlan := harness.plan(t, dir, "first")
		if firstPlan.exitCode != 2 {
			t.Fatalf("expected initial data-source plan to contain output changes, got exit code %d", firstPlan.exitCode)
		}
		const zeroResult = "rg-app-cubyhizzkakb"
		const seededResult = "rg-app-hrukpttuezpt"
		assertPlannedString(t, firstPlan.document, "zero_result", zeroResult)
		assertPlannedString(t, firstPlan.document, "seeded_result", seededResult)
		firstUnseeded := plannedString(t, firstPlan.document, "unseeded_result")
		assertRandomResourceName(t, firstUnseeded)
		harness.applyPlan(t, dir, firstPlan.path)

		secondPlan := harness.plan(t, dir, "second")
		if secondPlan.exitCode != 2 {
			t.Fatalf("expected unseeded data source to produce an output change, got exit code %d", secondPlan.exitCode)
		}
		assertPlannedString(t, secondPlan.document, "zero_result", zeroResult)
		assertPlannedString(t, secondPlan.document, "seeded_result", seededResult)
		secondUnseeded := plannedString(t, secondPlan.document, "unseeded_result")
		if secondUnseeded == firstUnseeded {
			t.Fatalf("expected omitted data-source seed to generate a fresh value, both plans produced %q", firstUnseeded)
		}
		harness.applyPlan(t, dir, secondPlan.path)

		thirdPlan := harness.plan(t, dir, "third")
		assertPlannedString(t, thirdPlan.document, "zero_result", zeroResult)
		assertPlannedString(t, thirdPlan.document, "seeded_result", seededResult)
		thirdUnseeded := plannedString(t, thirdPlan.document, "unseeded_result")
		if thirdUnseeded == secondUnseeded {
			t.Fatalf("expected omitted data-source seed to refresh again, both plans produced %q", secondUnseeded)
		}
	})

	t.Run("import converges without resource drift", func(t *testing.T) {
		dir := harness.newWorkspace(t, "import", importLifecycleConfig())

		harness.mustRun(
			t,
			dir,
			"import",
			"-input=false",
			"-no-color",
			"azurecaf_name.imported",
			"azurerm_storage_account:stmyexistingapp",
		)
		importPlan := harness.plan(t, dir, "post-import")
		assertResourceActions(t, importPlan.document, "azurecaf_name.imported", []string{"no-op"})
		harness.applyPlan(t, dir, importPlan.path)
		if got := harness.outputString(t, dir, "result"); got != "stmyexistingapp" {
			t.Fatalf("expected imported result %q, got %q", "stmyexistingapp", got)
		}
		harness.assertNoChanges(t, dir, "post-import-apply")
		harness.destroyAndAssertEmpty(t, dir)
	})
}

func newLifecycleHarness(t *testing.T) *lifecycleHarness {
	t.Helper()

	terraformPath, err := findTerraformBinary()
	if err != nil {
		t.Fatal(err)
	}
	goPath, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("go binary not found in PATH: %v", err)
	}
	projectRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("failed to resolve project root: %v", err)
	}

	rootDir := t.TempDir()
	mirrorDir := filepath.Join(rootDir, "mirror")
	providerDir := filepath.Join(
		mirrorDir,
		"aztfmod.com",
		"local",
		"azurecaf",
		"2.0.0",
		fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH),
	)
	if err := os.MkdirAll(providerDir, 0750); err != nil {
		t.Fatalf("failed to create provider mirror: %v", err)
	}
	providerBinary := filepath.Join(providerDir, "terraform-provider-azurecaf_v2.0.0")
	build := exec.Command(goPath, "build", "-o", providerBinary, ".")
	build.Dir = projectRoot
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("failed to build provider: %v\n%s", err, output)
	}
	if err := os.Chmod(providerBinary, 0750); err != nil {
		t.Fatalf("failed to make provider executable: %v", err)
	}

	cliConfigPath := filepath.Join(rootDir, "terraform.rc")
	cliConfig := fmt.Sprintf(`provider_installation {
  filesystem_mirror {
    path    = %q
    include = [%q]
  }
  direct {
    exclude = [%q]
  }
}
`, mirrorDir, lifecycleProviderSource, lifecycleProviderSource)
	if err := os.WriteFile(cliConfigPath, []byte(cliConfig), 0600); err != nil {
		t.Fatalf("failed to write Terraform CLI config: %v", err)
	}

	return &lifecycleHarness{
		terraformPath: terraformPath,
		cliConfigPath: cliConfigPath,
		rootDir:       rootDir,
	}
}

func (h *lifecycleHarness) newWorkspace(t *testing.T, name, config string) string {
	t.Helper()
	dir := filepath.Join(h.rootDir, name)
	if err := os.MkdirAll(dir, 0750); err != nil {
		t.Fatalf("failed to create workspace: %v", err)
	}
	h.writeConfig(t, dir, config)
	h.mustRun(t, dir, "init", "-input=false", "-no-color")
	return dir
}

func (h *lifecycleHarness) writeConfig(t *testing.T, dir, config string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "main.tf"), []byte(config), 0644); err != nil {
		t.Fatalf("failed to write Terraform configuration: %v", err)
	}
}

func (h *lifecycleHarness) run(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, h.terraformPath, args...)
	cmd.Dir = dir
	cmd.Env = append(
		os.Environ(),
		"TF_CLI_CONFIG_FILE="+h.cliConfigPath,
		"CHECKPOINT_DISABLE=1",
		"TF_IN_AUTOMATION=1",
		"TF_CLI_ARGS_init=-upgrade=false",
	)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("terraform %s timed out\n%s", strings.Join(args, " "), output)
	}
	return string(output), err
}

func (h *lifecycleHarness) mustRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	output, err := h.run(t, dir, args...)
	if err != nil {
		t.Fatalf("terraform %s failed: %v\n%s", strings.Join(args, " "), err, output)
	}
	return output
}

func (h *lifecycleHarness) plan(t *testing.T, dir, name string, extraArgs ...string) lifecyclePlanResult {
	t.Helper()
	planPath := filepath.Join(dir, name+".tfplan")
	args := []string{
		"plan",
		"-input=false",
		"-no-color",
		"-detailed-exitcode",
		"-out=" + planPath,
	}
	args = append(args, extraArgs...)
	output, err := h.run(t, dir, args...)
	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 2 {
			t.Fatalf("terraform plan failed: %v\n%s", err, output)
		}
		exitCode = 2
	}

	showOutput := h.mustRun(t, dir, "show", "-json", planPath)
	var document lifecyclePlan
	if err := json.Unmarshal([]byte(showOutput), &document); err != nil {
		t.Fatalf("failed to decode Terraform plan JSON: %v", err)
	}
	return lifecyclePlanResult{
		exitCode: exitCode,
		path:     planPath,
		output:   output,
		document: document,
	}
}

func (h *lifecycleHarness) applyPlan(t *testing.T, dir, planPath string) {
	t.Helper()
	h.mustRun(t, dir, "apply", "-input=false", "-no-color", planPath)
}

func (h *lifecycleHarness) assertNoChanges(t *testing.T, dir, name string) {
	t.Helper()
	plan := h.plan(t, dir, name)
	if plan.exitCode != 0 {
		t.Fatalf("expected no changes, got exit code %d\n%s", plan.exitCode, plan.output)
	}
	if !strings.Contains(plan.output, "No changes.") {
		t.Fatalf("expected Terraform no-change confirmation, got:\n%s", plan.output)
	}
}

func (h *lifecycleHarness) outputValues(t *testing.T, dir string) map[string]lifecycleOutput {
	t.Helper()
	output := h.mustRun(t, dir, "output", "-json")
	values := map[string]lifecycleOutput{}
	if err := json.Unmarshal([]byte(output), &values); err != nil {
		t.Fatalf("failed to decode Terraform outputs: %v", err)
	}
	return values
}

func (h *lifecycleHarness) outputString(t *testing.T, dir, name string) string {
	t.Helper()
	output, ok := h.outputValues(t, dir)[name]
	if !ok {
		t.Fatalf("Terraform output %q not found", name)
	}
	var value string
	if err := json.Unmarshal(output.Value, &value); err != nil {
		t.Fatalf("Terraform output %q is not a string: %v", name, err)
	}
	return value
}

func (h *lifecycleHarness) outputStringMap(t *testing.T, dir, name string) map[string]string {
	t.Helper()
	output, ok := h.outputValues(t, dir)[name]
	if !ok {
		t.Fatalf("Terraform output %q not found", name)
	}
	var value map[string]string
	if err := json.Unmarshal(output.Value, &value); err != nil {
		t.Fatalf("Terraform output %q is not a string map: %v", name, err)
	}
	return value
}

func (h *lifecycleHarness) destroyAndAssertEmpty(t *testing.T, dir string) {
	t.Helper()
	h.mustRun(t, dir, "destroy", "-input=false", "-no-color", "-auto-approve")
	stateList := strings.TrimSpace(h.mustRun(t, dir, "state", "list"))
	if stateList != "" {
		t.Fatalf("expected empty Terraform state after destroy, got:\n%s", stateList)
	}
}

func plannedString(t *testing.T, plan lifecyclePlan, name string) string {
	t.Helper()
	output, ok := plan.PlannedValues.Outputs[name]
	if !ok || len(output.Value) == 0 {
		t.Fatalf("planned output %q is unknown or missing", name)
	}
	var value string
	if err := json.Unmarshal(output.Value, &value); err != nil {
		t.Fatalf("planned output %q is not a string: %v", name, err)
	}
	return value
}

func plannedStringMap(t *testing.T, plan lifecyclePlan, name string) map[string]string {
	t.Helper()
	output, ok := plan.PlannedValues.Outputs[name]
	if !ok || len(output.Value) == 0 {
		t.Fatalf("planned output %q is unknown or missing", name)
	}
	var value map[string]string
	if err := json.Unmarshal(output.Value, &value); err != nil {
		t.Fatalf("planned output %q is not a string map: %v", name, err)
	}
	return value
}

func assertPlannedString(t *testing.T, plan lifecyclePlan, name, expected string) {
	t.Helper()
	if got := plannedString(t, plan, name); got != expected {
		t.Fatalf("expected planned output %q to be %q, got %q", name, expected, got)
	}
}

func assertPlannedUnknown(t *testing.T, plan lifecyclePlan, name string) {
	t.Helper()
	output, ok := plan.PlannedValues.Outputs[name]
	if !ok {
		t.Fatalf("planned output %q not found", name)
	}
	if len(output.Value) != 0 && string(output.Value) != "null" {
		t.Fatalf("expected planned output %q to be unknown, got %s", name, output.Value)
	}
}

func resourceChange(t *testing.T, plan lifecyclePlan, address string) struct {
	Actions      []string
	After        map[string]json.RawMessage
	AfterUnknown map[string]json.RawMessage
} {
	t.Helper()
	for _, change := range plan.ResourceChanges {
		if change.Address == address {
			return struct {
				Actions      []string
				After        map[string]json.RawMessage
				AfterUnknown map[string]json.RawMessage
			}{
				Actions:      change.Change.Actions,
				After:        change.Change.After,
				AfterUnknown: change.Change.AfterUnknown,
			}
		}
	}
	t.Fatalf("resource change %q not found", address)
	return struct {
		Actions      []string
		After        map[string]json.RawMessage
		AfterUnknown map[string]json.RawMessage
	}{}
}

func assertResourceActions(t *testing.T, plan lifecyclePlan, address string, expected []string) {
	t.Helper()
	change := resourceChange(t, plan, address)
	if !reflect.DeepEqual(change.Actions, expected) {
		t.Fatalf("expected %s actions %v, got %v", address, expected, change.Actions)
	}
}

func assertReplacementActions(t *testing.T, plan lifecyclePlan, address string) {
	t.Helper()
	actions := resourceChange(t, plan, address).Actions
	if len(actions) != 2 ||
		!containsString(actions, "create") ||
		!containsString(actions, "delete") {
		t.Fatalf("expected %s to be replaced, got actions %v", address, actions)
	}
}

func assertResourceResultUnknown(t *testing.T, plan lifecyclePlan, address string) {
	t.Helper()
	change := resourceChange(t, plan, address)
	raw, ok := change.AfterUnknown["result"]
	if !ok {
		t.Fatalf("result unknown marker missing for %s", address)
	}
	var unknown bool
	if err := json.Unmarshal(raw, &unknown); err != nil || !unknown {
		t.Fatalf("expected %s result to be unknown, got %s", address, raw)
	}
}

func assertRandomResourceName(t *testing.T, value string) {
	t.Helper()
	if !strings.HasPrefix(value, "rg-app-") || len(value) != len("rg-app-")+12 {
		t.Fatalf("expected generated resource-group name with 12 random characters, got %q", value)
	}
}

func containsString(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func lifecycleProviderConfig(body string) string {
	return fmt.Sprintf(`terraform {
  required_providers {
    azurecaf = {
      source  = %q
      version = "2.0.0"
    }
  }
}

provider "azurecaf" {}

%s
`, lifecycleProviderSource, body)
}

func seededLifecycleConfig(seed int) string {
	return lifecycleProviderConfig(fmt.Sprintf(`resource "azurecaf_name" "seeded" {
  name          = "checkout"
  resource_type = "azurerm_resource_group"
  prefixes      = ["prod"]
  suffixes      = ["sea"]
  random_length = 5
  random_seed   = %d
}

output "result" {
  value = azurecaf_name.seeded.result
}
`, seed))
}

func unseededLifecycleConfig() string {
	return lifecycleProviderConfig(`resource "azurecaf_name" "omitted" {
  name          = "app"
  resource_type = "azurerm_resource_group"
  random_length = 12
}

resource "azurecaf_name" "zero" {
  name          = "app"
  resource_type = "azurerm_resource_group"
  random_length = 12
  random_seed   = 0
}

output "omitted_result" {
  value = azurecaf_name.omitted.result
}

output "zero_result" {
  value = azurecaf_name.zero.result
}
`)
}

func unknownInputLifecycleConfig() string {
	return lifecycleProviderConfig(`resource "terraform_data" "prefix" {
  input = "dev"
}

resource "azurecaf_name" "deferred" {
  name          = "app"
  resource_type = "azurerm_resource_group"
  prefixes      = [terraform_data.prefix.output]
  random_length = 5
  random_seed   = 42
}

output "result" {
  value = azurecaf_name.deferred.result
}
`)
}

func multipleTypesLifecycleConfig() string {
	return lifecycleProviderConfig(`resource "azurecaf_name" "multiple" {
  name           = "checkout"
  resource_type  = "azurerm_resource_group"
  resource_types = ["azurerm_key_vault", "azurerm_storage_account"]
  prefixes       = ["prod"]
  random_length  = 4
  random_seed    = 42
}

output "result" {
  value = azurecaf_name.multiple.result
}

output "results" {
  value = azurecaf_name.multiple.results
}
`)
}

func dataSourceLifecycleConfig() string {
	return lifecycleProviderConfig(`data "azurecaf_name" "zero" {
  name          = "app"
  resource_type = "azurerm_resource_group"
  random_length = 12
  random_seed   = 0
}

data "azurecaf_name" "seeded" {
  name          = "app"
  resource_type = "azurerm_resource_group"
  random_length = 12
  random_seed   = 42
}

data "azurecaf_name" "unseeded" {
  name          = "app"
  resource_type = "azurerm_resource_group"
  random_length = 12
}

output "zero_result" {
  value = data.azurecaf_name.zero.result
}

output "seeded_result" {
  value = data.azurecaf_name.seeded.result
}

output "unseeded_result" {
  value = data.azurecaf_name.unseeded.result
}
`)
}

func importLifecycleConfig() string {
	return lifecycleProviderConfig(`resource "azurecaf_name" "imported" {
  name          = "stmyexistingapp"
  resource_type = "azurerm_storage_account"
  clean_input   = false
  passthrough   = true
  random_length = 0
  separator     = ""
  use_slug      = false
}

output "result" {
  value = azurecaf_name.imported.result
}
`)
}
