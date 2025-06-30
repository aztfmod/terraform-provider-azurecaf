package azurecaf

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceMatrix generates a comprehensive test matrix for all resource types
// This helps identify patterns and edge cases across different resource categories
func TestResourceMatrix(t *testing.T) {
	// Group resources by category
	categories := map[string][]string{
		"Storage":       {},
		"Compute":       {},
		"Networking":    {},
		"Database":      {},
		"Security":      {},
		"Monitoring":    {},
		"Web":          {},
		"AI_ML":        {},
		"Integration":   {},
		"Other":        {},
	}

	// Categorize resources based on their names
	for resourceType := range ResourceDefinitions {
		switch {
		case containsAny(resourceType, []string{"storage", "blob", "file", "disk", "backup"}):
			categories["Storage"] = append(categories["Storage"], resourceType)
		case containsAny(resourceType, []string{"vm", "virtual_machine", "scale_set", "batch", "container", "kubernetes", "aks"}):
			categories["Compute"] = append(categories["Compute"], resourceType)
		case containsAny(resourceType, []string{"network", "vnet", "subnet", "gateway", "dns", "lb", "traffic", "firewall", "bastion"}):
			categories["Networking"] = append(categories["Networking"], resourceType)
		case containsAny(resourceType, []string{"sql", "database", "mysql", "postgresql", "cosmos", "redis", "mariadb"}):
			categories["Database"] = append(categories["Database"], resourceType)
		case containsAny(resourceType, []string{"key_vault", "security", "managed_identity", "role"}):
			categories["Security"] = append(categories["Security"], resourceType)
		case containsAny(resourceType, []string{"monitor", "log", "insight", "alert", "diagnostic"}):
			categories["Monitoring"] = append(categories["Monitoring"], resourceType)
		case containsAny(resourceType, []string{"app_service", "web", "function", "logic_app", "api"}):
			categories["Web"] = append(categories["Web"], resourceType)
		case containsAny(resourceType, []string{"cognitive", "machine_learning", "search", "bot"}):
			categories["AI_ML"] = append(categories["AI_ML"], resourceType)
		case containsAny(resourceType, []string{"servicebus", "eventhub", "relay", "notification"}):
			categories["Integration"] = append(categories["Integration"], resourceType)
		default:
			categories["Other"] = append(categories["Other"], resourceType)
		}
	}

	// Sort categories for consistent output
	for category := range categories {
		sort.Strings(categories[category])
	}

	// Test each category
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("azurecaf_name resource not found")
	}

	totalTested := 0
	for category, resources := range categories {
		if len(resources) == 0 {
			continue
		}

		t.Run(fmt.Sprintf("Category_%s", category), func(t *testing.T) {
			t.Logf("Testing %s category with %d resources", category, len(resources))
			
			categorySuccesses := 0
			for _, resourceType := range resources {
				t.Run(sanitizeResourceType(resourceType), func(t *testing.T) {
					// Test the resource with lowercase category prefix
					categoryPrefix := strings.ToLower(category)
					if categoryPrefix == "ai_ml" {
						categoryPrefix = "ai" // Simplify for compatibility
					}
					
					resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
						"name":          "test",
						"resource_type": resourceType,
						"prefixes":      []interface{}{categoryPrefix},
						"clean_input":   true,
					})

					err := nameResource.Create(resourceData, nil)
					if err != nil {
						t.Errorf("Failed for %s: %v", resourceType, err)
						return
					}

					result := resourceData.Get("result").(string)
					if result == "" {
						t.Errorf("Empty result for %s", resourceType)
						return
					}

					categorySuccesses++
					t.Logf("✓ %s: %s", resourceType, result)
				})
			}
			
			successRate := float64(categorySuccesses) / float64(len(resources)) * 100
			t.Logf("Category %s: %d/%d successful (%.1f%%)", category, categorySuccesses, len(resources), successRate)
		})
		
		totalTested += len(resources)
	}

	// Summary
	t.Logf("\n=== RESOURCE MATRIX SUMMARY ===")
	for category, resources := range categories {
		if len(resources) > 0 {
			t.Logf("%-15s: %3d resources", category, len(resources))
		}
	}
	t.Logf("%-15s: %3d resources", "TOTAL", totalTested)
}

// containsAny checks if the string contains any of the specified substrings
func containsAny(s string, substrings []string) bool {
	for _, substring := range substrings {
		if len(s) > 0 && len(substring) > 0 {
			// Simple case-insensitive contains check
			for i := 0; i <= len(s)-len(substring); i++ {
				match := true
				for j := 0; j < len(substring); j++ {
					if s[i+j] != substring[j] && s[i+j] != substring[j]-32 && s[i+j] != substring[j]+32 {
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
		}
	}
	return false
}

// TestResourceConstraints validates that all resource definitions meet expected constraints
func TestResourceConstraints(t *testing.T) {
	constraints := map[string]func(string, ResourceStructure) bool{
		"ValidMinLength": func(name string, def ResourceStructure) bool {
			return def.MinLength >= 1 && def.MinLength <= 128
		},
		"ValidMaxLength": func(name string, def ResourceStructure) bool {
			return def.MaxLength >= def.MinLength && def.MaxLength <= 2048
		},
		"ValidCafPrefix": func(name string, def ResourceStructure) bool {
			// General resource types are designed not to use CAF prefixes
			if name == "general" || name == "general_safe" {
				return def.CafPrefix == ""
			}
			return len(def.CafPrefix) >= 1 && len(def.CafPrefix) <= 20
		},
		"ValidScope": func(name string, def ResourceStructure) bool {
			validScopes := []string{"global", "resourceGroup", "parent", "region", "subscription", "managementGroup", "tenant", "location", "assignment", "definition"}
			for _, scope := range validScopes {
				if def.Scope == scope {
					return true
				}
			}
			return false
		},
	}

	for constraintName, constraint := range constraints {
		t.Run(constraintName, func(t *testing.T) {
			violations := []string{}
			for resourceType, definition := range ResourceDefinitions {
				if !constraint(resourceType, definition) {
					violations = append(violations, resourceType)
				}
			}
			
			if len(violations) > 0 {
				t.Errorf("Constraint %s violated by %d resources: %v", constraintName, len(violations), violations[:min(10, len(violations))])
			} else {
				t.Logf("✓ All %d resources pass constraint %s", len(ResourceDefinitions), constraintName)
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
