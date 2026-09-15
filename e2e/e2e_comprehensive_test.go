package e2e

import (
	"regexp"
	"strings"
	"testing"
)

const e2ETwitterProviderConfig = `
terraform {
  required_providers {
    azurecaf = {
      source = "registry.terraform.io/aztfmod/azurecaf"
    }
  }
}

provider "azurecaf" {}
`

func TestE2EDataSource(t *testing.T) {
	testDir := setupTerraformTest(t, "data_source", e2ETwitterProviderConfig+`
data "azurecaf_name" "test" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  prefixes      = ["data"]
}

output "result" {
  value = data.azurecaf_name.test.result
}
`)

	applyTerraform(t, testDir)
	assertOutputEquals(t, testDir, "result", "datastmyapp")
	assertOutputMatches(t, testDir, "result", "^[a-z0-9]{3,24}$")
	assertOutputContains(t, testDir, "result", "st")
	assertPlanNoChanges(t, testDir)
}

func TestE2ENamingConventions(t *testing.T) {
	testDir := setupTerraformTest(t, "naming_conventions", e2ETwitterProviderConfig+`
resource "azurecaf_name" "passthrough" {
  name          = "exactname"
  resource_type = "azurerm_storage_account"
  passthrough   = true
}

resource "azurecaf_name" "with_random" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  random_length = 5
  random_seed   = 12345
}

output "passthrough_result" {
  value = azurecaf_name.passthrough.result
}

output "random_result" {
  value = azurecaf_name.with_random.result
}
`)

	applyTerraform(t, testDir)
	assertOutputEquals(t, testDir, "passthrough_result", "exactname")
	assertOutputMatches(t, testDir, "random_result", "^[a-z0-9]{3,24}$")
	assertOutputContains(t, testDir, "random_result", "stmyapp")
	assertPlanNoChanges(t, testDir)
}

func TestE2EMultipleResourceTypes(t *testing.T) {
	testDir := setupTerraformTest(t, "multiple_types", e2ETwitterProviderConfig+`
resource "azurecaf_name" "storage" {
  name          = "data"
  resource_type = "azurerm_storage_account"
  prefixes      = ["st"]
}

resource "azurecaf_name" "keyvault" {
  name          = "secrets"
  resource_type = "azurerm_key_vault"
  prefixes      = ["kv"]
  suffixes      = ["prod"]
}

resource "azurecaf_name" "vm" {
  name          = "webserver"
  resource_type = "azurerm_linux_virtual_machine"
  prefixes      = ["vm"]
  random_length = 3
  random_seed   = 12345
}

output "storage_result" {
  value = azurecaf_name.storage.result
}

output "keyvault_result" {
  value = azurecaf_name.keyvault.result
}

output "vm_result" {
  value = azurecaf_name.vm.result
}
`)

	applyTerraform(t, testDir)
	assertOutputMatches(t, testDir, "storage_result", "^[a-z0-9]{3,24}$")
	assertOutputContains(t, testDir, "storage_result", "st")
	assertOutputContains(t, testDir, "storage_result", "data")
	assertOutputMatches(t, testDir, "keyvault_result", "^[a-z][0-9a-z-]{1,23}$")
	assertOutputContains(t, testDir, "keyvault_result", "kv")
	assertOutputContains(t, testDir, "keyvault_result", "prod")
	assertOutputMatches(t, testDir, "vm_result", `^[^/"\[\]:|<>+=;,?*@&_][^/"\[\]:|<>+=;,?*@&]{0,62}[^/"\[\]:|<>+=;,?*@&.-]$`)
	assertOutputContains(t, testDir, "vm_result", "vm")
	assertOutputContains(t, testDir, "vm_result", "webserver")
	assertPlanNoChanges(t, testDir)
}

func TestE2EImportFunctionality(t *testing.T) {
	testDir := setupTerraformTest(t, "import_functionality", e2ETwitterProviderConfig+`
resource "azurecaf_name" "imported_storage" {
  name          = "stmyexistingapp"
  resource_type = "azurerm_storage_account"
  clean_input   = false
  passthrough   = true
  random_length = 0
  separator     = ""
  use_slug      = false
}

output "imported_result" {
  value = azurecaf_name.imported_storage.result
}
`)

	importOutput := runTerraformExpectSuccess(t, testDir, "import", "azurecaf_name.imported_storage", "azurerm_storage_account:stmyexistingapp")
	if !strings.Contains(importOutput, "Import successful") && !strings.Contains(importOutput, "imported successfully") {
		t.Fatalf("import output does not indicate success\nOutput:\n%s", importOutput)
	}

	applyTerraform(t, testDir)
	assertOutputEquals(t, testDir, "imported_result", "stmyexistingapp")
	assertPlanNoChanges(t, testDir)
}

