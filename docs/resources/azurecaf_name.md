# azurecaf_name (Resource)

The `azurecaf_name` resource generates Azure-compliant resource names following the Cloud Adoption Framework guidelines. This resource provides more flexibility and comprehensive resource type support compared to the legacy `azurecaf_naming_convention` resource.

> **Note**: For most use cases, the [`azurecaf_name` data source](../data-sources/azurecaf_name.md) is recommended as it evaluates names at plan time, making them visible before resource creation.

## Key Features

- **300+ Resource Types** - Comprehensive coverage of Azure services with accurate validation
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

# Name Composition and Truncation

This section provides detailed information about how the Azure CAF provider composes resource names, handles length constraints, and applies truncation when necessary.

## Name Composition Order

The provider follows a specific order when composing resource names, controlled by the **name precedence** algorithm. The default precedence order is:

1. **`name`** - The base name parameter
2. **`slug`** - The resource type abbreviation (when `use_slug = true`)
3. **`random`** - Random characters (when `random_length > 0`)
4. **`suffixes`** - Suffix strings (applied in order)
5. **`prefixes`** - Prefix strings (applied in reverse order)

### Component Placement

- **Prefixes**: Added to the **beginning** of the name (in reverse order: last prefix first)
- **Slug**: Added to the **beginning** after prefixes
- **Name**: The core name component
- **Suffixes**: Added to the **end** (in order: first suffix first)
- **Random**: Added to the **end** after suffixes

### Example Composition

```hcl
resource "azurecaf_name" "example" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  prefixes      = ["corp", "prod"]
  suffixes      = ["web", "001"]
  random_length = 3
  separator     = "-"
}
```

**Composition process:**
1. Start with empty name: `""`
2. Add prefixes (reverse order): `"prod-corp"`
3. Add slug: `"st-prod-corp"`
4. Add name: `"st-prod-corp-myapp"`
5. Add suffixes (forward order): `"st-prod-corp-myapp-web-001"`
6. Add random: `"st-prod-corp-myapp-web-001-abc"`

**Final result:** `"stprodcorpmyappweb001abc"` (after separator processing and lowercase conversion)

## Length Constraints and Truncation

### Maximum Length Enforcement

Each Azure resource type has specific length constraints defined in the provider. When the composed name exceeds the maximum length, the provider applies intelligent truncation.

### Truncation Algorithm

The provider uses a **priority-based truncation** system that respects the name precedence order:

1. **Calculate space**: Determine available space within the maximum length
2. **Add components by precedence**: Add each component only if it fits within remaining space
3. **Skip if no space**: If a component doesn't fit, it's skipped entirely
4. **Final trim**: Apply final length trimming if necessary

### Truncation Priority

Components are added in this priority order (higher priority = added first):

1. **`name`** (highest priority)
2. **`slug`** 
3. **`random`**
4. **`suffixes`**
5. **`prefixes`** (lowest priority)

This means if space is limited:
- The core `name` is always preserved
- `prefixes` are the first to be dropped
- `suffixes` are dropped before `random` or `slug`

### Truncation Examples

#### Example 1: Prefix Truncation
```hcl
# Storage account max length: 24 characters
resource "azurecaf_name" "example" {
  name          = "verylongapplicationname"  # 23 chars
  resource_type = "azurerm_storage_account"
  prefixes      = ["corporate"]              # 9 chars + separator
  use_slug      = true                       # "st" = 2 chars
}
```

**Process:**
- Available space: 24 characters
- Core name: "verylongapplicationname" (23 chars) - **added**
- Slug: "st" (2 chars) - would exceed limit, **skipped**
- Prefix: "corporate" - would exceed limit, **skipped**

**Result:** `"verylongapplicationname"` (23 chars)

#### Example 2: Suffix Truncation
```hcl
resource "azurecaf_name" "example" {
  name          = "myapp"                    # 5 chars
  resource_type = "azurerm_storage_account"
  suffixes      = ["production", "web", "001"] # Multiple suffixes
  random_length = 8                          # 8 chars
  use_slug      = true                       # "st" = 2 chars
}
```

**Process:**
- Available space: 24 characters
- Name: "myapp" (5 chars) - **added** (total: 5)
- Slug: "st" (2 chars) - **added** (total: 7)
- Random: 8 chars - **added** (total: 15)
- Suffix "production" (10 chars) - **added** (total: 25) - exceeds limit, **skipped**
- Suffix "web" (3 chars) - **added** (total: 18)
- Suffix "001" (3 chars) - **added** (total: 21)

**Result:** `"stmyappweb001abcdefgh"` (21 chars)

## Component Processing Rules

### Separator Handling

- Separators are only added between components when both components are present
- No leading or trailing separators
- Separator length is included in total length calculations

### Case Conversion

Many Azure resource types require lowercase names:

```hcl
# Input with mixed case
resource "azurecaf_name" "example" {
  name          = "MyApp"
  resource_type = "azurerm_storage_account"  # Requires lowercase
}
# Result: "stmyapp" (converted to lowercase)
```

### Input Cleaning

When `clean_input = true`, the provider sanitizes inputs:

- Removes invalid characters for the specific resource type
- Applies character restrictions (e.g., alphanumeric only)
- Removes characters that don't match the resource's validation pattern

### Passthrough Mode

When `passthrough = true`:

- **Composition is bypassed** - only the `name` parameter is used
- Prefixes, suffixes, slug, and random components are ignored
- Length trimming and validation still apply
- Useful for using pre-composed names while still getting validation

```hcl
resource "azurecaf_name" "example" {
  name          = "mycustomstorageaccount"
  resource_type = "azurerm_storage_account"
  passthrough   = true
  # prefixes, suffixes, etc. are ignored
}
# Result: "mycustomstorageaccount"
```

## Validation and Error Handling

### Length Validation

- Names that exceed maximum length after truncation will cause errors
- Random length is validated against resource type constraints
- Minimum length requirements are enforced

### Pattern Validation

After composition and truncation, names must match the resource type's validation pattern:

```hcl
# This might fail validation if the pattern doesn't allow certain characters
resource "azurecaf_name" "example" {
  name          = "my-app_name"
  resource_type = "azurerm_storage_account"  # Only allows alphanumeric
  clean_input   = false  # Won't clean invalid characters
}
# Error: Pattern validation failed
```

### Best Practices for Avoiding Truncation

1. **Keep base names short** - The `name` parameter should be concise
2. **Limit prefixes/suffixes** - Use only essential prefixes and suffixes
3. **Consider resource constraints** - Check maximum lengths for your resource types
4. **Use abbreviations** - Consider shorter alternatives for common terms
5. **Test composition** - Use the data source version to preview names during planning

```hcl
# Good: Short, descriptive components
resource "azurecaf_name" "example" {
  name          = "api"
  resource_type = "azurerm_storage_account"
  prefixes      = ["prod"]
  suffixes      = ["001"]
  random_length = 3
}
# Result: "stprodapi001abc" (15 chars - well within 24 char limit)
```

This systematic approach ensures that generated names are always valid, predictable, and comply with Azure resource naming requirements while maximizing the use of available character space.

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

This resource supports **300+ Azure resource types** with accurate naming validation rules. 

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

For a complete list of supported resource types with their constraints and validation rules, see the [Provider Index](../index.md#supported-azure_resource_types) documentation.
