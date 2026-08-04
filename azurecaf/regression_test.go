package azurecaf

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var updateSlugs = flag.Bool("update-slugs", false, "rewrite azurecaf/testdata/known_slugs.json from current ResourceDefinitions")

const knownSlugsPath = "testdata/known_slugs.json"

// --- Slug Drift Detection ---

// TestRegression_SlugConsistency ensures no slug silently changes between releases.
// If a slug needs to change, the golden file must be explicitly updated.
func TestRegression_SlugConsistency(t *testing.T) {
	current := currentKnownSlugs()
	if *updateSlugs {
		writeKnownSlugsFile(t, current)
	}

	expected := readKnownSlugsFile(t)
	issues := make([]string, 0)

	for resourceType, expectedSlug := range expected {
		actualSlug, ok := current[resourceType]
		if !ok {
			issues = append(issues, fmt.Sprintf("resource %q was removed from ResourceDefinitions (expected slug %q)", resourceType, expectedSlug))
			continue
		}
		if actualSlug != expectedSlug {
			issues = append(issues, fmt.Sprintf("resource %q slug drifted: expected %q, got %q", resourceType, expectedSlug, actualSlug))
		}
	}

	for resourceType, actualSlug := range current {
		if _, ok := expected[resourceType]; !ok {
			issues = append(issues, fmt.Sprintf("resource %q was added with slug %q; rerun with -update-slugs to accept it", resourceType, actualSlug))
		}
	}

	sort.Strings(issues)
	if len(issues) > 0 {
		t.Fatalf("slug snapshot mismatch:\n%s", strings.Join(issues, "\n"))
	}
}

// --- Known Issue Regressions ---

