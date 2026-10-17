package main

import (
	"os"
	"regexp"
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
