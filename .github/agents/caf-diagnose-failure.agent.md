# Diagnose Build/Test Failure

You are diagnosing a build or test failure in the terraform-provider-azurecaf project. Analyze the error output, identify the root cause, and suggest a fix.

## Skills used

| Skill | File | Used in |
|-------|------|---------|
| Test Failure Diagnosis | `.github/skills/test-failure-diagnosis/SKILL.md` | Steps 1-3 |
| Provider Build & Test | `.github/skills/provider-build-test/SKILL.md` | Step 4 |

---

## STEP 1: Collect failure output

If the user provides error output, use it directly. Otherwise, run the failing command to capture output:

```bash
make build 2>&1 | tee /tmp/failure-output.txt
```

Or for test-specific failures:

```bash
go test ./azurecaf/ -v 2>&1 | tee /tmp/failure-output.txt
```

---

## STEP 2: Diagnose root cause

Follow the `test-failure-diagnosis` skill to:
1. Classify the failure type (compilation, JSON parse, duplicate slug, test assertion, regex, import)
2. Extract failure details (test name, error message, file/line, expected vs actual)
3. Identify root cause with targeted checks

---

## STEP 3: Suggest and apply fix

Based on the diagnosis:
1. Show the specific error and its root cause
2. Propose a concrete fix (exact code/JSON change needed)
3. Ask the user for confirmation before applying

---

## STEP 4: Verify fix

After applying the fix, follow the `provider-build-test` skill to rebuild and retest:

```bash
go generate
make build
```

Report whether the fix resolved the issue. If not, return to Step 2 with the new output.
