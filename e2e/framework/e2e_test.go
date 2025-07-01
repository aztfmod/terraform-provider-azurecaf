// Package framework provides the core end-to-end testing framework for the terraform-provider-azurecaf.
//
// This framework implements comprehensive e2e tests that cover the complete workflow:
// 1. Build the provider from source
// 2. Configure Terraform to use the built provider
// 3. Generate and validate Azure CAF-compliant resource names
// 4. Test integration with azurerm provider using mock testing
// 5. Validate deployment scenarios
package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// Setup prepares the e2e test environment
func (suite *E2ETestSuite) Setup(t *testing.T) error {
	// Create working directory
	if err := os.MkdirAll(suite.WorkingDir, 0750); err != nil {
		return fmt.Errorf("failed to create working directory: %w", err)
	}

	// Build the provider
	if err := suite.buildProvider(t); err != nil {
		return fmt.Errorf("failed to build provider: %w", err)
	}

	// Setup Terraform with dev_overrides (no need to install provider locally)
	if err := suite.setupTerraform(t); err != nil {
		return fmt.Errorf("failed to setup Terraform: %w", err)
	}

	t.Logf("E2E test suite setup completed in: %s", suite.WorkingDir)
	return nil
}

// buildProvider compiles the terraform-provider-azurecaf from source
func (suite *E2ETestSuite) buildProvider(t *testing.T) error {
	t.Log("Building terraform-provider-azurecaf from source...")
	
	// Get the project root directory
	projectRoot, err := suite.getProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	// Set provider binary path
	suite.ProviderBinaryPath = filepath.Join(projectRoot, "terraform-provider-azurecaf")
	
	// Build the provider using make
	makePath, err := findMakeBinary()
	if err != nil {
		return fmt.Errorf("failed to find make binary: %w", err)
	}
	cmd := exec.Command(makePath, "build")
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(),
		"CHECKPOINT_DISABLE=1",
		"TF_IN_AUTOMATION=1",
		"TF_CLI_ARGS_init=-upgrade=false",
	)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %w\nOutput: %s", err, output)
	}

	// Verify the binary exists
	if _, err := os.Stat(suite.ProviderBinaryPath); os.IsNotExist(err) {
		return fmt.Errorf("provider binary not found at: %s", suite.ProviderBinaryPath)
	}

	t.Logf("Successfully built provider: %s", suite.ProviderBinaryPath)
	return nil
}

// setupTerraform configures Terraform for local provider testing
func (suite *E2ETestSuite) setupTerraform(t *testing.T) error {
	t.Log("Setting up Terraform for local provider testing...")
	
	// Find terraform executable
	terraformPath, err := exec.LookPath("terraform")
	if err != nil {
		return fmt.Errorf("terraform executable not found: %w", err)
	}
	suite.TerraformPath = terraformPath

	// Create provider override configuration
	if err := suite.createProviderOverride(t); err != nil {
		return fmt.Errorf("failed to create provider override: %w", err)
	}

	t.Logf("Terraform configured with executable: %s", suite.TerraformPath)
	return nil
}

// installProviderLocally installs the built provider in the local plugin directory
func (suite *E2ETestSuite) installProviderLocally(t *testing.T) error {
	t.Log("Installing provider in local plugin directory...")
	
	// Get OS and architecture
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	
	// Create local plugin directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	pluginDir := filepath.Join(homeDir, ".terraform.d", "plugins", suite.ProviderSource, "1.0.0", fmt.Sprintf("%s_%s", goos, goarch))
	if err := os.MkdirAll(pluginDir, 0750); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Copy provider binary
	destPath := filepath.Join(pluginDir, "terraform-provider-azurecaf")
	if err := suite.copyFile(suite.ProviderBinaryPath, destPath); err != nil {
		return fmt.Errorf("failed to copy provider binary: %w", err)
	}

	// Make executable
	if err := os.Chmod(destPath, 0750); err != nil {
		return fmt.Errorf("failed to make provider executable: %w", err)
	}

	t.Logf("Provider installed in: %s", pluginDir)
	return nil
}

