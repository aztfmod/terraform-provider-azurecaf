#!/usr/bin/env python3
"""Generate per-resource ``terraform test`` workspaces that validate every
CAF-generated name against the corresponding ``azurerm_*`` resource schema
using ``mock_provider azurerm``.

For each resource type chosen by ``--all`` / ``--diff-against`` / ``--resources``,
the script emits a directory under ``--out-dir`` containing:

* ``main.tf`` — three ``azurecaf_name`` variants (``default``, ``with_prefix``,
  ``with_random``) wired into three clones of the ``azurerm_*`` resource with
  every required attribute populated by deterministic fake values.
* ``terraform.rc`` — ``dev_overrides`` pointing at the locally built
  ``terraform-provider-azurecaf`` binary.
* ``tests/validate.tftest.hcl`` — five assertions per variant:
  output non-empty, length within ``min_length``/``max_length``, the
  CAF ``validation_regex`` matches, and the ``azurerm_*`` resource's ``name``
  attribute equals the CAF result.

Designed to be run from CI under ``mock_provider "azurerm" {}`` so no Azure
credentials are required.
"""
from __future__ import annotations

import argparse
import json
import shutil
import subprocess
import sys
from pathlib import Path

# ---------------------------------------------------------------------------
# Deterministic fake values used to satisfy required azurerm attributes that
# do not relate to the CAF name. The goal is to keep the azurerm provider
# happy at plan/apply time so the mock_provider can complete and the assertions
# on the CAF-generated `name` actually fire.
# ---------------------------------------------------------------------------

FAKE_RG = "rg-test"
FAKE_LOCATION = "westeurope"
FAKE_SUB_ID = "00000000-0000-0000-0000-000000000000"
FAKE_TENANT_ID = "11111111-1111-1111-1111-111111111111"
FAKE_OBJ_ID = "22222222-2222-2222-2222-222222222222"
RG_SCOPE = f"/subscriptions/{FAKE_SUB_ID}/resourceGroups/{FAKE_RG}"

# Repeated HCL literals lifted to constants to keep the override tables readable
# and to avoid magic-string duplication. Values are raw HCL literals (including
# surrounding quotes for strings) so they can be inlined verbatim into rendered
# terraform files.
HCL_LINUX = '"Linux"'
HCL_CAN_NOT_DELETE = '"CanNotDelete"'
HCL_CUSTOM = '"Custom"'
HCL_PRIVATE_IP = '"10.0.0.1"'  # documentation-only address used in mock plans
HCL_PT5M = '"PT5M"'
HCL_TEST_STR = '"test"'

# Per-attribute fake-value overrides. azurerm 4.x validates many string
# attributes against enums/regex at plan time, and the literal ``"test"`` is
# typically rejected. Each entry is the raw HCL literal (including quotes).
ATTR_OVERRIDES: dict[str, str] = {
    "os_type":                         HCL_LINUX,
    "kind":                            '"ServiceCatalog"',
    "lock_level":                      HCL_CAN_NOT_DELETE,
    "evaluator_type":                  '"AllowedValuesPolicy"',
    "time_zone_id":                    '"UTC"',
    "create_option":                   '"Empty"',
    "severity":                        '"Medium"',
    "product_filter":                  '"Azure Security Center"',
    "detector_type":                   '"FailureAnomaliesDetector"',
    "security_provider_name":          '"ZScaler"',
    "policy_type":                     HCL_CUSTOM,
    "mode":                            '"All"',
    "resource":                        '"directory"',
    "source":                          '"Microsoft.KeyVault"',
    "type":                            '"IPsec"',
    "priority":                        '100',
    "cache_size_in_gb":                '3072',
    "peer_ip":                         HCL_PRIVATE_IP,
    "ip_address":                      HCL_PRIVATE_IP,
    "namespace_path":                  '"/test"',
    "target_path":                     '"/test"',
    "sku_name":                        '"Standard_2G"',
    "automation_account_name":         '"aa-test"',
    "iothub_name":                     '"iot-test"',
    "key_vault_name":                  '"kv-test"',
    "endpoint_uri":                    '"sb://test.servicebus.windows.net/"',
    "entity_path":                     '"queue-test"',
    "connection_string":               '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k="',
    "scope_kind":                      '"resourceGroup"',
}

