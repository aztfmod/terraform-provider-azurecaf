---
name: azure-caf-sync
description: "Fetch the latest CAF abbreviations page from Microsoft Learn and compare against resourceDefinition.json official.slug values. Reports drift where the provider's slug differs from the official CAF slug. Triggers on: weekly scheduled audit, manual sync check."
---

# Azure CAF Sync

## Procedure

### 1. Fetch CAF abbreviations

Use Microsoft Docs MCP tools (`microsoft_docs_fetch`) to fetch the latest version:

**URL**: `https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations`

Extract all rows from the abbreviation table: resource name, abbreviation (slug), resource provider namespace.

### 2. Extract current provider slugs

```bash
python3 -c "
import json
with open('resourceDefinition.json') as f:
    resources = json.load(f)
for r in resources:
    official = r.get('official', {})
    ns = official.get('resource_provider_namespace', '')
    official_slug = official.get('slug', '')
    print(f'{r[\"name\"]}|{r[\"slug\"]}|{official_slug}|{ns}')
" > /tmp/caf-current.txt
```

### 3. Compare

For each resource in the provider that has an `official.resource_provider_namespace`:
- Find the matching entry in the CAF abbreviations page
- Compare the provider's `slug` against the official CAF abbreviation
- Flag mismatches

### 4. Report

```
## CAF Slug Drift Report

### Mismatches (<count>)
| Resource | Provider Slug | Official CAF Slug | Namespace |
|----------|--------------|-------------------|-----------|
| <name>   | <current>    | <official>        | <ns>      |

### New in CAF (not in provider) (<count>)
| Resource | CAF Slug | Namespace |
|----------|----------|-----------|
| <name>   | <slug>   | <ns>      |

### Summary
- Resources checked: <count>
- Matches: <count>
- Mismatches: <count>
- New resources in CAF: <count>
```
