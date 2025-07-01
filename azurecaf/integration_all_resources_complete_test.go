package azurecaf

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestAcc_AllResourceTypes tests ALL Azure resource types defined in the provider
// This test systematically goes through all 395+ resource types in batches
func TestAcc_AllResourceTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping comprehensive integration test in short mode")
	}

	provider := Provider()
	totalBatches := calculateTotalBatches()

	t.Logf("Testing ALL %d resource types across %d batches", len(ResourceDefinitions), totalBatches)

	for batchNum := 1; batchNum <= totalBatches; batchNum++ {
		t.Run(fmt.Sprintf("Batch_%d", batchNum), func(t *testing.T) {
			testResourceBatch(t, provider, batchNum)
		})
	}
}

// calculateTotalBatches returns the number of batches needed to test all resources
func calculateTotalBatches() int {
	totalResources := len(ResourceDefinitions)
	return (totalResources + maxResourcesPerBatch - 1) / maxResourcesPerBatch
}

// testResourceBatch tests a specific batch of resources
func testResourceBatch(t *testing.T, provider *schema.Provider, batchNumber int) {
	resourceBatch := generateResourceBatch(batchNumber)

	if len(resourceBatch) == 0 {
		t.Skipf("No resources in batch %d", batchNumber)
		return
	}

	t.Logf("Testing batch %d with %d resource types", batchNumber, len(resourceBatch))

	// Test azurecaf_name resource for each resource type in the batch
	t.Run("NameResource", func(t *testing.T) {
		nameResource := provider.ResourcesMap["azurecaf_name"]
		if nameResource == nil {
			t.Fatal("azurecaf_name resource not found")
		}

		for _, resourceType := range resourceBatch {
			t.Run(sanitizeResourceType(resourceType), func(t *testing.T) {
				testNameResource(t, nameResource, resourceType)
			})
		}
	})

	// Test azurecaf_name data source for each resource type in the batch
	t.Run("NameDataSource", func(t *testing.T) {
		nameDataSource := provider.DataSourcesMap["azurecaf_name"]
		if nameDataSource == nil {
			t.Fatal("azurecaf_name data source not found")
		}

		for _, resourceType := range resourceBatch {
			t.Run(sanitizeResourceType(resourceType), func(t *testing.T) {
				testNameDataSource(t, nameDataSource, resourceType)
			})
		}
	})
}

// testNameResource tests the azurecaf_name resource for a specific resource type
func testNameResource(t *testing.T, nameResource *schema.Resource, resourceType string) {
	// Test with various configurations
	testCases := []map[string]interface{}{
		// Basic configuration
		{
			"name":          "testname",
			"resource_type": resourceType,
			"prefixes":      []interface{}{"dev"},
			"suffixes":      []interface{}{"001"},
			"random_length": 5,
			"clean_input":   true,
		},
		// Configuration with separators
		{
			"name":          "testname",
			"resource_type": resourceType,
			"prefixes":      []interface{}{"prod", "web"},
			"suffixes":      []interface{}{"001", "east"},
			"separator":     "-",
			"random_length": 3,
			"clean_input":   true,
			"use_slug":      true,
		},
		// Configuration without random
		{
			"name":          "testname",
			"resource_type": resourceType,
			"prefixes":      []interface{}{"test"},
			"clean_input":   true,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Config_%d", i+1), func(t *testing.T) {
			// Create ResourceData for the azurecaf_name resource
			resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, testCase)

			// Execute create function
			err := nameResource.Create(resourceData, nil)
			if err != nil {
				t.Errorf("Failed to create name resource for %s with config %d: %v", resourceType, i+1, err)
				return
			}

			// Validate results
			result := resourceData.Get("result").(string)
			if result == "" {
				t.Errorf("Expected non-empty result for %s config %d", resourceType, i+1)
			}

			if resourceData.Id() == "" {
				t.Errorf("Expected non-empty ID for %s config %d", resourceType, i+1)
			}

			// Validate result follows naming constraints
			resourceDef, exists := ResourceDefinitions[resourceType]
			if exists {
				if len(result) < resourceDef.MinLength {
					t.Errorf("Result '%s' is shorter than min length %d for %s", result, resourceDef.MinLength, resourceType)
				}
				if len(result) > resourceDef.MaxLength {
					t.Errorf("Result '%s' is longer than max length %d for %s", result, resourceDef.MaxLength, resourceType)
				}
			}

			t.Logf("✓ Config %d for %s: %s", i+1, resourceType, result)
		})
	}
}

// testNameDataSource tests the azurecaf_name data source for a specific resource type
func testNameDataSource(t *testing.T, nameDataSource *schema.Resource, resourceType string) {
	// Create ResourceData for the azurecaf_name data source
	dataSourceData := schema.TestResourceDataRaw(t, nameDataSource.Schema, map[string]interface{}{
		"name":          "testname",
		"resource_type": resourceType,
		"prefixes":      []interface{}{"dev"},
		"suffixes":      []interface{}{"001"},
		"random_length": 5,
		"clean_input":   true,
	})

	// Execute read function
	diags := nameDataSource.ReadContext(context.Background(), dataSourceData, nil)
	if diags.HasError() {
		t.Errorf("Failed to read name data source for %s: %v", resourceType, diags)
		return
	}

	// Validate results
	result := dataSourceData.Get("result").(string)
	if result == "" {
		t.Errorf("Expected non-empty result for %s data source", resourceType)
	}

	if dataSourceData.Id() == "" {
		t.Errorf("Expected non-empty ID for %s data source", resourceType)
	}

	t.Logf("✓ Data source for %s: %s", resourceType, result)
}
