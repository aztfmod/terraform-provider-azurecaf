provider "caf" {

}

resource "caf_naming_convention" "nc" {
    name            = "log"
    prefix          = "rdmi"
    resource_type   = "st"
}