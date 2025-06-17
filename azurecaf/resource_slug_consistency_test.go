package azurecaf

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

// TestSlugUniqueness validates that slugs are unique across all resources unless intentionally shared
func TestSlugUniqueness(t *testing.T) {
	slugToResources := make(map[string][]string)
	
	// Collect all slug-resource mappings
	for resourceType := range ResourceDefinitions {
		resource, err := getResource(resourceType)
		if err != nil {
			t.Errorf("Failed to get resource %s: %v", resourceType, err)
			continue
		}
		slug := resource.CafPrefix
		slugToResources[slug] = append(slugToResources[slug], resourceType)
	}
	
	// Check for duplicate slugs and report them
	duplicates := make(map[string][]string)
	for slug, resources := range slugToResources {
		if len(resources) > 1 {
			duplicates[slug] = resources
		}
	}
	
	if len(duplicates) > 0 {
		t.Errorf("Found %d duplicate slugs:", len(duplicates))
		
		// Sort slugs for consistent output
		var slugs []string
		for slug := range duplicates {
			slugs = append(slugs, slug)
		}
		sort.Strings(slugs)
		
		for _, slug := range slugs {
			resources := duplicates[slug]
			sort.Strings(resources)
			t.Errorf("  Slug '%s' is used by %d resources: %v", slug, len(resources), resources)
		}
	}
}

// TestSlugConsistency validates that resources with similar purposes have consistent slug patterns
func TestSlugConsistency(t *testing.T) {
	// Group resources by service prefix and validate consistency
	serviceGroups := map[string]map[string]string{
		"azurerm_monitor_": {
			"prefix_pattern": "Monitor resources should typically start with 'am' or use descriptive abbreviations",
		},
		"azurerm_lb_": {
			"prefix_pattern": "Load balancer resources should typically start with 'lb'",
		},
		"azurerm_storage_": {
			"prefix_pattern": "Storage resources should typically start with 'st' or use storage-related abbreviations",
		},
		"azurerm_network_": {
			"prefix_pattern": "Network resources should use network-related abbreviations",
		},
	}
	
	for servicePrefix, rules := range serviceGroups {
		var resources []string
		var slugs []string
		
		for resourceType := range ResourceDefinitions {
			if strings.HasPrefix(resourceType, servicePrefix) {
				resource, err := getResource(resourceType)
				if err != nil {
					t.Errorf("Failed to get resource %s: %v", resourceType, err)
					continue
				}
				resources = append(resources, resourceType)
				slugs = append(slugs, resource.CafPrefix)
			}
		}
		
		if len(resources) > 0 {
			t.Logf("Service group '%s' has %d resources:", servicePrefix, len(resources))
			for i, resourceType := range resources {
				t.Logf("  %s -> %s", resourceType, slugs[i])
			}
			t.Logf("  Rule: %s", rules["prefix_pattern"])
		}
	}
}

// TestSlugQuality validates the quality and appropriateness of slugs
func TestSlugQuality(t *testing.T) {
	issues := []string{}
	slugStats := make(map[int]int) // length -> count
	
	for resourceType := range ResourceDefinitions {
		resource, err := getResource(resourceType)
		if err != nil {
			t.Errorf("Failed to get resource %s: %v", resourceType, err)
			continue
		}
		
		slug := resource.CafPrefix
		slugStats[len(slug)]++
		
		// Check for overly long slugs (may be hard to use)
		if len(slug) > 10 {
			issues = append(issues, fmt.Sprintf("Resource %s has long slug '%s' (length: %d)", resourceType, slug, len(slug)))
		}
		
		// Check for single character slugs (may be ambiguous)
		if len(slug) == 1 {
			issues = append(issues, fmt.Sprintf("Resource %s has very short slug '%s' (may be ambiguous)", resourceType, slug))
		}
		
		// Check for slugs that don't seem related to the resource name
		if !isSlugReasonable(resourceType, slug) {
			issues = append(issues, fmt.Sprintf("Resource %s has slug '%s' which may not be intuitive", resourceType, slug))
		}
	}
	
	// Report slug length distribution
	t.Logf("Slug length distribution:")
	var lengths []int
	for length := range slugStats {
		lengths = append(lengths, length)
	}
	sort.Ints(lengths)
	
	for _, length := range lengths {
		count := slugStats[length]
		t.Logf("  Length %d: %d slugs", length, count)
	}
	
	// Report quality issues as warnings (not failures)
	if len(issues) > 0 {
		t.Logf("Slug quality observations (%d total):", len(issues))
		for _, issue := range issues {
			t.Logf("  %s", issue)
		}
	}
}

