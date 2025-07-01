# CI/CD Integration for E2E Tests

This document explains how End-to-End (E2E) tests are integrated into the CI/CD pipeline for terraform-provider-azurecaf.

## Overview

The E2E tests are integrated into GitHub Actions workflows to ensure comprehensive validation of the provider functionality in automated environments.

## CI Workflows

### 1. Main CI Workflow (`go.yml`)

**Triggers:**
- Push to `main` branch
- Pull requests to `main` branch  
- Tagged releases

**E2E Test Integration:**
```yaml
- name: E2E Tests (Quick)
  run: make test_e2e_quick
  
- name: E2E Tests (Full) - Pull Requests Only  
  if: github.event_name == 'pull_request'
  run: make test_e2e
```

**Strategy:**
- **Fast Feedback**: Quick E2E tests run on every push
- **Comprehensive Validation**: Full E2E tests run on pull requests
- **Release Safety**: Quick tests validate release builds

### 2. Dedicated E2E Workflow (`e2e.yml`)

**Triggers:**
- Pull requests (when E2E-related files change)
- Manual workflow dispatch with options

**Features:**
- **Selective Testing**: Choose specific test categories
- **Smart Triggers**: Only runs when relevant files change
- **Manual Control**: On-demand testing with different configurations

**Test Categories:**
- `quick` - Basic functionality validation
- `all` - Complete test suite
- `import_only` - Import functionality only  
- `naming_only` - Naming convention tests only

## Test Execution Strategy

### Development Flow
```
1. Local Development
   ├── make test_e2e_quick (fast feedback)
   └── make test_e2e (before PR)

2. Pull Request
   ├── Quick E2E tests (main workflow)
   ├── Full E2E tests (main workflow) 
   └── Selective E2E tests (e2e workflow)

3. Main Branch Push
   └── Quick E2E tests (main workflow)

4. Release
   └── Quick E2E tests (main workflow)
```

### Performance Optimization

| Test Type | Duration | When Run | Purpose |
|-----------|----------|----------|---------|
| **Quick E2E** | ~10-15s | Every push/PR | Fast validation |
| **Full E2E** | ~25-30s | PR only | Comprehensive testing |
| **Selective** | Variable | Manual | Targeted testing |

## Environment Setup

### Automatic Dependencies
The CI automatically installs and configures:

```yaml
- name: Setup Terraform
  uses: hashicorp/setup-terraform@v3
  with:
    terraform_version: "~> 1.0"

- name: Set up Go  
  uses: actions/setup-go@v5
  with:
    go-version-file: './go.mod'
```

### Environment Variables
```yaml
env:
  CHECKPOINT_DISABLE: 1           # Disable Terraform checkpoint
  TF_IN_AUTOMATION: 1            # Terraform automation mode
  TF_CLI_ARGS_init: "-upgrade=false"  # Prevent upgrade prompts
```

## Make Targets for CI

### Primary Targets
```bash
make test_e2e_quick        # Fast validation (~10s)
make test_e2e             # Complete suite (~25s)
```

### Specific Test Categories  
```bash
make test_e2e_data_source     # Data source tests
make test_e2e_naming          # Naming convention tests  
make test_e2e_multiple_types  # Multiple resource types
make test_e2e_import          # Import functionality
```

### Combined Targets
```bash
make test_ci_with_e2e         # CI tests + quick E2E
make test_complete_with_e2e   # All tests + full E2E
```

## Failure Handling

### Test Failure Scenarios
1. **Provider Build Failure**: CI fails early, no E2E tests run
2. **Terraform Setup Failure**: E2E tests are skipped  
3. **Individual Test Failure**: Detailed output provided
4. **Timeout**: Tests fail after reasonable time limit

### Debugging CI Failures
```bash
# Local reproduction of CI environment
CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 make test_e2e_quick

# Debug specific test
cd e2e && go test -v -run TestE2EBasic

# Check provider build
make build
```

## Best Practices

### For Developers
✅ **Run Quick Tests Locally**: Before pushing changes  
✅ **Full Tests Before PR**: Ensure comprehensive validation  
✅ **Check CI Logs**: Review detailed test output on failures  
✅ **Manual E2E Runs**: Use dedicated workflow for specific testing  

### For Maintainers  
✅ **Monitor CI Performance**: Keep test execution times reasonable  
✅ **Update Dependencies**: Keep Terraform version current  
✅ **Review Test Output**: Ensure comprehensive coverage  
✅ **Optimize Selectively**: Balance speed vs coverage  

## Integration Benefits

🚀 **Fast Feedback**: Developers get quick validation  
🔍 **Comprehensive Testing**: Full validation on important changes  
⚡ **Selective Testing**: Targeted testing when needed  
🛡️ **Release Safety**: Validated provider builds  
📊 **Clear Results**: Detailed test output and summaries  

## Future Enhancements

Potential improvements for CI integration:
- **Matrix Testing**: Multiple Terraform versions
- **Parallel Execution**: Split tests across multiple runners  
- **Artifact Upload**: Test reports and logs
- **Performance Tracking**: Test execution time monitoring
- **Notification Integration**: Slack/Teams alerts for failures
