package azurecaf

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Test resourceAction function (data_environment_variable.go)
func TestResourceAction(t *testing.T) {
	tests := []struct {
		name         string
		set          bool
		value        string
		failsIfEmpty bool
		wantError    bool
	}{
		{name: "missing default", wantError: true},
		{name: "missing strict", failsIfEmpty: true, wantError: true},
		{name: "empty default", set: true},
		{name: "empty strict", set: true, failsIfEmpty: true, wantError: true},
		{name: "non-empty default", set: true, value: "test_value"},
		{name: "non-empty strict", set: true, value: "test_value", failsIfEmpty: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const variableName = "AZURECAF_TEST_ENV_VAR"
			if tt.set {
				t.Setenv(variableName, tt.value)
			} else if err := os.Unsetenv(variableName); err != nil {
				t.Fatalf("failed to clear test environment variable: %v", err)
			}

			rd := schema.TestResourceDataRaw(t, dataEnvironmentVariable().Schema, map[string]interface{}{
				"name":           variableName,
				"fails_if_empty": tt.failsIfEmpty,
			})
			diags := resourceAction(context.Background(), rd, nil)
			if diags.HasError() != tt.wantError {
				t.Fatalf("expected error=%t, got diagnostics: %v", tt.wantError, diags)
			}
			if tt.wantError {
				return
			}
			if rd.Id() != variableName {
				t.Fatalf("expected ID %q, got %q", variableName, rd.Id())
			}
			if got := rd.Get("value").(string); got != tt.value {
				t.Fatalf("expected value %q, got %q", tt.value, got)
			}
		})
	}
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

func TestDataNameRandomSeedPresence(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]interface{}
		wantNil  bool
		wantSeed int64
	}{
		{name: "omitted", config: map[string]interface{}{}, wantNil: true},
		{name: "explicit zero", config: map[string]interface{}{"random_seed": 0}, wantSeed: 0},
		{name: "explicit non-zero", config: map[string]interface{}{"random_seed": 42}, wantSeed: 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd := schema.TestResourceDataRaw(t, dataName().Schema, tt.config)
			seed := dataNameRandomSeed(rd)
			if tt.wantNil {
				if seed != nil {
					t.Fatalf("expected nil seed, got %d", *seed)
				}
				return
			}
			if seed == nil || *seed != tt.wantSeed {
				t.Fatalf("expected seed %d, got %v", tt.wantSeed, seed)
			}
		})
	}
}

func TestGetNameReadResultRandomSeedSemantics(t *testing.T) {
	restore := failingReader()
	defer restore()

	unseeded := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "azurerm_resource_group",
		"random_length": 4,
	})
	if err := getNameReadResult(unseeded, nil); err == nil || !strings.Contains(err.Error(), "failed to generate random suffix") {
		t.Fatalf("expected omitted seed to use crypto/rand and propagate its error, got: %v", err)
	}

	seeded := schema.TestResourceDataRaw(t, dataName().Schema, map[string]interface{}{
		"name":          "test",
		"resource_type": "azurerm_resource_group",
		"random_length": 4,
		"random_seed":   0,
	})
	if err := getNameReadResult(seeded, nil); err != nil {
		t.Fatalf("expected explicit seed 0 to use deterministic generation, got: %v", err)
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

	diags := resourceNameDelete(context.Background(), rd, nil)
	if diags.HasError() {
		t.Errorf("Expected no error, got: %v", diags)
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
