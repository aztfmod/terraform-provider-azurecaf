
#Storage account test
resource "azurecaf_name" "classic_st" {
    convention      = "cafclassic"
    name            = "log"
    resource_type   = "azurerm_storage_account"
}

output "st_classic_id" {
  value       = azurecaf_name.classic_st.id
  description = "Id of the resource's name"
}

output "st_classic" {
  value       = azurecaf_name.classic_st.result
  description = "Random result based on the resource type"
}

# Azure Automation Account
resource "azurecaf_name" "azurerm_cognitive_account" {
    name            = "automation"
    resource_type   = "azurerm_cognitive_account"
    prefixes        = ["pre"]
    random          = 4
    suffixes        = ["prod"]
}


output "azurerm_cognitive_account" {
  value       = azurecaf_name.azurerm_cognitive_account.result
  description = "Random result based on the resource type"
}
