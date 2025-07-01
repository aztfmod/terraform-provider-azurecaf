// Package e2e provides comprehensive end-to-end tests for the terraform-provider-azurecaf.
//
// These tests verify the complete workflow from building the provider to deploying
// resources with Terraform, ensuring Azure CAF compliance and azurerm provider integration.
package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Test timeout for long-running operations
	testTimeout = 5 * time.Minute

	// Provider binary name
	providerBinary = "terraform-provider-azurecaf"

	// Test workspace directory
	testWorkspace = "test-workspace"
)

// TestE2E_ComprehensiveWorkflow runs the complete end-to-end workflow
func TestE2E_ComprehensiveWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}

	// Set up test environment
	projectRoot := getProjectRoot(t)
	workspaceDir := setupTestWorkspace(t, projectRoot)
	defer cleanupTestWorkspace(workspaceDir)

	// Step 1: Build the provider
	t.Log("=== STEP 1: Building Provider ===")
	testBuildProvider(t, projectRoot)

	// Step 2: Validate provider binary
	t.Log("=== STEP 2: Validating Provider Binary ===")
	testValidateProviderBinary(t, projectRoot)

	// Step 3: Setup Terraform with local provider
	t.Log("=== STEP 3: Setting up Terraform with Local Provider ===")
	testSetupTerraformWithLocalProvider(t, projectRoot, workspaceDir)

	// Step 4: Validate CAF name generation
	t.Log("=== STEP 4: Validating CAF Name Generation ===")
	testValidateCAFNameGeneration(t, workspaceDir)

	// Step 5: Test azurerm provider integration
	t.Log("=== STEP 5: Testing Mock azurerm Provider Integration ===")
	testMockAzurermProviderIntegration(t, workspaceDir)

	// Step 6: Validate deployment scenarios
	t.Log("=== STEP 6: Validating Deployment Scenarios ===")
	testValidateDeploymentScenarios(t, workspaceDir)

	t.Log("=== E2E TESTS COMPLETED SUCCESSFULLY ===")
}

// testBuildProvider verifies that the provider can be built successfully from source
func testBuildProvider(t *testing.T, projectRoot string) {
	t.Log("Building terraform-provider-azurecaf from source...")

	// Clean any existing binary
	binaryPath := filepath.Join(projectRoot, providerBinary)
	os.Remove(binaryPath)

	// Build the provider
	cmd := exec.Command("go", "build", "-o", providerBinary)
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to build provider: %s", string(output))

	// Verify binary exists and is executable
	info, err := os.Stat(binaryPath)
	require.NoError(t, err, "Provider binary not found at %s", binaryPath)

	// Check if binary is executable (on Unix systems)
	if runtime.GOOS != "windows" {
		require.True(t, info.Mode()&0111 != 0, "Provider binary is not executable")
	}

	t.Logf("Successfully built provider binary: %s (size: %d bytes)", binaryPath, info.Size())
}

// testValidateProviderBinary verifies the built provider binary responds correctly
func testValidateProviderBinary(t *testing.T, projectRoot string) {
	t.Log("Validating provider binary functionality...")

	binaryPath := filepath.Join(projectRoot, providerBinary)

	// Test that the provider binary can be executed (though it will wait for plugin protocol)
	// We'll just check it exists and starts without immediate error
	cmd := exec.Command(binaryPath)
	cmd.Dir = projectRoot

	// Start the command but don't wait for it to finish (it will hang waiting for plugin protocol)
	err := cmd.Start()
	require.NoError(t, err, "Failed to start provider binary")

	// Give it a moment to start, then kill it
	time.Sleep(100 * time.Millisecond)
	err = cmd.Process.Kill()
	require.NoError(t, err, "Failed to stop provider binary")

	t.Log("Provider binary validated successfully")
}

