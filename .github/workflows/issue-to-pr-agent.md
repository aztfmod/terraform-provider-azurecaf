---
description: |
  Automated workflow that processes issues labeled 'new-resource' to
  research the resource, create a branch, add the resource definition,
  build, test, and open a pull request.

on:
  issues:
    types: [labeled]

permissions:
  contents: write
  pull-requests: write
  issues: read

network: defaults

tools:
  github:
    lockdown: true
    toolsets: [issues, pull_requests]
  bash:
    - "git *"
    - "go *"
    - "make *"
    - "grep *"
    - "sed *"
    - "cat *"
    - "python3 *"
    - "curl *"

safe-outputs:
  create-pull-request:
    title-prefix: "[auto] "
    labels: [automated, new-resource]
  add-comment: {}

source: local
engine: copilot
---

# Issue to PR Agent

When an issue is labeled `new-resource`, automatically research the requested Azure resource and create a PR to add it.

## Trigger condition

Only proceed if the label added is `new-resource`. Otherwise, exit immediately.

## Process

### 1. Parse the issue
Extract the Azure resource type name from the issue title and body. Look for patterns like:
- `azurerm_<resource_name>`
- "Add support for <resource name>"
- Resource provider namespace mentions

### 2. Check if already supported
Search `resourceDefinition.json` for the resource name:
```bash
grep '"name": "<resource_name>"' resourceDefinition.json
```
If found, comment on the issue: "This resource is already supported." and close with label `already-supported`.

### 3. Research naming constraints
Follow the azure-naming-research skill procedure:
1. Look up CAF abbreviation on the Microsoft Learn page
2. Look up Azure naming rules
3. Derive field values

### 4. Create branch and add resource
```bash
git checkout -b add/<resource_name>
```
- Add the JSON entry to `resourceDefinition.json` in alphabetical order
- Run `go generate && make build`
- Update `CHANGELOG.md`
- Commit all changes

### 5. Open pull request
Create a PR with:
- Title: `Add <resource_name> resource type`
- Body: Include the researched constraints, JSON entry, and reference to the issue
- Link to the originating issue with `Closes #<issue_number>`

### 6. Comment on issue
Comment: "PR #<pr_number> has been created to add this resource."
