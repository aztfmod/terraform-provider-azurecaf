#!/usr/bin/env bash
set -euo pipefail

fail() {
  printf 'release readiness: %s\n' "$*" >&2
  exit 1
}

tag="${1:-}"
if [[ ! "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  fail "expected a stable semantic-version tag such as v2.0.0"
fi

root="${RELEASE_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
changelog="${CHANGELOG_FILE:-$root/CHANGELOG.md}"
notes="${RELEASE_NOTES_FILE:-$root/.github/releases/$tag.md}"

[[ -f "$changelog" ]] || fail "missing CHANGELOG.md at $changelog"
[[ -s "$notes" ]] || fail "missing or empty release notes at $notes"
grep -Eq '[^[:space:]]' "$notes" || fail "release notes contain only whitespace"

unreleased_count=$(grep -c '^## \[Unreleased\]$' "$changelog" || true)
[[ "$unreleased_count" -eq 1 ]] || fail "CHANGELOG.md must contain exactly one Unreleased heading"

unreleased=$(
  awk '
    /^## \[Unreleased\]$/ { in_section = 1; next }
    in_section && /^## \[/ { exit }
    in_section { print }
  ' "$changelog"
)
if grep -Eq '[^[:space:]]' <<<"$unreleased"; then
  fail "CHANGELOG.md Unreleased section must be empty before tagging"
fi

tag_pattern=${tag//./\\.}
tag_heading_count=$(grep -Ec "^## \\[$tag_pattern\\]( |$)" "$changelog" || true)
[[ "$tag_heading_count" -eq 1 ]] || fail "CHANGELOG.md must contain exactly one section for $tag"

header_pattern="^## \\[$tag_pattern\\] - [0-9]{4}-[0-9]{2}-[0-9]{2}$"
header_count=$(grep -Ec "$header_pattern" "$changelog" || true)
[[ "$header_count" -eq 1 ]] || fail "CHANGELOG.md must contain one dated section for $tag"

release_section=$(
  awk -v tag="$tag" '
    index($0, "## [" tag "] - ") == 1 { in_section = 1; next }
    in_section && /^## \[/ { exit }
    in_section { print }
  ' "$changelog"
)
if ! grep -Eq '^(### |\- )' <<<"$release_section"; then
  fail "CHANGELOG.md section for $tag has no release entries"
fi

printf 'Release metadata is ready for %s\n' "$tag"
