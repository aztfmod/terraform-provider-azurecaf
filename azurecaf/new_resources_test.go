package azurecaf

import (
	"regexp"
	"testing"
)

var newResourceNames = []string{
	"azurerm_api_management_api_version_set",
	"azurerm_api_management_authorization_server",
	"azurerm_api_management_named_value",
	"azurerm_api_management_openid_connect_provider",
	"azurerm_app_service_certificate_order",
	"azurerm_application_insights_analytics_item",
	"azurerm_application_insights_api_key",
	"azurerm_automation_connection_certificate",
	"azurerm_automation_connection_classic_certificate",
	"azurerm_automation_connection_service_principal",
	"azurerm_automation_dsc_nodeconfiguration",
	"azurerm_automation_variable_bool",
	"azurerm_automation_variable_datetime",
	"azurerm_automation_variable_int",
	"azurerm_automation_variable_string",
	"azurerm_backup_policy_file_share",
	"azurerm_backup_policy_vm",
	"azurerm_data_factory_linked_service_azure_file_storage",

	"azurerm_dns_srv_record",
	"azurerm_eventgrid_system_topic",
	"azurerm_key_vault_certificate_issuer",
	"azurerm_kusto_cluster_principal_assignment",
	"azurerm_kusto_database_principal_assignment",
	"azurerm_log_analytics_datasource_windows_event",
	"azurerm_log_analytics_datasource_windows_performance_counter",
	"azurerm_network_packet_capture",
	"azurerm_network_watcher_flow_log",
	"azurerm_site_recovery_network_mapping",
	"azurerm_site_recovery_protection_container_mapping",
}

var newResourceTestPrecedence = []string{"name", "slug", "random", "suffixes", "prefixes"}

func TestNewResources_NameGeneration(t *testing.T) {
	for _, resName := range newResourceNames {
		t.Run(resName, func(t *testing.T) {
			resource, ok := ResourceDefinitions[resName]
			if !ok {
				t.Fatalf("resource %s not found in ResourceDefinitions", resName)
			}

			// Generate a name using same params as TestAllResourceTypes
			name, err := getResourceName(resName, "-", []string{"dev"}, "app", []string{"001"}, "", ConventionCafClassic, true, false, true, newResourceTestPrecedence, false)
			if err != nil {
				t.Fatalf("failed to generate name: %v", err)
			}

			if name == "" {
				t.Fatal("generated empty name")
			}

			// Validate against regex
			validationRegex := resource.ValidationRegExp
			if validationRegex != "" {
				re, err := regexp.Compile(validationRegex)
				if err != nil {
					t.Fatalf("invalid validation regex %q: %v", validationRegex, err)
				}
				if !re.MatchString(name) {
					t.Errorf("generated name %q does not match validation regex %q", name, validationRegex)
				}
			}

			t.Logf("OK: %s -> %q (slug=%q, max=%d)", resName, name, resource.CafPrefix, resource.MaxLength)
		})
	}
}

func TestNewResources_Conventions(t *testing.T) {
	conventions := []struct {
		name       string
		convention string
	}{
		{"cafclassic", ConventionCafClassic},
		{"cafrandom", ConventionCafRandom},
		{"random", ConventionRandom},
		{"passthrough", ConventionPassThrough},
	}

	for _, conv := range conventions {
		t.Run(conv.name, func(t *testing.T) {
			for _, resName := range newResourceNames {
				t.Run(resName, func(t *testing.T) {
					passthrough := conv.convention == ConventionPassThrough
					useSlug := !passthrough
					_, err := getResourceName(resName, "-", []string{"dev"}, "app", []string{"001"}, "", conv.convention, true, passthrough, useSlug, newResourceTestPrecedence, false)
					if err != nil {
						t.Errorf("convention %s failed: %v", conv.name, err)
					}
				})
			}
		})
	}
}
