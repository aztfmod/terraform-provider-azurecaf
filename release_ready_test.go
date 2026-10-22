package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseReadyScript(t *testing.T) {
	t.Parallel()

	script, err := filepath.Abs("scripts/check-release-ready.sh")
	if err != nil {
		t.Fatalf("resolve release readiness script: %v", err)
	}

	validChangelog := `# Changelog

## [Unreleased]

## [v2.1.0] - 2026-10-20

### Fixed

- Hardened the release workflow.

## [v2.0.0] - 2026-10-19

### Added

- Previous release.
`

	tests := []struct {
		name      string
		tag       string
		changelog string
		notes     string
		wantError string
	}{
		{
			name:      "ready",
			tag:       "v2.1.0",
			changelog: validChangelog,
			notes:     "## What's Changed\n\n- Hardened the release workflow.\n",
		},
		{
			name: "unreleased entries remain",
			tag:  "v2.1.0",
			changelog: strings.Replace(
				validChangelog,
				"## [Unreleased]\n",
				"## [Unreleased]\n\n### Fixed\n\n- Not finalized.\n",
				1,
			),
			notes:     "## What's Changed\n",
			wantError: "Unreleased section must be empty",
		},
		{
			name:      "version heading has no date",
			tag:       "v2.1.0",
			changelog: strings.Replace(validChangelog, "## [v2.1.0] - 2026-10-20", "## [v2.1.0]", 1),
			notes:     "## What's Changed\n",
			wantError: "one dated section",
		},
		{
			name:      "version heading is duplicated",
			tag:       "v2.1.0",
			changelog: strings.Replace(validChangelog, "## [v2.0.0] - 2026-10-19", "## [v2.1.0] - invalid\n\n## [v2.0.0] - 2026-10-19", 1),
			notes:     "## What's Changed\n",
			wantError: "exactly one section",
		},
		{
			name:      "release section is empty",
			tag:       "v2.1.0",
			changelog: strings.Replace(validChangelog, "### Fixed\n\n- Hardened the release workflow.\n\n", "", 1),
			notes:     "## What's Changed\n",
			wantError: "has no release entries",
		},
		{
			name:      "release notes are empty",
			tag:       "v2.1.0",
			changelog: validChangelog,
			notes:     "",
			wantError: "missing or empty release notes",
		},
		{
			name:      "tag is not stable semver",
			tag:       "v2.1.0-preview.1",
			changelog: validChangelog,
			notes:     "## What's Changed\n",
			wantError: "stable semantic-version tag",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, ".github", "releases"), 0o755); err != nil {
				t.Fatalf("create release notes directory: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte(tt.changelog), 0o644); err != nil {
				t.Fatalf("write CHANGELOG fixture: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, ".github", "releases", tt.tag+".md"), []byte(tt.notes), 0o644); err != nil {
				t.Fatalf("write release notes fixture: %v", err)
			}

			cmd := exec.Command("bash", script, tt.tag)
			cmd.Env = append(os.Environ(), "RELEASE_ROOT="+root)
			output, err := cmd.CombinedOutput()
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("expected release metadata to pass: %v\n%s", err, output)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected release metadata failure containing %q", tt.wantError)
			}
			if !strings.Contains(string(output), tt.wantError) {
				t.Fatalf("expected error containing %q, got:\n%s", tt.wantError, output)
			}
		})
	}
}
