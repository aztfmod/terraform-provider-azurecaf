---
description: |
  Release validation workflow that runs comprehensive checks when a
  version tag is pushed. Writes a workflow summary and opens an issue
  when any validation check fails.

on:
  push:
    tags: ["v*"]

permissions:
  contents: read
  copilot-requests: write

network: defaults

tools:
  bash:
    - "cat *"
    - "grep *"
    - "head *"
    - "tail *"

safe-outputs:
  create-issue:
    title-prefix: "[Release Validation] "
    max: 1
    close-older-issues: true

steps:
  - name: Set up Go
    uses: actions/setup-go@v7
    with:
      go-version-file: './go.mod'
      cache: true

  - name: Install tfproviderlint
    run: |
      go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@v0.31.0
      echo "$(go env GOPATH)/bin" >> "$GITHUB_PATH"

  - name: Setup Terraform
    uses: hashicorp/setup-terraform@v4
    with:
      terraform_version: "~> 1.15"
      terraform_wrapper: false

  - name: Build, test, and validate
    env:
      CHECKPOINT_DISABLE: "1"
      TF_IN_AUTOMATION: "1"
      TF_CLI_ARGS_init: "-upgrade=false"
    run: |
      set -o pipefail
      RESULTS_DIR="/tmp/gh-aw/agent/release-validation"
      mkdir -p "$RESULTS_DIR"
      : > "$RESULTS_DIR/results.env"

      TAG_NAME="${GITHUB_REF_NAME:-${GITHUB_REF##*/}}"
      echo "TAG_NAME=${TAG_NAME}" >> "$RESULTS_DIR/results.env"

      echo "🔨 Building provider..."
      if make build 2>&1 | tee "$RESULTS_DIR/build-output.txt"; then
        BUILD_EXIT=0
      else
        BUILD_EXIT=$?
      fi
      echo "BUILD_EXIT=$BUILD_EXIT" >> "$RESULTS_DIR/results.env"

      echo "🧪 Running CI tests..."
      if make test_ci 2>&1 | tee "$RESULTS_DIR/test-ci-output.txt"; then
        TEST_EXIT=0
      else
        TEST_EXIT=$?
      fi
      echo "TEST_EXIT=$TEST_EXIT" >> "$RESULTS_DIR/results.env"

      echo "🔬 Running E2E tests..."
      if make test_e2e 2>&1 | tee "$RESULTS_DIR/e2e-output.txt"; then
        E2E_EXIT=0
      else
        E2E_EXIT=$?
      fi
      echo "E2E_EXIT=$E2E_EXIT" >> "$RESULTS_DIR/results.env"

      echo "📊 Running coverage..."
      if make test_coverage 2>&1 | tee "$RESULTS_DIR/coverage-output.txt"; then
        COV_EXIT=0
      else
        COV_EXIT=$?
      fi

      COV_PCT=$(grep -Eo "Provider coverage: [0-9.]+%" "$RESULTS_DIR/coverage-output.txt" \
        | awk '{print $3}' | tail -n 1 || true)
      if [ -z "$COV_PCT" ]; then
        COV_PCT=unknown
        COV_EXIT=1
      fi
      echo "COV_EXIT=$COV_EXIT" >> "$RESULTS_DIR/results.env"
      echo "COV_PCT=$COV_PCT" >> "$RESULTS_DIR/results.env"

      # Check CHANGELOG for the tag version (strip leading "v" for the version-only match too).
      VERSION_NO_V="${TAG_NAME#v}"
      if grep -qE "^## \[?${VERSION_NO_V}\]?" CHANGELOG.md 2>/dev/null \
         || grep -qE "^## \[?${TAG_NAME}\]?" CHANGELOG.md 2>/dev/null; then
        CHANGELOG_EXIT=0
        grep -nE "^## \[?(${VERSION_NO_V}|${TAG_NAME})\]?" CHANGELOG.md \
          > "$RESULTS_DIR/changelog-entry.txt" 2>/dev/null || true
      else
        CHANGELOG_EXIT=1
        echo "No section matching ${TAG_NAME} or ${VERSION_NO_V} found in CHANGELOG.md" \
          > "$RESULTS_DIR/changelog-entry.txt"
      fi
      echo "CHANGELOG_EXIT=$CHANGELOG_EXIT" >> "$RESULTS_DIR/results.env"

      # Extract test summary lines for the agent to consume.
      {
        grep -E "^(ok|FAIL|---)" "$RESULTS_DIR/test-ci-output.txt" 2>/dev/null || true
        grep -E "^(ok|FAIL|---)" "$RESULTS_DIR/e2e-output.txt" 2>/dev/null || true
      } > "$RESULTS_DIR/test-summary.txt"

      VALIDATION_EXIT=0
      for result in "$BUILD_EXIT" "$TEST_EXIT" "$E2E_EXIT" "$COV_EXIT" "$CHANGELOG_EXIT"; do
        if [ "$result" -ne 0 ]; then
          VALIDATION_EXIT=1
        fi
      done
      echo "VALIDATION_EXIT=$VALIDATION_EXIT" >> "$RESULTS_DIR/results.env"

      status_icon() {
        if [ "$1" -eq 0 ]; then
          printf '✅'
        else
          printf '❌'
        fi
      }

      {
        echo "### 🏷️ Release Validation: $TAG_NAME"
        echo
        echo "| Check | Status |"
        echo "|-------|--------|"
        echo "| Build | $(status_icon "$BUILD_EXIT") |"
        echo "| CI test suite | $(status_icon "$TEST_EXIT") |"
        echo "| E2E tests | $(status_icon "$E2E_EXIT") |"
        echo "| Coverage >= 95% | $(status_icon "$COV_EXIT") $COV_PCT |"
        echo "| CHANGELOG entry | $(status_icon "$CHANGELOG_EXIT") |"
        echo
        if [ "$VALIDATION_EXIT" -eq 0 ]; then
          echo "Overall: **VALIDATED**"
        else
          echo "Overall: **ISSUES FOUND**"
        fi
      } >> "$GITHUB_STEP_SUMMARY"

      # Always exit 0 so the agent runs and reports the results.
      exit 0

