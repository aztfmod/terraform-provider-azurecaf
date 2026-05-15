#!/usr/bin/env bash
# Fetch the azurerm provider schema as JSON for the mock-test generator.
# Writes "$OUT" (default: ./azurerm-schema.json) using the version of
# hashicorp/azurerm that satisfies ">= 4.0.0". Idempotent: skips fetch if
# the file already exists and --force is not set.
set -euo pipefail

OUT="${1:-azurerm-schema.json}"
FORCE=""
for arg in "$@"; do
  case "$arg" in
    --force) FORCE=1 ;;
  esac
done

if [[ -f "$OUT" && -z "$FORCE" ]]; then
  echo "schema already present at $OUT (use --force to refetch)"
  exit 0
fi

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

cat > "$WORK/main.tf" <<'EOF'
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0.0"
    }
  }
}
EOF

(
  cd "$WORK"
  terraform init -no-color -input=false > /dev/null
  terraform providers schema -json > schema.json
)

mv "$WORK/schema.json" "$OUT"
echo "wrote $OUT"
