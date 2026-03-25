---
name: issue-to-resource-spec
description: "Parse a 'new resource request' GitHub issue and extract the resource type, slug, and naming constraints to produce a draft resourceDefinition.json entry. Triggers on: issue labeled 'feature' or 'new-resource', new resource request."
---

# Issue to Resource Spec

## Procedure

### 1. Parse the issue

Extract from the issue title and body:
- **Resource type**: The Terraform resource name (e.g., `azurerm_container_app`)
- **Slug hint**: Any suggested CAF abbreviation
- **Constraints hint**: Any mentioned length/scope/character rules
- **Azure service**: The Azure service name for research

### 2. Research constraints

If the issue does not provide complete naming constraints, use the `azure-naming-research` skill procedure:
1. Look up CAF abbreviation on the Microsoft Docs page
2. Look up Azure naming rules
3. Derive field values (regex, dashes, lowercase, scope)

### 3. Draft JSON entry

Using the researched values, format a complete JSON entry following the `resource-definition-json` skill's **Format** rules:

```json
{
    "name": "<resource_name>",
    "min_length": <number>,
    "max_length": <number>,
    "validation_regex": "\"^<pattern>$\"",
    "scope": "<scope>",
    "slug": "<caf_abbreviation>",
    "dashes": <true|false>,
    "lowercase": <true|false>,
    "regex": "\"[^<allowed_chars>]\"",
    "official": {
        "slug": "<caf_abbreviation>",
        "resource": "<Official resource display name>",
        "resource_provider_namespace": "<Microsoft.Provider/resourceType>"
    }
}
```

### 4. Output

Present the draft entry and a summary:
```
Resource: <resource_name>
Slug: <slug>
Length: <min>-<max>
Scope: <scope>
Status: Ready for implementation
```

This output can be used directly by the `caf.add-resource` agent.
