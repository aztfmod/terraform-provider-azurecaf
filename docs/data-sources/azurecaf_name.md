# azurecaf_name (Data Source)

The `azurecaf_name` data source generates Azure-compliant resource names following the Cloud Adoption Framework guidelines. **This is the recommended approach** for name generation as data sources are evaluated at plan time, making the generated names visible before resource creation.

## Key Features

- **Plan-time evaluation** - Names are visible during `terraform plan`
- **Azure compliance** - Automatically validates against Azure naming requirements
- **CAF guidelines** - Follows Microsoft Cloud Adoption Framework recommendations
- **Flexible configuration** - Supports prefixes, suffixes, random generation, and custom separators
- **Input sanitization** - Cleans inputs to ensure compliance
- **200+ resource types** - Comprehensive coverage of Azure services

## Example Usage

### Basic Resource Naming

```hcl
data "azurecaf_name" "example" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  prefixes      = ["prod"]
  suffixes      = ["001"]
  random_length = 3
  clean_input   = true
}

resource "azurerm_storage_account" "example" {
  name                     = data.azurecaf_name.example.result
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# Output: "stprodmyapp001abc" (st=storage slug, prod=prefix, myapp=name, 001=suffix, abc=random)
```

### Combined with Azure Resources

```hcl
data "azurecaf_name" "rg_example" {
  name          = "demogroup"
  resource_type = "azurerm_resource_group"
  random_length = 5
  clean_input   = true
}

resource "azurerm_resource_group" "rg" {
  name     = data.azurecaf_name.rg_example.result
  location = "southeastasia"
}
```

**Terraform Plan Output:**
```
data.azurecaf_name.rg_example: Reading...
data.azurecaf_name.rg_example: Read complete after 0s [id=rg-demogroup-wjyhr]

Terraform will perform the following actions:

  # azurerm_resource_group.rg will be created
  + resource "azurerm_resource_group" "rg" {
      + id       = (known after apply)
      + location = "southeastasia"
      + name     = "rg-demogroup-wjyhr"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Environment-Based Naming

```hcl
data "azurecaf_name" "app_service" {
  name          = "webapp"
  resource_type = "azurerm_app_service"
  prefixes      = [var.environment, var.project]
  suffixes      = ["web", "001"]
  random_length = 4
  separator     = "-"
}

# Example output: "app-prod-myproject-webapp-web-001-a1b2"
```

### Passthrough Mode (Validation Only)

```hcl
data "azurecaf_name" "existing_name" {
  name          = "mystorageaccount123"
  resource_type = "azurerm_storage_account"
  passthrough   = true
}

# Validates that "mystorageaccount123" complies with storage account naming rules
```

### Custom Separator and No Slug

```hcl
data "azurecaf_name" "custom" {
  name          = "database-server"
  resource_type = "azurerm_linux_virtual_machine"
  prefixes      = ["corp", "prod"]
  suffixes      = ["db", "001"]
  separator     = "_"
  use_slug      = false
  random_length = 4
}

# Output: "corp_prod_database_server_db_001_a1b2"
```


## Argument Reference

The following arguments are supported:

### Required Arguments

* `resource_type` - (Required) The Azure resource type for name generation (e.g., `azurerm_storage_account`, `azurerm_resource_group`). See [supported resource types](../index.md#supported-azure-resource-types).

### Optional Arguments

* `name` - (Optional) The base name for the resource. Will be sanitized according to the resource type's allowed character set. Defaults to empty string.

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

* `id` - Unique identifier for the naming configuration (same as `result`)
* `result` - The generated Azure-compliant resource name

## Naming Pattern

The generated name follows this pattern (when using default settings):

```
[prefix]-[prefix]...-[resource-slug]-[name]-[suffix]-[suffix]...-[random]
```

**Examples:**
- With slug: `rg-prod-myapp-001-abc12` 
- Without slug: `prod-myapp-001-abc12`
- Passthrough: `validated-input-name`

## Validation Rules

The data source automatically validates:

- **Length constraints** - Ensures generated names meet Azure length requirements
- **Character restrictions** - Filters invalid characters per resource type
- **Pattern compliance** - Validates against Azure naming patterns
- **Case requirements** - Handles case-sensitive vs case-insensitive resources

## Notes

- **Plan Visibility**: Names are generated during `terraform plan`, making them visible before resource creation
- **Deterministic**: Given the same inputs, the data source produces the same output (except when using time-based random seed)
- **Resource Compliance**: All generated names are guaranteed to comply with Azure naming requirements
- **Migration**: This data source is the recommended replacement for the `azurecaf_naming_convention` resource
