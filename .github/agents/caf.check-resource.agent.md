# Check Azure Resource Definition: $ARGUMENTS

You are validating the resource definition for `$ARGUMENTS` in the terraform-provider-azurecaf project. This provider generates Azure-compliant resource names following the Cloud Adoption Framework (CAF). Follow every step below precisely.

---

## STEP 1: Look up the resource in resourceDefinition.json

Search for `$ARGUMENTS` in the file `resourceDefinition.json` at the project root.

- Use Grep to find `"name": "$ARGUMENTS"` in `resourceDefinition.json`
- If found, read the full JSON object for this resource (approximately 15 lines from the `"name"` field to the closing `}`)
- Extract and record every field:
  - `name`, `min_length`, `max_length`, `validation_regex`, `scope`, `slug`, `dashes`, `lowercase`, `regex`
  - `official.slug`, `official.resource`, `official.resource_provider_namespace`
  - `out_of_doc` (if present)

If the resource is NOT found in resourceDefinition.json:
- Note it as **MISSING**
- Also check if the user may have omitted the `azurerm_` prefix. Try searching for `azurerm_$ARGUMENTS` as well.
- Proceed to Step 2 to research what the correct definition should be.

---

## STEP 2: Query Microsoft documentation for the correct naming constraints

You need to check TWO Microsoft documentation sources. Perform these lookups in parallel when possible. If MCP tools for Azure or Microsoft Docs are available, prefer those. Otherwise, use WebFetch.

### 2a: CAF Resource Abbreviations

Fetch the CAF abbreviations page to find the official slug and resource provider namespace:

**URL**: `https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations`

Search the page for the resource by:
1. The Terraform resource name without the `azurerm_` prefix (e.g., for `azurerm_container_app`, search for "container app")
2. The resource provider namespace if already known from the existing definition
3. The resource's common name (e.g., "Container Apps", "Storage account")

Extract:
- **Official CAF abbreviation** (the recommended slug)
- **Resource display name**
- **Resource provider namespace** (e.g., `Microsoft.App/containerApps`)

### 2b: Azure Resource Name Rules

Fetch the naming rules page to find the actual naming constraints:

**URL**: `https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules`

This page is organized by resource provider namespace (e.g., `Microsoft.Storage`, `Microsoft.App`). To find the right section:
1. Use the resource provider namespace from Step 2a (e.g., for `Microsoft.App/containerApps`, look under "Microsoft.App")
2. Find the specific entity (e.g., `containerApps`)

Extract from the naming rules table:
- **Scope** (Resource group, Subscription, Global, etc.)
- **Length** (min-max, usually presented as "1-63" format)
- **Valid characters** (the prose description of allowed characters, start/end rules)

### 2c: Convert naming rules to JSON fields

From the naming rules prose, derive these fields:

**min_length** and **max_length**: Take directly from the length range in the docs.

**dashes**: `true` if hyphens/dashes are mentioned as valid characters, `false` otherwise.

**lowercase**: `true` if ONLY lowercase letters are allowed (no uppercase mentioned), `false` if mixed case is allowed.

**scope**: Map the docs scope using:
- "globally unique" or "global" -> `"global"`
- "resource group" -> `"resourceGroup"`
- "within parent", "within vault", "parent resource" -> `"parent"`
- "subscription" -> `"subscription"`
- "region" -> `"region"`

**validation_regex**: Construct a regex that directly matches the **documented constraints**: required start/end characters (if any) and the allowed character set, with `{MIN,MAX}` quantifiers based on the `min_length`/`max_length` values.

| Rule description | validation_regex (example) |
|-----------------|----------------------------|
| Lowercase alphanumeric only | `"\"^[a-z0-9]{MIN,MAX}$\""` |
| Alphanumeric only (mixed case) | `"\"^[a-zA-Z0-9]{MIN,MAX}$\""` |
| Alphanumeric + hyphens | `"\"^[a-zA-Z0-9-]{MIN,MAX}$\""` |
| Alphanumeric + hyphens + underscores | `"\"^[a-zA-Z0-9_-]{MIN,MAX}$\""` |
| Alphanumeric + hyphens + underscores + periods | `"\"^[a-zA-Z0-9_.-]{MIN,MAX}$\""` |

