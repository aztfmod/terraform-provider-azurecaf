---
name: naming-rules-drift-check
description: "Re-fetch Azure naming rules for a set of resources and compare against current regex/length/scope values in resourceDefinition.json. Detects when Azure has changed naming constraints. Triggers on: periodic audit, resource validation."
---

# Naming Rules Drift Check

## Procedure

### 1. Select resources to check

If a specific list is provided, use it. Otherwise, select a batch of resources:

```bash
python3 -c "
import json
with open('resourceDefinition.json') as f:
    resources = json.load(f)
# Select resources with official namespace (can verify against docs)
official = [r for r in resources if r.get('official', {}).get('resource_provider_namespace')]
for r in official[:20]:  # Check 20 at a time
    ns = r['official']['resource_provider_namespace']
    print(f'{r[\"name\"]}|{ns}|{r[\"min_length\"]}|{r[\"max_length\"]}|{r[\"scope\"]}')
"
```

### 2. Fetch current Azure naming rules

For each resource's `resource_provider_namespace`, look up the naming rules:

**URL**: `https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules`

Use Microsoft Docs MCP tools (`microsoft_docs_fetch`) to fetch and parse the rules page. Find the section matching the resource provider namespace.

### 3. Compare

For each resource, compare:
- `min_length` vs documented minimum
- `max_length` vs documented maximum
- `scope` vs documented scope
- `dashes` vs whether hyphens are documented as valid
- `lowercase` vs whether only lowercase is documented
- `validation_regex` vs pattern derived from documented valid characters

### 4. Report

```
## Naming Rules Drift Report

### Drifted Resources (<count>)
| Resource | Field | Current | Azure Docs | Action |
|----------|-------|---------|------------|--------|
| <name>   | max_length | 63 | 128 | Update |

### Verified Resources (<count>)
All other checked resources match Azure documentation.

### Summary
- Resources checked: <count>
- Matches: <count>
- Drift detected: <count>
```
