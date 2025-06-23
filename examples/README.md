# Azure CAF Terraform Provider Examples

This directory contains example configurations for the Azure Cloud Adoption Framework (CAF) Terraform Provider. These examples demonstrate different usage patterns and features to help you get started quickly.

## 📋 Quick Start

To run these examples locally:

1. **Build the provider:**
   ```bash
   make build
   ```

2. **Run the examples:**
   ```bash
   make test
   ```

The `make test` command will:
- Build the provider locally
- Set up development overrides for Terraform
- Initialize and apply all examples in this directory

## 📁 Example Files

### Core Examples

| File | Description | Key Features |
|------|-------------|--------------|
| `resource_name.tf` | Basic `azurecaf_name` resource usage | Resource-based naming, multiple resource types |
| `name_resource_and_datasource.tf` | Comparison of resource vs data source | Shows differences and use cases |

### Feature Demonstrations

Each example file demonstrates specific features:

- **Basic Naming**: Simple name generation with prefixes/suffixes
- **Random Generation**: Adding random characters for uniqueness
- **Input Cleaning**: Sanitizing inputs to meet Azure requirements
- **Validation**: Using passthrough mode to validate existing names
- **Multiple Resources**: Generating names for multiple resource types

## 🎯 Key Examples Explained

### Resource vs Data Source Usage

**Data Source (Recommended):**
```hcl
# Evaluated during plan phase - shows name before creation
data "azurecaf_name" "storage" {
  name          = "mydata"
  resource_type = "azurerm_storage_account"
  prefixes      = ["prod"]
  random_length = 3
}

resource "azurerm_storage_account" "example" {
  name = data.azurecaf_name.storage.result
  # ... other configuration
}
```

**Resource:**
```hcl
# Useful for generating multiple related names
resource "azurecaf_name" "multi_res" {
  name           = "myapp"
  resource_type  = "azurerm_app_service"
  resource_types = [
    "azurerm_app_service_plan",
    "azurerm_application_insights"
  ]
  prefixes      = ["prod"]
  random_length = 3
}

# Access names:
# Primary: azurecaf_name.multi_res.result
# All: azurecaf_name.multi_res.results
```

### Advanced Configuration Features

**Custom Separators and No Slug:**
```hcl
data "azurecaf_name" "custom" {
  name          = "database"
  resource_type = "azurerm_postgresql_server"
  prefixes      = ["corp", "prod"]
  suffixes      = ["primary"]
  separator     = "_"
  use_slug      = false
  random_length = 4
}
# Result: "corp_prod_database_primary_a1b2"
```

**Input Cleaning:**
```hcl
data "azurecaf_name" "cleaned" {
  name          = "my-app@company.com"
  resource_type = "azurerm_storage_account"
  clean_input   = true  # Removes invalid characters
}
# Input: "my-app@company.com" → Output: "stmyappcompanycom"
```

**Passthrough Validation:**
```hcl
data "azurecaf_name" "validate" {
  name          = "existingstorageaccount123"
  resource_type = "azurerm_storage_account"
  passthrough   = true  # Validates without modification
}
```

## 🔧 Advanced Usage Patterns

### Environment-Based Naming

```hcl
# variables.tf
variable "environment" {
  description = "Environment name"
  type        = string
  validation {
    condition = contains(["dev", "test", "prod"], var.environment)
    error_message = "Environment must be dev, test, or prod."
  }
}

variable "application_name" {
  description = "Application name"
  type        = string
}

# naming.tf
locals {
  environment_config = {
    dev = {
      prefix = "dev"
      random_length = 3
    }
    test = {
      prefix = "tst"
      random_length = 3
    }
    prod = {
      prefix = "prd"
      random_length = 5
    }
  }
  
  current_config = local.environment_config[var.environment]
}

data "azurecaf_name" "resources" {
  for_each = toset([
    "azurerm_resource_group",
    "azurerm_storage_account",
    "azurerm_key_vault",
    "azurerm_app_service"
  ])
  
  name          = var.application_name
  resource_type = each.key
  prefixes      = [local.current_config.prefix]
  random_length = local.current_config.random_length
}

# Output all generated names
output "resource_names" {
  value = {
    for resource_type, name_data in data.azurecaf_name.resources :
    resource_type => name_data.result
  }
}
```

### Consistent Multi-Tier Application Naming