# Per-resource per-attribute overrides (most specific wins). Add an entry here
# whenever a new resource fails because azurerm's CustomizeDiff rejects the
# generic placeholder. Keep entries minimal — only the required attributes the
# generic logic cannot infer.
RESOURCE_ATTR_OVERRIDES: dict[str, dict[str, str]] = {
    "azurerm_storage_account": {
        "account_tier":             '"Standard"',
        "account_replication_type": '"LRS"',
    },
    "azurerm_key_vault": {
        "sku_name": '"standard"',
        "tenant_id": f'"{FAKE_TENANT_ID}"',
    },
    "azurerm_container_group":                  {"os_type": HCL_LINUX},
    "azurerm_service_plan":                     {"os_type": HCL_LINUX, "sku_name": '"B1"'},
    "azurerm_dedicated_hardware_security_module": {"sku_name": '"SafeNet Luna Network HSM A790"'},
    "azurerm_eventhub_cluster":                 {"sku_name": '"Dedicated_1"'},
    "azurerm_hpc_cache":                        {"cache_size_in_gb": '3072', "sku_name": '"Standard_2G"'},
    "azurerm_managed_application":              {"kind": '"ServiceCatalog"'},
    "azurerm_managed_application_definition":   {"lock_level": HCL_CAN_NOT_DELETE},
    "azurerm_management_lock":                  {"lock_level": HCL_CAN_NOT_DELETE},
    "azurerm_snapshot":                         {"create_option": '"Copy"'},
    "azurerm_storage_data_lake_gen2_path":      {"resource": '"directory"'},
    "azurerm_storage_encryption_scope":         {"source": '"Microsoft.Storage"'},
    "azurerm_sentinel_alert_rule_scheduled":    {
        "severity": '"Medium"', "query": '"SecurityEvent | take 1"',
        "query_frequency": HCL_PT5M, "query_period": HCL_PT5M,
        "trigger_operator": '"GreaterThan"', "trigger_threshold": '0',
    },
    "azurerm_sentinel_alert_rule_ms_security_incident": {
        "product_filter": '"Azure Security Center"',
        "severity_filter": '["Medium"]',
    },
    "azurerm_monitor_smart_detector_alert_rule": {
        "detector_type": '"FailureAnomaliesDetector"',
        "severity": '"Sev1"', "frequency": HCL_PT5M,
    },
    "azurerm_virtual_hub_security_partner_provider": {"security_provider_name": '"ZScaler"'},
    "azurerm_virtual_network_gateway_connection": {"type": '"IPsec"'},
    "azurerm_policy_definition": {
        "mode": '"All"', "policy_type": HCL_CUSTOM,
        "policy_rule": '"{\\"if\\":{\\"field\\":\\"type\\",\\"equals\\":\\"Microsoft.Storage/storageAccounts\\"},\\"then\\":{\\"effect\\":\\"audit\\"}}"',
    },
    "azurerm_policy_set_definition":            {"policy_type": HCL_CUSTOM},
    "azurerm_dev_test_policy": {
        "evaluator_type": '"AllowedValuesPolicy"', "threshold": '"1"',
        "fact_data": '""', "fact_name": '"UserOwnedLabVmCount"',
    },
    "azurerm_api_management_api_version_set": {
        "versioning_scheme": '"Segment"',
    },
    "azurerm_api_management_authorization_server": {
        "authorization_methods": '["GET"]',
        "grant_types": '["authorizationCode"]',
        "authorization_endpoint": '"https://login.microsoftonline.com/common/oauth2/authorize"',
        "client_id": '"00000000-0000-0000-0000-000000000000"',
    },
    "azurerm_api_management_named_value": {
        "value": '"test-value"',
    },
    "azurerm_application_insights_analytics_item": {
        "scope": '"shared"',
        "type": '"query"',
        "content": '"requests | take 10"',
        "application_insights_id": f'"{RG_SCOPE}/providers/Microsoft.Insights/components/appi-test"',
    },
    "azurerm_application_insights_api_key": {
        "read_permissions": '["api"]',
        "application_insights_id": f'"{RG_SCOPE}/providers/Microsoft.Insights/components/appi-test"',
    },
    "azurerm_automation_connection_service_principal": {
        "application_id": f'"{FAKE_TENANT_ID}"',
        "tenant_id": f'"{FAKE_TENANT_ID}"',
        "certificate_thumbprint": '"0000000000000000000000000000000000000000"',
        "subscription_id": f'"{FAKE_SUB_ID}"',
    },
    "azurerm_backup_policy_file_share": {
        "timezone": '"UTC"',
        "frequency": '"Daily"',
        "time": '"23:00"',
    },
    "azurerm_backup_policy_vm": {
        "timezone": '"UTC"',
        "frequency": '"Daily"',
        "time": '"23:00"',
    },
    "azurerm_key_vault_certificate_issuer": {
        "provider_name": '"DigiCert"',
    },
    "azurerm_kusto_cluster_principal_assignment": {
        "principal_type": '"App"',
        "role": '"AllDatabasesAdmin"',
        "tenant_id": f'"{FAKE_TENANT_ID}"',
        "principal_id": f'"{FAKE_TENANT_ID}"',
    },
    "azurerm_kusto_database_principal_assignment": {
        "principal_type": '"App"',
        "role": '"Viewer"',
        "tenant_id": f'"{FAKE_TENANT_ID}"',
        "principal_id": f'"{FAKE_TENANT_ID}"',
    },
    "azurerm_log_analytics_datasource_windows_event": {
        "event_log_name": '"Application"',
        "event_types": '["Error", "Warning"]',
    },
    "azurerm_log_analytics_datasource_windows_performance_counter": {
        "object_name": '"Processor"',
        "instance_name": '"*"',
        "counter_name": '"% Processor Time"',
        "interval_seconds": '60',
    },
    "azurerm_network_packet_capture": {
        "network_watcher_name": '"nw-test"',
        "target_resource_id": f'"{RG_SCOPE}/providers/Microsoft.Compute/virtualMachines/vm-test"',
        "storage_account_id": f'"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest"',
    },
    "azurerm_network_watcher_flow_log": {
        "network_watcher_name": '"nw-test"',
        "network_security_group_id": f'"{RG_SCOPE}/providers/Microsoft.Network/networkSecurityGroups/nsg-test"',
        "storage_account_id": f'"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest"',
        "enabled": 'true',
    },
    "azurerm_dev_test_schedule": {
        "time_zone_id": '"UTC"', "task_type": '"LabVmsShutdownTask"',
    },
    "azurerm_iothub_route": {
        "source": '"DeviceMessages"', "condition": '"true"',
        "endpoint_names": '["endpoint1"]', "enabled": 'true',
    },
    "azurerm_advanced_threat_protection": {
        "enabled": 'true',
        "target_resource_id": f'"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest"',
    },
    "azurerm_express_route_circuit_peering": {
        "vlan_id": '100', "peering_type": '"AzurePrivatePeering"',
        "primary_peer_address_prefix": '"10.0.0.0/30"',
        "secondary_peer_address_prefix": '"10.0.0.4/30"',
        "peer_asn": '100',
    },
    "azurerm_virtual_hub_bgp_connection":       {"peer_ip": HCL_PRIVATE_IP, "peer_asn": '65515'},
    "azurerm_firewall_policy_rule_collection_group": {"priority": '500'},
    "azurerm_data_share":                       {"kind": '"CopyBased"'},
    "azurerm_data_share_dataset_blob_storage":  {"file_path": '"file.txt"', "container_name": '"ct-test"'},
    "azurerm_resource_group_template_deployment": {
        "deployment_mode": '"Incremental"',
        "template_content": '"{\\"$schema\\":\\"https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#\\",\\"contentVersion\\":\\"1.0.0.0\\",\\"resources\\":[]}"',
    },
    "azurerm_subscription_template_deployment": {
        "location": '"westeurope"',
        "template_content": '"{\\"$schema\\":\\"https://schema.management.azure.com/schemas/2019-08-01/subscriptionDeploymentTemplate.json#\\",\\"contentVersion\\":\\"1.0.0.0\\",\\"resources\\":[]}"',
    },
    "azurerm_shared_image_version": {
        "blob_uri": '"https://sttest.blob.core.windows.net/vhds/disk.vhd"',
        "storage_account_id": f'"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest"',
    },
    "azurerm_kusto_attached_database_configuration": {
        "default_principals_modification_kind": '"Union"',
    },
    "azurerm_log_analytics_linked_service": {
        "read_access_id": f'"{RG_SCOPE}/providers/Microsoft.OperationalInsights/clusters/lac-test"',
    },
    "azurerm_app_service_certificate": {
        "pfx_blob": '""', "password": '""',
        "key_vault_secret_id": '"https://kv-test.vault.azure.net/secrets/cert/00000000000000000000000000000000"',
    },
    "azurerm_spring_cloud_certificate": {
        "key_vault_certificate_id": '"https://kv-test.vault.azure.net/certificates/cert/00000000000000000000000000000000"',
    },
    "azurerm_iothub_endpoint_eventhub": {
        "connection_string": '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k=;EntityPath=eh"',
        "endpoint_uri": '"sb://test.servicebus.windows.net"',
        "entity_path": '"eh-test"', "authentication_type": '"keyBased"',
        "iothub_id": f'"{RG_SCOPE}/providers/Microsoft.Devices/IotHubs/iot-test"',
    },
    "azurerm_iothub_endpoint_servicebus_queue": {
        "connection_string": '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k=;EntityPath=q"',
        "iothub_id": f'"{RG_SCOPE}/providers/Microsoft.Devices/IotHubs/iot-test"',
    },
    "azurerm_iothub_endpoint_servicebus_topic": {
        "connection_string": '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k=;EntityPath=t"',
        "iothub_id": f'"{RG_SCOPE}/providers/Microsoft.Devices/IotHubs/iot-test"',
    },
    "azurerm_iothub_endpoint_storage_container": {
        "connection_string": '"DefaultEndpointsProtocol=https;AccountName=sttest;AccountKey=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==;EndpointSuffix=core.windows.net"',
        "container_name": '"ctr-test"',
        "iothub_id": f'"{RG_SCOPE}/providers/Microsoft.Devices/IotHubs/iot-test"',
    },
    "azurerm_iothub_fallback_route": {
        "source": '"DeviceMessages"', "endpoint_names": '["events"]',
        "condition": '"true"', "enabled": 'true',
    },
    "azurerm_log_analytics_cluster_customer_managed_key": {
        "key_vault_key_id": '"https://kv-test.vault.azure.net/keys/test/00000000000000000000000000000000"',
    },
    "azurerm_security_center_auto_provisioning":{"auto_provision": '"On"'},
    "azurerm_mssql_virtual_machine":            {"sql_license_type": '"PAYG"'},
    "azurerm_key_vault_access_policy": {
        "key_vault_id": f'"{RG_SCOPE}/providers/Microsoft.KeyVault/vaults/kv-test"',
    },
}

