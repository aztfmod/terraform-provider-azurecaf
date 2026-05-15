# mock-azurerm test harness

This directory contains the harness that validates every CAF-generated name
against the corresponding `azurerm_*` resource schema using `terraform test`
with `mock_provider "azurerm" {}`. No Azure credentials are required.

The harness complements the existing test layers:

| Layer | Validates name against | Where |
|---|---|---|
| `azurecaf/resource_matrix_test.go` etc. | The CAF `validation_regex` (self-check) | `make test_resource_matrix` |
| `e2e/` | A handful of hand-coded scenarios | `make test_e2e` |
| **mock-azurerm (this)** | **The `azurerm` provider's schema validators** | `make test_mock_azurerm_*` |

## Files

- `generate_tests.py` — emits one `terraform test` workspace per resource
  with three naming variations (`default`, `with_prefix`, `with_random`).
- `run_all.sh` — runs every workspace, writes a TSV report, exits non-zero on
  any failure.
- `fetch_schema.sh` — downloads the `azurerm` provider schema as JSON.

## Local usage

```bash
# 1. Build the provider so the harness can use it via dev_overrides.
make build

# 2. Install the locally built binary into the dev_overrides path.
GOOS=$(go env GOOS); GOARCH=$(go env GOARCH)
LOCAL_PLUGIN_DIR=~/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/${GOOS}_${GOARCH}
mkdir -p "$LOCAL_PLUGIN_DIR" && cp ./terraform-provider-azurecaf "$LOCAL_PLUGIN_DIR/"

# 3. Fetch the azurerm schema once (cached on disk).
scripts/mock-test/fetch_schema.sh /tmp/azurerm-schema.json

# 4. Run the diff-scoped sweep (only resources changed vs origin/main).
make test_mock_azurerm_changed

# Or run the full sweep (slow — see CI weekly job).
make test_mock_azurerm_all
```

The TSV report and per-resource logs land under `/tmp/azurecaf-mock/` by
default. Inspect a failure with:

```bash
less /tmp/azurecaf-mock/logs/azurerm_storage_encryption_scope.log
```

## Adding fake values for a new resource

If a new resource's required attributes trip the azurerm provider's
`CustomizeDiff` validators, add an entry to `RESOURCE_ATTR_OVERRIDES` in
`generate_tests.py`. Keep entries minimal — only the attributes that the
generic `fake_value_for` heuristic cannot infer correctly.
