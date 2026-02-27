# GitLab DevSecOps CI/CD Template Library

A comprehensive, enterprise-grade GitLab CI/CD template library implementing DevSecOps best practices with automated security scanning, testing, and deployment.

---

## Quick Start

### Using The Templates

Include the templates you need in your `.gitlab-ci.yml` and enable desired security scans:
```yaml
# your-project/.gitlab-ci.yml
include:
  - project: platform/devsecops-template
    ref: v1.0.1
    file:
      - /templates/gitlab/base.yml
      - /templates/gitlab/security/secrets.yml
      - /templates/gitlab/security/dependency.yml
      - /templates/gitlab/security/sast.yml
      - /templates/gitlab/security/dtrack.yml  # Optional: SBOM upload to Dependency-Track

variables:
  LANGUAGE: "node"  # or python, php
  ENABLE_DAST: "true"
  STAGING_URL: "https://staging.example.com"

  # Optional: Enable Dependency-Track integration
  ENABLE_DTRACK: "true"
  DTRACK_URL: "https://api.dtrack.example.com"
  DTRACK_API_KEY: "${DTRACK_API_KEY}"

# Your custom jobs here...
```
---

## What's Included

### Security Scanning Templates (`templates/gitlab/security/`)

**Unified Scanning (Trivy - Default)**
- **Secrets Detection** - Trivy secret scanning
- **Dependency Scanning** - Trivy vulnerability scanning (SCA)
- **SAST** - Trivy misconfiguration detection
- **Container Scanning** - Trivy image vulnerability scanning
- **IaC Security** - Trivy configuration scanning

**Specialized Tools (Alternative)**
- **Secrets Detection** - Gitleaks for detecting committed credentials
- **Dependency Scanning** - npm audit, pip-audit, composer audit
- **SAST** - Semgrep comprehensive static analysis
- **IaC Security** - Kubeconform, Kube-Score, Polaris

**Runtime Security**
- **DAST** - OWASP ZAP dynamic application security testing

**SBOM & Dependency Tracking**
- **Dependency-Track Integration** - Upload CycloneDX SBOMs to Dependency-Track for centralized dependency and vulnerability tracking

### AI-Powered Pipeline Reporting (`templates/gitlab/ai-report.yml`, `templates/github/ai-report.yml`)
- **AI Analysis** - Gemini-powered analysis of all pipeline stage outputs
- **Slack Notifications** - Consolidated, color-coded summary sent to Slack channels
- **Replaces Raw Logs** - Actionable summaries instead of 2000-line log dumps

**📖 Full Documentation:** [docs/AI_REPORTING.md](docs/AI_REPORTING.md)

### Pipeline Templates (`templates/gitlab/` and `templates/github/`)
- **Base Configuration** - Stages, rules, and variables
- **Workflow Rules** - When pipelines should run
- **Build Templates** - Node.js, Python, PHP
- **Test Templates** - Unit testing for all languages
- **GitOps Deployment** - Staging and production via GitOps
- **Monitoring** - Post-deployment health checks
- **Reporting** - Aggregate security findings

### Examples (`examples/`)
- Node.js/Nuxt application
- Python application
- PHP Symfony application
- PHP Drupal application
- **Monorepo examples** - Frontend (Node.js) + Backend (Python) with independent CI/CD (GitLab CI & GitHub Actions)

### GitHub Actions Support (`templates/github/`)
- **Reusable workflows** for GitHub Actions
- Build workflows (Node.js, Python, PHP)
- Test workflows with coverage reporting
- Security workflows (secrets, SAST, dependency, container, IaC, Dependency-Track)
- Full monorepo support with `project_path` input parameter

### Testing Tools
- **Dagger Module** (`dagger/`) - Local pipeline testing
- **gitlab-ci-local Config** - Alternative local testing
- **Interactive Test Script** (`test-local.sh`)

---

## Security Features

**Secrets detection** in preflight stage (blocks by default)
- **Dependency vulnerability scanning** for all package managers
- **Static code analysis** (SAST) with Semgrep
- **Dynamic security testing** (DAST) with OWASP ZAP
- **Container image scanning** with Trivy
- **Infrastructure as Code** validation
- **Security policy enforcement** (strict/permissive modes)
- **GitLab Security Dashboard** integration
- **Audit trail** with 7-day artifact retention

