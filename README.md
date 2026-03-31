# GitLab DevSecOps CI/CD Template Library

A comprehensive, enterprise-grade GitLab CI/CD template library implementing DevSecOps best practices with automated security scanning, testing, and deployment.

---

## Table of Contents

- [Quick Start](#quick-start)
- [What's Included](#whats-included)
- [Security Features](#security-features)
- [Testing Your Pipeline Locally](#testing-your-pipeline-locally)
- [Key Features](#key-features)
  - [Flexible Configuration](#flexible-configuration)
  - [GitOps Deployment](#gitops-deployment)
- [Architecture](#architecture)
  - [Pipeline Stages (OWASP SPVS Aligned)](#pipeline-stages-owasp-spvs-aligned)
  - [Template Inheritance](#template-inheritance)
- [Monorepo Support](#monorepo-support)
- [Configuration Variables](#configuration-variables)
  - [Container Scanning Configuration](#container-scanning-configuration)
  - [Dependency-Track Integration](#dependency-track-integration)
  - [AI-Powered Pipeline Reporting](#ai-powered-pipeline-reporting)
- [Handling Security Findings](#handling-security-findings)
  - [Security Policy Modes](#security-policy-modes)
  - [Ignore Files](#ignore-files)
  - [Pipeline Failure Control Variables](#pipeline-failure-control-variables)
- [Integration Guide: Using a `devsecops.yml` Wrapper](#integration-guide-using-a-devsecopsyml-wrapper)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

### Using The Templates

Include the templates you need in your `.gitlab-ci.yml`. Below is a full example with **every available template** and stage described:

```yaml
# your-project/.gitlab-ci.yml
include:
  - project: components/dev-sec-ops
    ref: v1.0.2
    file:
      # ── Core ────────────────────────────────────────────────────────────
      - /templates/gitlab/base.yml              # Pipeline stages, rules, default variables

      # ── Source stage — scan code before building ────────────────────────
      - /templates/gitlab/security/secrets.yml   # Secrets detection (Gitleaks or Trivy)
      - /templates/gitlab/security/dependency.yml # Dependency vulnerability scanning (SCA)
      - /templates/gitlab/security/sast.yml      # Static Application Security Testing

      # ── Build & Test stages — compile and verify ────────────────────────
      # (no template needed — you define build-application and unit-tests yourself)

      # ── Package stage — scan built artifacts ────────────────────────────
      - /templates/gitlab/security/container.yml # Container image scanning (Trivy)
      - /templates/gitlab/security/iac.yml       # Infrastructure-as-Code scanning (Terraform, K8s)
      - /templates/gitlab/security/dtrack.yml    # Upload SBOM to Dependency-Track

      # ── Stage & Deploy — GitOps deployment ──────────────────────────────
      - /templates/gitlab/deploy-staging.yml     # Deploy to staging via GitOps
      - /templates/gitlab/deploy-production.yml  # Deploy to production via GitOps (manual gate)

      # ── Verify stage — test the running application ─────────────────────
      - /templates/gitlab/security/dast.yml      # DAST with OWASP ZAP against staging

      # ── Operate stage — post-deploy checks ──────────────────────────────
      - /templates/gitlab/monitor.yml            # Health-check and availability monitoring

      # ── Report stage — aggregate results ────────────────────────────────
      - /templates/gitlab/report.yml             # Security summary + optional Mattermost notification

      # ── AI stages — AI-powered analysis + Slack ─────────────────────────
      - /templates/gitlab/ai-report.yml          # Gemini/OpenAI analysis → Slack summary

variables:
  DEVSECOPS_PROJECT_LANGUAGE: "node"  # node | python | php | generic

  # Security
  DEVSECOPS_SECURITY_SCANNER: "trivy"          # trivy | specialized
  DEVSECOPS_ENABLE_CONTAINER_SCAN: "true"
  DEVSECOPS_ENABLE_IAC_SCAN: "true"
  DEVSECOPS_ENABLE_DAST: "true"
  DEVSECOPS_STAGING_URL: "https://staging.example.com"

  # Dependency-Track
  DEVSECOPS_ENABLE_DTRACK: "true"
  DEVSECOPS_DTRACK_URL: "https://api.dtrack.example.com"
  DEVSECOPS_DTRACK_API_KEY: "${DEVSECOPS_DTRACK_API_KEY}"            # CI/CD secret

  # GitOps deployment
  GITOPS_REPO: "git@gitlab.example.com:gitops/myapp.git"

  # AI reporting
  DEVSECOPS_ENABLE_AI_REPORT: "true"
  # DEVSECOPS_AI_REPORT_API_KEY and DEVSECOPS_SLACK_WEBHOOK_URL → set as CI/CD secrets

# Your custom jobs here…
```

#### Stage-by-Stage Overview

| Stage | Template(s) | What Happens |
|-------|-------------|--------------|
| **source** | `secrets.yml`, `dependency.yml`, `sast.yml` | Scans source code for leaked secrets, vulnerable dependencies, and code-level security issues — runs **before** anything is built. |
| **build** | _(your job)_ | Build your application. Define a `build-application` job with the image and commands for your stack. |
| **test** | _(your job)_ | Run unit / integration tests. Define a `unit-tests` job that produces JUnit and coverage reports. |
| **package** | `container.yml`, `iac.yml`, `dtrack.yml` | Scans the built container image for OS and app vulnerabilities, validates IaC manifests, and uploads the SBOM to Dependency-Track. |
| **stage** | `deploy-staging.yml` | Deploys to the **staging** environment by updating the image tag in your GitOps repository. |
| **verify** | `dast.yml` | Runs OWASP ZAP against the staging URL to detect runtime vulnerabilities (SQL injection, XSS, CSRF, etc.). |
| **deploy** | `deploy-production.yml` | Deploys to **production** via GitOps — requires **manual approval** on the default branch. |
| **operate** | `monitor.yml` | Runs health-check requests against the deployed application and reports availability. |
| **report** | `report.yml` | Aggregates all security scan results into a single summary and optionally notifies Mattermost. |
| **ai-analysis** | `ai-report.yml` | Sends each scan report to Gemini or OpenAI for individual analysis. |
| **ai-summary** | `ai-report.yml` | Consolidates the AI analyses into one actionable summary and posts it to Slack. |

> You don't need all templates — start with `base.yml` plus the security scans you want and add more as your pipeline matures. See the [examples/](examples/) directory for minimal per-language setups.

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

# Run all security scans (node, python, php)
make test

# Run a single language
make test-node
make test-python
make test-php

# Validate GitHub Actions YAML
make validate

# Test AI reporting pipeline
make ai-report-test

# Or call dagger directly
cd dagger
dagger call test --source=../examples/node --language=node
dagger call secrets-detection --source=../examples/node
dagger call sast-scanning --source=../examples/node
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
  DEVSECOPS_SECURITY_SCANNER: "trivy"
  # Or use specialized tools for comprehensive analysis
  DEVSECOPS_SECURITY_SCANNER: "specialized"
  DEVSECOPS_SAST_TOOL: "semgrep"  # More thorough SAST

  # Feature toggles
  DEVSECOPS_ENABLE_SECRETS: "true"
  DEVSECOPS_ENABLE_DEPENDENCY_SCAN: "true"
  DEVSECOPS_ENABLE_SAST: "true"
  DEVSECOPS_ENABLE_CONTAINER_SCAN: "false"
  DEVSECOPS_ENABLE_IAC_SCAN: "false"
  DEVSECOPS_ENABLE_DAST: "false"

  # Security policy
  DEVSECOPS_SECURITY_POLICY: "strict"  # or "permissive"
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

## Architecture

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
│ AI-Powered Reporting (Optional - DEVSECOPS_ENABLE_AI_REPORT: "true")             │
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
    DEVSECOPS_NODE_VERSION: "20"
```

---

## Monorepo Support

The DevSecOps templates support monorepos with multiple projects in a single repository. Each project can run CI stages independently with automatic change detection.

### Quick Example

**GitLab CI (.gitlab-ci.yml):**
```yaml
include:
  - project: components/dev-sec-ops
    ref: v1.0.2
    file:
      - /templates/gitlab/base.yml
      - /templates/build.yml
      - /templates/gitlab/security/secrets.yml

# Frontend project - only runs when frontend/ changes
build:frontend:
  extends: .build:node
  variables:
    DEVSECOPS_PROJECT_PATH: frontend
    DEVSECOPS_NODE_VERSION: "20"
  rules:
    - changes:
        - frontend/**/*

secrets-detection:frontend:
  extends: secrets-detection
  variables:
    DEVSECOPS_PROJECT_PATH: frontend
  rules:
    - changes:
        - frontend/**/*

# Backend project - only runs when backend/ changes
build:backend:
  extends: .build:python
  variables:
    DEVSECOPS_PROJECT_PATH: backend
    DEVSECOPS_PYTHON_VERSION: "3.12"
  rules:
    - changes:
        - backend/**/*

secrets-detection:backend:
  extends: secrets-detection
  variables:
    DEVSECOPS_PROJECT_PATH: backend
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
    uses: components/dev-sec-ops/.github/workflows/build-node.yml@v1.0.2
    with:
      project_path: frontend
      node_version: '20'

  secrets:
    uses: components/dev-sec-ops/.github/workflows/security-secrets.yml@v1.0.2
    with:
      project_path: frontend
    permissions:
      contents: read
```

### Key Features

- **Per-project CI/CD:** Each project runs its own build, test, and security scans
- **Change detection:** Jobs only run when project files change
- **Independent deployment:** Projects can be deployed separately
- **DEVSECOPS_PROJECT_PATH variable:** Scopes all operations to project directory
- **Efficient:** Reduces CI minutes by only testing affected projects

### How It Works

**GitLab CI:**
- Use `DEVSECOPS_PROJECT_PATH` variable to specify project directory
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
DEVSECOPS_PROJECT_LANGUAGE: "node"  # node|python|php|generic
DEVSECOPS_NODE_VERSION: "20"
DEVSECOPS_PYTHON_VERSION: "3.12"
DEVSECOPS_PHP_VERSION: "8.3"
DEVSECOPS_PACKAGE_MANAGER: "npm"  # npm|yarn|pnpm
```

### Security Scanning
```yaml
# Security tool selection
DEVSECOPS_SECURITY_SCANNER: "trivy"  # trivy|specialized
# trivy: Use Trivy for unified security scanning (recommended)
# specialized: Use specialized tools (Gitleaks, Semgrep, etc.)

# SAST tool selection (when using specialized tools or mixed approach)
DEVSECOPS_SAST_TOOL: "semgrep"  # semgrep|trivy
# semgrep: Comprehensive SAST with Semgrep (recommended for thorough analysis)
# trivy: Basic SAST with Trivy misconfiguration detection

# Feature toggles
DEVSECOPS_ENABLE_SECRETS: "true"
DEVSECOPS_ENABLE_DEPENDENCY_SCAN: "true"
DEVSECOPS_ENABLE_SAST: "true"
DEVSECOPS_ENABLE_CONTAINER_SCAN: "false"
DEVSECOPS_ENABLE_IAC_SCAN: "false"
DEVSECOPS_ENABLE_DAST: "false"
DEVSECOPS_ENABLE_DTRACK: "false"  # Dependency-Track SBOM upload
DEVSECOPS_ENABLE_AI_REPORT: "false"  # AI pipeline analysis + Slack summary

# Security policy
DEVSECOPS_SECURITY_POLICY: "strict"  # strict|permissive
```

### Trivy Configuration
```yaml
DEVSECOPS_TRIVY_SEVERITY: "CRITICAL,HIGH"
DEVSECOPS_TRIVY_EXIT_CODE: "1"
DEVSECOPS_IMAGE_NAME: "${CI_REGISTRY_IMAGE}"
DEVSECOPS_IMAGE_TAG: "${CI_COMMIT_SHA}"
```

### Container Scanning Configuration

The container scanning template scans Docker images for vulnerabilities using Trivy. It supports various deployment patterns and custom image naming schemes.

**📖 Full Documentation:** [docs/CONTAINER_SCANNING.md](docs/CONTAINER_SCANNING.md)

#### Basic Configuration

```yaml
include:
  - project: components/dev-sec-ops
    file: /templates/gitlab/security/container.yml

variables:
  DEVSECOPS_ENABLE_CONTAINER_SCAN: "true"
  DEVSECOPS_IMAGE_NAME: "${CI_REGISTRY_IMAGE}"
  DEVSECOPS_IMAGE_TAG: "latest"
```

#### Advanced Configuration: Sub-Images

For projects that build multiple images (e.g., `myapp/staging:latest`, `myapp/prod:latest`):

```yaml
variables:
  DEVSECOPS_IMAGE_NAME: "${CI_REGISTRY_IMAGE}"
  DEVSECOPS_IMAGE_TAG: "latest"
  DEVSECOPS_CONTAINER_IMAGE_SUFFIX: "staging"  # Scans CI_REGISTRY_IMAGE/staging:latest

# Or override per-branch:
container-security-scan:staging:
  extends: container-security-scan
  variables:
    DEVSECOPS_CONTAINER_IMAGE_SUFFIX: "staging"
  rules:
    - if: '$DEVSECOPS_ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_BRANCH == "staging"'
  needs:
    - "build-staging-image"

container-security-scan:prod:
  extends: container-security-scan
  variables:
    DEVSECOPS_CONTAINER_IMAGE_SUFFIX: "prod"
  rules:
    - if: '$DEVSECOPS_ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_BRANCH == "prod"'
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
  DEVSECOPS_IMAGE_NAME: "docker.io/myorg/myapp"  # External registry
  DEVSECOPS_IMAGE_TAG: "${CI_COMMIT_SHORT_SHA}"
  
container-security-scan:
  needs:
    - "push-to-dockerhub"
```

#### Available Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DEVSECOPS_IMAGE_NAME` | `${CI_REGISTRY_IMAGE}` | Full image repository path |
| `DEVSECOPS_IMAGE_TAG` | `${CI_COMMIT_SHORT_SHA}` | Image tag to scan |
| `DEVSECOPS_CONTAINER_IMAGE_SUFFIX` | _(empty)_ | Sub-path for image (e.g., `staging`, `prod`) |
| `DEVSECOPS_TRIVY_SEVERITY` | `UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL` | Severity levels to report |
| `DEVSECOPS_TRIVY_EXIT_CODE` | `0` | Exit code for failures (set to `1` to fail pipeline) |
| `DEVSECOPS_TRIVY_NON_SSL` | `false` | Set to `true` for insecure registries |

#### Image Resolution

The template builds the image reference as follows:

1. Base: `DEVSECOPS_IMAGE_NAME` (defaults to `CI_REGISTRY_IMAGE`)
2. If `DEVSECOPS_CONTAINER_IMAGE_SUFFIX` is set: `DEVSECOPS_IMAGE_NAME/DEVSECOPS_CONTAINER_IMAGE_SUFFIX`
3. Tag: `:DEVSECOPS_IMAGE_TAG`

**Examples:**
- `DEVSECOPS_IMAGE_NAME=registry/project`, `DEVSECOPS_IMAGE_TAG=v1.0` → `registry/project:v1.0`
- `DEVSECOPS_IMAGE_NAME=registry/project`, `DEVSECOPS_CONTAINER_IMAGE_SUFFIX=staging`, `DEVSECOPS_IMAGE_TAG=latest` → `registry/project/staging:latest`

### Dependency-Track Integration

Upload Software Bill of Materials (SBOM) to Dependency-Track for centralized dependency and vulnerability tracking. Supports both single-project and monorepo configurations with automatic project creation.

**📖 Full Documentation:** [docs/DEPENDENCY_TRACK.md](docs/DEPENDENCY_TRACK.md)

#### Quick Start (GitLab CI)

```yaml
include:
  - project: components/dev-sec-ops
    ref: v1.0.2
    file:
      - /templates/gitlab/base.yml
      - /templates/gitlab/security/dtrack.yml

variables:
  DEVSECOPS_ENABLE_DTRACK: "true"
  DEVSECOPS_DTRACK_URL: "https://api.dtrack.example.com"
  DEVSECOPS_DTRACK_API_KEY: "${DEVSECOPS_DTRACK_API_KEY}"  # Set as CI/CD secret
```

#### Quick Start (GitHub Actions)

```yaml
jobs:
  dtrack:
    uses: components/dev-sec-ops/.github/workflows/security-dtrack.yml@v1.0.2
    with:
      dtrack_url: "https://api.dtrack.example.com"
    secrets:
      dtrack_api_key: ${{ secrets.DEVSECOPS_DTRACK_API_KEY }}
```

#### Key Variables

| Variable (GitLab) / Input (GitHub) | Description |
|------------------------------------|-------------|
| `DEVSECOPS_DTRACK_URL` / `dtrack_url` | DTrack base URL (required) |
| `DEVSECOPS_DTRACK_API_KEY` / `dtrack_api_key` | DTrack API key (required, secret) |
| `DEVSECOPS_PROJECT_PATH` / `project_path` | Monorepo subproject path (e.g., `frontend`) |
| `DEVSECOPS_DTRACK_PROJECT_UUID` / `dtrack_project_uuid` | Explicit project UUID (optional) |

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
DEVSECOPS_STAGING_URL: "https://staging.example.com"
DEVSECOPS_PRODUCTION_URL: "https://production.example.com"
DEVSECOPS_HEALTHCHECK_URL: "${DEVSECOPS_STAGING_URL}/health"
```

### AI-Powered Pipeline Reporting

Automated pipeline analysis using Google Gemini with Slack notifications. Analyzes all non-deployment stage outputs and sends a single, actionable summary to Slack.

**📖 Full Documentation:** [docs/AI_REPORTING.md](docs/AI_REPORTING.md)

#### Quick Start (GitLab CI)

```yaml
include:
  - project: components/dev-sec-ops
    file:
      - /templates/gitlab/base.yml
      - /templates/gitlab/ai-report.yml

variables:
  DEVSECOPS_ENABLE_AI_REPORT: "true"
  # DEVSECOPS_AI_REPORT_PROVIDER: "openai"  # Optional: switch to OpenAI (default: "gemini")
  # Set DEVSECOPS_AI_REPORT_API_KEY as CI/CD secret (Gemini or OpenAI key)
  # Set DEVSECOPS_SLACK_WEBHOOK_URL as CI/CD secret (optional)
```

#### Quick Start (GitHub Actions)

```yaml
jobs:
  ai-report:
    needs: [build, test, sast, dependency-scan]
    if: always()
    uses: ./.github/workflows/ai-report.yml
    secrets:
      ai_api_key: ${{ secrets.DEVSECOPS_AI_REPORT_API_KEY }}
      slack_webhook_url: ${{ secrets.DEVSECOPS_SLACK_WEBHOOK_URL }}
```

#### Key Variables

| Variable / Secret | Description |
|-------------------|-------------|
| `DEVSECOPS_ENABLE_AI_REPORT` | Feature toggle (default: `"false"`) |
| `DEVSECOPS_AI_REPORT_API_KEY` | API key for Gemini or OpenAI (CI/CD secret) |
| `DEVSECOPS_AI_REPORT_PROVIDER` | `"gemini"` (default) or `"openai"` |
| `DEVSECOPS_SLACK_WEBHOOK_URL` | Slack incoming webhook URL (CI/CD secret, optional) |
| `DEVSECOPS_AI_REPORT_MODEL` | Model override (default: auto per provider) |

### Notifications
```yaml
DEVSECOPS_SLACK_WEBHOOK_URL: "https://mattermost.com/hooks/xxx"
```

---

## Handling Security Findings

### Security Policy Modes

The templates support two security policy modes that control how pipelines react to findings:

| Mode | Behavior | When to Use |
|------|----------|-------------|
| `strict` (default) | Pipeline **fails** on any security finding | Production branches, merge requests |
| `permissive` | Security jobs **warn** but do not block the pipeline on feature branches; still strict on main/MR | Early adoption, initial rollout |

```yaml
variables:
  DEVSECOPS_SECURITY_POLICY: "permissive"  # or "strict"
```

### Ignore Files

Suppress known false positives by adding ignore files to your repository root:

| Tool | Ignore File | Example Entry |
|------|-------------|---------------|
| Trivy | `.trivyignore` | `CVE-2024-12345` |
| Gitleaks | `.gitleaksignore` | `test-fixtures/fake-secret.txt` |
| Semgrep | `.semgrepignore` | `tests/` |

### Pipeline Failure Control Variables

Fine-tune which findings block the pipeline:

| Variable | Default | Description |
|----------|---------|-------------|
| `DEVSECOPS_TRIVY_EXIT_CODE` | `"1"` | `"0"` = warn only, `"1"` = fail on findings |
| `DEVSECOPS_TRIVY_SEVERITY` | `"CRITICAL,HIGH"` | Severity levels to report (`UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL`) |
| `DEVSECOPS_TRIVY_IGNORE_UNFIXED` | `"false"` | `"true"` = skip vulnerabilities without an available fix |

### Recommended Approach

1. **Start permissive** — set `DEVSECOPS_SECURITY_POLICY: "permissive"` so pipelines keep running while you assess the findings
2. **Triage** — review findings in the GitLab Security Dashboard, add genuine false positives to ignore files
3. **Go strict** — once the backlog is clean, switch to `DEVSECOPS_SECURITY_POLICY: "strict"` and set `DEVSECOPS_TRIVY_EXIT_CODE: "1"` to enforce a clean gate

---

## Integration Guide: Using a `devsecops.yml` Wrapper

The recommended way to add DevSecOps scanning to a project is through a local `devsecops.yml` wrapper file. This keeps the template version pinned in one place and your `.gitlab-ci.yml` focused on your app.

### Step 1 — Create `devsecops.yml` in your repository root

Pick the security templates you need. The base set covers secrets, dependencies, and SAST:

```yaml
# devsecops.yml
include:
  - project: components/dev-sec-ops
    ref: v1.0.2
    file:
      - /templates/gitlab/base.yml
      - /templates/gitlab/security/secrets.yml
      - /templates/gitlab/security/dependency.yml
      - /templates/gitlab/security/sast.yml
      # Uncomment as needed:
      # - /templates/gitlab/security/container.yml  # Container image scanning
      # - /templates/gitlab/security/iac.yml        # IaC scanning (Terraform, K8s)
      # - /templates/gitlab/security/dtrack.yml     # SBOM upload to Dependency-Track
      # - /templates/gitlab/ai-report.yml           # AI analysis + Slack summary
      # - /templates/gitlab/deploy-staging.yml      # GitOps staging deployment
      # - /templates/gitlab/deploy-production.yml   # GitOps production deployment
      # - /templates/gitlab/monitor.yml             # Post-deploy health checks
```

### Step 2 — Reference the wrapper from `.gitlab-ci.yml`

```yaml
# .gitlab-ci.yml
include:
  - local: devsecops.yml

variables:
  DEVSECOPS_PROJECT_LANGUAGE: "node"           # node | python | php | generic
  DEVSECOPS_NODE_VERSION: "20"
  DEVSECOPS_PACKAGE_MANAGER: "pnpm"   # npm | yarn | pnpm

build-application:
  stage: build
  image: node:${DEVSECOPS_NODE_VERSION}-alpine
  script:
    - corepack enable
    - pnpm install --frozen-lockfile
    - pnpm build
  artifacts:
    paths:
      - .output/

unit-tests:
  stage: test
  image: node:${DEVSECOPS_NODE_VERSION}-alpine
  script:
    - corepack enable
    - pnpm install --frozen-lockfile
    - pnpm test -- --ci
```

### Step 3 — Add CI/CD secrets (if using optional features)

In **Settings > CI/CD > Variables**, add any secrets required by the templates you enabled:

| Variable | When Needed |
|----------|-------------|
| `DEVSECOPS_DTRACK_API_KEY` | Dependency-Track integration |
| `DEVSECOPS_AI_REPORT_API_KEY` | AI pipeline analysis (Gemini or OpenAI key) |
| `DEVSECOPS_SLACK_WEBHOOK_URL` | Slack notifications for AI reports |
| `GITOPS_SSH_KEY` | GitOps deployments |

### Step 4 — Commit and push

```bash
git add devsecops.yml .gitlab-ci.yml
git commit -m "Add DevSecOps security scanning"
git push
```

The pipeline will now run secrets detection, dependency scanning, and SAST on every merge request and push to the default branch.

### Upgrading

To upgrade all templates, change the `ref:` in `devsecops.yml` — that's the only file you need to touch:

```diff
  - project: components/dev-sec-ops
-   ref: v1.0.2
+   ref: v1.1.0
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
1. Verify `DEVSECOPS_STAGING_URL` is correct
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

1. **Set correct DEVSECOPS_IMAGE_TAG variable:**
   ```yaml
   variables:
     DEVSECOPS_IMAGE_TAG: "latest"  # Must match your actual build tag
   ```

2. **For sub-images (e.g., staging/prod):**
   ```yaml
   variables:
     DEVSECOPS_CONTAINER_IMAGE_SUFFIX: "staging"  # Scans CI_REGISTRY_IMAGE/staging:TAG
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
       DEVSECOPS_CONTAINER_IMAGE_SUFFIX: "staging"
       DEVSECOPS_IMAGE_TAG: "latest"
     rules:
       - if: '$DEVSECOPS_ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_BRANCH == "staging"'
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

**Solution:** Verify DEVSECOPS_PROJECT_PATH (GitLab) or project_path (GitHub) is set:

**GitLab CI:**
```yaml
secrets-detection:frontend:
  extends: secrets-detection
  variables:
    DEVSECOPS_PROJECT_PATH: frontend  # Must be set
```

**GitHub Actions:**
```yaml
secrets:
  uses: components/dev-sec-ops/.github/workflows/security-secrets.yml@main
  with:
    project_path: frontend  # Must be set
```

### Monorepo: Build Artifacts Not Found

**Problem:** Downstream jobs can't find artifacts from build stage

**GitLab CI Solution:** Artifact paths must include DEVSECOPS_PROJECT_PATH prefix:
```yaml
artifacts:
  paths:
    - ${DEVSECOPS_PROJECT_PATH:-.}/dist/  # Correct
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

**GitLab CI:** Jobs should use `cd "${DEVSECOPS_PROJECT_PATH:-.}"` in before_script
**GitHub Actions:** Jobs should use `working-directory: ${{ inputs.project_path }}` in steps

---
