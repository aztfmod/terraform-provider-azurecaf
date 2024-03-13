# Contributing to the CAF provider

👍🎉 First off, thanks for taking the time to contribute! 🎉👍

## Code of conduct

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## What should I know before I get started?

The current goal of the CAFprovider is to support the CAF landing zones but can also be used to standardize the naming convention of your projects. It is important to keep in mind that many of the design decisions on the provider have been made to accommodate the needs of CAF.

To contribute to this project you are required to have at least go 1.13 installed in your system

## Adding a new resource

Please, find below the steps that should be followed to contribute:

1. Check if the resource has been implemented already

    You can find a list of resources implemented in the [README.md#resource-status](../README.md) under the resource status section.

2. Create an issue for the missing resource

    If there is no [issue created already](https://github.com/aztfmod/terraform-provider-azurecaf/issues) for the implementation of this resource you should [create an issue](https://docs.github.com/en/free-pro-team@latest/github/managing-your-work-on-github/creating-an-issue) requesting the implementation of the resource.

3. Check the requirements for your resource Name

    You can check the requirements for your resource name in the [docs](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules) or by checking the error message returned when trying to create the resource on Azure with an invalid name. Slug value can also be checked in the CAF [docs](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations).

4. Choose the slug for the resource

    Every resource in CAF does have a slug that associate with this resource this is 2 to 5 letters that identify that resource, for example, the slug for a `key vault` is `kv` for a storage account `st` What is important here it is to try to keep this short but meaningful and avoid collision with existing ones. Don't worry about knowing all existing ones if you choose one that exists already the tests will fail. You can also check if the resource has a example abbreviation on this page: [doc](https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations)

5. Modify the `resourceDefinition.json`

    You should now add your resource to the resource definitions in the just add another resource on the list. You can use the existing resources as a template for your resource implementation

6. Generate the definitions based on the `resourcedefinition.json` and test

    You can run `make build` in case you have make installed in your system in case you don't you can run from the repository root `go generate` followed by `go fmt ./...` and then `go test ./...`

7. Update the README.MD with coverage
 
    For quick reference, update the [README.md#resource-status](../README.md) at the root of the provider to mention the coverage you just added:
    ```|azurerm_api_management_custom_domain | ✔ |```

7. Commit and submit PR

    Now you should commit, remembering to put a meaningful commit message. After that, you should [make pull request](https://docs.github.com/en/free-pro-team@latest/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request) remembering to link in the PR the issue that it is solving.

### The `resourceDefinition.json`

Once you have all the information and have created an issue if one doesn't exist yet you can start to fill up the resource in the `resourceDefinition.json`

Each resource in the `resourceDefinitions.json` follow the following schema:

```json
{
    "name": "azurerm_snapshots", //Azurerm name of the resource
    "min_length": 1, // Minumum number of chars that this resource requires
    "max_length": 80, // Maximum number of chars that this resource can have
    "validation_regex": "\"^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$\"", // A regex expression that will match only a valid resource name
    "scope": "parent", // Where this name must be unique. global means that only one resource with this name it is allowed in azure. parent means that only one resource of this name based in the parent resource. Resource group means only one resource with this name per resource group.
    "slug": "snap", // This are the letters that identify the resource type
    "dashes": true, // if this resource allows you to use dashes '-'
    "lowercase": false, // if this resource will ONLY allow lowercase
    "regex": "\"[^0-9A-Za-z_.-]\"" // This is the 'cleaning' regex anything that is matched by this regex will be removed from the resource name that is why you normally use the negation of all the allowed chars in this regex.
}
```

## Legal

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.
