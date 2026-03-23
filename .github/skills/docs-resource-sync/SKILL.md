---
name: docs-resource-sync
description: "Ensure documentation files under docs/resources/ and docs/data-sources/ reflect current resource types with accurate examples. Triggers on: after adding/updating resources, documentation audit."
---

# Docs Resource Sync

## Procedure

### 1. Get current resource count

```bash
python3 -c "
import json
with open('resourceDefinition.json') as f:
    resources = json.load(f)
print(f'Total resources: {len(resources)}')
# Count by category
categories = {}
for r in resources:
    prefix = r['name'].replace('azurerm_', '').split('_')[0]
    categories[prefix] = categories.get(prefix, 0) + 1
for cat, count in sorted(categories.items(), key=lambda x: -x[1])[:20]:
    print(f'  {cat}: {count}')
"
```

### 2. Check docs accuracy

Read the current documentation files:
- `docs/resources/azurecaf_name.md` — main resource documentation
- `docs/data-sources/azurecaf_name.md` — data source documentation
- `docs/index.md` — provider index

Verify:
- Resource count matches `resourceDefinition.json`
- Example resource types mentioned in docs actually exist
- Supported resource type tables/lists are current

### 3. Update resource type tables

If the docs include a resource type listing (table or details block), update it to match the current `resourceDefinition.json` entries. Include:
- Resource type name
- Slug
- Min/Max length
- Example generated name

### 4. Verify examples

Check that example Terraform configurations in docs use valid resource types and produce valid outputs:
- Resource type names must exist in `resourceDefinition.json`
- Slug values must match
- Example outputs should be realistic

### 5. Report

```
Docs sync complete:
- docs/resources/azurecaf_name.md: <updated|current>
- docs/data-sources/azurecaf_name.md: <updated|current>
- docs/index.md: <updated|current>
- Resource count in docs: <count>
```
