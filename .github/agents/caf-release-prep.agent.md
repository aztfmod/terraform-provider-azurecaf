# Prepare Release

You are preparing a new release of the terraform-provider-azurecaf project. This agent validates everything is ready, determines the version, and generates release notes.

## Skills used

| Skill | File | Used in |
|-------|------|---------|
| Pre-Release Validation | `.github/skills/pre-release-validation/SKILL.md` | Step 1 |
| Semver Assessment | `.github/skills/semver-assessment/SKILL.md` | Step 2 |
| Release Notes Generator | `.github/skills/release-notes-generator/SKILL.md` | Step 3 |
| Changelog Update | `.github/skills/changelog-update/SKILL.md` | Step 4 |

---

## STEP 1: Pre-release validation

Follow the `pre-release-validation` skill to run all checks:
- CHANGELOG has unreleased entries
- Generated code is fresh
- Full test suite passes
- E2E tests pass
- Coverage >= 95%

If any check fails, STOP and report what needs to be fixed before release.

---

## STEP 2: Determine version

Follow the `semver-assessment` skill to analyze changes and recommend the version bump.

Present the recommendation to the user and ask for confirmation.

---

## STEP 3: Generate release notes

Follow the `release-notes-generator` skill to:
1. Extract unreleased CHANGELOG entries
2. Format them as GitHub Release notes
3. Save them as `.github/releases/vX.Y.Z.md`
4. Prepare the CHANGELOG update with the new version section

Show the release notes to the user for review.

---

## STEP 4: Finalize CHANGELOG

If the user confirms:
1. Update `CHANGELOG.md` — move entries from `## [Unreleased]` to `## [vX.Y.Z] - YYYY-MM-DD`
2. Add an empty `## [Unreleased]` section at the top
3. Commit `.github/releases/vX.Y.Z.md` with the CHANGELOG update
4. Run `./scripts/check-release-ready.sh vX.Y.Z`

If the release readiness script fails, STOP and fix the metadata before
requesting final review.

---

## STEP 5: Verify merge readiness

Follow `.github/RELEASE_RUNBOOK.md` and verify the final pull request SHA,
required checks, review decision, merge state, and branch freshness. Do not
lower branch protection. Record explicit maintainer authorization for any
admin bypass.

---

## Summary

```
✅ Release Preparation Complete

Version: vX.Y.Z
Changes: <count> entries
Tests: ALL PASSED
Coverage: <percentage>%

Next steps:
1. Merge the reviewed release pull request
2. Wait for the main-branch workflows to pass
3. Create tag: git tag vX.Y.Z
4. Push tag: git push origin vX.Y.Z
5. GoReleaser will run the full release gates, use the committed release
   notes, sign the checksum, and publish the artifacts automatically
```