func TestE2E_OutputValueAssertions(t *testing.T) {
	testDir := setupTerraformTest(t, "output_value_assertions", e2ETwitterProviderConfig+`
resource "azurecaf_name" "storage" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  prefixes      = ["dev"]
}

resource "azurecaf_name" "rg" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
}

output "storage_result" {
  value = azurecaf_name.storage.result
}

output "resource_group_result" {
  value = azurecaf_name.rg.result
}
`)

	applyTerraform(t, testDir)
	assertOutputMatches(t, testDir, "storage_result", "^[a-z0-9]{3,24}$")
	assertOutputContains(t, testDir, "storage_result", "dev")
	assertOutputContains(t, testDir, "storage_result", "st")
	assertOutputContains(t, testDir, "storage_result", "myapp")
	assertOutputEquals(t, testDir, "resource_group_result", "dev-rg-myapp-001")
	assertOutputMatches(t, testDir, "resource_group_result", `^[-\w\._\(\)]{1,80}$`)
}

func TestE2E_AllNamingConventions(t *testing.T) {
	testDir := setupTerraformTest(t, "all_naming_conventions", e2ETwitterProviderConfig+`
resource "azurecaf_naming_convention" "cafclassic" {
  convention    = "cafclassic"
  name          = "myapp"
  resource_type = "st"
  prefix        = "dev"
}

resource "azurecaf_naming_convention" "cafrandom" {
  convention    = "cafrandom"
  name          = "myapp"
  resource_type = "st"
  prefix        = "dev"
}

resource "azurecaf_naming_convention" "random" {
  convention    = "random"
  resource_type = "st"
  prefix        = "dev"
}

resource "azurecaf_naming_convention" "passthrough" {
  convention    = "passthrough"
  name          = "custom-rg-name"
  resource_type = "rg"
}

output "cafclassic_result" {
  value = azurecaf_naming_convention.cafclassic.result
}

output "cafrandom_result" {
  value = azurecaf_naming_convention.cafrandom.result
}

output "random_result" {
  value = azurecaf_naming_convention.random.result
}

output "passthrough_result" {
  value = azurecaf_naming_convention.passthrough.result
}
`)

	applyTerraform(t, testDir)
	assertOutputEquals(t, testDir, "cafclassic_result", "devstmyapp")
	assertOutputMatches(t, testDir, "cafclassic_result", "^[a-z0-9]{3,24}$")
	assertOutputContains(t, testDir, "cafclassic_result", "st")
	assertOutputMatches(t, testDir, "cafrandom_result", "^[a-z0-9]{24}$")
	assertOutputContains(t, testDir, "cafrandom_result", "devstmyapp")
	assertOutputMatches(t, testDir, "random_result", "^[a-z0-9]{24}$")
	assertOutputContains(t, testDir, "random_result", "dev")
	assertOutputEquals(t, testDir, "passthrough_result", "custom-rg-name")
	assertOutputMatches(t, testDir, "passthrough_result", `^[-\w\._\(\)]{1,80}$`)
}

func TestE2E_LengthConstraints(t *testing.T) {
	testDir := setupTerraformTest(t, "length_constraints", e2ETwitterProviderConfig+`
resource "azurecaf_name" "storage" {
  name          = "superlongapplicationname"
  resource_type = "azurerm_storage_account"
  prefixes      = ["dev", "platform"]
  suffixes      = ["prod", "001"]
}

output "storage_result" {
  value = azurecaf_name.storage.result
}
`)

	applyTerraform(t, testDir)
	storageResult := getOutputString(t, testDir, "storage_result")
	if len(storageResult) != 24 {
		t.Fatalf("expected storage account name to be truncated to 24 characters, got %d (%q)", len(storageResult), storageResult)
	}
	assertOutputEquals(t, testDir, "storage_result", "superlongapplicationname")
	assertOutputMatches(t, testDir, "storage_result", "^[a-z0-9]{24}$")
}