# Extra raw HCL lines to inject into resource blocks for constraints
# that cannot be expressed via attribute overrides alone.
RESOURCE_EXTRA_HCL: dict[str, list[str]] = {}

# Common parent-id / contextual attributes that map cleanly to a fake ARM id.
FAKE_IDS: dict[str, str] = {
    "key_vault_id":                    f"{RG_SCOPE}/providers/Microsoft.KeyVault/vaults/kv-test",
    "iothub_id":                       f"{RG_SCOPE}/providers/Microsoft.Devices/IotHubs/iot-test",
    "data_factory_id":                 f"{RG_SCOPE}/providers/Microsoft.DataFactory/factories/adf-test",
    "namespace_id":                    f"{RG_SCOPE}/providers/Microsoft.ServiceBus/namespaces/sb-test",
    "spring_cloud_service_id":         f"{RG_SCOPE}/providers/Microsoft.AppPlatform/spring/spc-test",
    "spring_cloud_app_id":             f"{RG_SCOPE}/providers/Microsoft.AppPlatform/spring/spc-test/apps/app-test",
    "lab_id":                          f"{RG_SCOPE}/providers/Microsoft.DevTestLab/labs/lab-test",
    "log_analytics_workspace_id":      f"{RG_SCOPE}/providers/Microsoft.OperationalInsights/workspaces/la-test",
    "workspace_id":                    f"{RG_SCOPE}/providers/Microsoft.OperationalInsights/workspaces/la-test",
    "workspace_resource_id":           f"{RG_SCOPE}/providers/Microsoft.OperationalInsights/workspaces/la-test",
    "log_analytics_cluster_id":        f"{RG_SCOPE}/providers/Microsoft.OperationalInsights/clusters/lac-test",
    "policy_definition_id":            "/providers/Microsoft.Authorization/policyDefinitions/00000000-0000-0000-0000-000000000000",
    "policy_assignment_id":            f"{RG_SCOPE}/providers/Microsoft.Authorization/policyAssignments/pa-test",
    "scope":                           RG_SCOPE,
    "target_resource_id":              f"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest",
    "target_resource_group_id":        RG_SCOPE,
    "automation_account_id":           f"{RG_SCOPE}/providers/Microsoft.Automation/automationAccounts/aa-test",
    "data_share_account_id":           f"{RG_SCOPE}/providers/Microsoft.DataShare/accounts/dsh-test",
    "share_id":                        f"{RG_SCOPE}/providers/Microsoft.DataShare/accounts/dsh-test/shares/dshr-test",
    "storage_account_id":              f"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest",
    "storage_account_resource_id":     f"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest",
    "storage_data_lake_gen2_filesystem_id": f"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest/blobServices/default/containers/fs-test",
    "image_id":                        f"{RG_SCOPE}/providers/Microsoft.Compute/galleries/gal/images/img",
    "managed_disk_id":                 f"{RG_SCOPE}/providers/Microsoft.Compute/disks/disk-test",
    "snapshot_id":                     f"{RG_SCOPE}/providers/Microsoft.Compute/snapshots/snap-test",
    "source_resource_id":              f"{RG_SCOPE}/providers/Microsoft.Compute/disks/disk-test",
    "platform_fault_domain_count":     "1",
    "subnet_id":                       f"{RG_SCOPE}/providers/Microsoft.Network/virtualNetworks/vnet-test/subnets/snet-test",
    "virtual_network_id":              f"{RG_SCOPE}/providers/Microsoft.Network/virtualNetworks/vnet-test",
    "virtual_hub_id":                  f"{RG_SCOPE}/providers/Microsoft.Network/virtualHubs/vh-test",
    "vpn_gateway_id":                  f"{RG_SCOPE}/providers/Microsoft.Network/vpnGateways/vpng-test",
    "express_route_circuit_id":        f"{RG_SCOPE}/providers/Microsoft.Network/expressRouteCircuits/erc-test",
    "express_route_circuit_peering_id": f"{RG_SCOPE}/providers/Microsoft.Network/expressRouteCircuits/erc-test/peerings/AzurePrivatePeering",
    "express_route_gateway_id":        f"{RG_SCOPE}/providers/Microsoft.Network/expressRouteGateways/erg-test",
    "virtual_network_gateway_id":      f"{RG_SCOPE}/providers/Microsoft.Network/virtualNetworkGateways/vng-test",
    "peer_virtual_network_gateway_id": f"{RG_SCOPE}/providers/Microsoft.Network/virtualNetworkGateways/vng-peer",
    "local_network_gateway_id":        f"{RG_SCOPE}/providers/Microsoft.Network/localNetworkGateways/lng-test",
    "firewall_policy_id":              f"{RG_SCOPE}/providers/Microsoft.Network/firewallPolicies/fwp-test",
    "route_filter_id":                 f"{RG_SCOPE}/providers/Microsoft.Network/routeFilters/rf-test",
    "network_profile_id":              f"{RG_SCOPE}/providers/Microsoft.Network/networkProfiles/np-test",
    "container_group_id":              f"{RG_SCOPE}/providers/Microsoft.ContainerInstance/containerGroups/aci-test",
    "cosmosdb_account_id":             f"{RG_SCOPE}/providers/Microsoft.DocumentDB/databaseAccounts/cosmos-test",
    "cassandra_keyspace_id":           f"{RG_SCOPE}/providers/Microsoft.DocumentDB/databaseAccounts/cosmos-test/cassandraKeyspaces/ks-test",
    "gremlin_database_id":             f"{RG_SCOPE}/providers/Microsoft.DocumentDB/databaseAccounts/cosmos-test/gremlinDatabases/db-test",
    "sql_database_id":                 f"{RG_SCOPE}/providers/Microsoft.DocumentDB/databaseAccounts/cosmos-test/sqlDatabases/db-test",
    "sql_container_id":                f"{RG_SCOPE}/providers/Microsoft.DocumentDB/databaseAccounts/cosmos-test/sqlDatabases/db-test/containers/c-test",
    "mongo_database_id":               f"{RG_SCOPE}/providers/Microsoft.DocumentDB/databaseAccounts/cosmos-test/mongodbDatabases/db-test",
    "managed_application_definition_id": f"{RG_SCOPE}/providers/Microsoft.Solutions/applicationDefinitions/mad-test",
    "recovery_vault_name":             "rsv-test",
    "recovery_fabric_name":            "fabric-test",
    "recovery_replication_policy_id":  f"{RG_SCOPE}/providers/Microsoft.RecoveryServices/vaults/rsv-test/replicationPolicies/rp-test",
    "source_recovery_fabric_id":       f"{RG_SCOPE}/providers/Microsoft.RecoveryServices/vaults/rsv-test/replicationFabrics/fab-test",
    "source_recovery_protection_container_id": f"{RG_SCOPE}/providers/Microsoft.RecoveryServices/vaults/rsv-test/replicationFabrics/fab-test/replicationProtectionContainers/pc-test",
    "target_recovery_fabric_id":       f"{RG_SCOPE}/providers/Microsoft.RecoveryServices/vaults/rsv-test/replicationFabrics/fab2-test",
    "target_recovery_protection_container_id": f"{RG_SCOPE}/providers/Microsoft.RecoveryServices/vaults/rsv-test/replicationFabrics/fab2-test/replicationProtectionContainers/pc2-test",
    "source_vm_id":                    f"{RG_SCOPE}/providers/Microsoft.Compute/virtualMachines/vm-test",
    "kusto_cluster_id":                f"{RG_SCOPE}/providers/Microsoft.Kusto/clusters/kc-test",
    "cluster_resource_id":             f"{RG_SCOPE}/providers/Microsoft.Kusto/clusters/kc-test",
    "database_name":                   "db-test",
    "kusto_database_id":               f"{RG_SCOPE}/providers/Microsoft.Kusto/clusters/kc-test/databases/db-test",
    "service_plan_id":                 f"{RG_SCOPE}/providers/Microsoft.Web/serverfarms/asp-test",
    "app_service_id":                  f"{RG_SCOPE}/providers/Microsoft.Web/sites/app-test",
    "app_service_name":                "app-test",
    "app_service_plan_id":             f"{RG_SCOPE}/providers/Microsoft.Web/serverfarms/asp-test",
    "linked_service_resource_id":      f"{RG_SCOPE}/providers/Microsoft.OperationalInsights/workspaces/la-test/linkedServices/Cluster",
    "hpc_cache_id":                    f"{RG_SCOPE}/providers/Microsoft.StorageCache/caches/hpcc-test",
    "target_host_name":                "nfs.test.local",
    "iothub_name":                     "iot-test",
    "endpoint_uri":                    "sb://test.servicebus.windows.net",
    "entity_path":                     "queue-test",
    "container_name":                  "ct-test",
    "connection_string":               "Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k=",
    "vault_id":                        f"{RG_SCOPE}/providers/Microsoft.KeyVault/vaults/kv-test",
    "vault_uri":                       "https://kv-test.vault.azure.net/",
    "key_vault_key_id":                "https://kv-test.vault.azure.net/keys/test/00000000000000000000000000000000",
    "addressing_family":               "IPv4",
    "address_prefix_cidr":             "10.0.0.0/24",
    "object_id":                       FAKE_OBJ_ID,
    "tenant_id":                       FAKE_TENANT_ID,
    "principal_id":                    FAKE_OBJ_ID,
    "license_type":                    "PAYG",
    "publisher_name":                  "publisher",
    "product_name":                    "product",
    "offer_name":                      "offer",
    "plan_name":                       "plan",
    "plan_version":                    "1.0.0",
    "managed_resource_group_name":     "rg-managed-test",
    "package_file_uri":                "https://example.com/package.zip",
    "function_app_id":                 f"{RG_SCOPE}/providers/Microsoft.Web/sites/func-test",
    "linux_function_app_id":           f"{RG_SCOPE}/providers/Microsoft.Web/sites/func-test",
}


