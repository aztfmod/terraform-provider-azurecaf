# Azure CAF Terraform Provider

[![Terraform](https://img.shields.io/badge/terraform-%235835CC.svg?style=for-the-badge&logo=terraform&logoColor=white)](https://registry.terraform.io/providers/aztfmod/azurecaf/latest)
[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

> :information_source: This solution is offered and supported by the Open-Source community

## Overview

The Azure CAF (Cloud Adoption Framework) provider is a *logical provider* that operates entirely within Terraform's logic without interacting with external services. It provides helper methods for implementing Azure landing zones using Terraform with consistent, compliant resource naming.

## Key Features

- **🏗️ Generate compliant Azure resource names** following CAF guidelines and Azure naming restrictions
- **🧹 Clean and sanitize inputs** to ensure compliance with allowed patterns for each Azure resource type
- **🎲 Add random characters** for uniqueness when required
- **🏷️ Handle prefixes and suffixes** (manual or CAF-compliant)
- **✅ Validate existing names** using passthrough mode
- **🔄 Generate one or multiple related names** from the same naming inputs
- **📋 Support 520 Azure resource types** with accurate validation rules

## Quick Start

### Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = "~> 2.0"
    }
  }
}

provider "azurecaf" {
  # Configuration options
}
```

### Basic Example

```hcl
# Data source (recommended - evaluated at plan time)
data "azurecaf_name" "example" {
  name          = "myproject"
  resource_type = "azurerm_resource_group"
  prefixes      = ["prod"]
  suffixes      = ["001"]
  random_length = 5
  random_seed   = 12345
  clean_input   = true
}

resource "azurerm_resource_group" "example" {
  name     = data.azurecaf_name.example.result
  location = "East US"
}

