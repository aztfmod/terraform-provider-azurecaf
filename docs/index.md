# Azurecaf provider

The Azurecaf provider is a *logical provider* which means that it works entirely within Terraform's logic, and doesn't interact with any other services. The goal of this provider is to provider helper methods in implementing Azure landing zones using Terraform.

The Azurecaf provider currently contains a two resources based on the Terraform Random_string provider. The naming_convention resources enforce is the first iteration of our naming convention implementation enforcing Azure Cloud Adoption Framework naming convention.

As per the growing number of azure resources a new implementation is now available using the azurecaf_name resource to avoid breaking changes. The new implementation supports an extensive list of resource types and will be updated on a regular basis as new services are released

## Resource types

We define resource types as per: https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging 

Current supported resource types:

| Resource type           | Resource type code (short)  | minimum length  |  maximum length | lowercase only | validation regex                          |
| ------------------------| ----------------------------|-----------------|-----------------|----------------|-------------------------------------------|
| azurerm_analysis_services_server| as| 3| 63| true| "^[a-z][a-z0-9]{2,62}$" |
| azurerm_api_management_service| apim| 1| 50| false| "^[a-z][a-zA-Z0-9-][a-zA-Z0-9]{0,48}$"|
| azurerm_app_configuration| appcg| 5| 50| false| "^[a-zA-Z0-9_-]{5,50}$"|
| azurerm_role_assignment| ra| 1| 64| false| "^[^%]{0,63}[^ %.]$"|
| azurerm_role_definition| rd| 1| 64| false| "^[^%]{0,63}[^ %.]$"|
| azurerm_automation_account| aa| 6| 50| false| "^[a-zA-Z][a-zA-Z0-9-]{4,48}[a-zA-Z0-9]$"|
| azurerm_automation_certificate| aacert| 1| 128| false| "^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$"|
| azurerm_automation_credential| aacred| 1| 128| false| "^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$"|
| azurerm_automation_runbook| aarun| 1| 63| false| "^[a-zA-Z][a-zA-Z0-9-]{0,62}$"|
| azurerm_automation_schedule| aasched| 1| 128| false| "^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$"|
| azurerm_automation_variable| aavar| 1| 128| false| "^[^<>*%:.?\\+\\/]{0,127}[^<>*%:.?\\+\\/ ]$"|
| azurerm_batch_account| ba| 3| 24| true| "^[a-z0-9]{3,24}$"|
| azurerm_batch_application| baapp| 1| 64| false| "^[a-zA-Z0-9_-]{1,64}$"|
| azurerm_batch_certificate| bacert| 5| 45| false| "^[a-zA-Z0-9_-]{5,45}$"|
| azurerm_batch_pool| bapool| 3| 24| false| "^[a-zA-Z0-9_-]{1,24}$"|
| azurerm_bot_web_app| bot| 2| 64| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,63}$"|
| azurerm_bot_channel_Email| botmail| 2| 64| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,63}$"|
| azurerm_bot_channel_ms_teams| botteams| 2| 64| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,63}$"|
| azurerm_bot_channel_slack| botslack| 2| 64| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,63}$"|
| azurerm_bot_channel_directline| botline| 2| 64| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,63}$"|
| azurerm_bot_channels_registration| botchan| 2| 64| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,63}$"|
| azurerm_bot_connection| botcon| 2| 64| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,63}$"|
| azurerm_redis_cache| redis| 1| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]$"|
| azurerm_redis_firewall_rule| redisfw| 1| 256| false| "^[a-zA-Z0-9]{1,256}$"|
| azurerm_cdn_profile| cdnprof| 1| 260| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{0,258}[a-zA-Z0-9]$"|
| azurerm_cdn_endpoint| cdn| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{0,48}[a-zA-Z0-9]$"|
| azurerm_cognitive_account| cog| 2| 64| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{0,63}$"|
| azurerm_availability_set| avail| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{0,78}[a-zA-Z0-9_]$"|
| azurerm_disk_encryption_set| des| 1| 80| false| "^[a-zA-Z0-9_]{1,80}$"|
| azurerm_image| img| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{0,78}[a-zA-Z0-9_]$"|
| azurerm_linux_virtual_machine| vm| 1| 64| false| "^[^\\/\"\\[\\]:|<>+=;,?*@&_][^\\/\"\\[\\]:|<>+=;,?*@&]{0,62}[^\\/\"\\[\\]:|<>+=;,?*@&.-]$"|
| azurerm_linux_virtual_machine_scale_set| vmss| 1| 64| false| "^[^\\/\"\\[\\]:|<>+=;,?*@&_][^\\/\"\\[\\]:|<>+=;,?*@&]{0,62}[^\\/\"\\[\\]:|<>+=;,?*@&.-]$"|
| azurerm_managed_disk| dsk| 1| 80| false| "^[a-zA-Z0-9_]{1,80}$"|
| azurerm_virtual_machine| vm| 1| 15| false| "^[^\\/\"\\[\\]:|<>+=;,?*@&_][^\\/\"\\[\\]:|<>+=;,?*@&]{0,13}[^\\/\"\\[\\]:|<>+=;,?*@&.-]$"|
| azurerm_virtual_machine_scale_set| vmss| 1| 15| false| "^[^\\/\"\\[\\]:|<>+=;,?*@&_][^\\/\"\\[\\]:|<>+=;,?*@&]{0,13}[^\\/\"\\[\\]:|<>+=;,?*@&.-]$"|
| azurerm_windows_virtual_machine| vm| 1| 15| false| "^[^\\/\"\\[\\]:|<>+=;,?*@&_][^\\/\"\\[\\]:|<>+=;,?*@&]{0,13}[^\\/\"\\[\\]:|<>+=;,?*@&.-]$"|
| azurerm_windows_virtual_machine_scale_set| vmss| 1| 15| false| "^[^\\/\"\\[\\]:|<>+=;,?*@&_][^\\/\"\\[\\]:|<>+=;,?*@&]{0,13}[^\\/\"\\[\\]:|<>+=;,?*@&.-]$"|
| azurerm_containerGroups| cg| 1| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]$"|
| azurerm_container_registry| cr| 1| 63| true| "^[a-zA-Z0-9]{1,63}$"|
| azurerm_container_registry_webhook| crwh| 1| 50| false| "^[a-zA-Z0-9]{1,50}$"|
| azurerm_kubernetes_cluster| aks| 1| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{0,61}[a-zA-Z0-9]$"|
| azurerm_cosmosdb_account| cosmos| 1| 63| false| "^[a-z0-9][a-zA-Z0-9-_.]{0,61}[a-zA-Z0-9]$"|
| azurerm_custom_provider| prov| 3| 64| false| "^[^&%?\\/]{2,63}[^&%.?\\/ ]$"|
| azurerm_mariadb_server| maria| 3| 63| false| "^[a-z0-9][a-zA-Z0-9-]{1,61}[a-z0-9]$"|
| azurerm_mariadb_firewall_rule| mariafw| 1| 128| false| "^[a-zA-Z0-9-_]{1,128}$"|
| azurerm_mariadb_database| mariadb| 1| 63| false| "^[a-zA-Z0-9-_]{1,63}$"|
| azurerm_mariadb_virtual_network_rule| mariavn| 1| 128| false| "^[a-zA-Z0-9-_]{1,128}$"|
| azurerm_mysql_server| mysql| 3| 63| false| "^[a-z0-9][a-zA-Z0-9-]{1,61}[a-z0-9]$"|
| azurerm_mysql_firewall_rule| mysqlfw| 1| 128| false| "^[a-zA-Z0-9-_]{1,128}$"|
| azurerm_mysql_database| mysqldb| 1| 63| false| "^[a-zA-Z0-9-_]{1,63}$"|
| azurerm_mysql_virtual_network_rule| mysqlvn| 1| 128| false| "^[a-zA-Z0-9-_]{1,128}$"|
| azurerm_postgresql_server| psql| 3| 63| false| "^[a-z0-9][a-zA-Z0-9-]{1,61}[a-z0-9]$"|
| azurerm_postgresql_firewall_rule| psqlfw| 1| 128| false| "^[a-zA-Z0-9-_]{1,128}$"|
| azurerm_postgresql_database| psqldb| 1| 63| false| "^[a-zA-Z0-9-_]{1,63}$"|
| azurerm_postgresql_virtual_network_rule| psqlvn| 1| 128| false| "^[a-zA-Z0-9-_]{1,128}$"|
| azurerm_database_migration_project| migr| 2| 57| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,56}$"|
| azurerm_database_migration_service| dms| 2| 62| false| "^[a-zA-Z0-9][a-zA-Z0-9-_.]{1,61}$"|
| azurerm_databricks_workspace| dbw| 3| 30| false| "^[a-zA-Z0-9-_]{3,30}$"|
| azurerm_kusto_cluster| kc| 4| 22| false| "^[a-z][a-z0-9]{3,21}$"|
| azurerm_kusto_database| kdb| 1| 260| false| "^[a-zA-Z0-9- .]{1,260}$"|
| azurerm_kusto_eventhub_data_connection| kehc| 1| 40| false| "^[a-zA-Z0-9- .]{1,40}$"|
| azurerm_data_factory| adf| 3| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]$"|
| azurerm_data_factory_dataset_mysql| adfmysql| 1| 260| false| "^[a-zA-Z0-9][^<>*%:.?\\+\\/]{0,258}[a-zA-Z0-9]$"|
| azurerm_data_factory_dataset_postgresql| adfpsql| 1| 260| false| "^[a-zA-Z0-9][^<>*%:.?\\+\\/]{0,258}[a-zA-Z0-9]$"|
| azurerm_data_factory_dataset_sql_server_table| adfmssql| 1| 260| false| "^[a-zA-Z0-9][^<>*%:.?\\+\\/]{0,258}[a-zA-Z0-9]$"|
| azurerm_data_factory_integration_runtime_managed| adfir| 3| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]$"|
| azurerm_data_factory_pipeline| adfpl| 1| 260| false| "^[a-zA-Z0-9][^<>*%:.?\\+\\/]{0,258}[a-zA-Z0-9]$"|
| azurerm_data_factory_linked_service_data_lake_storage_gen2| adfsvst| 1| 260| false| "^[a-zA-Z0-9][^<>*%:.?\\+\\/]{0,259}$"|
| azurerm_data_factory_linked_service_key_vault| adfsvkv| 1| 260| false| "^[a-zA-Z0-9][^<>*%:.?\\+\\/]{0,259}$"|
| azurerm_data_factory_linked_service_mysql| adfsvmysql| 1| 260| false| "^[a-zA-Z0-9][^<>*%:.?\\+\\/]{0,259}$"|
| azurerm_data_factory_linked_service_postgresql| adfsvpsql| 1| 260| false| "^[a-zA-Z0-9][^<>*%:.?\\+\\/]{0,259}$"|
| azurerm_data_factory_linked_service_sql_server| adfsvmssql| 1| 260| false| "^[a-zA-Z0-9][^<>*%:.?\\+\\/]{0,259}$"|
| azurerm_data_factory_trigger_schedule| adftg| 1| 260| false| "^[a-zA-Z0-9][^<>*%:.?\\+\\/]{0,259}$"|
| azurerm_data_lake_analytics_account| dla| 3| 24| false| "^[a-z0-9]{3,24}$"|
| azurerm_data_lake_analytics_firewall_rule| dlfw| 3| 50| false| "^[a-z0-9-_]{3,50}$"|
| azurerm_data_lake_store| dls| 3| 24| false| "^[a-z0-9]{3,24}$"|
| azurerm_data_lake_store_firewall_rule| dlsfw| 3| 50| false| "^[a-zA-Z0-9-_]{3,50}$"|
| azurerm_dev_test_lab| lab| 1| 50| false| "^[a-zA-Z0-9-_]{1,50}$"|
| azurerm_dev_test_linux_virtual_machine| labvm| 1| 64| false| "^[a-zA-Z0-9-]{1,64}$"|
| azurerm_dev_test_windows_virtual_machine| labvm| 1| 15| false| "^[a-zA-Z0-9-]{1,15}$"|
| azurerm_frontdoor| fd| 5| 64| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{3,62}[a-zA-Z0-9]$"|
| azurerm_frontdoor_firewall_policy| fdfw| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_hdinsight_hadoop_cluster| hadoop| 3| 59| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,57}[a-zA-Z0-9]$"|
| azurerm_hdinsight_hbase_cluster| hbase| 3| 59| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,57}[a-zA-Z0-9]$"|
| azurerm_hdinsight_kafka_cluster| kafka| 3| 59| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,57}[a-zA-Z0-9]$"|
| azurerm_hdinsight_interactive_query_cluster| iqr| 3| 59| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,57}[a-zA-Z0-9]$"|
| azurerm_hdinsight_ml_services_cluster| mls| 3| 59| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,57}[a-zA-Z0-9]$"|
| azurerm_hdinsight_rserver_cluster| rser| 3| 59| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,57}[a-zA-Z0-9]$"|
| azurerm_hdinsight_spark_cluster| spark| 3| 59| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,57}[a-zA-Z0-9]$"|
| azurerm_hdinsight_storm_cluster| storm| 3| 59| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,57}[a-zA-Z0-9]$"|
| azurerm_iotcentral_application| iotapp| 2| 63| true| "^[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$"|
| azurerm_iothub| iot| 3| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,48}[a-z0-9]$"|
| azurerm_iothub_consumer_group| iotcg| 1| 50| false| "^[a-zA-Z0-9-._]{1,50}$"|
| azurerm_iothub_dps| dps| 3| 64| false| "^[a-zA-Z0-9-]{1,63}[a-zA-Z0-9]$"|
| azurerm_iothub_dps_certificate| dpscert| 1| 64| false| "^[a-zA-Z0-9-._]{1,64}$"|
| azurerm_key_vault| kv| 3| 24| false| "^[a-zA-Z][a-zA-Z0-9-]{1,22}[a-zA-Z0-9]$"|
| azurerm_key_vault_key| kvk| 1| 127| false| "^[a-zA-Z0-9-]{1,127}$"|
| azurerm_key_vault_secret| kvs| 1| 127| false| "^[a-zA-Z0-9-]{1,127}$"|
| azurerm_key_vault_certificate| kvc| 1| 127| false| "^[a-zA-Z0-9-]{1,127}$"|
| azurerm_lb| lb| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_lb_nat_rule| lbnatrl| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_public_ip| pip| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_public_ip_prefix| pippf| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_route| rt| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_route_table| route| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_subnet| snet| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_traffic_manager_profile| traf| 1| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-.]{0,61}[a-zA-Z0-9_]$"|
| azurerm_virtual_wan| vwan| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_virtual_hub| vhub| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_virtual_network| vnet| 2| 64| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,62}[a-zA-Z0-9_]$"|
| azurerm_virtual_network_gateway| vgw| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_virtual_network_peering| vpeer| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_network_interface| nic| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_firewall| fw| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_eventhub| evh| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,48}[a-zA-Z0-9]$"|
| azurerm_eventhub_namespace| ehn| 1| 50| false| "^[a-zA-Z][a-zA-Z0-9-]{0,48}[a-zA-Z0-9]$"|
| azurerm_eventhub_authorization_rule| ehar| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,48}[a-zA-Z0-9]$"|
| azurerm_eventhub_namespace_authorization_rule| ehnar| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,48}[a-zA-Z0-9]$"|
| azurerm_eventhub_namespace_disaster_recovery_config| ehdr| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,48}[a-zA-Z0-9]$"|
| azurerm_eventhub_consumer_group| ehcg| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,48}[a-zA-Z0-9]$"|
| azurerm_stream_analytics_job| asa| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_stream_analytics_function_javascript_udf| asafunc| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_stream_analytics_output_blob| asaoblob| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_stream_analytics_output_mssql| asaomssql| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_stream_analytics_output_eventhub| asaoeh| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_stream_analytics_output_servicebus_queue| asaosbq| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_stream_analytics_output_servicebus_topic| asaosbt| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_stream_analytics_reference_input_blob| asarblob| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_stream_analytics_stream_input_blob| asaiblob| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_stream_analytics_stream_input_eventhub| asaieh| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_stream_analytics_stream_input_iothub| asaiiot| 3| 63| false| "^[a-zA-Z0-9-_]{3,63}$"|
| azurerm_shared_image_gallery| sig| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9.]{0,78}[a-zA-Z0-9]$"|
| azurerm_shared_image| si| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9]$"|
| azurerm_snapshots| snap| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_storage_account| st| 3| 24| true| "^[a-z0-9]{3,24}$"|
| azurerm_storage_container| stct| 3| 63| false| "^[a-z0-9][a-z0-9-]{2,62}$"|
| azurerm_storage_data_lake_gen2_filesystem| stdl| 3| 63| false| "^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$"|
| azurerm_storage_queue| stq| 3| 63| false| "^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$"|
| azurerm_storage_table| stt| 3| 63| false| "^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$"|
| azurerm_storage_share| sts| 3| 63| false| "^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$"|
| azurerm_storage_share_directory| sts| 3| 63| false| "^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$"|
| azurerm_machine_learning_workspace| mlw| 1| 260| false| "^[^<>*%:.?\\+\\/]{0,259}[^<>*%:.?\\+\\/ ]$"|
| azurerm_storage_blob| blob| 1| 1024| false| "^[^\\s\\/$#&]{1,1000}[^\\s\\/$#&]{0,24}$"|
| azurerm_bastion_host| snap| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_local_network_gateway| lgw| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_application_gateway| agw| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_express_route_gateway| ergw| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_express_route_circuit| erc| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_point_to_site_vpn_gateway| vpngw| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_template_deployment| deploy| 1| 64| false| "^[a-zA-Z0-9-._\\(\\)]{1,64}$"|
| azurerm_sql_server| sql| 1| 63| true| "^[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$"|
| azurerm_mssql_server| sql| 1| 63| true| "^[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$"|
| azurerm_mssql_database| sqldb| 1| 128| false| "^[^<>*%:.?\\+\\/]{1,127}[^<>*%:.?\\+\\/ ]$"|
| azurerm_sql_elasticpool| sqlep| 1| 128| false| "^[^<>*%:.?\\+\\/]{1,127}[^<>*%:.?\\+\\/ ]$"|
| azurerm_mssql_elasticpool| sqlep| 1| 128| false| "^[^<>*%:.?\\+\\/]{1,127}[^<>*%:.?\\+\\/ ]$"|
| azurerm_sql_failover_group| sqlfg| 1| 63| true| "^[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$"|
| azurerm_sql_firewall_rule| sqlfw| 1| 128| false| "^[^<>*%:?\\+\\/]{1,127}[^<>*%:.?\\+\\/]$"|
| azurerm_log_analytics_workspace| log| 4| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{2,61}[a-zA-Z0-9]$"|
| azurerm_service_fabric_cluster| sf| 4| 23| true| "^[a-z][a-z0-9-]{2,21}[a-z0-9]$"|
| azurerm_maps_account| map| 1| 98| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,97}$"|
| azurerm_network_watcher| nw| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_resource_group| rg| 1| 90| false| "^[a-zA-Z0-9-._\\(\\)]{0,89}[a-zA-Z0-9-_\\(\\)]$"|
| azurerm_network_security_group| nsg| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_network_security_group_rule| nsgr| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_network_security_rule| nsgr| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_application_security_group| asg| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_dns_zone| dns| 1| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,61}[a-zA-Z0-9_]$"|
| azurerm_private_dns_zone| pdns| 1| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,61}[a-zA-Z0-9_]$"|
| azurerm_notification_hub| nh| 1| 260| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,259}$"|
| azurerm_notification_hub_namespace| dnsrec| 6| 50| false| "^[a-zA-Z][a-zA-Z0-9-]{4,48}[a-zA-Z0-9]$"|
| azurerm_notification_hub_authorization_rule| dnsrec| 1| 256| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,255}$"|
| azurerm_servicebus_namespace| sb| 6| 50| false| "^[a-zA-Z][a-zA-Z0-9-]{4,48}[a-zA-Z0-9]$"|
| azurerm_servicebus_namespace_authorization_rule| sbar| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,48}[a-zA-Z0-9]$"|
| azurerm_servicebus_queue| sbq| 1| 260| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,258}[a-zA-Z0-9_]$"|
| azurerm_servicebus_queue_authorization_rule| sbqar| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,48}[a-zA-Z0-9]$"|
| azurerm_servicebus_subscription| sbs| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,48}[a-zA-Z0-9]$"|
| azurerm_servicebus_subscription_rule| sbsr| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,48}[a-zA-Z0-9]$"|
| azurerm_servicebus_topic| sbt| 1| 260| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,258}[a-zA-Z0-9]$"|
| azurerm_servicebus_topic_authorization_rule| dnsrec| 1| 50| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,48}[a-zA-Z0-9]$"|
| azurerm_powerbi_embedded| pbi| 3| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{2,62}$"|
| azurerm_dashboard| dsb| 3| 160| false| "^[a-zA-Z0-9-]{3,160}$"|
| azurerm_signalr_service| sgnlr| 3| 63| false| "^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]$"|
| azurerm_eventgrid_domain| egd| 3| 50| false| "^[a-zA-Z0-9-]{3,50}$"|
| azurerm_eventgrid_domain_topic| egdt| 3| 50| false| "^[a-zA-Z0-9-]{3,50}$"|
| azurerm_eventgrid_event_subscription| egs| 3| 64| false| "^[a-zA-Z0-9-]{3,64}$"|
| azurerm_eventgrid_topic| egt| 3| 50| false| "^[a-zA-Z0-9-]{3,50}$"|
| azurerm_relay_namespace| rln| 6| 50| false| "^[a-zA-Z][a-zA-Z0-9-]{4,48}[a-zA-Z0-9]$"|
| azurerm_relay_hybrid_connection| rlhc| 1| 260| false| "^[a-zA-Z0-9][a-zA-Z0-9-._]{0,258}[a-zA-Z0-9]$"|
cat resourceDefinition_out_of_docs.json | jq -r '.[] | "| \(.name)| \(.slug)| \(.min_length)| \(.max_length)| \(.lowercase)| \(.validation_regex)|"'
| azurerm_private_endpoint| pe| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_private_service_connection| psc| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_firewall_ip_configuration| fwipconf| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_firewall_application_rule_collection| fwapp| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_firewall_nat_rule_collection| fwnatrc| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_firewall_network_rule_collection| fwnetrc| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_dns_a_record| dnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_dns_aaaa_record| dnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_dns_caa_record| dnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_dns_cname_record| dnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_dns_mx_record| dnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_dns_ns_record| dnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_dns_ptr_record| dnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_dns_txt_record| dnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_private_dns_a_record| pdnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_private_dns_aaaa_record| pdnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_private_dns_cname_record| pdnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_private_dns_mx_record| pdnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_private_dns_ptr_record| pdnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_private_dns_srv_record| pdnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_private_dns_txt_record| pdnsrec| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_virtual_machine_extension| vmx| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_virtual_machine_scale_set_extension| vmssx| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_network_ddos_protection_plan| ddospp| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_private_dns_zone_group| pdnszg| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_proximity_placement_group| ppg| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| azurerm_private_link_service| pls| 1| 80| false| "^[a-zA-Z0-9][a-zA-Z0-9\\-\\._]{0,78}[a-zA-Z0-9_]$"|
| databricks_cluster| dbc| 3| 30| false| "^[a-zA-Z0-9-_]{3,30}$"|
| databricks_standard_cluster| dbsc| 3| 30| false| "^[a-zA-Z0-9-_]{3,30}$"|
| databricks_high_concurrency_cluster| dbhcc| 3| 30| false| "^[a-zA-Z0-9-_]{3,30}$"|
