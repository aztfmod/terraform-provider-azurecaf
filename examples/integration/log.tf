
#Create the resource groups to host the blueprint
resource "azurecaf_naming_convention" "pz_ops_name" {  
  name    = "demo_rg_log"
  resource_type    = "rg"
  max_length = 50
  convention  = "passthrough"
}

resource "azurerm_resource_group" "pz_ops" {
  name     = azurecaf_naming_convention.pz_ops_name.result
  location = "southeastasia"
}

locals {
     solution_plan_map = {
        NetworkMonitoring = {
            "publisher" = "Microsoft"
            "product"   = "OMSGallery/NetworkMonitoring"
        },
        ADAssessment = {
            "publisher" = "Microsoft"
            "product"   = "OMSGallery/ADAssessment"
        },
        ADReplication = {
            "publisher" = "Microsoft"
            "product"   = "OMSGallery/ADReplication"
        },
        AgentHealthAssessment = {
            "publisher" = "Microsoft"
            "product"   = "OMSGallery/AgentHealthAssessment"
        },
        DnsAnalytics = {
            "publisher" = "Microsoft"
            "product"   = "OMSGallery/DnsAnalytics"
        }
    }
}


#Create the Azure Monitor - Log Analytics workspace
module "log_analytics" {
  # source  = "aztfmod/caf-log-analytics/azurerm"
  # version = "1.0.0"
  source = "git://github.com/aztfmod/terraform-azurerm-caf-log-analytics?ref=v2.0.1"

  convention          = "passthrough"
  prefix              = ""
  name                = "validname"
  solution_plan_map   = local.solution_plan_map
  resource_group_name = azurerm_resource_group.pz_ops.name
  location            = "southeastasia"
  tags                = {}
}
