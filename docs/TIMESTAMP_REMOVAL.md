# Code Generation Timestamp Removal

## Problem
The `models_generated.go` file contained a timestamp that was dynamically generated during the build process. This caused the file to change on every generation, leading to a "git dirty state" error in CI/CD pipelines, particularly with GoReleaser.

## Solution
Removed the timestamp from the generated code by:

1. **Updated `gen.go`**:
   - Removed `GeneratedTime` field from `templateData` struct
   - Removed `time` import (no longer needed)
   - Simplified template execution to not pass timestamp

2. **Updated `templates/model.tmpl`**:
   - Removed timestamp line from generated file header
   - Simplified header comment

## Result
- Generated files are now completely stable across multiple runs
- No more git dirty state issues in CI/CD
- GoReleaser can successfully create releases without modification conflicts

## Impact
- **High**: Fixes critical CI/CD pipeline failures
- **Low Risk**: No functional changes to the provider logic
- **Maintenance**: Reduces build-time variability

## Testing
- Generated files are identical across multiple generation runs
- No timestamp information present in generated content
- Build process remains functional
