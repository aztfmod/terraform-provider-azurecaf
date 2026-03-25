---
name: pr-compliance-check
description: "Validate a pull request against the project's contribution checklist. Checks that resourceDefinition.json changes trigger models_generated.go regeneration, README is updated, CHANGELOG is updated, and tests pass. Triggers on: PR opened, PR updated."
---

# PR Compliance Check

## Procedure

### 1. Identify changed files

```bash
git diff --name-only origin/main...HEAD
```

### 2. Check compliance rules

| Rule | Condition | Check |
|------|-----------|-------|
| **Generated code fresh** | `resourceDefinition.json` changed | `models_generated.go` must also be changed |
| **CHANGELOG updated** | Any code change | `CHANGELOG.md` must be modified |
| **README updated** | Resource added/removed | `README.md` must be modified |
| **Tests pass** | Any code change | `make build` must succeed |

### 3. Verify generated code

If `resourceDefinition.json` is in the diff:

```bash
go generate
git diff --name-only azurecaf/models_generated.go
```

If there are uncommitted changes to `models_generated.go`, the PR has stale generated code.

### 4. Report

**If compliant:**
```
✅ PR Compliance: PASSED
   - [x] Generated code is fresh
   - [x] CHANGELOG updated
   - [x] README updated (if applicable)
   - [x] Tests pass
```

**If issues found:**
```
⚠️ PR Compliance: ISSUES FOUND
   - [ ] <issue description>
   - [x] <passing check>
```

Provide specific remediation steps for each failing check.
