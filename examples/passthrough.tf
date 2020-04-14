# Log Analytics Workspace
resource "azurecaf_naming_convention" "la_passthrough" {
    convention      = "passthrough"
    name            = "logs_invalid"
    prefix          = "rdmi"
    resource_type   = "la"
}

output "la_passthrough_id" {
  value       = azurecaf_naming_convention.la_passthrough.id
  description = "Id of the resource's name"
}

output "la_passthrough_random" {
  value       = azurecaf_naming_convention.la_passthrough.result
  description = "Random result based on the resource type"
}
