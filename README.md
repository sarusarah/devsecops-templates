# GitLab DevSecOps CI/CD Template Library

A comprehensive, enterprise-grade GitLab CI/CD template library implementing DevSecOps best practices with automated security scanning, testing, and deployment.

---

## Quick Start

### For Teams Using These Templates

Include the templates you need in your `.gitlab-ci.yml` and enable desired security scans:
```yaml
# your-project/.gitlab-ci.yml
include:
  - project: platform/devsecops-template
    ref: v1.0.1
    file:
      - /templates/base.yml
      - /templates/security/secrets.yml
      - /templates/security/dependency.yml
      - /templates/security/sast.yml

variables:
  LANGUAGE: "node"  # or python, php
  ENABLE_DAST: "true"
  STAGING_URL: "https://staging.example.com"

# Your custom jobs here...
```

### For Template Developers

```bash
# Clone repository

# Test locally with Dagger (recommended)
dagger call test --source=./examples/node --language=node

# Or use gitlab-ci-local
./test-local.sh
```

---

## What's Included

### Security Scanning Templates (`templates/security/`)

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

### Pipeline Templates (`templates/`)
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

**See [TESTING.md](./TESTING.md) for complete testing guide**

---

## Documentation

| Document | Description |
|----------|-------------|
| [TESTING.md](./TESTING.md) | Complete guide to local testing |
| [dagger/README.md](./dagger/README.md) | Dagger module usage guide |

---

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

### Template Inheritance

```yaml
# Base template provides stages and rules
include:
  - local: /templates/base.yml

# Extend with specific security scans
include:
  - local: /templates/security/secrets.yml
  - local: /templates/security/sast.yml

# Use pre-built job templates
build-app:
  extends: .build:node
  variables:
    NODE_VERSION: "20"
```

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

---
