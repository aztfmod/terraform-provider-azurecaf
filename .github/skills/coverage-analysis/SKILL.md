---
name: coverage-analysis
description: "Run test coverage analysis (make test_coverage), parse coverage percentage, compare against the 95% threshold, and flag regressions. Triggers on: PR checks, post-build validation, release prep."
---

# Coverage Analysis

## Procedure

### 1. Run coverage

```bash
make test_coverage 2>&1 | tee /tmp/coverage-output.txt
```

### 2. Extract coverage percentage

```bash
grep "coverage:" /tmp/coverage-output.txt | tail -1
```

Parse the percentage value (e.g., `99.3%`).

### 3. Check threshold

The minimum coverage threshold is **95%**.

### 4. Report

**If above threshold:**
```
✅ Coverage: <percentage>% (threshold: 95%)
```

**If below threshold:**
```
⚠️ Coverage: <percentage>% — BELOW THRESHOLD (95%)
   Consider adding tests for uncovered code paths.
```

### 5. Detailed breakdown (optional)

If requested, generate HTML coverage report:

```bash
make test_coverage_html
```

This creates a coverage report that can be opened in a browser for line-by-line analysis.

### 6. Cleanup

```bash
rm -f /tmp/coverage-output.txt
```
