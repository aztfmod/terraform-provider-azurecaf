package azurecaf

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Test resourceAction function (data_environment_variable.go)
func TestResourceAction(t *testing.T) {
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

		// Debug: print actual diagnostic details
		if len(diags) > 0 {
			t.Logf("Diagnostic: Severity=%d, Summary=%s", diags[0].Severity, diags[0].Summary)
		}

		if diags[0].Severity != 0 { // Error severity should be 0 based on actual behavior
			t.Errorf("Expected Error diagnostic with severity 0, got severity: %d", diags[0].Severity)
		}
	})
}

// Test dataNameRead function (data_name.go)
func TestDataNameRead(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "azurerm_storage_account",
	})

	diags := dataNameRead(context.Background(), rd, nil)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got: %v", diags)
	}

	result := rd.Get("result").(string)
	if result == "" {
		t.Error("Expected result to be set")
	}
}

// Test getNameReadResult function (data_name.go)
func TestGetNameReadResult(t *testing.T) {
	testCases := []struct {
		name         string
		resourceData map[string]interface{}
		expectedErr  bool
	}{
		{
			name: "valid_basic_case",
			resourceData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_storage_account",
			},
			expectedErr: false,
		},
		{
			name: "with_prefixes_and_suffixes",
			resourceData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_storage_account",
				"prefixes":      []interface{}{"prefix1", "prefix2"},
				"suffixes":      []interface{}{"suffix1"},
				"separator":     "-",
			},
			expectedErr: false,
		},
		{
			name: "with_random_settings",
			resourceData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_storage_account",
				"random_length": 5,
				"random_seed":   12345,
			},
			expectedErr: false,
		},
		{
			name: "with_clean_input",
			resourceData: map[string]interface{}{
				"name":          "test-with-special-chars!@#",
				"resource_type": "azurerm_storage_account",
				"clean_input":   true,
			},
			expectedErr: false,
		},
		{
			name: "with_passthrough",
			resourceData: map[string]interface{}{
				"name":        "test",
				"passthrough": true,
			},
			expectedErr: false,
		},
		{
			name: "with_use_slug",
			resourceData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_storage_account",
				"use_slug":      true,
			},
			expectedErr: false,
		},
		{
			name: "invalid_negative_random_length",
			resourceData: map[string]interface{}{
				"name":          "test",
				"resource_type": "azurerm_storage_account",
				"random_length": -5,
			},
			expectedErr: false, // Data source doesn't validate negative length currently
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rd := schema.TestResourceDataRaw(t, dataName().Schema, tc.resourceData)
			err := getNameReadResult(rd, nil)

			if tc.expectedErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectedErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tc.expectedErr {
				result := rd.Get("result").(string)
				if result == "" {
					t.Error("Expected result to be set")
				}
			}
		})
	}
}

// Test resourceNameDelete function (resource_name.go)
func TestResourceNameDelete(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, resourceName().Schema, map[string]interface{}{
		"name": "test",
	})

	err := resourceNameDelete(rd, nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Test resourceNamingConventionDelete function (resource_naming_convention.go)
func TestResourceNamingConventionDelete(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, resourceNamingConvention().Schema, map[string]interface{}{
		"name": "test",
	})

	err := resourceNamingConventionDelete(rd, nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Test getResource function edge cases
func TestGetResourceEdgeCases(t *testing.T) {
	// Test with ResourceMaps lookup (like "st" -> "azurerm_storage_account")
	t.Run("resource_maps_lookup", func(t *testing.T) {
		resource, err := getResource("st")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if resource == nil {
			t.Error("Expected resource to be found")
		}
		if resource.ResourceTypeName != "azurerm_storage_account" {
			t.Errorf("Expected name 'azurerm_storage_account', got: %s", resource.ResourceTypeName)
		}
	})

	// Test with direct ResourceDefinitions lookup
	t.Run("direct_resource_lookup", func(t *testing.T) {
		resource, err := getResource("azurerm_storage_account")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if resource == nil {
			t.Error("Expected resource to be found")
		}
	})

	// Test with invalid resource type
	t.Run("invalid_resource_type", func(t *testing.T) {
		resource, err := getResource("invalid_resource")
		if err == nil {
			t.Error("Expected error for invalid resource type")
		}
		if resource != nil {
			t.Error("Expected nil resource for invalid type")
		}
	})
}

// Test trimResourceName function edge cases
func TestTrimResourceNameEdgeCases(t *testing.T) {
	testCases := []struct {
		name         string
		resourceName string
		maxLength    int
		expected     string
	}{
		{
			name:         "name_shorter_than_max",
			resourceName: "short",
			maxLength:    10,
			expected:     "short",
		},
		{
			name:         "name_equal_to_max",
			resourceName: "exactten12",
			maxLength:    10,
			expected:     "exactten12",
		},
		{
			name:         "name_longer_than_max",
			resourceName: "verylongresourcename",
			maxLength:    10,
			expected:     "verylongre",
		},
		{
			name:         "zero_max_length",
			resourceName: "test",
			maxLength:    0,
			expected:     "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := trimResourceName(tc.resourceName, tc.maxLength)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// Test getSlug function with different conventions
func TestGetSlugExtended(t *testing.T) {
	testCases := []struct {
		name         string
		resourceType string
		convention   string
		expected     string
	}{
		{
			name:         "cafclassic_convention",
			resourceType: "azurerm_storage_account",
			convention:   ConventionCafClassic,
			expected:     "st",
		},
		{
			name:         "cafrandom_convention",
			resourceType: "azurerm_storage_account",
			convention:   ConventionCafRandom,
			expected:     "st",
		},
		{
			name:         "random_convention",
			resourceType: "azurerm_storage_account",
			convention:   ConventionRandom,
			expected:     "",
		},
		{
			name:         "passthrough_convention",
			resourceType: "azurerm_storage_account",
			convention:   ConventionPassThrough,
			expected:     "",
		},
		{
			name:         "unknown_resource_type",
			resourceType: "unknown_resource",
			convention:   ConventionCafClassic,
			expected:     "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getSlug(tc.resourceType, tc.convention)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}
