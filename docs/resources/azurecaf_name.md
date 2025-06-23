# azurecaf_name (Resource)

The `azurecaf_name` resource generates Azure-compliant resource names following the Cloud Adoption Framework guidelines. This resource provides more flexibility and comprehensive resource type support compared to the legacy `azurecaf_naming_convention` resource.

> **Note**: For most use cases, the [`azurecaf_name` data source](../data-sources/azurecaf_name.md) is recommended as it evaluates names at plan time, making them visible before resource creation.

## Key Features

- **200+ Resource Types** - Comprehensive coverage of Azure services with accurate validation
- **CAF Compliance** - Follows Microsoft Cloud Adoption Framework recommendations
- **Multi-Resource Support** - Generate names for multiple related resource types simultaneously
- **Flexible Configuration** - Supports prefixes, suffixes, random generation, and custom patterns
- **Input Sanitization** - Automatically cleans inputs to ensure Azure compliance
- **Validation Rules** - Enforces length, character, and pattern requirements per resource type

## Example Usage

### Basic Resource Naming

```hcl
resource "azurecaf_name" "example" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  prefixes      = ["prod"]
  suffixes      = ["001"]
  random_length = 3
  clean_input   = true
}

resource "azurerm_storage_account" "example" {
  name                     = azurecaf_name.example.result
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# Output: "stprodmyapp001abc"
```

### Multi-Resource Naming

Generate names for multiple related resource types with consistent settings:

```hcl
resource "azurecaf_name" "webapp_resources" {
  name           = "webapp"
  resource_type  = "azurerm_app_service"
  resource_types = [
    "azurerm_app_service_plan",
    "azurerm_application_insights"
  ]
  prefixes       = ["prod"]
  suffixes       = ["web"]
  random_length  = 3
  clean_input    = true
}

# Access names:
# Primary: azurecaf_name.webapp_resources.result
# Additional: azurecaf_name.webapp_resources.results["azurerm_app_service_plan"]
# Additional: azurecaf_name.webapp_resources.results["azurerm_application_insights"]

resource "azurerm_app_service_plan" "example" {
  name                = azurecaf_name.webapp_resources.results["azurerm_app_service_plan"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  
  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = azurecaf_name.webapp_resources.result
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id
}
```

### Complex Naming Pattern

```hcl
resource "azurecaf_name" "rg_example" {
  name          = "demogroup"
  resource_type = "azurerm_resource_group"
  prefixes      = ["corp", "proj"]
  suffixes      = ["web", "001"]
  random_length = 5
  separator     = "_"
  clean_input   = true
}

resource "azurerm_resource_group" "demo" {
  name     = azurecaf_name.rg_example.result
  location = "southeastasia"
}

# Output: "corp_proj_rg_demogroup_web_001_abc12"
```

### Passthrough Mode (Validation)

```hcl
resource "azurecaf_name" "existing_name" {
  name          = "mystorageaccount123"
  resource_type = "azurerm_storage_account"
  passthrough   = true
}

# Validates and returns "mystorageaccount123" if compliant
```

### Custom Patterns

