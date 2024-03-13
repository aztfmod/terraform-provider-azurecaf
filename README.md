# Azure Cloud Adoption Framework - Terraform provider

This provider implements a set of methodologies for naming convention implementation including the default Microsoft Cloud Adoption Framework for Azure recommendations as per https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging.

## Using the Provider

You can simply consume the provider from the Terraform registry from the following URL: [https://registry.terraform.io/providers/aztfmod/azurecaf/latest](https://registry.terraform.io/providers/aztfmod/azurecaf/latest), then add it in your provider declaration as follow:

```hcl
terraform {
  required_providers {
    azurecaf = {
      source = "aztfmod/azurecaf"
      version = "1.2.10"
    }
  }
}
```

The azurecaf_name resource allows you to:

* Clean inputs to make sure they remain compliant with the allowed patterns for each Azure resource.
* Generate random characters to append at the end of the resource name.
* Handle prefix, suffixes (either manual or as per the Azure cloud adoption framework resource conventions).
* Allow passthrough mode (simply validate the output).

## Example usage

This example outputs one name, the result of the naming convention query. The result attribute returns the name based on the convention and parameters input.

The example generates a 23 characters name compatible with the specification for an Azure Resource Group
dev-aztfmod-001

```hcl
data "azurecaf_name" "rg_example" {
  name          = "demogroup"
  resource_type = "azurerm_resource_group"
  prefixes      = ["a", "b"]
  suffixes      = ["y", "z"]
  random_length = 5
  clean_input   = true
}

output "rg_example" {
  value = data.azurecaf_name.rg_example.result
}
```
```bash
data.azurecaf_name.rg_example: Reading...
data.azurecaf_name.rg_example: Read complete after 0s [id=a-b-rg-demogroup-sjdeh-y-z]

Changes to Outputs:
  + rg_example = "a-b-rg-demogroup-sjdeh-y-z"
```

The provider generates a name using the input parameters and automatically appends a prefix (if defined), a caf prefix (resource type) and postfix (if defined) in addition to a generated padding string based on the selected naming convention.

The example above would generate a name using the pattern [prefix]-[cafprefix]-[name]-[postfix]-[5_random_chars]:

## Argument Reference

The following arguments are supported:

* **name** - (optional) the basename of the resource to create, the basename will be sanitized as per supported characters set for each Azure resources.
* **prefixes** (optional) - a list of prefix to append as the first characters of the generated name - prefixes will be separated by the separator character
* **suffixes** (optional) -  a list of additional suffix added after the basename, this is can be used to append resource index (eg. vm-001). Suffixes are separated by the separator character
* **random_length** (optional) - default to ``0`` : configure additional characters to append to the generated resource name. Random characters will remain compliant with the set of allowed characters per resources and will be appended before suffix(ess).
* **random_seed** (optional) - default to ``0`` : Define the seed to be used for random generator. 0 will not be respected and will generate a seed based in the unix time of the generation.
* **resource_type** (optional) -  describes the type of azure resource you are requesting a name from (eg. azure container registry: azurerm_container_registry). See the Resource Type section
* **resource_types** (optional) -  a list of additional resource type should you want to use the same settings for a set of resources
* **separator** (optional) - defaults to ``-``. The separator character to use between prefixes, resource type, name, suffixes, random character
* **clean_input** (optional) - defaults to ``true``. remove any noncompliant character from the name, suffix or prefix.
* **passthrough** (optional) - defaults to ``false``. Enables the passthrough mode - in that case only the clean input option is considered and the prefixes, suffixes, random, and are ignored. The resource prefixe is not added either to the resulting string
* **use_slug** (optional) - defaults to ``true``. If a slug should be added to the name - If you put false no slug (the few letters that identify the resource type) will be added to the name.

## Attributes Reference

The following attributes are exported:

* **id** - The id of the naming convention object
* **result** - The generated named for an Azure Resource based on the input parameter and the selected naming convention
* **results** - The generated name for the Azure resources based in the resource_types list

## Resource types

We define resource types as per [naming-and-tagging](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging)
The comprehensive list of resource type can be found [here](./docs/resources/azurecaf_name.md)


## Building the provider

Clone repository to: $GOPATH/src/github.com/aztfmod/terraform-provider-azurecaf

```
$ mkdir -p $GOPATH/src/github.com/aztfmod; cd $GOPATH/src/github.com/aztfmod
$ git clone https://github.com/aztfmod/terraform-provider-azurecaf.git

```
Enter the provider directory and build the provider

```
$ cd $GOPATH/src/github.com/aztfmod/terraform-provider-azurecaf
$ make build

```

## Developing the provider

If you wish to work on the provider, you'll first need Go installed on your machine (version 1.13+ is required). You'll also need to correctly setup a GOPATH, as well as adding $GOPATH/bin to your $PATH.

To display the makefile help run `make` or `make help`.

To compile the provider, run make build. This will build the provider and put the provider binary in the $GOPATH/bin directory.

```
$ make build
...
$ $GOPATH/bin/terraform-provider-azurecaf
...

```
## Testing

Running the acceptance test suite requires does not require an Azure subscription. 

to run the unit test:
```
make unittest
```

to run the integration test

```
make test
```

## Related repositories

| Repo                                                                                             | Description                                                |
|--------------------------------------------------------------------------------------------------|------------------------------------------------------------|
| [caf-terraform-landingzones](https://github.com/azure/caf-terraform-landingzones)                | landing zones repo with sample and core documentations     |
| [rover](https://github.com/aztfmod/rover)                                                        | devops toolset for operating landing zones                 |
| [azure_caf_provider](https://github.com/aztfmod/terraform-provider-azurecaf)                     | custom provider for naming conventions                     |
| [module](https://registry.terraform.io/modules/aztfmod)                                          | official CAF module available in the Terraform registry    |


## Community

Feel free to open an issue for feature or bug, or to submit a PR.

In case you have any question, you can reach out to tf-landingzones at microsoft dot com.

You can also reach us on [Gitter](https://gitter.im/aztfmod/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Contributing

information about contributing can be found at [CONTRIBUTING.md](.github/CONTRIBUTING.md)

## Resource Status

This is the current compreheensive status of the implemented resources in the provider comparing with the current list of resources in the azurerm terraform provider.

|resource | status |
|---|---|
|azurerm_aadb2c_directory | ✔ |
|azurerm_advanced_threat_protection | ❌ |
|azurerm_advisor_recommendations | ❌ |
|azurerm_analysis_services_server | ✔ |
|azurerm_api_management | ✔ |
|azurerm_api_management_api | ✔ |
|azurerm_api_management_api_diagnostic | ❌ |
|azurerm_api_management_api_operation | ❌ |
|azurerm_api_management_api_operation_policy | ❌ |
|azurerm_api_management_api_operation_tag | ✔ |
|azurerm_api_management_api_policy | ❌ |
|azurerm_api_management_api_schema | ❌ |
|azurerm_api_management_api_version_set | ❌ |
|azurerm_api_management_authorization_server | ❌ |
|azurerm_api_management_backend | ✔ |
|azurerm_api_management_certificate | ✔ |
|azurerm_api_management_custom_domain | ✔ |
|azurerm_api_management_diagnostic | ❌ |
|azurerm_api_management_gateway | ✔ |
|azurerm_api_management_group | ✔ |
|azurerm_api_management_group_user | ✔ |
|azurerm_api_management_identity_provider_aad | ❌ |
|azurerm_api_management_identity_provider_facebook | ❌ |
|azurerm_api_management_identity_provider_google | ❌ |
|azurerm_api_management_identity_provider_microsoft | ❌ |
|azurerm_api_management_identity_provider_twitter | ❌ |
|azurerm_api_management_logger | ✔ |
|azurerm_api_management_named_value | ❌ |
|azurerm_api_management_openid_connect_provider | ❌ |
|azurerm_api_management_product | ❌ |
|azurerm_api_management_product_api | ❌ |
|azurerm_api_management_product_group | ❌ |
|azurerm_api_management_product_policy | ❌ |
|azurerm_api_management_property | ❌ |
|azurerm_api_management_subscription | ❌ |
|azurerm_api_management_user | ✔ |
|azurerm_app_configuration | ✔ |
|azurerm_app_service | ✔ |
|azurerm_app_service_active_slot | ❌ |
|azurerm_app_service_certificate | ❌ |
|azurerm_app_service_certificate_order | ❌ |
|azurerm_app_service_custom_hostname_binding | ❌ |
|azurerm_app_service_environment | ✔ |
|azurerm_app_service_hybrid_connection | ❌ |
|azurerm_app_service_plan | ✔ |
|azurerm_app_service_slot | ❌ |
|azurerm_app_service_slot_virtual_network_swift_connection | ❌ |
|azurerm_app_service_source_control_token | ❌ |
|azurerm_app_service_virtual_network_swift_connection | ❌ |
|azurerm_application_gateway | ✔ |
|azurerm_application_insights | ✔ |
|azurerm_application_insights_analytics_item | ❌ |
|azurerm_application_insights_api_key | ❌ |
|azurerm_application_insights_web_test | ✔ |
|azurerm_application_security_group | ✔ |
|azurerm_attestation | ❌ |
|azurerm_automation_account | ✔ |
|azurerm_automation_certificate | ✔ |
|azurerm_automation_connection | ❌ |
|azurerm_automation_connection_certificate | ❌ |
|azurerm_automation_connection_classic_certificate | ❌ |
|azurerm_automation_connection_service_principal | ❌ |
|azurerm_automation_credential | ✔ |
|azurerm_automation_dsc_configuration | ❌ |
|azurerm_automation_dsc_nodeconfiguration | ❌ |
|azurerm_automation_hybrid_runbook_worker_group | ✔ |
|azurerm_automation_job_schedule | ✔ |
|azurerm_automation_module | ❌ |
|azurerm_automation_runbook | ✔ |
|azurerm_automation_schedule | ✔ |
|azurerm_automation_variable_bool | ❌ |
|azurerm_automation_variable_datetime | ❌ |
|azurerm_automation_variable_int | ❌ |
|azurerm_automation_variable_string | ❌ |
|azurerm_availability_set | ✔ |
|azurerm_backup_container_storage_account | ❌ |
|azurerm_backup_policy_file_share | ❌ |
|azurerm_backup_policy_vm | ❌ |
|azurerm_backup_protected_file_share | ❌ |
|azurerm_backup_protected_vm | ❌ |
|azurerm_bastion_host | ✔ |
|azurerm_batch_account | ✔ |
|azurerm_batch_application | ✔ |
|azurerm_batch_certificate | ✔ |
|azurerm_batch_pool | ✔ |
|azurerm_blueprint_assignment | ❌ |
|azurerm_blueprint_definition | ❌ |
|azurerm_blueprint_published_version | ❌ |
|azurerm_bot_channel_directline | ✔ |
|azurerm_bot_channel_email | ❌ |
|azurerm_bot_channel_ms_teams | ✔ |
|azurerm_bot_channel_slack | ✔ |
|azurerm_bot_channels_registration | ✔ |
|azurerm_bot_connection | ✔ |
|azurerm_bot_web_app | ✔ |
|azurerm_cdn_endpoint | ✔ |
|azurerm_cdn_frontdoor_custom_domain | ✔ |
|azurerm_cdn_frontdoor_endpoint | ✔ |
|azurerm_cdn_frontdoor_firewall_policy | ✔ |
|azurerm_cdn_frontdoor_origin | ✔ |
|azurerm_cdn_frontdoor_origin_group | ✔ |
|azurerm_cdn_frontdoor_profile | ✔ |
|azurerm_cdn_frontdoor_route | ✔ |
|azurerm_cdn_frontdoor_rule | ✔ |
|azurerm_cdn_frontdoor_rule_set | ✔ |
|azurerm_cdn_frontdoor_secret | ✔ |
|azurerm_cdn_frontdoor_security_policy | ✔ |
|azurerm_cdn_profile | ✔ |
|azurerm_client_config | ❌ |
|azurerm_cognitive_account | ✔ |
|azurerm_communication_service | ✔ |
|azurerm_consumption_budget_resource_group | ✔ |
|azurerm_consumption_budget_subscription | ✔ |
|azurerm_container_app | ✔ |
|azurerm_container_app_environment | ✔ |
|azurerm_container_group | ❌ |
|azurerm_container_registry | ✔ |
|azurerm_container_registry_webhook | ✔ |
|azurerm_cosmosdb_account | ✔ |
|azurerm_cosmosdb_cassandra_keyspace | ❌ |
|azurerm_cosmosdb_gremlin_database | ❌ |
|azurerm_cosmosdb_gremlin_graph | ❌ |
|azurerm_cosmosdb_mongo_collection | ❌ |
|azurerm_cosmosdb_mongo_database | ❌ |
|azurerm_cosmosdb_sql_container | ❌ |
|azurerm_cosmosdb_sql_database | ❌ |
|azurerm_cosmosdb_sql_stored_procedure | ❌ |
|azurerm_cosmosdb_table | ❌ |
|azurerm_cost_management_export_resource_group | ❌ |
|azurerm_custom_provider | ✔ |
|azurerm_dashboard | ✔ |
|azurerm_portal_dashboard | ✔ |
|azurerm_data_factory | ✔ |
|azurerm_data_factory_dataset_azure_blob | ✔ |
|azurerm_data_factory_dataset_cosmosdb_sqlapi | ✔ |
|azurerm_data_factory_dataset_delimited_text | ✔ |
|azurerm_data_factory_dataset_http | ✔ |
|azurerm_data_factory_dataset_json | ✔ |
|azurerm_data_factory_dataset_mysql | ✔ |
|azurerm_data_factory_dataset_postgresql | ✔ |
|azurerm_data_factory_dataset_sql_server_table | ✔ |
|azurerm_data_factory_integration_runtime_managed | ✔ |
|azurerm_data_factory_integration_runtime_self_hosted | ❌ |
|azurerm_data_factory_linked_service_azure_blob_storage | ✔ |
|azurerm_data_factory_linked_service_azure_databricks | ✔ |
|azurerm_data_factory_linked_service_azure_file_storage | ❌ |
|azurerm_data_factory_linked_service_azure_function | ✔ |
|azurerm_data_factory_linked_service_azure_sql_database | ✔ |
|azurerm_data_factory_linked_service_cosmosdb | ✔ |
|azurerm_data_factory_linked_service_data_lake_storage_gen2 | ✔ |
|azurerm_data_factory_linked_service_key_vault | ✔ |
|azurerm_data_factory_linked_service_mysql | ✔ |
|azurerm_data_factory_linked_service_postgresql | ✔ |
|azurerm_data_factory_linked_service_sftp | ✔ |
|azurerm_data_factory_linked_service_sql_server | ✔ |
|azurerm_data_factory_linked_service_web | ✔ |
|azurerm_data_factory_pipeline | ✔ |
|azurerm_data_factory_trigger_schedule | ✔ |
|azurerm_data_lake_analytics_account | ✔ |
|azurerm_data_lake_analytics_firewall_rule | ✔ |
|azurerm_data_lake_store | ✔ |
|azurerm_data_lake_store_file | ❌ |
|azurerm_data_lake_store_firewall_rule | ✔ |
|azurerm_data_protection_backup_policy_blob_storage | ✔ |
|azurerm_data_protection_backup_policy_disk | ✔ |
|azurerm_data_protection_backup_policy_postgresql | ✔ |
|azurerm_data_protection_backup_vault | ✔ |
|azurerm_data_share | ❌ |
|azurerm_data_share_account | ❌ |
|azurerm_data_share_dataset_blob_storage | ❌ |
|azurerm_data_share_dataset_data_lake_gen1 | ❌ |
|azurerm_data_share_dataset_data_lake_gen2 | ❌ |
|azurerm_data_share_dataset_kusto_cluster | ❌ |
|azurerm_data_share_dataset_kusto_database | ❌ |
|azurerm_database_migration_project | ✔ |
|azurerm_database_migration_service | ✔ |
|azurerm_databricks_workspace | ✔ |
|azurerm_dedicated_hardware_security_module | ❌ |
|azurerm_dedicated_host | ✔ |
|azurerm_dedicated_host_group | ✔ |
|azurerm_dev_test_global_vm_shutdown_schedule | ❌ |
|azurerm_dev_test_lab | ✔ |
|azurerm_dev_test_linux_virtual_machine | ✔ |
|azurerm_dev_test_policy | ❌ |
|azurerm_dev_test_schedule | ❌ |
|azurerm_dev_test_virtual_network | ❌ |
|azurerm_dev_test_windows_virtual_machine | ✔ |
|azurerm_devspace_controller | ❌ |
|azurerm_digital_twins_endpoint_eventgrid | ✔ |
|azurerm_digital_twins_endpoint_eventhub | ✔ |
|azurerm_digital_twins_endpoint_servicebus | ✔ |
|azurerm_digital_twins_instance | ✔ |
|azurerm_disk_encryption_set | ✔ |
|azurerm_dns_a_record | ❌ |
|azurerm_dns_aaaa_record | ❌ |
|azurerm_dns_caa_record | ❌ |
|azurerm_dns_cname_record | ❌ |
|azurerm_dns_mx_record | ❌ |
|azurerm_dns_ns_record | ❌ |
|azurerm_dns_ptr_record | ❌ |
|azurerm_dns_srv_record | ❌ |
|azurerm_dns_txt_record | ❌ |
|azurerm_dns_zone | ✔ |
|azurerm_eventgrid_domain | ✔ |
|azurerm_eventgrid_domain_topic | ✔ |
|azurerm_eventgrid_event_subscription | ✔ |
|azurerm_eventgrid_system_topic | ❌ |
|azurerm_eventgrid_topic | ✔ |
|azurerm_eventhub | ✔ |
|azurerm_eventhub_authorization_rule | ✔ |
|azurerm_eventhub_cluster | ❌ |
|azurerm_eventhub_consumer_group | ✔ |
|azurerm_eventhub_namespace | ✔ |
|azurerm_eventhub_namespace_authorization_rule | ✔ |
|azurerm_eventhub_namespace_disaster_recovery_config | ✔ |
|azurerm_express_route_circuit | ✔ |
|azurerm_express_route_circuit_authorization | ❌ |
|azurerm_express_route_circuit_peering | ❌ |
|azurerm_express_route_gateway | ✔ |
|azurerm_federated_identity_credential | ✔ |
|azurerm_firewall | ✔ |
|azurerm_firewall_application_rule_collection | ❌ |
|azurerm_firewall_nat_rule_collection | ❌ |
|azurerm_firewall_network_rule_collection | ❌ |
|azurerm_firewall_policy | ✔ |
|azurerm_firewall_policy_rule_collection_group | ❌ |
|azurerm_frontdoor | ✔ |
|azurerm_frontdoor_custom_https_configuration | ❌ |
|azurerm_frontdoor_firewall_policy | ✔ |
|azurerm_function_app | ✔ |
|azurerm_function_app_host_keys | ❌ |
|azurerm_function_app_slot | ✔ |
|azurerm_hdinsight_cluster | ❌ |
|azurerm_hdinsight_hadoop_cluster | ✔ |
|azurerm_hdinsight_hbase_cluster | ✔ |
|azurerm_hdinsight_interactive_query_cluster | ✔ |
|azurerm_hdinsight_kafka_cluster | ✔ |
|azurerm_hdinsight_ml_services_cluster | ✔ |
|azurerm_hdinsight_rserver_cluster | ✔ |
|azurerm_hdinsight_spark_cluster | ✔ |
|azurerm_hdinsight_storm_cluster | ✔ |
|azurerm_healthcare_dicom_service | ✔ |
|azurerm_healthcare_fhir_service | ✔ |
|azurerm_healthcare_medtech_service | ✔ |
|azurerm_healthcare_service | ✔ |
|azurerm_healthcare_workspace | ✔ |
|azurerm_hpc_cache | ❌ |
|azurerm_hpc_cache_blob_target | ❌ |
|azurerm_hpc_cache_nfs_target | ❌ |
|azurerm_image | ✔ |
|azurerm_images | ❌ |
|azurerm_integration_service_environment | ✔ |
|azurerm_iot_security_device_group | ✔ |
|azurerm_iot_security_solution | ✔ |
|azurerm_iot_time_series_insights_access_policy | ❌ |
|azurerm_iot_time_series_insights_reference_data_set | ❌ |
|azurerm_iot_time_series_insights_standard_environment | ❌ |
|azurerm_iotcentral_application | ✔ |
|azurerm_iothub | ✔ |
|azurerm_iothub_certificate | ✔ |
|azurerm_iothub_consumer_group | ✔ |
|azurerm_iothub_dps | ✔ |
|azurerm_iothub_dps_certificate | ✔ |
|azurerm_iothub_dps_shared_access_policy | ✔ |
|azurerm_iothub_endpoint_eventhub | ❌ |
|azurerm_iothub_endpoint_servicebus_queue | ❌ |
|azurerm_iothub_endpoint_servicebus_topic | ❌ |
|azurerm_iothub_endpoint_storage_container | ❌ |
|azurerm_iothub_fallback_route | ❌ |
|azurerm_iothub_route | ❌ |
|azurerm_iothub_shared_access_policy | ✔ |
|azurerm_ip_group | ✔ |
|azurerm_key_vault | ✔ |
|azurerm_key_vault_access_policy | ❌ |
|azurerm_key_vault_certificate | ✔ |
|azurerm_key_vault_certificate_issuer | ❌ |
|azurerm_key_vault_key | ✔ |
|azurerm_key_vault_secret | ✔ |
|azurerm_kubernetes_cluster | ✔ |
|azurerm_kubernetes_cluster_node_pool | ❌ |
|azurerm_kubernetes_fleet_manager | ✔ |
|azurerm_kubernetes_service_versions | ❌ |
|azurerm_kusto_attached_database_configuration | ❌ |
|azurerm_kusto_cluster | ✔ |
|azurerm_kusto_cluster_customer_managed_key | ❌ |
|azurerm_kusto_cluster_principal_assignment | ❌ |
|azurerm_kusto_database | ✔ |
|azurerm_kusto_database_principal | ❌ |
|azurerm_kusto_database_principal_assignment | ❌ |
|azurerm_kusto_eventhub_data_connection | ✔ |
|azurerm_lb | ✔ |
|azurerm_lb_backend_address_pool | ✔ |
|azurerm_lb_backend_address_pool_address | ✔ |
|azurerm_lb_nat_pool | ✔ |
|azurerm_lb_nat_rule | ✔ |
|azurerm_lb_outbound_rule | ✔ |
|azurerm_lb_probe | ✔ |
|azurerm_lb_rule | ✔ |
|azurerm_lighthouse_assignment | ❌ |
|azurerm_lighthouse_definition | ❌ |
|azurerm_linux_virtual_machine | ✔ |
|azurerm_linux_virtual_machine_scale_set | ✔ |
|azurerm_linux_web_app | ✔ |
|azurerm_linux_web_app_slot | ⚠ |
|azurerm_load_test | ✔ |
|azurerm_local_network_gateway | ✔ |
|azurerm_log_analytics_cluster | ✔ |
|azurerm_log_analytics_data_export_rule | ❌ |
|azurerm_log_analytics_datasource_windows_event | ❌ |
|azurerm_log_analytics_datasource_windows_performance_counter | ❌ |
|azurerm_log_analytics_linked_service | ❌ |
|azurerm_log_analytics_linked_storage_account | ❌ |
|azurerm_log_analytics_saved_search | ❌ |
|azurerm_log_analytics_solution | ❌ |
|azurerm_log_analytics_storage_insights | ✔ |
|azurerm_log_analytics_workspace | ✔ |
|azurerm_logic_app_action_custom | ✔ |
|azurerm_logic_app_action_http | ✔ |
|azurerm_logic_app_integration_account | ✔ |
|azurerm_logic_app_trigger_custom | ✔ |
|azurerm_logic_app_trigger_http_request | ✔ |
|azurerm_logic_app_trigger_recurrence | ✔ |
|azurerm_logic_app_workflow | ✔ |
|azurerm_machine_learning_compute_instance  | ✔ |
|azurerm_machine_learning_workspace | ✔ |
|azurerm_maintenance_assignment_dedicated_host | ❌ |
|azurerm_maintenance_assignment_virtual_machine | ❌ |
|azurerm_maintenance_configuration | ✔ |
|azurerm_managed_application | ❌ |
|azurerm_managed_application_definition | ❌ |
|azurerm_managed_disk | ✔ |
|azurerm_management_group | ❌ |
|azurerm_management_lock | ❌ |
|azurerm_maps_account | ✔ |
|azurerm_mariadb_configuration | ❌ |
|azurerm_mariadb_database | ✔ |
|azurerm_mariadb_firewall_rule | ✔ |
|azurerm_mariadb_server | ✔ |
|azurerm_mariadb_virtual_network_rule | ✔ |
|azurerm_marketplace_agreement | ❌ |
|azurerm_media_services_account | ❌ |
|azurerm_monitor_action_group | ✔ |
|azurerm_monitor_action_rule_action_group | ❌ |
|azurerm_monitor_action_rule_suppression | ❌ |
|azurerm_monitor_activity_log_alert | ❌ |
|azurerm_monitor_autoscale_setting | ✔ |
|azurerm_monitor_data_collection_endpoint | ✔ |
|azurerm_monitor_diagnostic_categories | ❌ |
|azurerm_monitor_diagnostic_setting | ✔ |
|azurerm_monitor_log_profile | ❌ |
|azurerm_monitor_metric_alert | ✔ |
|azurerm_monitor_private_link_scope | ✔ |
|azurerm_monitor_scheduled_query_rules_alert | ✔ |
|azurerm_monitor_scheduled_query_rules_log | ❌ |
|azurerm_monitor_smart_detector_alert_rule | ❌ |
|azurerm_mssql_database | ✔ |
|azurerm_mssql_database_extended_auditing_policy | ❌ |
|azurerm_mssql_database_vulnerability_assessment_rule_baseline | ❌ |
|azurerm_mssql_elasticpool | ✔ |
|azurerm_mssql_mi | ✔ |
|azurerm_mssql_server | ✔ |
|azurerm_mssql_server_extended_auditing_policy | ❌ |
|azurerm_mssql_server_security_alert_policy | ❌ |
|azurerm_mssql_server_vulnerability_assessment | ❌ |
|azurerm_mssql_virtual_machine | ❌ |
|azurerm_mysql_active_directory_administrator | ❌ |
|azurerm_mysql_configuration | ❌ |
|azurerm_mysql_database | ✔ |
|azurerm_mysql_firewall_rule | ✔ |
|azurerm_mysql_flexible_server | ✔ |
|azurerm_mysql_flexible_server_database | ✔ |
|azurerm_mysql_flexible_server_firewall_rule | ✔ |
|azurerm_mysql_server | ✔ |
|azurerm_mysql_server_key | ❌ |
|azurerm_mysql_virtual_network_rule | ✔ |
|azurerm_nat_gateway | ❌ |
|azurerm_nat_gateway_public_ip_association | ❌ |
|azurerm_netapp_account | ✔ |
|azurerm_netapp_pool | ✔ |
|azurerm_netapp_snapshot | ✔ |
|azurerm_netapp_volume | ✔ |
|azurerm_network_connection_monitor | ❌ |
|azurerm_network_ddos_protection_plan | ❌ |
|azurerm_network_interface | ✔ |
|azurerm_network_interface_application_gateway_backend_address_pool_association | ❌ |
|azurerm_network_interface_application_security_group_association | ❌ |
|azurerm_network_interface_backend_address_pool_association | ❌ |
|azurerm_network_interface_nat_rule_association | ❌ |
|azurerm_network_interface_security_group_association | ❌ |
|azurerm_network_packet_capture | ❌ |
|azurerm_network_profile | ❌ |
|azurerm_network_security_group | ✔ |
|azurerm_network_security_rule | ✔ |
|azurerm_network_service_tags | ❌ |
|azurerm_network_watcher | ✔ |
|azurerm_network_watcher_flow_log | ❌ |
|azurerm_nginx_deployment | ✔ |
|azurerm_notification_hub | ✔ |
|azurerm_notification_hub_authorization_rule | ✔ |
|azurerm_notification_hub_namespace | ✔ |
|azurerm_orchestrated_virtual_machine_scale_set | ❌ |
|azurerm_packet_capture | ❌ |
|azurerm_platform_image | ❌ |
|azurerm_point_to_site_vpn_gateway | ✔ |
|azurerm_policy_assignment | ❌ |
|azurerm_policy_definition | ❌ |
|azurerm_policy_remediation | ❌ |
|azurerm_policy_set_definition | ❌ |
|azurerm_postgresql_active_directory_administrator | ❌ |
|azurerm_postgresql_configuration | ❌ |
|azurerm_postgresql_database | ✔ |
|azurerm_postgresql_firewall_rule | ✔ |
|azurerm_postgresql_flexible_server | ✔ |
|azurerm_postgresql_flexible_server_configuration | ❌ |
|azurerm_postgresql_flexible_server_database | ✔ |
|azurerm_postgresql_flexible_server_firewall_rule | ✔ |
|azurerm_postgresql_server | ✔ |
|azurerm_postgresql_server_key | ❌ |
|azurerm_postgresql_virtual_network_rule | ✔ |
|azurerm_powerbi_embedded | ✔ |
|azurerm_private_dns_a_record | ❌ |
|azurerm_private_dns_aaaa_record | ❌ |
|azurerm_private_dns_cname_record | ❌ |
|azurerm_private_dns_mx_record | ❌ |
|azurerm_private_dns_ptr_record | ❌ |
|azurerm_private_dns_resolver | ✔ |
|azurerm_private_dns_resolver_dns_forwarding_ruleset | ✔ |
|azurerm_private_dns_resolver_forwarding_rule | ✔ |
|azurerm_private_dns_resolver_inbound_endpoint | ✔ |
|azurerm_private_dns_resolver_outbound_endpoint | ✔ |
|azurerm_private_dns_resolver_virtual_network_link | ✔ |
|azurerm_private_dns_srv_record | ❌ |
|azurerm_private_dns_txt_record | ❌ |
|azurerm_private_dns_zone | ✔ |
|azurerm_private_dns_zone_virtual_network_link | ✔ |
|azurerm_private_endpoint | ✔ |
|azurerm_private_endpoint_connection | ❌ |
|azurerm_private_link_service | ❌ |
|azurerm_private_link_service_endpoint_connections | ❌ |
|azurerm_proximity_placement_group | ❌ |
|azurerm_public_ip | ✔ |
|azurerm_public_ip_prefix | ✔ |
|azurerm_public_ips | ❌ |
|azurerm_purview_account | ✔ |
|azurerm_recovery_services_vault | ✔ |
|azurerm_redhat_openshift_cluster | ✔ |
|azurerm_redhat_openshift_domain | ✔ |
|azurerm_redis_cache | ✔ |
|azurerm_redis_firewall_rule | ✔ |
|azurerm_redis_linked_server | ❌ |
|azurerm_relay_hybrid_connection | ✔ |
|azurerm_relay_namespace | ✔ |
|azurerm_resource_group | ✔ |
|azurerm_resource_group_policy_assignment | ✔ |
|azurerm_resource_group_template_deployment | ❌ |
|azurerm_role_assignment | ✔ |
|azurerm_role_definition | ✔ |
|azurerm_route | ✔ |
|azurerm_route_filter | ❌ |
|azurerm_route_server | ✔ |
|azurerm_route_table | ✔ |
|azurerm_search_service | ✔ |
|azurerm_security_center_auto_provisioning | ❌ |
|azurerm_security_center_automation | ❌ |
|azurerm_security_center_contact | ❌ |
|azurerm_security_center_setting | ❌ |
|azurerm_security_center_subscription_pricing | ❌ |
|azurerm_security_center_workspace | ❌ |
|azurerm_sentinel_alert_rule | ❌ |
|azurerm_sentinel_alert_rule_ms_security_incident | ❌ |
|azurerm_sentinel_alert_rule_scheduled | ❌ |
|azurerm_service_fabric_cluster | ✔ |
|azurerm_service_fabric_mesh_application | ❌ |
|azurerm_service_fabric_mesh_local_network | ❌ |
|azurerm_service_fabric_mesh_secret | ❌ |
|azurerm_service_fabric_mesh_secret_value | ❌ |
|azurerm_servicebus_namespace | ✔ |
|azurerm_servicebus_namespace_authorization_rule | ✔ |
|azurerm_servicebus_namespace_network_rule_set | ❌ |
|azurerm_servicebus_queue | ✔ |
|azurerm_servicebus_queue_authorization_rule | ✔ |
|azurerm_servicebus_subscription | ✔ |
|azurerm_servicebus_subscription_rule | ✔ |
|azurerm_servicebus_topic | ✔ |
|azurerm_servicebus_topic_authorization_rule | ✔ |
|azurerm_shared_image | ✔ |
|azurerm_shared_image_gallery | ✔ |
|azurerm_shared_image_version | ❌ |
|azurerm_shared_image_versions | ❌ |
|azurerm_signalr_service | ✔ |
|azurerm_site_recovery_fabric | ❌ |
|azurerm_site_recovery_network_mapping | ❌ |
|azurerm_site_recovery_protection_container | ❌ |
|azurerm_site_recovery_protection_container_mapping | ❌ |
|azurerm_site_recovery_replicated_vm | ❌ |
|azurerm_site_recovery_replication_policy | ❌ |
|azurerm_snapshot | ❌ |
|azurerm_spatial_anchors_account | ❌ |
|azurerm_spring_cloud_app | ❌ |
|azurerm_spring_cloud_certificate | ❌ |
|azurerm_spring_cloud_service | ❌ |
|azurerm_sql_active_directory_administrator | ❌ |
|azurerm_sql_database | ❌ |
|azurerm_sql_elasticpool | ✔ |
|azurerm_sql_failover_group | ✔ |
|azurerm_sql_firewall_rule | ✔ |
|azurerm_sql_server | ✔ |
|azurerm_sql_virtual_network_rule | ❌ |
|azurerm_static_site | ✔ |
|azurerm_storage_account | ✔ |
|azurerm_storage_account_blob_container_sas | ❌ |
|azurerm_storage_account_customer_managed_key | ❌ |
|azurerm_storage_account_network_rules | ❌ |
|azurerm_storage_account_sas | ❌ |
|azurerm_storage_blob | ✔ |
|azurerm_storage_container | ✔ |
|azurerm_storage_data_lake_gen2_filesystem | ✔ |
|azurerm_storage_data_lake_gen2_path | ❌ |
|azurerm_storage_encryption_scope | ❌ |
|azurerm_storage_management_policy | ❌ |
|azurerm_storage_queue | ✔ |
|azurerm_storage_share | ✔ |
|azurerm_storage_share_directory | ✔ |
|azurerm_storage_sync | ✔ |
|azurerm_storage_sync_group | ✔ |
|azurerm_storage_table | ✔ |
|azurerm_storage_table_entity | ❌ |
|azurerm_stream_analytics_function_javascript_udf | ✔ |
|azurerm_stream_analytics_job | ✔ |
|azurerm_stream_analytics_output_blob | ✔ |
|azurerm_stream_analytics_output_eventhub | ✔ |
|azurerm_stream_analytics_output_mssql | ✔ |
|azurerm_stream_analytics_output_servicebus_queue | ✔ |
|azurerm_stream_analytics_output_servicebus_topic | ✔ |
|azurerm_stream_analytics_reference_input_blob | ✔ |
|azurerm_stream_analytics_stream_input_blob | ✔ |
|azurerm_stream_analytics_stream_input_eventhub | ✔ |
|azurerm_stream_analytics_stream_input_iothub | ✔ |
|azurerm_subnet | ✔ |
|azurerm_subnet_nat_gateway_association | ❌ |
|azurerm_subnet_network_security_group_association | ❌ |
|azurerm_subnet_route_table_association | ❌ |
|azurerm_subscription | ❌ |
|azurerm_subscription_policy_assignment | ✔ |
|azurerm_subscription_template_deployment | ❌ |
|azurerm_subscriptions | ❌ |
|azurerm_synapse_firewall_rule | ✔ |
|azurerm_synapse_integration_runtime_azure | ✔ |
|azurerm_synapse_integration_runtime_self_hosted | ✔ |
|azurerm_synapse_linked_service | ✔ |
|azurerm_synapse_managed_private_endpoint | ✔ |
|azurerm_synapse_private_link_hub | ✔ |
|azurerm_synapse_role_assignment | ❌ |
|azurerm_synapse_spark_pool | ✔ |
|azurerm_synapse_sql_pool | ✔ |
|azurerm_synapse_sql_pool_vulnerability_assessment_baseline | ✔ |
|azurerm_synapse_sql_pool_workload_classifier | ✔ |
|azurerm_synapse_sql_pool_workload_group | ✔ |
|azurerm_synapse_workspace | ✔ |
|azurerm_template_deployment | ✔ |
|azurerm_traffic_manager_endpoint | ❌ |
|azurerm_traffic_manager_geographical_location | ❌ |
|azurerm_traffic_manager_profile | ✔ |
|azurerm_user_assigned_identity | ✔ |
|azurerm_virtual_desktop_application_group | ✔ |
|azurerm_virtual_desktop_host_pool | ✔ |
|azurerm_virtual_desktop_workspace | ✔ |
|azurerm_virtual_desktop_workspace_application_group_association | ❌ |
|azurerm_virtual_hub | ✔ |
|azurerm_virtual_hub_bgp_connection | ❌ |
|azurerm_virtual_hub_connection | ✔ |
|azurerm_virtual_hub_ip | ❌ |
|azurerm_virtual_hub_route_table | ❌ |
|azurerm_virtual_hub_security_partner_provider | ❌ |
|azurerm_virtual_machine | ✔ |
|azurerm_virtual_machine_data_disk_attachment | ❌ |
|azurerm_virtual_machine_extension | ❌ |
|azurerm_virtual_machine_scale_set | ✔ |
|azurerm_virtual_machine_scale_set_extension | ❌ |
|azurerm_virtual_network | ✔ |
|azurerm_virtual_network_gateway | ✔ |
|azurerm_virtual_network_gateway_connection | ❌ |
|azurerm_virtual_network_peering | ✔ |
|azurerm_virtual_wan | ✔ |
|azurerm_vmware_cluster | ✔ |
|azurerm_vmware_express_route_authorization | ✔ |
|azurerm_vmware_private_cloud | ✔ |
|azurerm_vpn_gateway | ❌ |
|azurerm_vpn_gateway_connection | ✔ |
|azurerm_vpn_server_configuration | ❌ |
|azurerm_vpn_site | ✔ |
|azurerm_web_application_firewall_policy | ✔ |
|azurerm_web_pubsub | ✔ |
|azurerm_web_pubsub_hub | ✔ |
|azurerm_windows_virtual_machine | ✔ |
|azurerm_windows_virtual_machine_scale_set | ✔ |
|azurerm_windows_web_app | ✔ |
|azurerm_windows_web_app_slot | ⚠ |


❌ = Not yet implemented
✔  = Already implemented
⚠  = Will not be implemented
