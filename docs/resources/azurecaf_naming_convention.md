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
