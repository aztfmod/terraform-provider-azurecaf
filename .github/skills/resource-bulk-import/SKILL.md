---
name: resource-bulk-import
description: "Accept a list of Azure resource types, batch-research each via azure-naming-research conventions, generate JSON entries, and insert all into resourceDefinition.json. Triggers on: batch resource addition, bulk import requests."
---

# Resource Bulk Import

## Procedure

### 1. Accept resource list

The caller provides a list of resource type names (e.g., `azurerm_container_app`, `azurerm_container_registry`).

### 2. Filter out existing resources

```bash
for RESOURCE in <list>; do
  if grep -q "\"name\": \"$RESOURCE\"" resourceDefinition.json; then
    echo "SKIP: $RESOURCE (already exists)"
  else
    echo "ADD: $RESOURCE"
  fi
done
```

### 3. Batch research

For each resource to add, follow the `azure-naming-research` skill procedure:
1. Look up CAF abbreviation
2. Look up Azure naming rules
3. Derive field values

Collect all results before making any edits.

### 4. Generate JSON entries

For each researched resource, format the JSON entry following the `resource-definition-json` skill format rules. Present all entries to the user for review:

```
## Proposed Resources (<count>)

| # | Resource | Slug | Length | Scope | Dashes | Lowercase |
|---|----------|------|--------|-------|--------|-----------|
| 1 | <name>   | <slug> | <min>-<max> | <scope> | <yes/no> | <yes/no> |
...
```

### 5. Insert all entries

After user confirmation, insert all entries into `resourceDefinition.json`:
- Maintain alphabetical order by `name`
- Ensure proper comma separation
- Use 4-space indentation

### 6. Build and test once

After inserting all entries, run a single build/test cycle:

```bash
go generate
make build
```

This is more efficient than building after each individual insertion.

### 7. Report

```
## Bulk Import Complete

Added: <count> resources
Skipped: <count> (already existed)
Build: PASSED / FAILED
Tests: PASSED / FAILED
```
