provider "caf" {

}

resource "caf_naming_convention" "nc" {
    convention      = "cafclassic"
    name            = "log"
    prefix          = "rdmi-"
    resource_type   = "rg"
}

output "Id" {
  value       = caf_naming_convention.nc.id
  description = "Id of the resource's name"
}