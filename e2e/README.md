# End-to-End (E2E) Tests for terraform-provider-azurecaf

This directory contains end-to-end tests for the Azure Cloud Adoption Framework (CAF) Terraform provider.

## Overview

The E2E tests validate the complete functionality of the terraform-provider-azurecaf by:

1. **Building the provider from source**
2. **Setting up local provider overrides** using Terraform's `dev_overrides` feature
3. **Running actual Terraform commands** (plan, init, import, etc.) against real terraform configurations
4. **Validating the output** to ensure the provider works correctly

## Test Structure

### `e2e_test.go`
Contains the basic E2E test that validates core functionality:
- Provider builds successfully
- Terraform can load and use the provider
- Basic resource creation works

### `e2e_comprehensive_test.go`
Contains comprehensive test scenarios:
- **TestE2EDataSource**: Tests the `azurecaf_name` data source
- **TestE2ENamingConventions**: Tests different naming configurations (passthrough, random, etc.)
- **TestE2EMultipleResourceTypes**: Tests multiple Azure resource types in one configuration
- **TestE2EImportFunctionality**: Tests terraform import functionality for existing resources

## Running Tests

### Via Makefile (Recommended)

```bash
# Run all E2E tests
make test_e2e

# Run quick E2E tests (basic only)
make test_e2e_quick

# Run specific test categories
make test_e2e_data_source
make test_e2e_naming
make test_e2e_multiple_types
make test_e2e_import        # NEW: Import functionality tests
```

### Direct Go Commands

```bash
# Run all tests
cd e2e && go test -v

# Run specific test
cd e2e && go test -v -run TestE2EImportFunctionality
```

## Test Scenarios

### Basic Test
- Simple storage account name generation
- Validates core provider functionality

### Data Source Test
- Tests the `azurecaf_name` data source
- Validates data source reads complete successfully

### Naming Conventions Test
- Tests `passthrough` mode (exact name preservation)
- Tests random character generation
- Tests different resource configurations

### Multiple Resource Types Test
- Tests storage account, key vault, and virtual machine naming
- Tests prefixes, suffixes, and random length configurations
- Validates multiple resources in single configuration

### Import Functionality Test ✨ **NEW**
- Tests `terraform import` functionality for existing resources
- Validates import format: `azurerm_storage_account:stmyexistingapp`
- Verifies imported resource appears correctly in terraform state
- Confirms import process completes successfully

## Import Test Details

The import test demonstrates the ability to import existing Azure resource names into Terraform state:

```bash
# Example import command tested:
terraform import azurecaf_name.imported_storage azurerm_storage_account:stmyexistingapp
```

**What it validates:**
✅ Import command succeeds without errors  
✅ Resource is successfully added to terraform state  
✅ Imported attributes are correctly populated  
✅ Post-import terraform operations work properly  

## Expected Output

When tests run successfully, you'll see output like:

```
✅ E2E test data_source passed!
--- PASS: TestE2EDataSource (4.40s)
✅ E2E test naming_conventions passed!
--- PASS: TestE2ENamingConventions (3.43s)
✅ E2E test multiple_types passed!
--- PASS: TestE2EMultipleResourceTypes (3.34s)
✅ E2E import test import_functionality passed!
--- PASS: TestE2EImportFunctionality (8.98s)
✅ E2E test passed!
--- PASS: TestE2EBasic (7.76s)

All 5 tests passing in ~28 seconds
```

## Key Features

✅ **Real Provider Testing**: Uses actual built provider binary  
✅ **Isolated Test Environment**: Each test runs in isolated temporary directory  
✅ **Comprehensive Coverage**: Tests resources, data sources, and import functionality  
✅ **Dev Overrides**: Uses Terraform's official development override mechanism  
✅ **Import Testing**: Validates terraform import functionality  
✅ **Cleanup**: Automatic cleanup of temporary directories  
✅ **Fast Execution**: Simple and efficient test execution  

## Integration with CI/CD

The E2E tests are fully integrated into GitHub Actions CI/CD pipelines:

### Main CI Workflow (`.github/workflows/go.yml`)
- **Quick E2E Tests**: Run on every push and PR for fast feedback
- **Full E2E Tests**: Run on pull requests for comprehensive validation
- **Automatic Setup**: Terraform is automatically installed and configured

### Dedicated E2E Workflow (`.github/workflows/e2e.yml`)
- **Comprehensive Testing**: Full suite of E2E tests
- **Manual Triggers**: Can be run manually with different test types
- **Conditional Execution**: Smart execution based on file changes
- **Test Categories**: 
  - `quick` - Basic functionality only
  - `all` - Complete test suite  
  - `import_only` - Import functionality only
  - `naming_only` - Naming convention tests only

### CI Integration Features

✅ **Automatic Terraform Setup**: CI automatically installs Terraform  
✅ **Provider Building**: Builds provider from source in CI  
✅ **Fast Feedback**: Quick tests run on every commit  
✅ **Comprehensive Validation**: Full tests on pull requests  
✅ **Manual Testing**: On-demand test execution  
✅ **Smart Triggers**: Only runs when relevant files change  

### Running E2E Tests in CI

```bash
# CI automatically runs these targets:
make test_e2e_quick    # On every push/PR
make test_e2e          # On pull requests

# Available CI make targets:
make test_ci_with_e2e  # Run CI tests + quick E2E tests
make test_complete_with_e2e  # Run complete test suite + E2E tests
```

### Local Development vs CI

| Environment | Quick Tests | Full Tests | Manual Triggers |
|-------------|------------|------------|-----------------|
| **Local Dev** | `make test_e2e_quick` | `make test_e2e` | All targets available |
| **CI Push** | ✅ Auto | ❌ | ❌ |
| **CI PR** | ✅ Auto | ✅ Auto | ❌ |
| **CI Manual** | ✅ Auto | Configurable | ✅ All options |

### CI Environment Setup

The CI automatically configures:
- Go environment from `go.mod`
- Terraform latest 1.x version
- Provider dev_overrides for testing
- All required environment variables