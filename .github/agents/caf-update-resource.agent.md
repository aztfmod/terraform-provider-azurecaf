# Update Azure Resource: $ARGUMENTS

You are updating an existing resource definition for `$ARGUMENTS` in the terraform-provider-azurecaf project. Follow every step below precisely.

## Skills used

| Skill | File | Used in |
|-------|------|---------|
| Resource Definition JSON | `.github/skills/resource-definition-json/SKILL.md` | Steps 1, 3 |
| Azure Naming Research | `.github/skills/azure-naming-research/SKILL.md` | Step 2 |
| Provider Build & Test | `.github/skills/provider-build-test/SKILL.md` | Step 4 |
| Terraform Mock Test | `.github/skills/terraform-mock-test/SKILL.md` | Step 5 |
| Changelog Update | `.github/skills/changelog-update/SKILL.md` | Step 6 |

---

## STEP 1: Look up the current definition

Follow the **Lookup** operation in the `resource-definition-json` skill.

Search for `$ARGUMENTS` in `resourceDefinition.json`. If not found, also try `azurerm_$ARGUMENTS`.

- If **NOT found**: STOP and inform the user. Suggest using the `caf.add-resource` agent instead.
- If **found**: extract all current field values and proceed to Step 2.

---

## STEP 2: Research current naming constraints

Follow the `azure-naming-research` skill to query Microsoft documentation for the latest naming constraints for `$ARGUMENTS`.

---

## STEP 3: Compare and apply changes

Follow the **Compare** operation in the `resource-definition-json` skill.

Present the comparison table:

```
| Field            | Current    | Expected   | Status   |
|------------------|------------|------------|----------|
| slug             | ...        | ...        | MATCH    |
| min_length       | ...        | ...        | MISMATCH |
| ...              | ...        | ...        | ...      |
```

- If **all fields match**: inform the user that the resource is up to date. STOP.
- If **mismatches found**: show the proposed changes and ask for confirmation.
- If confirmed, apply the changes using the **Update** operation.

---

## STEP 4: Build and test

Follow the `provider-build-test` skill:

```bash
go generate
make build
```

All tests must pass. If build fails, diagnose and fix, then retry.

---

## STEP 5: Validate with mocked azurerm provider

**This step is MANDATORY** — always perform it after Step 4 succeeds.

Follow the `terraform-mock-test` skill for end-to-end validation.

---

## STEP 6: Update CHANGELOG

**Only after ALL tests pass (Steps 4 and 5).**

Follow the `changelog-update` skill to add an entry under `## [Unreleased]` → `### Changed`.

Use the "Updated resource" template:
```markdown
- **RESOURCE**: Updated `$ARGUMENTS` -- <describe field changes>
  - Impact: <Low|Medium|High> -- <brief justification>
```

---

## Summary

After all steps complete, present a summary:

```
✅ Resource: $ARGUMENTS
   Changes: <list of changed fields>
   Build: PASSED
   Mock test: PASSED
   CHANGELOG: Updated
```
