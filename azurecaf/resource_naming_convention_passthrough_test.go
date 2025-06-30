package azurecaf

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccCafNamingConvention_Passthrough(t *testing.T) {
	provider := Provider()
	namingConventionResource := provider.ResourcesMap["azurecaf_naming_convention"]
	if namingConventionResource == nil {
		t.Fatal("azurecaf_naming_convention resource not found")
	}

	// Test case 1: Log Analytics with invalid characters (should clean)
	t.Run("LogAnalyticsInvalid", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "passthrough",
			"name":          "logs_invalid",
			"resource_type": "la",
		})

		err := namingConventionResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expected := "logsinvalid" // underscores should be cleaned
		if result != expected {
			t.Errorf("Expected result '%s', got '%s'", expected, result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["la"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	// Test case 2: Application Gateway passthrough
	t.Run("ApplicationGatewayPassthrough", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "passthrough",
			"name":          "TEST-DEV-AGW-RG",
			"resource_type": "azurerm_application_gateway",
		})

		err := namingConventionResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		expected := "TEST-DEV-AGW-RG"
		if result != expected {
			t.Errorf("Expected result '%s', got '%s'", expected, result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["agw"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	// Test case 3: API Management passthrough
	t.Run("APIManagementPassthrough", func(t *testing.T) {
		resourceData := schema.TestResourceDataRaw(t, namingConventionResource.Schema, map[string]interface{}{
			"convention":    "passthrough",
			"name":          "TEST-DEV-APIM-RG",
			"resource_type": "azurerm_api_management",
		})

		err := namingConventionResource.Create(resourceData, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		result := resourceData.Get("result").(string)
		// Should be cleaned for APIM which doesn't allow hyphens
		expected := "TESTDEVAPIMRG"
		if result != expected {
			t.Errorf("Expected result '%s', got '%s'", expected, result)
		}

		// Validate against Azure naming requirements if Resources map exists
		if resource, exists := Resources["apim"]; exists && resource.ValidationRegExp != "" {
			if !regexp.MustCompile(resource.ValidationRegExp).MatchString(result) {
				t.Errorf("Result '%s' does not match Azure naming requirements", result)
			}
		}
	})

	t.Log("CAF Passthrough naming convention tests completed successfully")
}

const testAccResourcePassthroughConfig = `
#Storage account test
resource "azurecaf_naming_convention" "logs_inv" {
    convention      = "passthrough"
    name            = "logs_invalid"
    resource_type   = "la"
}

# Application Gateway
resource "azurecaf_naming_convention" "passthrough_agw" {
    convention      = "passthrough"
    name            = "TEST-DEV-AGW-RG"
    resource_type   = "azurerm_application_gateway"
}

# API Management
resource "azurecaf_naming_convention" "passthrough_apim" {
    convention      = "passthrough"
    name            = "TEST-DEV-APIM-RG"
    resource_type   = "azurerm_api_management"
}

# App Service
resource "azurecaf_naming_convention" "passthrough_app" {
    convention      = "passthrough"
    name            = "TEST-DEV-APP-RG"
    resource_type   = "azurerm_app_service"
}

# Application Insights
resource "azurecaf_naming_convention" "passthrough_appi" {
    convention      = "passthrough"
    name            = "TEST-DEV-APPI-RG"
    resource_type   = "azurerm_application_insights"
}

# Azure Kubernetes Services
resource "azurecaf_naming_convention" "passthrough_aks" {
    convention      = "passthrough"
    name            = "kubedemo"
    resource_type   = "azurerm_kubernetes_cluster"
}

# Azure Kubernetes Services DNS Prefix
resource "azurecaf_naming_convention" "passthrough_aksdns" {
    convention      = "passthrough"
    name            = "kubedemodns"
    resource_type   = "aksdns"
}

# Azure Kubernetes Services Node pool Linux 
resource "azurecaf_naming_convention" "passthrough_aksnpl" {
    convention      = "passthrough"
    name            = "knplinux"
    resource_type   = "aksnpl"
}

# Azure Kubernetes Services Node Pool Windows
resource "azurecaf_naming_convention" "passthrough_aksnpw" {
    convention      = "passthrough"
    name            = "knpwindows" #expecting 6 chars
    resource_type   = "aksnpw"
}

# App Service Environment
resource "azurecaf_naming_convention" "passthrough_ase" {
    convention      = "passthrough"
    name            = "TEST-DEV-ASE-RG"
    resource_type   = "azurerm_app_service_environment"
}

# App Service Plan
resource "azurecaf_naming_convention" "passthrough_plan" {
    convention      = "passthrough"
    name            = "TEST-DEV-PLAN-RG"
    resource_type   = "azurerm_app_service_plan"
}

# Azure SQL DB Server
resource "azurecaf_naming_convention" "passthrough_sql" {
    convention      = "passthrough"
    name            = "TEST-DEV-SQL-RG"
    resource_type   = "azurerm_sql_server"
}

# Azure SQL DB
resource "azurecaf_naming_convention" "passthrough_sqldb" {
    convention      = "passthrough"
    name            = "TEST-DEV-SQLDB-RG"
    resource_type   = "azurerm_sql_database"
}
`
