provider "azurerm" {
    version = "=2.2.0"
    features {}
}

provider "azurecaf" {
}


data "azurerm_client_config" "current" {}

resource "azurecaf_naming_convention" "cafrandom_st" {  
  name    = "aztfmod"
  prefix  = "dev"
  resource_type    = "st"
  postfix = "001"
  max_length = 20
  convention  = "cafrandom"
}

resource "azurecaf_naming_convention" "cafrandom_rg" {  
  name    = "aztfmod"
  prefix  = "dev"
  resource_type    = "rg"
  postfix = "001"
  max_length = 23
  convention  = "cafrandom"
}


resource "azurecaf_naming_convention" "cafrandom_kv" {  
  name    = "mykvaztfmod"
  prefix  = "dev"
  resource_type    = "kv"
  postfix = "01"
  max_length = 23
  convention  = "cafrandom"
}

resource "azurerm_resource_group" "random_rg" {
  name     = azurecaf_naming_convention.cafrandom_rg.result
  location = "southeastasia"
}

resource "azurerm_storage_account" "log" {
  name                      = azurecaf_naming_convention.cafrandom_st.result
  resource_group_name       =  azurerm_resource_group.random_rg.name
  location                  = "southeastasia"
  account_kind              = "StorageV2"
  account_tier              = "Standard"
  account_replication_type  = "GRS"
  access_tier               = "Hot"
  enable_https_traffic_only = true
}