# Named-attribute dispatch table: simple "attribute name -> HCL literal" rules
# whose value depends only on the attribute name. Lambdas keep callsites uniform
# with the suffix rules below.
_NAMED_ATTR_RULES: dict[str, "callable"] = {
    "resource_group_name": lambda: f'"{FAKE_RG}"',
    "location":            lambda: f'"{FAKE_LOCATION}"',
    "subscription_id":     lambda: f'"{FAKE_SUB_ID}"',
}

# Primitive-type dispatch table for scalar attribute types.
_PRIMITIVE_DEFAULTS: dict[str, str] = {
    "string": HCL_TEST_STR,
    "number": '1',
    "bool":   'true',
}

# Collection-inner-type dispatch table for set/list attribute types.
_COLLECTION_DEFAULTS: dict[str, str] = {
    "string": '["test"]',
    "number": '[1]',
    "bool":   '[true]',
}

_DEP_ID_LIT = f'"{RG_SCOPE}/providers/Microsoft.Resources/deployments/dep-test"'
_DEP_IDS_LIT = f'[{_DEP_ID_LIT}]'


def _override_lookup(attr_name: str, resource_type: str | None) -> str | None:
    """Return the most-specific HCL literal override for ``attr_name`` or None."""
    if resource_type and resource_type in RESOURCE_ATTR_OVERRIDES:
        ov = RESOURCE_ATTR_OVERRIDES[resource_type].get(attr_name)
        if ov is not None:
            return ov
    if attr_name in ATTR_OVERRIDES:
        return ATTR_OVERRIDES[attr_name]
    if attr_name in FAKE_IDS:
        return f'"{FAKE_IDS[attr_name]}"'
    return None


