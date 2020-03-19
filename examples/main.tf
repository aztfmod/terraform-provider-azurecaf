provider "caf" {

}

#Storage account test
resource "caf_naming_convention" "st" {
    convention      = "cafrandom"
    name            = "log"
    prefix          = "rdmi-"
    resource_type   = "st"
}

output "st_id" {
  value       = caf_naming_convention.st.id
  description = "Id of the resource's name"
}

output "st_random" {
  value       = caf_naming_convention.st.result
  description = "Random result based on the resource type"
}

# Azure Automation Account
resource "caf_naming_convention" "aaa" {
    convention      = "cafrandom"
    name            = "automation"
    prefix          = "rdmi-"
    resource_type   = "aaa"
}

output "aaa_id" {
  value       = caf_naming_convention.aaa.id
  description = "Id of the resource's name"
}

output "aaa_random" {
  value       = caf_naming_convention.aaa.result
  description = "Random result based on the resource type"
}


# Azure Container registry
resource "caf_naming_convention" "acr" {
    convention      = "cafrandom"
    name            = "registry"
    prefix          = "rdmi-"
    resource_type   = "acr"
}

output "acr_id" {
  value       = caf_naming_convention.acr.id
  description = "Id of the resource's name"
}

output "acr_random" {
  value       = caf_naming_convention.acr.result
  description = "Random result based on the resource type"
}

# Resource Group
resource "caf_naming_convention" "rg" {
    convention      = "cafrandom"
    name            = "myrg"
    prefix          = "(_124)-"
    resource_type   = "rg"
}

output "rg_id" {
  value       = caf_naming_convention.rg.id
  description = "Id of the resource's name"
}

output "rg_random" {
  value       = caf_naming_convention.rg.result
  description = "Random result based on the resource type"
}

# Azure Firewall
resource "caf_naming_convention" "afw" {
    convention      = "cafrandom"
    name            = "fire"
    prefix          = "rdmi-"
    resource_type   = "afw"
}

output "afw_id" {
  value       = caf_naming_convention.afw.id
  description = "Id of the resource's name"
}

output "afw_random" {
  value       = caf_naming_convention.afw.result
  description = "Random result based on the resource type"
}

# Azure Recovery Vault
resource "caf_naming_convention" "asr" {
    convention      = "cafrandom"
    name            = "recov"
    prefix          = "rdmi-"
    resource_type   = "asr"
}

output "asr_id" {
  value       = caf_naming_convention.asr.id
  description = "Id of the resource's name"
}

output "asr_random" {
  value       = caf_naming_convention.asr.result
  description = "Random result based on the resource type"
}


# Event Hub
resource "caf_naming_convention" "evh" {
    convention      = "cafrandom"
    name            = "hub"
    prefix          = "rdmi-"
    resource_type   = "evh"
}

output "evh_id" {
  value       = caf_naming_convention.evh.id
  description = "Id of the resource's name"
}

output "evh_random" {
  value       = caf_naming_convention.evh.result
  description = "Random result based on the resource type"
}

# Key Vault
resource "caf_naming_convention" "kv" {
    convention      = "cafrandom"
    name            = "passepartout"
    prefix          = "rdmi-"
    resource_type   = "kv"
}

output "kv_id" {
  value       = caf_naming_convention.kv.id
  description = "Id of the resource's name"
}

output "kv_random" {
  value       = caf_naming_convention.kv.result
  description = "Random result based on the resource type"
}

# Log Analytics Workspace
resource "caf_naming_convention" "la" {
    convention      = "cafrandom"
    name            = "logs"
    prefix          = "rdmi-"
    resource_type   = "la"
}

output "la_id" {
  value       = caf_naming_convention.la.id
  description = "Id of the resource's name"
}

output "la_random" {
  value       = caf_naming_convention.la.result
  description = "Random result based on the resource type"
}

# Network Interface
resource "caf_naming_convention" "nic" {
    convention      = "cafrandom"
    name            = "mynetcard"
    prefix          = "rdmi-"
    resource_type   = "nic"
}

output "nic_id" {
  value       = caf_naming_convention.nic.id
  description = "Id of the resource's name"
}

output "nic_random" {
  value       = caf_naming_convention.nic.result
  description = "Random result based on the resource type"
}

# Network Security Group
resource "caf_naming_convention" "nsg" {
    convention      = "cafrandom"
    name            = "sec"
    prefix          = "rdmi-"
    resource_type   = "nsg"
}

output "nsg_id" {
  value       = caf_naming_convention.nsg.id
  description = "Id of the resource's name"
}

output "nsg_random" {
  value       = caf_naming_convention.nsg.result
  description = "Random result based on the resource type"
}

# Public Ip
resource "caf_naming_convention" "pip" {
    convention      = "cafrandom"
    name            = "pip"
    prefix          = "rdmi-"
    resource_type   = "pip"
}

output "pip_id" {
  value       = caf_naming_convention.pip.id
  description = "Id of the resource's name"
}

output "pip_random" {
  value       = caf_naming_convention.pip.result
  description = "Random result based on the resource type"
}

# subnet
resource "caf_naming_convention" "snet" {
    convention      = "cafrandom"
    name            = "snet"
    prefix          = "rdmi-"
    resource_type   = "snet"
}

output "snet_id" {
  value       = caf_naming_convention.snet.id
  description = "Id of the resource's name"
}

output "snet_random" {
  value       = caf_naming_convention.snet.result
  description = "Random result based on the resource type"
}

# Virtual Network
resource "caf_naming_convention" "vnet" {
    convention      = "cafrandom"
    name            = "vnet"
    prefix          = "rdmi-"
    resource_type   = "vnet"
}

output "vnet_id" {
  value       = caf_naming_convention.vnet.id
  description = "Id of the resource's name"
}

output "vnet_random" {
  value       = caf_naming_convention.vnet.result
  description = "Random result based on the resource type"
}

# VM Windows
resource "caf_naming_convention" "vmw" {
    convention      = "cafrandom"
    name            = "winVM"
    prefix          = "rdmi-"
    resource_type   = "vmw"
}

output "vmw_id" {
  value       = caf_naming_convention.vmw.id
  description = "Id of the resource's name"
}

output "vmw_random" {
  value       = caf_naming_convention.vmw.result
  description = "Random result based on the resource type"
}

# VM Linux
resource "caf_naming_convention" "vml" {
    convention      = "cafrandom"
    name            = "linuxVM"
    prefix          = "rdmi-"
    resource_type   = "vml"
}

output "vml_id" {
  value       = caf_naming_convention.vml.id
  description = "Id of the resource's name"
}

output "vml_random" {
  value       = caf_naming_convention.vml.result
  description = "Random result based on the resource type"
}

