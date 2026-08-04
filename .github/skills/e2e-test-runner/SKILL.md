---
name: e2e-test-runner
description: "Run end-to-end tests (make test_e2e or make test_e2e_quick), parse results, and produce a structured summary. Use after build succeeds to validate real Terraform workflows. Triggers on: post-build validation, PR checks, release verification."
---

# E2E Test Runner

## Procedure

### 1. Choose test scope

| Scope | Command | Duration | When to use |
|-------|---------|----------|-------------|
| Quick | `make test_e2e_quick` | ~10-15s | Every change, fast feedback |
| Full | `make test_e2e` | ~25-30s | PR validation, pre-release |
| Data source | `make test_e2e_data_source` | ~10s | Data source changes |
| Naming | `make test_e2e_naming` | ~15s | Naming convention changes |
| Import | `make test_e2e_import` | ~10s | Import functionality changes |

Default to **quick** unless the caller specifies otherwise.

### 2. Run tests

```bash
make test_e2e_quick 2>&1 | tee /tmp/e2e-output.txt
```

### 3. Parse results

Check exit code and extract test results:

```bash
grep -E "^(ok|FAIL|---)" /tmp/e2e-output.txt
```

### 4. Report

**If all tests pass:**
```
✅ E2E tests: PASSED (<scope>)
   Tests: <count> passed
   Duration: <time>
```

**If failures detected:**
```
❌ E2E tests: FAILED (<scope>)
   Failed tests:
   - <test_name>: <failure summary>
   
   Details:
   <relevant output>
```

### 5. Lifecycle consistency check

After E2E tests pass, verify plan→apply→plan consistency:

```bash
mkdir -p /tmp/e2e-lifecycle && cat > /tmp/e2e-lifecycle/main.tf << 'LCEOF'
terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
    }
  }
}

resource "azurecaf_name" "lifecycle_test" {
  name          = "lctest"
  resource_type = "azurerm_resource_group"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 4
  random_seed   = 12345
  clean_input   = true
}

output "result" { value = azurecaf_name.lifecycle_test.result }
LCEOF
cd /tmp/e2e-lifecycle
terraform plan -out=tfplan
terraform apply tfplan
terraform plan  # Must show "No changes"
rm -rf /tmp/e2e-lifecycle
```

The final `terraform plan` must report **"No changes"**. A drift here indicates plan-apply inconsistency.

Also verify that `result` shows an actual name during the first plan (not `(known after apply)`).

### 6. Cleanup

```bash
rm -f /tmp/e2e-output.txt
```
