# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Weekly Azure Sync 2026-05-14 — 105 new resource types**: Added definitions for resources highlighted by the weekly `azure-sync` report. All entries are additive (no slug or behavior changes to existing resources) and follow CAF-style abbreviations and Azure naming rules. Categories covered:
  - Compute: `azurerm_container_group` (slug `aci`), `azurerm_orchestrated_virtual_machine_scale_set` (`ovmss`), `azurerm_snapshot` (`snp`), `azurerm_dedicated_hardware_security_module` (`hsm`).
  - Networking: `azurerm_network_profile` (`npr`), `azurerm_route_filter` (`rf`), `azurerm_vpn_gateway` (`vpng`), `azurerm_virtual_network_gateway_connection` (`vngc`), `azurerm_express_route_circuit_authorization` (`erca`), `azurerm_express_route_circuit_peering` (`ercp`), `azurerm_firewall_policy_rule_collection_group` (`fwprcg`), `azurerm_virtual_hub_route_table` (`vhrt`), `azurerm_virtual_hub_ip` (`vhip`), `azurerm_virtual_hub_bgp_connection` (`vhbgp`), `azurerm_virtual_hub_security_partner_provider` (`vhspp`).
  - Storage & Databases: `azurerm_sql_database` (`sqld`), `azurerm_mssql_virtual_machine` (`sqlvm`), `azurerm_cosmosdb_cassandra_keyspace` (`coscas`), `azurerm_cosmosdb_gremlin_database` (`cosgrm`), `azurerm_cosmosdb_gremlin_graph` (`cosgrmg`), `azurerm_cosmosdb_mongo_collection` (`cosmonc`), `azurerm_cosmosdb_mongo_database` (`cosmondb`), `azurerm_cosmosdb_sql_container` (`cosqlc`), `azurerm_cosmosdb_sql_database` (`cosqldb`), `azurerm_cosmosdb_sql_stored_procedure` (`cosqlsp`), `azurerm_cosmosdb_table` (`costbl`), `azurerm_storage_encryption_scope` (`stes`), `azurerm_storage_data_lake_gen2_path` (`stdlp`), `azurerm_shared_image_version` (`siv`).
  - IoT & Messaging: `azurerm_iothub_endpoint_eventhub` (`iothepeh`), `azurerm_iothub_endpoint_servicebus_queue` (`iothepsbq`), `azurerm_iothub_endpoint_servicebus_topic` (`iothepsbt`), `azurerm_iothub_endpoint_storage_container` (`iothepsc`), `azurerm_iothub_fallback_route` (`iothfr`), `azurerm_iothub_route` (`iothr`), `azurerm_iot_time_series_insights_standard_environment` (`tsise`), `azurerm_iot_time_series_insights_reference_data_set` (`tsirds`), `azurerm_eventhub_cluster` (`ehc`).
  - Analytics & Data: `azurerm_data_share` (`dshr`), `azurerm_data_share_account` (`dshra`), `azurerm_data_share_dataset_blob_storage` (`dshrdsb`), `azurerm_data_share_dataset_data_lake_gen1` (`dshrdsg1`), `azurerm_data_share_dataset_data_lake_gen2` (`dshrdsg2`), `azurerm_data_share_dataset_kusto_cluster` (`dshrdskc`), `azurerm_data_share_dataset_kusto_database` (`dshrdskd`), `azurerm_hpc_cache` (`hpcc`), `azurerm_hpc_cache_blob_target` (`hpcbt`), `azurerm_hpc_cache_nfs_target` (`hpcnt`), `azurerm_hdinsight_cluster` (`hdi`), `azurerm_data_factory_integration_runtime_self_hosted` (`adfirsh`), `azurerm_kusto_attached_database_configuration` (`kadc`).
  - Security & Governance: `azurerm_advanced_threat_protection` (`atp`), `azurerm_attestation` (`atst`), `azurerm_security_center_automation` (`sca`), `azurerm_security_center_auto_provisioning` (`scap`), `azurerm_security_center_contact` (`scc`), `azurerm_sentinel_alert_rule` (`sentar`), `azurerm_sentinel_alert_rule_ms_security_incident` (`sentarms`), `azurerm_sentinel_alert_rule_scheduled` (`sentars`), `azurerm_lighthouse_assignment` (`lha`), `azurerm_lighthouse_definition` (`lhd`), `azurerm_policy_definition` (`pold`), `azurerm_policy_set_definition` (`polsd`), `azurerm_policy_remediation` (`polr`), `azurerm_management_lock` (`mgl`), `azurerm_management_group` (`mg`), `azurerm_key_vault_access_policy` (`kvap`).
  - App Services & Containers: `azurerm_app_service_slot` (`apps`), `azurerm_app_service_certificate` (`appcert`), `azurerm_service_plan` (`asp`), `azurerm_spring_cloud_app` (`spca`), `azurerm_spring_cloud_certificate` (`spcert`), `azurerm_spring_cloud_service` (`spcs`), `azurerm_spatial_anchors_account` (`spaa`).
  - Monitoring & Operations: `azurerm_monitor_log_profile` (`mlp`), `azurerm_monitor_action_rule_action_group` (`marag`), `azurerm_monitor_action_rule_suppression` (`mars`), `azurerm_monitor_scheduled_query_rules_log` (`msqrl`), `azurerm_monitor_smart_detector_alert_rule` (`msdar`), `azurerm_log_analytics_data_export_rule` (`laer`), `azurerm_log_analytics_linked_service` (`lals`), `azurerm_log_analytics_saved_search` (`lass`), `azurerm_log_analytics_cluster_customer_managed_key` (`laccmk`).
  - Automation & DevOps: `azurerm_automation_connection` (`aacon`), `azurerm_automation_module` (`aamod`), `azurerm_automation_dsc_configuration` (`aadsc`), `azurerm_blueprint_assignment` (`bpa`), `azurerm_blueprint_definition` (`bpd`), `azurerm_blueprint_published_version` (`bppv`), `azurerm_dev_test_virtual_network` (`dtlvn`), `azurerm_dev_test_policy` (`dtlp`), `azurerm_dev_test_schedule` (`dtls`), `azurerm_resource_group_template_deployment` (`rgtd`), `azurerm_subscription_template_deployment` (`subtd`).
  - Other / Misc: `azurerm_managed_application` (`manapp`), `azurerm_managed_application_definition` (`manappd`), `azurerm_media_services_account` (`ams`), `azurerm_devspace_controller` (`dsc`), `azurerm_service_fabric_mesh_application` (`sfmesha`), `azurerm_service_fabric_mesh_local_network` (`sfmeshln`), `azurerm_service_fabric_mesh_secret` (`sfmeshs`), `azurerm_site_recovery_fabric` (`asrf`), `azurerm_site_recovery_replicated_vm` (`asrrvm`), `azurerm_site_recovery_replication_policy` (`asrrp`), `azurerm_site_recovery_protection_container` (`asrpc`).
  - Impact: Low - additive only. Existing resource slugs are unchanged.
  - Note: The 3 CAF slug drifts reported in the same issue (`ehn`→`evhns`, `syws`→`synw`, `dpbv`→`bvault`) are **not** included here. Changing existing slugs is breaking for downstream users and should be handled in a separate change with an aliasing/deprecation strategy.