```hcl
# Multi-tier application with consistent naming
locals {
  app_name = "ecommerce"
  env      = "prod"
  instance = "001"
}

# Frontend tier
data "azurecaf_name" "frontend" {
  for_each = toset([
    "azurerm_app_service",
    "azurerm_app_service_plan"
  ])
  
  name          = local.app_name
  resource_type = each.key
  prefixes      = [local.env, "web"]
  suffixes      = [local.instance]
}

# Backend tier  
data "azurecaf_name" "backend" {
  for_each = toset([
    "azurerm_app_service",
    "azurerm_app_service_plan"
  ])
  
  name          = local.app_name
  resource_type = each.key
  prefixes      = [local.env, "api"]
  suffixes      = [local.instance]
}

# Data tier
data "azurecaf_name" "data" {
  for_each = toset([
    "azurerm_postgresql_server",
    "azurerm_storage_account",
    "azurerm_key_vault"
  ])
  
  name          = local.app_name
  resource_type = each.key
  prefixes      = [local.env, "data"]
  suffixes      = [local.instance]
}
```

### Module Integration Example

```hcl
# modules/app-service/main.tf
variable "application_name" {
  description = "Name of the application"
  type        = string
}

variable "environment" {
  description = "Environment (dev/test/prod)"
  type        = string
}

variable "instance" {
  description = "Instance number"
  type        = string
  default     = "001"
}

# Generate names for all app service resources
data "azurecaf_name" "app_service_resources" {
  for_each = toset([
    "azurerm_resource_group",
    "azurerm_app_service_plan", 
    "azurerm_app_service",
    "azurerm_application_insights"
  ])
  
  name          = var.application_name
  resource_type = each.key
  prefixes      = [var.environment]
  suffixes      = [var.instance]
}

# Create resources with generated names
resource "azurerm_resource_group" "main" {
  name     = data.azurecaf_name.app_service_resources["azurerm_resource_group"].result
  location = "East US"
}

resource "azurerm_app_service_plan" "main" {
  name                = data.azurecaf_name.app_service_resources["azurerm_app_service_plan"].result
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  
  sku {
    tier = "Standard"
    size = "S1"
  }
}

# Output generated names for reference
output "resource_names" {
  description = "Generated resource names"
  value = {
    for key, name_data in data.azurecaf_name.app_service_resources :
    key => name_data.result
  }
}
```

## 🧪 Testing Your Examples

### Local Testing

1. **Validate configuration:**
   ```bash
   terraform validate
   ```

2. **Plan to see generated names:**
   ```bash
   terraform plan
   ```

3. **Apply (if testing with real resources):**
   ```bash
   terraform apply
   ```

### Unit Testing Examples

Create tests for your naming patterns:

```hcl
# test/naming_test.go
func TestNamingConventions(t *testing.T) {
    tests := []struct {
        name         string
        config       map[string]interface{}
        expectedMatch string
    }{
        {
            name: "production storage account",
            config: map[string]interface{}{
                "name":          "myapp",
                "resource_type": "azurerm_storage_account",
                "prefixes":      []string{"prod"},
                "random_length": 3,
            },
            expectedMatch: `^stprodmyapp[a-z0-9]{3}$`,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 📖 Best Practices Demonstrated

### 1. **Use Data Sources for Single Names**
- Evaluated during plan phase
- Names visible before resource creation
- Better for most use cases

### 2. **Use Resources for Multiple Names** 
- When generating names for related resources
- Leverages `resource_types` parameter
- Maintains consistency across resource types

### 3. **Environment-Specific Configuration**
- Use locals or variables for environment-specific settings
- Consistent prefixes and random lengths per environment
- Validation rules for environment values

### 4. **Input Validation and Cleaning**
- Always enable `clean_input` unless you need strict control
- Use validation for environment and application names
- Test with various input patterns

### 5. **Consistent Naming Patterns**
- Establish patterns for prefixes, suffixes, and separators
- Document naming conventions in your organization
- Use modules to enforce consistency

## 🔗 Additional Resources

- **Provider Documentation**: [Terraform Registry](https://registry.terraform.io/providers/aztfmod/azurecaf/latest)
- **Azure Naming Rules**: [Microsoft Documentation](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules)
- **CAF Guidelines**: [Cloud Adoption Framework](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging)
- **Testing Guide**: [TESTING.md](../TESTING.md)

---

💡 **Need more examples?** Check the [main README](../README.md) for additional usage patterns and troubleshooting tips.