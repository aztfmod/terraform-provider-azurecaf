package azurecaf

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccCafNamingConventionFull_Random(t *testing.T) {
	provider := Provider()
	namingConventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if namingConventionResource == nil {
		t.Fatal("azurecaf_naming_convention resource not found")
	}

	// Test case 1: Storage Account with random convention
	t.Run("StorageAccountRandom", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"name":          "catest",
			"prefix":        "utest",
			"resource_type": "st",
			"convention":    "random",
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
		if !strings.Contains(result, "utest") {
			t.Errorf("Expected result to contain 'utest', got '%s'", result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["st"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	// Test case 2: Application Gateway with random convention
	t.Run("ApplicationGatewayRandom", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "random",
			"name":          "TEST-DEV-AGW-RG",
			"prefix":        "utest",
			"resource_type": "azurerm_application_gateway",
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
		if !strings.Contains(result, "utest") {
			t.Errorf("Expected result to contain 'utest', got '%s'", result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["agw"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	// Test case 3: API Management with random convention
	t.Run("APIManagementRandom", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "random",
			"name":          "TEST-DEV-APIM-RG",
			"prefix":        "utest",
			"resource_type": "azurerm_api_management",
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
		if !strings.Contains(result, "utest") {
			t.Errorf("Expected result to contain 'utest', got '%s'", result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["apim"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	t.Log("CAF Full Random naming convention tests completed successfully")
}

const testAccResourceRandomConfig = `

#Storage account test
resource "azurecaf_naming_convention" "random_st" {  
	name    		= "catest"
	prefix  		= "utest"
	resource_type   = "st"
	convention  	= "random"
}

# Application Gateway
resource "azurecaf_naming_convention" "random_agw" {
    convention      = "random"
    name            = "TEST-DEV-AGW-RG"
    prefix          = "utest"
    resource_type   = "azurerm_application_gateway"
}

# API Management
resource "azurecaf_naming_convention" "random_apim" {
    convention      = "random"
    name            = "TEST-DEV-APIM-RG"
    prefix          = "utest"
    resource_type   = "azurerm_api_management"
}

# App Service
resource "azurecaf_naming_convention" "random_app" {
    convention      = "random"
    name            = "TEST-DEV-APP-RG"
    prefix          = "utest"
    resource_type   = "azurerm_app_service"
}

# Application Insights
resource "azurecaf_naming_convention" "random_appi" {
    convention      = "random"
    name            = "TEST-DEV-APPI-RG"
    prefix          = "utest"
    resource_type   = "azurerm_application_insights"
}

# Azure Kubernetes Service
resource "azurecaf_naming_convention" "random_aks" {
    convention      = "random"
    name            = "TEST-DEV-AKS-RG"
    prefix          = "utest"
    resource_type   = "azurerm_kubernetes_cluster"
}

# AKS DNS prefix
resource "azurecaf_naming_convention" "random_aksdns" {
    convention      = "random"
    name            = "myaksdnsdemo"
    prefix          = "utest"
    resource_type   = "aks_dns_prefix"
}
# AKS Node Pool Linux
resource "azurecaf_naming_convention" "random_aksnpl" {
    convention      = "random"
    name            = "np1"
    prefix          = "pr"
    resource_type   = "aksnpl"
}
# AKS Node Pool Windows
resource "azurecaf_naming_convention" "random_aksnpw" {
    convention      = "random"
    name            = "np2"
    prefix          = "pr"
    resource_type   = "aksnpw"
}

# App Service Environment
resource "azurecaf_naming_convention" "random_ase" {
    convention      = "random"
    name            = "TEST-DEV-ASE-RG"
    prefix          = "utest"
    resource_type   = "azurerm_app_service_environment"
}

# App Service Plan
resource "azurecaf_naming_convention" "random_plan" {
    convention      = "random"
    name            = "TEST-DEV-PLAN-RG"
    prefix          = "utest"
    resource_type   = "azurerm_app_service_plan"
}

# Azure SQL DB Server
resource "azurecaf_naming_convention" "random_sql" {
    convention      = "random"
    name            = "TEST-DEV-SQL-RG"
    prefix          = "utest"
    resource_type   = "azurerm_sql_server"
}

# Azure SQL DB
resource "azurecaf_naming_convention" "random_sqldb" {
    convention      = "random"
    name            = "TEST-DEV-SQLDB-RG"
    prefix          = "utest"
    resource_type   = "azurerm_sql_database"
}
`