---

## Testing Your Pipeline Locally

### Option 1: Dagger (Recommended)

```bash
# Install Dagger
curl -L https://dl.dagger.io/dagger/install.sh | sh

# Run all security scans
dagger call test --source=. --language=node

# Run individual scans
dagger call secrets-detection --source=.
dagger call sast-scanning --source=.
```

### Option 2: gitlab-ci-local

```bash
# Install
npm install -g gitlab-ci-local

# Interactive testing
./test-local.sh

# Manual testing
cd examples/node
gitlab-ci-local --preview
```
## Key Features

### Security by Default
- All security scans enabled by default
- Secrets detection blocks pipelines
- GitLab Security Dashboard integration
- Compliance-ready (ISO 27001, NIS2, OWASP ASVS)

### Polyglot Support
- Node.js (npm, yarn, pnpm)
- Python (pip, requirements.txt, pyproject.toml)
- PHP (Composer)
- Extensible for other languages

### Flexible Configuration

**Choose Your Security Scanner:**
```yaml
variables:
  # Unified scanning with Trivy (recommended for consistency)
  SECURITY_SCANNER: "trivy"
  # Or use specialized tools for comprehensive analysis
  SECURITY_SCANNER: "specialized"
  SAST_TOOL: "semgrep"  # More thorough SAST

  # Feature toggles
  ENABLE_SECRETS: "true"
  ENABLE_DEPENDENCY_SCAN: "true"
  ENABLE_SAST: "true"
  ENABLE_CONTAINER_SCAN: "false"
  ENABLE_IAC_SCAN: "false"
  ENABLE_DAST: "false"

  # Security policy
  SECURITY_POLICY: "strict"  # or "permissive"
```

**Trivy vs Specialized Tools:**
- **Trivy**: Single tool, faster pipelines, consistent output format, easier maintenance
- **Specialized**: More comprehensive detection, language-specific optimizations, mature tooling

### GitOps Deployment
```yaml
variables:
  GITOPS_REPO: "git@gitlab.com:gitops/myapp.git"
  GITOPS_PATH: "values.yaml"
  GITOPS_IMAGE_TAG_YQ_PATH: ".image.tag"
```

---

##  Architecture

### Pipeline Stages (OWASP SPVS Aligned)

