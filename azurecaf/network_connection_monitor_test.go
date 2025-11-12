package azurecaf

import (
	"regexp"
	"testing"
)

// TestNetworkConnectionMonitorDefinition tests that the network connection monitor resource is properly defined
func TestNetworkConnectionMonitorDefinition(t *testing.T) {
	// Get the resource definition
	resourceDef, exists := ResourceDefinitions["azurerm_network_connection_monitor"]
	if !exists {
		t.Fatal("azurerm_network_connection_monitor not found in ResourceDefinitions")
	}

	// Verify basic properties
	if resourceDef.CafPrefix != "cm" {
		t.Errorf("Expected slug 'cm', got '%s'", resourceDef.CafPrefix)
	}

	if resourceDef.MinLength != 1 {
		t.Errorf("Expected min_length 1, got %d", resourceDef.MinLength)
	}

	if resourceDef.MaxLength != 80 {
		t.Errorf("Expected max_length 80, got %d", resourceDef.MaxLength)
	}

	if resourceDef.Scope != "parent" {
		t.Errorf("Expected scope 'parent', got '%s'", resourceDef.Scope)
	}

	if !resourceDef.Dashes {
		t.Error("Expected dashes to be allowed")
	}

	if resourceDef.LowerCase {
		t.Error("Expected lowercase not to be required")
	}

	// Verify validation regex
	expectedRegex := "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"
	if resourceDef.ValidationRegExp != expectedRegex {
		t.Errorf("Expected validation regex '%s', got '%s'", expectedRegex, resourceDef.ValidationRegExp)
	}

	// Verify cleanup regex
	expectedCleanupRegex := "[^0-9A-Za-z_.-]"
	if resourceDef.RegEx != expectedCleanupRegex {
		t.Errorf("Expected cleanup regex '%s', got '%s'", expectedCleanupRegex, resourceDef.RegEx)
	}

	t.Log("Network connection monitor resource definition validated successfully")
}

// TestNetworkConnectionMonitorValidation tests the validation regex for network connection monitor
func TestNetworkConnectionMonitorValidation(t *testing.T) {
	resourceDef, exists := ResourceDefinitions["azurerm_network_connection_monitor"]
	if !exists {
		t.Fatal("azurerm_network_connection_monitor not found in ResourceDefinitions")
	}

	// Compile the validation regex
	validationRegex, err := regexp.Compile(resourceDef.ValidationRegExp)
	if err != nil {
		t.Fatalf("Failed to compile validation regex: %v", err)
	}

	validNames := []string{
		"ab",                  // min practical length (2 chars due to regex requiring start and end)
		"monitor1",            // simple name
		"test-monitor",        // with hyphen
		"test.monitor",        // with period
		"test_monitor",        // with underscore
		"Test-Monitor-123",    // mixed case with numbers
		"cm-prod-monitor-001", // full name with segments
		"a1234567890123456789012345678901234567890123456789012345678901234567890123456789", // max length (80 chars)
	}

	invalidNames := []string{
		"",         // empty
		"-monitor", // starts with hyphen
		"monitor-", // ends with hyphen
		".monitor", // starts with period
		"monitor.", // ends with period
		"mon@itor", // invalid character
		"mon itor", // space
		"a12345678901234567890123456789012345678901234567890123456789012345678901234567890", // too long (81 chars)
	}

	for _, name := range validNames {
		t.Run("valid_"+name, func(t *testing.T) {
			if !validationRegex.MatchString(name) {
				t.Errorf("Expected '%s' to be valid", name)
			}
		})
	}

	for _, name := range invalidNames {
		t.Run("invalid_"+name, func(t *testing.T) {
			if validationRegex.MatchString(name) {
				t.Errorf("Expected '%s' to be invalid", name)
			}
		})
	}
}
