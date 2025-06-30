package azurecaf

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceCoverage analyzes which resources are tested and which are missing
func TestResourceCoverage(t *testing.T) {
	// Get all resource types from ResourceDefinitions
	allResourceTypes := make([]string, 0, len(ResourceDefinitions))
	for resourceType := range ResourceDefinitions {
		allResourceTypes = append(allResourceTypes, resourceType)
	}
	sort.Strings(allResourceTypes)

	t.Logf("Total resource types defined: %d", len(allResourceTypes))

	// Test each resource type to ensure it works
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("azurecaf_name resource not found")
	}

	var successfulResources []string
	var failedResources []string

	for _, resourceType := range allResourceTypes {
		t.Run(fmt.Sprintf("Coverage_%s", sanitizeResourceType(resourceType)), func(t *testing.T) {
			// Test basic functionality
			resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
				"name":          "test",
				"resource_type": resourceType,
				"clean_input":   true,
			})

			err := nameResource.Create(resourceData, nil)
			if err != nil {
				failedResources = append(failedResources, resourceType)
				t.Errorf("Failed to create name for %s: %v", resourceType, err)
				return
			}

			result := resourceData.Get("result").(string)
			if result == "" {
				failedResources = append(failedResources, resourceType)
				t.Errorf("Empty result for %s", resourceType)
				return
			}

			successfulResources = append(successfulResources, resourceType)
			t.Logf("✓ %s: %s", resourceType, result)
		})
	}

	// Generate coverage report
	coverageReport := map[string]interface{}{
		"total_resources":     len(allResourceTypes),
		"successful_resources": len(successfulResources),
		"failed_resources":    len(failedResources),
		"coverage_percentage": float64(len(successfulResources)) / float64(len(allResourceTypes)) * 100,
		"successful_list":     successfulResources,
		"failed_list":         failedResources,
	}

	// Write coverage report to file
	reportData, _ := json.MarshalIndent(coverageReport, "", "  ")
	os.WriteFile("resource_coverage_report.json", reportData, 0644)

	t.Logf("Coverage Report:")
	t.Logf("  Total Resources: %d", len(allResourceTypes))
	t.Logf("  Successful: %d", len(successfulResources))
	t.Logf("  Failed: %d", len(failedResources))
	t.Logf("  Coverage: %.2f%%", float64(len(successfulResources))/float64(len(allResourceTypes))*100)
	
	if len(failedResources) > 0 {
		t.Logf("  Failed resources: %v", failedResources)
	}
}

// TestResourceDefinitionCompleteness verifies all resource definitions are valid
func TestResourceDefinitionCompleteness(t *testing.T) {
	var issues []string

	for resourceType, definition := range ResourceDefinitions {
		// Check minimum required fields
		if definition.ResourceTypeName == "" {
			issues = append(issues, fmt.Sprintf("%s: missing ResourceTypeName", resourceType))
		}
		// Skip CafPrefix check for general resource types that are designed not to use prefixes
		if definition.CafPrefix == "" && resourceType != "general" && resourceType != "general_safe" {
			issues = append(issues, fmt.Sprintf("%s: missing CafPrefix", resourceType))
		}
		if definition.MinLength <= 0 {
			issues = append(issues, fmt.Sprintf("%s: invalid MinLength %d", resourceType, definition.MinLength))
		}
		if definition.MaxLength <= 0 {
			issues = append(issues, fmt.Sprintf("%s: invalid MaxLength %d", resourceType, definition.MaxLength))
		}
		if definition.MinLength > definition.MaxLength {
			issues = append(issues, fmt.Sprintf("%s: MinLength (%d) > MaxLength (%d)", resourceType, definition.MinLength, definition.MaxLength))
		}
		if definition.ValidationRegExp == "" {
			issues = append(issues, fmt.Sprintf("%s: missing ValidationRegExp", resourceType))
		}
		if definition.RegEx == "" {
			issues = append(issues, fmt.Sprintf("%s: missing RegEx", resourceType))
		}
		if definition.Scope == "" {
			issues = append(issues, fmt.Sprintf("%s: missing Scope", resourceType))
		}
	}

	if len(issues) > 0 {
		t.Errorf("Found %d issues in resource definitions:", len(issues))
		for _, issue := range issues {
			t.Errorf("  - %s", issue)
		}
	} else {
		t.Logf("All %d resource definitions are complete", len(ResourceDefinitions))
	}
}
