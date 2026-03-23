---
description: |
  Welcome first-time contributors with a friendly message, link to
  CONTRIBUTING.md, and run basic compliance checks on their PR.

on:
  pull_request:
    types: [opened]

permissions:
  contents: read
  pull-requests: read

network: defaults

tools:
  github:
    lockdown: false
    min-integrity: none
    toolsets: [pull_requests]

safe-outputs:
  add-comment: {}

source: local
engine: copilot
---

# Contributor Welcome

When a pull request is opened, check if the author is a first-time contributor to this repository.

## Process

### 1. Check contributor status
Look at the PR author's previous contributions to this repository. If this is their first PR:

### 2. Welcome message
Post a comment:

```markdown
### 👋 Welcome to terraform-provider-azurecaf!

Thanks for your first contribution, @<author>! We're excited to have you here.

Here are some resources to help:
- 📋 [Contributing Guide](.github/CONTRIBUTING.md) — step-by-step instructions
- 🧪 [Testing Guide](TESTING.md) — how to run tests locally
- ✅ [PR Checklist](.github/PULL_REQUEST_TEMPLATE.md) — make sure all boxes are checked

**Quick checklist for resource changes:**
- [ ] Updated `resourceDefinition.json`
- [ ] Ran `go generate` to regenerate `models_generated.go`
- [ ] Updated `CHANGELOG.md`
- [ ] Updated `README.md` resource table
- [ ] All tests pass (`make build`)

A maintainer will review your PR shortly. Feel free to ask questions! 🚀
```

### 3. For returning contributors
If the author has contributed before, no welcome message is needed. Exit silently.
