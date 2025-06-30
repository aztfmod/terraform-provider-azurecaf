package azurecaf

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestAcc_ErrorHandling tests error handling of the azurecaf provider
// This test uses direct provider schema testing to avoid Terraform CLI dependency
func TestAcc_ErrorHandling(t *testing.T) {
	provider := Provider()
	
	// Test handling of invalid resource type
	t.Run("InvalidResourceType", func(t *testing.T) {
		namingConventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
		if namingConventionResource == nil {
			t.Fatal("azurecaf_naming_convention resource not found")
		}

		// Create ResourceData with invalid resource type
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"name":          "test",
			"prefix":        "dev",
			"resource_type": "not_a_valid_type",
			"convention":    "cafclassic",
		})

		// Try to create the resource - should fail validation
		err := namingConventionResource.Create(resourceData, nil)
		if err == nil {
			t.Error("Expected error for invalid resource type, but got none")
		}
		if err != nil && !strings.Contains(err.Error(), "Invalid resource type") {
			t.Errorf("Expected error about resource type validation, got: %v", err)
		}
	})

	// Test handling of excessive random length for azurecaf_name resource
	t.Run("ExcessiveRandomLength", func(t *testing.T) {
		nameResource := provider.ResourcesMap["azurecaf_name"]
		if nameResource == nil {
			t.Fatal("azurecaf_name resource not found")
		}

		// Create ResourceData with excessive random length
		resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
			"name":          "test",
			"prefixes":      []interface{}{"dev"},
			"resource_type": "azurerm_storage_account",
			"random_length": 30, // Too long for storage account
		})

		// Try to create the resource - should fail validation
		err := nameResource.Create(resourceData, nil)
		if err == nil {
			t.Error("Expected error for excessive random length, but got none")
		}
		if err != nil && !strings.Contains(err.Error(), "random_length") {
			t.Errorf("Expected error about random_length, got: %v", err)
		}
	})

	// Test handling of negative random length
	t.Run("NegativeRandomLength", func(t *testing.T) {
		nameResource := provider.ResourcesMap["azurecaf_name"]
		
		// Test at schema validation level
		schema := nameResource.Schema["random_length"]
		if schema == nil {
			t.Fatal("random_length schema not found")
		}

		// Validate that negative values are rejected by schema validation
		_, errors := schema.ValidateFunc(-5, "random_length")
		if len(errors) == 0 {
			t.Error("Expected schema validation error for negative random_length")
		}
		
		found := false
		for _, err := range errors {
			if strings.Contains(err.Error(), "expected random_length to be at least") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected validation error about minimum value, got: %v", errors)
		}
	})

	// Test handling of invalid convention type
	t.Run("InvalidConvention", func(t *testing.T) {
		namingConventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
		
		// Test at schema validation level first
		conventionSchema := namingConventionResource.Schema["convention"]
		if conventionSchema == nil {
			t.Fatal("convention schema not found")
		}

		// Test that the schema validation rejects invalid values
		_, errors := conventionSchema.ValidateFunc("invalid_convention", "convention")
		if len(errors) == 0 {
			t.Error("Expected schema validation error for invalid convention")
		} else {
			found := false
			for _, err := range errors {
				if strings.Contains(err.Error(), "expected convention to be one of") {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected validation error about convention values, got: %v", errors)
			}
		}
	})

	t.Log("Error handling tests completed successfully")
}

// Invalid resource type configuration
const testAccInvalidResourceTypeConfig = `
resource "azurecaf_naming_convention" "invalid_type" {
  name           = "test"
  prefix         = "dev"
  resource_type  = "not_a_valid_type"
  convention     = "cafclassic"
}
`

// Configuration with excessive random length for a resource type
const testAccExcessiveRandomLengthConfig = `
resource "azurecaf_name" "excessive_length" {
  name           = "test"
  prefixes       = ["dev"]
  resource_type  = "azurerm_storage_account"
  random_length  = 30
}
`

// Configuration with negative random length
const testAccNegativeRandomLengthConfig = `
resource "azurecaf_name" "negative_length" {
  name           = "test"
  prefixes       = ["dev"]
  resource_type  = "azurerm_resource_group"
  random_length  = -5
}
`

// Configuration with invalid convention type
const testAccInvalidConventionConfig = `
resource "azurecaf_naming_convention" "invalid_convention" {
  name           = "test"
  prefix         = "dev"
  resource_type  = "rg"
  convention     = "invalid_convention"
}
`
