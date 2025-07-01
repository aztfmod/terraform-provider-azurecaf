# Testing Guide for azurecaf Terraform Provider

This comprehensive guide covers testing strategies, tools, and best practices for the Azure CAF Terraform Provider. The project maintains high test coverage (>99%) to ensure reliability and correctness.

## 📋 Table of Contents

- [Quick Start](#quick-start)
- [Test Organization](#test-organization)
- [Running Tests](#running-tests)
- [Understanding Test Files](#understanding-test-files)
- [Writing New Tests](#writing-new-tests)
- [Test Coverage](#test-coverage)
- [Testing Best Practices](#testing-best-practices)
- [Continuous Integration](#continuous-integration)

## 🚀 Quick Start

### Prerequisites

- Go 1.19+
- Make (for using Makefile targets)
- Terraform CLI (for integration tests)

### Run All Tests

```bash
# Unit tests only (fast)
make unittest

# All tests including integration
make test_all

# CI-friendly tests (unit + coverage, no integration)
make test_ci
```

## 🗂️ Test Organization

The project includes a comprehensive test suite designed to ensure the proper functioning of the provider. Tests are organized into several categories:

### Test Organization

### 1. **Provider Tests**
   - `provider_test.go` - Basic provider configuration and initialization tests

### 2. **Resource Name Tests**
   - `resource_name_test.go` - Tests for the resource name resource functionality:
     - String cleaning operations
     - Name concatenation
     - Resource slug retrieval
     - Acceptance tests for the resource name resource

### 3. **Naming Convention Tests**
   - `resource_naming_convention_test.go` - Helper functions for testing naming conventions
   - `resource_naming_convention_cafclassic_test.go` - CAF Classic naming convention tests
   - `resource_naming_convention_cafrandom_test.go` - CAF Random naming convention tests
   - `resource_naming_convention_random_test.go` - Random naming convention tests
   - `resource_naming_convention_passthrough_test.go` - Passthrough naming convention tests

### 4. **Data Source Tests**
   - Tests for environment variables and data source functionality
   - Validation of data source behavior vs resource behavior

### 5. **Model Tests**
   - `models_generated_test.go` - Tests for generated resource models:
     - Regex validations
     - String processing
     - Length constraints
     - Resource definition validation

### 6. **Integration Tests**
   - `integration_data_sources_test.go` - Data source integration with Terraform
   - `integration_error_cases_test.go` - Error handling and edge cases
   - `integration_all_resource_types_test.go` - Comprehensive resource type testing
   - `integration_cross_resource_test.go` - Cross-resource interactions
   - `integration_naming_convention_types_test.go` - Naming convention validation

### 7. **Coverage and Edge Case Tests**
   - `complete_coverage_test.go` - Tests to improve code coverage
   - `final_coverage_test.go` - Additional coverage tests
   - `final_edge_cases_test.go` - Edge cases and error paths
   - `remaining_coverage_test.go` - Remaining untested code paths
   - `enhanced_tests_test.go` - Enhanced test cases with structured test data

### 8. **End-to-End (E2E) Tests**
   - `e2e/e2e_test.go` - Comprehensive end-to-end tests that validate the complete workflow:
     - Provider build from source
     - Terraform integration with built provider
     - Azure CAF name validation
     - Mock azurerm provider integration
     - Deployment scenario validation

## 🏃‍♂️ Running Tests

### Available Test Commands

The project provides several testing targets via Makefile:

```bash
# Unit Tests (Fast - No Terraform Required)
make unittest                 # Run all unit tests without coverage
make test_coverage           # Run unit tests with coverage reporting
make test_coverage_html      # Generate HTML coverage report

# Integration Tests (Slower - Requires Terraform)
make test_integration        # Run all integration tests
make test_data_sources       # Run data source integration tests  
make test_error_handling     # Run error handling integration tests

# End-to-End Tests (Comprehensive - Requires Terraform CLI)
make test_e2e               # Run full end-to-end tests
make test_e2e_ci            # Run E2E tests in CI mode

# Comprehensive Testing
make test_all               # Run unit, integration, and e2e tests
make test_ci                # Run CI tests (unit + coverage, no integration)

# Specialized Tests
make test_resource_naming   # Run naming convention tests specifically

# Build and Test
make build                  # Build project and run unit tests
make clean                  # Clean up build artifacts and test results
```

### Running Specific Test Categories

**Unit Tests Only:**
```bash
# Standard unit tests
go test ./azurecaf/...

# With verbose output
go test -v ./azurecaf/...

# Specific test patterns
go test ./azurecaf/... -run="TestResourceName"
go test ./azurecaf/... -run="TestNamingConvention"
```

**Integration Tests:**
```bash
# All integration tests (requires TF_ACC=1)
TF_ACC=1 go test -v ./azurecaf/... -run="TestAcc"

# Specific integration test categories
TF_ACC=1 go test -v ./azurecaf/... -run="TestAccDataSourcesIntegration"
TF_ACC=1 go test -v ./azurecaf/... -run="TestAccErrorHandling"
```

**End-to-End Tests:**
```bash
# Full E2E workflow (requires Terraform CLI)
go test -v ./e2e/...

# E2E tests in CI mode
CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 go test -v ./e2e/...

# Skip E2E tests in fast testing
go test -short ./...  # E2E tests are automatically skipped
```

**Coverage Analysis:**
```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./azurecaf/...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Environment Variables for Testing

| Variable | Purpose | Default |
|----------|---------|---------|
| `TF_ACC` | Enable Terraform acceptance tests | `""` (disabled) |
| `TF_LOG` | Terraform logging level | `""` |
| `TF_LOG_PATH` | Terraform log file path | `""` |
| `CHECKPOINT_DISABLE` | Disable Terraform update checks | `""` (enabled) |
| `TF_IN_AUTOMATION` | Disable interactive prompts | `""` (disabled) |
| `TF_CLI_ARGS_init` | Additional arguments for terraform init | `""` |

### Firewall and Connectivity Issues

The tests may encounter connectivity issues when trying to access external services like `checkpoint-api.hashicorp.com`. This typically manifests as DNS block errors during test execution. To resolve this:

1. **Automatic Resolution**: The Makefile and GitHub Actions have been configured to automatically set the required environment variables to disable these checks.

2. **Manual Resolution**: If running tests manually, use these environment variables:
   ```bash
   export CHECKPOINT_DISABLE=1
   export TF_IN_AUTOMATION=1
   export TF_CLI_ARGS_init="-upgrade=false"
   ```

3. **Integration Tests**: For tests that require Terraform CLI (prefixed with `TestAcc`), ensure all three variables are set to prevent external connectivity requirements.

## 📁 Understanding Test Files

### Test File Naming Convention

| Pattern | Purpose | Example |
|---------|---------|---------|
| `*_test.go` | Standard unit tests | `resource_name_test.go` |
| `integration_*_test.go` | Integration tests with Terraform | `integration_data_sources_test.go` |
| `*_coverage_test.go` | Tests specifically for coverage improvement | `complete_coverage_test.go` |
| `enhanced_*_test.go` | Comprehensive test suites | `enhanced_tests_test.go` |

### Test Structure Examples

**Unit Test Structure:**
```go
func TestResourceName_BasicFunctionality(t *testing.T) {
    tests := []struct {
        name     string
        input    map[string]interface{}
        expected string
        wantErr  bool
    }{
        {
            name: "basic storage account name",
            input: map[string]interface{}{
                "name":          "myapp",
                "resource_type": "azurerm_storage_account",
                "prefixes":      []string{"prod"},
            },
            expected: "stprodmyapp",
            wantErr:  false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

**Integration Test Structure:**
```go
func TestAccDataSourceAzureCAFName_StorageAccount(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccDataSourceAzureCAFName_StorageAccount(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("data.azurecaf_name.test", "result", "stmyapp123"),
                    resource.TestCheckResourceAttr("data.azurecaf_name.test", "id", "stmyapp123"),
                ),
            },
        },
    })
}
```

## ✍️ Writing New Tests

### Guidelines for New Tests

1. **Test Naming**: Use descriptive names that explain what is being tested
2. **Table-Driven Tests**: Use table-driven tests for multiple scenarios
3. **Edge Cases**: Include boundary conditions and error scenarios
4. **Resource Coverage**: Test all supported resource types when adding features
5. **Integration Tests**: Add integration tests for new resources or data sources

### Test Categories to Consider

**When adding a new resource type:**
```go
func TestResourceName_NewResourceType(t *testing.T) {
    tests := []struct {
        name     string
        config   map[string]interface{}
        expected string
        wantErr  bool
    }{
        {
            name: "basic name generation",
            config: map[string]interface{}{
                "name":          "test",
                "resource_type": "azurerm_new_resource",
            },
            expected: "newres-test",
            wantErr:  false,
        },
        {
            name: "with prefixes and suffixes",
            config: map[string]interface{}{
                "name":          "test",
                "resource_type": "azurerm_new_resource",
                "prefixes":      []string{"dev"},
                "suffixes":      []string{"001"},
            },
            expected: "newres-dev-test-001",
            wantErr:  false,
        },
        {
            name: "name too long",
            config: map[string]interface{}{
                "name":          strings.Repeat("a", 100),
                "resource_type": "azurerm_new_resource",
            },
            expected: "",
            wantErr:  true,
        },
        // Test edge cases, validation, etc.
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := generateName(tt.config)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

**Integration test for new resource:**
```go
func TestAccDataSourceAzureCAFName_NewResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: `
                    data "azurecaf_name" "test" {
                        name          = "mytest"
                        resource_type = "azurerm_new_resource"
                        prefixes      = ["dev"]
                        random_length = 3
                    }
                `,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttrSet("data.azurecaf_name.test", "result"),
                    resource.TestCheckResourceAttrSet("data.azurecaf_name.test", "id"),
                    resource.TestMatchResourceAttr("data.azurecaf_name.test", "result", 
                        regexp.MustCompile(`^newres-dev-mytest-[a-z0-9]{3}$`)),
                ),
            },
        },
    })
}
```

### Adding Tests for Coverage

When adding tests specifically for coverage improvement:

1. **Identify Uncovered Code**: Use coverage reports to find untested paths
2. **Focus on Error Paths**: Often error handling code is uncovered
3. **Test Edge Cases**: Boundary conditions and unusual inputs
4. **Mock Dependencies**: Use mocks for external dependencies

Example coverage test:
```go
func TestCoverageImprovement_ErrorHandling(t *testing.T) {
    tests := []struct {
        name        string
        input       interface{}
        expectedErr string
    }{
        {
            name:        "nil input",
            input:       nil,
            expectedErr: "input cannot be nil",
        },
        {
            name:        "invalid resource type",
            input:       map[string]interface{}{"resource_type": "invalid"},
            expectedErr: "unsupported resource type",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := processInput(tt.input)
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.expectedErr)
        })
    }
}
```

## 📊 Test Coverage

### Current Coverage Status

The project aims for high test coverage (>95%) across all components. Current metrics:

- **Overall Coverage**: 99.3%
- **Functions Requiring Additional Tests**: `getNameResult` (94.7%)

### Checking Coverage

**Generate Coverage Report:**
```bash
# Basic coverage
make test_coverage

# HTML coverage report
make test_coverage_html
```

**Detailed Coverage Analysis:**
```bash
# Coverage by function
go test -coverprofile=coverage.out ./azurecaf/...
go tool cover -func=coverage.out

# Coverage by file
go tool cover -html=coverage.out
```

### Coverage Goals and Maintenance

**Target Metrics:**
- Minimum 95% statement coverage
- All public functions covered
- All error paths tested
- All resource types validated

**Improving Coverage:**
1. Run coverage reports regularly
2. Identify uncovered code paths
3. Add specific tests for uncovered areas
4. Focus on error handling and edge cases
5. Review coverage in pull requests

## 🎯 Testing Best Practices

### Unit Testing Best Practices

1. **Test One Thing**: Each test should focus on a single behavior
2. **Descriptive Names**: Test names should clearly describe what they test
3. **Arrange-Act-Assert**: Structure tests with clear setup, execution, and verification
4. **Use Table-Driven Tests**: For testing multiple scenarios with similar logic
5. **Test Error Cases**: Don't just test the happy path

### Integration Testing Best Practices

1. **Real Terraform**: Use actual Terraform configurations
2. **Minimal Configs**: Keep test configurations as simple as possible
3. **Check Outputs**: Verify both success and expected values
4. **Clean State**: Ensure tests don't depend on external state
5. **Resource Validation**: Test that generated names work with actual Azure resources

### Naming Convention Testing

1. **All Resource Types**: Test all supported Azure resource types
2. **Boundary Conditions**: Test minimum and maximum length constraints
3. **Character Validation**: Test allowed and disallowed characters
4. **Regex Patterns**: Verify regex patterns match expected names
5. **Case Sensitivity**: Test both case-sensitive and case-insensitive resources

### Performance Testing

While not currently implemented, consider these for performance-critical changes:

```go
func BenchmarkNameGeneration(b *testing.B) {
    config := map[string]interface{}{
        "name":          "benchmark",
        "resource_type": "azurerm_storage_account",
        "random_length": 5,
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := generateName(config)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 🔄 Continuous Integration

The project uses GitHub Actions for continuous integration. Tests are run on:

- **Pull Requests**: All tests run on every PR
- **Main Branch**: Full test suite including integration tests
- **Multiple Platforms**: Tests run on Linux, Windows, and macOS
- **Multiple Go Versions**: Tests against supported Go versions

### CI Test Strategy

1. **Fast Feedback**: Unit tests run first for quick feedback
2. **Comprehensive Validation**: Integration tests run for thorough validation
3. **Coverage Reporting**: Coverage reports are generated and tracked
4. **Artifact Generation**: Test reports and coverage data are preserved

### Local CI Simulation

To simulate CI locally:

```bash
# Run the same tests as CI
make test_ci

# Or run everything like the full CI pipeline
make clean
make build
make test_all
```

## 🐛 Debugging Test Failures

### Common Test Failure Scenarios

1. **Resource Definition Changes**: Update tests when adding new resource types
2. **Validation Rule Changes**: Update expected patterns when changing validation
3. **Coverage Regression**: Add tests for new code paths
4. **Integration Issues**: Check Terraform version compatibility

### Debugging Techniques

**Verbose Test Output:**
```bash
go test -v ./azurecaf/... -run="TestSpecificFunction"
```

**Test with Detailed Logging:**
```bash
TF_LOG=DEBUG go test -v ./azurecaf/... -run="TestAccIntegration"
```

**Run Single Test:**
```bash
go test -v ./azurecaf/... -run="TestResourceName_SpecificCase"
```

### Test Data Inspection

For debugging failing tests, you can inspect intermediate values:

```go
func TestDebugExample(t *testing.T) {
    config := map[string]interface{}{
        "name":          "debug",
        "resource_type": "azurerm_storage_account",
    }
    
    result, err := generateName(config)
    
    // Debug output
    t.Logf("Config: %+v", config)
    t.Logf("Result: %s", result)
    t.Logf("Error: %v", err)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
}
```

---

🧪 **Happy Testing!** Remember that good tests are an investment in the future maintainability and reliability of the provider.

### Understanding Test Files

When working with the test files, it's important to understand their purpose:

1. **Unit Tests**:
   - Focus on testing individual components in isolation
   - Located in `*_test.go` files that match their corresponding source files
   - Use standard Go testing patterns

2. **Integration Tests**:
   - Use the `integration_*_test.go` naming pattern
   - Focus on interactions between multiple components
   - Often use the Terraform Plugin SDK's `resource.UnitTest` framework
   - Require the `TF_ACC=1` environment variable to run

3. **Coverage Tests**:
   - The `*_coverage_test.go` and `*_edge_cases_test.go` files target specific untested code paths
   - Help achieve high code coverage by testing edge cases and error conditions
   - Often modify resource definitions temporarily to test error handling

### Adding New Tests

When adding new tests:

1. **Unit Tests**: Add to the appropriate `*_test.go` file based on the component being tested.
   - For core resource functionality, add to `resource_name_test.go` or `resource_naming_convention_test.go`
   - For model validation, add to `models_generated_test.go`

2. **Integration Tests**: Add to an appropriate `integration_*_test.go` file based on the feature being tested.
   - For data source testing, use `integration_data_sources_test.go`
   - For error handling, use `integration_error_cases_test.go`
   - For cross-resource interactions, use `integration_cross_resource_test.go`

3. **Coverage Tests**: If testing edge cases or error conditions, add to the relevant coverage test file.
   - `enhanced_tests_test.go` - For structured test cases with multiple scenarios
   - `final_edge_cases_test.go` - For specific edge cases and error paths
   - Consider consolidating coverage tests if you're adding multiple new test cases

### Test Coverage Goals

The project aims for high test coverage (>95%) across all components. Current test coverage metrics:

- Overall coverage: 99.3%
- Functions requiring additional tests: `getNameResult` (94.7%)

To check current coverage, run:

```bash
make test_coverage_html
```

This will generate an HTML report showing coverage details for each file and function.

### Test Organization Recommendations

While the current test organization serves its purpose, future improvements could include:

1. **Consolidation of Coverage Tests**: The various coverage test files (`complete_coverage_test.go`, `final_coverage_test.go`, etc.) could be consolidated into more focused categories.

2. **Consistent Naming Patterns**: Adopting a more consistent naming pattern for test files would make it easier to understand their purpose.

3. **Table-Driven Tests**: Make greater use of table-driven tests for testing multiple scenarios with similar code, as seen in `enhanced_tests_test.go`.

4. **Integration Test Organization**: Continue to keep integration tests in separate files with descriptive names as already done.
