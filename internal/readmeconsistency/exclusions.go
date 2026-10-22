package readmeconsistency

// Excluded lists resourceDefinition.json entries that intentionally have no row
// in the README resource status table, because they are not AzureRM schema
// resources (compatibility aliases, sub-resource naming helpers, or generic
// definitions). Add an entry here only with a reason; stale entries fail the check.
var Excluded = map[string]string{
	"aks_node_pool_linux":                           "AKS node pool naming helper (Linux)",
	"aks_node_pool_windows":                         "AKS node pool naming helper (Windows)",
	"azurerm_api_management_service":                "compatibility alias of azurerm_api_management",
	"azurerm_automation_variable":                   "generic helper for the typed azurerm_automation_variable_* resources",
	"azurerm_bot_channel_Email":                     "legacy casing alias of azurerm_bot_channel_email",
	"azurerm_cognitive_deployment":                  "naming helper for Cognitive Services deployments",
	"azurerm_containerGroups":                       "legacy alias of azurerm_container_group",
	"azurerm_firewall_ip_configuration":             "sub-resource naming helper (firewall IP configuration block)",
	"azurerm_network_security_group_rule":           "legacy alias of azurerm_network_security_rule",
	"azurerm_private_dns_zone_group":                "sub-resource naming helper (private endpoint DNS zone group block)",
	"azurerm_private_service_connection":            "sub-resource naming helper (private endpoint connection block)",
	"azurerm_recovery_services_vault_backup_police": "legacy misspelled alias of a backup policy",
	"azurerm_snapshots":                             "legacy alias of azurerm_snapshot",
	"azurerm_virtual_machine_portal_name":           "naming helper for the VM portal display name",
	"azurerm_vm_windows_computer_name_prefix":       "naming helper for Windows computer name prefixes",
	"databricks_cluster":                            "Databricks provider resource, not AzureRM",
	"databricks_high_concurrency_cluster":           "Databricks provider resource, not AzureRM",
	"databricks_standard_cluster":                   "Databricks provider resource, not AzureRM",
	"general":                                       "generic naming definition",
	"general_safe":                                  "generic naming definition (safe characters)",
}
