package azurecaf

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourcesIntegration tests integration between data sources and resources
func TestAcc_DataSourcesIntegration(t *testing.T) {
	// Skip this test if we can't access external network resources
	// This test requires Terraform CLI which needs to connect to checkpoint-api.hashicorp.com
	t.Skip("Skipping acceptance test - requires network access to Terraform CLI")
	
	// Set environment variable for testing
	os.Setenv("TEST_ENV_VAR", "test-env-value")
	defer os.Unsetenv("TEST_ENV_VAR")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcesIntegrationConfig,
				Check: resource.ComposeTestCheckFunc(
					// Check environment variable data source
					resource.TestCheckResourceAttr(
						"data.azurecaf_environment_variable.test_env", "value", "test-env-value"),

					// Check name data source basic
					resource.TestCheckResourceAttrSet(
						"data.azurecaf_name.simple", "result"),

					// Check name data source with prefixes
					resource.TestCheckResourceAttr(
						"data.azurecaf_name.with_prefixes", "result", "devstoragedata"),

					// Check name data source with environment variable
					resource.TestCheckResourceAttrSet(
						"data.azurecaf_name.with_env_var", "result"),

					// Check resource that uses output from data sources
					resource.TestCheckResourceAttrSet(
						"azurecaf_naming_convention.combined", "result"),
				),
			},
		},
	})
}

// Configuration for data sources integration test
const testAccDataSourcesIntegrationConfig = `
# Test environment variable
data "azurecaf_environment_variable" "test_env" {
  name = "TEST_ENV_VAR"
}

# Basic name data source
data "azurecaf_name" "simple" {
  name          = "myapp"
  resource_type = "azurerm_app_service"
  random_length = 5
}

# Name data source with prefixes and suffixes
data "azurecaf_name" "with_prefixes" {
  name          = "storage-data"
  prefixes      = ["dev"]
  resource_type = "azurerm_storage_account"
  use_slug      = false
  clean_input   = true
  separator     = "-"
}

# Name data source using environment variable
data "azurecaf_name" "with_env_var" {
  name          = "${data.azurecaf_environment_variable.test_env.value}-resource"
  resource_type = "azurerm_resource_group"
}

# Naming convention resource that uses data source outputs
resource "azurecaf_naming_convention" "combined" {
  name          = "${data.azurecaf_name.with_prefixes.result}-combined"
  resource_type = "rg"
  convention    = "random"
}
`
