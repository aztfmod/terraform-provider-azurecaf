package azurecaf

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestAcc_DataSourcesIntegration tests integration between data sources and resources
// This test uses direct provider schema testing to avoid Terraform CLI dependency
func TestAcc_DataSourcesIntegration(t *testing.T) {
	provider := Provider()
	
	// Set environment variable for testing
	os.Setenv("TEST_ENV_VAR", "test-env-value")
	defer os.Unsetenv("TEST_ENV_VAR")

	// Test environment variable data source
	t.Run("EnvironmentVariableDataSource", func(t *testing.T) {
		envVarDataSource := provider.DataSourcesMap["azurecaf_environment_variable"]
		if envVarDataSource == nil {
			t.Fatal("azurecaf_environment_variable data source not found")
		}

		// Create ResourceData for the data source
		envVarData := schema.TestResourceDataRaw(t, envVarDataSource.Schema, map[string]interface{}{
			"name": "TEST_ENV_VAR",
		})

		// Execute the read function
		diags := envVarDataSource.ReadContext(context.Background(), envVarData, nil)
		if diags.HasError() {
			t.Fatalf("Failed to read environment variable: %v", diags)
		}

		// Check the result
		value := envVarData.Get("value").(string)
		if value != "test-env-value" {
			t.Errorf("Expected value 'test-env-value', got '%s'", value)
		}
	})

	// Test name data source basic functionality
	t.Run("NameDataSourceBasic", func(t *testing.T) {
		nameDataSource := provider.DataSourcesMap["azurecaf_name"]
		if nameDataSource == nil {
			t.Fatal("azurecaf_name data source not found")
		}

		// Create ResourceData for the data source
		nameData := schema.TestResourceDataRaw(t, nameDataSource.Schema, map[string]interface{}{
			"name":          "myapp",
			"resource_type": "azurerm_app_service",
			"random_length": 5,
		})

		// Execute the read function
		diags := nameDataSource.ReadContext(context.Background(), nameData, nil)
		if diags.HasError() {
			t.Fatalf("Failed to read name data source: %v", diags)
		}

		// Check that result is set and contains expected content
		result := nameData.Get("result").(string)
		if result == "" {
			t.Error("Expected non-empty result")
		}
		if len(result) < 5 { // Should have at least the base name
			t.Errorf("Expected result length >= 5, got %d", len(result))
		}
	})

	// Test name data source with prefixes
	t.Run("NameDataSourceWithPrefixes", func(t *testing.T) {
		nameDataSource := provider.DataSourcesMap["azurecaf_name"]
		
		nameData := schema.TestResourceDataRaw(t, nameDataSource.Schema, map[string]interface{}{
			"name":          "storage-data",
			"prefixes":      []interface{}{"dev"},
			"resource_type": "azurerm_storage_account",
			"use_slug":      false,
			"clean_input":   true,
			"separator":     "-",
		})

		diags := nameDataSource.ReadContext(context.Background(), nameData, nil)
		if diags.HasError() {
			t.Fatalf("Failed to read name data source with prefixes: %v", diags)
		}

		result := nameData.Get("result").(string)
		// For storage account with clean_input=true, should contain "devstoragedata"
		expectedResult := "devstoragedata"
		if result != expectedResult {
			t.Errorf("Expected result '%s', got '%s'", expectedResult, result)
		}
	})

	// Test integration between environment variable and name data source
	t.Run("IntegrationEnvVarAndName", func(t *testing.T) {
		envVarDataSource := provider.DataSourcesMap["azurecaf_environment_variable"]
		nameDataSource := provider.DataSourcesMap["azurecaf_name"]
		
		// First get the environment variable
		envVarData := schema.TestResourceDataRaw(t, envVarDataSource.Schema, map[string]interface{}{
			"name": "TEST_ENV_VAR",
		})
		
		diags := envVarDataSource.ReadContext(context.Background(), envVarData, nil)
		if diags.HasError() {
			t.Fatalf("Failed to read environment variable: %v", diags)
		}
		
		envValue := envVarData.Get("value").(string)
		
		// Use the environment variable value in the name data source
		nameData := schema.TestResourceDataRaw(t, nameDataSource.Schema, map[string]interface{}{
			"name":          envValue + "-resource",
			"resource_type": "azurerm_resource_group",
		})
		
		diags = nameDataSource.ReadContext(context.Background(), nameData, nil)
		if diags.HasError() {
			t.Fatalf("Failed to read name data source with env var: %v", diags)
		}
		
		result := nameData.Get("result").(string)
		if result == "" {
			t.Error("Expected non-empty result from integrated test")
		}
		// Should contain the environment variable value
		if !strings.Contains(result, "test-env-value") {
			t.Errorf("Expected result to contain 'test-env-value', got '%s'", result)
		}
	})

	t.Log("Data sources integration tests completed successfully")
}

// Configuration for data sources integration test
const testAccDataSourcesIntegrationConfig = `
# Test environment variable
data "azurecaf_environment_variable" "test_env" {
  name = "TEST_ENV_VAR"
}

# Basic name data source
data "azurecaf_name" "simple" {
  name          = "myapp"
  resource_type = "azurerm_app_service"
  random_length = 5
}

# Name data source with prefixes and suffixes
data "azurecaf_name" "with_prefixes" {
  name          = "storage-data"
  prefixes      = ["dev"]
  resource_type = "azurerm_storage_account"
  use_slug      = false
  clean_input   = true
  separator     = "-"
}

# Name data source using environment variable
data "azurecaf_name" "with_env_var" {
  name          = "${data.azurecaf_environment_variable.test_env.value}-resource"
  resource_type = "azurerm_resource_group"
}

# Naming convention resource that uses data source outputs
resource "azurecaf_naming_convention" "combined" {
  name          = "${data.azurecaf_name.with_prefixes.result}-combined"
  resource_type = "rg"
  convention    = "random"
}
`
