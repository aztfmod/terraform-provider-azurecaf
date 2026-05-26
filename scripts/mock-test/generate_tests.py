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
    "lock_level":                      '"CanNotDelete"',
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
    # Enum defaults for common attributes
    "sku":                             '"Standard"',
    "action":                          '"Allow"',
    "frequency":                       '"Daily"',
    "time":                            '"23:00"',
    "allocation_method":               '"Static"',
    "application_type":                '"web"',
    "authentication_type":             '"Basic"',
    "data_residency_location":         '"Europe"',
    "data_location":                   '"Europe"',
    "datastore_type":                  '"VaultStore"',
    "domain_join_type":                '"AzureADJoin"',
    "filter_type":                     '"SqlFilter"',
    "format":                          '"Cer"',
    "load_balancer_type":              '"BreadthFirst"',
    "namespace_type":                  '"Messaging"',
    "next_hop_type":                   '"VnetLocal"',
    "offer_type":                      '"Standard"',
    "protocol":                        '"http"',
    "reliability_level":               '"Silver"',
    "service_level":                   '"Standard"',
    "source_platform":                 '"SQL"',
    "storage_account_type":            '"Standard_LRS"',
    "storage_type":                    '"Standard"',
    "traffic_routing_method":          '"Performance"',
    "upgrade_policy_mode":             '"Manual"',
    "collation":                       '"SQL_Latin1_General_CP1_CI_AS"',
    "license_type":                    '"PAYG"',
    "family":                          '"C"',
    "capacity":                        '1',
    "private_ip_address_allocation":   '"Dynamic"',
    "operator":                        '"GreaterThan"',
    "category":                        '"Administrative"',
    # Format-specific defaults
    "end_ip_address":                  '"10.0.0.255"',
    "start_ip_address":                '"10.0.0.1"',
    "gateway_address":                 '"10.0.0.1"',
    "network_subnet_cidr":             '"10.0.0.0/24"',
    "address_prefix":                  '"10.0.0.0/24"',
    "address_space":                   '["10.0.0.0/16"]',
    "address_prefixes":                '["10.0.0.0/24"]',
    "dns_servers":                     '[]',
    "email":                           '"admin@test.local"',
    "email_address":                   '"admin@test.local"',
    "record":                          '"test.example.com"',
    "value":                           '"test-value"',
    "body":                            '"{\\"test\\":\\"value\\"}"',
    "node_count":                      '3',
    "cluster_version":                 '"5.1"',
    "vpn_authentication_types":        '["Certificate"]',
    "key_opts":                        '["encrypt", "decrypt"]',
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
        "iothub_id": f'"{RG_SCOPE}/providers/Microsoft.Devices/iotHubs/iot-test"',
    },
    "azurerm_iothub_endpoint_servicebus_queue": {
        "connection_string": '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k=;EntityPath=q"',
        "iothub_id": f'"{RG_SCOPE}/providers/Microsoft.Devices/iotHubs/iot-test"',
    },
    "azurerm_iothub_endpoint_servicebus_topic": {
        "connection_string": '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k=;EntityPath=t"',
        "iothub_id": f'"{RG_SCOPE}/providers/Microsoft.Devices/iotHubs/iot-test"',
    },
    "azurerm_iothub_endpoint_storage_container": {
        "connection_string": '"DefaultEndpointsProtocol=https;AccountName=sttest;AccountKey=dGVzdA==;EndpointSuffix=core.windows.net"',
        "container_name": '"ctr-test"',
        "iothub_id": f'"{RG_SCOPE}/providers/Microsoft.Devices/iotHubs/iot-test"',
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
    # --- At-least-one-of constraint fixes ---
    "azurerm_app_service_certificate": {
        "key_vault_secret_id": '"https://kv-test.vault.azure.net/secrets/cert/00000000000000000000000000000000"',
    },
    "azurerm_data_factory_dataset_delimited_text": {
        "linked_service_name": '"ls-test"',
    },
    "azurerm_data_factory_linked_service_azure_blob_storage": {
        "connection_string": '"DefaultEndpointsProtocol=https;AccountName=sttest;AccountKey=dGVzdA==;EndpointSuffix=core.windows.net"',
    },
    "azurerm_data_factory_linked_service_azure_databricks": {
        "access_token": '"dapifaketoken123"',
        "existing_cluster_id": '"0101-010101-test1"',
        "adb_domain": '"https://adb-0000000000000000.0.azuredatabricks.net"',
    },
    "azurerm_data_factory_linked_service_azure_function": {
        "key": '"test-function-key"',
        "url": '"https://func-test.azurewebsites.net"',
    },
    "azurerm_data_factory_linked_service_azure_sql_database": {
        "connection_string": '"Server=tcp:sql-test.database.windows.net,1433;Database=db-test;User ID=admin;Password=p;Encrypt=true;"',
    },
    "azurerm_data_factory_linked_service_sql_server": {
        "connection_string": '"Server=tcp:sql-test.database.windows.net,1433;Database=db-test;User ID=admin;Password=p;Encrypt=true;"',
    },
    "azurerm_data_factory_trigger_schedule": {
        "pipeline_name": '"pl-test"',
    },
    "azurerm_disk_encryption_set": {
        "key_vault_key_id": '"https://kv-test.vault.azure.net/keys/test/00000000000000000000000000000000"',
    },
    "azurerm_dns_cname_record": {
        "record": '"cname.example.com"',
    },
    "azurerm_federated_identity_credential": {
        "parent_id": f'"{RG_SCOPE}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/uai-test"',
        "audience": '["api://AzureADTokenExchange"]',
        "issuer": '"https://token.actions.githubusercontent.com"',
        "subject": '"repo:test/test:ref:refs/heads/main"',
    },
    "azurerm_key_vault_secret": {
        "value": '"secret-value"',
    },
    "azurerm_kusto_attached_database_configuration": {
        "cluster_resource_id": f'"{RG_SCOPE}/providers/Microsoft.Kusto/clusters/kc-test"',
    },
    "azurerm_linux_function_app_slot": {
        "storage_account_name": '"sttest"',
        "storage_account_access_key": '"dGVzdGtleQ=="',
    },
    "azurerm_windows_function_app_slot": {
        "storage_account_name": '"sttest"',
        "storage_account_access_key": '"dGVzdGtleQ=="',
    },
    "azurerm_local_network_gateway": {
        "gateway_address": '"10.0.0.1"',
    },
    "azurerm_monitor_metric_alert": {
        "scopes": f'["{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest"]',
    },
    "azurerm_network_security_rule": {
        "destination_address_prefix": '"*"',
        "source_address_prefix": '"*"',
        "source_port_range": '"*"',
        "destination_port_range": '"443"',
        "access": '"Allow"',
        "direction": '"Inbound"',
        "protocol": '"Tcp"',
    },
    "azurerm_private_endpoint": {
        "private_connection_resource_id": f'"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest"',
    },
    "azurerm_private_link_service": {
        "load_balancer_frontend_ip_configuration_ids": f'["{RG_SCOPE}/providers/Microsoft.Network/loadBalancers/lb-test/frontendIPConfigurations/fe-test"]',
    },
    "azurerm_role_assignment": {
        "role_definition_name": '"Reader"',
        "principal_id": f'"{FAKE_OBJ_ID}"',
    },
    "azurerm_storage_container": {
        "storage_account_name": '"sttest"',
    },
    "azurerm_storage_queue": {
        "storage_account_name": '"sttest"',
    },
    "azurerm_storage_share": {
        "storage_account_name": '"sttest"',
        "quota": '50',
    },
    "azurerm_storage_share_directory": {
        "storage_share_id": f'"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest/fileServices/default/fileshares/fs-test"',
    },
    "azurerm_subnet": {
        "virtual_network_name": '"vnet-test"',
        "address_prefixes": '["10.0.1.0/24"]',
    },
    "azurerm_synapse_spark_pool": {
        "node_count": '3',
        "node_size": '"Small"',
        "node_size_family": '"MemoryOptimized"',
    },
    "azurerm_virtual_network": {
        "address_space": '["10.0.0.0/16"]',
    },
    # --- ID parsing fixes (resources needing specific parent IDs) ---
    "azurerm_cdn_frontdoor_custom_domain": {
        "cdn_frontdoor_profile_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test"',
        "host_name": '"custom.example.com"',
    },
    "azurerm_cdn_frontdoor_endpoint": {
        "cdn_frontdoor_profile_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test"',
    },
    "azurerm_cdn_frontdoor_origin_group": {
        "cdn_frontdoor_profile_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test"',
    },
    "azurerm_cdn_frontdoor_origin": {
        "cdn_frontdoor_origin_group_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test/originGroups/og-test"',
        "host_name": '"origin.example.com"',
        "certificate_name_check_enabled": 'false',
    },
    "azurerm_cdn_frontdoor_route": {
        "cdn_frontdoor_endpoint_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test/afdEndpoints/ep-test"',
        "cdn_frontdoor_origin_group_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test/originGroups/og-test"',
        "cdn_frontdoor_origin_ids": f'["{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test/originGroups/og-test/origins/o-test"]',
        "patterns_to_match": '["/*"]',
        "supported_protocols": '["Http", "Https"]',
    },
    "azurerm_cdn_frontdoor_rule_set": {
        "cdn_frontdoor_profile_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test"',
    },
    "azurerm_cdn_frontdoor_rule": {
        "cdn_frontdoor_rule_set_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test/ruleSets/rs-test"',
        "order": '1',
    },
    "azurerm_cdn_frontdoor_secret": {
        "cdn_frontdoor_profile_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test"',
    },
    "azurerm_container_app": {
        "container_app_environment_id": f'"{RG_SCOPE}/providers/Microsoft.App/managedEnvironments/cae-test"',
        "revision_mode": '"Single"',
    },
    "azurerm_cognitive_deployment": {
        "cognitive_account_id": f'"{RG_SCOPE}/providers/Microsoft.CognitiveServices/accounts/cog-test"',
    },
    "azurerm_digital_twins_endpoint_eventgrid": {
        "digital_twins_id": f'"{RG_SCOPE}/providers/Microsoft.DigitalTwins/digitalTwinsInstances/dt-test"',
        "eventgrid_topic_endpoint": '"https://eg-test.westeurope-1.eventgrid.azure.net/api/events"',
        "eventgrid_topic_primary_access_key": '"dGVzdGtleQ=="',
    },
    "azurerm_digital_twins_endpoint_eventhub": {
        "digital_twins_id": f'"{RG_SCOPE}/providers/Microsoft.DigitalTwins/digitalTwinsInstances/dt-test"',
        "eventhub_primary_connection_string": '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k=;EntityPath=eh"',
        "eventhub_secondary_connection_string": '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k=;EntityPath=eh"',
    },
    "azurerm_digital_twins_endpoint_servicebus": {
        "digital_twins_id": f'"{RG_SCOPE}/providers/Microsoft.DigitalTwins/digitalTwinsInstances/dt-test"',
        "servicebus_primary_connection_string": '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k="',
        "servicebus_secondary_connection_string": '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=k;SharedAccessKey=k="',
    },
    "azurerm_dedicated_host": {
        "dedicated_host_group_id": f'"{RG_SCOPE}/providers/Microsoft.Compute/hostGroups/hg-test"',
        "sku_name": '"DSv3-Type1"',
        "platform_fault_domain": '0',
    },
    "azurerm_dev_center_catalog": {
        "dev_center_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test"',
    },
    "azurerm_dev_center_dev_box_definition": {
        "dev_center_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test"',
        "image_reference_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test/galleries/Default/images/img-test"',
        "sku_name": '"general_i_8c32gb256ssd_v2"',
    },
    "azurerm_dev_center_environment_type": {
        "dev_center_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test"',
    },
    "azurerm_dev_center_project_environment_type": {
        "dev_center_project_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/projects/proj-test"',
        "deployment_target_id": f'"/subscriptions/{FAKE_SUB_ID}"',
    },
    "azurerm_healthcare_dicom_service": {
        "workspace_id": f'"{RG_SCOPE}/providers/Microsoft.HealthcareApis/workspaces/hcw-test"',
    },
    "azurerm_healthcare_fhir_service": {
        "workspace_id": f'"{RG_SCOPE}/providers/Microsoft.HealthcareApis/workspaces/hcw-test"',
        "kind": '"fhir-R4"',
        "authentication_authority": f'"https://login.microsoftonline.com/{FAKE_TENANT_ID}"',
        "authentication_audience": '"https://hcw-test.fhir.azurehealthcareapis.com"',
    },
    "azurerm_lb_backend_address_pool": {
        "loadbalancer_id": f'"{RG_SCOPE}/providers/Microsoft.Network/loadBalancers/lb-test"',
    },
    "azurerm_lb_nat_pool": {
        "loadbalancer_id": f'"{RG_SCOPE}/providers/Microsoft.Network/loadBalancers/lb-test"',
        "frontend_ip_configuration_name": '"fe-test"',
        "protocol": '"Tcp"',
        "frontend_port_start": '80',
        "frontend_port_end": '81',
        "backend_port": '8080',
    },
    "azurerm_lb_nat_rule": {
        "loadbalancer_id": f'"{RG_SCOPE}/providers/Microsoft.Network/loadBalancers/lb-test"',
        "frontend_ip_configuration_name": '"fe-test"',
        "protocol": '"Tcp"',
        "frontend_port": '443',
        "backend_port": '8443',
    },
    "azurerm_lb_outbound_rule": {
        "loadbalancer_id": f'"{RG_SCOPE}/providers/Microsoft.Network/loadBalancers/lb-test"',
        "protocol": '"All"',
        "backend_address_pool_id": f'"{RG_SCOPE}/providers/Microsoft.Network/loadBalancers/lb-test/backendAddressPools/bap-test"',
    },
    "azurerm_lb_probe": {
        "loadbalancer_id": f'"{RG_SCOPE}/providers/Microsoft.Network/loadBalancers/lb-test"',
        "port": '443',
    },
    "azurerm_lb_rule": {
        "loadbalancer_id": f'"{RG_SCOPE}/providers/Microsoft.Network/loadBalancers/lb-test"',
        "frontend_ip_configuration_name": '"fe-test"',
        "protocol": '"Tcp"',
        "frontend_port": '80',
        "backend_port": '80',
    },
    "azurerm_logic_app_action_http": {
        "logic_app_id": f'"{RG_SCOPE}/providers/Microsoft.Logic/workflows/la-test"',
        "method": '"GET"',
        "uri": '"https://example.com"',
    },
    "azurerm_logic_app_trigger_http_request": {
        "logic_app_id": f'"{RG_SCOPE}/providers/Microsoft.Logic/workflows/la-test"',
        "schema": '"{{}}',
    },
    "azurerm_machine_learning_compute_instance": {
        "machine_learning_workspace_id": f'"{RG_SCOPE}/providers/Microsoft.MachineLearningServices/workspaces/mlw-test"',
        "virtual_machine_size": '"Standard_DS2_v2"',
    },
    "azurerm_machine_learning_workspace": {
        "application_insights_id": f'"{RG_SCOPE}/providers/Microsoft.Insights/components/appi-test"',
    },
    "azurerm_hpc_cache_blob_target": {
        "storage_container_id": f'"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest/blobServices/default/containers/ct-test"',
    },
    "azurerm_api_management_api_operation_tag": {
        "api_operation_id": f'"{RG_SCOPE}/providers/Microsoft.ApiManagement/service/apim-test/apis/api-test/operations/op-test"',
    },
    "azurerm_api_management_gateway": {
        "api_management_id": f'"{RG_SCOPE}/providers/Microsoft.ApiManagement/service/apim-test"',
    },
    "azurerm_application_insights_web_test": {
        "application_insights_id": f'"{RG_SCOPE}/providers/Microsoft.Insights/components/appi-test"',
        "geo_locations": '["us-tx-sn1-azr"]',
        "configuration": '"<WebTest Name=\\"test\\" Url=\\"https://example.com\\" />"',
    },
    "azurerm_blueprint_assignment": {
        "target_subscription_id": f'"/subscriptions/{FAKE_SUB_ID}"',
        "version_id": f'"/providers/Microsoft.Management/managementGroups/mg-test/providers/Microsoft.Blueprint/blueprints/bp-test/versions/1.0"',
    },
    "azurerm_data_protection_backup_policy_postgresql_flexible_server": {
        "vault_id": f'"{RG_SCOPE}/providers/Microsoft.DataProtection/backupVaults/bv-test"',
    },
    "azurerm_linux_web_app": {
        "service_plan_id": f'"{RG_SCOPE}/providers/Microsoft.Web/serverFarms/asp-test"',
    },
    "azurerm_network_connection_monitor": {
        "network_watcher_id": f'"{RG_SCOPE}/providers/Microsoft.Network/networkWatchers/nw-test"',
    },
    "azurerm_point_to_site_vpn_gateway": {
        "virtual_hub_id": f'"{RG_SCOPE}/providers/Microsoft.Network/virtualHubs/vh-test"',
        "vpn_server_configuration_id": f'"{RG_SCOPE}/providers/Microsoft.Network/vpnServerConfigurations/vpnsc-test"',
        "scale_unit": '1',
    },
    "azurerm_portal_dashboard": {
        "dashboard_properties": '"{\\"lenses\\": {{}}}"',
    },
    "azurerm_postgresql_flexible_server_database": {
        "server_id": f'"{RG_SCOPE}/providers/Microsoft.DBforPostgreSQL/flexibleServers/psql-test"',
    },
    "azurerm_private_dns_resolver_dns_forwarding_ruleset": {
        "private_dns_resolver_outbound_endpoint_ids": f'["{RG_SCOPE}/providers/Microsoft.Network/dnsResolvers/dnsr-test/outboundEndpoints/oep-test"]',
    },
    "azurerm_private_dns_resolver_forwarding_rule": {
        "dns_forwarding_ruleset_id": f'"{RG_SCOPE}/providers/Microsoft.Network/dnsForwardingRulesets/frs-test"',
        "domain_name": '"example.com."',
    },
    "azurerm_private_dns_resolver_inbound_endpoint": {
        "private_dns_resolver_id": f'"{RG_SCOPE}/providers/Microsoft.Network/dnsResolvers/dnsr-test"',
    },
    "azurerm_private_dns_resolver_outbound_endpoint": {
        "private_dns_resolver_id": f'"{RG_SCOPE}/providers/Microsoft.Network/dnsResolvers/dnsr-test"',
    },
    "azurerm_kusto_eventhub_data_connection": {
        "cluster_name": '"kc-test"',
        "database_name": '"db-test"',
        "eventhub_id": f'"{RG_SCOPE}/providers/Microsoft.EventHub/namespaces/ehn-test/eventhubs/eh-test"',
    },
    "azurerm_monitor_smart_detector_alert_rule": {
        "scope_resource_ids": f'["{RG_SCOPE}/providers/Microsoft.Insights/components/appi-test"]',
        "detector_type": '"FailureAnomaliesDetector"',
        "action_group_resource_ids": f'["{RG_SCOPE}/providers/Microsoft.Insights/actionGroups/ag-test"]',
    },
    "azurerm_data_share": {
        "account_id": f'"{RG_SCOPE}/providers/Microsoft.DataShare/accounts/dsh-test"',
        "kind": '"CopyBased"',
    },
    # --- SKU and enum fixes ---
    "azurerm_api_management": {
        "sku_name": '"Developer_1"',
        "publisher_name": '"Test Publisher"',
        "publisher_email": '"admin@test.local"',
    },
    "azurerm_automation_account": {"sku_name": '"Basic"'},
    "azurerm_cdn_frontdoor_profile": {"sku_name": '"Standard_AzureFrontDoor"'},
    "azurerm_database_migration_service": {
        "sku_name": '"Standard_1vCores"',
        "subnet_id": f'"{RG_SCOPE}/providers/Microsoft.Network/virtualNetworks/vnet-test/subnets/snet-test"',
    },
    "azurerm_firewall": {
        "sku_name": '"AZFW_VNet"',
        "sku_tier": '"Standard"',
    },
    "azurerm_logic_app_integration_account": {"sku_name": '"Standard"'},
    "azurerm_managed_redis": {"sku_name": '"Balanced_B1"'},
    "azurerm_maps_account": {"sku_name": '"S0"'},
    "azurerm_postgresql_server": {
        "sku_name": '"GP_Gen5_2"',
        "version": '"11"',
        "ssl_enforcement_enabled": 'true',
        "administrator_login": '"adminuser"',
        "administrator_login_password": '"P@ssw0rd1234!"',
    },
    "azurerm_relay_namespace": {"sku_name": '"Standard"'},
    "azurerm_synapse_sql_pool": {"sku_name": '"DW100c"'},
    "azurerm_analysis_services_server": {"sku": '"S0"'},
    "azurerm_cdn_profile": {"sku": '"Standard_Microsoft"'},
    "azurerm_container_registry": {"sku": '"Standard"'},
    "azurerm_databricks_workspace": {"sku": '"standard"'},
    "azurerm_eventhub_namespace": {"sku": '"Standard"'},
    "azurerm_recovery_services_vault": {"sku": '"Standard"'},
    "azurerm_search_service": {"sku": '"standard"'},
    "azurerm_servicebus_namespace": {"sku": '"Standard"'},
    "azurerm_virtual_network_gateway": {
        "sku": '"Basic"',
        "type": '"Vpn"',
        "vpn_type": '"RouteBased"',
    },
    "azurerm_web_pubsub": {"sku": '"Free_F1"', "capacity": '1'},
    # --- Additional enum fixes ---
    "azurerm_aadb2c_directory": {"data_residency_location": '"Europe"'},
    "azurerm_application_insights": {"application_type": '"web"'},
    "azurerm_bot_channel_registration": {"sku": '"F0"', "microsoft_app_id": f'"{FAKE_OBJ_ID}"'},
    "azurerm_bot_channels_registration": {"sku": '"F0"', "microsoft_app_id": f'"{FAKE_OBJ_ID}"'},
    "azurerm_bot_web_app": {"sku": '"F0"', "microsoft_app_id": f'"{FAKE_OBJ_ID}"'},
    "azurerm_cognitive_account": {"sku_name": '"S0"', "kind": '"Face"'},
    "azurerm_communication_service": {"data_location": '"Europe"'},
    "azurerm_consumption_budget_resource_group": {"amount": '100'},
    "azurerm_cosmosdb_account": {"offer_type": '"Standard"', "consistency_level": '"Session"'},
    "azurerm_data_factory_dataset_delimited_text": {
        "linked_service_name": '"ls-test"',
    },
    "azurerm_dns_caa_record": {"tag": '"issue"'},
    "azurerm_express_route_circuit": {
        "bandwidth_in_mbps": '50',
        "peering_location": '"London"',
        "service_provider_name": '"Equinix"',
    },
    "azurerm_frontdoor": {
        "backend_pool_health_probe_name": '"hp-test"',
        "backend_pool_load_balancing_name": '"lb-test"',
        "priority": '1',
    },
    "azurerm_hdinsight_hadoop_cluster": {
        "cluster_version": '"5.1"',
        "tier": '"Standard"',
    },
    "azurerm_hdinsight_hbase_cluster": {
        "cluster_version": '"5.1"',
        "tier": '"Standard"',
    },
    "azurerm_hdinsight_interactive_query_cluster": {
        "cluster_version": '"5.1"',
        "tier": '"Standard"',
    },
    "azurerm_hdinsight_kafka_cluster": {
        "cluster_version": '"5.1"',
        "tier": '"Standard"',
    },
    "azurerm_hdinsight_spark_cluster": {
        "cluster_version": '"5.1"',
        "tier": '"Standard"',
    },
    "azurerm_key_vault": {
        "sku_name": '"standard"',
        "tenant_id": f'"{FAKE_TENANT_ID}"',
    },
    "azurerm_key_vault_key": {
        "key_type": '"RSA"',
        "key_size": '2048',
        "key_opts": '["encrypt", "decrypt"]',
    },
    "azurerm_kusto_cluster": {
        "sku": '{{ name = "Dev(No SLA)_Standard_D11_v2" capacity = 1 }}',
    },
    "azurerm_monitor_action_group": {
        "short_name": '"test"',
    },
    "azurerm_mssql_database": {
        "collation": '"SQL_Latin1_General_CP1_CI_AS"',
        "server_id": f'"{RG_SCOPE}/providers/Microsoft.Sql/servers/sql-test"',
    },
    "azurerm_mssql_elasticpool": {
        "sku": '{{ name = "BasicPool" tier = "Basic" capacity = 50 }}',
    },
    "azurerm_netapp_pool": {
        "service_level": '"Standard"',
        "size_in_tb": '4',
    },
    "azurerm_netapp_volume": {
        "service_level": '"Standard"',
        "storage_quota_in_gb": '100',
        "volume_path": '"vol-test"',
        "protocols": '["NFSv3"]',
    },
    "azurerm_notification_hub_namespace": {
        "namespace_type": '"Messaging"',
        "sku_name": '"Free"',
    },
    "azurerm_powerbi_embedded": {
        "sku_name": '"A1"',
        "administrators": '["admin@test.local"]',
    },
    "azurerm_route": {
        "next_hop_type": '"VnetLocal"',
        "address_prefix": '"10.0.0.0/24"',
    },
    "azurerm_service_fabric_cluster": {
        "reliability_level": '"Silver"',
        "upgrade_mode": '"Automatic"',
        "vm_image": '"Windows"',
        "management_endpoint": '"https://sf-test.westeurope.cloudapp.azure.com:19080"',
    },
    "azurerm_signalr_service": {
        "sku": '{{ name = "Free_F1" capacity = 1 }}',
    },
    "azurerm_stream_analytics_output_eventhub": {
        "eventhub_name": '"eh-test"',
        "servicebus_namespace": '"ehn-test"',
        "shared_access_policy_key": '"dGVzdGtleQ=="',
        "shared_access_policy_name": '"RootManageSharedAccessKey"',
        "type": '"Json"',
    },
    "azurerm_stream_analytics_output_blob": {
        "path_pattern": '"{yyyy}/{MM}/{dd}"',
        "date_format": '"yyyy/MM/dd"',
        "time_format": '"HH"',
        "type": '"Json"',
    },
    "azurerm_stream_analytics_output_servicebus_queue": {
        "type": '"Json"',
    },
    "azurerm_stream_analytics_output_servicebus_topic": {
        "type": '"Json"',
    },
    "azurerm_stream_analytics_reference_input_blob": {
        "path_pattern": '"{yyyy}/{MM}/{dd}"',
        "date_format": '"yyyy/MM/dd"',
        "time_format": '"HH"',
        "type": '"Json"',
    },
    "azurerm_stream_analytics_stream_input_blob": {
        "path_pattern": '"{yyyy}/{MM}/{dd}"',
        "date_format": '"yyyy/MM/dd"',
        "time_format": '"HH"',
        "type": '"Json"',
    },
    "azurerm_stream_analytics_stream_input_eventhub": {
        "eventhub_name": '"eh-test"',
        "servicebus_namespace": '"ehn-test"',
        "shared_access_policy_key": '"dGVzdGtleQ=="',
        "shared_access_policy_name": '"RootManageSharedAccessKey"',
        "type": '"Json"',
    },
    "azurerm_stream_analytics_stream_input_iothub": {
        "type": '"Json"',
        "eventhub_consumer_group_name": '"$Default"',
        "shared_access_policy_key": '"dGVzdGtleQ=="',
        "shared_access_policy_name": '"iothubowner"',
        "endpoint": '"messages/events"',
    },
    "azurerm_traffic_manager_profile": {
        "traffic_routing_method": '"Performance"',
    },
    "azurerm_virtual_machine_scale_set": {
        "upgrade_policy_mode": '"Manual"',
    },
    "azurerm_virtual_desktop_host_pool": {
        "load_balancer_type": '"BreadthFirst"',
        "type": '"Pooled"',
    },
    "azurerm_vpn_gateway_connection": {
        "vpn_gateway_id": f'"{RG_SCOPE}/providers/Microsoft.Network/vpnGateways/vpng-test"',
    },
    # --- IoT Security ---
    "azurerm_iot_security_device_group": {
        "iothub_id": f'"{RG_SCOPE}/providers/Microsoft.Devices/iotHubs/iot-test"',
    },
    "azurerm_iot_security_solution": {
        "iothub_ids": f'["{RG_SCOPE}/providers/Microsoft.Devices/iotHubs/iot-test"]',
        "display_name": '"IoT Security"',
    },
    # --- Synapse resources ---
    "azurerm_synapse_integration_runtime_azure": {
        "synapse_workspace_id": f'"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test"',
    },
    "azurerm_synapse_integration_runtime_self_hosted": {
        "synapse_workspace_id": f'"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test"',
    },
    "azurerm_synapse_linked_service": {
        "synapse_workspace_id": f'"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test"',
        "type": '"AzureBlobStorage"',
        "type_properties_json": '"{\\"connectionString\\":\\"DefaultEndpointsProtocol=https;AccountName=sttest;\\"}"',
    },
    "azurerm_synapse_managed_private_endpoint": {
        "synapse_workspace_id": f'"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test"',
        "target_resource_id": f'"{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest"',
        "subresource_name": '"blob"',
    },
    "azurerm_synapse_sql_pool_vulnerability_assessment_baseline": {
        "sql_pool_vulnerability_assessment_id": f'"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test/sqlPools/sp-test/vulnerabilityAssessments/Default"',
        "rule_id": '"VA1234"',
    },
    "azurerm_synapse_sql_pool_workload_classifier": {
        "workload_group_id": f'"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test/sqlPools/sp-test/workloadGroups/wg-test"',
        "member_name": '"dbo"',
    },
    "azurerm_synapse_sql_pool_workload_group": {
        "sql_pool_id": f'"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test/sqlPools/sp-test"',
        "min_resource_percent": '0',
        "max_resource_percent": '100',
        "min_resource_percent_per_request": '3',
    },
    # --- Range/format fixes ---
    "azurerm_vmware_cluster": {"cluster_node_count": '3', "sku_name": '"av36"'},
    "azurerm_monitor_scheduled_query_rules_alert": {
        "frequency": '5',
        "time_window": '5',
        "severity": '3',
    },
    "azurerm_redis_firewall_rule": {
        "start_ip": '"10.0.0.1"',
        "end_ip": '"10.0.0.255"',
    },
    "azurerm_postgresql_database": {
        "collation": '"en_US.utf8"',
        "charset": '"UTF8"',
    },
    "azurerm_mssql_server": {
        "administrator_login": '"adminuser"',
        "administrator_login_password": '"P@ssw0rd1234!"',
        "version": '"12.0"',
    },
    "azurerm_resource_group_policy_assignment": {
        "policy_definition_id": '"/providers/Microsoft.Authorization/policyDefinitions/00000000-0000-0000-0000-000000000000"',
    },
    # HDInsight clusters — vm_size enum
    "azurerm_hdinsight_hadoop_cluster": {
        "cluster_version": '"5.1"',
        "tier": '"Standard"',
        "vm_size": '"Standard_D3_V2"',
    },
    "azurerm_hdinsight_hbase_cluster": {
        "cluster_version": '"5.1"',
        "tier": '"Standard"',
        "vm_size": '"Standard_D3_V2"',
    },
    "azurerm_hdinsight_interactive_query_cluster": {
        "cluster_version": '"5.1"',
        "tier": '"Standard"',
        "vm_size": '"Standard_D3_V2"',
    },
    "azurerm_hdinsight_kafka_cluster": {
        "cluster_version": '"5.1"',
        "tier": '"Standard"',
        "vm_size": '"Standard_D3_V2"',
        "number_of_disks_per_node": "3",
    },
    "azurerm_hdinsight_spark_cluster": {
        "cluster_version": '"5.1"',
        "tier": '"Standard"',
        "vm_size": '"Standard_D3_V2"',
    },
    # Stream analytics — serialization.format
    "azurerm_stream_analytics_output_blob": {
        "type": '"Json"',
        "format": '"LineSeparated"',
    },
    "azurerm_stream_analytics_output_eventhub": {
        "type": '"Json"',
        "format": '"Array"',
    },
    "azurerm_stream_analytics_output_servicebus_queue": {
        "type": '"Json"',
        "format": '"Array"',
    },
    "azurerm_stream_analytics_output_servicebus_topic": {
        "type": '"Json"',
        "format": '"Array"',
    },
    # EventHub — needs namespace_id instead of resource_group_name
    "azurerm_eventhub": {
        "namespace_id": f'"{RG_SCOPE}/providers/Microsoft.EventHub/namespaces/evhns-test"',
        "partition_count": "2",
        "message_retention": "1",
    },
    "azurerm_eventhub_authorization_rule": {
        "namespace_name": '"evhnstest"',
        "eventhub_name": '"evh-test"',
    },
    "azurerm_eventhub_consumer_group": {
        "namespace_name": '"evhnstest"',
        "eventhub_name": '"evh-test"',
    },
    "azurerm_eventhub_namespace_authorization_rule": {
        "namespace_name": '"evhnstest"',
    },
    "azurerm_eventhub_namespace_disaster_recovery_config": {
        "namespace_name": '"evhnstest"',
    },
    # Linux/Windows VM — needs source_image_reference block (at-least-one-of)
    "azurerm_linux_virtual_machine": {
        "source_image_id": f'"{RG_SCOPE}/providers/Microsoft.Compute/images/img-test"',
        "admin_username": '"adminuser"',
        "admin_password": '"P@ssword1234!"',
        "disable_password_authentication": "false",
        "size": '"Standard_B1s"',
        "network_interface_ids": f'["{RG_SCOPE}/providers/Microsoft.Network/networkInterfaces/nic-test"]',
        "caching": '"ReadWrite"',
        "storage_account_type": '"Standard_LRS"',
    },
    "azurerm_windows_virtual_machine": {
        "source_image_id": f'"{RG_SCOPE}/providers/Microsoft.Compute/images/img-test"',
        "admin_username": '"adminuser"',
        "admin_password": '"P@ssword1234!"',
        "size": '"Standard_B1s"',
        "network_interface_ids": f'["{RG_SCOPE}/providers/Microsoft.Network/networkInterfaces/nic-test"]',
        "caching": '"ReadWrite"',
        "storage_account_type": '"Standard_LRS"',
    },
    "azurerm_linux_virtual_machine_scale_set": {
        "source_image_id": f'"{RG_SCOPE}/providers/Microsoft.Compute/images/img-test"',
        "admin_username": '"adminuser"',
        "admin_password": '"P@ssword1234!"',
        "disable_password_authentication": "false",
        "sku": '"Standard_B1s"',
        "instances": "1",
        "caching": '"ReadWrite"',
        "storage_account_type": '"Standard_LRS"',
    },
    "azurerm_windows_virtual_machine_scale_set": {
        "source_image_id": f'"{RG_SCOPE}/providers/Microsoft.Compute/images/img-test"',
        "admin_username": '"adminuser"',
        "admin_password": '"P@ssword1234!"',
        "sku": '"Standard_B1s"',
        "instances": "1",
        "caching": '"ReadWrite"',
        "storage_account_type": '"Standard_LRS"',
    },
    # Kubernetes cluster — needs identity or service_principal
    "azurerm_kubernetes_cluster": {
        "dns_prefix": '"aks-test"',
    },
    # API Management certificate — needs data or key_vault_secret_id
    "azurerm_api_management_certificate": {
        "data": '"dGVzdA=="',
    },
    # Automation runbook — needs content or publish_content_link
    "azurerm_automation_runbook": {
        "content": '"Write-Output Hello"',
    },
    # Application gateway — needs gateway_ip_configuration
    "azurerm_application_gateway": {
        "sku_name": '"Standard_v2"',
        "sku_tier": '"Standard_v2"',
    },
    # Custom provider — needs action or resource_type
    "azurerm_custom_provider": {
        "action": "[{}]",
    },
    # Data factory linked service — needs url
    "azurerm_data_factory_linked_service_data_lake_storage_gen2": {
        "url": '"https://test.dfs.core.windows.net"',
    },
    # Data factory dataset delimited text — needs one location
    "azurerm_data_factory_dataset_delimited_text": {
        "column_delimiter": '","',
        "row_delimiter": '"\\n"',
    },
    # Data protection backup policy blob — needs retention_duration
    "azurerm_data_protection_backup_policy_blob_storage": {
        "retention_duration": '"P30D"',
    },
    # Data protection backup policy disk
    "azurerm_data_protection_backup_policy_disk": {
        "default_retention_duration": '"P7D"',
        "backup_repeating_time_intervals": '["R/2021-05-19T06:33:16+00:00/PT4H"]',
    },
    # Data protection backup policy postgresql
    "azurerm_data_protection_backup_policy_postgresql": {
        "default_retention_duration": '"P7D"',
        "backup_repeating_time_intervals": '["R/2021-05-19T06:33:16+00:00/PT4H"]',
    },
    # Key vault certificate — needs certificate or certificate_policy
    "azurerm_key_vault_certificate": {
        "key_vault_id": f'"{RG_SCOPE}/providers/Microsoft.KeyVault/vaults/kv-test"',
    },
    # Monitor data collection rule — needs data_flow
    "azurerm_monitor_data_collection_rule": {
        "kind": '"Linux"',
    },
    # Monitor diagnostic setting — needs at least one target
    "azurerm_monitor_diagnostic_setting": {
        "log_analytics_workspace_id": f'"{RG_SCOPE}/providers/Microsoft.OperationalInsights/workspaces/law-test"',
        "target_resource_id": f'"{RG_SCOPE}/providers/Microsoft.KeyVault/vaults/kv-test"',
    },
    # Monitor metric alert — needs criteria
    "azurerm_monitor_metric_alert": {
        "scopes": f'["{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest"]',
    },
    # Monitor smart detector alert rule — needs frequency as ISO 8601 duration
    "azurerm_monitor_smart_detector_alert_rule": {
        "frequency": '"PT1H"',
        "severity": '"Sev0"',
        "detector_type": '"FailureAnomaliesDetector"',
        "scope_resource_ids": f'["{RG_SCOPE}/providers/Microsoft.Insights/components/ai-test"]',
    },
    # IoT Hub endpoint eventhub — needs connection_string or endpoint_uri+entity_path
    "azurerm_iothub_endpoint_eventhub": {
        "connection_string": '"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=test;SharedAccessKey=dGVzdA==;EntityPath=test"',
    },
    # Storage table — name cannot start with number and no hyphens
    "azurerm_storage_table": {
        "storage_account_name": '"sttest"',
    },
    # Consumption budget subscription — needs subscription_id as full resource ID
    "azurerm_consumption_budget_subscription": {
        "subscription_id": '"/subscriptions/00000000-0000-0000-0000-000000000000"',
        "amount": "1000",
        "time_grain": '"Monthly"',
    },
    # ServiceBus subscription — needs topic_id
    "azurerm_servicebus_subscription": {
        "topic_id": f'"{RG_SCOPE}/providers/Microsoft.ServiceBus/namespaces/sbns-test/topics/sbt-test"',
        "max_delivery_count": "10",
    },
    "azurerm_servicebus_subscription_rule": {
        "subscription_id": f'"{RG_SCOPE}/providers/Microsoft.ServiceBus/namespaces/sbns-test/topics/sbt-test/subscriptions/sbsub-test"',
        "filter_type": '"SqlFilter"',
        "sql_filter": '"1=1"',
    },
    "azurerm_servicebus_queue_authorization_rule": {
        "queue_id": f'"{RG_SCOPE}/providers/Microsoft.ServiceBus/namespaces/sbns-test/queues/sbq-test"',
    },
    "azurerm_servicebus_topic_authorization_rule": {
        "topic_id": f'"{RG_SCOPE}/providers/Microsoft.ServiceBus/namespaces/sbns-test/topics/sbt-test"',
    },
    # Dev Center — needs dev_center_id (proper)
    "azurerm_dev_center_catalog": {
        "dev_center_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test"',
    },
    "azurerm_dev_center_dev_box_definition": {
        "dev_center_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test"',
        "image_reference_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test/galleries/Default/images/img-test"',
        "sku_name": '"general_i_8c32gb256ssd_v2"',
    },
    "azurerm_dev_center_environment_type": {
        "dev_center_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test"',
    },
    # Linux/Windows function apps — needs service_plan_id
    "azurerm_linux_function_app": {
        "service_plan_id": f'"{RG_SCOPE}/providers/Microsoft.Web/serverFarms/asp-test"',
        "storage_account_name": '"sttest"',
        "storage_account_access_key": '"dGVzdGtleQ=="',
    },
    "azurerm_linux_web_app": {
        "service_plan_id": f'"{RG_SCOPE}/providers/Microsoft.Web/serverFarms/asp-test"',
    },
    "azurerm_windows_function_app": {
        "service_plan_id": f'"{RG_SCOPE}/providers/Microsoft.Web/serverFarms/asp-test"',
        "storage_account_name": '"sttest"',
        "storage_account_access_key": '"dGVzdGtleQ=="',
    },
    "azurerm_windows_web_app": {
        "service_plan_id": f'"{RG_SCOPE}/providers/Microsoft.Web/serverFarms/asp-test"',
    },
    # Portal dashboard — needs dashboard_properties as valid JSON
    "azurerm_portal_dashboard": {
        "dashboard_properties": '"{\\"lenses\\": {}}"',
    },
    # Network connection monitor — needs network_watcher_id
    "azurerm_network_connection_monitor": {
        "network_watcher_id": f'"{RG_SCOPE}/providers/Microsoft.Network/networkWatchers/nw-test"',
    },
    # PostgreSQL flexible server firewall rule — needs server_id
    "azurerm_postgresql_flexible_server_firewall_rule": {
        "server_id": f'"{RG_SCOPE}/providers/Microsoft.DBforPostgreSQL/flexibleServers/psql-test"',
        "start_ip_address": '"10.0.0.1"',
        "end_ip_address": '"10.0.0.255"',
    },
    # Data share dataset blob storage
    "azurerm_data_share_dataset_blob_storage": {
        "data_share_id": f'"{RG_SCOPE}/providers/Microsoft.DataShare/accounts/dsa-test/shares/ds-test"',
        "container_name": '"testcontainer"',
    },
    # Virtual hub connection
    "azurerm_virtual_hub_connection": {
        "virtual_hub_id": f'"{RG_SCOPE}/providers/Microsoft.Network/virtualHubs/vhub-test"',
        "remote_virtual_network_id": f'"{RG_SCOPE}/providers/Microsoft.Network/virtualNetworks/vnet-test"',
    },
    # Virtual machine extension
    "azurerm_virtual_machine_extension": {
        "virtual_machine_id": f'"{RG_SCOPE}/providers/Microsoft.Compute/virtualMachines/vm-test"',
        "publisher": '"Microsoft.Azure.Extensions"',
        "type_handler_version": '"1.0"',
    },
    # Virtual machine scale set extension
    "azurerm_virtual_machine_scale_set_extension": {
        "virtual_machine_scale_set_id": f'"{RG_SCOPE}/providers/Microsoft.Compute/virtualMachineScaleSets/vmss-test"',
        "publisher": '"Microsoft.Azure.Extensions"',
        "type_handler_version": '"1.0"',
    },
    # Virtual network peering
    "azurerm_virtual_network_peering": {
        "virtual_network_id": f'"{RG_SCOPE}/providers/Microsoft.Network/virtualNetworks/vnet-test"',
        "remote_virtual_network_id": f'"{RG_SCOPE}/providers/Microsoft.Network/virtualNetworks/vnet-remote"',
    },
    # VMware cluster
    "azurerm_vmware_cluster": {
        "vmware_cloud_id": f'"{RG_SCOPE}/providers/Microsoft.AVS/privateClouds/avs-test"',
        "cluster_node_count": "3",
        "sku_name": '"av36"',
    },
    # Storage sync group
    "azurerm_storage_sync_group": {
        "storage_sync_id": f'"{RG_SCOPE}/providers/Microsoft.StorageSync/storageSyncServices/ss-test"',
    },
    # Healthcare medtech service — needs namespace_name format
    "azurerm_healthcare_medtech_service": {
        "eventhub_namespace_name": '"evhnstest"',
        "eventhub_name": '"evh-test"',
        "eventhub_consumer_group_name": '"$Default"',
        "workspace_id": f'"{RG_SCOPE}/providers/Microsoft.HealthcareApis/workspaces/hw-test"',
    },
    # Storage blob — type enum
    "azurerm_storage_blob": {
        "type": '"Block"',
        "storage_account_name": '"sttest"',
        "storage_container_name": '"testcontainer"',
    },
    # Virtual desktop application group
    "azurerm_virtual_desktop_application_group": {
        "type": '"Desktop"',
        "host_pool_id": f'"{RG_SCOPE}/providers/Microsoft.DesktopVirtualization/hostPools/hp-test"',
    },
    # Synapse SQL pool — storage_account_type enum
    "azurerm_synapse_sql_pool": {
        "storage_account_type": '"GRS"',
        "synapse_workspace_id": f'"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test"',
        "sku_name": '"DW100c"',
    },
    # Security center automation
    "azurerm_security_center_automation": {
        "scopes": f'["/subscriptions/00000000-0000-0000-0000-000000000000"]',
    },
    # Monitor metric alert — needs criteria or dynamic_criteria
    "azurerm_monitor_metric_alert": {
        "scopes": f'["{RG_SCOPE}/providers/Microsoft.Storage/storageAccounts/sttest"]',
    },
    # CDN frontdoor firewall policy
    "azurerm_cdn_frontdoor_firewall_policy": {
        "sku_name": '"Standard_AzureFrontDoor"',
        "mode": '"Prevention"',
    },
    # Enum fixes batch
    "azurerm_aadb2c_directory": {
        "sku_name": '"PremiumP1"',
        "data_residency_location": '"United States"',
    },
    "azurerm_application_insights_web_test": {
        "kind": '"ping"',
        "geo_locations": '["us-tx-sn1-azr"]',
        "configuration": '"<WebTest><Items><Request Method=\\\"GET\\\" Url=\\\"https://example.com\\\"/></Items></WebTest>"',
    },
    "azurerm_automation_runbook": {
        "runbook_type": '"PowerShell"',
        "content": '"Write-Output Hello"',
    },
    "azurerm_automation_schedule": {
        "frequency": '"Day"',
    },
    "azurerm_batch_certificate": {
        "thumbprint_algorithm": '"SHA1"',
        "thumbprint": '"312D31A8B5EAD0C15B4C4A369B0EE8E62B28EF25"',
        "certificate": '"dGVzdA=="',
    },
    "azurerm_container_registry_webhook": {
        "actions": '["push"]',
        "service_uri": '"https://example.com/webhook"',
    },
    "azurerm_data_protection_backup_vault": {
        "redundancy": '"LocallyRedundant"',
        "datastore_type": '"VaultStore"',
    },
    "azurerm_database_migration_project": {
        "target_platform": '"SQLDB"',
        "source_platform": '"SQL"',
    },
    "azurerm_express_route_circuit": {
        "family": '"MeteredData"',
    },
    "azurerm_firewall_nat_rule_collection": {
        "action": '"Dnat"',
    },
    "azurerm_firewall_network_rule_collection": {
        "action": '"Allow"',
        "protocols": '["TCP"]',
    },
    "azurerm_frontdoor": {
        "protocol": '"Https"',
    },
    "azurerm_iothub_dps": {
        "sku_name": '"S1"',
    },
    "azurerm_iothub": {
        "name": '"S1"',
        "capacity": "1",
    },
    "azurerm_iothub_dps": {
        "name": '"S1"',
        "capacity": "1",
    },
    "azurerm_kusto_cluster": {
        "name": '"Dev(No SLA)_Standard_D11_v2"',
    },
    "azurerm_signalr_service": {
        "name": '"Standard_S1"',
        "capacity": "1",
    },
    "azurerm_logic_app_trigger_recurrence": {
        "frequency": '"Day"',
    },
    "azurerm_maintenance_configuration": {
        "scope": '"Host"',
    },
    "azurerm_monitor_scheduled_query_rules_log": {
        "operator": '"Include"',
    },
    "azurerm_mssql_elasticpool": {
        "family": '"Gen5"',
    },
    "azurerm_redhat_openshift_cluster": {
        "visibility": '"Public"',
    },
    "azurerm_redis_cache": {
        "sku_name": '"Standard"',
        "family": '"C"',
        "capacity": "1",
    },
    "azurerm_security_center_automation": {
        "scopes": '["/subscriptions/00000000-0000-0000-0000-000000000000"]',
    },
    "azurerm_signalr_service": {
        "sku_name": '"Standard_S1"',
    },
    "azurerm_stream_analytics_function_javascript_udf": {
        "type": '"any"',
    },
    "azurerm_synapse_spark_pool": {
        "spark_version": '"3.4"',
        "node_size": '"Small"',
        "node_size_family": '"MemoryOptimized"',
        "synapse_workspace_id": f'"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test"',
    },
    "azurerm_traffic_manager_profile": {
        "protocol": '"HTTPS"',
        "traffic_routing_method": '"Performance"',
    },
    "azurerm_vmware_private_cloud": {
        "sku_name": '"av36"',
        "management_cluster_size": "3",
        "network_subnet_cidr": '"192.168.48.0/22"',
    },
    "azurerm_web_application_firewall_policy": {
        "type": '"OWASP"',
        "version": '"3.2"',
    },
    # Other specific fixes
    "azurerm_hpc_cache_nfs_target": {
        "nfs_export": '"/export"',
        "target_path": '"/nfs"',
    },
    # Additional fixes - batch 2
    "azurerm_batch_certificate": {
        "account_name": '"batchacct1"',
        "thumbprint_algorithm": '"SHA1"',
        "thumbprint": '"312D31A8B5EAD0C15B4C4A369B0EE8E62B28EF25"',
        "certificate": '"dGVzdA=="',
    },
    "azurerm_bot_service_azure_bot": {
        "microsoft_app_id": '"00000000-0000-0000-0000-000000000000"',
    },
    "azurerm_cdn_frontdoor_secret": {
        "cdn_frontdoor_profile_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnp-test"',
    },
    "azurerm_cdn_frontdoor_security_policy": {
        "cdn_frontdoor_profile_id": f'"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnp-test"',
    },
    "azurerm_consumption_budget_resource_group": {
        "resource_group_id": f'"{RG_SCOPE}"',
        "amount": "1000",
        "time_grain": '"Monthly"',
        "start_date": '"2024-01-01T00:00:00Z"',
    },
    "azurerm_consumption_budget_subscription": {
        "subscription_id": '"/subscriptions/00000000-0000-0000-0000-000000000000"',
        "amount": "1000",
        "time_grain": '"Monthly"',
        "start_date": '"2024-01-01T00:00:00Z"',
    },
    "azurerm_container_registry_webhook": {
        "registry_name": '"testregistry"',
        "actions": '["push"]',
        "service_uri": '"https://example.com/webhook"',
    },
    "azurerm_data_protection_backup_policy_postgresql_flexible_server": {
        "default_retention_duration": '"P7D"',
        "vault_id": f'"{RG_SCOPE}/providers/Microsoft.DataProtection/backupVaults/bv-test"',
    },
    "azurerm_dedicated_hardware_security_module": {
        "network_interface_private_ip_addresses": '["10.0.0.5"]',
    },
    "azurerm_firewall_nat_rule_collection": {
        "action": '"Dnat"',
        "protocols": '["TCP"]',
    },
    "azurerm_frontdoor": {
        "protocol": '"Https"',
        "priority": "1",
    },
    "azurerm_healthcare_medtech_service": {
        "eventhub_namespace_name": '"evhnstest"',
        "eventhub_name": '"evhtest"',
        "eventhub_consumer_group_name": '"consumergroup1"',
        "workspace_id": f'"{RG_SCOPE}/providers/Microsoft.HealthcareApis/workspaces/hw-test"',
    },
    "azurerm_hpc_cache_nfs_target": {
        "nfs_export": '"/export"',
        "target_path": '"nfs"',
    },
    "azurerm_lighthouse_definition": {
        "role_definition_id": '"00000000-0000-0000-0000-000000000000"',
        "principal_id": '"22222222-2222-2222-2222-222222222222"',
    },
    "azurerm_monitor_scheduled_query_rules_alert": {
        "action_group": f'["{RG_SCOPE}/providers/Microsoft.Insights/actionGroups/ag-test"]',
    },
    "azurerm_monitor_smart_detector_alert_rule": {
        "frequency": '"PT1H"',
        "severity": '"Sev0"',
        "detector_type": '"FailureAnomaliesDetector"',
        "scope_resource_ids": f'["{RG_SCOPE}/providers/Microsoft.Insights/components/ai-test"]',
        "action_group_ids": f'["{RG_SCOPE}/providers/Microsoft.Insights/actionGroups/ag-test"]',
    },
    "azurerm_mssql_elasticpool": {
        "name": '"StandardPool"',
        "tier": '"Standard"',
        "family": '"Gen5"',
        "capacity": "50",
    },
    "azurerm_redhat_openshift_cluster": {
        "visibility": '"Public"',
        "version": '"4.13.23"',
        "pull_secret": '""',
    },
    "azurerm_role_assignment": {
        "name": '"00000000-0000-0000-0000-000000000000"',
        "role_definition_id": '"/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000"',
    },
    "azurerm_security_center_automation": {
        "type": '"LogicApp"',
        "resource_id": f'"{RG_SCOPE}/providers/Microsoft.Logic/workflows/la-test"',
        "uri": '"https://example.com"',
        "scopes": '["/subscriptions/00000000-0000-0000-0000-000000000000"]',
    },
    "azurerm_signalr_service": {
        "name": '"Standard_S1"',
        "capacity": "1",
    },
    "azurerm_subscription_policy_assignment": {
        "subscription_id": '"/subscriptions/00000000-0000-0000-0000-000000000000"',
        "policy_definition_id": '"/providers/Microsoft.Authorization/policyDefinitions/00000000-0000-0000-0000-000000000000"',
    },
    "azurerm_synapse_spark_pool": {
        "spark_version": '"3.4"',
        "node_size": '"Small"',
        "node_size_family": '"MemoryOptimized"',
        "node_count": "3",
        "synapse_workspace_id": f'"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test"',
    },
    "azurerm_virtual_network_gateway": {
        "subnet_id": f'"{RG_SCOPE}/providers/Microsoft.Network/virtualNetworks/vnet-test/subnets/GatewaySubnet"',
    },
    "azurerm_vmware_private_cloud": {
        "sku_name": '"av36"',
        "management_cluster_size": "3",
        "size": "3",
        "network_subnet_cidr": '"192.168.48.0/22"',
    },
    "azurerm_vmware_express_route_authorization": {
        "private_cloud_id": f'"{RG_SCOPE}/providers/Microsoft.AVS/privateClouds/avs-test"',
    },
    "azurerm_vpn_gateway_connection": {
        "vpn_gateway_id": f'"{RG_SCOPE}/providers/Microsoft.Network/vpnGateways/vpng-test"',
    },
    "azurerm_vpn_site": {
        "virtual_wan_id": f'"{RG_SCOPE}/providers/Microsoft.Network/virtualWans/vwan-test"',
    },
    "azurerm_web_pubsub_hub": {
        "web_pubsub_id": f'"{RG_SCOPE}/providers/Microsoft.SignalRService/webPubSub/wps-test"',
    },
    "azurerm_network_connection_monitor": {
        "network_watcher_id": f'"{RG_SCOPE}/providers/Microsoft.Network/networkWatchers/nw-test"',
    },
    "azurerm_application_insights_web_test": {
        "kind": '"ping"',
        "geo_locations": '["us-tx-sn1-azr"]',
        "configuration": '"<WebTest><Items><Request Method=\\\"GET\\\" Url=\\\"https://example.com\\\"/></Items></WebTest>"',
        "application_insights_id": f'"{RG_SCOPE}/providers/Microsoft.Insights/components/ai-test"',
    },
    "azurerm_data_protection_backup_policy_disk": {
        "vault_id": f'"{RG_SCOPE}/providers/Microsoft.DataProtection/backupVaults/bv-test"',
        "default_retention_duration": '"P7D"',
        "backup_repeating_time_intervals": '["R/2021-05-19T06:33:16+00:00/PT4H"]',
    },
    "azurerm_dev_center_catalog": {
        "dev_center_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test"',
    },
    "azurerm_dev_center_dev_box_definition": {
        "dev_center_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test"',
        "image_reference_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test/galleries/Default/images/img-test"',
        "sku_name": '"general_i_8c32gb256ssd_v2"',
    },
    "azurerm_dev_center_environment_type": {
        "dev_center_id": f'"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test"',
    },
}

