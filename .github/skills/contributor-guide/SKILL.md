---
name: contributor-guide
description: "Provide step-by-step guidance for contributors based on their contribution type. References CONTRIBUTING.md and project conventions. Triggers on: new contributor questions, 'how do I contribute', contribution guidance requests."
---

# Contributor Guide

## Procedure

### 1. Identify contribution type

Ask the contributor (or infer from context) what they want to do:

| Type | Description |
|------|-------------|
| **new-resource** | Add a new Azure resource type to the provider |
| **update-resource** | Update an existing resource's naming constraints |
| **bug-fix** | Fix a bug in the provider logic |
| **docs** | Improve documentation |
| **test** | Add or improve tests |

### 2. Provide guidance

**For new-resource contributions:**

1. Check if the resource already exists: `grep '"name": "azurerm_<name>"' resourceDefinition.json`
2. Research naming constraints (see `azure-naming-research` skill)
3. Create an issue describing the resource (or reference existing issue)
4. Edit `resourceDefinition.json` — add the entry in alphabetical order
5. Run `go generate && make build` to regenerate code and test
6. Update `CHANGELOG.md` under `## [Unreleased]` → `### Added`
7. Update `README.md` resource status table
8. Submit PR using the template in `.github/PULL_REQUEST_TEMPLATE.md`

**For update-resource contributions:**

1. Find current definition: `grep -A 15 '"name": "azurerm_<name>"' resourceDefinition.json`
2. Research current Azure naming rules
3. Update the entry in `resourceDefinition.json`
4. Run `go generate && make build`
5. Update `CHANGELOG.md` under `## [Unreleased]` → `### Changed`
6. Submit PR

**For bug-fix contributions:**

1. Reproduce the bug with a test case
2. Fix the code
3. Verify fix with `make build`
4. Update `CHANGELOG.md` under `## [Unreleased]` → `### Fixed`
5. Submit PR

**For docs contributions:**

1. Edit files under `docs/` or `README.md`
2. Verify markdown renders correctly
3. Submit PR

### 3. Reference

Point the contributor to:
- `.github/CONTRIBUTING.md` for full guidelines
- `.github/PULL_REQUEST_TEMPLATE.md` for the PR checklist
- `TESTING.md` for test instructions
