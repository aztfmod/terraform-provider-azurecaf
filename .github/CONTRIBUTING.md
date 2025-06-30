# Contributing to the Azure CAF Terraform Provider

👍🎉 First off, thanks for taking the time to contribute! 🎉👍

We're excited to have you help improve the Azure Cloud Adoption Framework (CAF) Terraform Provider. This document provides guidelines and information to help you contribute effectively.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [What Should I Know Before I Get Started?](#what-should-i-know-before-i-get-started)
- [How Can I Contribute?](#how-can-i-contribute)
- [Adding a New Resource](#adding-a-new-resource)
- [Development Setup](#development-setup)
- [Testing Guidelines](#testing-guidelines)
- [Submitting Changes](#submitting-changes)
- [Documentation Standards](#documentation-standards)

## Code of Conduct

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## What Should I Know Before I Get Started?

### Project Overview

The Azure CAF Terraform Provider is designed to support the Microsoft Cloud Adoption Framework for Azure by providing:

- **Naming Convention Enforcement**: Ensures Azure resource names follow CAF guidelines
- **Resource Validation**: Validates names against Azure-specific constraints
- **Flexibility**: Supports multiple naming patterns and customization options
- **Comprehensive Coverage**: Supports 300+ Azure resource types

### Architecture

The provider consists of several key components:

- **Resource Definitions** (`resourceDefinition.json`): JSON file containing all Azure resource naming rules
- **Code Generation** (`gen.go`): Generates Go code from resource definitions
- **Provider Core** (`azurecaf/`): Main provider logic and resources
- **Templates** (`templates/`): Go templates for code generation

## How Can I Contribute?

### 🐛 Reporting Bugs

When reporting bugs, please include:

- **Clear Description**: What you expected vs. what happened
- **Reproduction Steps**: Minimal Terraform configuration to reproduce the issue
- **Environment Details**: Terraform version, provider version, OS
- **Error Messages**: Full error output if applicable

### 💡 Suggesting Enhancements

Enhancement suggestions are welcome! Please:

- **Search Existing Issues**: Check if your idea has already been suggested
- **Provide Clear Use Case**: Explain why this enhancement would be valuable
- **Include Examples**: Show how the enhancement would be used

### 📝 Improving Documentation

Documentation improvements are always appreciated:

- **Code Comments**: Improve inline documentation
- **Examples**: Add more usage examples
- **Guides**: Create tutorials or best practices guides
- **API Documentation**: Enhance resource and data source documentation

## Adding a New Resource

Contributing new Azure resource types is one of the most valuable ways to help the project. Here's a comprehensive guide:

### Step 1: Verify the Resource Isn't Already Implemented

1. **Check the Resource Status**: Look at the [resource status table](../README.md#-resource-status) in the README
2. **Search Issues**: Check if there's already an [open issue](https://github.com/aztfmod/terraform-provider-azurecaf/issues) for this resource
3. **Review Recent PRs**: Make sure someone isn't already working on it

### Step 2: Research Resource Requirements

Before implementing, gather these details:

1. **Naming Rules**: Check the [Azure resource naming rules](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules)
2. **CAF Abbreviation**: Look up the recommended abbreviation in the [CAF resource abbreviations guide](https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations)
3. **Validation Patterns**: Test Azure's actual validation by attempting to create resources with various names

### Step 3: Create an Issue

[Create a new issue](https://github.com/aztfmod/terraform-provider-azurecaf/issues/new) with:

- **Resource Type**: Full Terraform resource name (e.g., `azurerm_synapse_workspace`)
- **Proposed Slug**: CAF abbreviation (e.g., `syws`)
- **Naming Requirements**: Min/max length, allowed characters, case sensitivity
- **Use Case**: Why this resource is needed

### Step 4: Choose the Resource Slug

The slug is a 2-5 character abbreviation that identifies the resource type:

**Guidelines:**
- Keep it short but meaningful
- Follow CAF recommendations when available
- Avoid conflicts with existing slugs
- Use lowercase letters and numbers only

**Examples:**
- Storage Account: `st`
- Key Vault: `kv`
- Virtual Machine: `vm`
- Synapse Workspace: `syws`

### Step 5: Update `resourceDefinition.json`

Add your resource definition to the JSON file:

```json
{
  "name": "azurerm_example_resource",
  "slug": "example",
  "min_length": 3,
  "max_length": 63,
  "lowercase": false,
  "validation_regex": "^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]$",
  "dashes": true,
  "scope": "resourceGroup",
  "regex": "[^a-zA-Z0-9-]",
  "out_of_doc": true,
  "official": {
    "slug": "example",
    "resource": "Example Resource",
    "resource_provider_namespace": "Microsoft.Example/resources"
  }
}
```

**Field Descriptions:**
- `name`: Full Terraform resource type name
- `slug`: CAF abbreviation/prefix
- `min_length`/`max_length`: Azure's length constraints
- `lowercase`: Whether the resource name must be lowercase
- `validation_regex`: Pattern the final name must match
- `dashes`: Whether dashes are allowed
- `scope`: Uniqueness scope (`global`, `resourceGroup`, `parent`)
- `regex`: Pattern to remove invalid characters (cleaning regex)
- `out_of_doc`: (Optional) Set to `true` if the resource is not in the official Azure CAF documentation
- `official`: Official Azure CAF documentation attributes
  - `slug`: Official CAF abbreviation (usually same as root level slug)
  - `resource`: Official resource name from Azure CAF documentation
  - `resource_provider_namespace`: Azure resource provider namespace

### Step 6: Generate and Test

1. **Generate the code:**
   ```bash
   make build
   ```
   
   This runs:
   - `go generate` (regenerates `models_generated.go`)
   - `go fmt ./...` (formats the code)
   - `go test ./...` (runs tests)

2. **Write tests:**
   ```bash
   # Test your specific resource type
   go test ./azurecaf/... -run="TestResourceName.*YourResourceType"
   ```

3. **Verify the implementation:**
   ```bash
   # Create a simple Terraform configuration to test
   echo 'data "azurecaf_name" "test" {
     name = "mytest"
     resource_type = "azurerm_your_new_resource"
   }
   
   output "name" {
     value = data.azurecaf_name.test.result
   }' > test.tf
   
   terraform init && terraform plan
   ```

### Step 7: Update Documentation

1. **Update README**: The resource status table should automatically reflect your addition
2. **Add Examples**: Include usage examples in documentation if the resource type is commonly used
3. **Test Documentation**: Ensure all links and references work correctly

### Step 8: Submit Your Contribution

1. **Commit Changes**: Use descriptive commit messages
   ```bash
   git commit -m "Add support for azurerm_example_resource with 'example' slug"
   ```

2. **Create Pull Request**: Include:
   - Link to the issue you're addressing
   - Description of the resource and its naming constraints
   - Test results showing the implementation works
   - Any special considerations or limitations

## Development Setup

### Prerequisites

- **Go 1.19+**: [Download Go](https://golang.org/dl/)
- **Terraform 1.0+**: [Download Terraform](https://www.terraform.io/downloads.html)
- **Make**: For running build commands
- **Git**: For version control

### Local Development Environment

1. **Clone the repository:**
   ```bash
   git clone https://github.com/aztfmod/terraform-provider-azurecaf.git
   cd terraform-provider-azurecaf
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build the provider:**
   ```bash
   make build
   ```

4. **Run tests:**
   ```bash
   # Unit tests only
   make unittest
   
   # Integration tests (optional)
   make test_integration
   ```

### Development Workflow

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/add-new-resource-type
   ```

2. **Make your changes:**
   - Update `resourceDefinition.json`
   - Run `make build` to regenerate code
   - Add/update tests as needed

3. **Test your changes:**
   ```bash
   make test_all
   ```

4. **Commit and push:**
   ```bash
   git add .
   git commit -m "Descriptive commit message"
   git push origin feature/add-new-resource-type
   ```

## Testing Guidelines

### Unit Tests

Unit tests should cover:
- **Name Generation**: Verify correct name patterns
- **Validation**: Test regex patterns and constraints
- **Edge Cases**: Test boundary conditions and error cases
- **Different Configurations**: Test various parameter combinations

Example test structure:
```go
func TestResourceName_ExampleResource(t *testing.T) {
    tests := []struct {
        name     string
        input    map[string]interface{}
        expected string
        wantErr  bool
    }{
        {
            name: "basic name generation",
            input: map[string]interface{}{
                "name":          "mytest",
                "resource_type": "azurerm_example_resource",
            },
            expected: "example-mytest",
            wantErr:  false,
        },
        // Add more test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

Integration tests verify that the provider works correctly with Terraform:

```go
func TestAccDataSourceAzureCAFName_ExampleResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccDataSourceAzureCAFName_ExampleResource(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("data.azurecaf_name.test", "result", "example-mytest"),
                ),
            },
        },
    })
}
```

### Test Coverage

Maintain high test coverage by:
- Testing all new functionality
- Adding edge case tests
- Testing error conditions
- Verifying regex patterns work correctly

## Submitting Changes

### Pull Request Guidelines

**Before submitting:**
- ✅ All tests pass (`make test_all`)
- ✅ Code is properly formatted (`go fmt`)
- ✅ Documentation is updated
- ✅ Commit messages are descriptive
- ✅ Changes are focused and atomic

**Pull Request Template:**
```markdown
## Description
Brief description of the changes

## Related Issue
Fixes #123

## Type of Change
- [ ] Bug fix
- [ ] New feature (new resource type)
- [ ] Documentation update
- [ ] Code refactoring

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
```

### Review Process

1. **Automated Checks**: CI/CD will run tests and checks
2. **Code Review**: Maintainers will review your changes
3. **Feedback**: Address any requested changes
4. **Merge**: Once approved, changes will be merged

## Documentation Standards

### Code Documentation

- **Function Comments**: Document all public functions
- **Complex Logic**: Explain non-obvious code sections
- **Examples**: Include usage examples where helpful
- **Parameter Documentation**: Document all parameters and return values

### Markdown Standards

- **Consistent Formatting**: Use standard markdown conventions
- **Code Blocks**: Always specify language for syntax highlighting
- **Links**: Use relative links for internal documentation
- **Tables**: Use consistent table formatting

## Recognition

Contributors are recognized in:
- Git commit history
- Release notes for significant contributions
- Community discussions and announcements

Thank you for contributing to the Azure CAF Terraform Provider! 🚀
