---
description: |
  Weekly full sweep of the mock-azurerm validation. Runs the harness against
  every resource in resourceDefinition.json (not just the diff) and opens a
  single categorized GitHub issue summarizing failures. Closes older issues
  with the same prefix so the backlog stays tidy.

on:
  schedule:
    - cron: "0 9 * * 1"  # Mondays at 9:00 UTC
  workflow_dispatch:

permissions:
  contents: read
  issues: read

network: defaults

tools:
  github:
    lockdown: true
    toolsets: [issues]
  bash:
    - "make *"
    - "go *"
    - "scripts/mock-test/*"
    - "python3 *"
    - "grep *"
    - "awk *"
    - "cut *"
    - "sort *"
    - "head *"
    - "tail *"
    - "wc *"
    - "cat *"

safe-outputs:
  create-issue:
    title-prefix: "[mock-azurerm] "
    labels: [bug, automated, mock-azurerm]
    close-older-issues: true

steps:
  - name: Set up Go
    uses: actions/setup-go@v6
    with:
      go-version-file: './go.mod'
      cache: true

  - name: Set up Python
    uses: actions/setup-python@v6
    with:
      python-version: '3.x'

  - name: Setup Terraform
    uses: hashicorp/setup-terraform@v4
    with:
      terraform_version: "~> 1.15"
      terraform_wrapper: false

  - name: Cache azurerm plugin downloads
    uses: actions/cache@v4
    with:
      path: ~/.terraform.d/plugin-cache
      key: tf-plugin-cache-${{ runner.os }}-azurerm-v4

  - name: Run full mock-azurerm sweep
    env:
      CHECKPOINT_DISABLE: "1"
      TF_IN_AUTOMATION: "1"
    run: |
      set -o pipefail
      : > /tmp/results.env

      echo "🔨 Building provider..."
      if make build 2>&1 | tee /tmp/build-output.txt; then
        BUILD_EXIT=0
      else
        BUILD_EXIT=$?
      fi
      echo "BUILD_EXIT=$BUILD_EXIT" >> /tmp/results.env

      echo "🧪 Running full mock-azurerm sweep (this is the slow one)..."
      if make test_mock_azurerm_all 2>&1 | tee /tmp/sweep-output.txt; then
        SWEEP_EXIT=0
      else
        SWEEP_EXIT=$?
      fi
      echo "SWEEP_EXIT=$SWEEP_EXIT" >> /tmp/results.env

      # Quick numeric summary
      if [[ -f /tmp/azurecaf-mock/report.tsv ]]; then
        TOTAL=$(($(wc -l < /tmp/azurecaf-mock/report.tsv) - 1))
        PASS=$(awk -F'\t' 'NR>1 && $2=="PASS"' /tmp/azurecaf-mock/report.tsv | wc -l)
        FAIL=$(awk -F'\t' 'NR>1 && $2=="FAIL"' /tmp/azurecaf-mock/report.tsv | wc -l)
        INIT_FAIL=$(awk -F'\t' 'NR>1 && $2=="INIT_FAIL"' /tmp/azurecaf-mock/report.tsv | wc -l)
        echo "TOTAL=$TOTAL"           >> /tmp/results.env
        echo "PASS=$PASS"             >> /tmp/results.env
        echo "FAIL=$FAIL"             >> /tmp/results.env
        echo "INIT_FAIL=$INIT_FAIL"   >> /tmp/results.env

        # Top 30 failing resources with their one-line error summary
        awk -F'\t' 'NR>1 && ($2=="FAIL" || $2=="INIT_FAIL") {print $1 "\t" $2 "\t" $5}' \
          /tmp/azurecaf-mock/report.tsv | head -30 > /tmp/top-failures.tsv
      else
        echo "REPORT_MISSING=1" >> /tmp/results.env
      fi

      # Always exit 0 so the agent can run and report.
      exit 0

  - name: Upload artifacts
    if: always()
    uses: actions/upload-artifact@v4
    with:
      name: mock-azurerm-weekly
      path: |
        /tmp/azurecaf-mock/report.tsv
        /tmp/azurecaf-mock/logs/
        /tmp/sweep-output.txt
      if-no-files-found: ignore
      retention-days: 30

source: local
engine: copilot
---

# Weekly mock-azurerm sweep

Run the full mock-azurerm validation once a week and open a categorized issue
when failures are present.

## Inputs available to you

- `/tmp/results.env` — `BUILD_EXIT`, `SWEEP_EXIT`, `TOTAL`, `PASS`, `FAIL`, `INIT_FAIL` (or `REPORT_MISSING=1`)
- `/tmp/azurecaf-mock/report.tsv` — full TSV report (columns: `resource`, `status`, `pass`, `fail`, `error_summary`)
- `/tmp/top-failures.tsv` — top 30 failing rows
- `/tmp/azurecaf-mock/logs/<resource>.log` — per-resource `terraform init` + `terraform test` output
- `/tmp/sweep-output.txt`, `/tmp/build-output.txt` — full make output

## Process

1. Read `/tmp/results.env`.
2. **If the build failed** (`BUILD_EXIT != 0`) **or the report is missing** (`REPORT_MISSING=1`): open a single critical issue titled `[mock-azurerm] Weekly sweep could not run — <YYYY-MM-DD>` with the build output tail and stop.
3. **If `FAIL == 0` and `INIT_FAIL == 0`**: do nothing (no issue, no comment). Exit quietly.
4. **Otherwise**, classify each failing row from `/tmp/top-failures.tsv` into exactly one of three buckets by reading the per-resource log when ambiguous:
   - **Real CAF bug** — the regex / length constraint in `resourceDefinition.json` produces a name that azurerm rejects with an explicit naming error (e.g. "must match", "exceeds maximum length", "invalid characters"). These need a fix in `resourceDefinition.json`.
   - **Scaffolding gap** — azurerm rejects an *attribute other than the name* (e.g. `account_replication_type` got `"test"` instead of `"LRS"`). Needs a new entry in `RESOURCE_ATTR_OVERRIDES` in `scripts/mock-test/generate_tests.py`.
   - **Deprecated upstream resource** — `terraform init` fails because the `azurerm_*` type no longer exists in the current `azurerm` provider. The CAF entry should be marked deprecated or removed.
5. Open a **single** issue (the workflow's `close-older-issues: true` will automatically close last week's). Title: `[mock-azurerm] Weekly sweep — <PASS>/<TOTAL> passed (<FAIL+INIT_FAIL> failures) — <YYYY-MM-DD>`.
6. Body must include, in this order:
   - One-line summary (`<PASS>/<TOTAL> passed; <FAIL> assertion failures; <INIT_FAIL> init failures`)
   - Three H3 sections (`### Real CAF bugs`, `### Scaffolding gaps`, `### Deprecated upstream resources`), each with a markdown table `| Resource | Error |`
   - Skip a section entirely if its bucket is empty (write `_(none this week)_` instead)
   - A short "Suggested next steps" section pointing at the right file to edit per bucket
   - A "Reproduce locally" snippet:
     ````
     make build
     make test_mock_azurerm_all
     less /tmp/azurecaf-mock/logs/<resource>.log
     ````
   - A link to the workflow run: use `${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}`

## Rules

- Open **at most one** issue per run. Use `safe-outputs.create-issue` exactly once.
- Do **not** modify any source files in this run — your only side effect is opening the issue.
- Do **not** invent failure entries. Only report rows that appear in `/tmp/top-failures.tsv` (which is already capped at 30). If `FAIL+INIT_FAIL > 30`, mention the truncation in the issue body.
- Keep the issue body under ~6 KB so notifications stay readable.
