# End-to-End (E2E) Tests for terraform-provider-azurecaf

This directory contains comprehensive end-to-end tests that validate the complete workflow of the terraform-provider-azurecaf from build to deployment.

## Test Coverage

The E2E tests cover the following scenarios:

### 1. Provider Build and Validation
- ✅ **Build Provider**: Compiles the terraform-provider-azurecaf from source
- ✅ **Binary Validation**: Verifies the provider binary is functional and executable

### 2. Terraform Integration
- ✅ **Local Provider Setup**: Configures Terraform to use the locally built provider
- ✅ **Terraform Commands**: Tests `terraform init`, `plan`, and `apply` with the local provider

### 3. Azure CAF Name Generation
- ✅ **CAF Compliance**: Validates that generated names follow Azure CAF standards
- ✅ **Resource Constraints**: Tests various Azure resource types and their naming constraints
- ✅ **Edge Cases**: Validates special character handling and input cleaning

### 4. Azure Provider Integration (Mock)
- ✅ **azurerm Compatibility**: Tests that generated names work with azurerm provider patterns
- ✅ **Resource Mapping**: Validates name generation for real Azure resource scenarios

### 5. Deployment Scenarios
- ✅ **Multiple Resource Types**: Tests generating names for multiple related resources
- ✅ **Data Source Validation**: Tests the data source functionality
- ✅ **Complex Configurations**: Validates advanced naming scenarios

## Running the Tests

### Prerequisites
- Go 1.21 or later
- Terraform CLI installed
- Internet connection (for Terraform provider downloads)

### Running All E2E Tests
```bash
# From the e2e directory
go test -v

# From the project root
go test -v ./e2e/...
```

### Running with Build Context
```bash
# From the project root - ensures provider is built first
make build && go test -v ./e2e/...
```

### Running in CI/CD
The E2E tests are designed to run in CI/CD environments:
```bash
# CI-friendly execution
CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 go test -v ./e2e/...
```

### Skipping E2E Tests
E2E tests can be skipped in fast feedback scenarios:
```bash
# Skip E2E tests in short mode
go test -short ./...
```

## Test Architecture

### Test Flow
1. **Setup**: Creates isolated test workspace
2. **Build**: Compiles provider from source
3. **Configure**: Sets up Terraform with local provider
4. **Execute**: Runs Terraform commands to test functionality
5. **Validate**: Checks outputs against Azure CAF standards
6. **Cleanup**: Removes test workspace

### Key Components

#### `TestE2E_ComprehensiveWorkflow`
Main test function that orchestrates the complete E2E workflow with subtests for each major component.

#### Helper Functions
- `getProjectRoot()`: Locates the project root directory
- `setupTestWorkspace()`: Creates isolated test environment
- `runTerraformCommand()`: Executes Terraform commands safely
- `validateCAFCompliantNames()`: Validates naming compliance

#### Test Data
- **Terraform Configurations**: Multiple `.tf` files testing different scenarios
- **Mock Integrations**: Simulated azurerm provider interactions
- **Validation Rules**: Azure resource naming constraint checks

## Test Scenarios

### Basic Name Generation
Tests fundamental name generation with prefixes, suffixes, and random components.

### Resource-Specific Constraints
- **Storage Accounts**: 3-24 chars, lowercase alphanumeric only
- **Key Vaults**: 3-24 chars, start with letter, alphanumeric + hyphens
- **Resource Groups**: Up to 90 chars, specific character set

### Advanced Scenarios
- **Multi-Resource**: Single configuration generating multiple resource names
- **Data Sources**: Using data sources for name validation
- **Edge Cases**: Special character handling and input cleaning

## Integration with CI/CD

### Makefile Integration
Add to the main Makefile:
```makefile
test_e2e: ## Run end-to-end tests
	go test -v ./e2e/...

test_e2e_ci: ## Run E2E tests in CI mode
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 go test -v ./e2e/...
```

### GitHub Actions Integration
The tests work with the existing `.github/workflows/go.yml` by adding:
```yaml
- name: Test E2E
  run: make test_e2e_ci
  env:
    CHECKPOINT_DISABLE: 1
    TF_IN_AUTOMATION: 1
```

## Troubleshooting

### Common Issues

#### Provider Binary Not Found
Ensure the provider is built before running tests:
```bash
go build -o terraform-provider-azurecaf
```

#### Terraform Not Found
Install Terraform CLI or ensure it's in PATH:
```bash
which terraform
```

#### Network Issues
Set environment variables for offline/restricted environments:
```bash
export CHECKPOINT_DISABLE=1
export TF_IN_AUTOMATION=1
```

### Debug Mode
Run tests with verbose output for debugging:
```bash
go test -v -run TestE2E_ComprehensiveWorkflow ./e2e/...
```

## Extending the Tests

### Adding New Test Scenarios
1. Create new test configuration in the appropriate test function
2. Add validation logic for the new scenario
3. Update documentation

### Adding New Resource Types
1. Add resource configuration to test scenarios
2. Add specific validation rules for the resource type
3. Test with actual Azure naming constraints

### Mock Provider Integration
To test with other providers:
1. Create mock provider configurations
2. Add provider-specific validation logic
3. Test integration scenarios

## Benefits

### Development Confidence
- Validates complete provider functionality
- Catches integration issues early
- Ensures CAF compliance

### CI/CD Reliability
- Automated validation in pipelines
- Clear pass/fail feedback
- Isolated test environments

### Documentation
- Living examples of provider usage
- Validation of documented features
- Integration patterns