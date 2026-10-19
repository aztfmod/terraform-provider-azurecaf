# Release Runbook

Use this checklist for stable `vX.Y.Z` releases.

## Prepare the release

1. Run the pre-release validation procedure:
   - generated code freshness
   - `make test_ci`
   - `make test_e2e`
   - `make test_coverage`
   - plan/apply/plan lifecycle convergence
2. Generate `.github/releases/vX.Y.Z.md`.
3. Move the `CHANGELOG.md` entries into `## [vX.Y.Z] - YYYY-MM-DD`
   and leave `## [Unreleased]` empty.
4. Run:

   ```bash
   ./scripts/check-release-ready.sh vX.Y.Z
   ```

5. Commit the release notes and finalized CHANGELOG before requesting the
   final review.

## Merge readiness

Before merging the release pull request, record:

- the final head SHA
- required checks and their conclusions
- `reviewDecision`
- `mergeStateStatus`
- whether the branch is up to date

Use:

```bash
gh pr view <number> \
  --json headRefOid,reviewDecision,mergeStateStatus,statusCheckRollup
```

Do not lower branch protection. An admin bypass requires explicit maintainer
authorization and a recorded reason.

## Publish

1. Merge the reviewed release pull request.
2. Wait for the `main` push workflows to pass.
3. Create the stable tag at the merge commit and push it.
4. Monitor:
   - `Go`
   - `Release Validation`
   - `Release Labeler`
5. Do not retag or overwrite a published release.

## Verify

Confirm all of the following:

- the GitHub Release body matches `.github/releases/vX.Y.Z.md`
- all expected platform archives are present
- every archive matches `SHA256SUMS`
- `SHA256SUMS.sig` verifies with the trusted release key
- the Terraform Registry lists the released version
- the version label exists and is applied to the complete reviewed issue/PR set

The publishing workflow repairs a missing or mismatched release body from the
committed release-notes file and fails if the repair does not hold.

## Recovery

- **Release body missing:** update the existing release from the committed
  `.github/releases/vX.Y.Z.md`; do not rerun GoReleaser.
- **Version label missing:** create the exact `vX.Y.Z` label and rerun Release
  Labeler. Do not bulk-label issues until the Issue-to-PR workflow includes
  native `new-resource` event gating.
- **Validation failure:** keep the tag and release under investigation; use
  the Release Validation issue and job summary to determine whether a new
  patch release is required.