This template follows the **[OWASP Secure Pipeline Verification Standard (SPVS)](https://owasp.org/www-project-spvs/)** framework, which defines five key phases for secure software delivery: **Develop**, **Integrate**, **Release**, and **Operate**.

```
┌─────────────────────────────────────────────────────────────────────────┐
│ OWASP SPVS Phase: DEVELOP - Secure coding and continuous reviews       │
├─────────────────────────────────────────────────────────────────────────┤
│ 1. source   → Secrets detection, dependency scanning, SAST             │
│ 2. build    → Application build with artifacts                         │
│ 3. test     → Unit tests with coverage                                 │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ OWASP SPVS Phase: INTEGRATE - Automated validation & artifact integrity│
├─────────────────────────────────────────────────────────────────────────┤
│ 4. package  → Container scanning, IaC scanning                         │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ OWASP SPVS Phase: RELEASE - Final validations & secure deployment      │
├─────────────────────────────────────────────────────────────────────────┤
│ 5. stage    → GitOps deployment to staging                             │
│ 6. verify   → DAST against staging, E2E tests                          │
│ 7. deploy   → Manual deployment to production                          │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ OWASP SPVS Phase: OPERATE - Continuous monitoring & incident response  │
├─────────────────────────────────────────────────────────────────────────┤
│ 8. operate  → Health checks and monitoring                             │
│ 9. report   → Aggregate security findings and metrics                  │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│ AI-Powered Reporting (Optional - ENABLE_AI_REPORT: "true")             │
├─────────────────────────────────────────────────────────────────────────┤
│ 10. ai-analysis → Gemini analyzes each stage's output                  │
│ 11. ai-summary  → Consolidated summary → Slack notification            │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Stage Mapping & Security Activities

| Stage | OWASP SPVS Phase | Security Activities |
|-------|------------------|---------------------|
| `source` | **Develop** | Secret scanning, SCA (dependencies), SAST |
| `build` | **Develop** | Artifact creation, build integrity |
| `test` | **Develop** | Unit tests, integration tests, code coverage |
| `package` | **Integrate** | Container image scanning, IaC validation |
| `stage` | **Release** | Staging deployment via GitOps |
| `verify` | **Release** | DAST, E2E testing, acceptance tests |
| `deploy` | **Release** | Production deployment via GitOps |
| `operate` | **Operate** | Health monitoring, availability checks |
| `report` | **Operate** | Security metrics, compliance reporting |
| `ai-analysis` | **Reporting** | AI-powered per-stage analysis (Gemini) |
| `ai-summary` | **Reporting** | Consolidated summary + Slack notification |

### Template Inheritance

```yaml
# Base template provides stages and rules
include:
  - local: /templates/gitlab/base.yml

# Extend with specific security scans
include:
  - local: /templates/gitlab/security/secrets.yml
  - local: /templates/gitlab/security/sast.yml

# Use pre-built job templates
build-app:
  extends: .build:node
  variables:
    NODE_VERSION: "20"
```

---

## Monorepo Support

The DevSecOps templates support monorepos with multiple projects in a single repository. Each project can run CI stages independently with automatic change detection.

### Quick Example

**GitLab CI (.gitlab-ci.yml):**
```yaml
include:
  - project: platform/devsecops-template
    ref: v1.0.1
    file:
      - /templates/gitlab/base.yml
      - /templates/build.yml
      - /templates/gitlab/security/secrets.yml

# Frontend project - only runs when frontend/ changes
build:frontend:
  extends: .build:node
  variables:
    PROJECT_PATH: frontend
    NODE_VERSION: "20"
  rules:
    - changes:
        - frontend/**/*

secrets-detection:frontend:
  extends: secrets-detection
  variables:
    PROJECT_PATH: frontend
  rules:
    - changes:
        - frontend/**/*

# Backend project - only runs when backend/ changes
build:backend:
  extends: .build:python
  variables:
    PROJECT_PATH: backend
    PYTHON_VERSION: "3.12"
  rules:
    - changes:
        - backend/**/*

secrets-detection:backend:
  extends: secrets-detection
  variables:
    PROJECT_PATH: backend
  rules:
    - changes:
        - backend/**/*
```

**GitHub Actions (.github/workflows/frontend.yml):**
```yaml
name: Frontend CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'frontend/**'

jobs:
  build:
    uses: platform/devsecops-template/.github/workflows/build-node.yml@v1.0.1
    with:
      project_path: frontend
      node_version: '20'

  secrets:
    uses: platform/devsecops-template/.github/workflows/security-secrets.yml@v1.0.1
    with:
      project_path: frontend
    permissions:
      contents: read
```

### Key Features

- **Per-project CI/CD:** Each project runs its own build, test, and security scans
- **Change detection:** Jobs only run when project files change
- **Independent deployment:** Projects can be deployed separately
- **PROJECT_PATH variable:** Scopes all operations to project directory
- **Efficient:** Reduces CI minutes by only testing affected projects

### How It Works

**GitLab CI:**
- Use `PROJECT_PATH` variable to specify project directory
- Use `rules:changes` for automatic change detection
- Extend base jobs with project-specific configuration

**GitHub Actions:**
- Create separate workflow files per project
- Use `paths:` filter for change detection
- Call reusable workflows with `project_path` input

### Full Example

See complete working examples:
- **GitLab CI:** [examples/monorepo-gitlab/](examples/monorepo-gitlab/)
- **GitHub Actions:** [examples/monorepo-github/](examples/monorepo-github/)

Each example includes:
- Frontend (Node.js) and Backend (Python) projects
- Complete CI/CD configurations
- Change detection rules
- Documentation and best practices

### Benefits

**Faster CI/CD** - Only affected projects run
**Clear separation** - Each project has explicit CI configuration
**Independent deployment** - Different release cadences per project
**Scalable** - Easy to add new projects
**KISS principle** - Simple, no complex logic or dynamic generation

---

## Configuration Variables

### Language Selection
```yaml
LANGUAGE: "node"  # node|python|php|generic
NODE_VERSION: "20"
PYTHON_VERSION: "3.12"
PHP_VERSION: "8.3"
PACKAGE_MANAGER: "npm"  # npm|yarn|pnpm
```

### Security Scanning
```yaml
# Security tool selection
SECURITY_SCANNER: "trivy"  # trivy|specialized
# trivy: Use Trivy for unified security scanning (recommended)
# specialized: Use specialized tools (Gitleaks, Semgrep, etc.)

# SAST tool selection (when using specialized tools or mixed approach)
SAST_TOOL: "semgrep"  # semgrep|trivy
# semgrep: Comprehensive SAST with Semgrep (recommended for thorough analysis)
# trivy: Basic SAST with Trivy misconfiguration detection

# Feature toggles
ENABLE_SECRETS: "true"
ENABLE_DEPENDENCY_SCAN: "true"
ENABLE_SAST: "true"
ENABLE_CONTAINER_SCAN: "false"
ENABLE_IAC_SCAN: "false"
ENABLE_DAST: "false"
ENABLE_DTRACK: "false"  # Dependency-Track SBOM upload
ENABLE_AI_REPORT: "false"  # AI pipeline analysis + Slack summary

# Security policy
SECURITY_POLICY: "strict"  # strict|permissive
```

### Trivy Configuration
```yaml
TRIVY_SEVERITY: "CRITICAL,HIGH"
TRIVY_EXIT_CODE: "1"
IMAGE_NAME: "${CI_REGISTRY_IMAGE}"
IMAGE_TAG: "${CI_COMMIT_SHA}"
```

### Container Scanning Configuration

The container scanning template scans Docker images for vulnerabilities using Trivy. It supports various deployment patterns and custom image naming schemes.

**📖 Full Documentation:** [docs/CONTAINER_SCANNING.md](docs/CONTAINER_SCANNING.md)

#### Basic Configuration

```yaml
include:
  - project: platform/devsecops-template
    file: /templates/gitlab/security/container.yml

variables:
  ENABLE_CONTAINER_SCAN: "true"
  IMAGE_NAME: "${CI_REGISTRY_IMAGE}"
  IMAGE_TAG: "latest"
```

#### Advanced Configuration: Sub-Images

For projects that build multiple images (e.g., `myapp/staging:latest`, `myapp/prod:latest`):

```yaml
variables:
  IMAGE_NAME: "${CI_REGISTRY_IMAGE}"
  IMAGE_TAG: "latest"
  CONTAINER_IMAGE_SUFFIX: "staging"  # Scans CI_REGISTRY_IMAGE/staging:latest

# Or override per-branch:
container-security-scan:staging:
  extends: container-security-scan
  variables:
    CONTAINER_IMAGE_SUFFIX: "staging"
  rules:
    - if: '$ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_BRANCH == "staging"'
  needs:
    - "build-staging-image"

container-security-scan:prod:
  extends: container-security-scan
  variables:
    CONTAINER_IMAGE_SUFFIX: "prod"
  rules:
    - if: '$ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_BRANCH == "prod"'
  needs:
    - "build-prod-image"

# Disable default job
container-security-scan:
  rules:
    - when: never
```

#### Custom Image Naming

```yaml
variables:
  IMAGE_NAME: "docker.io/myorg/myapp"  # External registry
  IMAGE_TAG: "${CI_COMMIT_SHORT_SHA}"
  
container-security-scan:
  needs:
    - "push-to-dockerhub"
```

#### Available Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `IMAGE_NAME` | `${CI_REGISTRY_IMAGE}` | Full image repository path |
| `IMAGE_TAG` | `${CI_COMMIT_SHORT_SHA}` | Image tag to scan |
| `CONTAINER_IMAGE_SUFFIX` | _(empty)_ | Sub-path for image (e.g., `staging`, `prod`) |
| `TRIVY_SEVERITY` | `UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL` | Severity levels to report |
| `TRIVY_EXIT_CODE` | `0` | Exit code for failures (set to `1` to fail pipeline) |
| `TRIVY_NON_SSL` | `false` | Set to `true` for insecure registries |

#### Image Resolution

The template builds the image reference as follows:

1. Base: `IMAGE_NAME` (defaults to `CI_REGISTRY_IMAGE`)
2. If `CONTAINER_IMAGE_SUFFIX` is set: `IMAGE_NAME/CONTAINER_IMAGE_SUFFIX`
3. Tag: `:IMAGE_TAG`

**Examples:**
- `IMAGE_NAME=registry/project`, `IMAGE_TAG=v1.0` → `registry/project:v1.0`
- `IMAGE_NAME=registry/project`, `CONTAINER_IMAGE_SUFFIX=staging`, `IMAGE_TAG=latest` → `registry/project/staging:latest`

### Dependency-Track Integration

Upload Software Bill of Materials (SBOM) to Dependency-Track for centralized dependency and vulnerability tracking. Supports both single-project and monorepo configurations with automatic project creation.

**📖 Full Documentation:** [docs/DEPENDENCY_TRACK.md](docs/DEPENDENCY_TRACK.md)

#### Quick Start (GitLab CI)

```yaml
include:
  - project: platform/devsecops-template
    ref: v1.0.1
    file:
      - /templates/gitlab/base.yml
      - /templates/gitlab/security/dtrack.yml

variables:
  ENABLE_DTRACK: "true"
  DTRACK_URL: "https://api.dtrack.example.com"
  DTRACK_API_KEY: "${DTRACK_API_KEY}"  # Set as CI/CD secret
```

#### Quick Start (GitHub Actions)

```yaml
jobs:
  dtrack:
    uses: platform/devsecops-template/.github/workflows/security-dtrack.yml@v1.0.1
    with:
      dtrack_url: "https://api.dtrack.example.com"
    secrets:
      dtrack_api_key: ${{ secrets.DTRACK_API_KEY }}
```

#### Key Variables

| Variable (GitLab) / Input (GitHub) | Description |
|------------------------------------|-------------|
| `DTRACK_URL` / `dtrack_url` | DTrack base URL (required) |
| `DTRACK_API_KEY` / `dtrack_api_key` | DTrack API key (required, secret) |
| `PROJECT_PATH` / `project_path` | Monorepo subproject path (e.g., `frontend`) |
| `DTRACK_PROJECT_UUID` / `dtrack_project_uuid` | Explicit project UUID (optional) |

**For complete configuration options, monorepo examples, and troubleshooting, see [docs/DEPENDENCY_TRACK.md](docs/DEPENDENCY_TRACK.md)**

### GitOps Deployment
```yaml
GITOPS_REPO: "git@gitlab.com:gitops/app.git"
GITOPS_PATH: "values.yaml"
GITOPS_IMAGE_TAG_YQ_PATH: ".image.tag"
GITOPS_SSH_KEY: "${GITOPS_SSH_KEY}"
GITOPS_PRODUCTION_BRANCH: "production"
```

### URLs
```yaml
STAGING_URL: "https://staging.example.com"
PRODUCTION_URL: "https://production.example.com"
HEALTHCHECK_URL: "${STAGING_URL}/health"
```

### AI-Powered Pipeline Reporting

Automated pipeline analysis using Google Gemini with Slack notifications. Analyzes all non-deployment stage outputs and sends a single, actionable summary to Slack.

**📖 Full Documentation:** [docs/AI_REPORTING.md](docs/AI_REPORTING.md)

#### Quick Start (GitLab CI)

```yaml
include:
  - project: platform/devsecops-template
    file:
      - /templates/gitlab/base.yml
      - /templates/gitlab/ai-report.yml

variables:
  ENABLE_AI_REPORT: "true"
  # AI_PROVIDER: "openai"  # Optional: switch to OpenAI (default: "gemini")
  # Set AI_API_KEY as CI/CD secret (Gemini or OpenAI key)
  # Set SLACK_WEBHOOK_URL as CI/CD secret (optional)
```

#### Quick Start (GitHub Actions)

```yaml
jobs:
  ai-report:
    needs: [build, test, sast, dependency-scan]
    if: always()
    uses: ./.github/workflows/ai-report.yml
    secrets:
      ai_api_key: ${{ secrets.AI_API_KEY }}
      slack_webhook_url: ${{ secrets.SLACK_WEBHOOK_URL }}
```

#### Key Variables

| Variable / Secret | Description |
|-------------------|-------------|
| `ENABLE_AI_REPORT` | Feature toggle (default: `"false"`) |
| `AI_API_KEY` | API key for Gemini or OpenAI (CI/CD secret) |
| `AI_PROVIDER` | `"gemini"` (default) or `"openai"` |
| `SLACK_WEBHOOK_URL` | Slack incoming webhook URL (CI/CD secret, optional) |
| `AI_MODEL` | Model override (default: auto per provider) |

### Notifications
```yaml
MATTERMOST_WEBHOOK_URL: "https://mattermost.com/hooks/xxx"
```

---

## Troubleshooting

### Secrets Detection Failing

**Problem:** Gitleaks finds false positives
**Solution:** Create `.gitleaksignore`:
```
# .gitleaksignore
test-fixtures/fake-secret.txt
docs/examples/api-key-example.md
```

### Dependency Scan Not Finding Dependencies

**Problem:** Missing lock files
**Solution:** Ensure these exist:
- Node.js: `package-lock.json` or `yarn.lock`
- Python: `requirements.txt`
- PHP: `composer.lock`

### DAST Failing to Connect

**Problem:** Staging URL not accessible
**Solution:**
1. Verify `STAGING_URL` is correct
2. Ensure staging is deployed before DAST runs
3. Check network/firewall rules

### Container Scan Authentication Failed

**Problem:** Can't pull private images
**Solution:** Ensure CI/CD variables set:
- `CI_REGISTRY`
- `CI_REGISTRY_USER`
- `CI_REGISTRY_PASSWORD`

### Container Scan: Image Not Found

**Problem:** Trivy can't find the specified image
```
unable to find the specified image "registry/project:tag"
```

**Solution:** Ensure the image exists with the correct tag before scanning:

1. **Set correct IMAGE_TAG variable:**
   ```yaml
   variables:
     IMAGE_TAG: "latest"  # Must match your actual build tag
   ```

2. **For sub-images (e.g., staging/prod):**
   ```yaml
   variables:
     CONTAINER_IMAGE_SUFFIX: "staging"  # Scans CI_REGISTRY_IMAGE/staging:TAG
   ```

3. **Add needs dependency:**
   ```yaml
   container-security-scan:
     needs:
       - "build-and-push-image"  # Wait for image to be pushed
   ```

4. **For branch-specific images:**
   ```yaml
   container-security-scan:staging:
     extends: container-security-scan
     variables:
       CONTAINER_IMAGE_SUFFIX: "staging"
       IMAGE_TAG: "latest"
     rules:
       - if: '$ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_BRANCH == "staging"'
     needs:
       - "Tag staging deployment image"
   ```

See [Container Scanning Configuration Guide](#container-scanning-configuration) for detailed examples.

### Dependency-Track Issues

For Dependency-Track troubleshooting (authentication errors, project not found, network timeouts, empty SBOMs, monorepo configuration), see **[docs/DEPENDENCY_TRACK.md](docs/DEPENDENCY_TRACK.md#troubleshooting)**

### Monorepo: Jobs Running for All Changes

**Problem:** Jobs run even when their project hasn't changed

**GitLab CI Solution:**
```yaml
rules:
  - changes:
      - frontend/**/*  # Must include /**/* suffix
```

**GitHub Actions Solution:**
```yaml
on:
  push:
    paths:
      - 'frontend/**'  # Ensure quotes and correct path
```

### Monorepo: Security Scans Scanning Entire Repository

**Problem:** Security scans analyze all projects instead of just one

**Solution:** Verify PROJECT_PATH (GitLab) or project_path (GitHub) is set:

**GitLab CI:**
```yaml
secrets-detection:frontend:
  extends: secrets-detection
  variables:
    PROJECT_PATH: frontend  # Must be set
```

**GitHub Actions:**
```yaml
secrets:
  uses: platform/devsecops-template/.github/workflows/security-secrets.yml@main
  with:
    project_path: frontend  # Must be set
```

### Monorepo: Build Artifacts Not Found

**Problem:** Downstream jobs can't find artifacts from build stage

**GitLab CI Solution:** Artifact paths must include PROJECT_PATH prefix:
```yaml
artifacts:
  paths:
    - ${PROJECT_PATH:-.}/dist/  # Correct
    # NOT: - dist/  # Wrong - looks at repo root
```

**GitHub Actions Solution:** Ensure working-directory is set in all steps:
```yaml
- name: Build
  working-directory: ${{ inputs.project_path }}
  run: npm run build
```

### Monorepo: Dependencies Not Found During Build

**Problem:** Build fails with "module not found" or "package not found"

**Solution:** Ensure working directory is set before dependency installation:

**GitLab CI:** Jobs should use `cd "${PROJECT_PATH:-.}"` in before_script
**GitHub Actions:** Jobs should use `working-directory: ${{ inputs.project_path }}` in steps

---
