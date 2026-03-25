# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **CI Pipeline**: Replaced sequential `go.yml` with parallelized `ci.yml` + dedicated `release.yml`
  - 7 parallel test jobs via matrix strategies (unit, coverage, resource validation, 4 integration suites, 5 E2E suites)
  - Extracted shared setup into composite action (`.github/actions/setup-go-env`) to eliminate ~100 lines of YAML duplication
  - Build job saves Go module cache; test jobs restore-only — eliminates cache write contention warnings
  - Separated release workflow (`release.yml`) so CI jobs run with minimal `contents: read` permissions
  - Dynamic test summary — adding/removing jobs no longer requires editing the summary script
  - Pinned to latest GitHub Action versions: `actions/checkout@v6`, `actions/setup-go@v6`, `hashicorp/setup-terraform@v4`, `crazy-max/ghaction-import-gpg@v7`, `goreleaser/goreleaser-action@v7`
  - E2E full suite now included in test-summary gate (was previously excluded)
  - Removed redundant `test_ci` job that duplicated work already run by dedicated parallel jobs
  - Impact: Medium — CI/CD infrastructure only, no provider behavior changes

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