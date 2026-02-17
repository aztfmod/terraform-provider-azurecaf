---
name: terraform-mock-test
description: "Validate a resource definition end-to-end using terraform test with mock_provider azurerm. Proves the CAF-generated name is accepted by the azurerm provider schema without Azure credentials. Use after provider-build-test succeeds to run the mocked azurerm integration test."
---

# Terraform Mock Test

## Prerequisites

- Successfully built `terraform-provider-azurecaf` binary (from `provider-build-test` skill)
- Resource entry in `resourceDefinition.json` with known `min_length`, `max_length`, `validation_regex`

## Procedure

### 1. Install provider locally

```bash
GOOS=$(go env GOOS) GOARCH=$(go env GOARCH)
LOCAL_PLUGIN_DIR=~/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/${GOOS}_${GOARCH}
mkdir -p "$LOCAL_PLUGIN_DIR"
cp ./terraform-provider-azurecaf "$LOCAL_PLUGIN_DIR/"
```

### 2. Look up azurerm resource attributes

Use **Terraform MCP** `get_provider_details` (provider_name: `azurerm`, provider_namespace: `hashicorp`, provider_document_type: `resources`, resource_name: `<resource_name>`).

Extract only **Required** arguments and **Required nested blocks**. Do NOT use Context7.

### 3. Create test directory

```bash
mkdir -p /tmp/azurecaf-test-<resource_name>/tests
```

### 4. Create test files

See [references/test-templates.md](references/test-templates.md) for the `main.tf`, `terraform.rc`, and `tests/validate_name.tftest.hcl` templates.

Key rules:
- Use `aztfmod/azurecaf` as provider source (not `aztfmod.com/arnaudlh/azurecaf`)
- Use **hardcoded fake Azure resource IDs** for parent references — do NOT chain azurerm resources
- The azurerm resource MUST set `name = azurecaf_name.test.result`

### 5. Run test

```bash
cd /tmp/azurecaf-test-<resource_name>
terraform init   # downloads azurerm provider schema
TF_CLI_CONFIG_FILE=/tmp/azurecaf-test-<resource_name>/terraform.rc terraform test -verbose
```

Always use `-verbose` to show generated name, resource state, and output values.

### 6. Validate results

The `Outputs:` section must show non-empty `result` and `result_length`. If empty, assertions are not using `output.result` — fix them.

### 7. Clean up

```bash
rm -rf /tmp/azurecaf-test-<resource_name>
```
