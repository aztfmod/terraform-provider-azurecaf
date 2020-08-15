
#Storage account test
resource "azurecaf_naming_convention" "classic_st" {
    convention      = "cafclassic"
    name            = "log"
    resource_type   = "st"
}

output "st_classic_id" {
  value       = azurecaf_naming_convention.classic_st.id
  description = "Id of the resource's name"
}

output "st_classic" {
  value       = azurecaf_naming_convention.classic_st.result
  description = "Random result based on the resource type"
}

# Azure Automation Account
resource "azurecaf_naming_convention" "classic_aaa" {
    convention      = "cafclassic"
    name            = "automation"
    resource_type   = "aaa"
}

output "aaa_classic_id" {
  value       = azurecaf_naming_convention.classic_aaa.id
  description = "Id of the resource's name"
}

output "aaa_classic" {
  value       = azurecaf_naming_convention.classic_aaa.result
  description = "Random result based on the resource type"
}


# Azure Container registry
resource "azurecaf_naming_convention" "classic_acr" {
    convention      = "cafclassic"
    name            = "registry"
    resource_type   = "acr"
}

output "acr_classic_id" {
  value       = azurecaf_naming_convention.classic_acr.id
  description = "Id of the resource's name"
}

output "acr_classic" {
  value       = azurecaf_naming_convention.classic_acr.result
  description = "Random result based on the resource type"
}

# Resource Group
resource "azurecaf_naming_convention" "classic_rg" {
    convention      = "cafclassic"
    name            = "myrg"
    resource_type   = "rg"
}

output "rg_classic_id" {
  value       = azurecaf_naming_convention.classic_rg.id
  description = "Id of the resource's name"
}

output "rg_classic" {
  value       = azurecaf_naming_convention.classic_rg.result
  description = "Random result based on the resource type"
}

# Azure Firewall
resource "azurecaf_naming_convention" "classic_afw" {
    convention      = "cafclassic"
    name            = "fire"
    resource_type   = "afw"
}

output "afw_classic_id" {
  value       = azurecaf_naming_convention.classic_afw.id
  description = "Id of the resource's name"
}

output "afw_classic" {
  value       = azurecaf_naming_convention.classic_afw.result
  description = "Random result based on the resource type"
}

# Azure Recovery Vault
resource "azurecaf_naming_convention" "classic_asr" {
    convention      = "cafclassic"
    name            = "recov"
    resource_type   = "asr"
}

output "asr_classic_id" {
  value       = azurecaf_naming_convention.classic_asr.id
  description = "Id of the resource's name"
}

output "asr_classic" {
  value       = azurecaf_naming_convention.classic_asr.result
  description = "Random result based on the resource type"
}


# Event Hub
resource "azurecaf_naming_convention" "classic_evh" {
    convention      = "cafclassic"
    name            = "hub"
    resource_type   = "evh"
}

output "evh_classic_id" {
  value       = azurecaf_naming_convention.classic_evh.id
  description = "Id of the resource's name"
}

output "evh_classic" {
  value       = azurecaf_naming_convention.classic_evh.result
  description = "Random result based on the resource type"
}

# Key Vault
resource "azurecaf_naming_convention" "classic_kv" {
    convention      = "cafclassic"
    name            = "passepartout"
    resource_type   = "kv"
}

output "kv_classic_id" {
  value       = azurecaf_naming_convention.classic_kv.id
  description = "Id of the resource's name"
}

output "kv_classic" {
  value       = azurecaf_naming_convention.classic_kv.result
  description = "Random result based on the resource type"
}

# Log Analytics Workspace
resource "azurecaf_naming_convention" "classic_la" {
    convention      = "cafclassic"
    name            = "logs"
    resource_type   = "la"
}

output "la_classic_id" {
  value       = azurecaf_naming_convention.classic_la.id
  description = "Id of the resource's name"
}

output "la_classic" {
  value       = azurecaf_naming_convention.classic_la.result
  description = "Random result based on the resource type"
}

# Network Interface
resource "azurecaf_naming_convention" "classic_nic" {
    convention      = "cafclassic"
    name            = "mynetcard"
    resource_type   = "nic"
}

output "nic_classic_id" {
  value       = azurecaf_naming_convention.classic_nic.id
  description = "Id of the resource's name"
}

output "nic_classic" {
  value       = azurecaf_naming_convention.classic_nic.result
  description = "Random result based on the resource type"
}

# Network Security Group
resource "azurecaf_naming_convention" "classic_nsg" {
    convention      = "cafclassic"
    name            = "sec"
    resource_type   = "nsg"
}

output "nsg_classic_id" {
  value       = azurecaf_naming_convention.classic_nsg.id
  description = "Id of the resource's name"
}

