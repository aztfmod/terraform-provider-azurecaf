# azurecaf_environment_variable

The `azurecaf_environment_variable` data source provides a secure way to read environment variables from the system where Terraform is running. This is useful for integrating with CI/CD systems and external configuration management.

## Example Usage

### Basic Usage

```hcl
data "azurecaf_environment_variable" "subscription_id" {
  name = "ARM_SUBSCRIPTION_ID"
}

data "azurecaf_environment_variable" "environment" {
  name           = "ENVIRONMENT"
  fails_if_empty = true
}

output "subscription_id" {
  value = data.azurecaf_environment_variable.subscription_id.value
}
```

### With Default Values

```hcl
data "azurecaf_environment_variable" "log_level" {
  name = "LOG_LEVEL"
}

locals {
  log_level = data.azurecaf_environment_variable.log_level.value != "" ? 
              data.azurecaf_environment_variable.log_level.value : "INFO"
}
```

### Integration with Naming Convention

```hcl
data "azurecaf_environment_variable" "environment" {
  name           = "TERRAFORM_ENVIRONMENT"
  fails_if_empty = true
}

data "azurecaf_name" "storage_account" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  prefixes      = [data.azurecaf_environment_variable.environment.value]
  random_length = 3
}

resource "azurerm_storage_account" "example" {
  name                = data.azurecaf_name.storage_account.result
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the environment variable to read.

* `fails_if_empty` - (Optional) If set to `true`, Terraform will fail if the environment variable is not set or is empty. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `value` - The value of the environment variable. If the environment variable is not set, this will be an empty string.

## Security Considerations

**Important**: Environment variables retrieved through this data source will be stored in Terraform state. Be cautious when using this for sensitive values such as:

- API keys or tokens
- Database passwords
- Private keys or certificates
- Other secrets

For sensitive values, consider using:
- Terraform Cloud/Enterprise workspace variables marked as sensitive
- Azure Key Vault with the AzureRM provider's key vault data sources
- External secret management systems

## Common Use Cases

### CI/CD Integration

```hcl
# Read CI/CD provided environment variables
data "azurecaf_environment_variable" "build_number" {
  name = "BUILD_NUMBER"
}

data "azurecaf_environment_variable" "git_branch" {
  name = "GIT_BRANCH"
}

# Use in resource naming
data "azurecaf_name" "app_service" {
  name          = "myapp"
  resource_type = "azurerm_app_service"
  prefixes      = [data.azurecaf_environment_variable.git_branch.value]
  suffixes      = [data.azurecaf_environment_variable.build_number.value]
}
```

### Multi-Environment Configuration

```hcl
data "azurecaf_environment_variable" "environment" {
  name           = "TF_ENVIRONMENT"
  fails_if_empty = true
}

locals {
  environment_config = {
    dev = {
      location = "East US"
      sku      = "Basic"
    }
    prod = {
      location = "West US"
      sku      = "Standard"
    }
  }
  
  current_config = local.environment_config[data.azurecaf_environment_variable.environment.value]
}
```

### Feature Flags

```hcl
data "azurecaf_environment_variable" "enable_monitoring" {
  name = "ENABLE_MONITORING"
}

resource "azurerm_application_insights" "example" {
  count = data.azurecaf_environment_variable.enable_monitoring.value == "true" ? 1 : 0
  
  name                = "ai-${var.application_name}"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}
```