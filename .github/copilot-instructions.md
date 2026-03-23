DO:
- Update ./CHANGELOG.md with the changes you do and assess the impact of the changes.
- Write tests for the code you write.
- Write documentation for the code you write.

## Available Skills

Skills are small, composable building blocks. Each has a SKILL.md with step-by-step instructions.

### Resource Lifecycle
| Skill | Location | Purpose |
|-------|----------|---------|
| `azure-naming-research` | `.github/skills/azure-naming-research/` | Research CAF abbreviations and Azure naming constraints |
| `resource-definition-json` | `.github/skills/resource-definition-json/` | Lookup, compare, format, insert/update entries in resourceDefinition.json |
| `provider-build-test` | `.github/skills/provider-build-test/` | Regenerate Go code, build, and run unit tests |
| `terraform-mock-test` | `.github/skills/terraform-mock-test/` | Validate with mock azurerm provider (no Azure credentials) |
| `changelog-update` | `.github/skills/changelog-update/` | Add CHANGELOG.md entries with semver impact assessment |
| `readme-resource-table` | `.github/skills/readme-resource-table/` | Update README.md resource status table |
| `resource-bulk-import` | `.github/skills/resource-bulk-import/` | Batch-research and insert multiple resources |
| `resource-diff-report` | `.github/skills/resource-diff-report/` | Compare two versions of resourceDefinition.json |
| `resource-completeness-check` | `.github/skills/resource-completeness-check/` | Compare provider coverage against known azurerm resources |

### CI/CD & Testing
| Skill | Location | Purpose |
|-------|----------|---------|
| `regression-test-runner` | `.github/skills/regression-test-runner/` | Run full CI test suite and report results |
| `e2e-test-runner` | `.github/skills/e2e-test-runner/` | Run E2E tests and produce structured summary |
| `coverage-analysis` | `.github/skills/coverage-analysis/` | Run coverage and check against 95% threshold |
| `test-failure-diagnosis` | `.github/skills/test-failure-diagnosis/` | Analyze test failures and suggest fixes |

### Release Management
| Skill | Location | Purpose |
|-------|----------|---------|
| `semver-assessment` | `.github/skills/semver-assessment/` | Determine version bump from changes |
| `release-notes-generator` | `.github/skills/release-notes-generator/` | Generate release notes from CHANGELOG |
| `pre-release-validation` | `.github/skills/pre-release-validation/` | Run all pre-release checks |

### Community
| Skill | Location | Purpose |
|-------|----------|---------|
| `pr-compliance-check` | `.github/skills/pr-compliance-check/` | Validate PR against project checklist |
| `issue-to-resource-spec` | `.github/skills/issue-to-resource-spec/` | Parse issue into draft resource definition |
| `contributor-guide` | `.github/skills/contributor-guide/` | Step-by-step contribution guidance |

### Documentation
| Skill | Location | Purpose |
|-------|----------|---------|
| `docs-resource-sync` | `.github/skills/docs-resource-sync/` | Keep docs in sync with resource definitions |
| `example-generator` | `.github/skills/example-generator/` | Generate Terraform example configurations |

### Azure Sync
| Skill | Location | Purpose |
|-------|----------|---------|
| `azure-caf-sync` | `.github/skills/azure-caf-sync/` | Detect CAF slug drift from official docs |
| `azure-resource-discovery` | `.github/skills/azure-resource-discovery/` | Discover new azurerm resources not yet supported |
| `naming-rules-drift-check` | `.github/skills/naming-rules-drift-check/` | Detect when Azure naming rules change |

## Available Agents

Agents orchestrate skills into multi-step workflows.

### Interactive (Copilot Chat)
| Agent | File | Purpose |
|-------|------|---------|
| `caf-check-resource` | `.github/agents/caf-check-resource.agent.md` | Validate a resource definition (existing) |
| `caf-add-resource` | `.github/agents/caf-add-resource.agent.md` | End-to-end new resource addition |
| `caf-update-resource` | `.github/agents/caf-update-resource.agent.md` | Update an existing resource definition |
| `caf-bulk-add-resources` | `.github/agents/caf-bulk-add-resources.agent.md` | Add multiple resources in one session |
| `caf-audit-resources` | `.github/agents/caf-audit-resources.agent.md` | Full audit: completeness, drift, coverage |
| `caf-release-prep` | `.github/agents/caf-release-prep.agent.md` | Prepare a release: validate, version, notes |
| `caf-diagnose-failure` | `.github/agents/caf-diagnose-failure.agent.md` | Diagnose and fix build/test failures |

### Automated (GitHub Actions Agentic Workflows)
| Workflow | File | Trigger |
|----------|------|---------|
| Daily Repo Status | `.github/workflows/daily-repo-status.md` | Daily (existing) |
| Issue Triage | `.github/workflows/issue-triage-agent.md` | Weekdays (existing) |
| Issue Arborist | `.github/workflows/issue-arborist.md` | Daily (existing) |
| Nightly Regression | `.github/workflows/nightly-regression.md` | Nightly at 3 AM UTC |
| PR Review | `.github/workflows/pr-review-agent.md` | PR opened/updated |
| Issue to PR | `.github/workflows/issue-to-pr-agent.md` | Issue labeled `new-resource` |
| Contributor Welcome | `.github/workflows/contributor-welcome.md` | First-time contributor PR |
| Weekly Azure Sync | `.github/workflows/weekly-azure-sync.md` | Weekly Monday 9 AM UTC |
| Release Validation | `.github/workflows/release-validation.md` | Tag push (v*) |