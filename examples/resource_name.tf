
#Storage account test
resource "azurecaf_name" "classic_st" {
  name          = "log2"
  resource_type = "azurerm_storage_account"
}

output "caf_name_classic_st" {
  value       = azurecaf_name.classic_st.result
  description = "Random result based on the resource type"
}

resource "azurecaf_name" "azurerm_cognitive_account" {
  name          = "cogsdemo"
  resource_type = "azurerm_cognitive_account"
  prefixes      = ["a", "z"]
  suffixes      = ["prod"]
  random_length = 5
  random_seed   = 12343
  clean_input   = true
  separator     = "-"
}

output "azurerm_cognitive_account" {
  value       = azurecaf_name.azurerm_cognitive_account.result
  description = "Random result based on the resource type"
}

#Azure Open AI Deployment test
resource "azurecaf_name" "azurerm_synapse_workspace" {
  name          = "openai-deployment"
  resource_type = "azurerm_synapse_workspace"
  prefixes      = ["a", "b"]
  suffixes      = ["y", "z"]
  random_length = 5
  clean_input   = true
}

output "azurerm_synapse_workspace" {
  value = azurecaf_name.azurerm_synapse_workspace.result
  description = "Random result based on the resource type"
}

resource "azurecaf_name" "multiple_resources" {
  name           = "cogsdemo2"
  resource_type  = "azurerm_cognitive_account"
  resource_types = ["azurerm_storage_account"]
  prefixes       = ["a", "b"]
  suffixes       = ["prod"]
  random_length  = 4
  random_seed    = 12343
  clean_input    = true
  separator      = "-"
}

output "multiple_resources" {
  value = azurecaf_name.multiple_resources.results
}

output "multiple_resources_main" {
  value = azurecaf_name.multiple_resources.result
}
