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

// Test regex compilation error in getResult by handling panic (since there's a nil pointer issue)
func TestGetResultRegexError(t *testing.T) {
	// Some code paths in getResult don't properly handle invalid regex and cause panics
	// We can test this by using defer/recover

	// Save the original resource
	originalResource := Resources["st"]

	// Create a modified version with invalid regex pattern - this should error on compile
	modifiedResource := originalResource
	modifiedResource.RegEx = "[" // Invalid regex pattern that will cause compile error
	Resources["st"] = modifiedResource

	defer func() {
		// Restore original after test
		Resources["st"] = originalResource

		// Recover from panic
		if r := recover(); r == nil {
			t.Error("Expected panic but none occurred")
		}
	}()

	rd := schema.TestResourceDataRaw(t, resourceNamingConvention().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "st",
		"convention":    "random",
	})

	// This will cause a panic when it hits the invalid regex
	_ = getResult(rd, nil)
}

// Test getResult with validation regex error
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

		// Recover from panic
		if r := recover(); r == nil {
			t.Error("Expected panic but none occurred")
		}
	}()

	rd := schema.TestResourceDataRaw(t, resourceNamingConvention().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "st",
		"convention":    "random",
	})

	// This will cause a panic when it hits the invalid validation regex
	_ = getResult(rd, nil)
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
