---
name: test-failure-diagnosis
description: "Analyze test failure output to identify root cause and suggest fixes. Use when build or test failures occur. Triggers on: test failures, build errors, CI failures."
---

# Test Failure Diagnosis

## Procedure

### 1. Classify failure type

Read the failure output and classify:

| Failure type | Indicators | Common cause |
|-------------|------------|--------------|
| **Compilation error** | `cannot compile`, `undefined:`, `syntax error` | Bad JSON in resourceDefinition.json, invalid Go syntax |
| **JSON parse error** | `json: cannot unmarshal`, `unexpected end of JSON` | Missing comma, unmatched quotes, bad escapes in resourceDefinition.json |
| **Duplicate slug** | `duplicate key`, ambiguity errors | Two resources share the same slug |
| **Test assertion failure** | `FAIL`, `expected X got Y` | Naming logic bug, regex mismatch |
| **Regex error** | `error parsing regexp` | Invalid regex pattern in validation_regex or regex field |
| **Import error** | `cannot find package`, `module not found` | Missing dependency, go.mod issue |

### 2. Extract failure details

From the output, extract:
- **Failed test name(s)**: The specific `Test*` function(s) that failed
- **Error message**: The assertion or error text
- **File and line**: Where the failure occurred
- **Expected vs actual**: What was expected and what was produced

### 3. Identify root cause

Based on the classification, check the most likely source:

**For JSON errors:**
```bash
python3 -c "import json; json.load(open('resourceDefinition.json'))"
```

**For duplicate slugs:**
```bash
grep '"slug":' resourceDefinition.json | sort | uniq -d
```

**For regex issues:**
```bash
grep -n 'validation_regex\|"regex"' resourceDefinition.json | grep -v '"\"'
```

### 4. Suggest fix

Provide a specific, actionable fix:
- For JSON formatting: show the exact line and the correction
- For duplicate slugs: identify which resources conflict
- For regex: show the invalid pattern and the corrected version
- For test failures: show what the test expects and how to satisfy it

### 5. Verify fix

After applying the suggested fix, re-run the failed test:

```bash
go test ./azurecaf/ -run "<failed_test_name>" -v
```
