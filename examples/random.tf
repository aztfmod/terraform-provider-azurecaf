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