// createProviderOverride creates a Terraform configuration file for development overrides
func (suite *E2ETestSuite) createProviderOverride(t *testing.T) error {
	// Create a filesystem mirror structure for local provider testing
	mirrorDir := filepath.Join(suite.WorkingDir, "terraform-mirror")
	
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	
	// Create the provider directory structure
	providerDir := filepath.Join(mirrorDir, "aztfmod.com", "test", "azurecaf", "1.0.0", fmt.Sprintf("%s_%s", goos, goarch))
	if err := os.MkdirAll(providerDir, 0750); err != nil {
		return fmt.Errorf("failed to create provider directory: %w", err)
	}
	
	// Copy the provider binary to the mirror directory
	providerBinary := filepath.Join(providerDir, "terraform-provider-azurecaf")
	if err := suite.copyFile(suite.ProviderBinaryPath, providerBinary); err != nil {
		return fmt.Errorf("failed to copy provider binary: %w", err)
	}
	
	// Make the binary executable
	if err := os.Chmod(providerBinary, 0750); err != nil {
		return fmt.Errorf("failed to make provider binary executable: %w", err)
	}
	
	// Create provider installation configuration with filesystem_mirror
	overrideConfig := fmt.Sprintf(`provider_installation {
  filesystem_mirror {
    path    = "%s"
    include = ["aztfmod.com/test/*"]
  }
  direct {
    exclude = ["aztfmod.com/test/*"]
  }
}`, mirrorDir)

	configPath := filepath.Join(suite.WorkingDir, "terraform.rc")
	if err := ioutil.WriteFile(configPath, []byte(overrideConfig), 0644); err != nil {
		return fmt.Errorf("failed to write provider override config: %w", err)
	}

	t.Logf("Provider override configuration created: %s", configPath)
	t.Logf("Provider binary copied to mirror: %s", providerBinary)
	
	return nil
}

// copyFile copies a file from src to dst
func (suite *E2ETestSuite) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// RunScenario executes a single test scenario
func (suite *E2ETestSuite) RunScenario(t *testing.T, scenario TestScenario) error {
	t.Logf("Running scenario: %s", scenario.Name)
	
	// Create scenario directory
	scenarioDir := filepath.Join(suite.WorkingDir, "scenarios", scenario.Name)
	if err := os.MkdirAll(scenarioDir, 0750); err != nil {
		return fmt.Errorf("failed to create scenario directory: %w", err)
	}

	// Write Terraform configuration
	configPath := filepath.Join(scenarioDir, "main.tf")
	if err := ioutil.WriteFile(configPath, []byte(scenario.TerraformConfig), 0644); err != nil {
		return fmt.Errorf("failed to write terraform configuration: %w", err)
	}

	// Create Terraform executor
	tf, err := tfexec.NewTerraform(scenarioDir, suite.TerraformPath)
	if err != nil {
		return fmt.Errorf("failed to create terraform executor: %w", err)
	}

	// Set environment variables
	tf.SetEnv(map[string]string{
		"TF_CLI_CONFIG_FILE":   filepath.Join(suite.WorkingDir, "terraform.rc"),
		"CHECKPOINT_DISABLE":   "1",
		"TF_IN_AUTOMATION":     "1",
		"TF_CLI_ARGS_init":     "-upgrade=false",
	})

	// Initialize Terraform
	ctx, cancel := context.WithTimeout(context.Background(), suite.Timeout)
	defer cancel()

	if err := tf.Init(ctx); err != nil {
		return fmt.Errorf("terraform init failed: %w", err)
	}

	// Plan
	planChanges, err := tf.Plan(ctx)
	if err != nil {
		if len(scenario.ExpectedErrors) > 0 {
			// Check if this is an expected error
			errorMatched := false
			for _, expectedError := range scenario.ExpectedErrors {
				if strings.Contains(err.Error(), expectedError) {
					errorMatched = true
					t.Logf("Expected error matched: %s", expectedError)
					break
				}
			}
			if !errorMatched {
				return fmt.Errorf("unexpected error during plan: %w", err)
			}
			return nil // Expected error occurred, test passed
		}
		return fmt.Errorf("terraform plan failed: %w", err)
	}

	t.Logf("Plan completed, changes: %t", planChanges)

	// Apply (if no errors expected)
	if len(scenario.ExpectedErrors) == 0 {
		if err := tf.Apply(ctx); err != nil {
			return fmt.Errorf("terraform apply failed: %w", err)
		}

		// Get outputs
		outputs, err := tf.Output(ctx)
		if err != nil {
			return fmt.Errorf("failed to get terraform outputs: %w", err)
		}

		// Validate outputs
		if err := suite.validateOutputs(t, outputs, scenario.ExpectedOutputs); err != nil {
			return fmt.Errorf("output validation failed: %w", err)
		}

		// Run additional validation checks
		if err := suite.runValidationChecks(t, scenario.ValidationChecks, outputs); err != nil {
			return fmt.Errorf("validation checks failed: %w", err)
		}
	}

	t.Logf("Scenario completed successfully: %s", scenario.Name)
	return nil
}