def _name_based_value(attr_name: str) -> str | None:
    """Return an HCL literal derived from ``attr_name`` conventions, or None."""
    rule = _NAMED_ATTR_RULES.get(attr_name)
    if rule is not None:
        return rule()
    if attr_name.endswith("_ids"):
        return _DEP_IDS_LIT
    if attr_name.endswith("_id"):
        return _DEP_ID_LIT
    return None


def _type_based_value(attr_type) -> str:
    """Return an HCL literal derived solely from a Terraform attribute type."""
    if isinstance(attr_type, list):
        head = attr_type[0]
        if head in ("set", "list"):
            inner = attr_type[1] if len(attr_type) > 1 else "string"
            return _COLLECTION_DEFAULTS.get(inner, '[]')
        if head in ("map", "object"):
            return '{}'
        return HCL_TEST_STR
    return _PRIMITIVE_DEFAULTS.get(attr_type, HCL_TEST_STR)


def fake_value_for(attr_name: str, attr_type, _attr_schema: dict, resource_type: str | None = None) -> str:
    """Return an HCL literal for a required attribute. Most-specific override wins."""
    override = _override_lookup(attr_name, resource_type)
    if override is not None:
        return override
    name_based = _name_based_value(attr_name)
    if name_based is not None:
        return name_based
    return _type_based_value(attr_type)


