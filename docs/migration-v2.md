# Migrating from v1.x to v2.0.0

Version 2.0.0 removes the legacy naming resource, aligns CAF abbreviations with current Microsoft guidance, and corrects several Azure naming constraints. Review a saved plan before applying because changed generated names can force replacement of downstream Azure resources.

## Before upgrading

1. Upgrade to the latest v1 release and apply until `terraform plan` is clean.
2. Record generated names with `terraform state show`.
3. Search for removed resources and affected resource types:

   ```bash
   grep -rE 'azurecaf_naming_convention|azurerm_app_service_plan|azurerm_lb_backend_pool|azurerm_role_assignment' .
   ```

4. Decide whether each slug change should rename the Azure resource or retain its current name.

## Replace `azurecaf_naming_convention`

The `azurecaf_naming_convention` resource is no longer registered. Migrate it while still using provider v1.x.

| v1 argument | `azurecaf_name` equivalent |
|---|---|
| `resource_type = "rg"` | `resource_type = "azurerm_resource_group"` |
| `prefix = "dev"` | `prefixes = ["dev"]` |
| `postfix = "001"` | `suffixes = ["001"]` |
| `convention = "cafclassic"` | `random_length = 0` |
| `convention = "cafrandom"` | Set the required `random_length` |
| `convention = "random"` | Omit `name`, set `random_length`, and set `use_slug = false` |
| `convention = "passthrough"` | `passthrough = true` |

To retain an existing generated name without replacing downstream resources:

```hcl
resource "azurecaf_name" "current" {
  name          = "stdevmyapp001"
  resource_type = "azurerm_storage_account"
  passthrough   = true
}
```

Import the replacement while still on v1.x, update references, and then remove the old state:

```bash
terraform import azurecaf_name.current \
  'azurerm_storage_account:stdevmyapp001'
terraform plan
terraform state rm azurecaf_naming_convention.current
```

Do not use a `moved` block between these resource types; their state schemas are different.

## Replace removed resource definitions

| Removed `resource_type` | Replacement | Naming change |
|---|---|---|
| `azurerm_app_service_plan` | `azurerm_service_plan` | Slug changes from `plan` to `asp` |
| `azurerm_lb_backend_pool` | `azurerm_lb_backend_address_pool` | Uses the supported backend address pool definition |
| `azurerm_role_assignment` | Omit `name` or generate a UUID directly | AzureRM requires a tenant-unique UUID; a CAF-prefixed name is invalid |

## CAF abbreviation changes

The following generated names change when `use_slug = true`.

| Resource type | v1.x | v2.0.0 |
|---|---:|---:|
| `azurerm_app_configuration` | `appcg` | `appcs` |
| `azurerm_bastion_host` | `bast` | `bas` |
| `azurerm_cdn_endpoint` | `cdn` | `cdne` |
| `azurerm_cdn_profile` | `cdnprof` | `cdnp` |
| `azurerm_container_group` | `aci` | `ci` |
| `azurerm_dns_zone` | `dns` | `dnsz` |
| `azurerm_eventgrid_domain` | `egd` | `evgd` |
| `azurerm_eventgrid_domain_topic` | `egdt` | `evgdt` |
| `azurerm_eventgrid_event_subscription` | `egs` | `evgs` |
| `azurerm_eventgrid_topic` | `egt` | `evgt` |
| `azurerm_eventhub_authorization_rule` | `ehar` | `evhar` |
| `azurerm_eventhub_consumer_group` | `ehcg` | `evhnscg` |
| `azurerm_eventhub_namespace` | `ehn` | `evhns` |
| `azurerm_eventhub_namespace_authorization_rule` | `ehnar` | `evhnsar` |
| `azurerm_eventhub_namespace_disaster_recovery_config` | `ehdr` | `evhnsdr` |
| `azurerm_firewall` | `fw` | `afw` |
| `azurerm_frontdoor` | `fd` | `afd` |
| `azurerm_function_app` | `fa` | `func` |
| `azurerm_function_app_slot` | `fas` | `funcs` |
| `azurerm_kusto_cluster` | `kc` | `dec` |
| `azurerm_kusto_database` | `kdb` | `dedb` |
| `azurerm_kusto_eventhub_data_connection` | `kehc` | `deevhdcon` |
| `azurerm_lb_backend_address_pool` | `adt` | `lbbepool` |
| `azurerm_lb_nat_pool` | `adt` | `lbnatpool` |
| `azurerm_lb_nat_rule` | `lbnatrl` | `lbnatr` |
| `azurerm_lb_outbound_rule` | `adt` | `lbor` |
| `azurerm_lb_probe` | `adt` | `probe` |
| `azurerm_lb_rule` | `adt` | `rule` |
| `azurerm_logic_app_action_custom` | `lappac` | `logicac` |
| `azurerm_logic_app_action_http` | `lappah` | `logicah` |
| `azurerm_logic_app_integration_account` | `lappia` | `ia` |
| `azurerm_logic_app_trigger_custom` | `lapptc` | `logictc` |
| `azurerm_logic_app_trigger_http_request` | `lappth` | `logicth` |
| `azurerm_logic_app_trigger_recurrence` | `lapptc` | `logictc` |
| `azurerm_logic_app_workflow` | `lapp` | `logic` |
| `azurerm_managed_disk` | `dsk` | `disk` |
| `azurerm_monitor_action_group` | `amag` | `ag` |
| `azurerm_monitor_activity_log_alert` | `adfmysql` | `amala` |
| `azurerm_network_security_group_rule` | `nsgr` | `nsgsr` |
| `azurerm_notification_hub` | `nh` | `ntf` |
| `azurerm_notification_hub_authorization_rule` | `dnsrec` | `ntfar` |
| `azurerm_notification_hub_namespace` | `dnsrec` | `ntfns` |
| `azurerm_private_dns_zone` | `pdns` | `pdnsz` |
| `azurerm_private_link_service` | `pls` | `pl` |
| `azurerm_public_ip_prefix` | `pippf` | `ippre` |
| `azurerm_purview_account` | `purv` | `pview` |
| `azurerm_route` | `rt` | `udr` |
| `azurerm_route_table` | `route` | `rt` |
| `azurerm_servicebus_namespace` | `sb` | `sbns` |
| `azurerm_shared_image_gallery` | `sig` | `gal` |
| `azurerm_signalr_service` | `sgnlr` | `sigr` |
| `azurerm_synapse_firewall_rule` | `syfw` | `synfw` |
| `azurerm_synapse_spark_pool` | `sysp` | `synsp` |
| `azurerm_synapse_sql_pool` | `synsp` | `syndp` |
| `azurerm_synapse_sql_pool_vulnerability_assessment_baseline` | `synspvab` | `syndpvab` |
| `azurerm_synapse_sql_pool_workload_classifier` | `synspwc` | `syndpwc` |
| `azurerm_synapse_sql_pool_workload_group` | `synspwg` | `syndpwg` |
| `azurerm_synapse_workspace` | `syws` | `synw` |
| `azurerm_user_assigned_identity` | `msi` | `id` |
| `azurerm_virtual_desktop_application_group` | `dag` | `vdag` |
| `azurerm_virtual_desktop_host_pool` | `hpool` | `vdpool` |
| `azurerm_virtual_desktop_workspace` | `wvdws` | `vdws` |
| `azurerm_virtual_network_peering` | `vpeer` | `peer` |
| `azurerm_web_application_firewall_policy` | `wafw` | `waf` |

