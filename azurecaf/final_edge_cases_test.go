package azurecaf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Test getResourceName validation error path
func TestGetResourceNameValidationError(t *testing.T) {
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

// Test regex compilation error in getResult - now returns an error instead of panicking
func TestGetResultRegexError(t *testing.T) {
	// Save the original resource
	originalResource := Resources["st"]

	// Create a modified version with invalid regex pattern - this should now return an error
	modifiedResource := originalResource
	modifiedResource.RegEx = "[" // Invalid regex pattern that will cause compile error
	Resources["st"] = modifiedResource

	defer func() {
		// Restore original after test
		Resources["st"] = originalResource
	}()

	rd := schema.TestResourceDataRaw(t, resourceNamingConvention().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "st",
		"convention":    "random",
	})

	err := getResult(rd, nil)
	if err == nil {
		t.Error("Expected error for invalid regex pattern but got none")
	}
}

// Test getResult with validation regex error - now returns an error instead of panicking
func TestGetResultValidationRegexError(t *testing.T) {
	// Save the original resource
	originalResource := Resources["st"]

	// Create a modified version with valid regex but validation regex that's invalid
	modifiedResource := originalResource
	modifiedResource.ValidationRegExp = "[" // Invalid regex pattern
	Resources["st"] = modifiedResource

	defer func() {
		// Restore original after test
		Resources["st"] = originalResource
	}()

	rd := schema.TestResourceDataRaw(t, resourceNamingConvention().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "st",
		"convention":    "random",
	})

	err := getResult(rd, nil)
	if err == nil {
		t.Error("Expected error for invalid validation regex pattern but got none")
	}
}

// Test getResult error handling with validation match failure
func TestGetResultValidationMatchError(t *testing.T) {
	// Save the original resources
	original := Resources["st"]

	// Create a modified version with regex that won't match any input
	modified := original
	// Keep valid compilation but create a pattern that won't match our input
	modified.ValidationRegExp = "^$" // This will only match empty string

	Resources["st"] = modified

	defer func() {
		// Restore original after test
		Resources["st"] = original
	}()

	rd := schema.TestResourceDataRaw(t, resourceNamingConvention().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "st",
		"convention":    "random",
	})

	// This should fail validation
	err := getResult(rd, nil)

	if err == nil {
		t.Error("Expected validation match error but got none")
	}
}
