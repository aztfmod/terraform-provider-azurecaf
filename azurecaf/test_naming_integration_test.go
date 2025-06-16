package azurecaf

import (
	"strings"
	"testing"
)

// TestFixedSlugNaming validates that the fixed slugs work correctly in actual resource naming
func TestFixedSlugNaming(t *testing.T) {
	testCases := []struct {
		resourceType  string
		expectedSlug  string
		name          string
	}{
		{"azurerm_monitor_activity_log_alert", "amala", "test"},
		{"azurerm_lb_rule", "lbr", "test"},
		{"azurerm_lb_backend_pool", "lbbp", "test"},
		{"azurerm_digital_twins_instance", "adt", "test"},
	}

	for _, tc := range testCases {
		t.Run(tc.resourceType, func(t *testing.T) {
			namePrecedence := []string{"name", "slug", "random", "suffixes", "prefixes"}
			
			// Test with slug enabled
			resourceName, err := getResourceName(tc.resourceType, "-", nil, tc.name, nil, "", "cafclassic", true, false, true, namePrecedence)
			if err != nil {
				t.Errorf("getResourceName failed for %s: %v", tc.resourceType, err)
				return
			}
			
			// The generated name should contain the expected slug
			if !strings.Contains(resourceName, tc.expectedSlug) {
				t.Errorf("Resource name '%s' for %s should contain slug '%s'", resourceName, tc.resourceType, tc.expectedSlug)
			}
			
			t.Logf("✓ %s generates name '%s' containing slug '%s'", tc.resourceType, resourceName, tc.expectedSlug)
		})
	}
}