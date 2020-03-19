package caf

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCafNamingConventionRandom(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCafConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"caf_naming_convention.st",
						"log",
						Resources["st"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.st", regexp.MustCompile(Resources["st"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.aaa",
						"automation",
						Resources["aaa"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.aaa", regexp.MustCompile(Resources["aaa"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.acr",
						"registry",
						Resources["acr"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.acr", regexp.MustCompile(Resources["acr"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.rg",
						"myrg",
						Resources["rg"].MaxLength,
						"(_124)"),
					regexMatch("caf_naming_convention.rg", regexp.MustCompile(Resources["rg"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.afw",
						"fire",
						Resources["afw"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.afw", regexp.MustCompile(Resources["afw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.asr",
						"recov",
						Resources["asr"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.asr", regexp.MustCompile(Resources["asr"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.evh",
						"hub",
						Resources["evh"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.evh", regexp.MustCompile(Resources["evh"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.kv",
						"passepartout",
						Resources["kv"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.kv", regexp.MustCompile(Resources["kv"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.la",
						"logs",
						Resources["la"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.la", regexp.MustCompile(Resources["la"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.nic",
						"mynetcard",
						Resources["nic"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.nic", regexp.MustCompile(Resources["nic"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.nsg",
						"sec",
						Resources["nsg"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.nsg", regexp.MustCompile(Resources["nsg"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.pip",
						"mypip",
						Resources["pip"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.pip", regexp.MustCompile(Resources["pip"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.snet",
						"snet",
						Resources["snet"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.snet", regexp.MustCompile(Resources["snet"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.vnet",
						"vnet",
						Resources["vnet"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.vnet", regexp.MustCompile(Resources["vnet"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.vmw",
						"winVMT",
						Resources["vmw"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.vmw", regexp.MustCompile(Resources["vmw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.vml",
						"linuxVM",
						Resources["vml"].MaxLength,
						"rdmi"),
					regexMatch("caf_naming_convention.vml", regexp.MustCompile(Resources["vml"].ValidationRegExp), 1),
				),
			},
		},
	})
}
func TestAccCafNamingConventionClassic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCafClassicConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"caf_naming_convention.classic_st",
						"log",
						5,
						"st"),
					regexMatch("caf_naming_convention.classic_st", regexp.MustCompile(Resources["st"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_aaa",
						"automation",
						14,
						"aaa"),
					regexMatch("caf_naming_convention.classic_aaa", regexp.MustCompile(Resources["aaa"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_acr",
						"registry",
						11,
						"acr"),
					regexMatch("caf_naming_convention.classic_acr", regexp.MustCompile(Resources["acr"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_rg",
						"myrg",
						7,
						"rg"),
					regexMatch("caf_naming_convention.classic_rg", regexp.MustCompile(Resources["rg"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_afw",
						"fire",
						8,
						"afw"),
					regexMatch("caf_naming_convention.classic_afw", regexp.MustCompile(Resources["afw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_asr",
						"recov",
						9,
						"asr"),
					regexMatch("caf_naming_convention.classic_asr", regexp.MustCompile(Resources["asr"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_evh",
						"hub",
						7,
						"evh"),
					regexMatch("caf_naming_convention.classic_evh", regexp.MustCompile(Resources["evh"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_kv",
						"passepartout",
						15,
						"kv"),
					regexMatch("caf_naming_convention.classic_kv", regexp.MustCompile(Resources["kv"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_la",
						"logs",
						7,
						"la"),
					regexMatch("caf_naming_convention.classic_la", regexp.MustCompile(Resources["la"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_nic",
						"mynetcard",
						13,
						"nic"),
					regexMatch("caf_naming_convention.classic_nic", regexp.MustCompile(Resources["nic"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_nsg",
						"sec",
						7,
						"nsg"),
					regexMatch("caf_naming_convention.classic_nsg", regexp.MustCompile(Resources["nsg"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_pip",
						"mypip",
						9,
						"pip"),
					regexMatch("caf_naming_convention.classic_pip", regexp.MustCompile(Resources["pip"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_snet",
						"snet",
						9,
						"snet"),
					regexMatch("caf_naming_convention.classic_snet", regexp.MustCompile(Resources["snet"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_vnet",
						"vnet",
						9,
						"vnet"),
					regexMatch("caf_naming_convention.classic_vnet", regexp.MustCompile(Resources["vnet"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_vmw",
						"winVMT",
						15,
						"vmw"),
					regexMatch("caf_naming_convention.classic_vmw", regexp.MustCompile(Resources["vmw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"caf_naming_convention.classic_vml",
						"linuxVM",
						11,
						"vml"),
					regexMatch("caf_naming_convention.classic_vml", regexp.MustCompile(Resources["vml"].ValidationRegExp), 1),
				),
			},
		},
	})
}
func TestAccCafNamingConventionPassthrough(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePassthroughConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"caf_naming_convention.pass_st",
						"log",
						3,
						"log"),
					regexMatch("caf_naming_convention.pass_st", regexp.MustCompile(Resources["st"].ValidationRegExp), 1),
				),
			},
		},
	})
}

func TestAccCafNamingConventionFullRandom(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRandomConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"caf_naming_convention.random_st",
						"",
						24,
						""),
					regexMatch("caf_naming_convention.random_st", regexp.MustCompile(Resources["st"].ValidationRegExp), 1),
				),
			},
		},
	})
}

func testAccCafNamingValidation(id string, name string, expectedLength int, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		attrs := rs.Primary.Attributes

		result := attrs["result"]
		if len(result) != expectedLength {
			return fmt.Errorf("got %s %d result items; want %d", result, len(result), expectedLength)
		}
		if !strings.HasPrefix(result, prefix) {
			return fmt.Errorf("got %s which doesn't start with %s", result, prefix)
		}
		if !strings.Contains(result, name) {
			return fmt.Errorf("got %s which doesn't contain the name %s", result, name)
		}
		return nil
	}
}

func regexMatch(id string, exp *regexp.Regexp, requiredMatches int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		result := rs.Primary.Attributes["result"]

		if matches := exp.FindAllStringSubmatchIndex(result, -1); len(matches) != requiredMatches {
			return fmt.Errorf("result string is %s; did not match %s, got %d", result, exp, len(matches))
		}

		return nil
	}
}

const testAccResourceCafConfig = `
#Storage account test
resource "caf_naming_convention" "st" {
    convention      = "cafrandom"
    name            = "log"
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
    prefix          = "(_124)"
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
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
    name            = "mypip"
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
    name            = "winVMToolongShouldbetrimmed"
    prefix          = "rdmi"
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
    prefix          = "rdmi"
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
`
const testAccResourceCafClassicConfig = `
provider "caf" {

}


#Storage account test
resource "caf_naming_convention" "classic_st" {
    convention      = "cafclassic"
    name            = "log"
    resource_type   = "st"
}

output "st_classic_id" {
  value       = caf_naming_convention.classic_st.id
  description = "Id of the resource's name"
}

output "st_classic" {
  value       = caf_naming_convention.classic_st.result
  description = "Random result based on the resource type"
}

# Azure Automation Account
resource "caf_naming_convention" "classic_aaa" {
    convention      = "cafclassic"
    name            = "automation"
    resource_type   = "aaa"
}

output "aaa_classic_id" {
  value       = caf_naming_convention.classic_aaa.id
  description = "Id of the resource's name"
}

output "aaa_classic" {
  value       = caf_naming_convention.classic_aaa.result
  description = "Random result based on the resource type"
}


# Azure Container registry
resource "caf_naming_convention" "classic_acr" {
    convention      = "cafclassic"
    name            = "registry"
    resource_type   = "acr"
}

output "acr_classic_id" {
  value       = caf_naming_convention.classic_acr.id
  description = "Id of the resource's name"
}

output "acr_classic" {
  value       = caf_naming_convention.classic_acr.result
  description = "Random result based on the resource type"
}

# Resource Group
resource "caf_naming_convention" "classic_rg" {
    convention      = "cafclassic"
    name            = "myrg"
    resource_type   = "rg"
}

output "rg_classic_id" {
  value       = caf_naming_convention.classic_rg.id
  description = "Id of the resource's name"
}

output "rg_classic" {
  value       = caf_naming_convention.classic_rg.result
  description = "Random result based on the resource type"
}

# Azure Firewall
resource "caf_naming_convention" "classic_afw" {
    convention      = "cafclassic"
    name            = "fire"
    resource_type   = "afw"
}

output "afw_classic_id" {
  value       = caf_naming_convention.classic_afw.id
  description = "Id of the resource's name"
}

output "afw_classic" {
  value       = caf_naming_convention.classic_afw.result
  description = "Random result based on the resource type"
}

# Azure Recovery Vault
resource "caf_naming_convention" "classic_asr" {
    convention      = "cafclassic"
    name            = "recov"
    resource_type   = "asr"
}

output "asr_classic_id" {
  value       = caf_naming_convention.classic_asr.id
  description = "Id of the resource's name"
}

output "asr_classic" {
  value       = caf_naming_convention.classic_asr.result
  description = "Random result based on the resource type"
}


# Event Hub
resource "caf_naming_convention" "classic_evh" {
    convention      = "cafclassic"
    name            = "hub"
    resource_type   = "evh"
}

output "evh_classic_id" {
  value       = caf_naming_convention.classic_evh.id
  description = "Id of the resource's name"
}

output "evh_classic" {
  value       = caf_naming_convention.classic_evh.result
  description = "Random result based on the resource type"
}

# Key Vault
resource "caf_naming_convention" "classic_kv" {
    convention      = "cafclassic"
    name            = "passepartout"
    resource_type   = "kv"
}

output "kv_classic_id" {
  value       = caf_naming_convention.classic_kv.id
  description = "Id of the resource's name"
}

output "kv_classic" {
  value       = caf_naming_convention.classic_kv.result
  description = "Random result based on the resource type"
}

# Log Analytics Workspace
resource "caf_naming_convention" "classic_la" {
    convention      = "cafclassic"
    name            = "logs"
    resource_type   = "la"
}

output "la_classic_id" {
  value       = caf_naming_convention.classic_la.id
  description = "Id of the resource's name"
}

output "la_classic" {
  value       = caf_naming_convention.classic_la.result
  description = "Random result based on the resource type"
}

# Network Interface
resource "caf_naming_convention" "classic_nic" {
    convention      = "cafclassic"
    name            = "mynetcard"
    resource_type   = "nic"
}

output "nic_classic_id" {
  value       = caf_naming_convention.classic_nic.id
  description = "Id of the resource's name"
}

output "nic_classic" {
  value       = caf_naming_convention.classic_nic.result
  description = "Random result based on the resource type"
}

# Network Security Group
resource "caf_naming_convention" "classic_nsg" {
    convention      = "cafclassic"
    name            = "sec"
    resource_type   = "nsg"
}

output "nsg_classic_id" {
  value       = caf_naming_convention.classic_nsg.id
  description = "Id of the resource's name"
}

output "nsg_classic" {
  value       = caf_naming_convention.classic_nsg.result
  description = "Random result based on the resource type"
}

# Public Ip
resource "caf_naming_convention" "classic_pip" {
    convention      = "cafclassic"
    name            = "mypip"
    resource_type   = "pip"
}

output "pip_classic_id" {
  value       = caf_naming_convention.classic_pip.id
  description = "Id of the resource's name"
}

output "pip_classic" {
  value       = caf_naming_convention.classic_pip.result
  description = "Random result based on the resource type"
}

# subnet
resource "caf_naming_convention" "classic_snet" {
    convention      = "cafclassic"
    name            = "snet"
    resource_type   = "snet"
}

output "snet_classic_id" {
  value       = caf_naming_convention.classic_snet.id
  description = "Id of the resource's name"
}

output "snet_classic" {
  value       = caf_naming_convention.classic_snet.result
  description = "Random result based on the resource type"
}

# Virtual Network
resource "caf_naming_convention" "classic_vnet" {
    convention      = "cafclassic"
    name            = "vnet"
    resource_type   = "vnet"
}

output "vnet_classic_id" {
  value       = caf_naming_convention.classic_vnet.id
  description = "Id of the resource's name"
}

output "vnet_classic" {
  value       = caf_naming_convention.classic_vnet.result
  description = "Random result based on the resource type"
}

# VM Windows
resource "caf_naming_convention" "classic_vmw" {
    convention      = "cafclassic"
    name            = "winVMToolongShouldbetrimmed"
    resource_type   = "vmw"
}

output "vmw_classic_id" {
  value       = caf_naming_convention.classic_vmw.id
  description = "Id of the resource's name"
}

output "vmw_classic" {
  value       = caf_naming_convention.classic_vmw.result
  description = "Random result based on the resource type"
}

# VM Linux
resource "caf_naming_convention" "classic_vml" {
    convention      = "cafclassic"
    name            = "linuxVM"
    resource_type   = "vml"
}

output "vml_classic_id" {
  value       = caf_naming_convention.classic_vml.id
  description = "Id of the resource's name"
}

output "vml_classic" {
  value       = caf_naming_convention.classic_vml.result
  description = "Random result based on the resource type"
}
`

const testAccResourcePassthroughConfig = `
provider "caf" {

}

#Storage account test
resource "caf_naming_convention" "pass_st" {
    convention      = "passthrough"
    name            = "log"
    resource_type   = "st"
}
`

const testAccResourceRandomConfig = `
provider "caf" {

}

#Storage account test
resource "caf_naming_convention" "random_st" {
    convention      = "random"
    name            = "log"
    resource_type   = "st"
}
`
