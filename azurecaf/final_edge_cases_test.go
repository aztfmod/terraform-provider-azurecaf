package azurecaf

import (
	"testing"
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
	_, err := getResourceName("azurerm_storage_account", "-", []string{}, "test", []string{}, "", "cafclassic", false, false, true, []string{"name"}, false)

	if err == nil {
		t.Error("Expected validation error but got none")
	}
}
