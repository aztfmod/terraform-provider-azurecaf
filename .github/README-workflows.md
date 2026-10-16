# GitHub Workflows

This directory contains GitHub Actions workflows for the terraform-provider-azurecaf project.

## Workflows

### Agentic workflows

Agentic workflow sources are the Markdown files in `.github/workflows/`; their
generated `.lock.yml` files are committed for GitHub Actions to execute.

Copilot authentication uses the built-in GitHub Actions token by granting
`copilot-requests: write` in each workflow's `permissions` block. GitHub mints
an ephemeral token for each run and bills usage through the organization's
Copilot subscription, so the workflows do not require a
`COPILOT_GITHUB_TOKEN` personal access token or repository secret.

This requires centralized Copilot billing to be enabled under
**Organization settings → Copilot → Policies → Copilot CLI → Allow use of
Copilot CLI billed to the organization**. Without that policy, inference fails
and the workflows fail until the organization billing policy is enabled. When
`copilot-requests: write` is enabled, any existing `COPILOT_GITHUB_TOKEN`
secret is ignored for inference and can be removed after the migration is
verified.

`GH_AW_GITHUB_TOKEN` is separate from Copilot inference authentication. It is
an optional fallback credential for GitHub API operations that need access
beyond the built-in `GITHUB_TOKEN`; these same-repository workflows do not
require it for Copilot inference.

The workflows explicitly select `gpt-5.6-luna`. Do not replace it with the
`auto` model alias until upstream model routing no longer selects the
utility-only `gpt-5.6-luna-utility` model for the chat-completions endpoint.

After changing an agentic workflow, regenerate and validate the lock files:

```bash
gh aw compile --validate
```

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
