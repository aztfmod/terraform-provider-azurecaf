# JSON Entry Format

Every entry in `resourceDefinition.json` MUST follow this exact structure:

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

## Rules

- 4-space indentation (matching existing file)
- `validation_regex` and `regex`: wrap regex in escaped double quotes `"\"pattern\""`
- `official.slug` should match root-level `slug`
- If NOT in official CAF abbreviations page: omit `official.slug` and `official.resource_provider_namespace`, add `"out_of_doc": true`
