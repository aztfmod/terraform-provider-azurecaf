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

steps:
  - name: Set up Go
    uses: actions/setup-go@v6
    with:
      go-version-file: './go.mod'
      cache: true

  - name: Build provider
    id: build
    continue-on-error: true
    run: |
      set -o pipefail
      if make build 2>&1 | tee /tmp/pr-build.log; then
        echo "BUILD_RESULT=pass" > /tmp/pr-review-results.env
      else
        echo "BUILD_RESULT=fail" > /tmp/pr-review-results.env
      fi
      # Keep a short tail for the agent prompt
      tail -n 80 /tmp/pr-build.log > /tmp/pr-build-tail.log || true

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
- The `Build provider` pre-agent step has already run `make build` on the runner
  (where Go from `go.mod` is installed). Read the result from
  `/tmp/pr-review-results.env` — it contains `BUILD_RESULT=pass` or
  `BUILD_RESULT=fail`.
- If `BUILD_RESULT=fail`, include the tail of `/tmp/pr-build-tail.log` (last
  80 lines) in the comment as a fenced code block.
- Do NOT attempt to run `make build` or `go build` from inside the agent
  sandbox — the agent container does not have the Go toolchain. The pre-agent
  step is the source of truth.

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
