package azurecaf

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCafNamingConvention_Classic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
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
						"azurecaf_naming_convention.classic_aks",
						"kubedemo",
						11,
						"aks"),
					regexMatch("azurecaf_naming_convention.classic_aks", regexp.MustCompile(Resources["aks"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_aksdns",
						"kubedemodns",
						18,
						"aksdns"),
					regexMatch("azurecaf_naming_convention.classic_aksdns", regexp.MustCompile(Resources["aksdns"].ValidationRegExp), 1),
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
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_asg",
						"AppSecGroup",
						15,
						"asg"),
					regexMatch("azurecaf_naming_convention.classic_asg", regexp.MustCompile(Resources["asg"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.classic_cn",
						"My_VPN_Connection_",
						21,
						"cn"),
					regexMatch("azurecaf_naming_convention.classic_cn", regexp.MustCompile(Resources["cn"].ValidationRegExp), 1),
				),
			},
		},
	})
}

const testAccResourceCafClassicConfig = `

#Storage account test
resource "azurecaf_naming_convention" "classic_st" {
    convention      = "cafclassic"
    name            = "log"
    resource_type   = "st"
}

# Azure Automation Account
resource "azurecaf_naming_convention" "classic_aaa" {
    convention      = "cafclassic"
    name            = "automation"
    resource_type   = "aaa"
}

# Azure Container registry
resource "azurecaf_naming_convention" "classic_acr" {
    convention      = "cafclassic"
    name            = "registry"
    resource_type   = "acr"
}

# Resource Group
resource "azurecaf_naming_convention" "classic_rg" {
    convention      = "cafclassic"
    name            = "myrg"
    resource_type   = "rg"
}

# Azure Firewall
resource "azurecaf_naming_convention" "classic_afw" {
    convention      = "cafclassic"
    name            = "fire"
    resource_type   = "afw"
}

# Azure Recovery Vault
resource "azurecaf_naming_convention" "classic_asr" {
    convention      = "cafclassic"
    name            = "recov"
    resource_type   = "asr"
}

# Event Hub
resource "azurecaf_naming_convention" "classic_evh" {
    convention      = "cafclassic"
    name            = "hub"
    resource_type   = "evh"
}

# Key Vault
resource "azurecaf_naming_convention" "classic_kv" {
    convention      = "cafclassic"
    name            = "passepartout"
    resource_type   = "kv"
}

# Azure Kubernetes Service
resource "azurecaf_naming_convention" "classic_aks" {
    convention      = "cafclassic"
    name            = "kubedemo"
    resource_type   = "aks"
}
# Azure Kubernetes Service
resource "azurecaf_naming_convention" "classic_aksdns" {
    convention      = "cafclassic"
    name            = "kubedemodns"
    resource_type   = "aksdns"
}

# Log Analytics Workspace
resource "azurecaf_naming_convention" "classic_la" {
    convention      = "cafclassic"
    name            = "logs"
    resource_type   = "la"
}

# Network Interface
resource "azurecaf_naming_convention" "classic_nic" {
    convention      = "cafclassic"
    name            = "mynetcard"
    resource_type   = "nic"
}

# Network Security Group
resource "azurecaf_naming_convention" "classic_nsg" {
    convention      = "cafclassic"
    name            = "sec"
    resource_type   = "nsg"
}

# Public Ip
resource "azurecaf_naming_convention" "classic_pip" {
    convention      = "cafclassic"
    name            = "mypip"
    resource_type   = "pip"
}

# subnet
resource "azurecaf_naming_convention" "classic_snet" {
    convention      = "cafclassic"
    name            = "snet"
    resource_type   = "snet"
}

# Virtual Network
resource "azurecaf_naming_convention" "classic_vnet" {
    convention      = "cafclassic"
    name            = "vnet"
    resource_type   = "vnet"
}

# VM Windows
resource "azurecaf_naming_convention" "classic_vmw" {
    convention      = "cafclassic"
    name            = "winVMToolongShouldbetrimmed"
    resource_type   = "vmw"
}

# VM Linux
resource "azurecaf_naming_convention" "classic_vml" {
    convention      = "cafclassic"
    name            = "linuxVM"
    resource_type   = "vml"
}

#Application Security Group test
resource "azurecaf_naming_convention" "classic_asg" {
    convention      = "cafclassic"
    name            = "AppSecGroup"
    resource_type   = "asg"
}

#Azure VPN Connection test
resource "azurecaf_naming_convention" "classic_cn" {
    convention      = "cafclassic"
    name            = "My_VPN_Connection_"
    resource_type   = "cn"
}

`