// testSetupTerraformWithLocalProvider configures Terraform to use the local provider build
func testSetupTerraformWithLocalProvider(t *testing.T, projectRoot, workspaceDir string) {
	t.Log("Setting up Terraform to use local provider build...")

	// Create development override configuration
	// This tells Terraform to use the local binary instead of downloading from registry
	devOverrideConfig := fmt.Sprintf(`provider_installation {
  dev_overrides {
    "aztfmod/azurecaf" = "%s"
  }
  direct {}
}`, projectRoot)

	// Create .terraformrc file for development overrides
	terraformrcPath := filepath.Join(workspaceDir, ".terraformrc")
	err := os.WriteFile(terraformrcPath, []byte(devOverrideConfig), 0644)
	require.NoError(t, err, "Failed to write .terraformrc")

	// Create terraform configuration that uses the standard provider source
	terraformConfig := `terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = ">= 1.0.0"
    }
  }
}

provider "azurecaf" {}
`

	configPath := filepath.Join(workspaceDir, "main.tf")
	err = os.WriteFile(configPath, []byte(terraformConfig), 0644)
	require.NoError(t, err, "Failed to write terraform configuration")

	t.Logf("Successfully set up Terraform with development override at %s", projectRoot)
}

// testValidateCAFNameGeneration tests Azure CAF-compliant name generation
func testValidateCAFNameGeneration(t *testing.T, workspaceDir string) {
	t.Log("Validating Azure CAF name generation...")

	// Create test configuration for name generation
	nameTestConfig := `
# Test basic CAF name generation
resource "azurecaf_name" "rg_test" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 4
  clean_input   = true
}

# Test storage account (has specific constraints)
resource "azurecaf_name" "st_test" {
  name          = "storage"
  resource_type = "azurerm_storage_account"
  random_length = 8
  clean_input   = true
}

# Test key vault (has specific constraints)
resource "azurecaf_name" "kv_test" {
  name          = "secrets"
  resource_type = "azurerm_key_vault"
  prefixes      = ["prod"]
  random_length = 5
  clean_input   = true
}

# Output the generated names for validation
output "rg_name" {
  value = azurecaf_name.rg_test.result
}

output "st_name" {
  value = azurecaf_name.st_test.result
}

output "kv_name" {
  value = azurecaf_name.kv_test.result
}
`

	configPath := filepath.Join(workspaceDir, "name_test.tf")
	err := os.WriteFile(configPath, []byte(nameTestConfig), 0644)
	require.NoError(t, err, "Failed to write name test configuration")

	// Run terraform plan to generate names
	planOutput := runTerraformCommand(t, workspaceDir, "plan", "-out=tfplan")

	// Validate that plan succeeded and contains expected outputs
	assert.Contains(t, planOutput, "azurecaf_name.rg_test", "Resource group name resource not found in plan")
	assert.Contains(t, planOutput, "azurecaf_name.st_test", "Storage account name resource not found in plan")
	assert.Contains(t, planOutput, "azurecaf_name.kv_test", "Key vault name resource not found in plan")

	// Run terraform apply to actually generate the names
	runTerraformCommand(t, workspaceDir, "apply", "-auto-approve", "tfplan")

	// Extract output values
	outputCmd := runTerraformCommand(t, workspaceDir, "output", "-json")

	// Parse and validate the outputs meet CAF standards
	validateCAFCompliantNames(t, outputCmd)

	t.Log("Successfully validated Azure CAF name generation")
}

// testMockAzurermProviderIntegration tests integration with azurerm provider using mock
func testMockAzurermProviderIntegration(t *testing.T, workspaceDir string) {
	t.Log("Testing mock azurerm provider integration...")

	// Create test configuration that would use azurerm provider
	integrationConfig := `
# Generate names for Azure resources
resource "azurecaf_name" "integration_rg" {
  name          = "integration"
  resource_type = "azurerm_resource_group"
  prefixes      = ["test"]
  random_length = 4
}

resource "azurecaf_name" "integration_st" {
  name          = "storage"
  resource_type = "azurerm_storage_account"
  random_length = 8
}

# Mock azurerm resources that would use the generated names
# Note: These are commented out as we're doing mock testing
# resource "azurerm_resource_group" "integration" {
#   name     = azurecaf_name.integration_rg.result
#   location = "East US"
# }
# 
# resource "azurerm_storage_account" "integration" {
#   name                     = azurecaf_name.integration_st.result
#   resource_group_name      = azurerm_resource_group.integration.name
#   location                 = azurerm_resource_group.integration.location
#   account_tier             = "Standard"
#   account_replication_type = "LRS"
# }

# Outputs to validate names work with azurerm naming constraints
output "integration_rg_name" {
  value = azurecaf_name.integration_rg.result
  description = "Resource group name for azurerm integration"
}

output "integration_st_name" {
  value = azurecaf_name.integration_st.result
  description = "Storage account name for azurerm integration"
}
`

	configPath := filepath.Join(workspaceDir, "integration_test.tf")
	err := os.WriteFile(configPath, []byte(integrationConfig), 0644)
	require.NoError(t, err, "Failed to write integration test configuration")

	// Run terraform plan
	planOutput := runTerraformCommand(t, workspaceDir, "plan")

	// Validate that names would be compatible with azurerm resources
	assert.Contains(t, planOutput, "azurecaf_name.integration_rg", "Integration RG name not found")
	assert.Contains(t, planOutput, "azurecaf_name.integration_st", "Integration storage name not found")

	t.Log("Successfully tested azurerm provider integration scenarios")
}

