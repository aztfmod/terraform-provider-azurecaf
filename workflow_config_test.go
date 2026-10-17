package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestLifecycleWorkflowsUseNativeTerraformExitCodes(t *testing.T) {
	t.Parallel()

	versionPattern := regexp.MustCompile(`terraform_version:\s*["']?~>\s*1\.15["']?`)
	wrapperPattern := regexp.MustCompile(`terraform_wrapper:\s*false`)
	workflowPaths := []string{
		".github/workflows/e2e.yml",
		".github/workflows/go.yml",
		".github/workflows/release-validation.md",
		".github/workflows/release-validation.lock.yml",
	}

	for _, path := range workflowPaths {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read workflow: %v", err)
			}
			if !versionPattern.Match(content) {
				t.Fatalf("%s must use Terraform ~> 1.15 for lifecycle coverage", path)
			}
			if !wrapperPattern.Match(content) {
				t.Fatalf("%s must disable the Terraform wrapper to preserve detailed exit codes", path)
			}
		})
	}
}

func TestWeeklyAzureSyncExcludesSupersededAppServicePlan(t *testing.T) {
	t.Parallel()

	content, err := os.ReadFile("completness/non_nameable_resources.json")
	if err != nil {
		t.Fatalf("read non-nameable resource exclusions: %v", err)
	}
	var exclusions struct {
		DeprecatedInV4 []string `json:"deprecated_in_v4"`
	}
	if err := json.Unmarshal(content, &exclusions); err != nil {
		t.Fatalf("decode non-nameable resource exclusions: %v", err)
	}
	for _, resourceType := range exclusions.DeprecatedInV4 {
		if resourceType == "azurerm_app_service_plan" {
			return
		}
	}
	t.Fatal("azurerm_app_service_plan must remain excluded in favor of azurerm_service_plan")
}

func TestWeeklyAzureSyncLockIncludesSchemaFiltering(t *testing.T) {
	t.Parallel()

	workflowPaths := []string{
		".github/workflows/weekly-azure-sync.md",
		".github/workflows/weekly-azure-sync.lock.yml",
	}
	requiredMarkers := []string{
		"provider_schema.json",
		"SCHEMA_CHECK=ok",
		"AZURERM_TOTAL=$FILTERED_COUNT",
		"terraform_wrapper: false",
		"r not in static_excl",
	}

	for _, path := range workflowPaths {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read workflow: %v", err)
			}
			for _, marker := range requiredMarkers {
				if !strings.Contains(string(content), marker) {
					t.Fatalf("%s is missing generated schema-filtering marker %q", path, marker)
				}
			}
		})
	}
}
