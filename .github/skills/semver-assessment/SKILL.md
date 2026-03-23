---
name: semver-assessment
description: "Analyze changes since the last release tag and determine the appropriate semantic version bump (patch/minor/major) based on CHANGELOG entries and commit types. Triggers on: release preparation, version planning."
---

# Semver Assessment

## Procedure

### 1. Find last release tag

```bash
git describe --tags --abbrev=0
```

### 2. Analyze CHANGELOG entries

Read `CHANGELOG.md` and examine the `## [Unreleased]` section. Classify entries:

| Section | Semver impact |
|---------|--------------|
| `### Added` | Minor (new features) |
| `### Changed` | Patch (unless breaking → Major) |
| `### Deprecated` | Minor |
| `### Removed` | Major (breaking) |
| `### Fixed` | Patch |

### 3. Analyze commits

```bash
git log $(git describe --tags --abbrev=0)..HEAD --oneline
```

Look for:
- `BREAKING CHANGE` or `!` in commit messages → Major
- `feat:` → Minor
- `fix:` → Patch

### 4. Determine version

Apply the highest-impact rule:
- Any **Major** trigger → Major bump
- Any **Minor** trigger (no Major) → Minor bump
- Only **Patch** triggers → Patch bump

### 5. Report

```
## Version Assessment

Current version: <current_tag>
Recommended bump: <MAJOR|MINOR|PATCH>
Suggested next version: <x.y.z>

### Rationale
- <key change 1> → <impact>
- <key change 2> → <impact>
```
