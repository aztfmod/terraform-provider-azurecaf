# Complete Resource Testing Guide

This guide explains how to test all 395 Azure resource types supported by the terraform-provider-azurecaf.

## Overview

The provider supports **395 different Azure resource types** and provides comprehensive testing capabilities to ensure all resource types work correctly with the naming conventions.

## Current Test Coverage

- **98.6%** code coverage on the core provider functionality
- **395** Azure resource types defined
- **20 batches** of systematic testing (20 resources per batch)
- **Multiple test configurations** per resource type

## Testing Commands

### Quick Start
```bash
# Run the complete test suite
./test_all_resources.sh
```

### Individual Test Categories

```bash
# Test all resource types (comprehensive - takes time)
make test_all_resources

# Analyze which resources are tested successfully
make test_resource_coverage

# Validate all resource definitions are complete
make test_resource_definitions

# Test resources organized by category (Storage, Compute, etc.)
make test_resource_matrix

# Run standard test suite with coverage
make test_coverage

# Complete test suite (everything)
make test_complete
```

### Specific Batch Testing

The current batch system tests 20 resources at a time. To test a specific batch:

```bash
# Edit azurecaf/integration_all_resource_types_test.go
# Change: batchToRun = 1  // to the batch number you want (1-20)
make test_integration
```

## Test Files and Their Purpose

| File | Purpose |
|------|---------|
| `integration_all_resource_types_test.go` | Original batch testing (batch 1 only) |
| `integration_all_resources_complete_test.go` | **NEW**: Complete testing of all 395 resources |
| `resource_coverage_analysis_test.go` | **NEW**: Coverage analysis and validation |
| `resource_matrix_test.go` | **NEW**: Category-based testing and constraints |
| `models_generated_test.go` | Validation of resource definition structures |
| `complete_coverage_test.go` | Edge cases and specific scenarios |

## What Gets Tested

For each resource type, the tests verify:

1. **Basic Naming**: Simple name generation with prefixes/suffixes
2. **Complex Configuration**: Multiple prefixes, suffixes, separators
3. **Edge Cases**: Maximum/minimum lengths, special characters
4. **Validation**: Results match regex patterns and length constraints
5. **Data Sources**: Both resource and data source variants work

### Example Test Scenarios Per Resource

```hcl
# Configuration 1: Basic
resource "azurecaf_name" "example" {
  name          = "testname"
  resource_type = "azurerm_storage_account"
  prefixes      = ["dev"]
  suffixes      = ["001"]
  random_length = 5
  clean_input   = true
}

# Configuration 2: Complex
resource "azurecaf_name" "example" {
  name          = "testname"
  resource_type = "azurerm_storage_account"
  prefixes      = ["prod", "web"]
  suffixes      = ["001", "east"]
  separator     = "-"
  random_length = 3
  clean_input   = true
  use_slug      = true
}

# Configuration 3: Minimal
resource "azurecaf_name" "example" {
  name          = "testname"
  resource_type = "azurerm_storage_account"
  prefixes      = ["test"]
  clean_input   = true
}
```

## Resource Categories Tested

Resources are automatically categorized and tested by type:

- **Storage**: storage accounts, disks, backups (40+ resources)
- **Compute**: VMs, containers, AKS, batch (60+ resources)
- **Networking**: VNets, gateways, DNS, load balancers (50+ resources)
- **Database**: SQL, MySQL, PostgreSQL, Cosmos DB (30+ resources)
- **Security**: Key Vault, managed identities, roles (20+ resources)
- **Monitoring**: Application Insights, Log Analytics (15+ resources)
- **Web**: App Services, Functions, Logic Apps (25+ resources)
- **AI/ML**: Cognitive Services, Machine Learning (10+ resources)
- **Integration**: Service Bus, Event Hub (15+ resources)
- **Other**: Remaining Azure services (130+ resources)

## Understanding Test Results

### Success Indicators
- ✅ Resource generates valid names
- ✅ Names meet length constraints (min/max)
- ✅ Names match validation regex
- ✅ Both resource and data source work
- ✅ Multiple configurations work

### What to Look For
```bash
# Running tests shows:
✓ Config 1 for azurerm_storage_account: dev-st-testname-xj8kl-001
✓ Config 2 for azurerm_storage_account: prod-web-st-testname-x8k-001-east  
✓ Config 3 for azurerm_storage_account: test-st-testname
```

### Coverage Reports
Tests generate several reports:
- `resource_coverage_report.json` - Detailed success/failure analysis
- `coverage.html` - Code coverage visualization
- Terminal output with category summaries

## Troubleshooting

### Common Issues

1. **Resource Definition Missing Fields**
   ```bash
   # Run this to find incomplete definitions:
   make test_resource_definitions
   ```

2. **Regex Validation Failures**
   ```bash
   # Check regex patterns:
   go test ./azurecaf/... -run="TestCompileRegexValidation"
   ```

3. **Length Constraint Violations**
   ```bash
   # Validate length constraints:
   go test ./azurecaf/... -run="TestRegexValidation.*Length"
   ```

### Running Subset of Tests

To test only specific resource types, edit the test files or use:

```bash
# Test only storage resources
go test ./azurecaf/... -run="TestResourceMatrix.*Storage"

# Test specific resource type
go test ./azurecaf/... -run=".*storage_account"
```

## Continuous Integration

For CI/CD pipelines, use:

```bash
# Fast CI tests (unit + coverage, no comprehensive integration)
make test_ci

# Complete CI tests (includes resource validation)
make test_complete
```

## Performance Considerations

- **Complete test suite**: ~30+ minutes (395 resources × 3 configs each)
- **Single batch**: ~2-3 minutes (20 resources)
- **Coverage analysis**: ~5-10 minutes
- **Resource validation**: ~30 seconds

## Adding New Resources

When adding new resource types:

1. Add to `resourceDefinition.json`
2. Run `go generate` to update `models_generated.go`
3. Run `make test_resource_definitions` to validate
4. Run `make test_resource_coverage` to ensure it works
5. New resource automatically included in comprehensive tests

## Summary

This testing framework ensures that **all 395 Azure resource types** are thoroughly tested with multiple configurations, providing confidence that the provider works correctly for any Azure resource naming scenario.

The tests are organized in a hierarchical way:
- **Batch level**: Groups of 20 resources
- **Resource level**: Individual resource types  
- **Configuration level**: Multiple naming scenarios per resource
- **Category level**: Logical groupings (Storage, Compute, etc.)

This comprehensive approach guarantees that users can rely on the provider for any Azure resource naming needs.
