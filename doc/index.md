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
