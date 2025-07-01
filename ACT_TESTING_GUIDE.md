# Testing CI Workflows Locally with Act

This guide explains how to test GitHub Actions workflows locally using [act](https://github.com/nektos/act).

## Prerequisites

### Install Act
```bash
# macOS with Homebrew
brew install act

# Or download from GitHub releases
curl -sL https://raw.githubusercontent.com/nektos/act/master/install.sh | bash
```

### Docker Desktop
Act requires Docker to run containers locally:
- Install [Docker Desktop](https://www.docker.com/products/docker-desktop)
- Start Docker Desktop before running act

## Configuration

The repository includes `.actrc` configuration for optimal performance:

```
-P ubuntu-latest=catthehacker/ubuntu:act-latest
--container-architecture linux/amd64
```

This configuration:
- Uses optimized Ubuntu containers for faster startup
- Sets correct architecture for M-series Macs
- Reduces resource usage compared to full GitHub runners

## Available Workflows

List all available workflows:
```bash
act --list
```

Current workflows:
- **E2E Tests** (`e2e.yml`) - Dedicated E2E testing workflow
- **Go** (`go.yml`) - Main CI workflow with E2E integration
- **CodeQL** (`codeql.yml`) - Security analysis
- **Security** (`security.yml`) - Security scanning

## Testing E2E Workflows

### Quick Validation
```bash
# Test workflow structure (dry-run)
act pull_request --job e2e-tests -n

# Quick CI environment test
./scripts/quick-ci-test.sh
```

### Full E2E Testing
```bash
# Test complete E2E workflow
act pull_request --job e2e-tests

# Interactive CI testing
./scripts/test-ci-with-act.sh
```

### Specific Test Categories
```bash
# Test with manual dispatch (simulate manual trigger)
act workflow_dispatch -W .github/workflows/e2e.yml

# Test specific event types
act pull_request    # Simulate pull request
act push            # Simulate push to main
```

## Act Commands Reference

### Basic Commands
```bash
# List all workflows and jobs
act --list

# Run specific workflow
act -W .github/workflows/e2e.yml

# Run specific job
act pull_request --job e2e-tests

# Dry-run (validate without execution)
act pull_request -n

# Verbose output
act pull_request -v
```

### Environment Variables
```bash
# Set environment variables
act pull_request --env CHECKPOINT_DISABLE=1 --env TF_IN_AUTOMATION=1

# Use .env file
act pull_request --env-file .env
```

### Container Options
```bash
# Use specific platform
act --container-architecture linux/amd64

# Use custom image
act -P ubuntu-latest=catthehacker/ubuntu:act-latest

# Bind mount directory
act --bind
```

## Workflow Testing Strategy

### 1. Development Testing
```bash
# Quick validation during development
act pull_request --job e2e-tests -n

# Test specific changes
act pull_request --job e2e-tests
```

### 2. Pre-commit Testing
```bash
# Full workflow validation before commit
./scripts/quick-ci-test.sh

# Comprehensive testing before PR
./scripts/test-ci-with-act.sh
```

### 3. CI Debugging
```bash
# Debug failed CI jobs locally
act pull_request --job e2e-tests -v

# Test environment setup
act pull_request --job e2e-tests --env DEBUG=1
```

## Performance Optimization

### Container Images
- `act-latest` - Lightweight, faster startup
- `runner-latest` - Full GitHub runner, slower but more complete
- Custom images for specific needs

### Resource Management
```bash
# Limit resource usage
act --container-cap-add SYS_ADMIN --container-cap-drop ALL

# Use specific container options
act --container-options "-m 2g --cpus 2"
```

### Selective Testing
```bash
# Test only changed workflows
act --detect-event

# Skip time-consuming jobs
act --skip job-name
```

## Troubleshooting

### Common Issues

**Docker not found:**
```bash
# Check Docker status
docker info

# Start Docker Desktop if needed
```

**Container architecture errors on M-series Macs:**
```bash
# Use explicit architecture
act --container-architecture linux/amd64
```

**Permission errors:**
```bash
# Check Docker permissions
docker run hello-world

# Fix Docker permissions if needed
sudo chmod 666 /var/run/docker.sock
```

**Network issues:**
```bash
# Use host network
act --container-options "--network host"
```

### Debugging Tips

1. **Start with dry-run:** Always use `-n` first to validate
2. **Use verbose output:** Add `-v` for detailed logging
3. **Check logs:** Review container logs for failures
4. **Test incrementally:** Test individual jobs before full workflows

## Integration with Development Workflow

### Git Hooks
Add to `.git/hooks/pre-commit`:
```bash
#!/bin/bash
# Validate CI before commit
./scripts/quick-ci-test.sh
```

### VS Code Integration
Add to `.vscode/tasks.json`:
```json
{
    "label": "Test CI with Act",
    "type": "shell",
    "command": "./scripts/quick-ci-test.sh",
    "group": "test"
}
```

### Make Integration
Add to `Makefile`:
```makefile
test_ci_local:
	@echo "Testing CI locally with act..."
	act pull_request --job e2e-tests -n
```

## Benefits of Local CI Testing

✅ **Fast Feedback** - Test CI changes without pushing to GitHub  
✅ **Resource Efficient** - No GitHub Actions minutes consumed  
✅ **Debugging** - Debug CI issues in local environment  
✅ **Validation** - Ensure workflows work before committing  
✅ **Development Speed** - Iterate quickly on CI improvements  

## Limitations

⚠️ **Container Differences** - Local containers may differ from GitHub runners  
⚠️ **Secret Access** - Some secrets/services not available locally  
⚠️ **Performance** - Local resources may be different from CI  
⚠️ **Platform Differences** - Some platform-specific issues won't appear locally  

## Best Practices

1. **Always dry-run first** to validate structure
2. **Test incrementally** - individual jobs before full workflows
3. **Use lightweight containers** for faster iteration
4. **Combine with local testing** - act complements but doesn't replace local tests
5. **Document CI changes** - test and document workflow modifications
