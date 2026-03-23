# Audit Azure Resources

You are performing a comprehensive audit of the terraform-provider-azurecaf resource definitions. This audit checks completeness, naming rules drift, test coverage, and CAF slug alignment.

## Skills used

| Skill | File | Used in |
|-------|------|---------|
| Resource Completeness Check | `.github/skills/resource-completeness-check/SKILL.md` | Step 1 |
| Azure CAF Sync | `.github/skills/azure-caf-sync/SKILL.md` | Step 2 |
| Naming Rules Drift Check | `.github/skills/naming-rules-drift-check/SKILL.md` | Step 3 |
| Coverage Analysis | `.github/skills/coverage-analysis/SKILL.md` | Step 4 |

---

## STEP 1: Completeness check

Follow the `resource-completeness-check` skill to compare `resourceDefinition.json` against `completness/existing_tf_resources.txt`.

Report the coverage percentage and list missing resources.

---

## STEP 2: CAF slug alignment

Follow the `azure-caf-sync` skill to fetch the latest CAF abbreviations and compare against the provider's slugs.

Report any slug mismatches or new resources in the CAF documentation.

---

## STEP 3: Naming rules drift

Follow the `naming-rules-drift-check` skill to verify naming constraints against Azure documentation.

Check a batch of resources (prioritize official ones with `resource_provider_namespace`).

Report any drift in length, scope, valid characters, or regex patterns.

---

## STEP 4: Test coverage

Follow the `coverage-analysis` skill to run tests and verify coverage meets the 95% threshold.

---

## Summary

Present a combined audit report:

```
# 🔍 Resource Audit Report

## Completeness
- Coverage: <percentage>% (<supported>/<total> resources)
- Missing: <count> resources

## CAF Slug Alignment
- Checked: <count> resources
- Matches: <count>
- Mismatches: <count>
- New in CAF: <count>

## Naming Rules
- Checked: <count> resources
- Up to date: <count>
- Drifted: <count>

## Test Coverage
- Coverage: <percentage>%
- Threshold: 95%
- Status: PASS / FAIL

## Recommended Actions
1. <highest priority action>
2. <next priority action>
...
```
