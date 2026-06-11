---
description: |
  Labels closed issues and merged PRs with the release version they were
  included in. Cross-references git tags, merge commits, and CHANGELOG
  entries to determine which release shipped each item.

on:
  push:
    tags: ["v*"]
  workflow_dispatch:

permissions:
  contents: read
  issues: read
  pull-requests: read

network: defaults

tools:
  github:
    lockdown: false
    min-integrity: none
    allowed-repos: all
  bash:
    - "git *"
    - "grep *"
    - "cat *"
    - "jq *"
    - "sed *"
    - "sort *"
    - "comm *"

safe-outputs:
  mentions: false
  allowed-github-references: []
  add-labels: {}
  create-issue:
    title-prefix: "[Release Labeler] "
    max: 1

source: local
engine: copilot
timeout-minutes: 45
---

# Release Labeler

Label closed issues and merged pull requests with the release version tag they were included in.

## Goal

After a release is tagged (or on manual trigger for backfilling), determine which issues and PRs were resolved in each release and add a version label (e.g., `v1.2.32`) to them.

## Process

### 1. Determine releases to process

**On tag push (`v*`):**
- Process only the newly pushed tag.
- Find the previous tag to establish the commit range.

**On manual trigger (`workflow_dispatch`):**
- List all release tags (`git tag --sort=v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$'`).
- Process tags in batches, starting from the most recent. Focus on efficiency — extract all references with a single `git log` per tag pair using `--format` flags.
- Use bash to do bulk extraction of PR/issue numbers from commit messages rather than making individual API calls per commit.

### 2. For each release tag

Find the commit range between consecutive tags:

```bash
# Get ordered list of release tags (exclude previews, test tags)
git tag --sort=v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$'

# For each consecutive pair (prev_tag, curr_tag):
git log $prev_tag..$curr_tag --oneline
```

### 3. Extract issue and PR references

From the commits in each release range, extract:

1. **Merged PR numbers** from merge commit messages (e.g., `Merge pull request #399`)
2. **Referenced issue numbers** from commit messages (e.g., `fixes #379`, `resolves #291`, `#89`)
3. **PR numbers** from branch-based commits that were squash-merged

Use the GitHub API to also check:
- Each merged PR's linked/closing issues
- Issue cross-references in PR bodies

### 4. Create version labels

For each version tag being processed:
- Use the `add_labels` tool to apply the version label (e.g., `v1.2.32`) to issues and PRs. The tool will handle label creation if needed.
- If `add_labels` fails because the label doesn't exist, skip that item and note it in the report.

### 5. Apply labels

For each issue and PR identified in a release:
- Check if it already has the version label.
- If not, add the version label.
- Skip open issues/PRs — only label closed issues and merged PRs.

### 6. Report

After processing, create a summary:

```markdown
## Release Labeler Report

### vX.Y.Z
- Labeled N issues: #1, #2, #3
- Labeled N PRs: #10, #20, #30
- Skipped N (already labeled)

### vX.Y.Z-1
...
```

Post the report as a comment on the latest release (for tag push) or as a new issue (for backfill runs).

## Important rules

- Never label open issues or unmerged PRs.
- If a PR or issue appears in multiple release ranges (e.g., referenced in a backport), label it with the earliest release.
- Do not remove existing version labels — only add missing ones.
- Skip pre-release tags (e.g., `v1.2.24-preview`, `v2.0.0-preview-1`).
- Skip test tags (e.g., `v1.2.31-test`).

## Efficiency guidelines

- Use bash to extract all PR/issue references in bulk:
  ```bash
  git log $prev_tag..$curr_tag --oneline | grep -oE '#[0-9]+' | sort -u
  ```
- Minimize GitHub API calls — prefer `search_issues` with batch queries over individual issue lookups.
- When checking if a label is already applied, use `search_issues` with `label:vX.Y.Z` to find already-labeled items and skip them.
