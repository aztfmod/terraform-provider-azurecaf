# azurecaf_environment_variable

The data source azurecaf_environment_variable retrieve an OS environment variable.

## Exemple usage
This example shows how to get the value of an environment variable.

```hcl
# Retrieve the PATH variable
data "azurecaf_environment_variable" "path" {
  name = "PATH"
}

# Retreive the PAT_TOKEN variable as a sensitive data and through an error if it does not exist.
data "azurecaf_environment_variable" "PAT" {
    name           = "PAT_TOKEN"
    fails_if_empty = true
    sensitive      = true
}
```

## Argument Reference

The following arguments are supported:

* name - (required) Name of the environment variable.
* fails_if_empty (optional) - Through an error if the environment variable is not set (default: false).
* sensitive (optional) - Do not display the value in the log is the value is sensitive (default: false).

# Attributes Reference
The following attributes are exported:

* id - The id of the environment variable
* value - Value of the environment variable.
* value_sensitive - Value (sensitive) of the environment variable.

