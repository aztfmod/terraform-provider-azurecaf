package azurecaf

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCafNamingConventionCafRandom(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCafRandomConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"azurecaf_naming_convention.st",
						"log",
						Resources["st"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.st", regexp.MustCompile(Resources["st"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.aaa",
						"automation",
						Resources["aaa"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.aaa", regexp.MustCompile(Resources["aaa"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.acr",
						"registry",
						Resources["acr"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.acr", regexp.MustCompile(Resources["acr"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.acr_max",
						"acrlevel0",
						45,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.acr_max", regexp.MustCompile(Resources["acr"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.rg",
						"myrg",
						Resources["rg"].MaxLength,
						"(_124)"),
					regexMatch("azurecaf_naming_convention.rg", regexp.MustCompile(Resources["rg"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.afw",
						"fire",
						Resources["afw"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.afw", regexp.MustCompile(Resources["afw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.asr",
						"recov",
						Resources["asr"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.asr", regexp.MustCompile(Resources["asr"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.evh",
						"hub",
						Resources["evh"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.evh", regexp.MustCompile(Resources["evh"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.kv",
						"passepartout",
						Resources["kv"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.kv", regexp.MustCompile(Resources["kv"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.aks",
						"kubedemo",
						Resources["aks"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.aks", regexp.MustCompile(Resources["aks"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.la",
						"logs",
						Resources["la"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.la", regexp.MustCompile(Resources["la"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.nic",
						"mynetcard",
						Resources["nic"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.nic", regexp.MustCompile(Resources["nic"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.nsg",
						"sec",
						Resources["nsg"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.nsg", regexp.MustCompile(Resources["nsg"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.pip",
						"mypip",
						Resources["pip"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.pip", regexp.MustCompile(Resources["pip"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.snet",
						"snet",
						Resources["snet"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.snet", regexp.MustCompile(Resources["snet"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.vnet",
						"vnet",
						Resources["vnet"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.vnet", regexp.MustCompile(Resources["vnet"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.vmw",
						"winVMT",
						Resources["vmw"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.vmw", regexp.MustCompile(Resources["vmw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.vml",
						"linuxVM",
						Resources["vml"].MaxLength,
						"rdmi"),
					regexMatch("azurecaf_naming_convention.vml", regexp.MustCompile(Resources["vml"].ValidationRegExp), 1),
				),
			},
		},
	})
}

const testAccResourceCafRandomConfig = `
#Storage account test
resource "azurecaf_naming_convention" "st" {
    convention      = "cafrandom"
    name            = "log"
    prefix          = "rdmi"
    resource_type   = "st"
}

# Azure Automation Account
resource "azurecaf_naming_convention" "aaa" {
    convention      = "cafrandom"
    name            = "automation"
    prefix          = "rdmi"
    resource_type   = "aaa"
}

# Azure Container registry
resource "azurecaf_naming_convention" "acr" {
    convention      = "cafrandom"
    name            = "registry"
    prefix          = "rdmi"
    resource_type   = "acr"
}

resource "azurecaf_naming_convention" "acr_max" {
    convention      = "cafrandom"
    name            = "acrlevel0"
    prefix          = "rdmi"
    max_length      = 45
    resource_type   = "acr"
}

# Resource Group
resource "azurecaf_naming_convention" "rg" {
    convention      = "cafrandom"
    name            = "myrg"
    prefix          = "(_124)"
    resource_type   = "rg"
}

# Azure Firewall
resource "azurecaf_naming_convention" "afw" {
    convention      = "cafrandom"
    name            = "fire"
    prefix          = "rdmi"
    resource_type   = "afw"
}

# Azure Recovery Vault
resource "azurecaf_naming_convention" "asr" {
    convention      = "cafrandom"
    name            = "recov"
    prefix          = "rdmi"
    resource_type   = "asr"
}

# Event Hub
resource "azurecaf_naming_convention" "evh" {
    convention      = "cafrandom"
    name            = "hub"
    prefix          = "rdmi"
    resource_type   = "evh"
}

# Key Vault
resource "azurecaf_naming_convention" "kv" {
    convention      = "cafrandom"
    name            = "passepartout"
    prefix          = "rdmi"
    resource_type   = "kv"
}

# Azure Kubernetes Service
resource "azurecaf_naming_convention" "aks" {
    convention      = "cafrandom"
    name            = "kubedemo"
    prefix          = "rdmi"
    resource_type   = "aks"
}

# Log Analytics Workspace
resource "azurecaf_naming_convention" "la" {
    convention      = "cafrandom"
    name            = "logs"
    prefix          = "rdmi"
    resource_type   = "la"
}

# Network Interface
resource "azurecaf_naming_convention" "nic" {
    convention      = "cafrandom"
    name            = "mynetcard"
    prefix          = "rdmi"
    resource_type   = "nic"
}

# Network Security Group
resource "azurecaf_naming_convention" "nsg" {
    convention      = "cafrandom"
    name            = "sec"
    prefix          = "rdmi"
    resource_type   = "nsg"
}

# Public Ip
resource "azurecaf_naming_convention" "pip" {
    convention      = "cafrandom"
    name            = "mypip"
    prefix          = "rdmi"
    resource_type   = "pip"
}

# subnet
resource "azurecaf_naming_convention" "snet" {
    convention      = "cafrandom"
    name            = "snet"
    prefix          = "rdmi"
    resource_type   = "snet"
}

# Virtual Network
resource "azurecaf_naming_convention" "vnet" {
    convention      = "cafrandom"
    name            = "vnet"
    prefix          = "rdmi"
    resource_type   = "vnet"
}

# VM Windows
resource "azurecaf_naming_convention" "vmw" {
    convention      = "cafrandom"
    name            = "winVMToolongShouldbetrimmed"
    prefix          = "rdmi"
    resource_type   = "vmw"
}

# VM Linux
resource "azurecaf_naming_convention" "vml" {
    convention      = "cafrandom"
    name            = "linuxVM"
    prefix          = "rdmi"
    resource_type   = "vml"
}
`
