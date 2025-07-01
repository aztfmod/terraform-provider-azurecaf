package azurecaf

import (
	"testing"
)

func TestAccCafNamingConvention_Classic(t *testing.T) {
	testCases := []NamingConventionTestCase{
		{
			Name:             "log",
			Convention:       "cafclassic",
			ResourceType:     "st",
			ExpectedContains: []string{"st"},
		},
		{
			Name:             "automation",
			Convention:       "cafclassic",
			ResourceType:     "aaa",
			ExpectedContains: []string{"aaa"},
		},
		{
			Name:             "registry",
			Convention:       "cafclassic",
			ResourceType:     "acr",
			ExpectedContains: []string{"acr"},
		},
		{
			Name:             "myrg",
			Convention:       "cafclassic",
			ResourceType:     "rg",
			ExpectedContains: []string{"rg"},
		},
		{
			Name:             "passepartout",
			Convention:       "cafclassic",
			ResourceType:     "kv",
			ExpectedContains: []string{"kv"},
		},
		{
			Name:             "fire",
			Convention:       "cafclassic",
			ResourceType:     "afw",
			ExpectedContains: []string{"afw"},
		},
		{
			Name:             "recov",
			Convention:       "cafclassic",
			ResourceType:     "asr",
			ExpectedContains: []string{"asr"},
		},
		{
			Name:             "hub",
			Convention:       "cafclassic",
			ResourceType:     "evh",
			ExpectedContains: []string{"evh"},
		},
		{
			Name:             "kubedemo",
			Convention:       "cafclassic",
			ResourceType:     "aks",
			ExpectedContains: []string{"aks"},
		},
		{
			Name:             "kubedemodns",
			Convention:       "cafclassic",
			ResourceType:     "aksdns",
			ExpectedContains: []string{"aksdns"},
		},
	}

	runMultipleNamingConventionTests(t, testCases)
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
`
