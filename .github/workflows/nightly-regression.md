---
description: |
  Nightly regression test workflow. Runs the full test suite and E2E tests
  on the main branch. Creates a GitHub issue if any tests fail, with
  detailed failure information for rapid diagnosis.

on:
  schedule: daily
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
    - "grep *"
    - "cat *"
    - "head *"
    - "tail *"

safe-outputs:
  create-issue:
    title-prefix: "[nightly-regression] "
    labels: [bug, automated, nightly]
    close-older-issues: true

steps:
  - name: Build and test
    run: |
      echo "🔨 Building provider..."
      make build 2>&1 | tee /tmp/build-output.txt
      BUILD_EXIT=$?

      echo "🧪 Running CI tests..."
      make test_ci 2>&1 | tee /tmp/test-ci-output.txt
      TEST_EXIT=$?

      echo "🔬 Running E2E quick tests..."
      make test_e2e_quick 2>&1 | tee /tmp/e2e-output.txt
      E2E_EXIT=$?

      echo "BUILD_EXIT=$BUILD_EXIT" >> /tmp/results.env
      echo "TEST_EXIT=$TEST_EXIT" >> /tmp/results.env
      echo "E2E_EXIT=$E2E_EXIT" >> /tmp/results.env

      # Extract test summary lines
      grep -E "^(ok|FAIL|---)" /tmp/test-ci-output.txt > /tmp/test-summary.txt 2>/dev/null || true
      grep -E "^(ok|FAIL|---)" /tmp/e2e-output.txt >> /tmp/test-summary.txt 2>/dev/null || true

source: local
engine: copilot
---

# Nightly Regression

Run nightly regression tests on the main branch and report failures.

## Process

1. Read the test results from `/tmp/results.env` and `/tmp/test-summary.txt`
2. If ALL exit codes are 0: no action needed, exit quietly
3. If ANY exit code is non-zero:
   - Analyze the failure output files (`/tmp/build-output.txt`, `/tmp/test-ci-output.txt`, `/tmp/e2e-output.txt`)
   - Identify which tests failed and extract error details
   - Create a GitHub issue with:
     - Title: `[nightly-regression] Test failures on main — <date>`
     - Failed test names and error messages
     - Suggested investigation steps
     - Links to relevant test files
