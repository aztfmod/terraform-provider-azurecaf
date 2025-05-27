# This example demonstrates the use of both the resource and data source versions of azurecaf_name

# Example 1: Using azurecaf_name as a data source
data "azurecaf_name" "rg" {
  name          = "management"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 4
  clean_input   = true
}

output "name_data_source_rg" {
  value       = data.azurecaf_name.rg.result
  description = "Resource group name generated using data source"
}

# Example 2: Using azurecaf_name as a resource
resource "azurecaf_name" "kv" {
  name          = "secrets"
  resource_type = "azurerm_key_vault"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 4
  clean_input   = true
}

output "name_resource_kv" {
  value       = azurecaf_name.kv.result
  description = "Key vault name generated using resource"
}

# Example 3: Data source for an app service 
data "azurecaf_name" "app_service" {
  name          = "multiapp"
  resource_type = "azurerm_app_service"
  prefixes      = ["test"]
  suffixes      = ["demo"]
  random_length = 5
  clean_input   = true
}

output "name_data_source_app" {
  value       = data.azurecaf_name.app_service.result
  description = "App service name generated using data source"
}

# Example 4: Resource with multiple resource types
resource "azurecaf_name" "multi_res" {
  name           = "multiapp"
  resource_type  = "azurerm_app_service"
  resource_types = ["azurerm_function_app", "azurerm_app_service_plan"]
  prefixes       = ["prod"]
  suffixes       = ["demo"]
  random_length  = 5
  clean_input    = true
}

output "name_resource_multi_primary" {
  value       = azurecaf_name.multi_res.result
  description = "Primary resource name (app service)"
}

output "name_resource_multi_all" {
  value       = azurecaf_name.multi_res.results
  description = "All resource names from multi-resource resource"
}

# Example 5: Data source with passthrough
data "azurecaf_name" "passthrough_ds" {
  name          = "existing-name-to-validate"
  resource_type = "azurerm_virtual_machine"
  passthrough   = true
  clean_input   = true
}

output "name_data_source_passthrough" {
  value       = data.azurecaf_name.passthrough_ds.result
  description = "Validated name through passthrough data source"
}

# Example 6: Resource with passthrough
resource "azurecaf_name" "passthrough_res" {
  name          = "another-name-to-validate"
  resource_type = "azurerm_virtual_machine"
  passthrough   = true
  clean_input   = true
}

output "name_resource_passthrough" {
  value       = azurecaf_name.passthrough_res.result
  description = "Validated name through passthrough resource"
}

# Example 7: Data source with use_slug = false
data "azurecaf_name" "no_slug_ds" {
  name          = "noslug"
  resource_type = "azurerm_subnet"
  prefixes      = ["net"]
  suffixes      = ["private"]
  random_length = 3
  clean_input   = true
  use_slug      = false
}

output "name_data_source_no_slug" {
  value       = data.azurecaf_name.no_slug_ds.result
  description = "Name generated using data source without slug"
}

# Example 8: Resource with use_slug = false
resource "azurecaf_name" "no_slug_res" {
  name          = "noslug"
  resource_type = "azurerm_subnet"
  prefixes      = ["net"]
  suffixes      = ["private"]
  random_length = 3
  clean_input   = true
  use_slug      = false
}

output "name_resource_no_slug" {
  value       = azurecaf_name.no_slug_res.result
  description = "Name generated using resource without slug"
}
