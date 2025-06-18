package azurecaf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Constants to control testing batches
const (
	maxResourcesPerBatch = 20 // Maximum number of resources to test per batch
	batchToRun           = 1  // Which batch to run (1-indexed)
)

// TestAccResourceTypeBatch tests a batch of Azure resource types defined in the provider
// using both the azurecaf_name resource and data source
func TestAcc_ResourceTypeBatch(t *testing.T) {
	// Skip this test in short mode as it tests many resource types
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	resourceBatch, batchConfig := generateResourceBatchConfig(batchToRun)

	if len(resourceBatch) == 0 {
		t.Skipf("No resources in batch %d", batchToRun)
	}

	t.Logf("Testing batch %d with %d resource types", batchToRun, len(resourceBatch))

	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: batchConfig,
				Check: resource.ComposeTestCheckFunc(
					// Check that all resources and data sources in this batch have been created and have results
					testCheckResourceTypeBatch(resourceBatch),
				),
			},
		},
	})
}

// testCheckResourceTypeBatch returns a TestCheckFunc that validates the resources and data sources in a specific batch
func testCheckResourceTypeBatch(resourceBatch []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Validate each resource type in this batch has a resource and data source result
		for _, resourceType := range resourceBatch {
			sanitizedType := sanitizeResourceType(resourceType)

			// Check resource result
			resourcePath := fmt.Sprintf("azurecaf_name.%s", sanitizedType)
			rs, ok := s.RootModule().Resources[resourcePath]
			if !ok {
				return fmt.Errorf("Not found: %s", resourcePath)
			}
			if rs.Primary.ID == "" {
				return fmt.Errorf("No ID is set for %s", resourcePath)
			}
			if rs.Primary.Attributes["result"] == "" {
				return fmt.Errorf("No result attribute for %s", resourcePath)
			}

			// Check data source result
			dataPath := fmt.Sprintf("data.azurecaf_name.%s_data", sanitizedType)
			ds, ok := s.RootModule().Resources[dataPath]
			if !ok {
				return fmt.Errorf("Not found: %s", dataPath)
			}
			if ds.Primary.ID == "" {
				return fmt.Errorf("No ID is set for %s", dataPath)
			}
			if ds.Primary.Attributes["result"] == "" {
				return fmt.Errorf("No result attribute for %s", dataPath)
			}
		}
		return nil
	}
}

// sanitizeResourceType makes a resource type name suitable for use as a Terraform identifier
func sanitizeResourceType(resourceType string) string {
	// You could replace special characters, but for now we'll use the resource type as is
	// This function can be enhanced if needed
	return resourceType
}

// Generate the Terraform configuration header
const testAccResourceTypesHeader = `
# Test Azure resource types with both azurecaf_name resource and data source
`

// generateResourceBatchConfig creates a config for a specific batch of resource types
func generateResourceBatchConfig(batchNumber int) ([]string, string) {
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
		return []string{}, ""
	}

	// Adjust endIdx if it exceeds the length
	if endIdx > totalResources {
		endIdx = totalResources
	}

	// Get the batch of resource types to test
	resourceBatch := allResourceTypes[startIdx:endIdx]

	// Generate the config for this batch
	config := testAccResourceTypesHeader
	for _, resourceType := range resourceBatch {
		sanitizedType := sanitizeResourceType(resourceType)

		// Add resource configuration
		config += fmt.Sprintf(`
# Resource for %s
resource "azurecaf_name" "%s" {
  name          = "testname"
  resource_type = "%s"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 5
  clean_input   = true
}
`, resourceType, sanitizedType, resourceType)

		// Add data source configuration
		config += fmt.Sprintf(`
# Data source for %s
data "azurecaf_name" "%s_data" {
  name          = "testname"
  resource_type = "%s"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 5
  clean_input   = true
}
`, resourceType, sanitizedType, resourceType)
	}

	return resourceBatch, config
}
