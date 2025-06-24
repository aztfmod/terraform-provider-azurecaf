package azurecaf

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCafNamingConventionFull_Random(t *testing.T) {
	// Skip this test if we can't access external network resources
	// This test requires Terraform CLI which needs to connect to checkpoint-api.hashicorp.com
	t.Skip("Skipping acceptance test - requires network access to Terraform CLI")
	
	resource.UnitTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRandomConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_st",
						"",
						Resources["st"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_st", regexp.MustCompile(Resources["st"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_agw",
						"",
						Resources["agw"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_agw", regexp.MustCompile(Resources["agw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_apim",
						"",
						Resources["apim"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_apim", regexp.MustCompile(Resources["apim"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_app",
						"",
						Resources["app"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_app", regexp.MustCompile(Resources["app"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_appi",
						"",
						Resources["appi"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_appi", regexp.MustCompile(Resources["appi"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_aks",
						"",
						Resources["aks"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_aks", regexp.MustCompile(Resources["aks"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_aksdns",
						"",
						Resources["aksdns"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_aksdns", regexp.MustCompile(Resources["aksdns"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_aksnpl",
						"",
						Resources["aksnpl"].MaxLength,
						"pr"),
					regexMatch("azurecaf_naming_convention.random_aksnpl", regexp.MustCompile(Resources["aksnpl"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_aksnpw",
						"",
						Resources["aksnpw"].MaxLength,
						"pr"),
					regexMatch("azurecaf_naming_convention.random_aksnpl", regexp.MustCompile(Resources["aksnpl"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_aksnpw",
						"",
						Resources["aksnpw"].MaxLength,
						"pr"),
					regexMatch("azurecaf_naming_convention.random_aksnpw", regexp.MustCompile(Resources["aksnpw"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_ase",
						"",
						Resources["ase"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_ase", regexp.MustCompile(Resources["ase"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_plan",
						"",
						Resources["plan"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_plan", regexp.MustCompile(Resources["plan"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_sql",
						"",
						Resources["sql"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_sql", regexp.MustCompile(Resources["sql"].ValidationRegExp), 1),
					testAccCafNamingValidation(
						"azurecaf_naming_convention.random_sqldb",
						"",
						Resources["sqldb"].MaxLength,
						"utest"),
					regexMatch("azurecaf_naming_convention.random_sqldb", regexp.MustCompile(Resources["sqldb"].ValidationRegExp), 1),
				),
			},
		},
	})
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