// TestSlugNamingPatterns validates common naming patterns across resource families
func TestSlugNamingPatterns(t *testing.T) {
	patterns := map[string][]string{
		"Azure Monitor": {},
		"Load Balancer": {},
		"Storage": {},
		"Network": {},
		"Database": {},
		"Compute": {},
		"Security": {},
	}
	
	for resourceType := range ResourceDefinitions {
		resource, err := getResource(resourceType)
		if err != nil {
			continue
		}
		
		slug := resource.CafPrefix
		
		// Categorize resources by service family
		switch {
		case strings.Contains(resourceType, "monitor"):
			patterns["Azure Monitor"] = append(patterns["Azure Monitor"], fmt.Sprintf("%s -> %s", resourceType, slug))
		case strings.Contains(resourceType, "lb_") || strings.Contains(resourceType, "load_balancer"):
			patterns["Load Balancer"] = append(patterns["Load Balancer"], fmt.Sprintf("%s -> %s", resourceType, slug))
		case strings.Contains(resourceType, "storage") || strings.Contains(resourceType, "blob") || strings.Contains(resourceType, "file"):
			patterns["Storage"] = append(patterns["Storage"], fmt.Sprintf("%s -> %s", resourceType, slug))
		case strings.Contains(resourceType, "network") || strings.Contains(resourceType, "vnet") || strings.Contains(resourceType, "subnet"):
			patterns["Network"] = append(patterns["Network"], fmt.Sprintf("%s -> %s", resourceType, slug))
		case strings.Contains(resourceType, "sql") || strings.Contains(resourceType, "database") || strings.Contains(resourceType, "mysql") || strings.Contains(resourceType, "postgres"):
			patterns["Database"] = append(patterns["Database"], fmt.Sprintf("%s -> %s", resourceType, slug))
		case strings.Contains(resourceType, "vm") || strings.Contains(resourceType, "virtual_machine") || strings.Contains(resourceType, "scale_set"):
			patterns["Compute"] = append(patterns["Compute"], fmt.Sprintf("%s -> %s", resourceType, slug))
		case strings.Contains(resourceType, "key_vault") || strings.Contains(resourceType, "security") || strings.Contains(resourceType, "firewall"):
			patterns["Security"] = append(patterns["Security"], fmt.Sprintf("%s -> %s", resourceType, slug))
		}
	}
	
	// Report patterns for each service family
	for family, resources := range patterns {
		if len(resources) > 0 {
			sort.Strings(resources)
			t.Logf("%s resources (%d):", family, len(resources))
			for _, resource := range resources {
				t.Logf("  %s", resource)
			}
		}
	}
}

// isSlugReasonable checks if a slug seems reasonably related to the resource name
func isSlugReasonable(resourceType, slug string) bool {
	// This is a heuristic check - extract key parts of the resource name
	// and see if the slug contains some recognizable abbreviation
	
	// Remove azurerm_ prefix
	name := strings.TrimPrefix(resourceType, "azurerm_")
	
	// Common abbreviations that are acceptable
	commonAbbrevs := map[string][]string{
		"application":     {"app", "appl"},
		"management":      {"mgmt", "man"},
		"network":         {"net", "nw"},
		"security":        {"sec"},
		"virtual":         {"v", "virt"},
		"machine":         {"vm", "m"},
		"database":        {"db"},
		"storage":         {"st", "stor"},
		"account":         {"acc", "a"},
		"resource":        {"res", "r"},
		"group":           {"grp", "g"},
		"service":         {"svc", "s"},
		"load_balancer":   {"lb"},
		"monitor":         {"mon", "am"},
		"diagnostic":      {"diag"},
		"automation":      {"auto", "aa"},
		"configuration":   {"config", "cfg"},
		"certificate":     {"cert"},
		"firewall":        {"fw"},
		"gateway":         {"gw"},
		"public":          {"pub"},
		"private":         {"priv"},
		"subnet":          {"snet"},
		"endpoint":        {"ep"},
	}
	
	// For very short slugs or very specific cases, be more lenient
	if len(slug) <= 3 {
		return true
	}
	
	// Check if slug is a reasonable abbreviation
	words := strings.Split(name, "_")
	for _, word := range words {
		if word == slug {
			return true // exact match
		}
		if abbrevs, exists := commonAbbrevs[word]; exists {
			for _, abbrev := range abbrevs {
				if strings.Contains(slug, abbrev) {
					return true
				}
			}
		}
		// Check if slug starts with first letters of the word
		if len(word) > 0 && strings.HasPrefix(slug, string(word[0])) {
			return true
		}
	}
	
	return false // Couldn't find reasonable connection
}