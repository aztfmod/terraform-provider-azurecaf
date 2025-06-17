## Testing Guide

The project includes a comprehensive test suite designed to ensure the proper functioning of the provider. Tests are organized into several categories:

### Test Organization

1. **Provider Tests**:
   - `provider_test.go` - Basic provider configuration and initialization tests

2. **Resource Name Tests**:
   - `resource_name_test.go` - Tests for the resource name resource functionality, including:
     - String cleaning operations
     - Name concatenation
     - Resource slug retrieval
     - Acceptance tests for the resource name resource

3. **Naming Convention Tests**:
   - `resource_naming_convention_test.go` - Helper functions for testing naming conventions
   - `resource_naming_convention_cafclassic_test.go` - CAF Classic naming convention tests
   - `resource_naming_convention_cafrandom_test.go` - CAF Random naming convention tests
   - `resource_naming_convention_random_test.go` - Random naming convention tests
   - `resource_naming_convention_passthrough_test.go` - Passthrough naming convention tests

4. **Data Source Tests**:
   - Tests for environment variables and data source functionality

5. **Model Tests**:
   - `models_generated_test.go` - Tests for generated resource models, including:
     - Regex validations
     - String processing
     - Length constraints

6. **Integration Tests**:
   - `integration_data_sources_test.go` - Tests for data source integration
   - `integration_error_cases_test.go` - Error handling tests
   - `integration_all_resource_types_test.go` - Tests for all resource types
   - `integration_cross_resource_test.go` - Tests for cross-resource interactions
   - `integration_naming_convention_types_test.go` - Tests for naming conventions

7. **Coverage and Edge Case Tests**:
   - `complete_coverage_test.go` - Tests to improve code coverage
   - `final_coverage_test.go` - Additional coverage tests
   - `final_edge_cases_test.go` - Tests for edge cases and error paths
   - `remaining_coverage_test.go` - Tests for remaining untested code paths
   - `enhanced_tests_test.go` - Enhanced test cases with structured test data

### Running Tests

The project provides several testing targets in the Makefile:

```bash
# Run all unit tests
make unittest

# Run integration tests (requires TF_ACC=1)
make test_integration

# Run data source integration tests
make test_data_sources

# Run error handling tests
make test_error_handling

# Run naming convention tests
make test_resource_naming

# Run tests with coverage reporting
make test_coverage

# Generate HTML coverage report
make test_coverage_html

# Run all tests
make test_all

# Run CI tests (unit tests with coverage, no integration tests)
make test_ci
```

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
