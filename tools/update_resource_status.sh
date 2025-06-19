#!/bin/bash

# Script para actualizar la tabla de Resource Status en el README.md
# Cambia ❌ por ✔ si el recurso está presente en resourceDefinition.json

# Calcula rutas absolutas robustas
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
README="$SCRIPT_DIR/../README.md"
RESDEF="$SCRIPT_DIR/../resourceDefinition.json"

# Extrae todos los nombres de recursos del JSON
RESOURCES=$(jq -r '.[].name' "$RESDEF")

awk -v resources="$RESOURCES" '
  BEGIN {
    split(resources, arr, "\n");
    for (i in arr) present[arr[i]] = 1;
    in_table=0
  }
  /^\|resource \| status \|/ { in_table=1; print; next }
  in_table && /^\|/ {
    if ($0 ~ /⚠/) { print $0; next }
    split($0, a, "|")
    res=a[2]; gsub(/^ +| +$/, "", res)
    if (res == "resource" || res == "---") { print $0; next }
    status = (present[res]) ? "✔" : "❌"
    print "|" res " | " status " |"
    next
  }
  { print }
' "$README" > "$README.tmp" && mv "$README.tmp" "$README"

echo "Tabla de Resource Status actualizada en $README"