# caf naming convention Provider

The caf naming convention provider implements a set of methodologies for to apply consistent resource naming using the default Microsoft Cloud Adoption Framework for Azure recommendations as per https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging.

## Exemple usage
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
  name     = azurecaf_naming_convention.cafrandom_rg.
  location = "southeastasia"
}
```

```

The provider generates a name using the input parameters and automatically appends a prefix (if defined), a caf prefix (resource type) and postfix (if defined) in addition to a generated padding string based on the selected naming convention.


## Argument Reference
The following arguments are supported:

* name - (optional) the basename of the resource to create, the basename will be sanitized as per supported character set in Azure.
* convention (optional): one of the four naming convention supported. Defaults to cafrandom. Allowed values are cafclassic, cafrandom, random, passthrough
* prefix (optional)
* postfix (optional)
* max_length (optional)
* resource_type (optional)

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
