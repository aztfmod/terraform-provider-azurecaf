# azurecaf_naming_convention

!> **WARNING:** This resource is deprecated. Please use [`azurecaf_name`](azurecaf_name.md) instead, which provides more flexibility and supports a broader range of Azure resources.

The resource naming_convention implements a set of methodologies to apply consistent resource naming using the default Microsoft Cloud Adoption Framework for Azure recommendations as per https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging.

The naming_convention is the initial resource released as part of the azurecaf provider, the naming_convention supports a fixed set of resources as described in the documention. In order to provider more flexibility and support the large breadth of Azure resources available you can use the azurecaf_name resource.

## Example usage
This example outputs one name, the result of the naming convention query. The result attribute returns the name based on the convention and parameters input.

The example generates a 23 characters name compatible with the specification for an Azure Resource Group
dev-aztfmod-001

```hcl
resource "azurecaf_naming_convention" "cafrandom_rg" {  
  name    = "aztfmod"
  prefix  = "dev"
  resource_type    = "rg"
  postfix = "001"
  max_length = 23
  convention  = "cafrandom"
}

resource "azurerm_resource_group" "cafrandom" {
  name     = azurecaf_naming_convention.cafrandom_rg.result
  location = "southeastasia"
}


The provider generates a name using the input parameters and automatically appends a prefix (if defined), a caf prefix (resource type) and postfix (if defined) in addition to a generated padding string based on the selected naming convention.

```
The example above would generate a name using the pattern [prefix]-[cafprefix]-[name]-[postfix]-[padding]:

```
dev-aztfmod-rg-001-wxyz
```

## Argument Reference

The following arguments are supported:

* name - (optional) the basename of the resource to create, the basename will be sanitized as per supported character set in Azure.
* convention (optional) - one of the four naming convention supported. Defaults to cafrandom. Allowed values are cafclassic, cafrandom, random, passthrough
* prefix (optional) - prefix to append as the first characters of the generated name
* postfix (optional) -  additional postfix added after the basename, this is can be used to append resource index (eg. vm-001)
* max_length (optional) - configure the maximum length of the returned object name, is the specified length is longer than the supported length of the Azure resource the later applies
* resource_type (optional) -  describes the type of azure resource you are requesting a name from (eg. azure container registry: acr). See the Resource Type section

# Attributes Reference
The following attributes are exported:

* id - The id of the naming convention object
* result - The generated named for an Azure Resource based on the input parameter and the selected naming convention


# Methods for naming convention

The following methods are implemented for naming conventions:

| method name | description of the naming convention used |
| -- | -- |
| cafclassic | follows Cloud Adoption Framework for Azure recommendations as per https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging |
| cafrandom | follows Cloud Adoption Framework for Azure recommendations as per https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging and adds randomly generated characters up to maximum length of name |
| random | name will be generated automatically in full lengths of azure object |
| passthrough | naming convention is implemented manually, fields given as input will be same as the output (but lengths and forbidden chars will be filtered out) |

## Migration to azurecaf_name

To migrate from `azurecaf_naming_convention` to `azurecaf_name`, you need to update your configuration. Here's how the parameters map between the old and new resources:

### Parameter Mapping

| azurecaf_naming_convention | azurecaf_name | Notes |
|----------------------------|---------------|-------|
| `name` | `name` | Same parameter |
| `prefix` | `prefixes` | Single string becomes a list |
| `postfix` | `suffixes` | Single string becomes a list |
| `resource_type` | `resource_type` | Use the azurerm resource type (see mapping table below) |
| `convention` | N/A | Replaced by other parameters: use `random_length` for random characters |
| `max_length` | `resource_types` | Length validation is automatic based on resource type |

### Migration Example

**Old configuration:**
```hcl
resource "azurecaf_naming_convention" "cafrandom_rg" {  
  name          = "aztfmod"
  prefix        = "dev"
  resource_type = "rg"
  postfix       = "001"
  max_length    = 23
  convention    = "cafrandom"
}
```

**New configuration:**
```hcl
resource "azurecaf_name" "rg_example" {
  name          = "aztfmod"
  prefixes      = ["dev"]
  resource_type = "azurerm_resource_group"
  suffixes      = ["001"]
  random_length = 5  # If you want random characters
}
```

## Resource types

We define resource types as per: https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging 

Current prototype supports:

| Resource type                       | Resource type code (short)  | Resource type code (long)               |
| ----------------------------------- | ----------------------------|-----------------------------------------|
| Azure Automation                    | aaa                         | azurerm_automation_account              |
| Azure Container App                 | ac                          | azurerm_container_app                   |
| Azure Container App Environment     | ace                         | azurerm_container_app_environment       |
| Azure Container Registry            | acr                         | azurerm_container_registry              |
| Azure Firewall                      | afw                         | azurerm_firewall                        |
| Application Gateway                 | agw                         | azurerm_application_gateway             |
| API Management                      | apim                        | azurerm_api_management                  |
| App Service                         | app                         | azurerm_app_service                     |
| Application Insights                | appi                        | azurerm_application_insights            |
| App Service Environment             | ase                         | azurerm_app_service_environment         |
| Azure Kubernetes Service            | aks                         | azurerm_kubernetes_cluster              |
| Azure Kubernetes Service DNS prefix | aksdns                      | aks_dns_prefix                          |
| AKS Node Pool Linux                 | aksnpl                      | aks_node_pool_linux                     |
| AKS Node Pool Windows               | aksnpw                      | aks_node_pool_windows                   |
| Azure Site Recovery                 | asr                         | azurerm_recovery_services_vault         |
| Azure Event Hubs                    | evh                         | azurerm_eventhub_namespace              |
| generic                             | gen                         | generic                                 |
| Azure Key Vault                     | kv                          | azurerm_key_vault                       |
| Azure Monitor Log Analytics         | la                          | azurerm_log_analytics_workspace         |
| Virtual Network Interface Card      | nic                         | azurerm_network_interface               |
| Network Security Group              | nsg                         | azurerm_network_security_group          |
| Public IP                           | pip                         | azurerm_public_ip                       |
| App Service Plan                    | plan                        | azurerm_app_service_plan                |
| Service Plan                        | plan                        | azurerm_service_plan                    |
| Resource group                      | rg                          | azurerm_resource_group                  |
| Subnet                              | snet                        | azurerm_subnet                          |
| Azure SQL DB Server                 | sql                         | azurerm_sql_server                      |
| Azure SQL DB                        | sqldb                       | azurerm_sql_database                    |
| Azure Storage Account               | st                          | azurerm_storage_account                 |
| Linux Virtual Machine               | vml                         | azurerm_virtual_machine_linux           |
| Windows Virtual Machine             | vmw                         | azurerm_virtual_machine_windows         |
| Virtual Network                     | vnet                        | azurerm_virtual_network                 |
