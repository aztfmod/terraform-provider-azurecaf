---
name: resource-definition-json
description: "Manage resource entries in resourceDefinition.json for the terraform-provider-azurecaf naming provider. Use for lookup, comparison, formatting, and insert/update of resource naming definitions. Triggers on: resourceDefinition.json edits, resource definition lookup, JSON entry formatting, resource comparison."
---

# Resource Definition JSON

## Lookup

```bash
grep -n '"name": "<resource_name>"' resourceDefinition.json
```

If found, read ~15 lines to extract: `name`, `min_length`, `max_length`, `validation_regex`, `scope`, `slug`, `dashes`, `lowercase`, `regex`, `official.slug`, `official.resource`, `official.resource_provider_namespace`, `out_of_doc`.

If not found, also try `azurerm_<resource_name>`.

## Compare

Present a comparison table:

```
| Field | Current | Expected | Status |
|-------|---------|----------|--------|
| slug  | ...     | ...      | MATCH / MISMATCH |
```

Check all fields. Statuses: **MATCH**, **MISMATCH**, **MISSING**, **UNKNOWN**.

## Format

See [references/json-format.md](references/json-format.md) for the exact JSON entry template and formatting rules.

## Insert / Update

- Maintain **alphabetical ordering** by `name` field
- Ensure proper comma separation between entries
- After any edit, the caller MUST run the `provider-build-test` skill
