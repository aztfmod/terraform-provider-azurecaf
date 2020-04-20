package azurecaf

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCafNamingConventionFullRandom(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
provider "azurecaf" {

}

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