output "nsg_classic" {
  value       = azurecaf_naming_convention.classic_nsg.result
  description = "Random result based on the resource type"
}

# Public Ip
resource "azurecaf_naming_convention" "classic_pip" {
    convention      = "cafclassic"
    name            = "mypip"
    resource_type   = "pip"
}

output "pip_classic_id" {
  value       = azurecaf_naming_convention.classic_pip.id
  description = "Id of the resource's name"
}

output "pip_classic" {
  value       = azurecaf_naming_convention.classic_pip.result
  description = "Random result based on the resource type"
}

# subnet
resource "azurecaf_naming_convention" "classic_snet" {
    convention      = "cafclassic"
    name            = "snet"
    resource_type   = "snet"
}

output "snet_classic_id" {
  value       = azurecaf_naming_convention.classic_snet.id
  description = "Id of the resource's name"
}

output "snet_classic" {
  value       = azurecaf_naming_convention.classic_snet.result
  description = "Random result based on the resource type"
}

# Virtual Network
resource "azurecaf_naming_convention" "classic_vnet" {
    convention      = "cafclassic"
    name            = "vnet"
    resource_type   = "vnet"
}

output "vnet_classic_id" {
  value       = azurecaf_naming_convention.classic_vnet.id
  description = "Id of the resource's name"
}

output "vnet_classic" {
  value       = azurecaf_naming_convention.classic_vnet.result
  description = "Random result based on the resource type"
}

# VM Windows
resource "azurecaf_naming_convention" "classic_vmw" {
    convention      = "cafclassic"
    name            = "winVMToolongShouldbetrimmed"
    resource_type   = "vmw"
}

output "vmw_classic_id" {
  value       = azurecaf_naming_convention.classic_vmw.id
  description = "Id of the resource's name"
}

output "vmw_classic" {
  value       = azurecaf_naming_convention.classic_vmw.result
  description = "Random result based on the resource type"
}

# VM Linux
resource "azurecaf_naming_convention" "classic_vml" {
    convention      = "cafclassic"
    name            = "linuxVM"
    resource_type   = "vml"
}

output "vml_classic_id" {
  value       = azurecaf_naming_convention.classic_vml.id
  description = "Id of the resource's name"
}

output "vml_classic" {
  value       = azurecaf_naming_convention.classic_vml.result
  description = "Random result based on the resource type"
}

#Application Security Group test
resource "azurecaf_naming_convention" "classic_asg" {
    convention      = "cafclassic"
    name            = "AppSecGroup"
    resource_type   = "asg"
}

output "asg_classic_id" {
  value       = azurecaf_naming_convention.classic_asg.id
  description = "Id of the resource's name"
}

output "asg_classic" {
  value       = azurecaf_naming_convention.classic_asg.result
  description = "Random result based on the resource type"
}

#Azure VPN Connection
resource "azurecaf_naming_convention" "classic_cn" {
    convention      = "cafclassic"
    name            = "My_VPN_Connection_"
    resource_type   = "cn"
}

output "cn_classic_id" {
  value       = azurecaf_naming_convention.classic_cn.id
  description = "Id of the resource's name"
}

output "cn_classic" {
  value       = azurecaf_naming_convention.classic_cn.result
  description = "Random result based on the resource type"
}

#Azure Load Balancer (external)
resource "azurecaf_naming_convention" "classic_lbe" {
    convention      = "cafclassic"
    name            = "My_External.Load.Balancer_"
    resource_type   = "lbe"
}

output "lbe_classic_id" {
  value       = azurecaf_naming_convention.classic_lbe.id
  description = "Id of the resource's name"
}

output "lbe_classic" {
  value       = azurecaf_naming_convention.classic_lbe.result
  description = "Random result based on the resource type"
}

#Azure Load Balancer (internal)
resource "azurecaf_naming_convention" "classic_lbi" {
    convention      = "cafclassic"
    name            = "My_Internal.Load.Balancer_"
    resource_type   = "lbi"
}

output "lbi_classic_id" {
  value       = azurecaf_naming_convention.classic_lbi.id
  description = "Id of the resource's name"
}

output "lbi_classic" {
  value       = azurecaf_naming_convention.classic_lbi.result
  description = "Random result based on the resource type"
}

#Azure Local Network Gateway
resource "azurecaf_naming_convention" "classic_lgw" {
    convention      = "cafclassic"
    name            = "My_Local.Network.Gateway_"
    resource_type   = "lgw"
}

output "lgw_classic_id" {
  value       = azurecaf_naming_convention.classic_lgw.id
  description = "Id of the resource's name"
}

output "lgw_classic" {
  value       = azurecaf_naming_convention.classic_lgw.result
  description = "Random result based on the resource type"
}

