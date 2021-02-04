package azurecaf

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCafNamingConvention_Passthrough(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePassthroughConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"azurecaf_naming_convention.logs_inv",
						"logsinvalid",
						11,
						"log"),
					regexMatch("azurecaf_naming_convention.logs_inv", regexp.MustCompile(Resources["la"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_agw",
						"TEST-DEV-AGW-RG",
						15,
						"TEST"),
					regexMatch("azurecaf_naming_convention.passthrough_agw", regexp.MustCompile(Resources["agw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_apim",
						"TESTDEVAPIMRG",
						13,
						"TEST"),
					regexMatch("azurecaf_naming_convention.passthrough_apim", regexp.MustCompile(Resources["apim"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_app",
						"TEST-DEV-APP-RG",
						15,
						"TEST"),
					regexMatch("azurecaf_naming_convention.passthrough_app", regexp.MustCompile(Resources["app"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_appi",
						"TEST-DEV-APPI-RG",
						16,
						"TEST"),
					regexMatch("azurecaf_naming_convention.passthrough_appi", regexp.MustCompile(Resources["appi"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_aks",
						"kubedemo",
						8,
						"kube"),
					regexMatch("azurecaf_naming_convention.passthrough_aks", regexp.MustCompile(Resources["aks"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_aksdns",
						"kubedemodns",
						11,
						"kube"),
					regexMatch("azurecaf_naming_convention.passthrough_aksdns", regexp.MustCompile(Resources["aksdns"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_aksnpl",
						"knplinux",
						8,
						"knp"),
					regexMatch("azurecaf_naming_convention.passthrough_aksnpl", regexp.MustCompile(Resources["aksnpl"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_aksnpw",
						"knpwin",
						6,
						"knp"),
					regexMatch("azurecaf_naming_convention.passthrough_aksnpw", regexp.MustCompile(Resources["aksnpw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_ase",
						"TEST-DEV-ASE-RG",
						15,
						"TEST"),
					regexMatch("azurecaf_naming_convention.passthrough_ase", regexp.MustCompile(Resources["ase"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_plan",
						"TEST-DEV-PLAN-RG",
						16,
						"TEST"),
					regexMatch("azurecaf_naming_convention.passthrough_plan", regexp.MustCompile(Resources["plan"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_sql",
						"test-dev-sql-rg",
						15,
						"test"),
					regexMatch("azurecaf_naming_convention.passthrough_sql", regexp.MustCompile(Resources["sql"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.passthrough_sqldb",
						"TEST-DEV-SQLDB-RG",
						17,
						"TEST"),
					regexMatch("azurecaf_naming_convention.passthrough_sqldb", regexp.MustCompile(Resources["sqldb"].ValidationRegExp), 1),
				),
			},
		},
	})
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