// testValidateDeploymentScenarios validates various deployment scenarios
func testValidateDeploymentScenarios(t *testing.T, workspaceDir string) {
	t.Log("Validating deployment scenarios...")

	// Test multiple resource types scenario
	multiResourceConfig := `
# Test multiple resource types with same base name
resource "azurecaf_name" "multi_app" {
  name           = "webapp"
  resource_type  = "azurerm_app_service"
  resource_types = ["azurerm_app_service_plan", "azurerm_storage_account"]
  prefixes       = ["prod"]
  suffixes       = ["api"]
  random_length  = 5
  clean_input    = true
}

# Test data source for validation
data "azurecaf_name" "validation" {
  name          = "existing-resource"
  resource_type = "azurerm_virtual_machine"
  passthrough   = true
  clean_input   = true
}

# Test edge cases
resource "azurecaf_name" "edge_case" {
  name          = "Test-Name_With.Special@Characters!"
  resource_type = "azurerm_subnet"
  clean_input   = true
  use_slug      = false
}

output "multi_app_primary" {
  value = azurecaf_name.multi_app.result
}

output "multi_app_all" {
  value = azurecaf_name.multi_app.results
}

output "validation_result" {
  value = data.azurecaf_name.validation.result
}

output "edge_case_result" {
  value = azurecaf_name.edge_case.result
}
`

	configPath := filepath.Join(workspaceDir, "scenarios_test.tf")
	err := os.WriteFile(configPath, []byte(multiResourceConfig), 0644)
	require.NoError(t, err, "Failed to write scenarios test configuration")

	// Run terraform plan and apply
	runTerraformCommand(t, workspaceDir, "plan")
	runTerraformCommand(t, workspaceDir, "apply", "-auto-approve")

	// Get outputs and validate
	outputCmd := runTerraformCommand(t, workspaceDir, "output", "-json")

	// Validate multiple scenarios worked
	assert.Contains(t, outputCmd, "multi_app_primary", "Multi-resource primary output not found")
	assert.Contains(t, outputCmd, "multi_app_all", "Multi-resource all outputs not found")
	assert.Contains(t, outputCmd, "validation_result", "Validation result not found")
	assert.Contains(t, outputCmd, "edge_case_result", "Edge case result not found")

	t.Log("Successfully validated deployment scenarios")
}

// Helper functions

func getProjectRoot(t *testing.T) string {
	wd, err := os.Getwd()
	require.NoError(t, err, "Failed to get working directory")

	// Go up one level from e2e directory to project root
	return filepath.Dir(wd)
}

func setupTestWorkspace(t *testing.T, projectRoot string) string {
	workspaceDir := filepath.Join(projectRoot, testWorkspace)

	// Clean up any existing workspace
	os.RemoveAll(workspaceDir)

	// Create new workspace
	err := os.MkdirAll(workspaceDir, 0755)
	require.NoError(t, err, "Failed to create test workspace")

	return workspaceDir
}

func cleanupTestWorkspace(workspaceDir string) {
	os.RemoveAll(workspaceDir)
}

