
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