source: local
engine:
  id: copilot
  model: gpt-5.6-luna
---

# Release Validation

When a version tag is pushed, the pre-agent host steps build the provider,
run the full CI + E2E test suites, compute coverage, and verify the tag is
present in `CHANGELOG.md`. The host step writes the complete status table to
the GitHub Actions job summary. This agent reads the precomputed results and
opens one actionable issue only when validation fails.

## Inputs (precomputed by host steps)

- `/tmp/gh-aw/agent/release-validation/results.env` — `TAG_NAME`,
  `BUILD_EXIT`, `TEST_EXIT`, `E2E_EXIT`, `COV_EXIT`, `COV_PCT`,
  `CHANGELOG_EXIT`, and `VALIDATION_EXIT`.
- `/tmp/gh-aw/agent/release-validation/build-output.txt` — full `make build` log.
- `/tmp/gh-aw/agent/release-validation/test-ci-output.txt` — full `make test_ci` log.
- `/tmp/gh-aw/agent/release-validation/e2e-output.txt` — full `make test_e2e` log.
- `/tmp/gh-aw/agent/release-validation/coverage-output.txt` — full `make test_coverage` log.
- `/tmp/gh-aw/agent/release-validation/test-summary.txt` — `ok` / `FAIL` /
  `---` lines from test runs.
- `/tmp/gh-aw/agent/release-validation/changelog-entry.txt` — matched
  `CHANGELOG.md` heading(s) or a "not found" message.

## Process

1. `cat /tmp/gh-aw/agent/release-validation/results.env` to load the exit
   codes, tag name, and coverage percentage. Use these — do NOT try to run
   `make`, `go`, or unrestricted shell commands; they are unavailable here.
2. If `VALIDATION_EXIT=0`, call `noop` with a concise success message. The
   detailed table is already in the GitHub Actions job summary.
3. If `VALIDATION_EXIT` is non-zero, read the first 200 lines from each failed
   step's output file and create one issue titled `<TAG_NAME> validation
   failed`. Include the table below plus concise failure details and the
   workflow run URL.

## Failure issue template

```markdown
### 🏷️ Release Validation: <TAG_NAME>

| Check | Status |
|-------|--------|
| Build | ✅ / ❌ |
| CI test suite | ✅ / ❌ |
| E2E tests | ✅ / ❌ |
| Coverage | ✅ <COV_PCT> / ❌ <COV_PCT> |
| CHANGELOG entry | ✅ / ❌ |

Overall: **VALIDATED** / **ISSUES FOUND**
```

Append a short "Failure details" section with only the first relevant error
lines from each failed check.
