# Bulk Add Resources: $ARGUMENTS

You are adding multiple resource definitions to the terraform-provider-azurecaf project in a single session. This is more efficient than adding one at a time.

## Skills used

| Skill | File | Used in |
|-------|------|---------|
| Resource Bulk Import | `.github/skills/resource-bulk-import/SKILL.md` | Steps 1-4 |
| Provider Build & Test | `.github/skills/provider-build-test/SKILL.md` | Step 5 |
| Terraform Mock Test | `.github/skills/terraform-mock-test/SKILL.md` | Step 6 |
| Changelog Update | `.github/skills/changelog-update/SKILL.md` | Step 7 |
| README Resource Table | `.github/skills/readme-resource-table/SKILL.md` | Step 8 |

---

## STEP 1: Parse resource list

Extract the list of resources from `$ARGUMENTS`. Accept formats:
- Comma-separated: `azurerm_a, azurerm_b, azurerm_c`
- Space-separated: `azurerm_a azurerm_b azurerm_c`
- Newline-separated list

---

## STEP 2: Filter existing resources

Follow the `resource-bulk-import` skill to check which resources already exist in `resourceDefinition.json`. Report skipped resources.

---

## STEP 3: Batch research

For each new resource, research naming constraints using the `azure-naming-research` conventions (CAF slug, naming rules, regex patterns).

Present all researched entries in a summary table and ask for confirmation.

---

## STEP 4: Insert all entries

After confirmation, insert all entries into `resourceDefinition.json` in alphabetical order.

---

## STEP 5: Build and test (once)

Follow the `provider-build-test` skill — run a single build/test cycle for all added resources:

```bash
go generate
make build
```

---

## STEP 6: Mock test (sample)

Follow the `terraform-mock-test` skill to validate at least 2-3 representative resources from the batch. Pick resources with different constraint profiles (e.g., one lowercase-only, one with dashes, one global scope).

---

## STEP 7: Update CHANGELOG

**Only after ALL tests pass.**

Follow the `changelog-update` skill. Use the batch template:

```markdown
- **RESOURCE**: Added support for <count> new resource types
  - Added `<resource_1>` with slug `<slug_1>`
  - Added `<resource_2>` with slug `<slug_2>`
  - ...
  - Impact: <Low|Medium|High> -- <count> new resource types added
```

---

## STEP 8: Update README

Follow the `readme-resource-table` skill to add all new resources to the status table.

---

## Summary

```
✅ Bulk Import Complete

Resources added: <count>
Resources skipped: <count> (already existed)
Build: PASSED
Mock tests: PASSED (<tested_count>/<total_count> sampled)
CHANGELOG: Updated
README: Updated
```
