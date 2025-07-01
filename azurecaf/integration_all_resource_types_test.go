package azurecaf

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Constants to control testing batches
const (
	maxResourcesPerBatch = 20 // Maximum number of resources to test per batch
	batchToRun           = 1  // Which batch to run (1-indexed)
)

// TestAcc_ResourceTypeBatch tests a batch of Azure resource types defined in the provider
// This test uses direct provider schema testing to avoid Terraform CLI dependency
func TestAcc_ResourceTypeBatch(t *testing.T) {
	// Skip this test in short mode as it tests many resource types
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	provider := Provider()
	resourceBatch := generateResourceBatch(batchToRun)

	if len(resourceBatch) == 0 {
		t.Skipf("No resources in batch %d", batchToRun)
	}

	t.Logf("Testing batch %d with %d resource types", batchToRun, len(resourceBatch))

	// Test azurecaf_name resource for each resource type in the batch
	t.Run("NameResourceBatch", func(t *testing.T) {
		nameResource := provider.ResourcesMap["azurecaf_name"]
		if nameResource == nil {
			t.Fatal("azurecaf_name resource not found")
		}

		for _, resourceType := range resourceBatch {
			t.Run(fmt.Sprintf("Name_%s", sanitizeResourceType(resourceType)), func(t *testing.T) {
				// Create ResourceData for the azurecaf_name resource
				resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
					"name":          "testname",
					"resource_type": resourceType,
					"prefixes":      []interface{}{"dev"},
					"suffixes":      []interface{}{"001"},
					"random_length": 5,
					"clean_input":   true,
				})

				// Execute create function
				err := nameResource.Create(resourceData, nil)
				if err != nil {
					t.Errorf("Failed to create name resource for %s: %v", resourceType, err)
					return
				}

				// Check that result is set
				result := resourceData.Get("result").(string)
				if result == "" {
					t.Errorf("Expected non-empty result for %s", resourceType)
				}

				// Check that ID is set
				if resourceData.Id() == "" {
					t.Errorf("Expected non-empty ID for %s", resourceType)
				}

				t.Logf("Successfully tested %s with result: %s", resourceType, result)
			})
		}
	})

	// Test azurecaf_name data source for each resource type in the batch
	t.Run("NameDataSourceBatch", func(t *testing.T) {
		nameDataSource := provider.DataSourcesMap["azurecaf_name"]
		if nameDataSource == nil {
			t.Fatal("azurecaf_name data source not found")
		}

		for _, resourceType := range resourceBatch {
			t.Run(fmt.Sprintf("DataSource_%s", sanitizeResourceType(resourceType)), func(t *testing.T) {
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

				// Check that result is set
				result := dataSourceData.Get("result").(string)
				if result == "" {
					t.Errorf("Expected non-empty result for %s data source", resourceType)
				}

				// Check that ID is set
				if dataSourceData.Id() == "" {
					t.Errorf("Expected non-empty ID for %s data source", resourceType)
				}

				t.Logf("Successfully tested data source %s with result: %s", resourceType, result)
			})
		}
	})

	t.Logf("Batch %d testing completed successfully", batchToRun)
}

// sanitizeResourceType makes a resource type name suitable for use as a test identifier
func sanitizeResourceType(resourceType string) string {
	// Replace dots and other special characters with underscores for test names
	result := ""
	for _, char := range resourceType {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			result += string(char)
		} else {
			result += "_"
		}
	}
	return result
}

// generateResourceBatch creates a batch of resource types for testing
func generateResourceBatch(batchNumber int) []string {
	// Get all resource types
	allResourceTypes := make([]string, 0, len(ResourceDefinitions))
	for resourceType := range ResourceDefinitions {
		allResourceTypes = append(allResourceTypes, resourceType)
	}

	// Calculate start and end indices for this batch
	totalResources := len(allResourceTypes)
	startIdx := (batchNumber - 1) * maxResourcesPerBatch
	endIdx := startIdx + maxResourcesPerBatch

	// Check if this batch exists
	if startIdx >= totalResources {
		return []string{}
	}

	// Adjust endIdx if it exceeds the length
	if endIdx > totalResources {
		endIdx = totalResources
	}

	// Get the batch of resource types to test
	return allResourceTypes[startIdx:endIdx]
}
