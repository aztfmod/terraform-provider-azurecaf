# Azurecaf provider

The Azurecaf provider is a *logical provider* which means that it works entirely within Terraform's logic, and doesn't interact with any other services. The goal of this provider is to provider helper methods in implementing Azure landing zones using Terraform.

The Azurecaf provider currently contains a single resource based on the Terraform Random_string provider. The naming_convention resources enforce consistant naming covention for a set of supported Azure services.

You may select different type of naming convention (cafclassic,cafrandom,random,passthrough) based on the environment that you target or the naming style that you need to apply. Whichever convention you select, the naming convention ensures that the generated name is compliant with the Azure service that you target.

## Resource types

We define resource types as per: https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging 

Current prototype supports:

| Resource type                       | Resource type code (short)  | Resource type code (long)               |
| ----------------------------------- | ----------------------------|-----------------------------------------|
| Azure Automation                    | aaa                         | azurerm_automation_account              |
| Azure Container Registry            | acr                         | azurerm_container_registry              |
| Azure Firewall                      | afw                         | azurerm_firewall                        |
| Application Gateway                 | agw                         | azurerm_application_gateway             |
| API Management                      | apim                        | azurerm_api_management                  |
| App Service                         | app                         | azurerm_app_service                     |
| Application Insights                | appi                        | azurerm_application_insights            |
| App Service Environment             | ase                         | azurerm_app_service_environment         |
| Application Security Group          | asg                         | azurerm_app_security_group              |
| Azure Kubernetes Service            | aks                         | azurerm_kubernetes_cluster              |
| Azure Kubernetes Service DNS prefix | aksdns                      | aks_dns_prefix                          |
| AKS Node Pool Linux                 | aksnpl                      | aks_node_pool_linux                     |
| AKS Node Pool Windows               | aksnpw                      | aks_node_pool_windows                   |
| Azure Site Recovery                 | asr                         | azurerm_recovery_services_vault         |
| Azure Availability Set              | avail                       | azurerm_availability_set                |
| Azure Vpn Connection                | cn                          | azurerm_vpn_connection                  |
| Azure Event Hubs                    | evh                         | azurerm_eventhub_namespace              |
| generic                             | gen                         | generic                                 |
| Azure Key Vault                     | kv                          | azurerm_key_vault                       |
| Azure Monitor Log Analytics         | la                          | azurerm_log_analytics_workspace         |
| Azure Load Balancer (External)      | lbe                         | azurerm_load_balancer_external          |
| Azure Load Balancer (Internal)      | lbi                         | azurerm_load_balancer_internal          |
| Azure Local Network Gateway         | lgw                         | azurerm_local_network_gateway           |
| Azure Mysql Database                | mysql                       | azurerm_mysql_database                  |
| Virtual Network Interface Card      | nic                         | azurerm_network_interface               |
| Network Security Group              | nsg                         | azurerm_network_security_group          |
| Public IP                           | pip                         | azurerm_public_ip                       |
| App Service Plan                    | plan                        | azurerm_app_service_plan                |
| Resource group                      | rg                          | azurerm_resource_group                  |
| Azure Route Table                   | route                       | azurerm_route_table                     |
| Azure Service Bus                   | sb                          | azurerm_service_bus                     |
| Azure Service Bus Queue             | sbq                         | azurerm_service_bus_queue               |
| Azure Service Bus Topic             | sbt                         | azurerm_service_bus_topic               |
| Subnet                              | snet                        | azurerm_subnet                          |
| Azure SQL DB Server                 | sql                         | azurerm_sql_server                      |
| Azure SQL DB                        | sqldb                       | azurerm_sql_database                    |
| Azure Storage Account               | st                          | azurerm_storage_account                 |
| Azure Traffic Manager Profile       | traf                        | azurerm_traffic_manager_profile         |
| Virtual Network Gateway             | vgw                         | azurerm_virtual_network_gateway         |
| Linux Virtual Machine               | vml                         | azurerm_virtual_machine_linux           |
| Virtual Machine Scale Set Linux     | vmssl                       | azurerm_vm_scale_set_linux              |
| Virtual Machine Scale Set Windows   | vmssw                       | azurerm_vm_scale_set_windows            |
| Windows Virtual Machine             | vmw                         | azurerm_virtual_machine_windows         |
| Virtual Network                     | vnet                        | azurerm_virtual_network                 |