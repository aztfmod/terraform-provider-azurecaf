---
name: example-generator
description: "Generate Terraform example configurations for a given resource type using the azurecaf provider. Produces ready-to-use HCL code for documentation or onboarding. Triggers on: documentation requests, example creation, contributor guidance."
---

# Example Generator

## Procedure

### 1. Look up resource definition

```bash
grep -A 15 '"name": "<resource_name>"' resourceDefinition.json
```

Extract: `name`, `slug`, `min_length`, `max_length`, `scope`, `dashes`, `lowercase`.

### 2. Generate basic example

```hcl
resource "azurecaf_name" "example" {
  name          = "demo"
  resource_type = "<resource_name>"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 3
  clean_input   = true
}

output "name" {
  value = azurecaf_name.example.result
}
# Expected output format: <slug>-dev-demo-<random>-001
# Length: <min_length>-<max_length> characters
```

### 3. Generate data source example

```hcl
data "azurecaf_name" "example" {
  name          = "demo"
  resource_type = "<resource_name>"
  prefixes      = ["prod"]
  suffixes      = ["eu"]
}

output "name" {
  value = data.azurecaf_name.example.result
}
```

### 4. Generate multi-resource example

If the user wants examples showing multiple resources:

```hcl
resource "azurecaf_name" "multi" {
  name          = "myapp"
  resource_types = [
    "<resource_name_1>",
    "<resource_name_2>",
  ]
  prefixes      = ["dev"]
  random_length = 3
}

output "results" {
  value = azurecaf_name.multi.results
}
```

### 5. Generate passthrough example

```hcl
resource "azurecaf_name" "passthrough" {
  name          = "my-existing-name"
  resource_type = "<resource_name>"
  passthrough   = true
}
# Output: my-existing-name (validated against naming rules only)
```
