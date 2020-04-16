# Log Analytics Workspace
resource "azurecaf_naming_convention" "la_passthrough" {
    convention      = "passthrough"
    name            = "logs_invalid"
    prefix          = "rdmi"
    resource_type   = "la"
}

output "la_passthrough_id" {
  value       = azurecaf_naming_convention.la_passthrough.id
  description = "Id of the resource's name"
}

output "la_passthrough_random" {
  value       = azurecaf_naming_convention.la_passthrough.result
  description = "Random result based on the resource type"
}

# Resource Group
resource "azurecaf_naming_convention" "passthrough_rg" {
    convention      = "passthrough"
    name            = "TEST-DEV-ASE-RG"
    prefix          = "rdmi"
    resource_type   = "azurerm_resource_group"
}

output "rg_passthrough_id" {
  value       = azurecaf_naming_convention.passthrough_rg.id
  description = "Id of the resource's name"
}

output "rg_passthrough" {
  value       = azurecaf_naming_convention.passthrough_rg.result
  description = "Random result based on the resource type"
}

# Application Gateway
resource "azurecaf_naming_convention" "passthrough_agw" {
    convention      = "passthrough"
    name            = "TEST-DEV-AGW-RG"
    prefix          = "rdmi"
    resource_type   = "azurerm_application_gateway"
}

output "agw_passthrough_id" {
  value       = azurecaf_naming_convention.passthrough_agw.id
  description = "Id of the resource's name"
}

output "agw_passthrough" {
  value       = azurecaf_naming_convention.passthrough_agw.result
  description = "Random result based on the resource type"
}

# API Management
resource "azurecaf_naming_convention" "passthrough_apim" {
    convention      = "passthrough"
    name            = "TEST-DEV-APIM-RG"
    prefix          = "rdmi"
    resource_type   = "azurerm_api_management"
}

output "apim_passthrough_id" {
  value       = azurecaf_naming_convention.passthrough_apim.id
  description = "Id of the resource's name"
}

output "apim_passthrough" {
  value       = azurecaf_naming_convention.passthrough_apim.result
  description = "Random result based on the resource type"
}

# App Service
resource "azurecaf_naming_convention" "passthrough_app" {
    convention      = "passthrough"
    name            = "TEST-DEV-APP-RG"
    prefix          = "rdmi"
    resource_type   = "azurerm_app_service"
}

output "app_passthrough_id" {
  value       = azurecaf_naming_convention.passthrough_app.id
  description = "Id of the resource's name"
}

output "app_passthrough" {
  value       = azurecaf_naming_convention.passthrough_app.result
  description = "Random result based on the resource type"
}

# Application Insights
resource "azurecaf_naming_convention" "passthrough_appi" {
    convention      = "passthrough"
    name            = "TEST-DEV-APPI-RG"
    prefix          = "rdmi"
    resource_type   = "azurerm_application_insights"
}

output "appi_passthrough_id" {
  value       = azurecaf_naming_convention.passthrough_appi.id
  description = "Id of the resource's name"
}

output "appi_passthrough" {
  value       = azurecaf_naming_convention.passthrough_appi.result
  description = "Random result based on the resource type"
}

# App Service Environment
resource "azurecaf_naming_convention" "passthrough_ase" {
    convention      = "passthrough"
    name            = "TEST-DEV-ASE-RG"
    prefix          = "rdmi"
    resource_type   = "azurerm_app_service_environment"
}

output "ase_passthrough_id" {
  value       = azurecaf_naming_convention.passthrough_ase.id
  description = "Id of the resource's name"
}

output "ase_passthrough" {
  value       = azurecaf_naming_convention.passthrough_ase.result
  description = "Random result based on the resource type"
}

# App Service Plan
resource "azurecaf_naming_convention" "passthrough_plan" {
    convention      = "passthrough"
    name            = "TEST-DEV-PLAN-RG"
    prefix          = "rdmi"
    resource_type   = "azurerm_app_service_plan"
}

output "plan_passthrough_id" {
  value       = azurecaf_naming_convention.passthrough_plan.id
  description = "Id of the resource's name"
}

output "plan_passthrough" {
  value       = azurecaf_naming_convention.passthrough_plan.result
  description = "Random result based on the resource type"
}

# Azure SQL DB Server
resource "azurecaf_naming_convention" "passthrough_sql" {
    convention      = "passthrough"
    name            = "TEST-DEV-SQL-RG"
    prefix          = "rdmi"
    resource_type   = "azurerm_sql_server"
}

output "sql_passthrough_id" {
  value       = azurecaf_naming_convention.passthrough_sql.id
  description = "Id of the resource's name"
}

output "sql_passthrough" {
  value       = azurecaf_naming_convention.passthrough_sql.result
  description = "Random result based on the resource type"
}

# Azure SQL DB
resource "azurecaf_naming_convention" "passthrough_sqldb" {
    convention      = "passthrough"
    name            = "TEST-DEV-SQLDB-RG"
    prefix          = "rdmi"
    resource_type   = "azurerm_sql_database"
}

output "sqldb_passthrough_id" {
  value       = azurecaf_naming_convention.passthrough_sqldb.id
  description = "Id of the resource's name"
}

output "sqldb_passthrough" {
  value       = azurecaf_naming_convention.passthrough_sqldb.result
  description = "Random result based on the resource type"
}
