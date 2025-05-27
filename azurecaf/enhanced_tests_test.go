package azurecaf

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

// TestResourceNameInvalidInputs tests error handling for invalid inputs
func TestResourceNameInvalidInputs(t *testing.T) {
	testCases := []struct {
		name          string
		resourceData  map[string]interface{}
		expectedError bool
		description   string
	}{
		{
			name: "invalid_resource_type",
			resourceData: map[string]interface{}{
				"name":          "test",
				"resource_type": "nonexistent_type",
			},
			expectedError: true,
			description:   "Should fail with invalid resource type",
		},
		{
			name: "negative_random_length",
			resourceData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_storage_account",
				"random_length": -1,
			},
			expectedError: true,
			description:   "Should fail with negative random length",
		},
		{
			name: "excessive_random_length",
			resourceData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_storage_account",
				"random_length": 1000,
			},
			expectedError: true,
			description:   "Should fail when random length exceeds max length",
		},
		{
			name: "empty_prefixes_in_list",
			resourceData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_storage_account",
				"prefixes":      []interface{}{"valid", "", "another"},
			},
			expectedError: false, // Empty strings should be filtered or handled gracefully
			description:   "Should handle empty strings in prefixes list",
		},
		{
			name: "valid_input_baseline",
			resourceData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_storage_account",
			},
			expectedError: false,
			description:   "Valid input should succeed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a resource data object
			rd := schema.TestResourceDataRaw(t, resourceName().Schema, tc.resourceData)

			// Test the validation
			err := getNameResult(rd, nil)

			if tc.expectedError && err == nil {
				t.Errorf("%s: Expected error but got none", tc.description)
			}
			if !tc.expectedError && err != nil {
				t.Errorf("%s: Expected no error but got: %v", tc.description, err)
			}

			// If no error, check that result was set
			if err == nil {
				result := rd.Get("result").(string)
				if result == "" {
					t.Errorf("%s: Result should not be empty", tc.description)
				}
				t.Logf("%s: Result = %s", tc.description, result)
			}
		})
	}
}

// TestBoundaryConditions tests edge cases and boundary conditions
func TestBoundaryConditions(t *testing.T) {
	testCases := []struct {
		name         string
		resourceType string
		inputName    string
		prefixes     []interface{}
		suffixes     []interface{}
		randomLength int
		description  string
	}{
		{
			name:         "max_length_boundary",
			resourceType: "azurerm_storage_account",
			inputName:    "a", // Single character
			randomLength: 23,  // Should hit exactly max length (24)
			description:  "Test exact max length boundary",
		},
		{
			name:         "min_length_boundary",
			resourceType: "azurerm_storage_account",
			inputName:    "ab", // Two characters (min length is 3)
			description:  "Test min length boundary",
		},
		{
			name:         "unicode_characters",
			resourceType: "azurerm_resource_group",
			inputName:    "test-ñáñé-unicode",
			description:  "Test unicode character handling",
		},
		{
			name:         "multiple_separators",
			resourceType: "azurerm_resource_group",
			inputName:    "test--multiple---separators",
			description:  "Test multiple consecutive separators",
		},
		{
			name:         "long_prefixes",
			resourceType: "azurerm_storage_account",
			inputName:    "test",
			prefixes:     []interface{}{"verylongprefixthatmightcauseissues", "another"},
			description:  "Test with very long prefixes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceData := map[string]interface{}{
				"name":          tc.inputName,
				"resource_type": tc.resourceType,
				"random_length": tc.randomLength,
			}

			if tc.prefixes != nil {
				resourceData["prefixes"] = tc.prefixes
			}
			if tc.suffixes != nil {
				resourceData["suffixes"] = tc.suffixes
			}

			rd := schema.TestResourceDataRaw(t, resourceName().Schema, resourceData)
			err := getNameResult(rd, nil)

			// Log the result for debugging
			if err != nil {
				t.Logf("%s failed as expected: %v", tc.description, err)
			} else {
				result := rd.Get("result").(string)
				t.Logf("%s succeeded with result: %s (length: %d)", tc.description, result, len(result))

				// Validate result meets resource requirements
				if def, exists := ResourceDefinitions[tc.resourceType]; exists {
					if len(result) > def.MaxLength {
						t.Errorf("Result length %d exceeds max length %d", len(result), def.MaxLength)
					}
					if len(result) < def.MinLength {
						t.Errorf("Result length %d below min length %d", len(result), def.MinLength)
					}
				}
			}
		})
	}
}

// TestValidResourceTypes verifies all resource types work correctly
func TestValidResourceTypes(t *testing.T) {
	// Test a sample of resource types
	resourceTypes := []string{
		"azurerm_storage_account",
		"azurerm_resource_group",
		"azurerm_virtual_machine",
		"azurerm_key_vault",
		"azurerm_app_service",
	}

	for _, resourceType := range resourceTypes {
		t.Run(resourceType, func(t *testing.T) {
			resourceData := map[string]interface{}{
				"name":          "test",
				"resource_type": resourceType,
			}

			rd := schema.TestResourceDataRaw(t, resourceName().Schema, resourceData)
			err := getNameResult(rd, nil)

			if err != nil {
				t.Errorf("Valid resource type %s should not fail: %v", resourceType, err)
				return
			}

			result := rd.Get("result").(string)
			if result == "" {
				t.Errorf("Result should not be empty for %s", resourceType)
			}

			t.Logf("Resource type %s: %s", resourceType, result)
		})
	}
}

// TestNamingConventions tests all naming conventions
func TestNamingConventions(t *testing.T) {
	conventions := []string{"cafclassic", "cafrandom", "random", "passthrough"}

	for _, convention := range conventions {
		t.Run(convention, func(t *testing.T) {
			resourceData := map[string]interface{}{
				"name":              "test",
				"resource_type":     "azurerm_storage_account",
				"convention":        convention,
			}

			// Add random length for conventions that support it
			if convention == "cafrandom" || convention == "random" {
				resourceData["random_length"] = 5
			}

			rd := schema.TestResourceDataRaw(t, resourceName().Schema, resourceData)
			err := getNameResult(rd, nil)

			if err != nil {
				t.Errorf("Convention %s failed: %v", convention, err)
				return
			}

			result := rd.Get("result").(string)
			t.Logf("Convention %s: %s", convention, result)

			// Convention-specific assertions
			switch convention {
			case "passthrough":
				if !contains(result, "test") {
					t.Errorf("Passthrough should contain original name: %s", result)
				}
			case "cafclassic":
				// Should contain CAF prefix
				if def, exists := ResourceDefinitions["azurerm_storage_account"]; exists {
					if def.CafPrefix != "" && !contains(result, def.CafPrefix) {
						t.Errorf("CAF classic should contain prefix %s: %s", def.CafPrefix, result)
					}
				}
			case "cafrandom", "random":
				// Should contain random elements
				if len(result) <= len("test") {
					t.Errorf("Random convention should add random characters: %s", result)
				}
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			fmt.Sprintf("%s", s)[len(substr):len(s)-len(substr)] != s[len(substr):len(s)-len(substr)])))
}