# Extra raw HCL lines to inject into resource blocks for constraints
# that cannot be expressed via attribute overrides alone.
RESOURCE_EXTRA_HCL: dict[str, list[str]] = {}

# Block-specific attr overrides: {block_name: {attr_name: value}}
# Used when a generic ATTR_OVERRIDES value is wrong inside a specific block type.
BLOCK_ATTR_OVERRIDES: dict[str, dict[str, str]] = {
    "identity": {"type": '"SystemAssigned"'},
    "sku": {"name": '"Standard"', "tier": '"Standard"', "family": '"C"', "capacity": "1"},
    "ip_configuration": {"name": '"internal"'},
    "frontend_ip_configuration": {"name": '"internal"'},
    "site_config": {"type": '"Default"'},
    "serialization": {"type": '"Json"', "format": '"LineSeparated"'},
    "head_node": {"vm_size": '"Standard_D3_V2"'},
    "worker_node": {"vm_size": '"Standard_D3_V2"'},
    "zookeeper_node": {"vm_size": '"Standard_D3_V2"'},
    "os_disk": {"caching": '"ReadWrite"', "storage_account_type": '"Standard_LRS"'},
    "action": {"type": '"LogicApp"'},
    "managed_rule_set": {"type": '"OWASP"', "version": '"3.2"'},
    "life_cycle": {"data_store_type": '"VaultStore"'},
    "monitor_config": {"protocol": '"HTTPS"', "port": "443"},
    "backend_pool_health_probe": {"protocol": '"Https"'},
    "input": {"type": '"any"'},
    "output": {"type": '"any"'},
    "api_server_profile": {"visibility": '"Public"'},
    "ingress_profile": {"visibility": '"Public"'},
    "master_profile": {"vm_size": '"Standard_D8s_v3"'},
    "worker_profile": {"vm_size": '"Standard_D4s_v3"'},
    "default_node_pool": {"type": '"VirtualMachineScaleSets"', "vm_size": '"Standard_DS2_v2"'},
    "source": {"event_source": '"Alerts"'},
    "namespace_junction": {"target_path": '"nfs"', "nfs_export": '"/export"', "namespace_path": '"/nfs"'},
}

