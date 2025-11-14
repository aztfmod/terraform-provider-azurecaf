package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Test for getNameReadResult error handling.
func TestGetNameReadResultErrors(t *testing.T) {
	t.Run("getResourceName_error_path", func(t *testing.T) {
		// Setup with invalid resource type to trigger error in getResourceName
		rd := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
			"name":          "test",
			"resource_type": "invalid_resource_type",
			"random_length": 5,
		})

		err := getNameReadResult(rd)
		if err == nil {
			t.Error("Expected error with invalid resource type but got none")
		}
	})
}

// Test validateResourceType corner cases.
func TestValidateResourceTypeEdgeCases(t *testing.T) {
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

// Test getResourceName regex compilation error.
func TestGetResourceNameRegexError(t *testing.T) {
	// Save the original regex pattern
	originalResource := ResourceDefinitions["azurerm_storage_account"]

	// Modify the validation regex to be invalid for testing
	modifiedResource := originalResource
	modifiedResource.ValidationRegExp = "[" // Invalid regex pattern
	ResourceDefinitions["azurerm_storage_account"] = modifiedResource

	defer func() {
		// Restore the original after test
		ResourceDefinitions["azurerm_storage_account"] = originalResource
	}()

	_, err := getResourceName("azurerm_storage_account", "-", []string{}, "test", []string{}, "", "cafclassic", false, false, true, []string{"name"})
	if err == nil {
		t.Error("Expected regex compilation error but got none")
	}
}

// Test getNameResult with multiple resource types.
func TestGetNameResultMultipleResourceTypes(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{
		"name":           "test",
		"resource_types": []interface{}{"azurerm_storage_account", "azurerm_resource_group"},
		"random_length":  5,
	})

	err := getNameResult(rd, nil)
	if err != nil {
		t.Errorf("Unexpected error with multiple resource types: %v", err)
	}

	results := rd.Get("results").(map[string]interface{})
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got: %d", len(results))
	}

	if _, ok := results["azurerm_storage_account"]; !ok {
		t.Error("Missing result for azurerm_storage_account")
	}
	if _, ok := results["azurerm_resource_group"]; !ok {
		t.Error("Missing result for azurerm_resource_group")
	}
}

// Test getNameResult with only resource_types (no resource_type).
func TestGetNameResultOnlyResourceTypes(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{
		"name":           "test",
		"resource_type":  "",
		"resource_types": []interface{}{"azurerm_storage_account"},
		"random_length":  5,
	})

	err := getNameResult(rd, nil)
	if err != nil {
		t.Errorf("Unexpected error with only resource_types: %v", err)
	}
}

// Test getNameResult error handling.
func TestGetNameResultErrors(t *testing.T) {
	t.Run("invalid_resource_types", func(t *testing.T) {
		rd := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{
			"name":           "test",
			"resource_types": []interface{}{"invalid_resource_type"},
			"random_length":  5,
		})

		err := getNameResult(rd, nil)
		if err == nil {
			t.Error("Expected error with invalid resource types but got none")
		}
	})
}

// Test getResult with an invalid resource type.
func TestGetResultInvalidResource(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, resourceNamingConvention().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "invalid_type",
		"convention":    "random",
	})

	err := getResult(rd, nil)
	if err == nil {
		t.Error("Expected error with invalid resource type but got none")
	}
}

// Test getResult with an invalid resource mapping.
func TestGetResultInvalidResourceMapping(t *testing.T) {
	// Test with a non-existent key in ResourcesMapping
	rd := schema.TestResourceDataRaw(t, resourceNamingConvention().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "invalid_mapping",
		"convention":    "random",
	})

	err := getResult(rd, nil)
	if err == nil {
		t.Error("Expected error with invalid resource mapping but got none")
	}
}
