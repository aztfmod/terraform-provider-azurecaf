---
name: pre-release-validation
description: "Run comprehensive pre-release checks: full test suite, E2E tests, coverage verification, CHANGELOG validation, and generated code freshness. Triggers on: before creating a release tag, release preparation."
---

# Pre-Release Validation

## Procedure

### 1. Verify CHANGELOG

Check that `CHANGELOG.md` has entries under `## [Unreleased]`:

```bash
sed -n '/^## \[Unreleased\]/,/^## \[v/p' CHANGELOG.md | grep -c "^- "
```

If no entries, FAIL: "No unreleased changes documented in CHANGELOG.md."

### 2. Verify generated code is fresh

```bash
go generate
git diff --name-only azurecaf/models_generated.go
```

If there are uncommitted changes, FAIL: "Generated code is stale. Run `go generate` and commit."

### 3. Run full test suite

```bash
make test_ci 2>&1 | tee /tmp/prerelease-tests.txt
```

All tests must pass.

### 4. Run E2E tests

```bash
make test_e2e 2>&1 | tee /tmp/prerelease-e2e.txt
```

All E2E tests must pass.

### 5. Check coverage

```bash
make test_coverage 2>&1 | tee /tmp/prerelease-coverage.txt
```

Coverage must be >= 95%.

### 6. Lifecycle consistency check

Verify plan→apply→plan produces no drift. Build and install the provider locally, then:

```bash
mkdir -p /tmp/prerelease-lifecycle && cat > /tmp/prerelease-lifecycle/main.tf << 'LCEOF'
terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
    }
  }
}

resource "azurecaf_name" "lifecycle_test" {
  name          = "prerelease"
  resource_type = "azurerm_resource_group"
  prefixes      = ["prod"]
  suffixes      = ["001"]
  random_length = 5
  random_seed   = 42
  clean_input   = true
}

resource "azurecaf_name" "no_seed_test" {
  name          = "prerelease"
  resource_type = "azurerm_storage_account"
  prefixes      = ["prod"]
  random_length = 5
  clean_input   = true
}

output "rg_result" { value = azurecaf_name.lifecycle_test.result }
output "st_result" { value = azurecaf_name.no_seed_test.result }
LCEOF
cd /tmp/prerelease-lifecycle
terraform plan -out=tfplan   # result should be visible for rg (has seed)
terraform apply tfplan
terraform plan               # Must show "No changes"
rm -rf /tmp/prerelease-lifecycle
```

FAIL if the final plan shows any changes (plan-apply inconsistency).

Also verify that `rg_result` shows an actual name during the first plan (not `(known after apply)`).

### 7. Report

```
## Pre-Release Validation Report

| Check | Status |
|-------|--------|
| CHANGELOG has entries | ✅ / ❌ |
| Generated code fresh | ✅ / ❌ |
| Unit tests | ✅ / ❌ |
| E2E tests | ✅ / ❌ |
| Coverage >= 95% | ✅ (<percentage>%) / ❌ (<percentage>%) |
| Lifecycle consistency | ✅ / ❌ |

Overall: READY FOR RELEASE / NOT READY
```

### 8. Cleanup

```bash
rm -f /tmp/prerelease-*.txt
```
