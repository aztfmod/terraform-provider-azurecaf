# Test Templates

## main.tf

```hcl
terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = ">= 1.2.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0.0"
    }
  }
}

provider "azurecaf" {}
provider "azurerm" {
  features {}
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "azurecaf_name" "test" {
  name          = "testname"
  resource_type = "<resource_name>"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 3
  clean_input   = true
}

# Add the azurerm resource here with:
#   name = azurecaf_name.test.result
# Use ONLY minimum required attributes.
# Use HARDCODED fake Azure resource IDs for parent references.
# Example:
#   key_vault_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test/providers/Microsoft.KeyVault/vaults/kv-test"

output "result" {
  value = azurecaf_name.test.result
}

output "result_length" {
  value = length(azurecaf_name.test.result)
}
```

## terraform.rc

```hcl
provider_installation {
  dev_overrides {
    "aztfmod/azurecaf" = "<LOCAL_PLUGIN_DIR>"
  }
  direct {}
}
```

Replace `<LOCAL_PLUGIN_DIR>` with the path from the install step (e.g., `~/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/darwin_arm64`).

## tests/validate_name.tftest.hcl

```hcl
mock_provider "azurerm" {}

run "validate_generated_name" {
  command = apply

  assert {
    condition     = output.result != ""
    error_message = "Generated name output is empty"
  }

  assert {
    condition     = output.result_length >= <min_length>
    error_message = "Generated name is shorter than min_length (<min_length>)"
  }

  assert {
    condition     = output.result_length <= <max_length>
    error_message = "Generated name exceeds max_length (<max_length>)"
  }

  assert {
    condition     = can(regex("<validation_regex_pattern>", output.result))
    error_message = "Generated name does not match validation regex"
  }

  assert {
    condition     = <azurerm_resource>.test.name == output.result
    error_message = "Azure resource name does not match generated CAF name"
  }
}
```

## Critical rules

- Use `command = apply`, NOT `command = plan` (azurecaf result is computed, unknown at plan time)
- Use `output.result` / `output.result_length` in assertions (not `azurecaf_name.test.result` directly)
- Compare azurerm resource `name` against `output.result` to validate the full chain
