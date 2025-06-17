package azurecaf

import (
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestSlugValidation performs comprehensive validation of all resource slugs
func TestSlugValidation(t *testing.T) {
	// Collect all resources and their slugs
	allSlugs := make(map[string][]string)
	var allResourceTypes []string
	
	for resourceType := range ResourceDefinitions {
		resource, err := getResource(resourceType)
		if err != nil {
			t.Errorf("Failed to get resource %s: %v", resourceType, err)
			continue
		}
		
		slug := resource.CafPrefix
		allSlugs[slug] = append(allSlugs[slug], resourceType)
		allResourceTypes = append(allResourceTypes, resourceType)
	}
	
	t.Logf("Validating %d resource types with %d unique slugs", len(allResourceTypes), len(allSlugs))
	
	// Run all validation checks
	t.Run("UniquenessCheck", func(t *testing.T) {
		validateSlugUniqueness(t, allSlugs)
	})
	
	t.Run("FormatValidation", func(t *testing.T) {
		validateSlugFormat(t, allSlugs)
	})
	
	t.Run("LengthAnalysis", func(t *testing.T) {
		analyzeSlugLengths(t, allSlugs)
	})
	
	t.Run("NamingConsistency", func(t *testing.T) {
		validateNamingConsistency(t, allSlugs)
	})
}

// validateSlugUniqueness ensures slugs are not inappropriately duplicated
func validateSlugUniqueness(t *testing.T, allSlugs map[string][]string) {
	duplicateCount := 0
	var duplicatedSlugs []string
	
	for slug, resources := range allSlugs {
		if len(resources) > 1 {
			duplicateCount++
			duplicatedSlugs = append(duplicatedSlugs, slug)
		}
	}
	
	if duplicateCount > 0 {
		sort.Strings(duplicatedSlugs)
		t.Logf("Found %d duplicate slugs (this may cause naming conflicts):", duplicateCount)
		
		// Count critical vs non-critical duplicates
		criticalDuplicates := 0
		for _, slug := range duplicatedSlugs {
			resources := allSlugs[slug]
			sort.Strings(resources)
			
			// Determine if this is a critical duplicate (unrelated resources sharing a slug)
			isCritical := false
			if len(slug) == 0 {
				// Empty slug is always critical
				isCritical = true
			} else if !areResourcesRelated(resources) {
				// Unrelated resources sharing a slug is critical
				isCritical = true
			}
			
			if isCritical {
				criticalDuplicates++
				t.Errorf("  CRITICAL: '%s' used by unrelated resources: %v", slug, resources)
			} else {
				t.Logf("  INFO: '%s' used by related resources: %v", slug, resources)
			}
		}
		
		if criticalDuplicates > 0 {
			t.Errorf("Found %d critical slug duplicates that should be resolved", criticalDuplicates)
		} else {
			t.Logf("✓ No critical slug duplicates found (all duplicates are for related resources)")
		}
	} else {
		t.Logf("✓ All slugs are unique")
	}
}

// validateSlugFormat ensures slugs follow consistent formatting rules
func validateSlugFormat(t *testing.T, allSlugs map[string][]string) {
	var issues []string
	validSlugPattern := regexp.MustCompile(`^[a-z][a-z0-9]*$`)
	
	for slug := range allSlugs {
		// Check if slug follows lowercase alphanumeric pattern
		if !validSlugPattern.MatchString(slug) {
			issues = append(issues, "Slug '"+slug+"' contains invalid characters (should be lowercase alphanumeric)")
		}
		
		// Check for overly short slugs
		if len(slug) == 1 {
			issues = append(issues, "Slug '"+slug+"' is very short (may cause conflicts)")
		}
		
		// Check for overly long slugs
		if len(slug) > 12 {
			issues = append(issues, "Slug '"+slug+"' is very long (may be unwieldy)")
		}
		
		// Check for consecutive numbers at the end (may indicate versioning issues)
		if strings.HasSuffix(slug, "2") || strings.HasSuffix(slug, "3") {
			issues = append(issues, "Slug '"+slug+"' ends with number (may indicate duplicate)")
		}
	}
	
	if len(issues) > 0 {
		sort.Strings(issues)
		t.Logf("Slug format observations (%d):", len(issues))
		for _, issue := range issues {
			t.Logf("  %s", issue)
		}
	} else {
		t.Logf("✓ All slugs follow proper format")
	}
}

// analyzeSlugLengths provides statistics on slug length distribution
func analyzeSlugLengths(t *testing.T, allSlugs map[string][]string) {
	lengthDistribution := make(map[int][]string)
	
	for slug := range allSlugs {
		length := len(slug)
		lengthDistribution[length] = append(lengthDistribution[length], slug)
	}
	
	var lengths []int
	for length := range lengthDistribution {
		lengths = append(lengths, length)
	}
	sort.Ints(lengths)
	
	t.Logf("Slug length distribution:")
	for _, length := range lengths {
		slugs := lengthDistribution[length]
		sort.Strings(slugs)
		t.Logf("  Length %d (%d slugs): %v", length, len(slugs), slugs)
	}
	
	// Identify optimal length ranges
	optimalCount := len(lengthDistribution[3]) + len(lengthDistribution[4]) + len(lengthDistribution[5])
	totalCount := len(allSlugs)
	
	t.Logf("Optimal length (3-5 chars): %d/%d (%.1f%%)", optimalCount, totalCount, float64(optimalCount)/float64(totalCount)*100)
}

// validateNamingConsistency checks for consistent naming patterns within service families
func validateNamingConsistency(t *testing.T, allSlugs map[string][]string) {
	// Group resources by service families and analyze patterns
	servicePatterns := make(map[string]map[string][]string)
	
	for slug, resources := range allSlugs {
		for _, resourceType := range resources {
			family := getServiceFamily(resourceType)
			if servicePatterns[family] == nil {
				servicePatterns[family] = make(map[string][]string)
			}
			servicePatterns[family][slug] = append(servicePatterns[family][slug], resourceType)
		}
	}
	
	t.Logf("Service family consistency analysis:")
	
	var families []string
	for family := range servicePatterns {
		if family != "Other" { // Skip the catch-all category for now
			families = append(families, family)
		}
	}
	sort.Strings(families)
	
	for _, family := range families {
		patterns := servicePatterns[family]
		if len(patterns) < 2 {
			continue // Skip families with only one resource
		}
		
		var slugs []string
		for slug := range patterns {
			slugs = append(slugs, slug)
		}
		sort.Strings(slugs)
		
		t.Logf("  %s family (%d resources):", family, len(patterns))
		
		// Analyze prefix consistency
		prefixes := make(map[string]int)
		for _, slug := range slugs {
			if len(slug) >= 2 {
				prefix := slug[:2]
				prefixes[prefix]++
			}
		}
		
		var commonPrefixes []string
		for prefix, count := range prefixes {
			if count > 1 {
				commonPrefixes = append(commonPrefixes, prefix)
			}
		}
		
		if len(commonPrefixes) > 0 {
			sort.Strings(commonPrefixes)
			t.Logf("    Common prefixes: %v", commonPrefixes)
		}
		
		// Show all slugs in this family
		for _, slug := range slugs {
			resources := patterns[slug]
			t.Logf("    %s: %v", slug, resources)
		}
	}
}

// areResourcesRelated checks if a set of resources sharing a slug are logically related
func areResourcesRelated(resources []string) bool {
	if len(resources) <= 1 {
		return true
	}
	
	// Check for common patterns that indicate related resources
	basePatterns := []string{
		"azurerm_api_management",        // API Management services
		"azurerm_mssql",                 // Microsoft SQL services
		"azurerm_sql",                   // SQL services (legacy)
		"azurerm_virtual_machine",       // Virtual machine variants
		"azurerm_linux_virtual_machine", // Linux VM variants
		"azurerm_windows_virtual_machine", // Windows VM variants
		"virtual_machine_scale_set",     // All VM scale set variants
		"azurerm_dns_",                  // DNS record types
		"azurerm_private_dns_",          // Private DNS record types
		"azurerm_storage_share",         // Storage share related
		"azurerm_logic_app_trigger",     // Logic app triggers
		"azurerm_network_security",      // Network security rules
		"azurerm_dev_test_",             // Dev test lab VMs
		"azurerm_cognitive",             // Cognitive services
		"dashboard",                     // Dashboard services
		"general",                       // General naming resources
	}
	
	for _, pattern := range basePatterns {
		relatedCount := 0
		for _, resource := range resources {
			if strings.HasPrefix(resource, pattern) || strings.Contains(resource, pattern) {
				relatedCount++
			}
		}
		
		// If most resources match this pattern, they're related
		if relatedCount >= len(resources)*2/3 {
			return true
		}
	}
	
	// Check if resources share a common service prefix
	if len(resources) >= 2 {
		firstResource := resources[0]
		parts := strings.Split(firstResource, "_")
		if len(parts) >= 3 {
			commonPrefix := strings.Join(parts[:3], "_") // e.g., "azurerm_api_management"
			relatedCount := 0
			for _, resource := range resources {
				if strings.HasPrefix(resource, commonPrefix) {
					relatedCount++
				}
			}
			if relatedCount == len(resources) {
				return true
			}
		}
	}
	
	return false
}

// getServiceFamily categorizes a resource into a service family
func getServiceFamily(resourceType string) string {
	switch {
	case strings.Contains(resourceType, "monitor"):
		return "Azure Monitor"
	case strings.Contains(resourceType, "lb_") || strings.Contains(resourceType, "load_balancer"):
		return "Load Balancer"
	case strings.Contains(resourceType, "storage") || strings.Contains(resourceType, "blob") || strings.Contains(resourceType, "file"):
		return "Storage"
	case strings.Contains(resourceType, "network") || strings.Contains(resourceType, "vnet") || strings.Contains(resourceType, "subnet"):
		return "Networking"
	case strings.Contains(resourceType, "sql") || strings.Contains(resourceType, "database") || strings.Contains(resourceType, "mysql") || strings.Contains(resourceType, "postgres"):
		return "Database"
	case strings.Contains(resourceType, "vm") || strings.Contains(resourceType, "virtual_machine") || strings.Contains(resourceType, "scale_set"):
		return "Compute"
	case strings.Contains(resourceType, "key_vault") || strings.Contains(resourceType, "security") || strings.Contains(resourceType, "firewall"):
		return "Security"
	case strings.Contains(resourceType, "app_service") || strings.Contains(resourceType, "function") || strings.Contains(resourceType, "logic_app"):
		return "App Services"
	case strings.Contains(resourceType, "api_management"):
		return "API Management"
	case strings.Contains(resourceType, "automation"):
		return "Automation"
	case strings.Contains(resourceType, "container") || strings.Contains(resourceType, "kubernetes"):
		return "Containers"
	default:
		return "Other"
	}
}
