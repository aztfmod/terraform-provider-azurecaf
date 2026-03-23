---
description: |
  Weekly workflow that discovers new Azure resources and checks for CAF 
  slug drift. Creates GitHub issues for gaps and mismatches so maintainers
  can prioritize adding or updating resource definitions.

on:
  schedule: "0 9 * * 1"
  workflow_dispatch:

permissions:
  contents: read
  issues: write

network:
  allowed:
    - defaults
    - learn.microsoft.com
    - registry.terraform.io

tools:
  github:
    lockdown: true
    toolsets: [issues]
  bash:
    - "python3 *"
    - "grep *"
    - "sort *"
    - "comm *"
    - "wc *"
    - "cat *"
    - "curl *"

safe-outputs:
  create-issue:
    title-prefix: "[azure-sync] "
    labels: [automated, azure-sync, enhancement]
    close-older-issues: true

source: local
engine: copilot
---

# Weekly Azure Sync

Check for new Azure resources and CAF slug changes on a weekly basis.

## Process

### 1. Resource discovery
- Extract supported resource names from `resourceDefinition.json`
- Compare against the known Terraform azurerm resources in `completness/existing_tf_resources.txt`
- Identify any new resources that should be added

### 2. CAF slug check
- Fetch the latest CAF abbreviations page from Microsoft Learn
- Compare official slugs against the provider's current slug values
- Identify mismatches

### 3. Create issues
If gaps or drift are found:
- Create a GitHub issue summarizing the findings
- Title: `[azure-sync] <count> new resources and <count> slug changes detected — <date>`
- Include a prioritized list of resources to add or update
- Label each finding by category (compute, storage, networking, etc.)

If no changes detected:
- Exit silently, no issue needed
