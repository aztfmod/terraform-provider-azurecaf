# GitHub Workflows

This directory contains GitHub Actions workflows for the terraform-provider-azurecaf project.

## Workflows

### `go.yml`
Main CI/CD workflow that:
- Builds the provider
- Runs comprehensive tests (unit, integration, coverage)
- Creates releases when tags are pushed

### `codeql.yml`
CodeQL security analysis workflow that:
- Scans code for security vulnerabilities
- Runs on push/PR to main branch
- Scheduled to run weekly on Mondays

### `security.yml`
Additional security scanning workflow that:
- Runs Gosec security scanner
- Checks for dependency vulnerabilities with Nancy
- Runs daily at 2 AM UTC

## Dependabot

The `dependabot.yml` configuration file automatically:
- Updates Go module dependencies weekly
- Updates GitHub Actions weekly
- Creates PRs with dependency updates
- Assigns PRs to maintainers

## Security

See `SECURITY.md` in the root directory for security policy and vulnerability reporting process.
