---
description: |
  PR review agent that validates pull request compliance. Checks that
  resourceDefinition.json changes include regenerated models, CHANGELOG
  is updated, and tests pass. Runs on PR open and update.

on:
  pull_request:
    types: [opened, synchronize]

permissions:
  contents: read
  pull-requests: read

network: defaults

tools:
  github:
    lockdown: true
    toolsets: [pull_requests]
  bash:
    - "git *"
    - "go *"
    - "make *"
    - "grep *"
    - "diff *"
    - "cat *"
    - "wc *"

safe-outputs:
  add-comment: {}

source: local
engine: copilot
---

# PR Review Agent

Review the pull request for compliance with project conventions.

## Checks

### 1. Generated code freshness
If `resourceDefinition.json` is modified in this PR:
- Run `go generate` and check if `azurecaf/models_generated.go` has uncommitted changes
- If stale, comment with: "Generated code is stale. Please run `go generate` and commit the result."

### 2. CHANGELOG updated
If any `.go` or `.json` file is modified:
- Check that `CHANGELOG.md` is also in the diff
- If missing, comment with: "Please update CHANGELOG.md with a description of your changes."

### 3. Resource diff summary
If `resourceDefinition.json` is modified:
- Compare the PR version against main
- Count added, removed, and modified resources
- Comment with a summary table of changes

### 4. Build verification
- Run `make build` to verify the PR compiles and tests pass
- If failures, comment with the error output

## Comment format

Post a single review comment combining all results:

```markdown
### 🤖 PR Compliance Review

| Check | Status |
|-------|--------|
| Generated code fresh | ✅ / ❌ |
| CHANGELOG updated | ✅ / ❌ |
| Build passes | ✅ / ❌ |

<details if resource changes>
#### Resource Changes
| Resource | Change | Details |
|----------|--------|---------|
| ... | Added/Modified/Removed | ... |
</details>
```
