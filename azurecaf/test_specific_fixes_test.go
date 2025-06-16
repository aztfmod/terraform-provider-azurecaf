package azurecaf

import (
	"fmt"
	"testing"
)

// TestSpecificSlugFixes validates that the specific issues mentioned in the GitHub issue are fixed
func TestSpecificSlugFixes(t *testing.T) {
	// Test the specific resources mentioned in the issue
	testCases := map[string]string{
		"azurerm_monitor_activity_log_alert": "amala",  // Was "adfmysql", now should be "amala"
		"azurerm_lb_rule":                    "lbr",    // Was "adt", now should be "lbr"
		"azurerm_lb_backend_pool":            "lbbp",   // Was "adt", now should be "lbbp"
		"azurerm_lb_backend_address_pool":    "lbbap",  // Was "adt", now should be "lbbap"
		"azurerm_lb_nat_pool":                "lbnp",   // Was "adt", now should be "lbnp"
		"azurerm_lb_outbound_rule":           "lbor",   // Was "adt", now should be "lbor"
		"azurerm_lb_probe":                   "lbp",    // Was "adt", now should be "lbp"
		// Verify that these still have their correct slugs
		"azurerm_digital_twins_instance":     "adt",    // Should remain "adt"
		"azurerm_data_factory_dataset_mysql": "adfmysql", // Should remain "adfmysql"
	}

	fmt.Println("\nVerifying slug fixes:")
	fmt.Println("Resource Type -> Expected Slug -> Actual Slug -> Status")
	fmt.Println("=======================================================")
	
	for resourceType, expectedSlug := range testCases {
		resource, err := getResource(resourceType)
		if err != nil {
			t.Errorf("Error getting resource %s: %v", resourceType, err)
			continue
		}
		
		status := "✓ PASS"
		if resource.CafPrefix != expectedSlug {
			status = "✗ FAIL"
			t.Errorf("Resource %s has incorrect slug. Expected: %s, Got: %s", 
				resourceType, expectedSlug, resource.CafPrefix)
		}
		
		fmt.Printf("%-35s -> %-10s -> %-10s -> %s\n", 
			resourceType, expectedSlug, resource.CafPrefix, status)
	}
}