IMPORTANT:
- The regex value MUST be wrapped in escaped double quotes: `"\"^pattern$\""`. This is a project convention.
- If the docs **require a specific starting or ending character** (for example, “must start with a letter and end with an alphanumeric character”), construct the pattern explicitly from that rule instead of blindly applying a `{MIN-2,MAX-2}` formula.
  - In such cases, only split the pattern into first/middle/last segments when `MIN` and `MAX` are both **at least 2**.
  - Derive the middle quantifier from the actual lengths. For example, for “start with a letter, end with alphanumeric, length MIN–MAX, allowed [a-zA-Z0-9-] in the middle”, use: `"\"^[a-zA-Z][a-zA-Z0-9-]{MIN-2,MAX-2}[a-zA-Z0-9]$\""`, but **only** when `MIN >= 2` and `MAX >= 2`.
  - When `MIN < 2` or when there is **no** fixed start/end requirement, prefer a single character class with `{MIN,MAX}` as in the table above.

**regex** (cleaning regex): The inverse pattern that matches characters to REMOVE. Common mappings:

| Allowed characters | Cleaning regex |
|-------------------|---------------|
| Lowercase alphanumeric | `"\"[^0-9a-z]\""` |
| Alphanumeric (mixed case) | `"\"[^0-9A-Za-z]\""` |
| Alphanumeric + hyphens | `"\"[^0-9A-Za-z-]\""` |
| Alphanumeric + hyphens + underscores | `"\"[^0-9A-Za-z_-]\""` |
| Alphanumeric + hyphens + underscores + periods | `"\"[^0-9A-Za-z_.-]\""` |

---

## STEP 3: Compare and report

### If the resource EXISTS in resourceDefinition.json:

Present a comparison table:

```
## Resource: $ARGUMENTS

| Field | Current | Microsoft Docs | Status |
|-------|---------|---------------|--------|
| slug | ... | ... | MATCH / MISMATCH |
| min_length | ... | ... | MATCH / MISMATCH |
| max_length | ... | ... | MATCH / MISMATCH |
| dashes | ... | ... | MATCH / MISMATCH |
| lowercase | ... | ... | MATCH / MISMATCH |
| scope | ... | ... | MATCH / MISMATCH |
| validation_regex | ... | ... | MATCH / MISMATCH |
| regex | ... | ... | MATCH / MISMATCH |
| official.slug | ... | ... | MATCH / MISMATCH / MISSING |
| official.resource | ... | ... | MATCH / MISMATCH / MISSING |
| official.resource_provider_namespace | ... | ... | MATCH / MISMATCH / MISSING |
```

Use these status labels:
- **MATCH** -- values are identical or semantically equivalent
- **MISMATCH** -- values differ; the docs value should be preferred
- **MISSING** -- the field is absent but should be present
- **UNKNOWN** -- could not determine the correct value from docs

If all fields match, say: "All fields are up to date. No changes needed."

### If the resource is MISSING:

Say: "Resource `$ARGUMENTS` is NOT FOUND in resourceDefinition.json. A new entry is needed."

---

## STEP 4: Propose changes

If there are discrepancies or the resource is missing, output the complete JSON entry to add or replace in `resourceDefinition.json`.

The JSON entry MUST follow this exact format:

```json
{
    "name": "$ARGUMENTS",
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

Formatting rules:
- Use 4-space indentation (matching the existing file)
- The `validation_regex` and `regex` values must have their regex wrapped in escaped double quotes: `"\"pattern\""`
- `official.slug` should match the root-level `slug` when present
- If the resource is NOT in the official CAF abbreviations page, omit `official.slug` and `official.resource_provider_namespace`, and add `"out_of_doc": true`

For **existing resources** with discrepancies, show what should replace the current entry.
For **missing resources**, show the complete new entry to insert (maintaining alphabetical order by `name`).

---

## STEP 5: Apply changes and verify

Ask the user whether to apply the changes. If they agree:

1. Edit `resourceDefinition.json` to add or update the entry (maintain alphabetical ordering by the `name` field)
2. Run `go generate` from the project root to regenerate `azurecaf/models_generated.go`
3. Run `make build` to build and test the provider
4. Show a sample Terraform configuration to verify:

```hcl
data "azurecaf_name" "test" {
  name          = "testname"
  resource_type = "$ARGUMENTS"
  prefixes      = ["dev"]
  random_length = 3
}

output "result" {
  value = data.azurecaf_name.test.result
}
```

5. Update `CHANGELOG.md` with the change (under the `[Unreleased]` section)

---

## STEP 6: Test with the locally built provider

After the build succeeds in Step 5, validate the resource end-to-end using the locally built provider and Terraform CLI.

### 6a: Install the provider locally

Run these commands from the project root to install the freshly built binary where Terraform can find it:

```bash
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
LOCAL_PLUGIN_DIR=~/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/${GOOS}_${GOARCH}
mkdir -p "$LOCAL_PLUGIN_DIR"
cp ./terraform-provider-azurecaf "$LOCAL_PLUGIN_DIR/"
```

### 6b: Create a temporary Terraform configuration

Create a temporary working directory (e.g., `/tmp/azurecaf-test-$ARGUMENTS`) with the following files:

**providers.tf**
```hcl
terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = ">= 1.2.0"
    }
  }
}