def render_block(block_name: str, block_def: dict, indent: int, resource_type: str) -> str:
    sp = " " * indent
    lines = [f"{sp}{block_name} {{"]
    block = block_def.get("block", {})
    for attr_name, attr_def in block.get("attributes", {}).items():
        if attr_def.get("required"):
            t = attr_def.get("type")
            lines.append(f"{sp}  {attr_name} = {fake_value_for(attr_name, t, attr_def, resource_type)}")
        elif _override_lookup(attr_name, resource_type) is not None:
            # Emit optional attrs that have explicit overrides (for "at least one of" constraints)
            lines.append(f"{sp}  {attr_name} = {_override_lookup(attr_name, resource_type)}")
    for bn, bdef in block.get("block_types", {}).items():
        if bdef.get("min_items", 0) > 0:
            lines.append(render_block(bn, bdef, indent + 2, resource_type))
    lines.append(f"{sp}}}")
    return "\n".join(lines)


def render_required_attrs_and_blocks(schema: dict, name_attr_name: str, resource_type: str) -> list[str]:
    block = schema.get("block", {})
    out: list[str] = []
    for attr_name, attr_def in block.get("attributes", {}).items():
        if attr_name == name_attr_name:
            continue
        if attr_def.get("required"):
            t = attr_def.get("type")
            out.append(f"  {attr_name} = {fake_value_for(attr_name, t, attr_def, resource_type)}")
        elif not attr_def.get("required") and resource_type in RESOURCE_ATTR_OVERRIDES:
            # Emit optional attrs with resource-specific overrides (for "at least one of" constraints)
            ov = RESOURCE_ATTR_OVERRIDES[resource_type].get(attr_name)
            if ov is not None:
                out.append(f"  {attr_name} = {ov}")
    for bn, bdef in block.get("block_types", {}).items():
        if bdef.get("min_items", 0) > 0:
            out.append(render_block(bn, bdef, indent=2, resource_type=resource_type))
    return out


