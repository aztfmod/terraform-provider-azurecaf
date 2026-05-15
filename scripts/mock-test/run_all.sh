#!/usr/bin/env bash
# Run `terraform test` for every generated workspace under OUT_DIR and write
# an aggregate TSV report to REPORT_FILE. Exits non-zero if any workspace
# fails its assertions or fails to initialize.
#
# Environment / args:
#   OUT_DIR     (required)  Directory containing one subdirectory per resource.
#   REPORT_FILE (required)  Path to the aggregate TSV report to write.
#   LOG_DIR     (optional)  Where per-resource logs are written. Defaults to
#                           "$OUT_DIR/../logs".
set -uo pipefail

usage() {
  cat <<'EOF'
Usage: run_all.sh --out-dir <dir> --report <file> [--log-dir <dir>]
EOF
}

OUT_DIR=""
REPORT_FILE=""
LOG_DIR=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --out-dir)  OUT_DIR="$2"; shift 2 ;;
    --report)   REPORT_FILE="$2"; shift 2 ;;
    --log-dir)  LOG_DIR="$2"; shift 2 ;;
    -h|--help)  usage; exit 0 ;;
    *)          echo "unknown arg: $1" >&2; usage; exit 2 ;;
  esac
done

if [[ -z "$OUT_DIR" || -z "$REPORT_FILE" ]]; then
  usage; exit 2
fi
if [[ ! -d "$OUT_DIR" ]]; then
  echo "out-dir does not exist: $OUT_DIR" >&2
  exit 2
fi

LOG_DIR="${LOG_DIR:-$(dirname "$OUT_DIR")/logs}"
mkdir -p "$LOG_DIR"

# Use a shared plugin cache so azurerm is downloaded once for the whole sweep.
export TF_PLUGIN_CACHE_DIR="${TF_PLUGIN_CACHE_DIR:-${HOME}/.terraform.d/plugin-cache}"
mkdir -p "$TF_PLUGIN_CACHE_DIR"
export TF_IN_AUTOMATION=1
export CHECKPOINT_DISABLE=1

printf 'resource\tstatus\tpass\tfail\terror_summary\n' > "$REPORT_FILE"

total=0
pass=0
fail=0
init_fail=0

for d in "$OUT_DIR"/*/ ; do
  [[ -d "$d" ]] || continue
  rt="$(basename "$d")"
  total=$((total+1))
  log="$LOG_DIR/${rt}.log"

  pushd "$d" > /dev/null

  if ! terraform init -no-color -input=false > "$log" 2>&1 ; then
    printf '%s\tINIT_FAIL\t0\t0\tinit failed\n' "$rt" >> "$REPORT_FILE"
    init_fail=$((init_fail+1))
    popd > /dev/null
    continue
  fi

  if TF_CLI_CONFIG_FILE="$PWD/terraform.rc" terraform test -no-color >> "$log" 2>&1 ; then
    line=$(grep -E "Success!|Failure!" "$log" | tail -1)
    p=$(echo "$line" | sed -nE 's/.*Success! ([0-9]+) passed.*/\1/p')
    [[ -z "$p" ]] && p=0
    printf '%s\tPASS\t%s\t0\t\n' "$rt" "$p" >> "$REPORT_FILE"
    pass=$((pass+1))
  else
    line=$(grep -E "Success!|Failure!" "$log" | tail -1)
    p=$(echo "$line" | sed -nE 's/.*Failure! ([0-9]+) passed.*/\1/p')
    f=$(echo "$line" | sed -nE 's/.*Failure! [0-9]+ passed, ([0-9]+) failed.*/\1/p')
    [[ -z "$p" ]] && p=0
    [[ -z "$f" ]] && f=1
    summary=$(grep -E "^Error:" "$log" | head -1 | tr -d '\n' | head -c 160)
    [[ -z "$summary" ]] && summary=$(grep -m1 -E "error_message" "$log" | head -c 160)
    printf '%s\tFAIL\t%s\t%s\t%s\n' "$rt" "$p" "$f" "$summary" >> "$REPORT_FILE"
    fail=$((fail+1))
  fi
  popd > /dev/null
done

echo "==============================="
echo "Total:     $total"
echo "Pass:      $pass"
echo "Fail:      $fail"
echo "InitFail:  $init_fail"
echo "Report:    $REPORT_FILE"
echo "Logs:      $LOG_DIR"

# Exit non-zero if any failure, so CI reports a red check.
if (( fail > 0 || init_fail > 0 )); then
  exit 1
fi
exit 0
