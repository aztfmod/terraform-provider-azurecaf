package framework

import (
	"fmt"
	"os/exec"
	"strings"
)

// findTerraformBinary safely locates the terraform binary and validates it
func findTerraformBinary() (string, error) {
	path, err := exec.LookPath("terraform")
	if err != nil {
		return "", fmt.Errorf("terraform binary not found in PATH: %w", err)
	}
	
	// Verify it's actually terraform by checking version
	cmd := exec.Command(path, "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to verify terraform binary at %s: %w", path, err)
	}
	
	if !strings.Contains(string(output), "Terraform") {
		return "", fmt.Errorf("binary at %s doesn't appear to be terraform", path)
	}
	
	return path, nil
}

// findMakeBinary safely locates the make binary
func findMakeBinary() (string, error) {
	path, err := exec.LookPath("make")
	if err != nil {
		return "", fmt.Errorf("make binary not found in PATH: %w", err)
	}
	return path, nil
}
