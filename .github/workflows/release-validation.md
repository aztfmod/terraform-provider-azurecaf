---
description: |
  Release validation workflow that runs comprehensive checks when a
  version tag is pushed. Posts validation results as a comment on the
  associated GitHub Release.

on:
  push:
    tags: ["v*"]

permissions:
  contents: read
  pull-requests: read

network: defaults

tools:
  github:
    lockdown: true
    toolsets: [pull_requests]
  bash:
    - "make *"
    - "go *"
    - "grep *"
    - "cat *"

safe-outputs:
  add-comment: {}

source: local
engine: copilot
---

# Release Validation

When a version tag is pushed, run the full validation suite and report results.

## Process

### 1. Build and test
- Run `make build` to verify compilation
- Run `make test_ci` for the full test suite
- Run `make test_e2e` for E2E validation
- Run `make test_coverage` for coverage check

### 2. Verify CHANGELOG
- Check that the tag version appears in `CHANGELOG.md`
- Verify the version section has entries

### 3. Report
Create a summary of all validation results. If all checks pass, the release is validated. If any fail, flag the issues for immediate attention.

```markdown
### 🏷️ Release Validation: vX.Y.Z

| Check | Status |
|-------|--------|
| Build | ✅ / ❌ |
| Unit tests | ✅ / ❌ |
| E2E tests | ✅ / ❌ |
| Coverage | ✅ <percentage>% / ❌ <percentage>% |
| CHANGELOG entry | ✅ / ❌ |

Overall: **VALIDATED** / **ISSUES FOUND**
```
