---
name: regression-test-runner
description: "Run the full CI test suite (make test_ci), parse output, and report pass/fail with failure details. Use after code changes to verify nothing is broken. Triggers on: code changes, PR validation, post-build verification."
---

# Regression Test Runner

## Procedure

### 1. Run test suite

```bash
make test_ci 2>&1 | tee /tmp/test-ci-output.txt
```

This runs: unit tests, coverage, resource validation, and matrix testing.

### 2. Parse results

Check exit code:
- **0**: All tests passed
- **Non-zero**: Failures detected

Extract summary from output:
```bash
grep -E "^(ok|FAIL|---)" /tmp/test-ci-output.txt
```

### 3. Report

**If all tests pass:**
```
✅ Regression tests: PASSED
   Packages: <count> passed
   Coverage: <percentage>%
```

**If failures detected:**
```
❌ Regression tests: FAILED
   Failed packages:
   - <package>: <failure summary>
   
   Failure details:
   <relevant test output>
```

Extract the specific test function names and error messages for each failure.

### 4. Coverage check

Extract coverage percentage:
```bash
grep "coverage:" /tmp/test-ci-output.txt | tail -1
```

Flag if coverage drops below 95%.

### 5. Cleanup

```bash
rm -f /tmp/test-ci-output.txt
```
