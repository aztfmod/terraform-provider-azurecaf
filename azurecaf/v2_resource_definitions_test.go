package azurecaf

import (
	"encoding/json"
	"regexp"
	"testing"
)

func TestV2ResourceDefinitions(t *testing.T) {
	tests := []struct {
		resourceType string
		slug         string
		minLength    int
		maxLength    int
		scope        string
		regex        string
		validation   string
		dashes       bool
		lowercase    bool
	}{
		{"azurerm_custom_provider", "prov", 3, 64, "resourceGroup", "[^A-Za-z0-9]", "^[A-Za-z0-9]{3,64}$", false, false},
		{"azurerm_frontdoor", "afd", 5, 64, "global", "[^0-9A-Za-z-]", "^[a-zA-Z0-9][a-zA-Z0-9-]{3,62}[a-zA-Z0-9]$", true, false},
		{"azurerm_private_link_service", "pl", 2, 64, "resourceGroup", "[^0-9A-Za-z_.-]", "^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,62}[a-zA-Z0-9_]$", true, false},
		{"azurerm_servicebus_namespace", "sbns", 6, 50, "global", "[^0-9A-Za-z-]", "^[a-zA-Z][a-zA-Z0-9-]{4,48}[a-zA-Z0-9]$", true, false},
		{"azurerm_storage_table", "stt", 3, 63, "parent", "[^A-Za-z0-9]", "^[A-Za-z][A-Za-z0-9]{2,62}$", false, false},
		{"azurerm_synapse_workspace", "synw", 1, 50, "global", "[^0-9a-z-]", "^[a-z0-9](?:[a-z0-9-]{0,48}[a-z0-9])?$", true, true},
		{"azurerm_web_application_firewall_policy", "waf", 1, 80, "resourceGroup", "[^0-9A-Za-z_.-]", "^[a-zA-Z0-9](?:[a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_])?$", true, false},
		{"azurerm_lb_backend_address_pool", "lbbepool", 1, 80, "parent", "[^0-9A-Za-z_.-]", "^[a-zA-Z0-9](?:[a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_])?$", true, false},
		{"azurerm_lb_nat_pool", "lbnatpool", 1, 80, "parent", "[^0-9A-Za-z_.-]", "^[a-zA-Z0-9](?:[a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_])?$", true, false},
		{"azurerm_lb_outbound_rule", "lbor", 1, 80, "parent", "[^0-9A-Za-z_.-]", "^[a-zA-Z0-9](?:[a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_])?$", true, false},
		{"azurerm_lb_probe", "probe", 1, 80, "parent", "[^0-9A-Za-z_.-]", "^[a-zA-Z0-9](?:[a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_])?$", true, false},
		{"azurerm_lb_rule", "rule", 1, 80, "parent", "[^0-9A-Za-z_.-]", "^[a-zA-Z0-9](?:[a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_])?$", true, false},
		{"azurerm_load_test", "load", 1, 64, "global", "[^0-9A-Za-z_-]", "^[a-zA-Z](?:[a-zA-Z0-9_-]{0,62}[a-zA-Z0-9])?$", true, false},
	}

	for _, test := range tests {
		t.Run(test.resourceType, func(t *testing.T) {
			definition, exists := ResourceDefinitions[test.resourceType]
			if !exists {
				t.Fatalf("resource definition %q does not exist", test.resourceType)
			}

			if definition.CafPrefix != test.slug ||
				definition.MinLength != test.minLength ||
				definition.MaxLength != test.maxLength ||
				definition.Scope != test.scope ||
				definition.RegEx != test.regex ||
				definition.ValidationRegExp != test.validation ||
				definition.Dashes != test.dashes ||
				definition.LowerCase != test.lowercase {
				t.Errorf("unexpected definition for %s: %#v", test.resourceType, definition)
			}
		})
	}

	for _, removed := range []string{"azurerm_app_service_plan", "azurerm_lb_backend_pool", "azurerm_role_assignment"} {
		t.Run("removed_"+removed, func(t *testing.T) {
			if _, exists := ResourceDefinitions[removed]; exists {
				t.Errorf("removed resource definition %q is still registered", removed)
			}
		})
	}
}

func TestV2OneCharacterMinimumValidation(t *testing.T) {
	resourceTypes := []string{
		"azurerm_api_management",
		"azurerm_api_management_api",
		"azurerm_api_management_api_operation_tag",
		"azurerm_api_management_api_version_set",
		"azurerm_api_management_authorization_server",
		"azurerm_api_management_backend",
		"azurerm_api_management_certificate",
		"azurerm_api_management_gateway",
		"azurerm_api_management_group",
		"azurerm_api_management_logger",
		"azurerm_api_management_named_value",
		"azurerm_api_management_openid_connect_provider",
		"azurerm_app_service_certificate_order",
		"azurerm_application_insights_analytics_item",
		"azurerm_application_insights_api_key",
		"azurerm_automation_connection_certificate",
		"azurerm_automation_connection_classic_certificate",
		"azurerm_automation_connection_service_principal",
		"azurerm_backup_policy_file_share",
		"azurerm_backup_policy_vm",
		"azurerm_data_factory_linked_service_azure_file_storage",
		"azurerm_dns_srv_record",
		"azurerm_load_test",
		"azurerm_log_analytics_datasource_windows_event",
		"azurerm_log_analytics_datasource_windows_performance_counter",
		"azurerm_network_packet_capture",
		"azurerm_network_watcher_flow_log",
		"azurerm_site_recovery_network_mapping",
		"azurerm_site_recovery_protection_container_mapping",
	}

	for _, resourceType := range resourceTypes {
		t.Run(resourceType, func(t *testing.T) {
			definition := ResourceDefinitions[resourceType]
			if definition.MinLength != 1 {
				t.Fatalf("expected minimum length 1, got %d", definition.MinLength)
			}
			validation := regexp.MustCompile(definition.ValidationRegExp)
			if !validation.MatchString("a") {
				t.Fatalf("%q rejects a valid one-character name", definition.ValidationRegExp)
			}
		})
	}
}

func TestV2APIManagementValidationRejectsTrailingPipe(t *testing.T) {
	resourceTypes := []string{
		"azurerm_api_management_api_version_set",
		"azurerm_api_management_authorization_server",
		"azurerm_api_management_named_value",
		"azurerm_api_management_openid_connect_provider",
	}

	for _, resourceType := range resourceTypes {
		t.Run(resourceType, func(t *testing.T) {
			validation := regexp.MustCompile(ResourceDefinitions[resourceType].ValidationRegExp)
			if validation.MatchString("a|") {
				t.Fatalf("%q accepts a trailing pipe", validation.String())
			}
		})
	}
}

func TestResourceStructureUnmarshalsValidationRegex(t *testing.T) {
	var definition ResourceStructure
	if err := json.Unmarshal([]byte(`{"name":"test","validation_regex":"^[a-z]+$"}`), &definition); err != nil {
		t.Fatalf("failed to unmarshal resource definition: %v", err)
	}
	if definition.ValidationRegExp != "^[a-z]+$" {
		t.Errorf("unexpected validation regex: %q", definition.ValidationRegExp)
	}
}
