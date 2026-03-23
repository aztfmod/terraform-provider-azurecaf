# Check Azure Resource Definition: $ARGUMENTS

You are validating the resource definition for `$ARGUMENTS` in the terraform-provider-azurecaf project. This provider generates Azure-compliant resource names following the Cloud Adoption Framework (CAF). Follow every step below precisely.

## Skills used

This agent uses the following skill instruction files. Read each skill file when you reach the step that references it:

| Skill | File | Used in |
|-------|------|---------|
| Azure Naming Research | `.github/instructions/azure-naming-research/SKILL.md` | Step 2 |
| Resource Definition JSON | `.github/instructions/resource-definition-json/SKILL.md` | Steps 1, 3, 4, 5 |
| Provider Build & Test | `.github/instructions/provider-build-test/SKILL.md` | Step 5 |
| Terraform Mock Test | `.github/instructions/terraform-mock-test/SKILL.md` | Step 6 |

---

## STEP 1: Look up the resource in resourceDefinition.json

Follow the **Lookup** operation in the `resource-definition-json` skill.

Search for `$ARGUMENTS` in `resourceDefinition.json`. If not found, also try `azurerm_$ARGUMENTS`.

- If **found**: extract all fields and proceed to Step 2 to verify against documentation.
- If **NOT found**: note as **MISSING** and proceed to Step 2.

---

## STEP 2: Research naming constraints

Follow the `azure-naming-research` skill to query Microsoft documentation for the correct naming constraints for `$ARGUMENTS`.

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
This will produce:
- CAF slug and resource provider namespace (from CAF abbreviations page)
- Length, scope, valid characters, and derived regex patterns (from Azure naming rules page)

---

## STEP 3: Compare and report

Follow the **Compare** operation in the `resource-definition-json` skill.

### If the resource EXISTS in resourceDefinition.json:

Present the comparison table showing current vs. researched values with MATCH/MISMATCH/MISSING status for each field.

If all fields match, say: "All fields are up to date. No changes needed." -- then STOP (no further steps needed).

### If the resource is MISSING:

Say: "Resource `$ARGUMENTS` is NOT FOUND in resourceDefinition.json. A new entry is needed."

---

## STEP 4: Propose changes

Follow the **Format** operation in the `resource-definition-json` skill to produce the complete JSON entry.

Show the proposed JSON entry to the user. Ask whether to apply the changes.

---

## STEP 5: Apply changes, build, and test

If the user agrees to apply the changes:

1. Follow the **Insert / Update** operation in the `resource-definition-json` skill to edit `resourceDefinition.json`.
2. Follow the `provider-build-test` skill to regenerate Go code, build, and run tests.

**NOTE**: Do NOT update `CHANGELOG.md` at this stage. The CHANGELOG is updated only in Step 7, after ALL tests pass.

---

## STEP 6: Validate with mocked azurerm provider

**This step is MANDATORY** -- always perform it after Step 5 succeeds.

Follow the `terraform-mock-test` skill to:
1. Install the locally built provider
2. Look up the azurerm resource required attributes (via Terraform MCP tools)
3. Create a test configuration with `mock_provider "azurerm"`
4. Run `terraform test -verbose`
5. Report results with the generated name and full test output
6. Clean up the temporary directory

---

## STEP 7: Update CHANGELOG

**This step is performed ONLY after ALL previous steps have passed successfully** -- including the `make build` tests (Step 5) and the mocked azurerm provider test (Step 6).

Update `CHANGELOG.md` under the `[Unreleased]` section with the change. Use this format:

For new resources:

```markdown
- **RESOURCE**: Added `$ARGUMENTS` (<official resource display name>) -- slug: `<slug>`, length: <min>-<max>, scope: <scope>, dashes <allowed|not allowed>, <lowercase|mixed case>
  - Impact: Low -- new resource type added, no existing behavior changed
```

For updates to existing resources:

```markdown
- **RESOURCE**: Updated `$ARGUMENTS` -- <describe field changes>
  - Impact: Low -- naming constraint update, no breaking changes
```
