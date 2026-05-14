---
description: |
  Release validation workflow that runs comprehensive checks when a
  version tag is pushed. Posts validation results as a comment on the
  associated GitHub Release.

on:
  push:
    tags: ["v*"]

permissions:
  contents: read
  pull-requests: read

network: defaults

tools:
  github:
    lockdown: true
    toolsets: [pull_requests]
  bash:
    - "cat *"
    - "grep *"
    - "head *"
    - "tail *"

safe-outputs:
  add-comment: {}

steps:
  - name: Set up Go
    uses: actions/setup-go@v6
    with:
      go-version-file: './go.mod'
      cache: true

  - name: Install tfproviderlint
    run: |
      go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest
      echo "$(go env GOPATH)/bin" >> "$GITHUB_PATH"

  - name: Setup Terraform
    uses: hashicorp/setup-terraform@v4
    with:
      terraform_version: "~> 1.0"
      terraform_wrapper: false

  - name: Build, test, and validate
    env:
      CHECKPOINT_DISABLE: "1"
      TF_IN_AUTOMATION: "1"
      TF_CLI_ARGS_init: "-upgrade=false"
    run: |
      set -o pipefail
      : > /tmp/results.env

      TAG_NAME="${GITHUB_REF_NAME:-${GITHUB_REF##*/}}"
      echo "TAG_NAME=${TAG_NAME}" >> /tmp/results.env

      echo "🔨 Building provider..."
      if make build 2>&1 | tee /tmp/build-output.txt; then
        BUILD_EXIT=0
      else
        BUILD_EXIT=$?
      fi
      echo "BUILD_EXIT=$BUILD_EXIT" >> /tmp/results.env

      echo "🧪 Running CI tests..."
      if make test_ci 2>&1 | tee /tmp/test-ci-output.txt; then
        TEST_EXIT=0
      else
        TEST_EXIT=$?
      fi
      echo "TEST_EXIT=$TEST_EXIT" >> /tmp/results.env

      echo "🔬 Running E2E tests..."
      if make test_e2e 2>&1 | tee /tmp/e2e-output.txt; then
        E2E_EXIT=0
      else
        E2E_EXIT=$?
      fi
      echo "E2E_EXIT=$E2E_EXIT" >> /tmp/results.env

      echo "📊 Running coverage..."
      if make test_coverage 2>&1 | tee /tmp/coverage-output.txt; then
        COV_EXIT=0
      else
        COV_EXIT=$?
      fi
      echo "COV_EXIT=$COV_EXIT" >> /tmp/results.env

      # Extract overall coverage percentage from the "total:" line (go tool cover output).
      COV_PCT=$(grep -E "^total:" /tmp/coverage-output.txt | awk '{print $NF}' | tail -n 1)
      if [ -z "$COV_PCT" ]; then
        COV_PCT=$(grep -Eo "coverage: [0-9.]+% of statements" /tmp/coverage-output.txt | awk '{print $2}' | tail -n 1)
      fi
      echo "COV_PCT=${COV_PCT:-unknown}" >> /tmp/results.env

      # Check CHANGELOG for the tag version (strip leading "v" for the version-only match too).
      VERSION_NO_V="${TAG_NAME#v}"
      if grep -qE "^## \[?${VERSION_NO_V}\]?" CHANGELOG.md 2>/dev/null \
         || grep -qE "^## \[?${TAG_NAME}\]?" CHANGELOG.md 2>/dev/null; then
        CHANGELOG_EXIT=0
        grep -nE "^## \[?(${VERSION_NO_V}|${TAG_NAME})\]?" CHANGELOG.md > /tmp/changelog-entry.txt 2>/dev/null || true
      else
        CHANGELOG_EXIT=1
        echo "No section matching ${TAG_NAME} or ${VERSION_NO_V} found in CHANGELOG.md" > /tmp/changelog-entry.txt
      fi
      echo "CHANGELOG_EXIT=$CHANGELOG_EXIT" >> /tmp/results.env

      # Extract test summary lines for the agent to consume.
      {
        grep -E "^(ok|FAIL|---)" /tmp/test-ci-output.txt 2>/dev/null || true
        grep -E "^(ok|FAIL|---)" /tmp/e2e-output.txt 2>/dev/null || true
      } > /tmp/test-summary.txt

      # Always exit 0 so the agent runs and reports the results.
      exit 0

source: local
engine: copilot
---

# Release Validation

When a version tag is pushed, the pre-agent host steps build the provider,
run the full unit + E2E test suites, compute coverage, and verify the tag is
present in `CHANGELOG.md`. This agent reads the precomputed results and posts
a single summary comment on the release.

## Inputs (precomputed by host steps)

- `/tmp/results.env` — `TAG_NAME`, `BUILD_EXIT`, `TEST_EXIT`, `E2E_EXIT`,
  `COV_EXIT`, `COV_PCT`, `CHANGELOG_EXIT`.
- `/tmp/build-output.txt` — full `make build` log.
- `/tmp/test-ci-output.txt` — full `make test_ci` log.
- `/tmp/e2e-output.txt` — full `make test_e2e` log.
- `/tmp/coverage-output.txt` — full `make test_coverage` log.
- `/tmp/test-summary.txt` — `ok` / `FAIL` / `---` lines from test runs.
- `/tmp/changelog-entry.txt` — matched `CHANGELOG.md` heading(s) or a
  "not found" message.

## Process

1. `cat /tmp/results.env` to load the exit codes, tag name, and coverage
   percentage. Use these — do NOT try to run `make`, `go`, or `bash` commands;
   they are unavailable in this agent.
2. For each non-zero exit code, `head -200 /tmp/<step>-output.txt` to extract
   error context.
3. Build the validation summary using the table template below.
4. Post the summary as a single comment via the `add-comment` safe output.

## Report template

```markdown
### 🏷️ Release Validation: <TAG_NAME>

| Check | Status |
|-------|--------|
| Build | ✅ / ❌ |
| Unit tests | ✅ / ❌ |
| E2E tests | ✅ / ❌ |
| Coverage | ✅ <COV_PCT> / ❌ <COV_PCT> |
| CHANGELOG entry | ✅ / ❌ |

Overall: **VALIDATED** / **ISSUES FOUND**
```

If any check failed, append a short "Failure details" section under the table
with the first relevant error lines from the corresponding output file.
