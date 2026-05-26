---
description: |
  Weekly workflow that discovers new Azure resources and checks for CAF 
  slug drift. Creates GitHub issues for gaps and mismatches so maintainers
  can prioritize adding or updating resource definitions.

on:
  schedule: weekly on monday
  workflow_dispatch:

permissions:
  contents: read
  issues: read

network:
  allowed:
    - defaults
    - learn.microsoft.com
    - terraform

tools:
  github:
    lockdown: true
    toolsets: [issues]
  bash:
    - "grep *"
    - "sort *"
    - "comm *"
    - "wc *"
    - "cat *"
    - "head *"
    - "tail *"
    - "python3 *"

safe-outputs:
  create-issue:
    title-prefix: "[azure-sync] "
    labels: [automated, azure-sync, enhancement]
    close-older-issues: true

steps:
  - name: Pre-compute resource gap and slug data
    env:
      GH_AW_OUT: ${{ github.workspace }}/.gh-aw-data
    run: |
      set -o pipefail
      mkdir -p "$GH_AW_OUT"

      # 1. Supported resources from the provider's resourceDefinition.json
      python3 -c "import json; print('\n'.join(sorted({r['name'] for r in json.load(open('resourceDefinition.json'))})))" \
        > "$GH_AW_OUT/supported.txt"

      # 2. Known azurerm resources tracked in completness/existing_tf_resources.txt
      sort -u completness/existing_tf_resources.txt > "$GH_AW_OUT/azurerm_raw.txt"

      # 2b. Filter out non-nameable resources (associations, configs, deprecated)
      python3 -c "
import json, sys
with open('completness/non_nameable_resources.json') as f:
    excl = json.load(f)
suffix_pats = excl.get('suffix_patterns', [])
contains_pats = excl.get('contains_patterns', [])
exact = set(excl.get('exact_resources', []) + excl.get('deprecated_in_v4', []))
for line in open('$GH_AW_OUT/azurerm_raw.txt'):
    r = line.strip()
    if not r:
        continue
    if r in exact:
        continue
    if any(r.endswith(s) for s in suffix_pats):
        continue
    if any(p in r for p in contains_pats):
        continue
    print(r)
" > "$GH_AW_OUT/azurerm.txt"

      echo "FILTERED=$(( $(wc -l < "$GH_AW_OUT/azurerm_raw.txt" | tr -d ' ') - $(wc -l < "$GH_AW_OUT/azurerm.txt" | tr -d ' ') )) non-nameable resources excluded"

      # 3. Diff: resources in azurerm.txt but missing from supported.txt
      comm -23 "$GH_AW_OUT/azurerm.txt" "$GH_AW_OUT/supported.txt" > "$GH_AW_OUT/missing.txt"

      # 4. Counts
      {
        echo "AZURERM_TOTAL=$(wc -l < "$GH_AW_OUT/azurerm.txt" | tr -d ' ')"
        echo "SUPPORTED_TOTAL=$(wc -l < "$GH_AW_OUT/supported.txt" | tr -d ' ')"
        echo "MISSING_TOTAL=$(wc -l < "$GH_AW_OUT/missing.txt" | tr -d ' ')"
      } > "$GH_AW_OUT/counts.env"

      # 5. CAF abbreviations page (best-effort; agent should note if empty)
      if curl -fsSL --max-time 30 \
          "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations" \
          > "$GH_AW_OUT/caf-abbreviations.html"; then
        echo "CAF_FETCH=ok" >> "$GH_AW_OUT/counts.env"
      else
        echo "CAF_FETCH=failed" >> "$GH_AW_OUT/counts.env"
        : > "$GH_AW_OUT/caf-abbreviations.html"
      fi

      echo "----- counts -----"
      cat "$GH_AW_OUT/counts.env"
      echo "----- first 20 missing -----"
      head -n 20 "$GH_AW_OUT/missing.txt"

source: local
engine: copilot
---

# Weekly Azure Sync

Check for new Azure resources and CAF slug changes on a weekly basis.

## Pre-computed inputs

The pre-agent step has already produced the following files under
`${GITHUB_WORKSPACE}/.gh-aw-data/` — read them with `cat`, `head`, `tail`,
`grep`, `sort`, `comm`, or `wc`:

- `supported.txt` — sorted unique resource names from `resourceDefinition.json`.
- `azurerm_raw.txt` — ALL azurerm resource names from `completness/existing_tf_resources.txt`.
- `azurerm.txt` — filtered: only resources with user-controlled `name` fields (non-nameable excluded).
- `missing.txt` — azurerm resources NOT present in `supported.txt` (pre-filtered for nameability).
- `counts.env` — `AZURERM_TOTAL`, `SUPPORTED_TOTAL`, `MISSING_TOTAL`, `CAF_FETCH`.
- `caf-abbreviations.html` — Microsoft Learn CAF abbreviations page (empty if `CAF_FETCH=failed`).

Do NOT attempt to run `python3` or `curl` from inside the agent sandbox — those
tools are intentionally not allowed. The data above is the source of truth.

## Process

### 1. Resource discovery
- Read `counts.env` for totals.
- Read `missing.txt` for the gap list.
- Note: `missing.txt` is already pre-filtered to exclude non-nameable resources
  (associations, bindings, configs, policies, deprecated resources) using
  `completness/non_nameable_resources.json`. Only resources with a required
  user-controlled `name` field appear in the list.
- Skip any line beginning with `azurerm_` followed by a name that ends in
  `_data_source`, `_versions`, or otherwise represents a data source rather
  than a manageable resource.

### 2. CAF slug check
- If `CAF_FETCH=ok`, parse `caf-abbreviations.html` with `grep` to extract
  resource → abbreviation pairs and compare against the provider's current
  slug values from `resourceDefinition.json` (use `grep` for spot checks).
- If `CAF_FETCH=failed`, note the failure in the report and continue with the
  resource discovery section only.

### 3. Create issues
If gaps or drift are found:
- Create a GitHub issue summarizing the findings
- Title: `[azure-sync] <count> new resources and <count> slug changes detected — <date>`
- Include a prioritized list of resources to add or update
- Label each finding by category (compute, storage, networking, etc.)

If no changes detected:
- Exit silently, no issue needed