- **Network Connection Monitor Support**: Added support for `azurerm_network_connection_monitor` resource type
  - Resource slug: `cm`
  - Min length: 1, Max length: 80
  - Scope: parent (child resource of Network Watcher)
  - Allows alphanumerics, hyphens, periods, and underscores
  - Follows Azure naming conventions for Network Connection Monitor resources
  - Impact: Low - Adds new resource type support for Azure Network Watcher connection monitoring
- **Missing name resources from tracking issue [#432](https://github.com/aztfmod/terraform-provider-azurecaf/issues/432)**: Added four resource type entries that were previously requested but unsupported. Each entry was researched against the Azure naming-rules documentation (`Microsoft.Network/natGateways`) and CAF abbreviations page; resources without an official CAF abbreviation are flagged `out_of_doc: true`.
  - `azurerm_nat_gateway` — slug `ng` (official CAF), scope `resourceGroup`, min 1 / max 80, allows alphanumerics, hyphens, periods, and underscores. Resource provider namespace: `Microsoft.Network/natGateways`. Closes [#254](https://github.com/aztfmod/terraform-provider-azurecaf/issues/254).
  - `azurerm_monitor_workspace` — slug `amw`, scope `resourceGroup`, min 4 / max 63, allows alphanumerics and hyphens. Azure Monitor managed Prometheus workspace. Closes [#276](https://github.com/aztfmod/terraform-provider-azurecaf/issues/276).
  - `azurerm_email_communication_service` — slug `acsmail`, scope `global`, min 1 / max 63, allows alphanumerics and hyphens. Companion to the existing `azurerm_communication_service` entry (`acs`). Closes [#261](https://github.com/aztfmod/terraform-provider-azurecaf/issues/261).
  - `azurerm_vpn_server_configuration` — slug `vpnsc`, scope `resourceGroup`, min 1 / max 80, allows alphanumerics, hyphens, periods, and underscores. Companion to the existing `azurerm_vpn_gateway_connection` (`vcn`) and `azurerm_vpn_site` (`vst`) entries. Closes [#174](https://github.com/aztfmod/terraform-provider-azurecaf/issues/174).
  - Impact: Low - additive only, no breaking changes; existing names are unchanged.

### Changed
- **Dependencies**: Bumped `github.com/hashicorp/terraform-plugin-sdk/v2` from v2.38.2 to v2.40.0
  - Includes resource configuration generation logic for `-generate-config-out` flag (Terraform v1.14.0+)
  - Added deprecation message support for attributes and blocks
  - Go version updated from 1.24.4 to 1.25.0
  - Aligned `e2e/go.mod` dependencies (`terraform-exec` v0.25.0, `terraform-json` v0.27.2, `go-cty` v1.17.0)
  - Impact: Low -- dependency update only, no breaking changes
- **Dev Environment**: Bumped `.devcontainer/devcontainer.json` Go feature from `1.24.4` to `1.25.0` to match `go.mod` (`go 1.25.0`). Previous pin would have failed `go build` inside the dev container.
- **Documentation**: Updated `.github/CONTRIBUTING.md` prerequisite from `Go 1.24.4+` to `Go 1.25.0+` to match `go.mod`.
- **`.gitignore`**: Added `.copilot-tracking/` to exclude local Doc-Ops session tracking files from version control.

### Security
- **CI/Automation, top-level workflow permissions**: Added explicit `permissions: contents: read` block at the top of `.github/workflows/codeql.yml` and `.github/workflows/copilot-setup-steps.yml`. Both workflows previously declared per-job permissions only, which left the top-level `GITHUB_TOKEN` defaulting to the repo/org-level setting (historically `write-all`). The per-job blocks already declare the minimum-required elevated scopes (`security-events: write`, `packages: read`, `actions: read` for CodeQL; `contents: read` for Copilot Setup Steps), so the top-level default tightens the blast radius for any future job that omits its own `permissions:` block. Closes Checkov [CKV2_GHA_1](https://docs.bridgecrew.io/docs/ensure-top-level-permissions-are-not-set-to-write-all) code-scanning alert [#2](https://github.com/aztfmod/terraform-provider-azurecaf/security/code-scanning/2). Mitigates CWE-732 (Incorrect Permission Assignment).
- **CI/Automation**: SHA-pinned `actions/checkout` and `actions/setup-go` in hand-authored workflows (`codeql.yml`, `copilot-setup-steps.yml`, `e2e.yml`, `go.yml`, `security.yml`) to match the SHA-pinning posture of the gh-aw-managed agentic lock files. Pinned `actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd` (v6.0.2) and `actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c` (v6.4.0). Mitigates tag-mutability risk (CWE-829).

### Fixed
- **CI/Automation, `release-validation`**: Added a `hashicorp/setup-terraform@v4` (with `terraform_wrapper: false`) and `actions/setup-go@v6` pre-agent step block to `.github/workflows/release-validation.md`, plus a `Install tfproviderlint` step and a host-side run-script that builds the provider and runs `make test_ci`, `make test_e2e`, `make test_coverage`, and the `CHANGELOG.md` tag-section check. Exit codes, the parsed coverage percentage, and matched CHANGELOG headings are written to `/tmp/results.env` / `/tmp/changelog-entry.txt`, and full logs to `/tmp/{build,test-ci,e2e,coverage}-output.txt`. Removed `make *` and `go *` from the agent's allowed `bash` toolset; the agent now only consumes pre-computed files with `cat`/`grep`/`head`/`tail`. Previously the agent failed every tagged release with `bash command execution is blocked by security policy` and `Go toolchain (go, make) commands are blocked by runner security policy`, because the gh-aw firewall sandbox lacks a Go toolchain and denies shell command execution. Recompiled `release-validation.lock.yml` with `gh aw` v0.72.1; approved the new `hashicorp/setup-terraform@v4` action pin (`dfe3c3f87815947d99a8997f908cb6525fc44e9e`). Closes [#478](https://github.com/aztfmod/terraform-provider-azurecaf/issues/478).
- **CI/Automation, `nightly-regression`**: Install `tfproviderlint` on the host runner before the agent step runs `make test_ci`, and soften the `Makefile` `unittest` target to skip the lint step (with a warning) when `tfproviderlint` is not on `PATH`. The nightly run on 2026-05-14 failed with `make: tfproviderlint: No such file or directory` (exit 127) because `make test_ci` → `make unittest` shells out to `tfproviderlint ./...` but `nightly-regression.md` never installed the tool (unlike the hand-authored `go.yml` workflow, which already has an install step). Added a `go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest` step after `Set up Go` and recompiled `nightly-regression.lock.yml` with `gh aw` v0.72.1. The `Makefile` guard is defense-in-depth: developers on fresh checkouts and other agentic workflows (e.g. `release-validation.md`) that consume `make unittest` / `make test_ci` no longer hit a hard failure if the optional lint tool is missing. Closes [#472](https://github.com/aztfmod/terraform-provider-azurecaf/issues/472).
- **CI/Automation, `nightly-regression`**: Added a `hashicorp/setup-terraform@v4` (with `terraform_wrapper: false`) and `actions/setup-go@v6` pre-agent step to `.github/workflows/nightly-regression.md`. The nightly E2E suite (`make test_e2e_quick`) shells out to a real `terraform` binary (`e2e/e2e_test.go:80`); without it the suite failed every night with `terraform binary not found in PATH`. Also reworked the run-script to use `set -o pipefail` and capture each suite's exit code in an `if`/`else` (assignment after `tee` was always seeing `$?` of `tee`, so `E2E_EXIT` was reported as `0` even when E2E failed). The step now always exits `0` so the agent can run and decide whether to file an issue.
- **CI/Automation, `issue-arborist`**: Moved the pre-downloaded `issues.json` and `issues-schema.json` from `/tmp/gh-aw/issues-data/` to `${GITHUB_WORKSPACE}/.gh-aw-data/`. The gh-aw v0.72.1 agent firewall reserves `/tmp/gh-aw/` for runtime files and denies application reads from the `issues-data/` subpath, so the agent could not `cat` or `jq` the file the previous step had just written. Updated all path references in the agent prompt body. (Tracking the upstream `github/gh-aw` `issue-arborist.md` source-of-truth — diverges from `@852cb06ad…` until upstream adopts the same fix.)
- **CI/Automation, `weekly-azure-sync`**: Moved the resource-discovery compute (supported list, azurerm list, diff, counts, CAF abbreviations fetch) into a pre-agent step that writes `${GITHUB_WORKSPACE}/.gh-aw-data/{supported,azurerm,missing,counts.env,caf-abbreviations.html}`. Removed `python3 *` and `curl *` from the agent's allowed `bash` toolset. The agent now only consumes the pre-computed files with `cat`/`grep`/`head`/`tail`/`sort`/`comm`/`wc`. Previously the agent attempted to run `python3` inside the firewall sandbox, which is denied — see [#468](https://github.com/aztfmod/terraform-provider-azurecaf/issues/468).
- **CI/Automation**: Added a pre-agent `Set up Go` + `Build provider` step block to `pr-review-agent.md` so the PR review agent can verify `make build` against the repo's Go version (`go.mod` → `1.25.0`). Previously the agent reported "Build passes ⚠️ Unverified — Go toolchain unavailable in sandbox" because the gh-aw agent container (`ghcr.io/github/gh-aw-firewall/agent`) does not ship a Go toolchain. The build now runs on the host runner (where `setup-go` populates the tool cache) and the agent reads `BUILD_RESULT` and a 80-line log tail from `/tmp/`. **Requires recompiling** with `gh aw compile pr-review-agent` to regenerate `pr-review-agent.lock.yml`.
- **CI/Automation**: Recompiled all 10 GitHub Agentic Workflow lock files with `gh aw` v0.72.1 (previously v0.61.0)
  - Generated missing `.lock.yml` files for `contributor-welcome`, `issue-to-pr-agent`, `nightly-regression`, `pr-review-agent`, `release-validation`, and `weekly-azure-sync`
- **CI/Automation**: Migrated agentic workflows to current gh-aw schema
  - Replaced deprecated `tools.github.repos` with `tools.github.allowed-repos` in `daily-repo-status` and `release-labeler`
  - Removed redundant `contents: write` / `issues: write` permissions now handled by `safe-outputs` (strict-mode requirement) in `issue-to-pr-agent`, `nightly-regression`, `weekly-azure-sync`
  - Added missing toolset read permissions (`pull-requests: read`, `issues: read`) required by declared GitHub toolsets
  - Switched fixed cron expressions to fuzzy schedules (`daily`, `weekly on monday`) in `nightly-regression` and `weekly-azure-sync` to spread load
  - Replaced `registry.terraform.io` domain in `weekly-azure-sync` network allowlist with the `terraform` ecosystem identifier

### Documentation
- **Resource count alignment**: Updated documented Azure resource type count to the actual `405` (verified via `jq 'length' resourceDefinition.json`) across `README.md`, `COMPLETE_TESTING_GUIDE.md`, `docs/index.md`, `docs/resources/azurecaf_name.md`, `docs/resources/azurecaf_naming_convention.md`, `docs/data-sources/azurecaf_name.md`, and `.github/CONTRIBUTING.md`. Reflects the four new resources added under tracking issue [#432](https://github.com/aztfmod/terraform-provider-azurecaf/issues/432).
- **README.md run-on regressions**: Fixed three rendering bugs introduced by an earlier bulk substitution where a missing newline collapsed adjacent list items / fenced-code lines (lines 398, 418, and 465 — `Resource Matrix Tests` / `Constraint Tests`, `100% Resource Coverage` / `Naming Validation`, and the comprehensive testing framework code block). Items now render as separate bullets / lines.
- **README.md broken make targets**: Replaced two references to the non-existent `make test_resource_constraints` target with `make test_resource_matrix` (whose Makefile help string already says "Test resources by category and validate constraints"). The duplicate row was de-duplicated.
- **README.md missing argument**: Added `error_when_exceeding_max_length` to the Supported Parameters table; the schema (`resource_name.go:238`) exposes the field but the table omitted it.
- **TESTING.md stale file references**: Removed two references to integration test files that do not exist (`integration_cross_resource_test.go`, `integration_naming_convention_types_test.go`) and replaced them with `integration_all_resources_complete_test.go`. Updated the "where to add a new test" guidance accordingly.
- **TESTING.md / E2E_IMPLEMENTATION_SUMMARY.md**: Replaced two hedging phrases (`Don't just test the happy path` → `Cover edge cases, not only the happy path`; `easy to extend` → `extensible without refactoring`).
- **e2e/README.md**: Converted a dangling `- **Test Categories**:` empty-body bullet into a `#### Test Categories` subheading so the four sub-bullets render as a section instead of as orphan list children.

### Changed
- **Resource status table**: Added `azurerm_linux_function_app`, `azurerm_linux_function_app_slot`, `azurerm_managed_redis`, `azurerm_windows_function_app`, `azurerm_windows_function_app_slot`, and `azurerm_windows_web_app` to the README.md and docs/index.md resource tables. These resources were already supported in `resourceDefinition.json` but the user-facing tables were stale.
- **Missing argument doc**: Added `error_when_exceeding_max_length` to `docs/resources/azurecaf_name.md` (was already documented for the data source but absent from the resource doc since v1.2.32).
- **Go prerequisite**: TESTING.md updated `Go 1.19+` → `Go 1.25.0+` to match `go.mod` (`go 1.25.0`).
- **Version pin**: Bumped `version = "~> 1.2.28"` example pin to `~> 1.2.32` (latest tag) in README.md and docs/index.md.
- **SECURITY.md**: Replaced unmodified GitHub template placeholder text with an actual Supported Versions statement.
- **Encoding**: Fixed U+FFFD replacement character in E2E_IMPLEMENTATION_SUMMARY.md heading.
- **Deprecation**: Strengthened deprecation notice on `docs/resources/azurecaf_naming_convention.md` to a top-of-page Deprecated callout pointing at `azurecaf_name` and the migration guide (the source code already marks the resource `Deprecated:`; the doc framing was too soft).
- **`azurecaf_naming_convention` algorithm doc**: Removed the incorrect "Name Composition and Truncation" section that had been copy-pasted from `azurecaf_name.md`. It documented arguments (`use_slug`, `random_length`, `clean_input`, `passthrough` as flag, `separator`) that do not exist on this resource's schema, and described the `azurecaf_name` truncation pipeline rather than the legacy `getResult` algorithm. Replaced with a concise, accurate description of the `[prefix, cafprefix, name, postfix]` composition and per-`convention` behavior.
- **`azurecaf_naming_convention` dead arguments**: Documented that the `prefixes` and `suffixes` *list* arguments are accepted by the schema but **ignored** by `getResult`; only the singular `prefix` and `postfix` are honored. Users who need list semantics should migrate to `azurecaf_name`.
- **TESTING.md**: Removed an orphan `### Test Organization` heading that had no body and was sitting directly under the `## 🗂️ Test Organization` section heading.
- **docs/index.md**: Fixed a broken `<details>` block where leftover generator commentary (a bare `# Resources not in official Azure CAF documentation` heading and a stray `cat resourceDefinition.json | jq …` shell command) was rendering as a level-1 heading and a paragraph between two markdown tables. Replaced with a proper `#### Resources not in official Azure CAF documentation` subheading and a one-sentence intro.

### Security
- **CI/Automation**: Excluded the auto-generated `agentics-maintenance.yml` self-maintenance workflow that `gh aw compile` v0.72.1 emits. The compiler-emitted file inlined `${{ inputs.operation }}` and `${{ inputs.run_url }}` directly into shell `run:` blocks (CWE-94 script injection, SonarCloud rule `actions:S7631`). Even after refactoring those values through environment variables, SonarCloud's analyzer continues to flag any propagation of `${{ inputs.* }}` into a step that has a `run:` block. Since the maintenance workflow is optional, manually-triggered (`workflow_dispatch` / `workflow_call`), and not referenced by any other workflow in this repo, it has been removed from version control. If a future `gh aw compile` re-emits it, it must be deleted again or refactored to use `actions/github-script` (no shell `run:` block) before being committed.
- **CI/Automation**: Added `.checkov.yaml` to suppress Checkov rule `CKV_GHA_7` ("workflow_dispatch inputs MUST be empty") for the auto-generated agentic workflow lock files. The compiler emits an internal `aw_context` input on every `workflow_dispatch` trigger; the lock files carry "DO NOT EDIT" headers, so a global skip with a documented justification is the auto-regen-safe approach. The repo is not a SLSA-tracked build-artifact producer, so the rule does not apply.

## [v1.2.32] - 2026-03-23

### Added
- **Documentation**: Added Terraform Registry downloads badge to README.md
- **Feature**: Added `error_when_exceeding_max_length` attribute to `azurecaf_name` resource and data source
  - When enabled, returns an error if the composed name exceeds the resource type's maximum length
  - Supports both resource and data source usage
  - Includes comprehensive tests and documentation

### Fixed
- **Bug Fix**: Handle `regexp.Compile` errors to prevent nil pointer panics (#379)
  - Properly handle regex compilation errors in resource name validation
  - Prevents crashes when invalid regex patterns are encountered
- **Testing**: Moved test helper functions to `_test.go` file to exclude them from coverage metrics (92.8% → 98.3%)


- **Automation**: Added comprehensive Copilot skills and agents framework for repository automation
  - **18 new skills** across 6 domains: resource lifecycle, CI/CD & testing, release management, community, documentation, and Azure sync
    - `changelog-update` — automated CHANGELOG.md entry creation with semver impact assessment
    - `readme-resource-table` — README.md resource status table sync
    - `regression-test-runner` — full CI test suite execution and reporting
    - `e2e-test-runner` — E2E test execution with structured summaries
    - `coverage-analysis` — test coverage analysis against 95% threshold
    - `test-failure-diagnosis` — automated test failure root cause analysis
    - `pr-compliance-check` — PR checklist validation (generated code, CHANGELOG, README)
    - `resource-diff-report` — resourceDefinition.json diff between versions
    - `issue-to-resource-spec` — parse resource request issues into draft JSON entries
    - `contributor-guide` — step-by-step contribution guidance by type
    - `azure-caf-sync` — CAF slug drift detection from official Microsoft docs
    - `azure-resource-discovery` — discover new azurerm resources not yet supported
    - `naming-rules-drift-check` — detect Azure naming rule changes
    - `resource-completeness-check` — provider coverage vs known azurerm resources
    - `semver-assessment` — semantic version bump determination
    - `release-notes-generator` — GitHub Release notes from CHANGELOG
    - `pre-release-validation` — comprehensive pre-release checks
    - `docs-resource-sync` — keep documentation in sync with resource definitions
    - `example-generator` — generate Terraform example configurations
    - `resource-bulk-import` — batch-research and insert multiple resources
  - **6 new interactive agents** for Copilot Chat workflows:
    - `caf.add-resource` — end-to-end new resource addition (research → JSON → build → test → changelog → readme)
    - `caf.update-resource` — update existing resource with full validation
    - `caf.bulk-add-resources` — add multiple resources in one session
    - `caf.audit-resources` — full audit: completeness, drift, coverage
    - `caf.release-prep` — prepare releases: validate, version, generate notes
    - `caf.diagnose-failure` — diagnose and fix build/test failures
  - **6 new GitHub Actions agentic workflows**:
    - `nightly-regression` — nightly test suite on main, auto-creates issues on failure
    - `pr-review-agent` — automated PR compliance review on open/update
    - `issue-to-pr-agent` — auto-creates PR from issues labeled `new-resource`
    - `contributor-welcome` — welcomes first-time contributors with guidance
    - `weekly-azure-sync` — weekly Azure resource discovery and CAF drift detection
    - `release-validation` — validates releases on tag push
  - Updated `copilot-instructions.md` with complete skill and agent catalog
  - Impact: High — enables comprehensive repository automation across all development workflows

### Fixed

- **Nil pointer panic on invalid regex patterns**: `cleanString()` in `resource_name.go` and `getResult()` in `resource_naming_convention.go` previously discarded `regexp.Compile` errors, causing nil pointer dereference panics at runtime when a regex pattern in `resourceDefinition.json` was invalid. Both functions now handle compilation errors explicitly: `cleanString()` logs a warning and returns the input string unchanged, while `getResult()` returns a descriptive error to the caller.
- **Function App Resources**: Added support for new Azure Function App resource types
  - Added `azurerm_linux_function_app` with slug `fa`
  - Added `azurerm_linux_function_app_slot` with slug `fas`
  - Added `azurerm_windows_function_app` with slug `fa`
  - Added `azurerm_windows_function_app_slot` with slug `fas`
  - These replace the deprecated `azurerm_function_app` and `azurerm_function_app_slot`
  - Maintains consistency with existing function app naming rules (2-60 chars, global scope)
  - Impact: High - Enables support for modern Azure Function App deployment patterns
- **azurerm_managed_redis**: Added support for Azure Managed Redis resource type
  - Slug: `amr` (per Microsoft CAF documentation)
  - Length: 3–63 characters
  - Scope: `resourceGroup`
  - Valid characters: alphanumeric and hyphens; must start and end with alphanumeric; consecutive hyphens are permitted (matches current regex validation behavior)
  - Resource provider namespace: `Microsoft.Cache/redisEnterprise`
  - This resource supersedes `azurerm_redis_cache` (Azure Cache for Redis), which is being retired
  - Impact: Medium - Enables CAF-compliant naming for the new Azure Managed Redis offering

### Fixed
- **Go Version Alignment**: Resolved conflicting Go version declarations in go.mod
  - Changed from conflicting `go 1.23.0` and `toolchain go1.24.4` to unified `go 1.24`
  - Eliminates version mismatch errors during builds
  - Ensures consistent Go toolchain usage across all environments
  - Impact: Medium - Fixes build reliability and development environment consistency
- **Linting Issues**: Fixed non-constant format string errors in logging and error handling
  - Fixed `fmt.Errorf` call in `resource_name.go` to use proper format string
  - Fixed `log.Printf` call in `resource_naming_convention.go` to use proper format string
  - Resolves Go vet warnings and ensures build passes all checks
  - Impact: Low - Improves code quality and eliminates build warnings

## [v1.2.31] - 2025-07-03

### Fixed
- **CI/CD Pipeline**: Release provider in zip archives instead of tarballs for Terraform Registry compatibility

## [v1.2.30] - 2025-07-02

### Fixed
- **CI/CD Pipeline**: Fixed GoReleaser failure due to git tag mismatch and dirty state
  - Removed problematic auto-commit step that was creating commits during release process
  - Fixed generated file timestamp stability to prevent dirty git state in CI
  - Added `fetch-depth: 0` to GitHub Actions checkout for full git history
  - Stabilized `models_generated.go` timestamp format to be environment-independent
  - Resolves GoReleaser errors: "git tag was not made against commit" and "git is in a dirty state"
  - Impact: High - Fixes release automation and ensures reliable tag-based releases
- **GoReleaser Configuration**: Updated GoReleaser configuration to v2 format
  - Added `version: 2` to support GoReleaser v2.x
  - Changed `changelog.skip: true` to `changelog.disable: true`
  - Removed deprecated `archives.format` property to use automatic format selection
  - Fixes release pipeline compatibility with goreleaser-action@v6
- **GitHub Workflow**: Fixed workflow step ordering and improved GPG key import
  - Moved "Set up Go" step before "Install tfproviderlint" to resolve dependency issues
  - Enhanced GPG key import with additional configuration options
  - Added `continue-on-error: true` for GPG import to handle missing secrets gracefully
  - Improved Git signing configuration with proper trust levels
- **README Display**: Fixed GitHub repository homepage README display issue
  - Converted README.md line endings from Windows-style (CRLF) to Unix-style (LF)
  - Renamed .github/README.md to .github/README-workflows.md to prevent GitHub display conflict
  - Resolves issue where GitHub was showing workflows documentation instead of main project README
  - Ensures proper display of comprehensive project documentation on repository homepage
- **Code Generation**: Removed timestamp from generated `models_generated.go` file
  - Eliminated dynamic timestamp that was causing git dirty state during CI/CD
  - Removed `GeneratedTime` field from template data structure
  - Updated template to exclude timestamp comment from generated code
  - Impact: High - Prevents GoReleaser "git is in a dirty state" errors during releases
  - Resolves: CI builds no longer modify tracked files during generation process

### Security
- **CRITICAL**: Fixed security vulnerabilities in Go dependencies:
  - Updated `golang.org/x/net` from v0.23.0 to v0.38.0 to resolve:
    - GO-2025-3595: Cross-site Scripting vulnerability in html package
    - GO-2025-3503: HTTP Proxy bypass using IPv6 Zone IDs
    - GO-2024-3333: Non-linear parsing vulnerability in html package
  - Updated `golang.org/x/crypto` from v0.21.0 to v0.36.0
  - Updated `golang.org/x/sys` from v0.18.0 to v0.31.0
  - Updated `golang.org/x/text` from v0.14.0 to v0.23.0
- Updated Go toolchain from 1.20 to 1.23.0 with Go 1.24.4 for enhanced security
- **SECURITY**: Fixed loose POSIX file permissions in E2E testing framework:
  - Changed directory permissions from 0755 to 0750 (removed world access)
  - Changed executable file permissions from 0755 to 0750 (removed world access)
  - Affected files: `e2e/framework/e2e_test.go`, `e2e/framework/framework.go`

### Added
- **E2E Testing Infrastructure**: Complete end-to-end testing framework
  - Comprehensive test suite covering all provider functionality
  - Import functionality testing with real Terraform state management
  - Data source validation with cross-platform compatibility
  - Naming convention testing across multiple resource types
  - Multi-resource type testing for complex scenarios
- **CI/CD Integration**: Full GitHub Actions integration for automated testing
  - Quick E2E tests on every push (~10-15 seconds)
  - Full E2E tests on pull requests (~25-30 seconds) 
  - Manual workflow dispatch with selective test execution
  - Smart triggering based on file changes
- **Local CI Simulation**: Act integration for local CI environment testing
  - Complete workflow validation before pushing to GitHub
  - Docker-based CI environment simulation
  - Cross-platform testing (macOS M-series compatibility)
  - Comprehensive testing scripts for development workflow
- **Testing Scripts**: Production-ready testing automation
  - `scripts/complete-e2e-validation.sh` - Full validation pipeline
  - `scripts/quick-ci-test.sh` - Quick CI environment validation
  - `scripts/test-ci-with-act.sh` - Interactive CI simulation
  - `scripts/validate-ci-e2e.sh` - Enhanced local + CI testing
- **Documentation**: Complete testing and CI/CD documentation
  - `E2E_IMPLEMENTATION_SUMMARY.md` - Implementation overview
  - `ACT_TESTING_GUIDE.md` - Local CI testing guide
  - `CI_E2E_INTEGRATION.md` - CI/CD integration documentation
  - `e2e/README.md` - E2E testing framework documentation
- **GitHub Copilot Integration**: Enhanced development workflow automation
  - `copilot-setup-steps.yml` - GitHub Actions workflow for Copilot environment setup
  - Automated Go and Terraform environment configuration for Copilot sessions
  - Streamlined development environment preparation with proper versioning
- **MAJOR**: Comprehensive end-to-end (E2E) testing framework for real-world validation
  - Complete E2E test suite covering provider build → Terraform usage → Azure integration
  - Mock Azure RM provider integration for testing without actual Azure API calls
  - Automated provider compilation and local installation testing
  - Azure resource naming compliance validation for all supported resource types
  - Support for all naming conventions (cafclassic, cafrandom, random, passthrough)
  - Edge case testing including length limits, special characters, and error conditions
  - Integration testing with azurerm provider using mock scenarios
  - Test runner CLI with flexible execution options and debugging support
  - Comprehensive documentation and troubleshooting guides
- New Makefile targets for E2E testing:
  - `test_e2e` - Complete E2E test suite
  - `test_e2e_quick` - Fast E2E tests for CI/CD
  - `test_e2e_integration` - AzureRM integration tests
  - `test_e2e_naming` - Naming convention validation
  - `test_e2e_edge_cases` - Edge case scenarios
  - `test_e2e_verbose` - Verbose output for debugging
  - `test_complete_with_e2e` - Complete testing including E2E
- Official Azure Cloud Adoption Framework documentation mapping for 55 resources
- New nested `official` object structure containing Azure CAF documentation attributes
- Comprehensive official resource provider namespace mappings
- GitHub Copilot Agent firewall configuration for improved CI/CD testing
- Enhanced resource validation and testing framework
- Comprehensive CI testing pipeline with resource validation, matrix testing, and coverage analysis
- Advanced Makefile targets for comprehensive testing (`test_ci`, `test_ci_fast`, `test_ci_complete`)
- Shared testing utilities to reduce code duplication (SonarQube compliance)
- Refactored naming convention tests to use centralized test helpers

### Changed
- **BREAKING**: Consolidated `resourceDefinition.json` and `resourceDefinition_out_of_docs.json` into single unified file
- **BREAKING**: Refactored JSON structure to nest official Azure CAF attributes under `official` object
- Updated resource definitions to include proper Azure CAF documentation mapping for key resources:
  - API Management service instance (`apim`) - Microsoft.ApiManagement/service
  - AKS cluster (`aks`) - Microsoft.ContainerService/managedClusters
  - Container apps (`ca`) - Microsoft.App/containerApps
  - Application gateway (`agw`) - Microsoft.ApplicationGateway/applicationGateways
  - Virtual network (`vnet`) - Microsoft.Network/virtualNetworks
  - Storage account (`st`) - Microsoft.Storage/storageAccounts
  - And 49 additional resources with official mappings
- Simplified resource definition structure for non-official resources (only `resource` field in `official` object)
- Enhanced code generation logic to handle nested official attributes
- Updated documentation and contribution guidelines to reflect new structure

### Fixed
- DNS blocking issues with `checkpoint-api.hashicorp.com` during integration tests
- Resource provider namespace accuracy for officially documented Azure resources
- Resource generation and validation processes for unified file structure

### Removed
- `resourceDefinition_out_of_docs.json` file (consolidated into main file)
- Legacy flat structure for official documentation attributes

## Migration Guide

### For Contributors
- Use the new nested `official` object structure when adding or modifying resources
- Resources in official Azure CAF documentation should include `slug`, `resource`, and `resource_provider_namespace` in the `official` object
- Resources not in official documentation should only include the `resource` field in the `official` object

### For Consumers
- The root-level `slug` field remains unchanged for backward compatibility
- New official documentation data is available through the `official` object
- No breaking changes to existing provider functionality

## Statistics
- **Total Resources**: 395 (previously 364 + 31 across two files)
- **Official Azure CAF Mappings**: 55 resources with complete official documentation data
- **Non-Official Resources**: 340 resources with simplified official structure
- **Files Consolidated**: 2 → 1 resource definition file

## [v1.2.29] - 2025-06-16

### Added
- Support for Azure Dev Center resources (`azurerm_dev_center`, `azurerm_dev_center_project`, `azurerm_dev_center_gallery`)
- Support for `azurerm_service_plan`
- Support for `azurerm_servicebus_namespace_disaster_recovery_config`
- Support for Log Analytics Solution, Query Pack, and Monitor Data Collection Rule
- Support for `azurerm_bot_service_azure_bot`
- Support for `azurerm_data_protection_backup_policy_postgresql_flexible_server`
- Added LICENSE file

### Fixed
- Fixed typos across documentation
- Replaced Dockerfile with devcontainer.json configuration

### Security
- Bumped `golang.org/x/net` from 0.20.0 to 0.23.0
- Bumped `google.golang.org/protobuf` from 1.32.0 to 1.33.0

### Added
- Support for Azure OpenAI Deployment resource

## [v1.2.28] - 2024-03-13

### Added
- Support for `azurerm_load_test` resource

### Changed
- Updated documentation to add notice and remove Microsoft references

## [v1.2.27] - 2024-01-16

### Added
- Support for `azurerm_powerbi_embedded`
- Support for `azurerm_search_service`
- Support for `azurerm_monitor_data_collection_endpoint`
- Support for `azurerm_portal_dashboard`
- Support for `azurerm_route_server`

### Security
- Bumped `github.com/cloudflare/circl` from 1.3.3 to 1.3.7
- Bumped `golang.org/x/net` from 0.7.0 to 0.17.0
- Bumped `golang.org/x/crypto` from 0.1.0 to 0.17.0
- Bumped `google.golang.org/grpc` from 1.32.0 to 1.56.3

## [v1.2.26] - 2023-06-23

### Added
- Support for `azurerm_linux_web_app` and `azurerm_windows_web_app`

### Fixed
- Better error handling for empty environment variable values

## [v1.2.25] - 2023-05-03

### Added
- Support for IoT security resources
- Support for IoTHub and DPS shared access policies

## [v1.2.24] - 2023-03-09

### Added
- Support for `azurerm_container_app`
- Support for `azurerm_virtual_machine_portal_name`

### Fixed
- Fixed documentation issues
- Reverted #200

### Security
- Bumped `golang.org/x/crypto` from 0.0.0 to 0.1.0
- Bumped `golang.org/x/net` from 0.0.0 to 0.7.0
- Bumped `golang.org/x/text` from 0.3.5 to 0.3.8

## [v1.2.23] - 2022-11-29

### Added
- Support for FHIR service
- Support for `azurerm_federated_identity_credential`
- Support for log alert, Kubernetes fleet manager, DNS forwarding rule and VNet link
- Support for Application Insights web test

### Fixed
- Allow hyphens in `azurerm_shared_image` resource

## [v1.2.22] - 2022-11-16

### Added
- Support for `azurerm_maintenance_configuration`
- Data source `azurecaf_name` for reading existing naming conventions

### Fixed
- Fixed issues #194 and #204

## [v1.2.21] - 2022-11-01

### Added
- Support for DNS forwarding rulesets and private resolver endpoints
- Support for CDN FrontDoor route and custom domain
- Support for metric alert, NGINX, and DNS resolver resources
- Support for `azurerm_automation_job_schedule`
- Data source `azurecaf_environment_variable` for reading environment variables

### Changed
- Increased length limits for `azurerm_windows_virtual_machine` and `azurerm_windows_virtual_machine_scale_set`

## [v1.2.20] - 2022-09-28

### Added
- Support for Azure Red Hat OpenShift (ARO)
- Support for `azurerm_web_pubsub` and `azurerm_web_pubsub_hub`
- Support for CDN FrontDoor rule and secret

## [v1.2.19] - 2022-08-19

### Added
- Support for `azurerm_static_site`
- Support for CDN FrontDoor firewall and security policies

### Fixed
- Corrected slug for `azurerm_static_site`

## [v1.2.18] - 2022-08-03

### Added
- Support for `azurerm_iothub_certificate`

### Fixed
- Fixed issue #162
- Updated CI/CD pipeline

## [v1.2.17] - 2022-05-05

### Added
- Support for `azurerm_data_protection_backup_vault` and backup policies
- Support for `azurerm_virtual_hub_connection`

## [v1.2.16] - 2022-03-11

### Fixed
- Fixed `azurerm_synapse_sql_pool` naming

## [v1.2.15] - 2022-03-07

### Added
- Support for Synapse and Purview resources

## [v1.2.14] - 2022-03-01

### Added
- Support for `azurerm_log_analytics_storage_insights`

### Fixed
- Fixed issues #146 and #156

## [v1.2.13] - 2022-02-15

### Fixed
- Fixes for APIM naming conventions

## [v1.2.12] - 2022-02-14

### Added
- Support for `azurerm_mysql_flexible_server`
- Support for `azurerm_digital_twins_endpoint`
- Support for `azurerm_aadb2c_directory`

## [v1.2.11] - 2022-01-14

### Added
- Support for Load Balancer resources
- Support for APIM resources
- Support for `azurerm_digital_twins_instance`

## [v1.2.10] - 2021-12-02

### Added
- Support for `azurerm_elastic_cloud_deployment`
- Support for `azurerm_postgresql_flexible_server` resources
- Support for additional resources (#127, #136)

## [v1.2.9] - 2021-11-24

### Added
- Support for Data Factory resources

### Changed
- Updated Go and GoReleaser versions

## [v1.2.8] - 2021-11-17

### Added
- Support for Azure Communication Services
- Support for `azurerm_machine_learning_compute_instance`

## [v1.2.7] - 2021-11-10

### Added
- Support for `azurerm_storage_sync` and `azurerm_storage_sync_group`

### Fixed
- Updated regex for `azurerm_managed_disk`
- Added darwin_arm64 build target

## [v1.2.6] - 2021-08-24

### Added
- Support for `azurerm_web_application_firewall_policy`
- Support for `azurerm_vmware_cluster`
- Support for NetApp resources
- Support for issue #59

### Fixed
- Fixed issues #107 and #120

## [v1.2.5] - 2021-07-02

### Added
- Support for `azurerm_consumption_budget_subscription`
- Support for `azurerm_consumption_budget_resource_group`
- Support for `azurerm_monitor_action_group`

### Fixed
- Fixed tests with Terraform v1.0.1

## [v1.2.4] - 2021-06-18

### Added
- Support for `azurerm_vpn_gateway_connection`
- Support for `azurerm_vpn_site`
- Support for `azurerm_monitor_activity_log_alert`

### Fixed
- Fixed CosmosDB naming (#89)
- Fixed App Configuration by removing `_` support (#58)

## [v1.2.3] - 2021-05-04

### Added
- Support for computer name prefix resource

### Changed
- Upgraded dependencies and actions to latest Go version

## [v1.2.2] - 2021-02-18

### Fixed
- Added missing `models_generated.go`

## [v1.2.1] - 2021-02-18

### Added
- Support for `azurerm_ip_group`

### Changed
- Pinned GoReleaser version

## [v1.2.0] - 2021-02-04

### Added
- Support for Logic App
- Support for Functions App

### Changed
- Migrated provider to SDK v2

## [v1.1.9] - 2020-12-16

### Fixed
- Fixed naming convention for APIM (#74)

### Added
- Added contribution guidelines

## [v1.1.8] - 2020-11-24

### Added
- Added implementation reference documentation

## [v1.1.7] - 2020-11-17

### Fixed
- Fixed shared image documentation bug

## [v1.1.6] - 2020-11-17

### Added
- Added test for hyphen validation in regex

## [v1.1.5] - 2020-11-10

### Added
- Support for `azurerm_monitor_diagnostic_setting`
- Added space to valid characters for resource names

## [v1.1.4] - 2020-10-20

### Fixed
- Corrected name validation

## [v1.1.3] - 2020-10-19

### Fixed
- Fixed App Service Environment (ASE) slug and character limits

## [v1.1.2] - 2020-10-19

### Added
- Support for `azurerm_private_dns_zone_virtual_network_link`
- Support for Managed Identity (MSI)
- Support for Synapse SQL pools

## [v1.1.1] - 2020-09-17

### Fixed
- Corrected regex for name cleaning

## [v1.1.0] - 2020-09-17

### Added
- Support for Azure Synapse resources

### Fixed
- Removed slug name collision

## [v1.0.0] - 2020-09-16

### Added
- Support for `azurerm_recovery_services_vault`
- Additional missing resources
- Initial stable release with comprehensive Azure resource naming support

## [v0.4.3] - 2020-06-18

### Changed
- Added GoReleaser for automated releases
- Moved documentation to `docs/` directory for Terraform Registry compliance
- Added GPG signing for releases

## [v0.4.2] - 2020-05-15

### Added
- Support for AKS (cluster, node pool, DNS prefix)
- Support for Application Gateway, API Management, Application Insights
- Support for App Service, App Service Plan, SQL Server, SQL Database
- Support for Application Service Environment (ASE)
- Support for subnets

### Fixed
- Fixed filter regex patterns
- Fixed passthrough mode character stripping
- Fixed validation regex character limits based on Microsoft docs

## [v0.2.1] - 2020-03-27

### Fixed
- Bug fixes and improvements

## [v0.2] - 2020-03-23

### Added
- Initial release with basic Azure resource naming convention support