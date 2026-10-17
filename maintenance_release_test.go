package main

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"testing"
	"text/template"
)

func TestMaintenanceWorkflowBranches(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		path  string
		event string
	}{
		{".github/workflows/go.yml", "push"},
		{".github/workflows/go.yml", "pull_request"},
		{".github/workflows/e2e.yml", "pull_request"},
	} {
		t.Run(test.path+"/"+test.event, func(t *testing.T) {
			t.Parallel()

			content := maintenanceReleaseFile(t, test.path)
			pattern := regexp.MustCompile(`(?m)^  ` + test.event + `:\n    branches: \[main, release/1\.x\]$`)
			if !pattern.MatchString(content) {
				t.Fatalf("%s must check %s events for main and release/1.x", test.path, test.event)
			}
		})
	}
}

func TestMaintenanceReleaseLatestPolicy(t *testing.T) {
	t.Parallel()

	content := maintenanceReleaseFile(t, ".goreleaser.yml")
	if !strings.Contains(content, "  prerelease: auto\n") {
		t.Fatal("release prerelease status must follow the version tag")
	}
	match := regexp.MustCompile(`(?m)^  make_latest: '(.+)'$`).FindStringSubmatch(content)
	if len(match) != 2 {
		t.Fatal("release.make_latest must explicitly distinguish maintenance and prerelease versions")
	}
	policy, err := template.New("make_latest").Parse(match[1])
	if err != nil {
		t.Fatalf("parse release.make_latest: %v", err)
	}

	for _, test := range []struct {
		name       string
		major      uint64
		prerelease string
		want       string
	}{
		{"v1 stable", 1, "", "false"},
		{"v1 prerelease", 1, "rc.1", "false"},
		{"v2 stable", 2, "", "true"},
		{"v2 prerelease", 2, "rc.1", "false"},
		{"future major", 3, "", "true"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var result bytes.Buffer
			data := struct {
				Major      uint64
				Prerelease string
			}{test.major, test.prerelease}
			if err := policy.Execute(&result, data); err != nil {
				t.Fatalf("render release.make_latest: %v", err)
			}
			if result.String() != test.want {
				t.Fatalf("release.make_latest = %q, want %q", result.String(), test.want)
			}
		})
	}
}

func TestMaintenanceReleaseArtifacts(t *testing.T) {
	t.Parallel()

	for path, markers := range map[string][]string{
		".goreleaser.yml": {
			"archives:\n  - formats:\n      - zip\n",
			`binary: "{{ .ProjectName }}_v{{ .Version }}"`,
			`name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"`,
			`name_template: "{{ .ProjectName }}_{{ .Version }}_SHA256SUMS"`,
			"algorithm: sha256",
			"- artifacts: checksum",
			`"--detach-sign"`,
		},
		".github/workflows/go.yml": {
			"tags:\n      - \"v*\"",
			"gpg_private_key: ${{ secrets.GPG_PRIVATE_KEY }}",
			"GPG_FINGERPRINT: ${{ steps.import_gpg.outputs.fingerprint }}",
			"args: release --clean",
		},
	} {
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			content := maintenanceReleaseFile(t, path)
			for _, marker := range markers {
				if !strings.Contains(content, marker) {
					t.Errorf("%s is missing release contract %q", path, marker)
				}
			}
		})
	}
}

func maintenanceReleaseFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read release configuration %s: %v", path, err)
	}
	return string(content)
}
