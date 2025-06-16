package azurecaf

import (
	"testing"
)

// TestResourceSlugConsistency validates that resource slugs are appropriate for their resource types
func TestResourceSlugConsistency(t *testing.T) {
	// Test monitor resources have appropriate slugs
	monitorResources := map[string]string{
		"azurerm_monitor_action_group":                    "amag",
		"azurerm_monitor_activity_log_alert":              "amala", // Should be "amala", not "adfmysql"
		"azurerm_monitor_autoscale_setting":               "amas",
		"azurerm_monitor_data_collection_endpoint":        "dce",
		"azurerm_monitor_data_collection_rule":            "dcr",
		"azurerm_monitor_diagnostic_setting":              "amds",
		"azurerm_monitor_metric_alert":                     "ma",
		"azurerm_monitor_private_link_scope":               "ampls",
		"azurerm_monitor_scheduled_query_rules_alert":     "schqra",
	}

	for resourceType, expectedSlug := range monitorResources {
		resource, err := getResource(resourceType)
		if err != nil {
			t.Errorf("Failed to get resource %s: %v", resourceType, err)
			continue
		}
		if resource.CafPrefix != expectedSlug {
			t.Errorf("Resource %s has incorrect slug. Expected: %s, Got: %s", 
				resourceType, expectedSlug, resource.CafPrefix)
		}
	}

	// Test load balancer resources have appropriate slugs
	lbResources := map[string]string{
		"azurerm_lb":                         "lb",
		"azurerm_lb_nat_rule":                "lbnatrl", // Already correct
		"azurerm_lb_backend_pool":            "lbbp",    // Should be "lbbp", not "adt"
		"azurerm_lb_backend_address_pool":    "lbbap",   // Should be "lbbap", not "adt"
		"azurerm_lb_nat_pool":                "lbnp",    // Should be "lbnp", not "adt"
		"azurerm_lb_outbound_rule":           "lbor",    // Should be "lbor", not "adt"
		"azurerm_lb_probe":                   "lbp",     // Should be "lbp", not "adt"
		"azurerm_lb_rule":                    "lbr",     // Should be "lbr", not "adt"
	}

	for resourceType, expectedSlug := range lbResources {
		resource, err := getResource(resourceType)
		if err != nil {
			t.Errorf("Failed to get resource %s: %v", resourceType, err)
			continue
		}
		if resource.CafPrefix != expectedSlug {
			t.Errorf("Resource %s has incorrect slug. Expected: %s, Got: %s", 
				resourceType, expectedSlug, resource.CafPrefix)
		}
	}

	// Test that Azure Digital Twins is the only resource using "adt"
	resource, err := getResource("azurerm_digital_twins_instance")
	if err != nil {
		t.Errorf("Failed to get azurerm_digital_twins_instance: %v", err)
	} else if resource.CafPrefix != "adt" {
		t.Errorf("azurerm_digital_twins_instance should have slug 'adt', got: %s", resource.CafPrefix)
	}
}

// TestSlugUniqueness validates that slugs are not duplicated inappropriately
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
	
	// Check for specific problematic duplicates
	if resources, exists := slugToResources["adt"]; exists && len(resources) > 1 {
		t.Logf("Slug 'adt' is used by multiple resources: %v", resources)
		for _, resource := range resources {
			if resource != "azurerm_digital_twins_instance" {
				t.Errorf("Resource %s incorrectly uses slug 'adt' which should be reserved for Azure Digital Twins", resource)
			}
		}
	}
	
	if resources, exists := slugToResources["adfmysql"]; exists && len(resources) > 1 {
		t.Logf("Slug 'adfmysql' is used by multiple resources: %v", resources)
		for _, resource := range resources {
			if resource == "azurerm_monitor_activity_log_alert" {
				t.Errorf("Resource %s incorrectly uses slug 'adfmysql' which should be reserved for Data Factory MySQL dataset", resource)
			}
		}
	}
}