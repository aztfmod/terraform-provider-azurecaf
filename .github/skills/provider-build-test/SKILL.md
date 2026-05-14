---
name: provider-build-test
description: "Regenerate Go code from resourceDefinition.json, build the terraform-provider-azurecaf binary, and run unit tests. Use after any change to resourceDefinition.json to verify the provider compiles and tests pass."
---

# Provider Build & Test

Run from the project root after editing `resourceDefinition.json`:

```bash
go generate          # regenerates azurecaf/models_generated.go
make build           # runs go generate + go fmt + go build + go test
```

Verify the resource appears in generated code:

```bash
grep "<resource_name>" azurecaf/models_generated.go
```

Expected output: `ok` for all test packages, coverage ~90%+, no compilation errors.

If `go generate` fails: check JSON formatting (missing commas, unmatched quotes, bad escapes).
If tests fail: check for duplicate slugs causing ambiguity.
