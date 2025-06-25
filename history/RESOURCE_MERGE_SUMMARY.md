# Resource Definition Merge Summary

This document summarizes the changes made to combine `resourceDefinition.json` and `resourceDefinition_out_of_docs.json` with official Azure documentation mapping.

## Changes Made

### 1. File Consolidation
- **Before**: Two separate files
  - `resourceDefinition.json` (364 resources)
  - `resourceDefinition_out_of_docs.json` (31 resources)
- **After**: Single `resourceDefinition.json` (395 resources)

### 2. New Resource Attributes
Added three new fields to all resources:

- `out_of_doc` (boolean): Indicates resources not in official Azure CAF documentation
- `resource` (string): Official resource name from Azure CAF documentation  
- `resource_provider_namespace` (string): Azure resource provider namespace

### 3. Official Documentation Mapping
Implemented mapping for key Azure resources per the official documentation:

| Resource | Official Name | Resource Provider Namespace | Slug |
|----------|---------------|----------------------------|------|
| azurerm_kubernetes_cluster | AKS cluster | Microsoft.ContainerService/managedClusters | aks |
| azurerm_container_app | Container apps | Microsoft.App/containerApps | ca |
| azurerm_container_app_environment | Container apps environment | Microsoft.App/managedEnvironments | cae |
| azurerm_storage_account | Storage account | Microsoft.Storage/storageAccounts | st |
| azurerm_resource_group | Resource group | Microsoft.Resources/resourceGroups | rg |
| azurerm_virtual_machine | Virtual machine | Microsoft.Compute/virtualMachines | vm |
| azurerm_key_vault | Key Vault | Microsoft.KeyVault/vaults | kv |
| azurerm_app_service | App Service | Microsoft.Web/sites | app |
| azurerm_virtual_network | Virtual network | Microsoft.Network/virtualNetworks | vnet |
| azurerm_subnet | Subnet | Microsoft.Network/virtualNetworks/subnets | snet |

### 4. Code Generation Updates
- Updated `ResourceStructure` in `gen.go` to include new fields
- Simplified file reading logic to use single combined file
- Maintained backward compatibility with existing resource definitions

### 5. Out-of-Documentation Resources
31 resources marked with `out_of_doc: true`, including:
- azurerm_private_service_connection
- azurerm_firewall_ip_configuration  
- azurerm_firewall_application_rule_collection
- azurerm_dns_*_record types
- And others not yet in official Azure CAF documentation

### 6. Testing and Validation
- All existing unit tests pass
- Added new tests to validate the merge
- Build process works correctly
- Resource generation and code compilation successful

### 7. Documentation Updates
- Updated Makefile to reflect single file structure
- Added merge script to `scripts/` directory for future maintenance
- Enhanced .gitignore for better file management

## Files Changed
- `resourceDefinition.json` - Combined and enhanced resource definitions
- `gen.go` - Updated ResourceStructure and file reading logic
- `Makefile` - Updated resource table generation command
- `azurecaf/models_generated.go` - Regenerated with new structure
- `azurecaf/resource_definition_merge_test.go` - New validation tests
- `scripts/merge_resource_definitions.py` - Merge automation script
- `.gitignore` - Enhanced file exclusions

## Files Removed
- `resourceDefinition_out_of_docs.json` - No longer needed

## Future Maintenance
The `scripts/merge_resource_definitions.py` script can be used to:
1. Add new official documentation mappings
2. Update resource attributes
3. Handle future Azure CAF documentation changes

## Validation Results
- ✅ 395 total resources (364 + 31)
- ✅ 31 resources marked as out_of_doc
- ✅ 10 resources with official documentation mapping
- ✅ All unit tests passing
- ✅ Build and generation working correctly

This implementation satisfies all requirements specified in issue #331.