func TestE2E_MultipleResults(t *testing.T) {
	testDir := setupTerraformTest(t, "multiple_results", e2ETwitterProviderConfig+`
resource "azurecaf_name" "multi" {
  name           = "sharedapp"
  prefixes       = ["dev"]
  resource_types = ["azurerm_storage_account", "azurerm_resource_group", "azurerm_key_vault"]
}

output "results" {
  value = azurecaf_name.multi.results
}
`)

	applyTerraform(t, testDir)
	results := getOutputStringMap(t, testDir, "results")

	if len(results) != 3 {
		t.Fatalf("expected 3 generated results, got %d: %#v", len(results), results)
	}

	if storage := results["azurerm_storage_account"]; storage == "" {
		t.Fatalf("expected storage account entry in results map: %#v", results)
	} else {
		if !strings.Contains(storage, "st") {
			t.Fatalf("expected storage result %q to contain storage slug", storage)
		}
		if matched := regexpMustMatch(t, "^[a-z0-9]{3,24}$", storage); !matched {
			t.Fatalf("storage result %q does not match storage account regex", storage)
		}
	}

	if rg := results["azurerm_resource_group"]; rg == "" {
		t.Fatalf("expected resource group entry in results map: %#v", results)
	} else {
		if !strings.Contains(rg, "rg") {
			t.Fatalf("expected resource group result %q to contain rg slug", rg)
		}
		if matched := regexpMustMatch(t, `^[-\w\._\(\)]{1,80}$`, rg); !matched {
			t.Fatalf("resource group result %q does not match resource group regex", rg)
		}
	}

	if keyVault := results["azurerm_key_vault"]; keyVault == "" {
		t.Fatalf("expected key vault entry in results map: %#v", results)
	} else {
		if !strings.Contains(keyVault, "kv") {
			t.Fatalf("expected key vault result %q to contain kv slug", keyVault)
		}
		if matched := regexpMustMatch(t, "^[a-z][0-9a-z-]{1,23}$", keyVault); !matched {
			t.Fatalf("key vault result %q does not match key vault regex", keyVault)
		}
	}
}

func TestE2E_ErrorMessages(t *testing.T) {
	t.Run("missing_resource_type", func(t *testing.T) {
		testDir := setupTerraformTest(t, "error_missing_resource_type", e2ETwitterProviderConfig+`
resource "azurecaf_name" "invalid" {
  name = "missingtype"
}
`)

		applyOutput := runTerraformApplyExpectError(t, testDir)
		for _, expected := range []string{"resource_type and resource_types parameters are empty", "you must specify at least one resource type"} {
			if !strings.Contains(applyOutput, expected) {
				t.Fatalf("expected apply error to contain %q\nOutput:\n%s", expected, applyOutput)
			}
		}
	})

	t.Run("max_length_exceeded", func(t *testing.T) {
		testDir := setupTerraformTest(t, "error_max_length_exceeded", e2ETwitterProviderConfig+`
resource "azurecaf_name" "invalid" {
  name                           = "superlongapplicationname"
  prefixes                       = ["dev", "platform", "shared"]
  suffixes                       = ["prod", "001"]
  resource_type                  = "azurerm_storage_account"
  error_when_exceeding_max_length = true
}
`)

		applyOutput := runTerraformApplyExpectError(t, testDir)
		for _, expected := range []string{"exceeds maximum length of 24", "composed name"} {
			if !strings.Contains(applyOutput, expected) {
				t.Fatalf("expected apply error to contain %q\nOutput:\n%s", expected, applyOutput)
			}
		}
	})
}

func TestE2E_StateConsistency(t *testing.T) {
	testDir := setupTerraformTest(t, "state_consistency", e2ETwitterProviderConfig+`
resource "azurecaf_name" "rg" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
}

output "result" {
  value = azurecaf_name.rg.result
}
`)

	applyTerraform(t, testDir)
	assertOutputEquals(t, testDir, "result", "dev-rg-myapp-001")
	assertPlanNoChanges(t, testDir)
}

func regexpMustMatch(t *testing.T, pattern, value string) bool {
	t.Helper()

	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		t.Fatalf("invalid regex %q: %v", pattern, err)
	}
	return matched
}