# Output: "rg-prod-myproject-001-a1b2c"
```

## Provider Components

The Azure CAF provider includes:

### Resources
- **[azurecaf_name](resources/azurecaf_name.md)** - Generate Azure-compliant resource names (recommended)
- **[v2.0.0 migration guide](migration-v2.md)** - Breaking changes and upgrade steps

### Data Sources
- **[azurecaf_name](data-sources/azurecaf_name.md)** - Generate names at plan time (recommended approach)
- **[azurecaf_environment_variable](data-sources/azurecaf_environment_variable.md)** - Read environment variables securely

## Migration Guide

Version 2.0.0 removes `azurecaf_naming_convention` and updates CAF abbreviations and naming constraints. Complete the [v2.0.0 migration guide](migration-v2.md) before upgrading.

## Supported Azure Resource Types

The provider supports **520 Azure resource types** with accurate naming validation rules. Each resource type has specific constraints for:

- **Length requirements** (minimum and maximum)
- **Character restrictions** (allowed patterns)
- **Case sensitivity** requirements
- **Uniqueness scope** (global, resource group, or parent resource)

### Popular Resource Types

| Resource Type | Slug | Min | Max | Example Generated Name |
|---------------|------|-----|-----|----------------------|
| `azurerm_resource_group` | `rg` | 1 | 90 | `rg-prod-myapp-001` |
| `azurerm_storage_account` | `st` | 3 | 24 | `stprodmyapp001` |
| `azurerm_key_vault` | `kv` | 3 | 24 | `kv-prod-myapp-001` |
| `azurerm_linux_web_app` | `lwapp` | 2 | 60 | `lwapp-prod-myapp-001` |
| `azurerm_kubernetes_cluster` | `aks` | 1 | 63 | `aks-prod-myapp-001` |
| `azurerm_virtual_machine` | `vm` | 1 | 15 | `vm-prod-001` |
| `azurerm_sql_server` | `sql` | 1 | 63 | `sql-prod-myapp-001` |

<details>
<summary>📋 View Complete Resource Type List</summary>

### Complete Supported Resource Types

| Resource type | Slug | Min | Max | Scope | Example generated name |
|---|---:|---:|---:|---|---|
| `aks_node_pool_linux` | `npl` | 1 | 12 | parent | `nplexample` |
| `aks_node_pool_windows` | `npw` | 1 | 6 | parent | `npwexa` |
| `azurerm_aadb2c_directory` | `aadb2c` | 1 | 75 | global | `aadb2c-example` |
| `azurerm_analysis_services_server` | `as` | 3 | 63 | resourceGroup | `asexample` |
| `azurerm_api_management` | `apim` | 1 | 50 | global | `apim-example` |
| `azurerm_api_management_api` | `apimapi` | 1 | 80 | global | `apimapi-example` |
| `azurerm_api_management_api_operation_tag` | `apimapiopt` | 1 | 80 | global | `apimapiopt-example` |
| `azurerm_api_management_api_version_set` | `apimvs` | 1 | 80 | parent | `apimvs-example` |
| `azurerm_api_management_authorization_server` | `apimas` | 1 | 80 | parent | `apimas-example` |
| `azurerm_api_management_backend` | `apimbe` | 1 | 80 | global | `apimbe-example` |
| `azurerm_api_management_certificate` | `apimcer` | 1 | 80 | global | `apimcer-example` |
| `azurerm_api_management_gateway` | `apimgw` | 1 | 80 | global | `apimgw-example` |
| `azurerm_api_management_group` | `apimgr` | 1 | 80 | global | `apimgr-example` |
| `azurerm_api_management_logger` | `apimlg` | 1 | 80 | global | `apimlg-example` |
| `azurerm_api_management_named_value` | `apimnv` | 1 | 80 | parent | `apimnv-example` |
| `azurerm_api_management_openid_connect_provider` | `apimoidc` | 1 | 80 | parent | `apimoidc-example` |
| `azurerm_api_management_service` | `apim` | 1 | 50 | global | `apim-example` |
| `azurerm_app_configuration` | `appcs` | 5 | 50 | resourceGroup | `appcs-example` |
| `azurerm_app_service` | `app` | 2 | 60 | global | `app-example` |
| `azurerm_app_service_certificate` | `appcert` | 1 | 80 | resourceGroup | `appcert-example` |
| `azurerm_app_service_certificate_order` | `appco` | 1 | 80 | resourceGroup | `appco-example` |
| `azurerm_app_service_environment` | `ase` | 2 | 36 | resourceGroup | `ase-example` |
| `azurerm_app_service_slot` | `apps` | 1 | 59 | parent | `apps-example` |
| `azurerm_application_gateway` | `agw` | 1 | 80 | resourceGroup | `agw-example` |
| `azurerm_application_insights` | `appi` | 1 | 260 | resourceGroup | `appi-example` |
| `azurerm_application_insights_analytics_item` | `appiai` | 1 | 64 | parent | `appiai-example` |
| `azurerm_application_insights_api_key` | `appiak` | 1 | 64 | parent | `appiak-example` |
| `azurerm_application_insights_web_test` | `appiwt` | 1 | 64 | resourceGroup | `appiwt-example` |
| `azurerm_application_security_group` | `asg` | 1 | 80 | resourceGroup | `asg-example` |
| `azurerm_automation_account` | `aa` | 6 | 50 | resourceGroup | `aa-example` |
| `azurerm_automation_certificate` | `aacert` | 1 | 128 | parent | `aacert-example` |
| `azurerm_automation_connection` | `aacon` | 1 | 128 | parent | `aacon-example` |
| `azurerm_automation_connection_certificate` | `aaconcrt` | 1 | 128 | parent | `aaconcrt-example` |
| `azurerm_automation_connection_classic_certificate` | `aaconcc` | 1 | 128 | parent | `aaconcc-example` |
| `azurerm_automation_connection_service_principal` | `aaconsp` | 1 | 128 | parent | `aaconsp-example` |
| `azurerm_automation_credential` | `aacred` | 1 | 128 | parent | `aacred-example` |
| `azurerm_automation_dsc_configuration` | `aadsc` | 1 | 128 | parent | `aadscexample` |
| `azurerm_automation_dsc_nodeconfiguration` | `aadscnc` | 1 | 128 | parent | `aadscnc-example` |
| `azurerm_automation_hybrid_runbook_worker_group` | `aahwg` | 1 | 128 | parent | `aahwg-example` |
| `azurerm_automation_job_schedule` | `aajs` | 1 | 128 | parent | `aajs-example` |
| `azurerm_automation_module` | `aamod` | 1 | 128 | parent | `aamod-example` |
| `azurerm_automation_runbook` | `aarun` | 1 | 63 | parent | `aarunexample` |
| `azurerm_automation_schedule` | `aasched` | 1 | 128 | parent | `aasched-example` |
| `azurerm_automation_variable` | `aavar` | 1 | 128 | parent | `aavar-example` |
| `azurerm_automation_variable_bool` | `aavb` | 1 | 128 | parent | `aavb-example` |
| `azurerm_automation_variable_datetime` | `aavdt` | 1 | 128 | parent | `aavdt-example` |
| `azurerm_automation_variable_int` | `aavi` | 1 | 128 | parent | `aavi-example` |
| `azurerm_automation_variable_string` | `aavs` | 1 | 128 | parent | `aavs-example` |
| `azurerm_availability_set` | `avail` | 1 | 80 | resourceGroup | `avail-example` |
| `azurerm_backup_policy_file_share` | `bkpfs` | 1 | 150 | parent | `bkpfs-example` |
| `azurerm_backup_policy_vm` | `bkpvm` | 1 | 150 | parent | `bkpvm-example` |
| `azurerm_bastion_host` | `bas` | 1 | 80 | parent | `bas-example` |
| `azurerm_batch_account` | `ba` | 3 | 24 | region | `baexample` |
| `azurerm_batch_application` | `baapp` | 1 | 64 | parent | `baapp-example` |
| `azurerm_batch_certificate` | `bacert` | 5 | 45 | parent | `bacert-example` |
| `azurerm_batch_pool` | `bapool` | 3 | 24 | parent | `bapool-example` |
| `azurerm_blueprint_assignment` | `bpa` | 1 | 90 | parent | `bpa-example` |
| `azurerm_bot_channel_Email` | `botmail` | 2 | 64 | parent | `botmail-example` |
| `azurerm_bot_channel_directline` | `botline` | 2 | 64 | parent | `botline-example` |
| `azurerm_bot_channel_ms_teams` | `botteams` | 2 | 64 | parent | `botteams-example` |
| `azurerm_bot_channel_slack` | `botslack` | 2 | 64 | parent | `botslack-example` |
| `azurerm_bot_channels_registration` | `botchan` | 2 | 64 | parent | `botchan-example` |
| `azurerm_bot_connection` | `botcon` | 2 | 64 | parent | `botcon-example` |
| `azurerm_bot_service_azure_bot` | `botaz` | 2 | 64 | global | `botaz-example` |
| `azurerm_bot_web_app` | `bot` | 2 | 64 | global | `bot-example` |
| `azurerm_cdn_endpoint` | `cdne` | 1 | 50 | global | `cdne-example` |
| `azurerm_cdn_frontdoor_custom_domain` | `cfdcd` | 1 | 260 | parent | `cfdcd-example` |
| `azurerm_cdn_frontdoor_endpoint` | `cfde` | 1 | 46 | global | `cfde-example` |
| `azurerm_cdn_frontdoor_firewall_policy` | `cfdfp` | 1 | 128 | resourceGroup | `cfdfpexample` |
| `azurerm_cdn_frontdoor_origin` | `cfdo` | 1 | 90 | parent | `cfdo-example` |
| `azurerm_cdn_frontdoor_origin_group` | `cfdog` | 1 | 90 | parent | `cfdog-example` |
| `azurerm_cdn_frontdoor_profile` | `cfdp` | 1 | 90 | resourceGroup | `cfdp-example` |
| `azurerm_cdn_frontdoor_route` | `cfdroute` | 1 | 90 | parent | `cfdroute-example` |
| `azurerm_cdn_frontdoor_rule` | `cfdr` | 1 | 60 | parent | `cfdrexample` |
| `azurerm_cdn_frontdoor_rule_set` | `cfdrs` | 1 | 60 | parent | `cfdrsexample` |
| `azurerm_cdn_frontdoor_secret` | `cfds` | 2 | 260 | parent | `cfds-example` |
| `azurerm_cdn_frontdoor_security_policy` | `cfdsp` | 1 | 260 | parent | `cfdsp-example` |
| `azurerm_cdn_profile` | `cdnp` | 1 | 260 | resourceGroup | `cdnp-example` |
| `azurerm_cognitive_account` | `cog` | 2 | 64 | resourceGroup | `cog-example` |
| `azurerm_cognitive_deployment` | `cog` | 2 | 64 | resourceGroup | `cog-example` |
| `azurerm_communication_service` | `acs` | 1 | 64 | parent | `acs-example` |
| `azurerm_consumption_budget_resource_group` | `acbrg` | 1 | 63 | resourceGroup | `acbrg-example` |
| `azurerm_consumption_budget_subscription` | `acbs` | 1 | 63 | subscription | `acbs-example` |
| `azurerm_containerGroups` | `cg` | 1 | 63 | resourceGroup | `cg-example` |
| `azurerm_container_app` | `ca` | 1 | 32 | resourceGroup | `ca-example` |
| `azurerm_container_app_environment` | `cae` | 1 | 60 | resourceGroup | `cae-example` |
| `azurerm_container_group` | `ci` | 1 | 63 | resourceGroup | `ci-example` |
| `azurerm_container_registry` | `cr` | 1 | 63 | resourceGroup | `crexample` |
| `azurerm_container_registry_webhook` | `crwh` | 1 | 50 | resourceGroup | `crwhexample` |
| `azurerm_cosmosdb_account` | `cosmos` | 1 | 44 | resourceGroup | `cosmos-example` |
| `azurerm_cosmosdb_cassandra_keyspace` | `coscas` | 1 | 255 | parent | `coscas-example` |
| `azurerm_cosmosdb_gremlin_database` | `cosgrm` | 1 | 255 | parent | `cosgrm-example` |
| `azurerm_cosmosdb_gremlin_graph` | `cosgrmg` | 1 | 255 | parent | `cosgrmg-example` |
| `azurerm_cosmosdb_mongo_collection` | `cosmonc` | 1 | 255 | parent | `cosmonc-example` |
| `azurerm_cosmosdb_mongo_database` | `cosmondb` | 1 | 255 | parent | `cosmondb-example` |
| `azurerm_cosmosdb_sql_container` | `cosqlc` | 1 | 255 | parent | `cosqlc-example` |
| `azurerm_cosmosdb_sql_database` | `cosqldb` | 1 | 255 | parent | `cosqldb-example` |
| `azurerm_cosmosdb_sql_stored_procedure` | `cosqlsp` | 1 | 255 | parent | `cosqlsp-example` |
| `azurerm_cosmosdb_table` | `costbl` | 1 | 255 | parent | `costbl-example` |
| `azurerm_custom_provider` | `prov` | 3 | 64 | resourceGroup | `prov-example` |
| `azurerm_dashboard` | `dsb` | 3 | 160 | parent | `dsb-example` |
| `azurerm_data_factory` | `adf` | 3 | 63 | global | `adf-example` |
| `azurerm_data_factory_dataset_azure_blob` | `adfblob` | 1 | 260 | parent | `adfblobexample` |
| `azurerm_data_factory_dataset_cosmosdb_sqlapi` | `adfsqlapi` | 1 | 260 | parent | `adfsqlapiexample` |
| `azurerm_data_factory_dataset_delimited_text` | `adfdtext` | 1 | 260 | parent | `adfdtextexample` |
| `azurerm_data_factory_dataset_http` | `adfhttp` | 1 | 260 | parent | `adfhttpexample` |
| `azurerm_data_factory_dataset_json` | `adfjson` | 1 | 260 | parent | `adfjsonexample` |
| `azurerm_data_factory_dataset_mysql` | `adfmysql` | 1 | 260 | parent | `adfmysqlexample` |
| `azurerm_data_factory_dataset_postgresql` | `adfpsql` | 1 | 260 | parent | `adfpsqlexample` |
| `azurerm_data_factory_dataset_sql_server_table` | `adfmssql` | 1 | 260 | parent | `adfmssqlexample` |
| `azurerm_data_factory_integration_runtime_managed` | `adfir` | 3 | 63 | parent | `adfir-example` |
| `azurerm_data_factory_integration_runtime_self_hosted` | `adfirsh` | 3 | 63 | parent | `adfirsh-example` |
| `azurerm_data_factory_linked_service_azure_blob_storage` | `adflsabs` | 1 | 260 | parent | `adflsabsexample` |
| `azurerm_data_factory_linked_service_azure_databricks` | `adflsadb` | 1 | 260 | parent | `adflsadbexample` |
| `azurerm_data_factory_linked_service_azure_file_storage` | `adflsafs` | 1 | 260 | parent | `adflsafsexample` |
| `azurerm_data_factory_linked_service_azure_function` | `adflsaf` | 1 | 260 | parent | `adflsafexample` |
| `azurerm_data_factory_linked_service_azure_sql_database` | `adflsasdb` | 1 | 260 | parent | `adflsasdbexample` |
| `azurerm_data_factory_linked_service_cosmosdb` | `adflsacdb` | 1 | 260 | parent | `adflsacdbexample` |
| `azurerm_data_factory_linked_service_data_lake_storage_gen2` | `adfsvst` | 1 | 260 | parent | `adfsvstexample` |
| `azurerm_data_factory_linked_service_key_vault` | `adfsvkv` | 1 | 260 | parent | `adfsvkvexample` |
| `azurerm_data_factory_linked_service_mysql` | `adfsvmysql` | 1 | 260 | parent | `adfsvmysqlexample` |
| `azurerm_data_factory_linked_service_postgresql` | `adfsvpsql` | 1 | 260 | parent | `adfsvpsqlexample` |
| `azurerm_data_factory_linked_service_sftp` | `adflsaftp` | 1 | 260 | parent | `adflsaftpexample` |
| `azurerm_data_factory_linked_service_sql_server` | `adfsvmssql` | 1 | 260 | parent | `adfsvmssqlexample` |
| `azurerm_data_factory_linked_service_web` | `adfsvweb` | 1 | 260 | parent | `adfsvwebexample` |
| `azurerm_data_factory_pipeline` | `adfpl` | 1 | 260 | parent | `adfplexample` |
| `azurerm_data_factory_trigger_schedule` | `adftg` | 1 | 260 | parent | `adftgexample` |
| `azurerm_data_lake_analytics_account` | `dla` | 3 | 24 | global | `dlaexample` |
| `azurerm_data_lake_analytics_firewall_rule` | `dlfw` | 3 | 50 | parent | `dlfw-example` |
| `azurerm_data_lake_store` | `dls` | 3 | 24 | parent | `dlsexample` |
| `azurerm_data_lake_store_firewall_rule` | `dlsfw` | 3 | 50 | parent | `dlsfw-example` |
| `azurerm_data_protection_backup_policy_blob_storage` | `dpbpb` | 3 | 150 | resourceGroup | `dpbpb-example` |
| `azurerm_data_protection_backup_policy_disk` | `dpbpd` | 3 | 150 | resourceGroup | `dpbpd-example` |
| `azurerm_data_protection_backup_policy_postgresql` | `dpbpp` | 3 | 150 | resourceGroup | `dpbpp-example` |
| `azurerm_data_protection_backup_policy_postgresql_flexible_server` | `dpbppf` | 3 | 150 | resourceGroup | `dpbppf-example` |
| `azurerm_data_protection_backup_vault` | `dpbv` | 2 | 50 | resourceGroup | `dpbv-example` |
| `azurerm_data_share` | `dshr` | 2 | 90 | parent | `dshr-example` |
| `azurerm_data_share_account` | `dshra` | 2 | 90 | resourceGroup | `dshra-example` |
| `azurerm_data_share_dataset_blob_storage` | `dshrdsb` | 2 | 90 | parent | `dshrdsb-example` |
| `azurerm_data_share_dataset_data_lake_gen1` | `dshrdsg1` | 2 | 90 | parent | `dshrdsg1-example` |
| `azurerm_data_share_dataset_data_lake_gen2` | `dshrdsg2` | 2 | 90 | parent | `dshrdsg2-example` |
| `azurerm_data_share_dataset_kusto_cluster` | `dshrdskc` | 2 | 90 | parent | `dshrdskc-example` |
| `azurerm_data_share_dataset_kusto_database` | `dshrdskd` | 2 | 90 | parent | `dshrdskd-example` |
| `azurerm_database_migration_project` | `migr` | 2 | 57 | parent | `migr-example` |
| `azurerm_database_migration_service` | `dms` | 2 | 62 | resourceGroup | `dms-example` |
| `azurerm_databricks_workspace` | `dbw` | 3 | 64 | resourceGroup | `dbw-example` |
| `azurerm_dedicated_hardware_security_module` | `hsm` | 3 | 24 | resourceGroup | `hsm-example` |
| `azurerm_dedicated_host` | `dh` | 1 | 80 | resourceGroup | `dh-example` |
| `azurerm_dedicated_host_group` | `dhg` | 1 | 80 | resourceGroup | `dhg-example` |
| `azurerm_dev_center` | `dc` | 3 | 26 | resourceGroup | `dc-example` |
| `azurerm_dev_center_catalog` | `dcc` | 3 | 63 | parent | `dcc-example` |
| `azurerm_dev_center_dev_box_definition` | `dcdb` | 3 | 63 | parent | `dcdb-example` |
| `azurerm_dev_center_environment_type` | `dcet` | 3 | 63 | parent | `dcet-example` |
| `azurerm_dev_center_gallery` | `dcg` | 1 | 80 | parent | `dcgexample` |
| `azurerm_dev_center_network_connection` | `dcnc` | 3 | 63 | resourceGroup | `dcnc-example` |
| `azurerm_dev_center_project` | `dcp` | 3 | 63 | resourceGroup | `dcp-example` |
| `azurerm_dev_center_project_environment_type` | `dcpet` | 3 | 63 | parent | `dcpet-example` |
| `azurerm_dev_test_lab` | `lab` | 1 | 50 | resourceGroup | `lab-example` |
| `azurerm_dev_test_linux_virtual_machine` | `labvm` | 1 | 64 | parent | `labvm-example` |
| `azurerm_dev_test_schedule` | `dtls` | 1 | 64 | parent | `dtls-example` |
| `azurerm_dev_test_virtual_network` | `dtlvn` | 1 | 64 | parent | `dtlvn-example` |
| `azurerm_dev_test_windows_virtual_machine` | `labvm` | 1 | 15 | parent | `labvm-example` |
| `azurerm_devspace_controller` | `dsc` | 1 | 63 | resourceGroup | `dsc-example` |
| `azurerm_digital_twins_endpoint_eventgrid` | `adteg` | 3 | 50 | parent | `adteg-example` |
| `azurerm_digital_twins_endpoint_eventhub` | `adteh` | 3 | 50 | parent | `adteh-example` |
| `azurerm_digital_twins_endpoint_servicebus` | `adtsb` | 3 | 50 | parent | `adtsb-example` |
| `azurerm_digital_twins_instance` | `adt` | 4 | 63 | subscription | `adt-example` |
| `azurerm_disk_encryption_set` | `des` | 1 | 80 | resourceGroup | `des-example` |
| `azurerm_dns_a_record` | `dnsrec` | 1 | 80 | parent | `dnsrec-example` |
| `azurerm_dns_aaaa_record` | `dnsrec` | 1 | 80 | parent | `dnsrec-example` |
| `azurerm_dns_caa_record` | `dnsrec` | 1 | 80 | parent | `dnsrec-example` |
| `azurerm_dns_cname_record` | `dnsrec` | 1 | 80 | parent | `dnsrec-example` |
| `azurerm_dns_mx_record` | `dnsrec` | 1 | 80 | parent | `dnsrec-example` |
| `azurerm_dns_ns_record` | `dnsrec` | 1 | 80 | parent | `dnsrec-example` |
| `azurerm_dns_ptr_record` | `dnsrec` | 1 | 80 | parent | `dnsrec-example` |
| `azurerm_dns_srv_record` | `dnsrec` | 1 | 80 | parent | `dnsrec-example` |
| `azurerm_dns_txt_record` | `dnsrec` | 1 | 80 | parent | `dnsrec-example` |
| `azurerm_dns_zone` | `dnsz` | 1 | 63 | resourceGroup | `dnsz-example` |
| `azurerm_email_communication_service` | `acsmail` | 1 | 63 | global | `acsmail-example` |
| `azurerm_eventgrid_domain` | `evgd` | 3 | 50 | resourceGroup | `evgd-example` |
| `azurerm_eventgrid_domain_topic` | `evgdt` | 3 | 50 | parent | `evgdt-example` |
| `azurerm_eventgrid_event_subscription` | `evgs` | 3 | 64 | resourceGroup | `evgs-example` |
| `azurerm_eventgrid_system_topic` | `egst` | 3 | 50 | resourceGroup | `egst-example` |
| `azurerm_eventgrid_topic` | `evgt` | 3 | 50 | resourceGroup | `evgt-example` |
| `azurerm_eventhub` | `evh` | 1 | 50 | parent | `evh-example` |
| `azurerm_eventhub_authorization_rule` | `evhar` | 1 | 50 | parent | `evhar-example` |
| `azurerm_eventhub_cluster` | `ehc` | 1 | 50 | resourceGroup | `ehc-example` |
| `azurerm_eventhub_consumer_group` | `evhnscg` | 1 | 50 | parent | `evhnscg-example` |
| `azurerm_eventhub_namespace` | `evhns` | 1 | 50 | global | `evhns-example` |
| `azurerm_eventhub_namespace_authorization_rule` | `evhnsar` | 1 | 50 | parent | `evhnsar-example` |
| `azurerm_eventhub_namespace_disaster_recovery_config` | `evhnsdr` | 1 | 50 | parent | `evhnsdr-example` |
| `azurerm_express_route_circuit` | `erc` | 1 | 80 | resourceGroup | `erc-example` |
| `azurerm_express_route_circuit_authorization` | `erca` | 1 | 80 | parent | `erca-example` |
| `azurerm_express_route_gateway` | `ergw` | 1 | 80 | resourceGroup | `ergw-example` |
| `azurerm_federated_identity_credential` | `fedcred` | 3 | 120 | parent | `fedcred-example` |
| `azurerm_firewall` | `afw` | 1 | 80 | resourceGroup | `afw-example` |
| `azurerm_firewall_application_rule_collection` | `fwapp` | 1 | 80 | parent | `fwapp-example` |
| `azurerm_firewall_ip_configuration` | `fwipconf` | 1 | 80 | resourceGroup | `fwipconf-example` |
| `azurerm_firewall_nat_rule_collection` | `fwnatrc` | 1 | 80 | parent | `fwnatrc-example` |
| `azurerm_firewall_network_rule_collection` | `fwnetrc` | 1 | 80 | parent | `fwnetrc-example` |
| `azurerm_firewall_policy` | `afwp` | 1 | 80 | resourceGroup | `afwp-example` |
| `azurerm_firewall_policy_rule_collection_group` | `fwprcg` | 1 | 80 | parent | `fwprcg-example` |
| `azurerm_frontdoor` | `afd` | 5 | 64 | global | `afd-example` |
| `azurerm_frontdoor_firewall_policy` | `fdfw` | 1 | 80 | global | `fdfwexample` |
| `azurerm_function_app` | `func` | 2 | 60 | global | `func-example` |
| `azurerm_function_app_slot` | `funcs` | 2 | 59 | global | `funcs-example` |
| `azurerm_hdinsight_hadoop_cluster` | `hadoop` | 3 | 59 | global | `hadoop-example` |
| `azurerm_hdinsight_hbase_cluster` | `hbase` | 3 | 59 | global | `hbase-example` |
| `azurerm_hdinsight_interactive_query_cluster` | `iqr` | 3 | 59 | global | `iqr-example` |
| `azurerm_hdinsight_kafka_cluster` | `kafka` | 3 | 59 | global | `kafka-example` |
| `azurerm_hdinsight_ml_services_cluster` | `mls` | 3 | 59 | global | `mls-example` |
| `azurerm_hdinsight_rserver_cluster` | `rser` | 3 | 59 | global | `rser-example` |
| `azurerm_hdinsight_spark_cluster` | `spark` | 3 | 59 | global | `spark-example` |
| `azurerm_hdinsight_storm_cluster` | `storm` | 3 | 59 | global | `storm-example` |
| `azurerm_healthcare_dicom_service` | `dicom` | 3 | 24 | parent | `dicom-example` |
| `azurerm_healthcare_fhir_service` | `fhir` | 3 | 24 | parent | `fhir-example` |
| `azurerm_healthcare_medtech_service` | `medtech` | 3 | 24 | parent | `medtech-example` |
| `azurerm_healthcare_service` | `hcasvc` | 3 | 24 | global | `hcasvc-example` |
| `azurerm_healthcare_workspace` | `hcw` | 3 | 24 | global | `hcwexample` |
| `azurerm_hpc_cache` | `hpcc` | 2 | 31 | resourceGroup | `hpcc-example` |
| `azurerm_hpc_cache_blob_target` | `hpcbt` | 2 | 31 | parent | `hpcbt-example` |
| `azurerm_hpc_cache_nfs_target` | `hpcnt` | 2 | 31 | parent | `hpcnt-example` |
| `azurerm_image` | `img` | 1 | 80 | resourceGroup | `img-example` |
| `azurerm_integration_service_environment` | `lappise` | 1 | 80 | resourceGroup | `lappise-example` |
| `azurerm_iot_security_device_group` | `iotdg` | 1 | 32 | parent | `iotdg-example` |
| `azurerm_iot_security_solution` | `iotss` | 1 | 260 | resourceGroup | `iotss-example` |
| `azurerm_iot_time_series_insights_reference_data_set` | `tsirds` | 3 | 63 | parent | `tsirdsexample` |
| `azurerm_iot_time_series_insights_standard_environment` | `tsise` | 1 | 90 | resourceGroup | `tsise-example` |
| `azurerm_iotcentral_application` | `iotapp` | 2 | 63 | global | `iotapp-example` |
| `azurerm_iothub` | `iot` | 3 | 50 | global | `iot-example` |
| `azurerm_iothub_certificate` | `iotcert` | 1 | 64 | parent | `iotcert-example` |
| `azurerm_iothub_consumer_group` | `iotcg` | 1 | 50 | parent | `iotcg-example` |
| `azurerm_iothub_dps` | `dps` | 3 | 64 | resourceGroup | `dps-example` |
| `azurerm_iothub_dps_certificate` | `dpscert` | 1 | 64 | parent | `dpscert-example` |
| `azurerm_iothub_dps_shared_access_policy` | `dpssap` | 1 | 64 | parent | `dpssap-example` |
| `azurerm_iothub_endpoint_eventhub` | `iothepeh` | 1 | 64 | parent | `iothepeh-example` |
| `azurerm_iothub_endpoint_servicebus_queue` | `iothepsbq` | 1 | 64 | parent | `iothepsbq-example` |
| `azurerm_iothub_endpoint_servicebus_topic` | `iothepsbt` | 1 | 64 | parent | `iothepsbt-example` |
| `azurerm_iothub_endpoint_storage_container` | `iothepsc` | 1 | 64 | parent | `iothepsc-example` |
| `azurerm_iothub_route` | `iothr` | 1 | 64 | parent | `iothr-example` |
| `azurerm_iothub_shared_access_policy` | `iotsap` | 1 | 64 | parent | `iotsap-example` |
| `azurerm_ip_group` | `ipgr` | 1 | 80 | resourceGroup | `ipgr-example` |
| `azurerm_key_vault` | `kv` | 3 | 24 | global | `kv-example` |
| `azurerm_key_vault_certificate` | `kvc` | 1 | 127 | parent | `kvc-example` |
| `azurerm_key_vault_certificate_issuer` | `kvci` | 1 | 127 | parent | `kvci-example` |
| `azurerm_key_vault_key` | `kvk` | 1 | 127 | parent | `kvk-example` |
| `azurerm_key_vault_secret` | `kvs` | 1 | 127 | parent | `kvs-example` |
| `azurerm_kubernetes_cluster` | `aks` | 1 | 63 | resourceGroup | `aks-example` |
| `azurerm_kubernetes_fleet_manager` | `fleet` | 1 | 63 | resourceGroup | `fleet-example` |
| `azurerm_kusto_attached_database_configuration` | `kadc` | 1 | 260 | parent | `kadc-example` |
| `azurerm_kusto_cluster` | `dec` | 4 | 22 | global | `decexample` |
| `azurerm_kusto_cluster_principal_assignment` | `kcpa` | 1 | 260 | parent | `kcpa-example` |
| `azurerm_kusto_database` | `dedb` | 1 | 260 | parent | `dedb-example` |
| `azurerm_kusto_database_principal_assignment` | `kdpa` | 1 | 260 | parent | `kdpa-example` |
| `azurerm_kusto_eventhub_data_connection` | `deevhdcon` | 1 | 40 | parent | `deevhdcon-example` |
| `azurerm_lb` | `lb` | 1 | 80 | resourceGroup | `lb-example` |
| `azurerm_lb_backend_address_pool` | `lbbepool` | 1 | 80 | parent | `lbbepool-example` |
| `azurerm_lb_nat_pool` | `lbnatpool` | 1 | 80 | parent | `lbnatpool-example` |
| `azurerm_lb_nat_rule` | `lbnatr` | 1 | 80 | parent | `lbnatr-example` |
| `azurerm_lb_outbound_rule` | `lbor` | 1 | 80 | parent | `lbor-example` |
| `azurerm_lb_probe` | `probe` | 1 | 80 | parent | `probe-example` |
| `azurerm_lb_rule` | `rule` | 1 | 80 | parent | `rule-example` |
| `azurerm_lighthouse_definition` | `lhd` | 1 | 80 | parent | `lhd-example` |
| `azurerm_linux_function_app` | `fa` | 2 | 60 | global | `fa-example` |
| `azurerm_linux_function_app_slot` | `fas` | 2 | 59 | global | `fas-example` |
| `azurerm_linux_virtual_machine` | `vm` | 1 | 64 | resourceGroup | `vm-example` |
| `azurerm_linux_virtual_machine_scale_set` | `vmss` | 1 | 64 | resourceGroup | `vmss-example` |
| `azurerm_linux_web_app` | `lwapp` | 2 | 60 | global | `lwapp-example` |
| `azurerm_load_test` | `load` | 1 | 64 | global | `load-example` |
| `azurerm_local_network_gateway` | `lgw` | 1 | 80 | resourceGroup | `lgw-example` |
| `azurerm_log_analytics_cluster` | `logc` | 4 | 63 | resourceGroup | `logc-example` |
| `azurerm_log_analytics_data_export_rule` | `laer` | 4 | 63 | parent | `laer-example` |
| `azurerm_log_analytics_datasource_windows_event` | `ladswe` | 1 | 63 | parent | `ladswe-example` |
| `azurerm_log_analytics_datasource_windows_performance_counter` | `ladswpc` | 1 | 63 | parent | `ladswpc-example` |
| `azurerm_log_analytics_query_pack` | `laqp` | 4 | 63 | parent | `laqp-example` |
| `azurerm_log_analytics_saved_search` | `lass` | 1 | 80 | parent | `lass-example` |
| `azurerm_log_analytics_solution` | `las` | 4 | 63 | parent | `las-example` |
| `azurerm_log_analytics_storage_insights` | `lasi` | 4 | 63 | parent | `lasi-example` |
| `azurerm_log_analytics_workspace` | `log` | 4 | 63 | parent | `log-example` |
| `azurerm_logic_app_action_custom` | `logicac` | 1 | 80 | resourceGroup | `logicac-example` |
| `azurerm_logic_app_action_http` | `logicah` | 1 | 80 | resourceGroup | `logicah-example` |
| `azurerm_logic_app_integration_account` | `ia` | 1 | 80 | resourceGroup | `ia-example` |
| `azurerm_logic_app_trigger_custom` | `logictc` | 1 | 80 | resourceGroup | `logictc-example` |
| `azurerm_logic_app_trigger_http_request` | `logicth` | 1 | 80 | resourceGroup | `logicth-example` |
| `azurerm_logic_app_trigger_recurrence` | `logictc` | 1 | 80 | resourceGroup | `logictc-example` |
| `azurerm_logic_app_workflow` | `logic` | 1 | 80 | resourceGroup | `logic-example` |
| `azurerm_machine_learning_compute_instance` | `amlci` | 1 | 16 | parent | `amlci-example` |
| `azurerm_machine_learning_workspace` | `mlw` | 1 | 260 | resourceGroup | `mlw-example` |
| `azurerm_maintenance_configuration` | `mcf` | 1 | 60 | resourceGroup | `mcf-example` |
| `azurerm_managed_application` | `manapp` | 1 | 64 | resourceGroup | `manapp-example` |
| `azurerm_managed_application_definition` | `manappd` | 3 | 64 | resourceGroup | `manappdexample` |
| `azurerm_managed_disk` | `disk` | 1 | 80 | resourceGroup | `disk-example` |
| `azurerm_managed_redis` | `amr` | 3 | 63 | resourceGroup | `amr-example` |
| `azurerm_management_group` | `mg` | 1 | 90 | parent | `mg-example` |
| `azurerm_management_lock` | `mgl` | 1 | 260 | parent | `mgl-example` |
| `azurerm_maps_account` | `map` | 1 | 98 | resourceGroup | `map-example` |
| `azurerm_mariadb_database` | `mariadb` | 1 | 63 | parent | `mariadb-example` |
| `azurerm_mariadb_firewall_rule` | `mariafw` | 1 | 128 | parent | `mariafw-example` |
| `azurerm_mariadb_server` | `maria` | 3 | 63 | global | `maria-example` |
| `azurerm_mariadb_virtual_network_rule` | `mariavn` | 1 | 128 | parent | `mariavn-example` |
| `azurerm_media_services_account` | `ams` | 3 | 24 | resourceGroup | `amsexample` |
| `azurerm_monitor_action_group` | `ag` | 1 | 260 | resourceGroup | `ag-example` |
| `azurerm_monitor_action_rule_action_group` | `marag` | 1 | 260 | resourceGroup | `marag-example` |
| `azurerm_monitor_action_rule_suppression` | `mars` | 1 | 260 | resourceGroup | `mars-example` |
| `azurerm_monitor_activity_log_alert` | `amala` | 1 | 260 | parent | `amalaexample` |
| `azurerm_monitor_autoscale_setting` | `amas` | 2 | 64 | resourceGroup | `amas-example` |
| `azurerm_monitor_data_collection_endpoint` | `dce` | 3 | 44 | resourceGroup | `dce-example` |
| `azurerm_monitor_data_collection_rule` | `dcr` | 3 | 44 | resourceGroup | `dcr-example` |
| `azurerm_monitor_diagnostic_setting` | `amds` | 1 | 260 | parent | `amds-example` |
| `azurerm_monitor_log_profile` | `mlp` | 1 | 80 | parent | `mlp-example` |
| `azurerm_monitor_metric_alert` | `ma` | 1 | 251 | resourceGroup | `ma-example` |
| `azurerm_monitor_private_link_scope` | `ampls` | 1 | 255 | resourceGroup | `ampls-example` |
| `azurerm_monitor_scheduled_query_rules_alert` | `schqra` | 1 | 260 | resourceGroup | `schqra-example` |
| `azurerm_monitor_scheduled_query_rules_log` | `msqrl` | 1 | 260 | resourceGroup | `msqrl-example` |
| `azurerm_monitor_smart_detector_alert_rule` | `msdar` | 1 | 260 | resourceGroup | `msdar-example` |
| `azurerm_monitor_workspace` | `amw` | 4 | 63 | resourceGroup | `amw-example` |
| `azurerm_mssql_database` | `sqldb` | 1 | 128 | parent | `sqldb-example` |
| `azurerm_mssql_elasticpool` | `sqlep` | 1 | 128 | parent | `sqlep-example` |
| `azurerm_mssql_mi` | `sqlmi` | 1 | 63 | global | `sqlmi-example` |
| `azurerm_mssql_server` | `sql` | 1 | 63 | global | `sql-example` |
| `azurerm_mysql_database` | `mysqldb` | 1 | 63 | parent | `mysqldb-example` |
| `azurerm_mysql_firewall_rule` | `mysqlfw` | 1 | 128 | parent | `mysqlfw-example` |
| `azurerm_mysql_flexible_server` | `mysqlf` | 3 | 63 | global | `mysqlf-example` |
| `azurerm_mysql_flexible_server_database` | `mysqlfdb` | 1 | 63 | parent | `mysqlfdb-example` |
| `azurerm_mysql_flexible_server_firewall_rule` | `mysqlffw` | 1 | 128 | parent | `mysqlffw-example` |
| `azurerm_mysql_server` | `mysql` | 3 | 63 | global | `mysql-example` |
| `azurerm_mysql_virtual_network_rule` | `mysqlvn` | 1 | 128 | parent | `mysqlvn-example` |
| `azurerm_nat_gateway` | `ng` | 1 | 80 | resourceGroup | `ng-example` |
| `azurerm_netapp_account` | `ana` | 1 | 128 | resourceGroup | `ana-example` |
| `azurerm_netapp_pool` | `anp` | 1 | 63 | resourceGroup | `anp-example` |
| `azurerm_netapp_snapshot` | `ans` | 1 | 63 | resourceGroup | `ans-example` |
| `azurerm_netapp_volume` | `anv` | 1 | 63 | resourceGroup | `anv-example` |
| `azurerm_network_connection_monitor` | `cm` | 1 | 80 | parent | `cm-example` |
| `azurerm_network_ddos_protection_plan` | `ddospp` | 1 | 80 | parent | `ddospp-example` |
| `azurerm_network_interface` | `nic` | 1 | 80 | resourceGroup | `nic-example` |
| `azurerm_network_packet_capture` | `npc` | 1 | 80 | parent | `npc-example` |
| `azurerm_network_profile` | `npr` | 1 | 80 | resourceGroup | `npr-example` |
| `azurerm_network_security_group` | `nsg` | 1 | 80 | resourceGroup | `nsg-example` |
| `azurerm_network_security_group_rule` | `nsgsr` | 1 | 80 | parent | `nsgsr-example` |
| `azurerm_network_security_rule` | `nsgr` | 1 | 80 | parent | `nsgr-example` |
| `azurerm_network_watcher` | `nw` | 1 | 80 | resourceGroup | `nw-example` |
| `azurerm_network_watcher_flow_log` | `nwfl` | 1 | 80 | parent | `nwfl-example` |
| `azurerm_nginx_deployment` | `nginx` | 1 | 30 | resourceGroup | `nginx-example` |
| `azurerm_notification_hub` | `ntf` | 1 | 260 | parent | `ntf-example` |
| `azurerm_notification_hub_authorization_rule` | `ntfar` | 1 | 256 | parent | `ntfar-example` |
| `azurerm_notification_hub_namespace` | `ntfns` | 6 | 50 | global | `ntfns-example` |
| `azurerm_orchestrated_virtual_machine_scale_set` | `ovmss` | 1 | 64 | resourceGroup | `ovmss-example` |
| `azurerm_point_to_site_vpn_gateway` | `vpngw` | 1 | 80 | resourceGroup | `vpngw-example` |
| `azurerm_policy_definition` | `pold` | 1 | 64 | parent | `pold-example` |
| `azurerm_policy_remediation` | `polr` | 1 | 64 | parent | `polr-example` |
| `azurerm_policy_set_definition` | `polsd` | 1 | 64 | parent | `polsd-example` |
| `azurerm_portal_dashboard` | `dsb` | 3 | 160 | parent | `dsb-example` |
| `azurerm_postgresql_database` | `psqldb` | 1 | 63 | parent | `psqldb-example` |
| `azurerm_postgresql_firewall_rule` | `psqlfw` | 1 | 128 | parent | `psqlfw-example` |
| `azurerm_postgresql_flexible_server` | `psqlf` | 3 | 63 | global | `psqlf-example` |
| `azurerm_postgresql_flexible_server_database` | `psqlfdb` | 1 | 63 | parent | `psqlfdb-example` |
| `azurerm_postgresql_flexible_server_firewall_rule` | `psqlffw` | 1 | 128 | parent | `psqlffw-example` |
| `azurerm_postgresql_server` | `psql` | 3 | 63 | global | `psql-example` |
| `azurerm_postgresql_virtual_network_rule` | `psqlvn` | 1 | 128 | parent | `psqlvn-example` |
| `azurerm_powerbi_embedded` | `pbi` | 3 | 63 | region | `pbiexample` |
| `azurerm_private_dns_a_record` | `pdnsrec` | 1 | 80 | parent | `pdnsrec-example` |
| `azurerm_private_dns_aaaa_record` | `pdnsrec` | 1 | 80 | parent | `pdnsrec-example` |
| `azurerm_private_dns_cname_record` | `pdnsrec` | 1 | 80 | parent | `pdnsrec-example` |
| `azurerm_private_dns_mx_record` | `pdnsrec` | 1 | 80 | parent | `pdnsrec-example` |
| `azurerm_private_dns_ptr_record` | `pdnsrec` | 1 | 80 | parent | `pdnsrec-example` |
| `azurerm_private_dns_resolver` | `dnspr` | 3 | 80 | resourceGroup | `dnspr-example` |
| `azurerm_private_dns_resolver_dns_forwarding_ruleset` | `dnsfwrs` | 2 | 80 | resourceGroup | `dnsfwrs-example` |
| `azurerm_private_dns_resolver_forwarding_rule` | `dnsfwr` | 1 | 80 | parent | `dnsfwr-example` |
| `azurerm_private_dns_resolver_inbound_endpoint` | `dnsprie` | 3 | 80 | parent | `dnsprie-example` |
| `azurerm_private_dns_resolver_outbound_endpoint` | `dnsproe` | 3 | 80 | parent | `dnsproe-example` |
| `azurerm_private_dns_resolver_virtual_network_link` | `dnsfwrsvnetl` | 1 | 80 | parent | `dnsfwrsvnetl-example` |
| `azurerm_private_dns_srv_record` | `pdnsrec` | 1 | 80 | parent | `pdnsrec-example` |
| `azurerm_private_dns_txt_record` | `pdnsrec` | 1 | 80 | parent | `pdnsrec-example` |
| `azurerm_private_dns_zone` | `pdnsz` | 1 | 63 | resourceGroup | `pdnsz-example` |
| `azurerm_private_dns_zone_group` | `pdnszg` | 1 | 80 | resourceGroup | `pdnszg-example` |
| `azurerm_private_dns_zone_virtual_network_link` | `pnetlk` | 1 | 80 | parent | `pnetlk-example` |
| `azurerm_private_endpoint` | `pe` | 2 | 64 | resourceGroup | `pe-example` |
| `azurerm_private_link_service` | `pl` | 2 | 64 | resourceGroup | `pl-example` |
| `azurerm_private_service_connection` | `psc` | 1 | 80 | resourceGroup | `psc-example` |
| `azurerm_proximity_placement_group` | `ppg` | 1 | 80 | resourceGroup | `ppg-example` |
| `azurerm_public_ip` | `pip` | 1 | 80 | parent | `pip-example` |
| `azurerm_public_ip_prefix` | `ippre` | 1 | 80 | parent | `ippre-example` |
| `azurerm_purview_account` | `pview` | 3 | 63 | subscription | `pview-example` |
| `azurerm_recovery_services_vault` | `rsv` | 2 | 50 | resourceGroup | `rsv-example` |
| `azurerm_recovery_services_vault_backup_police` | `rsvbp` | 3 | 150 | resourceGroup | `rsvbp-example` |
| `azurerm_redhat_openshift_cluster` | `aroc` | 1 | 30 | resourceGroup | `arocexample` |
| `azurerm_redhat_openshift_domain` | `arod` | 1 | 30 | resourceGroup | `arodexample` |
| `azurerm_redis_cache` | `redis` | 1 | 63 | global | `redis-example` |
| `azurerm_redis_firewall_rule` | `redisfw` | 1 | 256 | parent | `redisfwexample` |
| `azurerm_relay_hybrid_connection` | `rlhc` | 1 | 260 | parent | `rlhc-example` |
| `azurerm_relay_namespace` | `rln` | 6 | 50 | global | `rln-example` |
| `azurerm_resource_group` | `rg` | 1 | 90 | subscription | `rg-example` |
| `azurerm_resource_group_policy_assignment` | `argpa` | 1 | 128 | resourceGroup | `argpa-example` |
| `azurerm_resource_group_template_deployment` | `rgtd` | 1 | 64 | resourceGroup | `rgtd-example` |
| `azurerm_role_definition` | `rd` | 1 | 64 | definition | `rd-example` |
| `azurerm_route` | `udr` | 1 | 80 | parent | `udr-example` |
| `azurerm_route_filter` | `rf` | 1 | 80 | resourceGroup | `rf-example` |
| `azurerm_route_server` | `rts` | 1 | 80 | resourceGroup | `rts-example` |
| `azurerm_route_table` | `rt` | 1 | 80 | resourceGroup | `rt-example` |
| `azurerm_search_service` | `srch` | 2 | 60 | global | `srch-example` |
| `azurerm_security_center_automation` | `sca` | 1 | 80 | resourceGroup | `sca-example` |
| `azurerm_security_center_contact` | `scc` | 1 | 80 | parent | `scc-example` |
| `azurerm_sentinel_alert_rule_ms_security_incident` | `sentarms` | 1 | 80 | parent | `sentarms-example` |
| `azurerm_sentinel_alert_rule_scheduled` | `sentars` | 1 | 80 | parent | `sentars-example` |
| `azurerm_service_fabric_cluster` | `sf` | 4 | 23 | region | `sf-example` |
| `azurerm_service_fabric_mesh_application` | `sfmesha` | 1 | 80 | resourceGroup | `sfmesha-example` |
| `azurerm_service_fabric_mesh_local_network` | `sfmeshln` | 1 | 80 | resourceGroup | `sfmeshln-example` |
| `azurerm_service_fabric_mesh_secret` | `sfmeshs` | 1 | 80 | resourceGroup | `sfmeshs-example` |
| `azurerm_service_plan` | `asp` | 1 | 40 | resourceGroup | `asp-example` |
| `azurerm_servicebus_namespace` | `sbns` | 6 | 50 | global | `sbns-example` |
| `azurerm_servicebus_namespace_authorization_rule` | `sbar` | 1 | 50 | parent | `sbar-example` |
| `azurerm_servicebus_namespace_disaster_recovery_config` | `sbdr` | 1 | 50 | parent | `sbdr-example` |
| `azurerm_servicebus_queue` | `sbq` | 1 | 260 | parent | `sbq-example` |
| `azurerm_servicebus_queue_authorization_rule` | `sbqar` | 1 | 50 | parent | `sbqar-example` |
| `azurerm_servicebus_subscription` | `sbs` | 1 | 50 | parent | `sbs-example` |
| `azurerm_servicebus_subscription_rule` | `sbsr` | 1 | 50 | parent | `sbsr-example` |
| `azurerm_servicebus_topic` | `sbt` | 1 | 260 | parent | `sbt-example` |
| `azurerm_servicebus_topic_authorization_rule` | `sbtar` | 1 | 50 | parent | `sbtar-example` |
| `azurerm_shared_image` | `si` | 1 | 80 | parent | `si-example` |
| `azurerm_shared_image_gallery` | `gal` | 1 | 80 | resourceGroup | `galexample` |
| `azurerm_signalr_service` | `sigr` | 3 | 63 | global | `sigr-example` |
| `azurerm_site_recovery_fabric` | `asrf` | 1 | 128 | parent | `asrf-example` |
| `azurerm_site_recovery_network_mapping` | `asrnm` | 1 | 128 | parent | `asrnm-example` |
| `azurerm_site_recovery_protection_container` | `asrpc` | 1 | 128 | parent | `asrpc-example` |
| `azurerm_site_recovery_protection_container_mapping` | `asrpcm` | 1 | 128 | parent | `asrpcm-example` |
| `azurerm_site_recovery_replicated_vm` | `asrrvm` | 1 | 128 | parent | `asrrvm-example` |
| `azurerm_site_recovery_replication_policy` | `asrrp` | 1 | 128 | parent | `asrrp-example` |
| `azurerm_snapshot` | `snp` | 1 | 80 | resourceGroup | `snp-example` |
| `azurerm_snapshots` | `snap` | 1 | 80 | parent | `snap-example` |
| `azurerm_spatial_anchors_account` | `spaa` | 1 | 90 | resourceGroup | `spaa-example` |
| `azurerm_spring_cloud_app` | `spca` | 4 | 32 | parent | `spca-example` |
| `azurerm_spring_cloud_certificate` | `spcert` | 1 | 64 | parent | `spcert-example` |
| `azurerm_spring_cloud_service` | `spcs` | 4 | 32 | resourceGroup | `spcs-example` |
| `azurerm_sql_database` | `sqld` | 1 | 128 | parent | `sqld-example` |
| `azurerm_sql_elasticpool` | `sqlep` | 1 | 128 | parent | `sqlep-example` |
| `azurerm_sql_failover_group` | `sqlfg` | 1 | 63 | global | `sqlfg-example` |
| `azurerm_sql_firewall_rule` | `sqlfw` | 1 | 128 | parent | `sqlfw-example` |
| `azurerm_sql_server` | `sql` | 1 | 63 | global | `sql-example` |
| `azurerm_static_site` | `stapp` | 1 | 40 | resourceGroup | `stapp-example` |
| `azurerm_storage_account` | `st` | 3 | 24 | global | `stexample` |
| `azurerm_storage_blob` | `blob` | 1 | 1024 | parent | `blob-example` |
| `azurerm_storage_container` | `stct` | 3 | 63 | parent | `stct-example` |
| `azurerm_storage_data_lake_gen2_filesystem` | `stdl` | 3 | 63 | parent | `stdl-example` |
| `azurerm_storage_data_lake_gen2_path` | `stdlp` | 1 | 1024 | parent | `stdlp-example` |
| `azurerm_storage_encryption_scope` | `stes` | 4 | 63 | parent | `stesexample` |
| `azurerm_storage_queue` | `stq` | 3 | 63 | parent | `stq-example` |
| `azurerm_storage_share` | `sts` | 3 | 63 | parent | `sts-example` |
| `azurerm_storage_share_directory` | `sts` | 3 | 63 | parent | `sts-example` |
| `azurerm_storage_sync` | `stsy` | 1 | 260 | resourceGroup | `stsy-example` |
| `azurerm_storage_sync_group` | `stsg` | 1 | 260 | resourceGroup | `stsg-example` |
| `azurerm_storage_table` | `stt` | 3 | 63 | parent | `stt-example` |
| `azurerm_stream_analytics_function_javascript_udf` | `asafunc` | 3 | 63 | parent | `asafunc-example` |
| `azurerm_stream_analytics_job` | `asa` | 3 | 63 | resourceGroup | `asa-example` |
| `azurerm_stream_analytics_output_blob` | `asaoblob` | 3 | 63 | parent | `asaoblob-example` |
| `azurerm_stream_analytics_output_eventhub` | `asaoeh` | 3 | 63 | parent | `asaoeh-example` |
| `azurerm_stream_analytics_output_mssql` | `asaomssql` | 3 | 63 | parent | `asaomssql-example` |
| `azurerm_stream_analytics_output_servicebus_queue` | `asaosbq` | 3 | 63 | parent | `asaosbq-example` |
| `azurerm_stream_analytics_output_servicebus_topic` | `asaosbt` | 3 | 63 | parent | `asaosbt-example` |
| `azurerm_stream_analytics_reference_input_blob` | `asarblob` | 3 | 63 | parent | `asarblob-example` |
| `azurerm_stream_analytics_stream_input_blob` | `asaiblob` | 3 | 63 | parent | `asaiblob-example` |
| `azurerm_stream_analytics_stream_input_eventhub` | `asaieh` | 3 | 63 | parent | `asaieh-example` |
| `azurerm_stream_analytics_stream_input_iothub` | `asaiiot` | 3 | 63 | parent | `asaiiot-example` |
| `azurerm_subnet` | `snet` | 1 | 80 | parent | `snet-example` |
| `azurerm_subscription_policy_assignment` | `aspa` | 1 | 128 | subscription | `aspa-example` |
| `azurerm_subscription_template_deployment` | `subtd` | 1 | 64 | parent | `subtd-example` |
| `azurerm_synapse_firewall_rule` | `synfw` | 1 | 128 | parent | `synfw-example` |
| `azurerm_synapse_integration_runtime_azure` | `synira` | 3 | 63 | subscription | `synira-example` |
| `azurerm_synapse_integration_runtime_self_hosted` | `synirsh` | 3 | 63 | subscription | `synirsh-example` |
| `azurerm_synapse_linked_service` | `synls` | 1 | 140 | subscription | `synlsexample` |
| `azurerm_synapse_managed_private_endpoint` | `synmpe` | 3 | 63 | subscription | `synmpe-example` |
| `azurerm_synapse_private_link_hub` | `synplh` | 1 | 45 | subscription | `synplhexample` |
| `azurerm_synapse_spark_pool` | `synsp` | 1 | 15 | parent | `synspexample` |
| `azurerm_synapse_sql_pool` | `syndp` | 3 | 15 | subscription | `syndpexample` |
| `azurerm_synapse_sql_pool_vulnerability_assessment_baseline` | `syndpvab` | 3 | 63 | subscription | `syndpvabexample` |
| `azurerm_synapse_sql_pool_workload_classifier` | `syndpwc` | 1 | 128 | subscription | `syndpwc-example` |
| `azurerm_synapse_sql_pool_workload_group` | `syndpwg` | 1 | 128 | subscription | `syndpwg-example` |
| `azurerm_synapse_workspace` | `synw` | 1 | 50 | global | `synw-example` |
| `azurerm_template_deployment` | `deploy` | 1 | 64 | resourceGroup | `deploy-example` |
| `azurerm_traffic_manager_profile` | `traf` | 1 | 63 | global | `traf-example` |
| `azurerm_user_assigned_identity` | `id` | 3 | 128 | parent | `id-example` |
| `azurerm_virtual_desktop_application_group` | `vdag` | 1 | 260 | resourceGroup | `vdag-example` |
| `azurerm_virtual_desktop_host_pool` | `vdpool` | 1 | 260 | resourceGroup | `vdpool-example` |
| `azurerm_virtual_desktop_workspace` | `vdws` | 1 | 260 | resourceGroup | `vdws-example` |
| `azurerm_virtual_hub` | `vhub` | 1 | 50 | parent | `vhub-example` |
| `azurerm_virtual_hub_bgp_connection` | `vhbgp` | 1 | 80 | parent | `vhbgp-example` |
| `azurerm_virtual_hub_connection` | `vhcon` | 1 | 80 | parent | `vhcon-example` |
| `azurerm_virtual_hub_ip` | `vhip` | 1 | 80 | parent | `vhip-example` |
| `azurerm_virtual_hub_route_table` | `vhrt` | 1 | 80 | parent | `vhrt-example` |
| `azurerm_virtual_hub_security_partner_provider` | `vhspp` | 1 | 80 | resourceGroup | `vhspp-example` |
| `azurerm_virtual_machine` | `vm` | 1 | 15 | resourceGroup | `vm-example` |
| `azurerm_virtual_machine_extension` | `vmx` | 1 | 80 | parent | `vmx-example` |
| `azurerm_virtual_machine_portal_name` | `vm` | 1 | 64 | resourceGroup | `vm-example` |
| `azurerm_virtual_machine_scale_set` | `vmss` | 1 | 15 | resourceGroup | `vmss-example` |
| `azurerm_virtual_machine_scale_set_extension` | `vmssx` | 1 | 80 | parent | `vmssx-example` |
| `azurerm_virtual_network` | `vnet` | 2 | 64 | resourceGroup | `vnet-example` |
| `azurerm_virtual_network_gateway` | `vgw` | 1 | 80 | resourceGroup | `vgw-example` |
| `azurerm_virtual_network_gateway_connection` | `vngc` | 1 | 80 | resourceGroup | `vngc-example` |
| `azurerm_virtual_network_peering` | `peer` | 1 | 80 | parent | `peer-example` |
| `azurerm_virtual_wan` | `vwan` | 1 | 80 | parent | `vwan-example` |
| `azurerm_vm_windows_computer_name_prefix` | `cn` | 1 | 9 | resourceGroup | `cn-exampl` |
| `azurerm_vmware_cluster` | `vwc` | 1 | 80 | resourceGroup | `vwc-example` |
| `azurerm_vmware_express_route_authorization` | `vwera` | 1 | 80 | resourceGroup | `vwera-example` |
| `azurerm_vmware_private_cloud` | `vwpc` | 1 | 80 | resourceGroup | `vwpc-example` |
| `azurerm_vpn_gateway` | `vpng` | 1 | 80 | resourceGroup | `vpng-example` |
| `azurerm_vpn_gateway_connection` | `vcn` | 1 | 80 | parent | `vcn-example` |
| `azurerm_vpn_server_configuration` | `vpnsc` | 1 | 80 | resourceGroup | `vpnsc-example` |
| `azurerm_vpn_site` | `vst` | 1 | 80 | parent | `vst-example` |
| `azurerm_web_application_firewall_policy` | `waf` | 1 | 80 | resourceGroup | `waf-example` |
| `azurerm_web_pubsub` | `ps` | 3 | 63 | resourceGroup | `ps-example` |
| `azurerm_web_pubsub_hub` | `pshub` | 1 | 128 | parent | `pshubexample` |
| `azurerm_windows_function_app` | `fa` | 2 | 60 | global | `fa-example` |
| `azurerm_windows_function_app_slot` | `fas` | 2 | 59 | global | `fas-example` |
| `azurerm_windows_virtual_machine` | `vm` | 1 | 15 | resourceGroup | `vm-example` |
| `azurerm_windows_virtual_machine_scale_set` | `vmss` | 1 | 15 | resourceGroup | `vmss-example` |
| `azurerm_windows_web_app` | `wwapp` | 2 | 60 | global | `wwapp-example` |
| `databricks_cluster` | `dbc` | 3 | 30 | parent | `dbc-example` |
| `databricks_high_concurrency_cluster` | `dbhcc` | 3 | 30 | parent | `dbhcc-example` |
| `databricks_standard_cluster` | `dbsc` | 3 | 30 | parent | `dbsc-example` |
| `general` | `` | 1 | 250 | global | `-example` |
| `general_safe` | `` | 1 | 250 | global | `example` |

</details>

*Resource types are defined according to [Azure Cloud Adoption Framework naming and tagging best practices](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging).*

## Configuration Examples

### Environment-Based Naming

```hcl
locals {
  environment_config = {
    dev = {
      prefix = "dev"
      random_length = 3
    }
    prod = {
      prefix = "prod" 
      random_length = 5
    }
  }
  
  current_env = local.environment_config[var.environment]
}

