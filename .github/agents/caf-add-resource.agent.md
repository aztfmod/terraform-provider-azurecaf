# Add Azure Resource: $ARGUMENTS

You are adding a new resource definition for `$ARGUMENTS` to the terraform-provider-azurecaf project. This provider generates Azure-compliant resource names following the Cloud Adoption Framework (CAF). Follow every step below precisely.

## Skills used

| Skill | File | Used in |
|-------|------|---------|
| Azure Naming Research | `.github/skills/azure-naming-research/SKILL.md` | Step 1 |
| Resource Definition JSON | `.github/skills/resource-definition-json/SKILL.md` | Steps 2, 3 |
| Provider Build & Test | `.github/skills/provider-build-test/SKILL.md` | Step 4 |
| Terraform Mock Test | `.github/skills/terraform-mock-test/SKILL.md` | Step 5 |
| Changelog Update | `.github/skills/changelog-update/SKILL.md` | Step 6 |
| README Resource Table | `.github/skills/readme-resource-table/SKILL.md` | Step 7 |

---

## STEP 1: Research naming constraints

Follow the `azure-naming-research` skill to query Microsoft documentation for the correct naming constraints for `$ARGUMENTS`.

1. Look up the CAF abbreviation (slug, resource display name, resource provider namespace)
2. Look up Azure naming rules (scope, length, valid characters)
3. Derive field values (regex patterns, dashes, lowercase, scope mapping)

---

## STEP 2: Check if resource already exists

Follow the **Lookup** operation in the `resource-definition-json` skill.

Search for `$ARGUMENTS` in `resourceDefinition.json`. If not found, also try `azurerm_$ARGUMENTS`.

- If **found**: STOP and inform the user. Suggest using the `caf.update-resource` agent instead.
- If **NOT found**: proceed to Step 3.

---

## STEP 3: Create the JSON entry

Follow the **Format** and **Insert** operations in the `resource-definition-json` skill.

1. Format the complete JSON entry using the researched values from Step 1.
2. Show the proposed JSON entry to the user and ask for confirmation.
3. If confirmed, insert the entry into `resourceDefinition.json` in alphabetical order by `name`.

---

## STEP 4: Build and test

Follow the `provider-build-test` skill:

```bash
go generate
make build
```

Verify the resource appears in `azurecaf/models_generated.go`. All tests must pass.

If build fails, diagnose and fix the JSON entry, then retry.

---

## STEP 5: Validate with mocked azurerm provider

**This step is MANDATORY** — always perform it after Step 4 succeeds.

Follow the `terraform-mock-test` skill to:
1. Install the locally built provider
2. Look up the azurerm resource required attributes
3. Create a test configuration with `mock_provider "azurerm"`
4. Run `terraform test -verbose`
5. Report results
6. Clean up

---

## STEP 6: Update CHANGELOG

**Only after ALL tests pass (Steps 4 and 5).**

Follow the `changelog-update` skill to add an entry under `## [Unreleased]` → `### Added`.

Use the "New resource" template:
```markdown
- **RESOURCE**: Added `$ARGUMENTS` (<official display name>) -- slug: `<slug>`, length: <min>-<max>, scope: <scope>, dashes <allowed|not allowed>, <lowercase|mixed case>
  - Impact: Low -- new resource type added, no existing behavior changed
```

---

## STEP 7: Update README

Follow the `readme-resource-table` skill to add the new resource to the status table in `README.md`.

---

## Summary

After all steps complete, present a summary:

```
✅ Resource: $ARGUMENTS
   Slug: <slug>
   Length: <min>-<max>
   Scope: <scope>
   Build: PASSED
   Mock test: PASSED
   CHANGELOG: Updated
   README: Updated
```