# Common parent-id / contextual attributes that map cleanly to a fake ARM id.
FAKE_IDS: dict[str, str] = {
    "resource_group_id":               RG_SCOPE,
    "key_vault_id":                    f"{RG_SCOPE}/providers/Microsoft.KeyVault/vaults/kv-test",
    "iothub_id":                       f"{RG_SCOPE}/providers/Microsoft.Devices/iotHubs/iot-test",
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
    "service_plan_id":                 f"{RG_SCOPE}/providers/Microsoft.Web/serverFarms/asp-test",
    "app_service_id":                  f"{RG_SCOPE}/providers/Microsoft.Web/sites/app-test",
    "app_service_name":                "app-test",
    "app_service_plan_id":             f"{RG_SCOPE}/providers/Microsoft.Web/serverFarms/asp-test",
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
    "cdn_frontdoor_profile_id":         f"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test",
    "cdn_frontdoor_origin_group_id":   f"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test/originGroups/og-test",
    "cdn_frontdoor_endpoint_id":       f"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test/afdEndpoints/ep-test",
    "cdn_frontdoor_rule_set_id":       f"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test/ruleSets/rs-test",
    "api_management_id":               f"{RG_SCOPE}/providers/Microsoft.ApiManagement/service/apim-test",
    "api_management_name":             "apim-test",
    "api_id":                          f"{RG_SCOPE}/providers/Microsoft.ApiManagement/service/apim-test/apis/api-test",
    "container_app_environment_id":    f"{RG_SCOPE}/providers/Microsoft.App/managedEnvironments/cae-test",
    "cognitive_account_id":            f"{RG_SCOPE}/providers/Microsoft.CognitiveServices/accounts/cog-test",
    "cognitive_deployment_id":         f"{RG_SCOPE}/providers/Microsoft.CognitiveServices/accounts/cog-test/deployments/dep-test",
    "digital_twins_id":                f"{RG_SCOPE}/providers/Microsoft.DigitalTwins/digitalTwinsInstances/dt-test",
    "dev_center_id":                   f"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test",
    "dev_center_project_id":           f"{RG_SCOPE}/providers/Microsoft.DevCenter/projects/proj-test",
    "dedicated_host_group_id":         f"{RG_SCOPE}/providers/Microsoft.Compute/hostGroups/hg-test",
    "healthcare_workspace_id":         f"{RG_SCOPE}/providers/Microsoft.HealthcareApis/workspaces/hcw-test",
    "loadbalancer_id":                 f"{RG_SCOPE}/providers/Microsoft.Network/loadBalancers/lb-test",
    "machine_learning_workspace_id":   f"{RG_SCOPE}/providers/Microsoft.MachineLearningServices/workspaces/mlw-test",
    "logic_app_id":                    f"{RG_SCOPE}/providers/Microsoft.Logic/workflows/la-test",
    "network_watcher_id":              f"{RG_SCOPE}/providers/Microsoft.Network/networkWatchers/nw-test",
    "dns_resolver_id":                 f"{RG_SCOPE}/providers/Microsoft.Network/dnsResolvers/dnsr-test",
    "dns_forwarding_ruleset_id":       f"{RG_SCOPE}/providers/Microsoft.Network/dnsForwardingRulesets/frs-test",
    "outbound_endpoint_id":            f"{RG_SCOPE}/providers/Microsoft.Network/dnsResolvers/dnsr-test/outboundEndpoints/oep-test",
    "private_dns_resolver_id":         f"{RG_SCOPE}/providers/Microsoft.Network/dnsResolvers/dnsr-test",
    "service_bus_namespace_id":        f"{RG_SCOPE}/providers/Microsoft.ServiceBus/namespaces/sb-test",
    "eventhub_namespace_id":           f"{RG_SCOPE}/providers/Microsoft.EventHub/namespaces/ehn-test",
    "eventhub_name":                   "eh-test",
    "eventhub_id":                     f"{RG_SCOPE}/providers/Microsoft.EventHub/namespaces/ehn-test/eventhubs/eh-test",
    "data_protection_backup_vault_id": f"{RG_SCOPE}/providers/Microsoft.DataProtection/backupVaults/bv-test",
    "backup_vault_id":                 f"{RG_SCOPE}/providers/Microsoft.DataProtection/backupVaults/bv-test",
    "synapse_workspace_id":            f"{RG_SCOPE}/providers/Microsoft.Synapse/workspaces/sw-test",
    "iothub_dps_id":                   f"{RG_SCOPE}/providers/Microsoft.Devices/provisioningServices/dps-test",
    "network_interface_id":            f"{RG_SCOPE}/providers/Microsoft.Network/networkInterfaces/nic-test",
    "public_ip_address_id":            f"{RG_SCOPE}/providers/Microsoft.Network/publicIPAddresses/pip-test",
    "application_gateway_id":          f"{RG_SCOPE}/providers/Microsoft.Network/applicationGateways/agw-test",
    "managed_environment_id":          f"{RG_SCOPE}/providers/Microsoft.App/managedEnvironments/cae-test",
    "redis_cache_id":                  f"{RG_SCOPE}/providers/Microsoft.Cache/redis/redis-test",
    "search_service_id":               f"{RG_SCOPE}/providers/Microsoft.Search/searchServices/srch-test",
    "signalr_service_id":              f"{RG_SCOPE}/providers/Microsoft.SignalRService/SignalR/sigr-test",
    "web_pubsub_id":                   f"{RG_SCOPE}/providers/Microsoft.SignalRService/WebPubSub/wps-test",
    "blueprint_id":                    f"/providers/Microsoft.Management/managementGroups/mg-test/providers/Microsoft.Blueprint/blueprints/bp-test/versions/1.0",
    "user_assigned_identity_id":       f"{RG_SCOPE}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/uai-test",
    "frontdoor_id":                    f"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnfd-test",
    "host_pool_id":                    f"{RG_SCOPE}/providers/Microsoft.DesktopVirtualization/hostPools/hp-test",
    "virtual_desktop_host_pool_id":    f"{RG_SCOPE}/providers/Microsoft.DesktopVirtualization/hostPools/hp-test",
    "batch_account_id":                f"{RG_SCOPE}/providers/Microsoft.Batch/batchAccounts/ba-test",
    "spring_cloud_id":                 f"{RG_SCOPE}/providers/Microsoft.AppPlatform/spring/spc-test",
    "server_id":                       f"{RG_SCOPE}/providers/Microsoft.Sql/servers/sql-test",
    "mssql_server_id":                 f"{RG_SCOPE}/providers/Microsoft.Sql/servers/sql-test",
    "sql_server_id":                   f"{RG_SCOPE}/providers/Microsoft.Sql/servers/sql-test",
    "dev_center_id":                   f"{RG_SCOPE}/providers/Microsoft.DevCenter/devCenters/dc-test",
    "vpn_gateway_id":                  f"{RG_SCOPE}/providers/Microsoft.Network/vpnGateways/vpng-test",
    "virtual_wan_id":                  f"{RG_SCOPE}/providers/Microsoft.Network/virtualWans/vwan-test",
    "web_pubsub_id":                   f"{RG_SCOPE}/providers/Microsoft.SignalRService/webPubSub/wps-test",
    "private_cloud_id":                f"{RG_SCOPE}/providers/Microsoft.AVS/privateClouds/avs-test",
    "vault_id":                        f"{RG_SCOPE}/providers/Microsoft.DataProtection/backupVaults/bv-test",
    "application_insights_id":         f"{RG_SCOPE}/providers/Microsoft.Insights/components/ai-test",
    "network_watcher_id":              f"{RG_SCOPE}/providers/Microsoft.Network/networkWatchers/nw-test",
    "cdn_frontdoor_profile_id":        f"{RG_SCOPE}/providers/Microsoft.Cdn/profiles/cdnp-test",
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
    # Block-level overrides (lower priority than resource-specific)
    blk_overrides = BLOCK_ATTR_OVERRIDES.get(block_name, {})
    for attr_name, attr_def in block.get("attributes", {}).items():
        if attr_def.get("required"):
            # Priority: resource-specific override > block override > generic fake_value_for
            res_ov = RESOURCE_ATTR_OVERRIDES.get(resource_type, {}).get(attr_name)
            if res_ov is not None:
                lines.append(f"{sp}  {attr_name} = {res_ov}")
            elif attr_name in blk_overrides:
                lines.append(f"{sp}  {attr_name} = {blk_overrides[attr_name]}")
            else:
                t = attr_def.get("type")
                lines.append(f"{sp}  {attr_name} = {fake_value_for(attr_name, t, attr_def, resource_type)}")
        elif attr_def.get("optional") and _override_lookup(attr_name, resource_type) is not None:
            # Emit optional attrs that have explicit overrides (for "at least one of" constraints)
            # Skip computed-only attrs (computed=true without optional=true)
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
        elif attr_def.get("optional") and resource_type in RESOURCE_ATTR_OVERRIDES:
            # Emit optional attrs with resource-specific overrides (for "at least one of" constraints)
            # Skip computed-only attrs (computed=true without optional=true)
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
    # For HCL: backslashes in regex need to be escaped once for HCL string literals.
    # The JSON source already has them as \\, which is correct for HCL.
    regex_lit = f'"{regex}"'
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
