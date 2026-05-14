---
name: changelog-update
description: "Update CHANGELOG.md with a new entry under the correct section. Parses existing structure, adds entries under [Unreleased], and assesses semver impact. Triggers on: after resource changes, bug fixes, documentation updates, or any notable project change."
---

# Changelog Update

## Procedure

### 1. Read current CHANGELOG

```bash
head -80 CHANGELOG.md
```

Identify the `## [Unreleased]` section and its subsections (`### Added`, `### Changed`, `### Fixed`, `### Removed`, `### Deprecated`).

### 2. Determine section

| Change type | Section | Semver impact |
|-------------|---------|---------------|
| New resource type added | `### Added` | Minor |
| New feature or capability | `### Added` | Minor |
| Updated resource constraints (slug, regex, length) | `### Changed` | Patch |
| Bug fix | `### Fixed` | Patch |
| Removed resource or feature | `### Removed` | Major |
| Deprecated resource or feature | `### Deprecated` | Minor |

### 3. Format entry

Use these templates based on change type:

**New resource:**
```markdown
- **RESOURCE**: Added `<resource_name>` (<official display name>) -- slug: `<slug>`, length: <min>-<max>, scope: <scope>, dashes <allowed|not allowed>, <lowercase|mixed case>
  - Impact: <Low|Medium|High> -- <brief justification>
```

**Updated resource:**
```markdown
- **RESOURCE**: Updated `<resource_name>` -- <describe field changes>
  - Impact: <Low|Medium|High> -- <brief justification>
```

**Bug fix:**
```markdown
- **<Component>**: <Brief description of fix>
  - <Detail of what was fixed and why>
  - Impact: <Low|Medium|High> -- <brief justification>
```

**Multiple resources (batch):**
```markdown
- **RESOURCE**: Added support for <category> resource types
  - Added `<resource_1>` with slug `<slug_1>`
  - Added `<resource_2>` with slug `<slug_2>`
  - <Summary of shared characteristics>
  - Impact: <Low|Medium|High> -- <brief justification>
```

### 4. Insert entry

Insert the new entry under the appropriate subsection within `## [Unreleased]`.

- If the subsection (e.g., `### Added`) does not exist, create it under `## [Unreleased]`.
- Maintain the conventional order: Added, Changed, Deprecated, Removed, Fixed.
- Add a blank line between entries for readability.

### 5. Assess overall impact

After inserting, summarize the semver impact:

- **Patch** (x.y.Z): Only fixes or constraint updates, no new features
- **Minor** (x.Y.0): New resources or features added, no breaking changes
- **Major** (X.0.0): Removed resources, changed slugs, or other breaking changes
