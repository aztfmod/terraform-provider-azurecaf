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
| Resource group                      | rg                          | azurerm_resource_group                  |
| Azure Storage Account               | st                          | azurerm_storage_account                 |
| Azure Event Hubs                    | evh                         | azurerm_eventhub_namespace              |
| Azure Monitor Log Analytics         | la                          | azurerm_log_analytics_workspace         |
| Azure Key Vault                     | kv                          | azurerm_key_vault                       |
| Windows Virtual Machine             | vmw                         | azurerm_windows_virtual_machine_windows |
| Linux Virtual Machine               | vml                         | azurerm_windows_virtual_machine_linux   |
| Public IP                           | pip                         | azurerm_public_ip                       |
| Network Security Group              | nsg                         | azurerm_network_security_group          |
| Virtual Network Interface Card      | nic                         | azurerm_network_interface               |
| Virtual Network                     | vnet                        | azurerm_virtual_network                 |
| Azure Firewall                      | afw                         | azurerm_firewall                        |
| Azure Container Registry            | acr                         | azurerm_container_registry              |
| Azure Site Recovery                 | asr                         | azurerm_recovery_services_vault         |
| Azure Automation                    | aaa                         | azurerm_automation_account              |
| generic                             | gen                         | generic                                 |

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

### type of object
describes the type of object you are requesting a name from, for instance if you are requesting a name for event hub:

```hcl
type = "evh"
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
