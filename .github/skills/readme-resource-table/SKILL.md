---
name: readme-resource-table
description: "Update the resource status table in README.md to reflect current resource support. Adds new rows, updates existing entries, and marks support status. Triggers on: after adding/updating resources in resourceDefinition.json."
---

# README Resource Table Update

## Procedure

### 1. Read current README

Locate the resource status table in `README.md`. It is typically inside a `<details>` block or a markdown table listing supported resource types with their slugs and status.

```bash
grep -n "resource_type\|slug\|azurerm_" README.md | head -40
```

### 2. Get current resources

Extract all resource names from `resourceDefinition.json`:

```bash
grep '"name":' resourceDefinition.json | sed 's/.*"name": "//;s/".*//' | sort
```

### 3. Compare and identify gaps

Cross-reference the README table against `resourceDefinition.json`. Identify:
- Resources in JSON but missing from README (need to add rows)
- Resources in README but not in JSON (need to mark as removed or verify)
- Resources with changed slugs or attributes (need to update rows)

### 4. Update table

For each new or changed resource, format the row following the existing table structure in README.md. Maintain alphabetical ordering within each category.

### 5. Verify

Run the deterministic checker. It is also enforced in CI and by `go test ./...`:

```bash
make test_readme_consistency   # fails on drift
make fix_readme_consistency    # auto-fix ✔ status, missing rows, and resource counts
```

Aliases and helpers that intentionally have no row are listed in `internal/readmeconsistency/exclusions.go`.

After updating, also ensure:
- No duplicate rows
- All rows follow the same column format
- Resource count matches `resourceDefinition.json` entry count
