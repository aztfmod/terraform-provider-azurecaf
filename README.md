[![VScodespaces](https://img.shields.io/endpoint?url=https%3A%2F%2Faka.ms%2Fvso-badge)](https://online.visualstudio.com/environments/new?name=caf%20terraform%20provider&repo=aztfmod/terraform-provider-azurecaf)

# Azure Cloud Adoption Framework - Terraform provider

This provider implements a set of methodologies for naming convention implementation including the default Microsoft Cloud Adoption Framework for Azure recommendations as per https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging.

## Using the Provider

You can simply consume the provider from the Terraform registry from the following URL: [https://registry.terraform.io/providers/aztfmod/azurecaf/latest](https://registry.terraform.io/providers/aztfmod/azurecaf/latest), then add it in your provider declaration as follow: 

```hcl
terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
      version = "1.1.3"
    }
  }
}
```

The azurecaf_name resource allows you to:

* Clean inputs to make sure they remain compliant with the allowed patterns for each Azure resource
* Generate random characters to append at the end of the resource name
* Handle prefix, suffixes (either manual or as per the Azure cloud adoption framework resource conventions)
* Allow passthrough mode (simply validate the output)

## Example usage

This example outputs one name, the result of the naming convention query. The result attribute returns the name based on the convention and parameters input.

The example generates a 23 characters name compatible with the specification for an Azure Resource Group
dev-aztfmod-001

```hcl
resource "azurecaf_name" "rg_example" {
  name            = "demogroup"
    resource_type   = "azurerm_resource_group"
    prefixes        = ["a", "b"]
    suffixes        = ["y", "z"]
    random_length   = 5
    clean_input     = true
}

resource "azurerm_resource_group" "demo" {
  name     = azurecaf_name.rg_example.result
  location = "southeastasia"
}
```

The provider generates a name using the input parameters and automatically appends a prefix (if defined), a caf prefix (resource type) and postfix (if defined) in addition to a generated padding string based on the selected naming convention.

The example above would generate a name using the pattern [prefix]-[cafprefix]-[name]-[postfix]-[5_random_chars]:

## Argument Reference

The following arguments are supported:

* **name** - (optional) the basename of the resource to create, the basename will be sanitized as per supported characters set for each Azure resources.
* **prefixes** (optional) - a list of prefix to append as the first characters of the generated name - prefixes will be separated by the separator character
* **suffixes** (optional) -  a list of additional suffix added after the basename, this is can be used to append resource index (eg. vm-001). Suffixes are separated by the separator character
* **random_length** (optional) - default to ``0`` : configure additional characters to append to the generated resource name. Random characters will remain compliant with the set of allowed characters per resources and will be appended after the suffixes
* **random_seed** (optional) - default to ``0`` : Define the seed to be used for random generator. 0 will not be respected and will generate a seed based in the unix time of the generation.
* **resource_type** (optional) -  describes the type of azure resource you are requesting a name from (eg. azure container registry: azurerm_container_registry). See the Resource Type section
* **resource_types** (optional) -  a list of additional resource type should you want to use the same settings for a set of resources
* **separator** (optional) - defaults to ``-``. The separator character to use between prefixes, resource type, name, suffixes, random character
* **clean_input** (optional) - defaults to ``true``. remove any noncompliant character from the name, suffix or prefix.
* **passthrough** (optional) - defaults to ``false``. Enables the passthrough mode - in that case only the clean input option is considered and the prefixes, suffixes, random, and are ignored. The resource prefixe is not added either to the resulting string
* **use_slug** (optional) - defaults to ``true``. If a slug should be added to the name - If you put false no slug (the few letters that identify the resource type) will be added to the name.

## Attributes Reference

The following attributes are exported:

* **id** - The id of the naming convention object
* **result** - The generated named for an Azure Resource based on the input parameter and the selected naming convention
* **results** - The generated name for the Azure resources based in the resource_types list

## Resource types

We define resource types as per [naming-and-tagging](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging)
The comprehensive list of resource type can be found [here](./docs/resources/azurecaf_name.md)


## Building the provider

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

## Developing the provider

If you wish to work on the provider, you'll first need Go installed on your machine (version 1.13+ is required). You'll also need to correctly setup a GOPATH, as well as adding $GOPATH/bin to your $PATH.

To compile the provider, run make build. This will build the provider and put the provider binary in the $GOPATH/bin directory.

```
$ make build
...
$ $GOPATH/bin/terraform-provider-azurecaf
...

```
## Testing

Running the acceptance test suite requires does not require an Azure subscription. 

to run the unit test:
```
make unittest
```

to run the integration test

```
make test
```

## Related repositories

| Repo                                                                                             | Description                                                |
|--------------------------------------------------------------------------------------------------|------------------------------------------------------------|
| [caf-terraform-landingzones](https://github.com/azure/caf-terraform-landingzones)                | landing zones repo with sample and core documentations     |
| [rover](https://github.com/aztfmod/rover)                                                        | devops toolset for operating landing zones                 |
| [azure_caf_provider](https://github.com/aztfmod/terraform-provider-azurecaf)                     | custom provider for naming conventions                     |
| [module](https://registry.terraform.io/modules/aztfmod)                                          | official CAF module available in the Terraform registry    |


## Community

Feel free to open an issue for feature or bug, or to submit a PR.

In case you have any question, you can reach out to tf-landingzones at microsoft dot com.

You can also reach us on [Gitter](https://gitter.im/aztfmod/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

## Code of conduct

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
