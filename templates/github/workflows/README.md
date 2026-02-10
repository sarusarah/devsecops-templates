# GitHub Actions Reusable Workflows

This directory contains reusable GitHub Actions workflows for DevSecOps CI/CD pipelines with monorepo support.

## Overview

These reusable workflows provide the same security scanning, build, and test capabilities as the GitLab CI templates, designed specifically for GitHub Actions.

## Monorepo Support

All workflows support monorepo configurations via the `project_path` input parameter:

- **Default behavior:** When `project_path` is not specified or set to `.`, workflows operate on the repository root
- **Monorepo mode:** Set `project_path` to a subdirectory (e.g., `frontend`, `backend`) to scope operations to that project

## Usage Pattern

### Single Workflow File (Recommended for Monorepo)

Create separate workflow files per project using native GitHub Actions path filtering:

```yaml
# .github/workflows/frontend.yml
name: Frontend CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'frontend/**'
  pull_request:
    branches: [main, develop]
    paths:
      - 'frontend/**'

jobs:
  build:
    uses: platform/devsecops-template/.github/workflows/build-node.yml@v1.0.1
    with:
      project_path: frontend
      node_version: '20'

  test:
    needs: build
    uses: platform/devsecops-template/.github/workflows/test-node.yml@v1.0.1
    with:
      project_path: frontend

  secrets:
    uses: platform/devsecops-template/.github/workflows/security-secrets.yml@v1.0.1
    with:
      project_path: frontend
    permissions:
      contents: read

  sast:
    uses: platform/devsecops-template/.github/workflows/security-sast.yml@v1.0.1
    with:
      project_path: frontend

  dependency-scan:
    uses: platform/devsecops-template/.github/workflows/security-dependency.yml@v1.0.1
    with:
      project_path: frontend
      language: node
```

```yaml
# .github/workflows/backend.yml
name: Backend CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'backend/**'
  pull_request:
    branches: [main, develop]
    paths:
      - 'backend/**'

jobs:
  build:
    uses: platform/devsecops-template/.github/workflows/build-python.yml@v1.0.1
    with:
      project_path: backend
      python_version: '3.12'

  test:
    needs: build
    uses: platform/devsecops-template/.github/workflows/test-python.yml@v1.0.1
    with:
      project_path: backend

  secrets:
    uses: platform/devsecops-template/.github/workflows/security-secrets.yml@v1.0.1
    with:
      project_path: backend
    permissions:
      contents: read
```

## Available Workflows

### Build Workflows
- `build-node.yml` - Node.js build with npm/yarn/pnpm support
- `build-python.yml` - Python build with pip/pyproject.toml support
- `build-php.yml` - PHP build with Composer support

### Test Workflows
- `test-node.yml` - Node.js testing
- `test-python.yml` - Python testing with pytest and coverage
- `test-php.yml` - PHP testing with PHPUnit

### Security Workflows
- `security-secrets.yml` - Secret scanning with Trivy/Gitleaks
- `security-sast.yml` - Static analysis with Trivy/Semgrep
- `security-dependency.yml` - Dependency vulnerability scanning
- `security-container.yml` - Container image scanning
- `security-iac.yml` - Infrastructure as Code scanning

## Common Input Parameters

### All Workflows
- `project_path` (string, default: `.`) - Project directory for monorepo support

### Build Workflows
- `node_version` (string, default: `'20'`) - Node.js version for build-node.yml
- `python_version` (string, default: `'3.12'`) - Python version for build-python.yml
- `php_version` (string, default: `'8.3'`) - PHP version for build-php.yml
- `package_manager` (string, default: `'npm'`) - Package manager for Node.js (npm/yarn/pnpm)

### Security Workflows
- `language` (string) - Language for dependency scanning (node/python/php)
- `security_scanner` (string, default: `'trivy'`) - Scanner to use (trivy/specialized)
- `severity` (string, default: `'HIGH,CRITICAL'`) - Vulnerability severity levels

## Change Detection

Use GitHub Actions native `paths` filter for change detection:

```yaml
on:
  push:
    paths:
      - 'frontend/**'      # Only trigger on frontend changes
  pull_request:
    paths:
      - 'frontend/**'
```

This is more efficient than using conditional `if` statements and provides better UI feedback in GitHub.

## Differences from GitLab CI Templates

| Feature | GitLab CI | GitHub Actions |
|---------|-----------|----------------|
| Job reuse | `extends` keyword | Reusable workflows (`uses`) |
| Change detection | `rules:changes` | `paths` filter in triggers |
| Working directory | `cd` in before_script | `working-directory` in steps |
| Variables | Job-level variables | Workflow inputs |
| Artifacts | Automatic between jobs | Manual upload/download |

## Migration from GitLab CI

If you're migrating from GitLab CI templates:

1. Create separate workflow files per project (instead of single .gitlab-ci.yml)
2. Replace `extends` with `uses` for reusable workflows
3. Replace `rules:changes` with `paths:` in workflow triggers
4. Replace `PROJECT_PATH` variable with `project_path` input parameter
5. Add explicit `needs` for job dependencies

## Examples

See the [examples/monorepo-github/](../../../examples/monorepo-github/) directory for complete working GitHub Actions examples.

For GitLab CI examples, see [examples/monorepo-gitlab/](../../../examples/monorepo-gitlab/).

## Contributing

When adding new workflows:
1. Follow the `workflow_call` pattern for reusability
2. Support `project_path` input for monorepo compatibility
3. Use `working-directory: ${{ inputs.project_path }}` in steps
4. Document all input parameters
5. Test with both single-project and monorepo configurations