#Azure Mysql Database
resource "azurecaf_naming_convention" "classic_mysql" {
    convention      = "cafclassic"
    name            = "My-MySQL-Database-001"
    resource_type   = "mysql"
}

output "mysql_classic_id" {
  value       = azurecaf_naming_convention.classic_mysql.id
  description = "Id of the resource's name"
}

output "mysql_classic" {
  value       = azurecaf_naming_convention.classic_mysql.result
  description = "Random result based on the resource type"
}

#Azure Route Table
resource "azurecaf_naming_convention" "classic_route" {
    convention      = "cafclassic"
    name            = "My-Route.Table-001_"
    resource_type   = "route"
}

output "route_classic_id" {
  value       = azurecaf_naming_convention.classic_route.id
  description = "Id of the resource's name"
}

output "route_classic" {
  value       = azurecaf_naming_convention.classic_route.result
  description = "Random result based on the resource type"
}

#Azure Service Bus
resource "azurecaf_naming_convention" "classic_sb" {
    convention      = "cafclassic"
    name            = "My-Service.Bus-001_"
    resource_type   = "sb"
}

output "sb_classic_id" {
  value       = azurecaf_naming_convention.classic_sb.id
  description = "Id of the resource's name"
}

output "sb_classic" {
  value       = azurecaf_naming_convention.classic_sb.result
  description = "Random result based on the resource type"
}

#Azure Service queue
resource "azurecaf_naming_convention" "classic_sbq" {
    convention      = "cafclassic"
    name            = "My-Service.Bus/001/queue/001"
    resource_type   = "sbq"
}

output "sbq_classic_id" {
  value       = azurecaf_naming_convention.classic_sbq.id
  description = "Id of the resource's name"
}

output "sbq_classic" {
  value       = azurecaf_naming_convention.classic_sbq.result
  description = "Random result based on the resource type"
}

#Azure Service Topic
resource "azurecaf_naming_convention" "classic_sbt" {
    convention      = "cafclassic"
    name            = "My-Service.Bus/001/queue/001/topic/001"
    resource_type   = "sbt"
}

output "sbt_classic_id" {
  value       = azurecaf_naming_convention.classic_sbt.id
  description = "Id of the resource's name"
}

output "sbt_classic" {
  value       = azurecaf_naming_convention.classic_sbt.result
  description = "Random result based on the resource type"
}

#Azure Virtual Network Gateway
resource "azurecaf_naming_convention" "classic_vgw" {
    convention      = "cafclassic"
    name            = "My-Virtual.Network.Gateway-001_"
    resource_type   = "vgw"
}

output "vgw_classic_id" {
  value       = azurecaf_naming_convention.classic_vgw.id
  description = "Id of the resource's name"
}

output "vgw_classic" {
  value       = azurecaf_naming_convention.classic_vgw.result
  description = "Random result based on the resource type"
}

#Azure Availability Set
resource "azurecaf_naming_convention" "classic_avail" {
    convention      = "cafclassic"
    name            = "My-AvailabilitySet-001_"
    resource_type   = "avail"
}

output "avail_classic_id" {
  value       = azurecaf_naming_convention.classic_avail.id
  description = "Id of the resource's name"
}

output "avail_classic" {
  value       = azurecaf_naming_convention.classic_avail.result
  description = "Random result based on the resource type"
}

#Azure Traffic Manager Profile
resource "azurecaf_naming_convention" "classic_traf" {
    convention      = "cafclassic"
    name            = "My-Traffic.Manager-001_"
    resource_type   = "traf"
}

output "traf_classic_id" {
  value       = azurecaf_naming_convention.classic_traf.id
  description = "Id of the resource's name"
}

output "traf_classic" {
  value       = azurecaf_naming_convention.classic_traf.result
  description = "Random result based on the resource type"
}

#Azure VM Scale Set Linux
resource "azurecaf_naming_convention" "classic_vmssl" {
    convention      = "cafclassic"
    name            = "My-VM-ScaleSet-001_"
    resource_type   = "vmssl"
}

output "vmssl_classic_id" {
  value       = azurecaf_naming_convention.classic_vmssl.id
  description = "Id of the resource's name"
}

output "vmssl_classic" {
  value       = azurecaf_naming_convention.classic_vmssl.result
  description = "Random result based on the resource type"
}

#Azure VM Scale Set Windows
resource "azurecaf_naming_convention" "classic_vmssw" {
    convention      = "cafclassic"
    name            = "VMScaleSet001"
    resource_type   = "vmssw"
}

output "vmssw_classic_id" {
  value       = azurecaf_naming_convention.classic_vmssw.id
  description = "Id of the resource's name"
}

output "vmssw_classic" {
  value       = azurecaf_naming_convention.classic_vmssw.result
  description = "Random result based on the resource type"
}