// TestRegression_Issues contains test cases for previously reported bugs.
// Each test is tagged with the issue number for traceability.
func TestRegression_Issues(t *testing.T) {
	namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
	cases := []struct {
		issue        int
		description  string
		resourceType string
		name         string
		prefixes     []string
		suffixes     []string
		separator    string
		expected     string // expected result or "ERROR" for expected failures
	}{
		{issue: 26, description: "event hub namespace still rejects numeric-only prefixes", resourceType: "azurerm_eventhub_namespace", name: "core", prefixes: []string{"123"}, separator: "-", expected: "ERROR"},
		{issue: 36, description: "bastion host keeps the corrected bast slug", resourceType: "azurerm_bastion_host", name: "jump", separator: "-", expected: "bast-jump"},
		{issue: 41, description: "recovery services vault keeps the rsv slug", resourceType: "azurerm_recovery_services_vault", name: "backup", separator: "-", expected: "rsv-backup"},
		{issue: 41, description: "hdinsight rserver cluster no longer reuses the recovery vault slug", resourceType: "azurerm_hdinsight_rserver_cluster", name: "cluster", separator: "-", expected: "rser-cluster"},
		{issue: 73, description: "key vault still errors when the name would start with a number", resourceType: "azurerm_key_vault", name: "vmsecrets", prefixes: []string{"5"}, separator: "-", expected: "ERROR"},
		{issue: 89, description: "cosmos db account keeps lowercase hyphenated names within the stricter rule set", resourceType: "azurerm_cosmosdb_account", name: "my-cosmos-01", separator: "-", expected: "cosmos-my-cosmos-01"},
		{issue: 107, description: "frontdoor firewall policy strips the default separator", resourceType: "azurerm_frontdoor_firewall_policy", name: "default", separator: "-", expected: "fdfwdefault"},
		{issue: 120, description: "service bus topic authorization rule keeps the corrected slug", resourceType: "azurerm_servicebus_topic_authorization_rule", name: "send", separator: "-", expected: "sbtar-send"},
		{issue: 162, description: "mysql flexible server database accepts the generated CAF name", resourceType: "azurerm_mysql_flexible_server_database", name: "exampledb", prefixes: []string{"qwud"}, separator: "-", expected: "qwud-mysqlfdb-exampledb"},
		{issue: 248, description: "powerbi embedded strips disallowed separators and underscores", resourceType: "azurerm_powerbi_embedded", name: "report_01", separator: "-", expected: "pbireport01"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("issue_%d_%s", tc.issue, strings.ReplaceAll(tc.resourceType, "azurerm_", "")), func(t *testing.T) {
			resourceName, err := getResourceName(
				tc.resourceType,
				tc.separator,
				append([]string(nil), tc.prefixes...),
				tc.name,
				append([]string(nil), tc.suffixes...),
				"",
				ConventionCafClassic,
				true,
				false,
				true,
				namePrecedence,
				false,
			)

			if tc.expected == "ERROR" {
				if err == nil {
					t.Fatalf("expected an error for issue #%d (%s), got result %q", tc.issue, tc.description, resourceName)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for issue #%d (%s): %v", tc.issue, tc.description, err)
			}

			if resourceName != tc.expected {
				t.Fatalf("unexpected result for issue #%d (%s): got %q want %q", tc.issue, tc.description, resourceName, tc.expected)
			}

			resource, err := getResource(tc.resourceType)
			if err != nil {
				t.Fatalf("failed to reload resource definition for %s: %v", tc.resourceType, err)
			}
			validationRegEx := regexp.MustCompile(resource.ValidationRegExp)
			if !validationRegEx.MatchString(resourceName) {
				t.Fatalf("result %q for issue #%d does not match validation regex %q", resourceName, tc.issue, resource.ValidationRegExp)
			}
		})
	}
}

// --- Resource Definition Integrity ---

// TestRegression_ResourceDefinitionIntegrity validates that all resource
// definitions have valid regex, sane lengths, and required fields.
func TestRegression_ResourceDefinitionIntegrity(t *testing.T) {
	for name, def := range ResourceDefinitions {
		t.Run(name, func(t *testing.T) {
			if def.RegEx != "" {
				if _, err := regexp.Compile(def.RegEx); err != nil {
					t.Errorf("invalid regex for %s: %v", name, err)
				}
			}

			if def.ValidationRegExp != "" {
				if _, err := regexp.Compile(def.ValidationRegExp); err != nil {
					t.Errorf("invalid validation regex for %s: %v", name, err)
				}
			}

			if def.MinLength > def.MaxLength {
				t.Errorf("min_length > max_length for %s", name)
			}

			if def.MinLength < 0 {
				t.Errorf("negative min_length %d for %s", def.MinLength, name)
			}

			if def.CafPrefix == "" && name != "general" && name != "general_safe" {
				t.Errorf("empty slug for %s", name)
			}

			if def.ResourceTypeName == "" {
				t.Errorf("empty resource type name for %s", name)
			}

			if def.MaxLength <= 0 || def.MaxLength > 4096 {
				t.Errorf("unreasonable max_length %d for %s", def.MaxLength, name)
			}
		})
	}
}

// --- State Migration ---

// TestRegression_StateV0ToV1Migration preserves the historical test name, but
// validates the current resource state upgrader that injects use_slug=true.
func TestRegression_StateV0ToV1Migration(t *testing.T) {
	resource := resourceName()
	if len(resource.StateUpgraders) != 1 {
		t.Fatalf("expected exactly one state upgrader, got %d", len(resource.StateUpgraders))
	}

	upgrader := resource.StateUpgraders[0]
	if upgrader.Version != 2 {
		t.Fatalf("expected state upgrader from version 2, got %d", upgrader.Version)
	}

	rawState := map[string]interface{}{
		"name":          "legacy",
		"resource_type": "azurerm_resource_group",
		"result":        "rg-legacy",
	}

	upgraded, err := upgrader.Upgrade(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("unexpected state upgrade error: %v", err)
	}

	if upgraded["use_slug"] != true {
		t.Fatalf("expected migrated state to enable use_slug, got %#v", upgraded["use_slug"])
	}

	if upgraded["name"] != rawState["name"] {
		t.Fatalf("expected migrated state to preserve name %q, got %#v", rawState["name"], upgraded["name"])
	}

	if upgraded["resource_type"] != rawState["resource_type"] {
		t.Fatalf("expected migrated state to preserve resource_type %q, got %#v", rawState["resource_type"], upgraded["resource_type"])
	}
}

func currentKnownSlugs() map[string]string {
	snapshot := make(map[string]string, len(ResourceDefinitions))
	for resourceType, def := range ResourceDefinitions {
		snapshot[resourceType] = def.CafPrefix
	}
	return snapshot
}

func readKnownSlugsFile(t *testing.T) map[string]string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("testdata", filepath.Base(knownSlugsPath)))
	if err != nil {
		t.Fatalf("failed to read %s: %v", knownSlugsPath, err)
	}

	var slugs map[string]string
	if err := json.Unmarshal(data, &slugs); err != nil {
		t.Fatalf("failed to unmarshal %s: %v", knownSlugsPath, err)
	}
	return slugs
}

func writeKnownSlugsFile(t *testing.T, slugs map[string]string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(knownSlugsPath), 0o755); err != nil {
		t.Fatalf("failed to create %s: %v", filepath.Dir(knownSlugsPath), err)
	}

	data, err := json.MarshalIndent(slugs, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal slug snapshot: %v", err)
	}

	data = append(data, '\n')
	if err := os.WriteFile(knownSlugsPath, data, 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", knownSlugsPath, err)
	}
}