To preserve an existing name, use `passthrough = true` with the exact current value or set `use_slug = false` and compose the legacy abbreviation explicitly.

## Naming constraint changes

| Resource type | v2.0.0 behavior |
|---|---|
| Load Balancer child resources | Scope is the parent Load Balancer; length is 1-80; periods are allowed |
| `azurerm_custom_provider` | AzureRM-compatible names contain 3-64 alphanumeric characters without separators |
| `azurerm_private_link_service` | Length is 2-64 |
| `azurerm_storage_table` | Starts with a letter and contains only 3-63 alphanumeric characters |
| `azurerm_synapse_workspace` | Globally unique, 1-50 characters, lowercase alphanumeric and hyphens |
| `azurerm_web_application_firewall_policy` | Resource-group scope, 1-80 characters, allows `-`, `.`, and `_` |
| API Management and Load Test names | A trailing pipe (`|`) is rejected |
| New and corrected definitions with `min_length = 1` | Validation accepts the documented one-character minimum |

If an existing generated name does not satisfy a corrected rule, v2.0.0 will return a validation error. Keep the v1 provider pinned until the downstream name is migrated or stop deriving that existing name through `azurecaf_name`.

## Seeded random output changes

Version 1.x incorrectly sampled only the first 25 characters of the 26-character random alphabet, so `z` could never appear. Version 2.0.0 corrects the range. As a result, the same explicit `random_seed` can produce a different suffix:

| `random_seed` | `random_length` | v1.x suffix | v2.0.0 suffix |
|---:|---:|---|---|
| `42` | `12` | `fmsaxuhbdser` | `hrukpttuezpt` |

This affects seeded `azurecaf_name` data sources as soon as they are read by v2. Existing `azurecaf_name` resource results remain stored during an unchanged refresh, but the new sequence is used when a resource is created, replaced, or recomputed after an input change. Record current seeded results before upgrading and inspect downstream replacements in the plan. To retain an exact v1 name, use `passthrough = true` with the recorded name while migrating the downstream resource; do not rely on the old seed reproducing it.

## Plan-time generated names

`azurecaf_name` resources now calculate deterministic names during planning. When `random_length > 0`, set a non-zero seed to expose the final name in the plan:

```hcl
resource "azurecaf_name" "example" {
  name          = "myapp"
  resource_type = "azurerm_storage_account"
  random_length = 5
  random_seed   = 12345
}
```

If the seed is omitted or is `0`, the resource remains non-deterministic and its result is known after apply, then preserved in state. For the `azurecaf_name` data source, an explicitly configured seed is deterministic and `0` is valid; omitting the seed uses fresh randomness on every read. Set an explicit data-source seed whenever its result must remain stable across plans. Resource results also remain known after apply when any naming input, including a list element, depends on a value that Terraform cannot know during planning.

## Upgrade checklist

1. Migrate every `azurecaf_naming_convention` resource on v1.x.
2. Replace removed resource types.
3. Review all affected slug, constraint, and seeded-random output changes.
4. Add non-zero resource seeds where plan-time visibility is required and explicit data-source seeds where stable names are required.
5. Update the provider constraint to `~> 2.0`.
6. Run `terraform init -upgrade` and inspect the complete plan.
7. Apply name changes during an appropriate maintenance window.
