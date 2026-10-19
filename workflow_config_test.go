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

func TestReleasePublishingRunsAllRequiredGates(t *testing.T) {
	t.Parallel()

	content, err := os.ReadFile(".github/workflows/go.yml")
	if err != nil {
		t.Fatalf("read release workflow: %v", err)
	}
	requiredMarkers := []string{
		`github.event_name == 'pull_request' || startsWith(github.ref, 'refs/tags/')`,
		`.github/releases/${GITHUB_REF_NAME}.md`,
		"--release-notes",
		`>> "$GITHUB_PATH"`,
	}
	for _, marker := range requiredMarkers {
		if !strings.Contains(string(content), marker) {
			t.Fatalf("go.yml is missing release gate marker %q", marker)
		}
	}
	if _, err := os.Stat(".github/releases/v2.0.0.md"); err != nil {
		t.Fatalf("v2.0.0 release notes must be committed for the tag workflow: %v", err)
	}
}

func TestCoverageGateAndReleaseReporting(t *testing.T) {
	t.Parallel()

	makefile, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatalf("read Makefile: %v", err)
	}
	for _, marker := range []string{
		"COVERAGE_THRESHOLD ?= 95",
		"go test -coverprofile=coverage.out ./azurecaf/...",
		"Provider coverage:",
		"coverage + 0 >= threshold + 0",
	} {
		if !strings.Contains(string(makefile), marker) {
			t.Fatalf("Makefile is missing coverage gate marker %q", marker)
		}
	}

	workflowMarkers := map[string][]string{
		".github/workflows/release-validation.md": {
			"/tmp/gh-aw/agent/release-validation",
			"Provider coverage:",
			"VALIDATION_EXIT",
			"create-issue",
		},
		".github/workflows/release-validation.lock.yml": {
			"/tmp/gh-aw/agent/release-validation",
			"Provider coverage:",
			"VALIDATION_EXIT",
			"create_issue",
		},
	}
	for path, markers := range workflowMarkers {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		for _, marker := range markers {
			if !strings.Contains(string(content), marker) {
				t.Fatalf("%s is missing release validation marker %q", path, marker)
			}
		}
	}

	labeler, err := os.ReadFile(".github/workflows/release-labeler.md")
	if err != nil {
		t.Fatalf("read release labeler workflow: %v", err)
	}
	for _, marker := range []string{
		"create-issue:",
		"GitHub Releases do not support comments",
		"For manual backfills",
		"If any label operation fails",
	} {
		if !strings.Contains(string(labeler), marker) {
			t.Fatalf("release-labeler.md is missing reporting marker %q", marker)
		}
	}
	if strings.Contains(string(labeler), "Post the report as a comment on the latest release") {
		t.Fatal("release-labeler.md must not instruct the agent to comment on a GitHub Release")
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