```hcl
resource "azurecaf_name" "custom_vm" {
  name          = "database-server"
  resource_type = "azurerm_linux_virtual_machine"
  prefixes      = ["corp", "prod"]
  suffixes      = ["db", "001"]
  separator     = "_"
  use_slug      = false  # No "vm" prefix
  random_length = 4
  clean_input   = true
}

# Output: "corp_prod_database_server_db_001_a1b2"
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `resource_type` - (Required) The Azure resource type for name generation (e.g., `azurerm_storage_account`, `azurerm_resource_group`). See [supported resource types](../index.md#supported-azure-resource-types).

### Optional Arguments

* `name` - (Optional) The base name for the resource. Will be sanitized according to the resource type's allowed character set. Defaults to empty string.

* `resource_types` - (Optional) List of additional resource types for generating multiple names with the same configuration. Used with the `results` attribute.

* `prefixes` - (Optional) List of prefixes to prepend to the generated name. Prefixes are separated by the separator character. Defaults to `[]`.

* `suffixes` - (Optional) List of suffixes to append to the generated name. Suffixes are separated by the separator character. Defaults to `[]`.

* `random_length` - (Optional) Number of random characters to append. Random characters comply with the resource's allowed character set. Defaults to `0`.

* `random_seed` - (Optional) Seed for random character generation. Use `0` for time-based seed (default behavior). Defaults to `0`.

* `separator` - (Optional) Character used to separate name components (prefixes, resource type slug, name, suffixes). Defaults to `"-"`.

* `clean_input` - (Optional) Remove non-compliant characters from name, prefixes, and suffixes. **Recommended to keep enabled.** Defaults to `true`.

* `passthrough` - (Optional) Enable passthrough mode for name validation only. When enabled, only input cleaning is applied; prefixes, suffixes, random characters, and resource slug are ignored. Defaults to `false`.

* `use_slug` - (Optional) Include resource type abbreviation (slug) in the generated name. When `false`, no resource type identifier is added. Defaults to `true`.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier for the naming configuration
* `result` - The generated Azure-compliant name for the primary resource type
* `results` - Map of generated names for all resource types specified in `resource_types` (includes the primary `resource_type`)

## Naming Pattern

The generated name follows this pattern (when using default settings):

```
[prefix]-[prefix]...-[resource-slug]-[name]-[suffix]-[suffix]...-[random]
```

**Examples:**
- Standard: `rg-prod-myapp-001-abc12`
- Without slug: `prod-myapp-001-abc12` 
- Passthrough: `validated-input-name`
- Multi-separator: `corp_prod_rg_myapp_db_001_abc12`

## Multi-Resource Usage

When using `resource_types`, the resource generates names for multiple resource types:

```hcl
resource "azurecaf_name" "multi" {
  name           = "webapp"
  resource_type  = "azurerm_app_service"        # Primary type
  resource_types = [
    "azurerm_app_service_plan",                 # Additional types
    "azurerm_application_insights"
  ]
  prefixes       = ["prod"]
  random_length  = 3
}

# Access the names:
output "app_service_name" {
  value = azurecaf_name.multi.result  # Primary resource type
}

output "all_names" {
  value = azurecaf_name.multi.results  # Map of all names
}
```

## Migration from azurecaf_naming_convention

```hcl
# Legacy approach (deprecated)
resource "azurecaf_naming_convention" "old" {
  name         = "myapp"
  resource_type = "rg"
  convention   = "cafrandom"
  prefix       = "prod"
  postfix      = "001"
}

# New approach (recommended)
resource "azurecaf_name" "new" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  prefixes      = ["prod"]
  suffixes      = ["001"]
  random_length = 5
}
```

## Supported Resource Types

This resource supports **200+ Azure resource types** with accurate naming validation rules. 

For the complete list of supported resource types, validation rules, and examples, see the [main provider documentation](../index.md#supported-azure-resource-types).

## Notes

### Data Source vs Resource

**Recommendation**: Use the [`azurecaf_name` data source](../data-sources/azurecaf_name.md) instead of this resource when possible, as data sources:
- Evaluate at plan time, showing generated names before resource creation
- Provide better visibility in Terraform plans
- Are generally preferred for name generation workflows

### State Management

Resource names are stored in Terraform state. Changes to naming parameters will trigger resource recreation, which may affect dependent resources.

### Validation

All generated names are automatically validated against:
- Azure naming requirements per resource type
- Length constraints (minimum and maximum)
- Character restrictions and allowed patterns
- Case sensitivity requirements

### Performance

The resource supports generating multiple names simultaneously using the `resource_types` argument, which is more efficient than creating multiple separate `azurecaf_name` resources.

## Related Resources

- [`azurecaf_name` data source](../data-sources/azurecaf_name.md) - Recommended approach for name generation
- [`azurecaf_environment_variable` data source](../data-sources/azurecaf_environment_variable.md) - Read environment variables for dynamic naming

For a complete list of supported resource types with their constraints and validation rules, see the [Provider Index](../index.md#supported-azure-resource-types) documentation.