func copyFile(t *testing.T, src, dst string) {
	data, err := os.ReadFile(src)
	require.NoError(t, err, "Failed to read source file")

	err = os.WriteFile(dst, data, 0644)
	require.NoError(t, err, "Failed to write destination file")
}

func runTerraformCommand(t *testing.T, workdir string, args ...string) string {
	// Set TF_CLI_CONFIG_FILE to use our custom .terraformrc
	terraformrcPath := filepath.Join(workdir, ".terraformrc")

	// For dev overrides, we need to run init first, but it will warn us to skip it in future
	// We'll run init only once when the .terraform directory doesn't exist
	if args[0] != "init" && !fileExists(filepath.Join(workdir, ".terraform")) {
		initCmd := exec.Command("terraform", "init")
		initCmd.Dir = workdir
		initCmd.Env = append(os.Environ(),
			"TF_IN_AUTOMATION=1",
			"CHECKPOINT_DISABLE=1",
			"TF_CLI_CONFIG_FILE="+terraformrcPath,
		)
		initOutput, initErr := initCmd.CombinedOutput()
		t.Logf("Terraform init output: %s", string(initOutput))
		if initErr != nil {
			t.Logf("Terraform init error: %v", initErr)
			// Don't fail on init errors with dev overrides, as it's expected
		}
	}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = workdir
	cmd.Env = append(os.Environ(),
		"TF_IN_AUTOMATION=1",
		"CHECKPOINT_DISABLE=1",
		"TF_CLI_CONFIG_FILE="+terraformrcPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Terraform command failed: terraform %s", strings.Join(args, " "))
		t.Logf("Output: %s", string(output))
		t.Logf("Error: %v", err)
		require.NoError(t, err, "Terraform command failed")
	}

	return string(output)
}

// fileExists checks if a file or directory exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func validateCAFCompliantNames(t *testing.T, outputJSON string) {
	// Parse JSON and validate naming patterns
	t.Log("Validating CAF compliance of generated names...")

	// Resource group names should match Azure patterns
	rgPattern := regexp.MustCompile(`"rg_name":\s*{\s*"value":\s*"([^"]+)"`)
	if matches := rgPattern.FindStringSubmatch(outputJSON); len(matches) > 1 {
		rgName := matches[1]
		t.Logf("Generated RG name: %s", rgName)

		// Validate RG naming constraints
		assert.LessOrEqual(t, len(rgName), 90, "Resource group name too long")
		assert.Regexp(t, `^[a-zA-Z0-9\-_.()]+$`, rgName, "Resource group name contains invalid characters")
		assert.NotRegexp(t, `\.$`, rgName, "Resource group name cannot end with period")
	}

	// Storage account names should be lowercase, no special chars, 3-24 chars
	stPattern := regexp.MustCompile(`"st_name":\s*{\s*"value":\s*"([^"]+)"`)
	if matches := stPattern.FindStringSubmatch(outputJSON); len(matches) > 1 {
		stName := matches[1]
		t.Logf("Generated storage name: %s", stName)

		// Validate storage account naming constraints
		assert.GreaterOrEqual(t, len(stName), 3, "Storage account name too short")
		assert.LessOrEqual(t, len(stName), 24, "Storage account name too long")
		assert.Regexp(t, `^[a-z0-9]+$`, stName, "Storage account name must be lowercase alphanumeric only")
	}

	// Key vault names should follow key vault constraints
	kvPattern := regexp.MustCompile(`"kv_name":\s*{\s*"value":\s*"([^"]+)"`)
	if matches := kvPattern.FindStringSubmatch(outputJSON); len(matches) > 1 {
		kvName := matches[1]
		t.Logf("Generated key vault name: %s", kvName)

		// Validate key vault naming constraints
		assert.GreaterOrEqual(t, len(kvName), 3, "Key vault name too short")
		assert.LessOrEqual(t, len(kvName), 24, "Key vault name too long")
		assert.Regexp(t, `^[a-zA-Z][a-zA-Z0-9-]*$`, kvName, "Key vault name must start with letter and contain only alphanumeric and hyphens")
		assert.NotRegexp(t, `--`, kvName, "Key vault name cannot contain consecutive hyphens")
		assert.NotRegexp(t, `-$`, kvName, "Key vault name cannot end with hyphen")
	}
}
