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
    - "terraform *"

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

      # 3. Schema-based nameability check using terraform provider schema
      #    Only resources with a required user-controlled `name` attribute are nameable
      mkdir -p "$GH_AW_OUT/_schema_check"
      cat > "$GH_AW_OUT/_schema_check/main.tf" <<'TF'
terraform {
  required_providers {
    azurerm = { source = "hashicorp/azurerm", version = "~> 4.0" }
  }
}
TF
      (cd "$GH_AW_OUT/_schema_check" && terraform init -input=false -no-color >/dev/null 2>&1 && \
       terraform providers schema -json > "$GH_AW_OUT/provider_schema.json" 2>/dev/null)
      SCHEMA_OK=$?

      if [ $SCHEMA_OK -eq 0 ] && [ -s "$GH_AW_OUT/provider_schema.json" ]; then
        echo "SCHEMA_CHECK=ok" >> "$GH_AW_OUT/counts.env"
        # Filter: keep only resources that have a required `name` attribute
        python3 -c "
import json, sys, os

with open('$GH_AW_OUT/provider_schema.json') as f:
    schema = json.load(f)

# Find azurerm provider
azurerm = None
for key, val in schema.get('provider_schemas', {}).items():
    if 'azurerm' in key:
        azurerm = val
        break

if not azurerm:
    sys.exit('No azurerm schema found')

resources = azurerm.get('resource_schemas', {})

# Also load static exclusions as fallback for resources not in current schema
excl_file = 'completness/non_nameable_resources.json'
static_excl = set()
not_caf_eligible = set()
if os.path.exists(excl_file):
    with open(excl_file) as f:
        excl = json.load(f)
    static_excl = set(excl.get('exact_resources', []) + excl.get('deprecated_in_v4', []))
    not_caf_eligible = set(excl.get('name_not_caf_eligible', []))
    suffix_pats = excl.get('suffix_patterns', [])
    contains_pats = excl.get('contains_patterns', [])
else:
    suffix_pats, contains_pats = [], []

for line in open('$GH_AW_OUT/azurerm_raw.txt'):
    r = line.strip()
    if not r:
        continue

    # Check via provider schema (definitive)
    if r in resources:
        attrs = resources[r].get('block', {}).get('attributes', {})
        name_attr = attrs.get('name', {})
        has_name = 'name' in attrs
        name_required = name_attr.get('required', False)
        # Keep only if name is required AND it's a CAF-eligible name
        if has_name and name_required and r not in not_caf_eligible:
            print(r)
        continue

    # Not in schema (deprecated/removed) — use static exclusions
    if r in static_excl:
        continue
    if any(r.endswith(s) for s in suffix_pats):
        continue
    if any(p in r for p in contains_pats):
        continue
    # Unknown resource not in schema and not excluded — skip (likely deprecated)
    # Don't flag it since we can't verify nameability
" > "$GH_AW_OUT/azurerm.txt"
      else
        echo "SCHEMA_CHECK=failed" >> "$GH_AW_OUT/counts.env"
        # Fallback: use static exclusion list only
        python3 -c "
import json, sys
with open('completness/non_nameable_resources.json') as f:
    excl = json.load(f)
suffix_pats = excl.get('suffix_patterns', [])
contains_pats = excl.get('contains_patterns', [])
exact = set(excl.get('exact_resources', []) + excl.get('deprecated_in_v4', []) + excl.get('name_not_caf_eligible', []))
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
      fi

      RAW_COUNT=$(wc -l < "$GH_AW_OUT/azurerm_raw.txt" | tr -d ' ')
      FILTERED_COUNT=$(wc -l < "$GH_AW_OUT/azurerm.txt" | tr -d ' ')
      echo "Filtered: $RAW_COUNT raw → $FILTERED_COUNT nameable ($((RAW_COUNT - FILTERED_COUNT)) excluded)"

      # 4. Diff: resources in azurerm.txt but missing from supported.txt
      comm -23 "$GH_AW_OUT/azurerm.txt" "$GH_AW_OUT/supported.txt" > "$GH_AW_OUT/missing.txt"

      # 5. Counts
      {
        echo "AZURERM_TOTAL=$FILTERED_COUNT"
        echo "SUPPORTED_TOTAL=$(wc -l < "$GH_AW_OUT/supported.txt" | tr -d ' ')"
        echo "MISSING_TOTAL=$(wc -l < "$GH_AW_OUT/missing.txt" | tr -d ' ')"
      } >> "$GH_AW_OUT/counts.env"

      # 6. CAF abbreviations page (best-effort; agent should note if empty)
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
- `azurerm.txt` — filtered: only resources with a **required** `name` attribute in the provider schema. Uses `terraform providers schema -json` for definitive nameability check, with static fallback if schema fetch fails.
- `missing.txt` — nameable azurerm resources NOT present in `supported.txt`.
- `counts.env` — `AZURERM_TOTAL`, `SUPPORTED_TOTAL`, `MISSING_TOTAL`, `CAF_FETCH`, `SCHEMA_CHECK`.
- `provider_schema.json` — full azurerm provider schema (if `SCHEMA_CHECK=ok`).
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
