package azurecaf

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccErrorHandling tests error handling of the azurecaf provider
// This includes invalid resource types, invalid constraints, etc.
func TestAcc_ErrorHandling(t *testing.T) {
	// Test handling of invalid resource type with a standard error message
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccInvalidResourceTypeConfig,
				ExpectError: regexp.MustCompile(`expected resource_type to be one of`),
			},
		},
	})

	// Test handling of resource length constraints
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccExcessiveRandomLengthConfig,
				ExpectError: regexp.MustCompile(`random_length`),
			},
		},
	})

	// Test handling of negative random length
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNegativeRandomLengthConfig,
				ExpectError: regexp.MustCompile(`expected random_length to be at least \(0\)`),
			},
		},
	})

	// Test handling of invalid convention type
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccInvalidConventionConfig,
				ExpectError: regexp.MustCompile(`expected convention to be one of`),
			},
		},
	})
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