// validateOutputs validates terraform outputs against expected values
func (suite *E2ETestSuite) validateOutputs(t *testing.T, outputs map[string]tfexec.OutputMeta, expected map[string]interface{}) error {
	for key, expectedValue := range expected {
		output, exists := outputs[key]
		if !exists {
			return fmt.Errorf("expected output %s not found", key)
		}

		var actualValue interface{}
		if err := json.Unmarshal(output.Value, &actualValue); err != nil {
			return fmt.Errorf("failed to unmarshal output %s: %w", key, err)
		}

		if !suite.compareValues(actualValue, expectedValue) {
			return fmt.Errorf("output %s mismatch: expected %v, got %v", key, expectedValue, actualValue)
		}

		t.Logf("Output validation passed: %s = %v", key, actualValue)
	}

	return nil
}

// runValidationChecks executes additional validation checks
func (suite *E2ETestSuite) runValidationChecks(t *testing.T, checks []ValidationCheck, outputs map[string]tfexec.OutputMeta) error {
	for _, check := range checks {
		switch check.Type {
		case "naming_convention":
			if err := suite.validateNamingConvention(t, check, outputs); err != nil {
				return fmt.Errorf("naming convention validation failed: %w", err)
			}
		case "azure_compliance":
			if err := suite.validateAzureCompliance(t, check, outputs); err != nil {
				return fmt.Errorf("azure compliance validation failed: %w", err)
			}
		default:
			t.Logf("Unknown validation check type: %s", check.Type)
		}
	}

	return nil
}

// validateNamingConvention validates that generated names follow Azure naming conventions
func (suite *E2ETestSuite) validateNamingConvention(t *testing.T, check ValidationCheck, outputs map[string]tfexec.OutputMeta) error {
	output, exists := outputs[check.Target]
	if !exists {
		return fmt.Errorf("output %s not found for naming convention validation", check.Target)
	}

	var name string
	if err := json.Unmarshal(output.Value, &name); err != nil {
		return fmt.Errorf("failed to unmarshal output %s: %w", check.Target, err)
	}

	// Perform naming convention validations based on expected criteria
	if expected, ok := check.Expected.(map[string]interface{}); ok {
		if maxLength, exists := expected["max_length"]; exists {
			if len(name) > int(maxLength.(float64)) {
				return fmt.Errorf("name %s exceeds max length %d", name, int(maxLength.(float64)))
			}
		}

		if mustContain, exists := expected["must_contain"]; exists {
			if !strings.Contains(name, mustContain.(string)) {
				return fmt.Errorf("name %s must contain %s", name, mustContain.(string))
			}
		}

		if pattern, exists := expected["pattern"]; exists {
			// Could add regex validation here
			t.Logf("Pattern validation for %s: %s (placeholder)", name, pattern.(string))
		}
	}

	t.Logf("Naming convention validation passed: %s", name)
	return nil
}

// validateAzureCompliance validates that names comply with Azure resource requirements
func (suite *E2ETestSuite) validateAzureCompliance(t *testing.T, check ValidationCheck, outputs map[string]tfexec.OutputMeta) error {
	output, exists := outputs[check.Target]
	if !exists {
		return fmt.Errorf("output %s not found for azure compliance validation", check.Target)
	}

	var name string
	if err := json.Unmarshal(output.Value, &name); err != nil {
		return fmt.Errorf("failed to unmarshal output %s: %w", check.Target, err)
	}

	// Basic Azure compliance checks
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("name %s cannot start or end with hyphen", name)
	}

	if strings.Contains(name, "--") {
		return fmt.Errorf("name %s cannot contain consecutive hyphens", name)
	}

	t.Logf("Azure compliance validation passed: %s", name)
	return nil
}

// Cleanup removes temporary test files and directories
func (suite *E2ETestSuite) Cleanup(t *testing.T) {
	if !suite.CleanupEnabled {
		t.Logf("Cleanup disabled, preserving test directory: %s", suite.WorkingDir)
		return
	}

	if err := os.RemoveAll(suite.WorkingDir); err != nil {
		t.Logf("Warning: failed to cleanup test directory %s: %v", suite.WorkingDir, err)
	} else {
		t.Logf("Cleaned up test directory: %s", suite.WorkingDir)
	}
}

// Helper functions

func (suite *E2ETestSuite) getProjectRoot() (string, error) {
	// Start from current working directory and walk up to find main.go and Makefile
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Look for main.go and Makefile which indicate the terraform provider root
		mainGoExists := false
		makefileExists := false
		
		if _, err := os.Stat(filepath.Join(dir, "main.go")); err == nil {
			mainGoExists = true
		}
		if _, err := os.Stat(filepath.Join(dir, "Makefile")); err == nil {
			makefileExists = true
		}
		
		if mainGoExists && makefileExists {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("project root not found (looking for main.go and Makefile)")
}

func (suite *E2ETestSuite) compareValues(actual, expected interface{}) bool {
	// Simple comparison - could be enhanced for more complex types
	actualJSON, _ := json.Marshal(actual)
	expectedJSON, _ := json.Marshal(expected)
	return string(actualJSON) == string(expectedJSON)
}
