# azurecaf_naming_convention

The `azurecaf_naming_convention` resource provides a legacy approach to generating Azure resource names following the Microsoft Cloud Adoption Framework (CAF) naming conventions. This resource implements predefined naming methodologies with a fixed set of Azure resource types.

> **Note**: For new projects, consider using the [`azurecaf_name` resource](azurecaf_name.md) which offers more flexibility and supports a broader range of Azure resource types.

## Key Features

- **CAF-compliant naming**: Follows Microsoft's official [Cloud Adoption Framework naming conventions](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging)
- **Multiple conventions**: Supports four different naming convention methods
- **Automatic validation**: Enforces Azure resource naming constraints and character limitations
- **Flexible patterns**: Combines prefix, basename, suffix, and convention-specific padding

## Example Usage

### Basic Resource Group Naming

```hcl
resource "azurecaf_naming_convention" "example_rg" {
  name          = "myapp"
  prefix        = "dev"
  resource_type = "rg"
  postfix       = "001"
  max_length    = 23
  convention    = "cafrandom"
}

resource "azurerm_resource_group" "example" {
  name     = azurecaf_naming_convention.example_rg.result
  location = "East US"
}

# Output: "dev-myapp-rg-001-wxyz" (where wxyz is randomly generated)
```

### Multiple Resources with Consistent Naming

```hcl
# Storage Account
resource "azurecaf_naming_convention" "storage" {
  name          = "myapp"
  prefix        = "prod"
  resource_type = "st"
  convention    = "cafrandom"
}

# Virtual Network
resource "azurecaf_naming_convention" "vnet" {
  name          = "myapp"
  prefix        = "prod"
  resource_type = "vnet"
  convention    = "cafrandom"
}

# Use in resources
resource "azurerm_storage_account" "example" {
  name                = azurecaf_naming_convention.storage.result
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  # ... other configuration
}
```

## Argument Reference

### Required Arguments

- `resource_type` - (Required) The type of Azure resource for naming. See [supported resource types](#supported-resource-types) below.

### Optional Arguments

- `name` - (Optional) The base name of the resource. Will be sanitized according to Azure naming requirements. Defaults to an empty string.
- `convention` - (Optional) The naming convention method to use. Defaults to `"cafrandom"`. Allowed values:
  - `"cafclassic"` - Standard CAF naming without random padding
  - `"cafrandom"` - CAF naming with random padding to max length
  - `"random"` - Fully random name within Azure constraints  
  - `"passthrough"` - Manual naming (filtered for length and invalid characters)
- `prefix` - (Optional) Prefix prepended to the generated name.
- `postfix` - (Optional) Suffix appended after the base name (useful for indexes like "001").
- `max_length` - (Optional) Maximum length of the generated name. If longer than the Azure resource limit, the resource limit applies.

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
resource "azurecaf_naming_convention" "example" {
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
resource "azurecaf_naming_convention" "example" {
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
resource "azurecaf_naming_convention" "example" {
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
resource "azurecaf_naming_convention" "example" {
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
resource "azurecaf_naming_convention" "example" {
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
resource "azurecaf_naming_convention" "example" {
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
resource "azurecaf_naming_convention" "example" {
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

- `id` - The unique identifier of the naming convention resource.
- `result` - The generated name for the Azure resource based on input parameters and the selected convention.

## Naming Convention Methods

| Method | Description |
|--------|-------------|
| `cafclassic` | Follows standard CAF recommendations with the pattern: `[prefix]-[name]-[resource_type]-[postfix]` |
| `cafrandom` | CAF pattern with random padding to reach maximum allowed length |
| `random` | Generates a completely random name within Azure resource constraints |
| `passthrough` | Uses provided components as-is, but validates length and character constraints |

## Name Pattern

The generated names follow this pattern (depending on convention):

```
[prefix]-[name]-[resource_type_abbreviation]-[postfix]-[padding]
```

**Example breakdown:**
- Input: `prefix="dev"`, `name="myapp"`, `resource_type="rg"`, `postfix="001"`, `convention="cafrandom"`
- Output: `"dev-myapp-rg-001-wxyz"` (where "wxyz" is random padding)

## Supported Resource Types

| Resource Type | Short Code | Long Code |
|---------------|------------|-----------|
| Azure Automation | `aaa` | `azurerm_automation_account` |
| Container App | `ca` | `azurerm_container_app` |
| Container App Environment | `cae` | `azurerm_container_app_environment` |
| Container Registry | `acr` | `azurerm_container_registry` |
| Azure Firewall | `afw` | `azurerm_firewall` |
| Application Gateway | `agw` | `azurerm_application_gateway` |
| API Management | `apim` | `azurerm_api_management` |
| App Service | `app` | `azurerm_app_service` |
| Application Insights | `appi` | `azurerm_application_insights` |
| App Service Environment | `ase` | `azurerm_app_service_environment` |
| Azure Kubernetes Service | `aks` | `azurerm_kubernetes_cluster` |
| AKS DNS Prefix | `aksdns` | `aks_dns_prefix` |
| AKS Node Pool (Linux) | `aksnpl` | `aks_node_pool_linux` |
| AKS Node Pool (Windows) | `aksnpw` | `aks_node_pool_windows` |
| Recovery Services Vault | `asr` | `azurerm_recovery_services_vault` |
| Event Hubs Namespace | `evh` | `azurerm_eventhub_namespace` |
| Key Vault | `kv` | `azurerm_key_vault` |
| Log Analytics Workspace | `la` | `azurerm_log_analytics_workspace` |
| Network Interface | `nic` | `azurerm_network_interface` |
| Network Security Group | `nsg` | `azurerm_network_security_group` |
| Public IP | `pip` | `azurerm_public_ip` |
| App Service Plan | `plan` | `azurerm_app_service_plan` |
| Service Plan | `plan` | `azurerm_service_plan` |
| Resource Group | `rg` | `azurerm_resource_group` |
| Subnet | `snet` | `azurerm_subnet` |
| SQL Server | `sql` | `azurerm_sql_server` |
| SQL Database | `sqldb` | `azurerm_sql_database` |
| Storage Account | `st` | `azurerm_storage_account` |
| Linux Virtual Machine | `vml` | `azurerm_virtual_machine_linux` |
| Windows Virtual Machine | `vmw` | `azurerm_virtual_machine_windows` |
| Virtual Network | `vnet` | `azurerm_virtual_network` |
| Generic Resource | `gen` | `generic` |

## Migration Notes

### Upgrading to azurecaf_name

For new deployments, consider migrating to the more flexible [`azurecaf_name` resource](azurecaf_name.md):

```hcl
# Legacy approach
resource "azurecaf_naming_convention" "legacy" {
  name          = "myapp"
  resource_type = "rg"
  convention    = "cafrandom"
}

# Modern approach (equivalent)
resource "azurecaf_name" "modern" {
  name          = "myapp"
  resource_type = "azurerm_resource_group"
  random_length = 4
}
```

## Related Resources

- [`azurecaf_name` resource](azurecaf_name.md) - Modern, flexible resource naming (recommended)
- [`azurecaf_name` data source](../data-sources/azurecaf_name.md) - Data source for name generation
- [`azurecaf_environment_variable` data source](../data-sources/azurecaf_environment_variable.md) - Environment variable support

For more resource types and advanced features, see the [Provider Index](../index.md) documentation.
