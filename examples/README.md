# Azure CAF Terraform Provider Examples

This directory contains example configurations for the Azure Cloud Adoption Framework (CAF) Terraform Provider. These examples demonstrate the different ways to use the provider to generate Azure resource names that comply with the naming conventions specified in the Cloud Adoption Framework.

## Example Files

- **resource_name.tf**: Demonstrates the use of the `azurecaf_name` resource
- **name_resource_and_datasource.tf**: Compares the use of both the resource and data source versions of `azurecaf_name`

## Key Examples

### Resource vs Data Source

Both resource and data source versions of `azurecaf_name` are available. The data source is evaluated before resources are created and can be viewed at plan time. This is the recommended approach for most scenarios.

Example using data source:
```hcl
data "azurecaf_name" "rg_example" {
  name          = "demogroup"
  resource_type = "azurerm_resource_group"
  random_length = 5
  clean_input   = true
}

resource "azurerm_resource_group" "rg" {
  name     = data.azurecaf_name.rg_example.result
  location = "southeastasia"
}
```

Example using resource:
```hcl
resource "azurecaf_name" "kv" {
  name          = "secrets"
  resource_type = "azurerm_key_vault"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 4
  clean_input   = true
}

resource "azurerm_key_vault" "example" {
  name                = azurecaf_name.kv.result
  location            = "eastus"
  resource_group_name = azurerm_resource_group.example.name
}
```

### Multiple Resource Types

The resource version (but not data source) supports generating names for multiple resource types at once:

```hcl
resource "azurecaf_name" "multi_res" {
  name           = "multiapp"
  resource_type  = "azurerm_app_service"
  resource_types = ["azurerm_function_app", "azurerm_app_service_plan"]
  prefixes       = ["prod"]
  suffixes       = ["demo"]
  random_length  = 5
}

# Access the primary resource name
primary_name = azurecaf_name.multi_res.result

# Access all resource names
all_names = azurecaf_name.multi_res.results
```

### Other Features

- **clean_input**: Controls whether the provider sanitizes inputs to comply with Azure naming restrictions
- **passthrough**: Validates a name without modifying it (useful for checking existing names)
- **use_slug**: Controls whether to add the resource type abbreviation (slug) to the name
- **separator**: Configures the character used between name components (default: "-")

## Running the Examples

To run these examples:

1. Build the provider locally:
   ```
   make build
   ```

2. Run the examples:
   ```
   make test
   ```

The `make test` command will:
- Build the provider locally
- Set up the correct development overrides for Terraform
- Initialize and apply all the examples in this directory
