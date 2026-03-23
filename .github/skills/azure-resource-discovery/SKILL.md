---
name: azure-resource-discovery
description: "Discover new Azure resources by fetching the latest azurerm provider resource list from the Terraform Registry and comparing against resourceDefinition.json. Identifies unsupported resources. Triggers on: weekly scheduled discovery, manual audit."
---

# Azure Resource Discovery

## Procedure

### 1. Fetch azurerm provider resources

Use Terraform MCP tools (`get_provider_capabilities`) to list all resources in the azurerm provider:

```
provider_name: azurerm
provider_namespace: hashicorp
provider_document_type: resources
```

Alternatively, query the Terraform Registry API:
```bash
curl -s "https://registry.terraform.io/v1/providers/hashicorp/azurerm" | python3 -c "
import json, sys
data = json.load(sys.stdin)
print(json.dumps(data, indent=2))
"
```

### 2. Extract current provider resources

```bash
grep '"name":' resourceDefinition.json | sed 's/.*"name": "//;s/".*//' | sort > /tmp/ard-supported.txt
```

### 3. Compare

```bash
# Identify resources in azurerm but not in our provider
comm -23 /tmp/ard-azurerm.txt /tmp/ard-supported.txt > /tmp/ard-missing.txt

# Count
MISSING=$(wc -l < /tmp/ard-missing.txt)
SUPPORTED=$(wc -l < /tmp/ard-supported.txt)
```

### 4. Categorize missing resources

Group missing resources by Azure service category:
- Compute (vm, vmss, disk, etc.)
- Storage (storage_account, blob, etc.)
- Networking (vnet, subnet, nsg, etc.)
- Database (sql, cosmosdb, mysql, etc.)
- Security (keyvault, etc.)
- Other

### 5. Report

```
## Azure Resource Discovery Report

### Summary
- azurerm provider resources: <total>
- Supported in CAF provider: <supported>
- Missing from CAF provider: <missing>
- Coverage: <percentage>%

### Missing Resources by Category
#### Compute (<count>)
- `azurerm_<resource>` — <description if available>

#### Storage (<count>)
- `azurerm_<resource>`

...

### Priority Recommendations
Resources that should be added first (commonly used, well-documented):
1. `azurerm_<resource>` — <reason>
2. ...
```

### 6. Cleanup

```bash
rm -f /tmp/ard-*.txt
```
