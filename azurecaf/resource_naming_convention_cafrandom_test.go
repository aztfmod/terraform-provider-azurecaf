package azurecaf

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccCafNamingConventionCaf_Random(t *testing.T) {
	provider := Provider()
	namingConventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if namingConventionResource == nil {
		t.Fatal("azurecaf_naming_convention resource not found")
	}

	// Test case 1: Storage Account with cafrandom convention
	t.Run("StorageAccountRandom", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "cafrandom",
			"name":          "log",
			"prefix":        "rdmi",
			"resource_type": "st",
		})

		err := namingConventionResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Validate the result contains the prefix
		if !strings.Contains(result, "rdmi") {
			t.Errorf("Expected result to contain 'rdmi', got '%s'", result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["st"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	// Test case 2: Resource Group with special prefix
	t.Run("ResourceGroupSpecialPrefix", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "cafrandom",
			"name":          "myrg",
			"prefix":        "(_124)",
			"resource_type": "rg",
		})

		err := namingConventionResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["rg"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	// Test case 3: Container Registry with max length
	t.Run("ContainerRegistryMaxLength", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "cafrandom",
			"name":          "acrlevel0",
			"prefix":        "rdmi",
			"max_length":    45,
			"resource_type": "acr",
		})

		err := namingConventionResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Check that result respects max length
		if len(result) > 45 {
			t.Errorf("Expected result length <= 45, got %d: '%s'", len(result), result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["acr"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	t.Log("CAF Random naming convention tests completed successfully")
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
# Azure Kubernetes Service
resource "azurecaf_naming_convention" "aksdns" {
    convention      = "cafrandom"
    name            = "kubedemodns"
    prefix          = "rdmi"
    resource_type   = "aksdns"
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