def find_name_attr(schema: dict) -> str | None:
    block = schema.get("block", {})
    attrs = block.get("attributes", {})
    for cand in ("name", "display_name", "rule_name", "policy_name", "alert_rule_name"):
        if cand in attrs and attrs[cand].get("required"):
            return cand
    for k, v in attrs.items():
        if k.endswith("_name") and v.get("required"):
            return k
    return None


def make_main_tf(resource_type: str, schema: dict) -> str:
    name_attr = find_name_attr(schema) or "name"
    required_lines = render_required_attrs_and_blocks(schema, name_attr_name=name_attr, resource_type=resource_type)
    required_block = "\n".join(required_lines)

    # Inject extra HCL for "at least one of" constraints
    extra_lines = RESOURCE_EXTRA_HCL.get(resource_type, [])
    if extra_lines:
        required_block += "\n" + "\n".join(extra_lines)

    variants = [
        ("default",     '  name          = "test"'),
        ("with_prefix", '  name          = "test"\n  prefixes      = ["dev"]'),
        ("with_random", '  name          = "test"\n  random_length = 5\n  random_seed   = 12345'),
    ]

    parts = [f'''terraform {{
  required_providers {{
    azurecaf = {{ source = "aztfmod/azurecaf", version = ">= 1.2.0" }}
    azurerm  = {{ source = "hashicorp/azurerm",  version = ">= 4.0.0" }}
  }}
}}

provider "azurecaf" {{}}
provider "azurerm" {{
  features {{}}
  subscription_id = "{FAKE_SUB_ID}"
}}

''']
    for variant_name, naming_block in variants:
        parts.append(f'''resource "azurecaf_name" "{variant_name}" {{
{naming_block}
  resource_type = "{resource_type}"
  clean_input   = true
}}

resource "{resource_type}" "{variant_name}" {{
  {name_attr} = azurecaf_name.{variant_name}.result
{required_block}
}}

output "{variant_name}_result" {{
  value = azurecaf_name.{variant_name}.result
}}

output "{variant_name}_length" {{
  value = length(azurecaf_name.{variant_name}.result)
}}

''')
    return "".join(parts)


def make_test_hcl(resource_type: str, res_def: dict, name_attr: str = "name") -> str:
    min_l = res_def.get("min_length", 1)
    max_l = res_def.get("max_length", 80)
    regex = res_def.get("validation_regex", ".*")
    # Project convention stores regex wrapped in literal double quotes; strip them.
    if regex.startswith('"') and regex.endswith('"'):
        regex = regex[1:-1]
    regex_lit = json.dumps(regex)
    body = ['mock_provider "azurerm" {}\n\n']
    for variant in ("default", "with_prefix", "with_random"):
        body.append(f'''run "{variant}" {{
  command = apply

  assert {{
    condition     = output.{variant}_result != ""
    error_message = "Generated name output is empty for {variant}"
  }}
  assert {{
    condition     = output.{variant}_length >= {min_l}
    error_message = "Generated name shorter than min_length ({min_l}) for {variant}"
  }}
  assert {{
    condition     = output.{variant}_length <= {max_l}
    error_message = "Generated name exceeds max_length ({max_l}) for {variant}"
  }}
  assert {{
    condition     = can(regex({regex_lit}, output.{variant}_result))
    error_message = "Generated name does not match validation regex for {variant}"
  }}
  assert {{
    condition     = {resource_type}.{variant}.{name_attr} == output.{variant}_result
    error_message = "azurerm name does not equal CAF result for {variant}"
  }}
}}

''')
    return "".join(body)


def make_terraform_rc(plugin_dir: str) -> str:
    return f'''provider_installation {{
  dev_overrides {{
    "aztfmod/azurecaf" = "{plugin_dir}"
  }}
  direct {{}}
}}
'''


def load_resource_definitions(path: Path) -> dict[str, dict]:
    with path.open() as fp:
        return {r["name"]: r for r in json.load(fp)}


def load_schema(path: Path) -> dict[str, dict]:
    with path.open() as fp:
        doc = json.load(fp)
    return doc["provider_schemas"]["registry.terraform.io/hashicorp/azurerm"]["resource_schemas"]