provider "azurecaf" {
}
```

**test.tf** — generate a name for the resource:
```hcl
resource "azurecaf_name" "test" {
  name          = "testname"
  resource_type = "$ARGUMENTS"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 3
  clean_input   = true
}

output "result" {
  value = azurecaf_name.test.result
}

output "result_length" {
  value = length(azurecaf_name.test.result)
}
```

### 6c: Run Terraform plan and apply

Run the following in the temporary directory, using the dev_overrides config from the `examples/terraform.rc` file:

```bash
export TF_CLI_CONFIG_FILE=<project_root>/examples/terraform.rc
terraform plan
terraform apply -auto-approve
```

Note: With `dev_overrides`, `terraform init` is not required.

### 6d: Validate the output

After `terraform apply` completes:

1. **Check the generated name** — the `result` output should:
   - Start with the prefix: `dev`
   - Contain the slug for the resource type (if the naming convention uses slugs)
   - Contain `testname` (or a cleaned/truncated version of it)
   - End with the suffix `001` followed by a random string
   - Respect the `max_length` constraint (verify with the `result_length` output)
   - Only contain characters allowed by the `validation_regex`

2. **Verify naming constraints** — confirm the generated name matches the regex from the resource definition:
   - If `lowercase` is `true`, the name must contain only lowercase characters
   - If `dashes` is `false`, the name must not contain hyphens
   - The total length must be between `min_length` and `max_length`

3. Report the test result:
   - Show the generated name
   - Show the length
   - Confirm it passes validation or describe the failure

### 6e: Validate with mocked azurerm provider

**This step is MANDATORY** — always perform it after step 6d. Use Terraform's `terraform test` framework with `mock_provider "azurerm"` to simulate deploying the actual Azure resource using the generated name. This verifies the name is accepted by the azurerm provider schema — without requiring Azure credentials.

#### Determine the matching azurerm resource

Map `$ARGUMENTS` to its corresponding azurerm resource. The resource type name IS the azurerm resource (e.g., `azurerm_container_app` maps to the Terraform resource `azurerm_container_app`). Look up the minimum required attributes for that resource in the azurerm provider documentation.

#### Create the test configuration

Create a single `main.tf` file in the temporary working directory with both providers and the azurerm resource. **IMPORTANT**: Use `aztfmod/azurecaf` as the provider source (not `aztfmod.com/arnaudlh/azurecaf`), and use hardcoded fake Azure resource IDs for parent references instead of chaining azurerm resources (mock provider generates random IDs that break cross-resource references like parsed resource IDs).

**main.tf** — use the generated name in the real azurerm resource:
```hcl
terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = ">= 1.2.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0.0"
    }
  }
}