data "azurecaf_name" "app_service" {
  name          = var.application_name
  resource_type = "azurerm_app_service"
  prefixes      = [local.current_env.prefix]
  random_length = local.current_env.random_length
  random_seed   = 12345
}
```

### Multiple Resource Generation

```hcl
data "azurecaf_name" "resources" {
  for_each = toset([
    "azurerm_resource_group",
    "azurerm_storage_account", 
    "azurerm_key_vault"
  ])
  
  name          = var.project_name
  resource_type = each.key
  prefixes      = [var.environment]
  random_length = 3
  random_seed   = 12345
}

output "resource_names" {
  value = { for k, v in data.azurecaf_name.resources : k => v.result }
}
```

## Best Practices

1. **Use Data Sources**: Prefer `data "azurecaf_name"` over `resource "azurecaf_name"` for better plan visibility
2. **Consistent Naming**: Use the same prefixes and patterns across your infrastructure
3. **Environment Separation**: Include environment identifiers in prefixes
4. **Random Length**: Use appropriate random length for uniqueness without excessive length
5. **Input Cleaning**: Keep `clean_input = true` (default) for compliance

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](https://github.com/aztfmod/terraform-provider-azurecaf/blob/main/CONTRIBUTING.md) for details.

## Support

- **Documentation**: [Terraform Registry](https://registry.terraform.io/providers/aztfmod/azurecaf/latest/docs)
- **Issues**: [GitHub Issues](https://github.com/aztfmod/terraform-provider-azurecaf/issues)
- **Discussions**: [GitHub Discussions](https://github.com/aztfmod/terraform-provider-azurecaf/discussions)

## Related Projects

| Project | Description |
|---------|-------------|
| [CAF Landing Zones](https://github.com/azure/caf-terraform-landingzones) | Azure landing zones implementation |
| [CAF Modules](https://registry.terraform.io/modules/aztfmod) | Official CAF modules |
| [Rover](https://github.com/aztfmod/rover) | DevOps toolset for landing zones |
