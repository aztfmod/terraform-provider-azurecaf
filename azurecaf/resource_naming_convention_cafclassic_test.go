package azurecaf

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccCafNamingConvention_Classic(t *testing.T) {
	provider := Provider()
	namingConventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if namingConventionResource == nil {
		t.Fatal("azurecaf_naming_convention resource not found")
	}

	// Test case 1: Storage Account
	t.Run("StorageAccount", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "cafclassic",
			"name":          "log",
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

		// Validate the result contains the prefix for storage account
		if !strings.Contains(result, "st") {
			t.Errorf("Expected result to contain 'st', got '%s'", result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["st"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	// Test case 2: Azure Automation Account
	t.Run("AutomationAccount", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "cafclassic",
			"name":          "automation",
			"resource_type": "aaa",
		})

		err := namingConventionResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Validate the result contains the prefix for automation account
		if !strings.Contains(result, "aaa") {
			t.Errorf("Expected result to contain 'aaa', got '%s'", result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["aaa"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	// Test case 3: Container Registry
	t.Run("ContainerRegistry", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "cafclassic",
			"name":          "registry",
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

		// Validate the result contains the prefix for container registry
		if !strings.Contains(result, "acr") {
			t.Errorf("Expected result to contain 'acr', got '%s'", result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["acr"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	// Test case 4: Resource Group
	t.Run("ResourceGroup", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "cafclassic",
			"name":          "myrg",
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

		// Validate the result contains the prefix for resource group
		if !strings.Contains(result, "rg") {
			t.Errorf("Expected result to contain 'rg', got '%s'", result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["rg"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	// Test case 5: Key Vault
	t.Run("KeyVault", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "cafclassic",
			"name":          "passepartout",
			"resource_type": "kv",
		})

		err := namingConventionResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Validate the result contains the prefix for key vault
		if !strings.Contains(result, "kv") {
			t.Errorf("Expected result to contain 'kv', got '%s'", result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["kv"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	t.Log("CAF Classic naming convention tests completed successfully")
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
