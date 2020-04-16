package azurecaf

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCafNamingConventionClassic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCafClassicConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_st",
						"log",
						5,
						"st"),
					regexMatch("azurecaf_naming_convention.classic_st", regexp.MustCompile(Resources["st"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_aaa",
						"automation",
						14,
						"aaa"),
					regexMatch("azurecaf_naming_convention.classic_aaa", regexp.MustCompile(Resources["aaa"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_acr",
						"registry",
						11,
						"acr"),
					regexMatch("azurecaf_naming_convention.classic_acr", regexp.MustCompile(Resources["acr"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_rg",
						"myrg",
						7,
						"rg"),
					regexMatch("azurecaf_naming_convention.classic_rg", regexp.MustCompile(Resources["rg"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_afw",
						"fire",
						8,
						"afw"),
					regexMatch("azurecaf_naming_convention.classic_afw", regexp.MustCompile(Resources["afw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_asr",
						"recov",
						9,
						"asr"),
					regexMatch("azurecaf_naming_convention.classic_asr", regexp.MustCompile(Resources["asr"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_evh",
						"hub",
						7,
						"evh"),
					regexMatch("azurecaf_naming_convention.classic_evh", regexp.MustCompile(Resources["evh"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_kv",
						"passepartout",
						15,
						"kv"),
					regexMatch("azurecaf_naming_convention.classic_kv", regexp.MustCompile(Resources["kv"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_la",
						"logs",
						7,
						"la"),
					regexMatch("azurecaf_naming_convention.classic_la", regexp.MustCompile(Resources["la"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_nic",
						"mynetcard",
						13,
						"nic"),
					regexMatch("azurecaf_naming_convention.classic_nic", regexp.MustCompile(Resources["nic"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_nsg",
						"sec",
						7,
						"nsg"),
					regexMatch("azurecaf_naming_convention.classic_nsg", regexp.MustCompile(Resources["nsg"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_pip",
						"mypip",
						9,
						"pip"),
					regexMatch("azurecaf_naming_convention.classic_pip", regexp.MustCompile(Resources["pip"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_snet",
						"snet",
						9,
						"snet"),
					regexMatch("azurecaf_naming_convention.classic_snet", regexp.MustCompile(Resources["snet"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_vnet",
						"vnet",
						9,
						"vnet"),
					regexMatch("azurecaf_naming_convention.classic_vnet", regexp.MustCompile(Resources["vnet"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_vmw",
						"winVMT",
						15,
						"vmw"),
					regexMatch("azurecaf_naming_convention.classic_vmw", regexp.MustCompile(Resources["vmw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_vml",
						"linuxVM",
						11,
						"vml"),
					regexMatch("azurecaf_naming_convention.classic_vml", regexp.MustCompile(Resources["vml"].ValidationRegExp), 1),
				),
			},
		},
	})
}

const testAccResourceCafClassicConfig = `
provider "azurecaf" {

}


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
`
