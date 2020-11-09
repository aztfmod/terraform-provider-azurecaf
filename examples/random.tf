resource "azurecaf_naming_convention" "random_st2" {  
	name    = "catest"
	prefix  = "test"
	resource_type    = "st"
	convention  = "random"
  }

  output "random_st2_id" {
  value       = azurecaf_naming_convention.random_st2.id
  description = "Id of the resource's name"
}

output "random_st2_random" {
  value       = azurecaf_naming_convention.random_st2.result
  description = "Random result based on the resource type"
}

resource "azurecaf_naming_convention" "random_st_fullname" {  
	name    = "catest"
	prefix  = "test"
	resource_type    = "azurerm_storage_account"
	convention  = "random"
  }

  output "random_st_fullname_id" {
  value       = azurecaf_naming_convention.random_st_fullname.id
  description = "Id of the resource's name"
}

output "random_st_fullname_random" {
  value       = azurecaf_naming_convention.random_st_fullname.result
  description = "Random result based on the resource type"
}

#Event hub Random
resource "azurecaf_naming_convention" "evh_fullrandom" {
    name            = "evhrandseed"
    convention      = "random"
    resource_type   = "azurerm_eventhub_namespace"
}

output "evh_fullrandom_id" {
  value       = azurecaf_naming_convention.evh_fullrandom.id
  description = "Id of the resource's name"
}

output "evh_fullrandom" {
  value       = azurecaf_naming_convention.evh_fullrandom.result
  description = "Random result based on the resource type"
}
