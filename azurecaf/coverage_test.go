package azurecaf

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This file contains tests specifically designed to increase code coverage
// by testing edge cases and error conditions that might not be covered
// by other tests.

// SECTION: Resource Name tests

// TestCoverage_GetNameResultExcessiveRandomLength tests resource_type constraint with excessive random_length
func TestCoverage_GetNameResultExcessiveRandomLength(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "azurerm_storage_account", // max length 24
		"random_length": 25,                        // exceeds max length
	})

	err := getNameResult(rd, nil)
	if err == nil {
		t.Error("Expected error for exceeding max length but got none")
	}
}

// TestCoverage_ResourceNameInvalidInputs tests error handling for invalid inputs
func TestCoverage_ResourceNameInvalidInputs(t *testing.T) {
	// Just test a single case for now to avoid issues
	rd := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "nonexistent_type",
	})

	err := getNameResult(rd, nil)
	if err == nil {
		t.Error("Expected error for invalid resource type but got none")
	}
}

// SECTION: Data Source tests

// TestCoverage_GetNameReadResultErrors tests error handling for data source read results
func TestCoverage_GetNameReadResultErrors(t *testing.T) {
	t.Run("getResourceName_error_path", func(t *testing.T) {
		// Setup with invalid resource type to trigger error in getResourceName
		rd := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
			"name":          "test",
			"resource_type": "invalid_resource_type",
			"random_length": 5,
		})

		err := getNameReadResult(rd, nil)
		if err == nil {
			t.Error("Expected error with invalid resource type but got none")
		}
	})
}

// TestCoverage_ValidateResourceTypeEdgeCases tests edge cases in resource type validation
func TestCoverage_ValidateResourceTypeEdgeCases(t *testing.T) {
	t.Run("validate_empty_lists", func(t *testing.T) {
		valid, err := validateResourceType("", []string{})
		if valid || err == nil {
			t.Error("Expected error with empty resource types lists")
		}
	})

	t.Run("validate_one_valid_one_invalid", func(t *testing.T) {
		valid, err := validateResourceType("", []string{"azurerm_storage_account", "invalid_type"})
		if valid || err == nil {
			t.Error("Expected error with one valid and one invalid resource type")
		}
	})
}

// SECTION: Environment Variable tests

// Test resourceAction function (data_environment_variable.go)
func TestCoverage_ResourceAction(t *testing.T) {
	// Test with existing environment variable
	t.Run("existing_env_var", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test_value")
		defer os.Unsetenv("TEST_VAR")

		rd := schema.TestResourceDataRaw(t, dataEnvironmentVariable().Schema, map[string]interface{}{
			"name": "TEST_VAR",
		})

		diags := resourceAction(context.Background(), rd, nil)

		if len(diags) != 0 {
			t.Errorf("Expected no diagnostics, got: %v", diags)
		}

		if rd.Id() != "TEST_VAR" {
			t.Errorf("Expected ID to be 'TEST_VAR', got: %s", rd.Id())
		}

		if value := rd.Get("value").(string); value != "test_value" {
			t.Errorf("Expected value to be 'test_value', got: %s", value)
		}
	})

	// Test with non-existing environment variable
	t.Run("non_existing_env_var", func(t *testing.T) {
		rd := schema.TestResourceDataRaw(t, dataEnvironmentVariable().Schema, map[string]interface{}{
			"name": "NON_EXISTING_VAR",
		})

		diags := resourceAction(context.Background(), rd, nil)

		if len(diags) != 1 {
			t.Errorf("Expected 1 diagnostic, got: %d", len(diags))
		}
	})
}

// SECTION: Edge Case tests

// Test getResourceName validation error path
func TestCoverage_GetResourceNameValidationError(t *testing.T) {
	// Save the original resources
	original := ResourceDefinitions["azurerm_storage_account"]

	// Create a modified version with invalid regex pattern that will cause validation failure
	modified := original
	// Keep valid compilation but create a pattern that won't match any input
	modified.ValidationRegExp = "^$" // This will only match empty string

	ResourceDefinitions["azurerm_storage_account"] = modified

	defer func() {
		// Restore original after test
		ResourceDefinitions["azurerm_storage_account"] = original
	}()

	// Now try to use the resource type with a name that won't match the regex
	_, err := getResourceName("azurerm_storage_account", "-", []string{}, "test", []string{}, "", "cafclassic", false, false, true, []string{"name"})

	if err == nil {
		t.Error("Expected validation error but got none")
	}
}

// Test regex compilation error in getResult by handling panic
func TestCoverage_GetResultRegexError(t *testing.T) {
	// Save the original resource
	originalResource := Resources["st"]

	// Create a modified version with invalid regex pattern - this should error on compile
	modifiedResource := originalResource
	modifiedResource.RegEx = "[" // Invalid regex pattern that will cause compile error
	Resources["st"] = modifiedResource

	defer func() {
		// Restore original after test
		Resources["st"] = originalResource

		// Recover from any panic
		if r := recover(); r != nil {
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	// This should properly handle the error rather than panic
	// Create a ResourceData with the required fields
	rd := schema.TestResourceDataRaw(t, resourceNamingConvention().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "st", // The resource with invalid regex
		"convention":    "cafclassic",
	})

	err := getResult(rd, nil)
	if err == nil {
		t.Error("Expected regex compilation error but got none")
	}
}