provider "azurecaf" {}
provider "azurerm" {
  features {}
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "azurecaf_name" "test" {
  name          = "testname"
  resource_type = "$ARGUMENTS"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 3
  clean_input   = true
}

# MANDATORY: The azurerm resource MUST use the generated name from azurecaf.
# Set: name = azurecaf_name.test.result
# This chains the CAF-generated name into the actual Azure resource,
# proving the name is accepted by the azurerm provider schema.
#
# Use ONLY the minimum required attributes for the resource.
# For parent resource references, use HARDCODED fake Azure resource IDs.
# Do NOT chain azurerm resources — mock provider generates random IDs that
# break cross-resource references (e.g., parsed storage account IDs).
#
# Example for azurerm_synapse_spark_pool:
# resource "azurerm_synapse_spark_pool" "test" {
#   name                 = azurecaf_name.test.result
#   synapse_workspace_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test/providers/Microsoft.Synapse/workspaces/synwstest"
#   node_size_family     = "MemoryOptimized"
#   node_size            = "Small"
#   node_count           = 3
#   spark_version        = "3.4"
# }

output "result" {
  value = azurecaf_name.test.result
}

output "result_length" {
  value = length(azurecaf_name.test.result)
}
```

Fill in the actual resource block based on the specific `$ARGUMENTS` resource type. Use the azurerm provider documentation to determine required attributes.

Also create a `terraform.rc` in the temporary directory to override the azurecaf provider with the local build:

```hcl
provider_installation {
  dev_overrides {
    "aztfmod/azurecaf" = "<LOCAL_PLUGIN_DIR>"
  }
  direct {}
}
```

Where `<LOCAL_PLUGIN_DIR>` is the path used in step 6a (e.g., `~/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/darwin_arm64`).

**tests/validate_name.tftest.hcl** — mock the azurerm provider and validate:
```hcl
mock_provider "azurerm" {}

run "validate_generated_name" {
  command = apply

  # Validate the output from the azurecaf_name resource is properly set
  assert {
    condition     = output.result != ""
    error_message = "Generated name output is empty"
  }

  # Validate length constraints via outputs
  assert {
    condition     = output.result_length >= <min_length>
    error_message = "Generated name is shorter than min_length (<min_length>)"
  }

  assert {
    condition     = output.result_length <= <max_length>
    error_message = "Generated name exceeds max_length (<max_length>)"
  }

  # Validate the generated name matches the validation regex
  assert {
    condition     = can(regex("<validation_regex_pattern>", output.result))
    error_message = "Generated name does not match validation regex"
  }

  # CRITICAL: Validate the azurerm resource received the generated name from the output
  # This proves the CAF-generated name flows through outputs and is accepted by the azurerm provider schema
  assert {
    condition     = <azurerm_resource>.test.name == output.result
    error_message = "Azure resource name does not match generated CAF name"
  }
}
```

Replace `<min_length>`, `<max_length>`, `<validation_regex_pattern>`, and `<azurerm_resource>` with the values from the resource definition.

**CRITICAL RULES**:
- Use `command = apply`, NOT `command = plan`. The `azurecaf_name.result` attribute is computed and not known during plan phase, so plan-time assertions will fail with "Unknown condition value".
- **ALWAYS use `output.result` and `output.result_length`** in assertions instead of `azurecaf_name.test.result` directly. This ensures:
  1. The outputs are evaluated and displayed in verbose test output
  2. The test validates the full data flow: `azurecaf_name` → `output` → assertion
  3. The `Outputs:` section in verbose mode shows the generated name and length
- **ALWAYS compare the azurerm resource's `name` attribute against `output.result`** (not against `azurecaf_name.test.result`). This validates the complete chain: azurecaf generates the name → output captures it → azurerm resource uses it.

#### Run the test

```bash
# First, run terraform init to download the azurerm provider schema
terraform init

# Then run the test with the local azurecaf override
# IMPORTANT: Always use -verbose to show the generated name and resource attributes
export TF_CLI_CONFIG_FILE=<temp_dir>/terraform.rc
terraform test -verbose
```

`terraform init` is required to download the azurerm provider binary (needed for schema validation even with mocking). The `dev_overrides` in `terraform.rc` ensure the locally built azurecaf provider is used instead of the registry version.

**IMPORTANT**: Always use `-verbose` flag. This shows:
- The full `azurecaf_name` resource state including the generated `result`
- The full `azurerm_*` resource state confirming the name was accepted
- The output values (`result` and `result_length`)

With `mock_provider "azurerm"`, the azurerm provider schema is loaded and validated but no real Azure API calls are made. The test confirms:
1. The azurerm provider accepts the generated name in its `name` attribute
2. The apply succeeds without schema validation errors
3. The name meets length and regex constraints
4. The azurerm resource's `name` attribute matches the generated CAF name (validated via `output.result`)
5. The outputs (`result` and `result_length`) are properly evaluated and displayed

#### Report the mock test result

ALWAYS show the full verbose test output to the user, including:
- The generated name value from the `result` output (must NOT be empty in the `Outputs:` section)
- The `result_length` output (must NOT be empty in the `Outputs:` section)
- The azurerm resource block showing the name was accepted
- The pass/fail status of each assertion

If the `Outputs:` section is empty, the test assertions are NOT using `output.result` / `output.result_length` — this is a bug. Fix the assertions to reference outputs instead of resource attributes directly.

- If all assertions pass: report success and show the full verbose test output
- If the apply fails with a schema error: the generated name may contain invalid characters or violate a provider-side constraint — report the error details
- If assertions fail: report which constraint was violated

### 6f: Clean up

Remove the temporary directory:
```bash
rm -rf /tmp/azurecaf-test-$ARGUMENTS
```

---

## Execution notes

- When fetching web pages, if an MCP tool for Azure or Microsoft Docs is available, prefer that over WebFetch
- If WebFetch or MCP calls fail, clearly state which source could not be reached and what information is missing
- Never guess at naming constraints -- if you cannot verify a value, mark it as UNKNOWN and explain what needs manual verification
- The `validation_regex` is the most critical field -- it must accurately reflect Azure's naming rules
- Check for the resource using both the full Terraform name (`azurerm_container_app`) and the short form (`container_app`) when searching documentation
