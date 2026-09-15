package e2e

import "testing"

func TestE2EBasic(t *testing.T) {
	testDir := setupTerraformTest(t, "basic", `
terraform {
  required_providers {
    azurecaf = {
      source = "registry.terraform.io/aztfmod/azurecaf"
    }
  }
}

provider "azurecaf" {}

resource "azurecaf_name" "test" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  prefixes      = ["test"]
}

output "result" {
  value = azurecaf_name.test.result
}
`)

	applyTerraform(t, testDir)
	assertOutputEquals(t, testDir, "result", "teststmyapp")
	assertOutputMatches(t, testDir, "result", "^[a-z0-9]{3,24}$")
	assertOutputContains(t, testDir, "result", "st")
	assertPlanNoChanges(t, testDir)
}