def resources_changed_between(base_ref: str, res_def_path: Path, repo_root: Path) -> list[str]:
    """Return the resource ``name`` values added or modified in ``res_def_path``
    between ``base_ref`` and ``HEAD``.

    Compares the parsed JSON at both refs so the result is insensitive to
    whitespace/formatting changes and immune to unified-diff edge cases.
    """
    rel = res_def_path.relative_to(repo_root)
    base_blob = subprocess.run(
        ["git", "-C", str(repo_root), "show", f"{base_ref}:{rel}"],
        capture_output=True, text=True, check=False,
    ).stdout
    head_blob = subprocess.run(
        ["git", "-C", str(repo_root), "show", f"HEAD:{rel}"],
        capture_output=True, text=True, check=False,
    ).stdout

    base_by_name = {r["name"]: r for r in (json.loads(base_blob) if base_blob.strip() else [])
                    if isinstance(r, dict) and "name" in r}
    head_by_name = {r["name"]: r for r in (json.loads(head_blob) if head_blob.strip() else [])
                    if isinstance(r, dict) and "name" in r}

    changed = {name for name, entry in head_by_name.items()
               if base_by_name.get(name) != entry}
    return sorted(changed)


def _build_arg_parser() -> argparse.ArgumentParser:
    p = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    p.add_argument("--plugin-dir", required=True,
                   help="Local plugin directory containing the built terraform-provider-azurecaf binary "
                        "(e.g. ~/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/<os_arch>).")
    p.add_argument("--out-dir", required=True, type=Path,
                   help="Directory under which one workspace per resource is generated.")
    p.add_argument("--res-def", type=Path, default=Path("resourceDefinition.json"),
                   help="Path to resourceDefinition.json (default: ./resourceDefinition.json).")
    p.add_argument("--schema-file", type=Path, required=True,
                   help="Path to azurerm provider schema JSON produced by `terraform providers schema -json`.")
    sel = p.add_mutually_exclusive_group(required=True)
    sel.add_argument("--all", action="store_true", help="Generate for every resource in resourceDefinition.json.")
    sel.add_argument("--diff-against", metavar="BASE_REF",
                     help="Generate only for resources added/modified in resourceDefinition.json since BASE_REF.")
    sel.add_argument("--resources", metavar="CSV",
                     help="Comma-separated explicit list of resource type names.")
    p.add_argument("--repo-root", type=Path, default=Path.cwd(),
                   help="Repository root for git-diff scoping (default: cwd).")
    return p


def _select_wanted(args, res_defs: dict[str, dict]) -> list[str] | None:
    """Return the resource list selected by CLI args, or None to signal a clean
    no-op (e.g. nothing to diff against)."""
    if args.all:
        return sorted(res_defs)
    if args.diff_against:
        wanted = resources_changed_between(
            args.diff_against, args.res_def.resolve(), args.repo_root.resolve()
        )
        if not wanted:
            print(f"No resourceDefinition.json changes vs {args.diff_against}; nothing to test.")
            return None
        return wanted
    return [r.strip() for r in args.resources.split(",") if r.strip()]


def _classify_resource(rt: str, schema: dict, res_defs: dict[str, dict]) -> str | None:
    """Return a skip reason for ``rt``, or None if it should be generated."""
    if not rt.startswith("azurerm_"):
        return "not-an-azurerm-resource"
    if rt not in schema:
        return "missing-from-azurerm-provider"
    if rt not in res_defs:
        return "missing-from-resourceDefinition.json"
    return None


def _emit_workspace(out_dir: Path, rt: str, schema_entry: dict, res_def: dict, plugin_dir: str) -> None:
    d = out_dir / rt
    if d.exists():
        shutil.rmtree(d)
    d.mkdir(parents=True)
    name_attr = find_name_attr(schema_entry) or "name"
    (d / "main.tf").write_text(make_main_tf(rt, schema_entry))
    (d / "terraform.rc").write_text(make_terraform_rc(plugin_dir))
    (d / "tests").mkdir()
    (d / "tests" / "validate.tftest.hcl").write_text(make_test_hcl(rt, res_def, name_attr))


def main() -> int:
    args = _build_arg_parser().parse_args()

    res_defs = load_resource_definitions(args.res_def)
    schema = load_schema(args.schema_file)

    wanted = _select_wanted(args, res_defs)
    if wanted is None:
        return 0
    if not wanted:
        print("No resources selected.", file=sys.stderr)
        return 2

    args.out_dir.mkdir(parents=True, exist_ok=True)

    generated: list[str] = []
    skipped: list[tuple[str, str]] = []
    for rt in wanted:
        reason = _classify_resource(rt, schema, res_defs)
        if reason is not None:
            skipped.append((rt, reason))
            continue
        _emit_workspace(args.out_dir, rt, schema[rt], res_defs[rt], args.plugin_dir)
        generated.append(rt)

    print(f"generated: {len(generated)} workspaces under {args.out_dir}")
    if skipped:
        print(f"skipped:   {len(skipped)}")
        for rt, reason in skipped:
            print(f"  - {rt}: {reason}")

    if not generated and skipped:
        # Requested resources were all unusable; surface a non-zero exit so CI
        # doesn't silently succeed when nothing was actually validated.
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
