---
name: release-notes-generator
description: "Generate release notes from CHANGELOG.md entries since the last release. Formats for GitHub Release publication. Triggers on: release preparation, tag creation."
---

# Release Notes Generator

## Procedure

### 1. Extract unreleased entries

Read `CHANGELOG.md` and extract everything under `## [Unreleased]` up to the next `## [vX.Y.Z]` header.

### 2. Format for GitHub Release

Structure the release notes:

```markdown
## What's Changed

### ✨ New Features
<entries from ### Added>

### 🔧 Changes
<entries from ### Changed>

### 🐛 Bug Fixes
<entries from ### Fixed>

### ⚠️ Deprecated
<entries from ### Deprecated>

### 💥 Breaking Changes
<entries from ### Removed>

---

**Full Changelog**: <link to compare between last tag and new tag>
```

Omit empty sections.

### 3. Generate CHANGELOG update

Prepare the CHANGELOG update — replace `## [Unreleased]` content with a new versioned section:

```markdown
## [Unreleased]

## [vX.Y.Z] - YYYY-MM-DD

### Added
<moved entries>

### Fixed
<moved entries>
```

### 4. Output

Write the formatted release notes to `.github/releases/vX.Y.Z.md` so the tag
workflow passes them to GoReleaser with `--release-notes`.

After updating both files, validate the finalized metadata:

```bash
./scripts/check-release-ready.sh vX.Y.Z
```

Return:
- The formatted release notes path and content
- The updated CHANGELOG.md content (for committing)
