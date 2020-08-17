[![VScodespaces](https://img.shields.io/endpoint?url=https%3A%2F%2Faka.ms%2Fvso-badge)](https://online.visualstudio.com/environments/new?name=caf%20terraform%20provider&repo=aztfmod/terraform-provider-azurecaf)

# Naming convention

This provider implements a set of methodologies for naming convention implementation including the default Microsoft Cloud Adoption Framework for Azure recommendations as per https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging.

# Building the  Provider

Clone repository to: $GOPATH/src/github.com/aztfmod/terraform-provider-azurecaf

```
$ mkdir -p $GOPATH/src/github.com/aztfmod; cd $GOPATH/src/github.com/aztfmod
$ git clone https://github.com/aztfmod/terraform-provider-azurecaf.git

```
Enter the provider directory and build the provider

```
$ cd $GOPATH/src/github.com/aztfmod/terraform-provider-azurecaf
$ make build

```

# Using the Provider

If you're building the provider, follow the [terraform instructions](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) to install it as a plugin. After placing it into your plugins directory, run terraform init to initialize it.

# Developing the Provider

If you wish to work on the provider, you'll first need Go installed on your machine (version 1.13+ is required). You'll also need to correctly setup a GOPATH, as well as adding $GOPATH/bin to your $PATH.

To compile the provider, run make build. This will build the provider and put the provider binary in the $GOPATH/bin directory.

```
$ make build
...
$ $GOPATH/bin/terraform-provider-azurecaf
...

```
# Testing

Running the acceptance test suite requires does not require an Azure subscription. 

to run the unit test:
```
make unittest
```

to run the integration test

```
make test
```

## Methods for naming convention

The following methods are implemented for naming conventions:

| method name | description of the naming convention used |
| -- | -- |
| cafclassic | follows Cloud Adoption Framework for Azure recommendations as per https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging |
| cafrandom | follows Cloud Adoption Framework for Azure recommendations as per https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging and adds randomly generated characters up to maximum length of name |
| random | name will be generated automatically in full lengths of azure object |
| passthrough | naming convention is implemented manually, fields given as input will be same as the output (but lengths and forbidden chars will be filtered out) |

## Resource types

We define resource types as per: https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging 

Current prototype supports:

| Resource type                       | Resource type code (short)  | Resource type code (long)               |
| ----------------------------------- | ----------------------------|-----------------------------------------|
| Azure Automation                    | aaa                         | azurerm_automation_account              |
| Azure Container Registry            | acr                         | azurerm_container_registry              |
| Azure Firewall                      | afw                         | azurerm_firewall                        |
| Application Gateway                 | agw                         | azurerm_application_gateway             |
| API Management                      | apim                        | azurerm_api_management                  |
| App Service                         | app                         | azurerm_app_service                     |
| Application Insights                | appi                        | azurerm_application_insights            |
| App Service Environment             | ase                         | azurerm_app_service_environment         |
| Application Security Group          | asg                         | azurerm_app_security_group              |
| Azure Kubernetes Service            | aks                         | azurerm_kubernetes_cluster              |
| Azure Kubernetes Service DNS prefix | aksdns                      | aks_dns_prefix                          |
| AKS Node Pool Linux                 | aksnpl                      | aks_node_pool_linux                     |
| AKS Node Pool Windows               | aksnpw                      | aks_node_pool_windows                   |
| Azure Site Recovery                 | asr                         | azurerm_recovery_services_vault         |
| Azure Availability Set              | avail                       | azurerm_availability_set                |
| Azure Vpn Connection                | cn                          | azurerm_vpn_connection                  |
| Azure Event Hubs                    | evh                         | azurerm_eventhub_namespace              |
| generic                             | gen                         | generic                                 |
| Azure Key Vault                     | kv                          | azurerm_key_vault                       |
| Azure Monitor Log Analytics         | la                          | azurerm_log_analytics_workspace         |
| Azure Load Balancer (External)      | lbe                         | azurerm_load_balancer_external          |
| Azure Load Balancer (Internal)      | lbi                         | azurerm_load_balancer_internal          |
| Azure Local Network Gateway         | lgw                         | azurerm_local_network_gateway           |
| Azure Mysql Database                | mysql                       | azurerm_mysql_database                  |
| Virtual Network Interface Card      | nic                         | azurerm_network_interface               |
| Network Security Group              | nsg                         | azurerm_network_security_group          |
| Public IP                           | pip                         | azurerm_public_ip                       |
| App Service Plan                    | plan                        | azurerm_app_service_plan                |
| Resource group                      | rg                          | azurerm_resource_group                  |
| Azure Route Table                   | route                       | azurerm_route_table                     |
| Azure Service Bus                   | sb                          | azurerm_service_bus                     |
| Azure Service Bus Queue             | sbq                         | azurerm_service_bus_queue               |
| Azure Service Bus Topic             | sbt                         | azurerm_service_bus_topic               |
| Subnet                              | snet                        | azurerm_subnet                          |
| Azure SQL DB Server                 | sql                         | azurerm_sql_server                      |
| Azure SQL DB                        | sqldb                       | azurerm_sql_database                    |
| Azure Storage Account               | st                          | azurerm_storage_account                 |
| Azure Traffic Manager Profile       | traf                        | azurerm_traffic_manager_profile         |
| Virtual Network Gateway             | vgw                         | azurerm_virtual_network_gateway         |
| Linux Virtual Machine               | vml                         | azurerm_virtual_machine_linux           |
| Virtual Machine Scale Set Linux     | vmssl                       | azurerm_vm_scale_set_linux              |
| Virtual Machine Scale Set Windows   | vmssw                       | azurerm_vm_scale_set_windows            |
| Windows Virtual Machine             | vmw                         | azurerm_virtual_machine_windows         |
| Virtual Network                     | vnet                        | azurerm_virtual_network                 |

## Parameters

### name
input name from the user (from landing zones settings or from blueprints)
name will be sanitized as per supported character set in Azure.

Example:

```hcl
name = "samplename"
```

### postfix
input postfix from the user (to managed instance numbers for instance)

Example:
```hcl
postfix = "001"
```

### convention
one of the four methods as described above:

Example:

```hcl
convention = "cafclassic"
```

### Resource type
describes the type of object you are requesting a name from, for instance if you are requesting a name for event hub:

```hcl
resource_type = "evh"
```

### Maximum length
configure the maximum length of the returned object name, is the specified length is longer than the supported length of the Azure resource the later applies

```hcl
max_length = 24
```

## Limitations and planned improvements

- Currently you can only get one name at a time, support for multiple names via input map of the same type coming.
- Filter for minimum size for passthrough method
- Support for hub_spoke landing zone components
- Support for VM components
Feel free to submit your PR to add capabilities

## Referencing the provider
To reference the provider simply include
```hcl
provider azurecaf {}
```
make sure that the terraform-provider-azurecaf binary is installed as a Terraform plugin

## Outputs

This provider outputs one name, the result of the naming convention query, you must specify the leverage the result output value, example for a storage account, you will get azurecaf_naming_convention.cafrandom_rg.result which returns the name based on the convention input.
This output will be consumed directly by a resource to name the component before calling the azurerm resource provider.

Example:
```hcl
resource "azurecaf_naming_convention" "cafrandom_rg" {  
  name    = "aztfmod"
  prefix  = "dev"
  resource_type    = "rg"
  postfix = "001"
  max_length = 23
  convention  = "cafrandom"
}

resource "azurerm_storage_account" "log" {
  name                      = azurecaf_naming_convention.cafrandom_rg.result
  resource_group_name       = var.resource_group_name
  location                  = var.location
  account_kind              = "StorageV2"
  account_tier              = "Standard"
  account_replication_type  = "GRS"
  access_tier               = "Hot"
  enable_https_traffic_only = true
}
```

