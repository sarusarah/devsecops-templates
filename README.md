# GitLab DevSecOps CI/CD Template Library

**Version:** 1.0.1
**Status:** ✅ Production Ready
**Last Updated:** 2025-12-12

A comprehensive, enterprise-grade GitLab CI/CD template library implementing DevSecOps best practices with automated security scanning, testing, and deployment.

---

## 🚀 Quick Start

### For Teams Using These Templates

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
cd /home/sarah/Projects/Liip/DevSecOps

# Test locally with Dagger (recommended)
dagger call test --source=./examples/node --language=node

# Or use gitlab-ci-local
./test-local.sh
```

---

## 📋 What's Included

### Security Scanning Templates (`templates/security/`)
- **Secrets Detection** - Gitleaks for detecting committed credentials
- **Dependency Scanning** - npm audit, pip-audit, composer audit
- **SAST** - Semgrep static code analysis
- **DAST** - OWASP ZAP dynamic scanning
- **Container Scanning** - Trivy for container vulnerabilities
- **IaC Security** - Kubernetes manifest validation with kubeconform, kube-score, Polaris

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

## 🔒 Security Features

- ✅ **Secrets detection** in preflight stage (blocks by default)
- ✅ **Dependency vulnerability scanning** for all package managers
- ✅ **Static code analysis** (SAST) with Semgrep
- ✅ **Dynamic security testing** (DAST) with OWASP ZAP
- ✅ **Container image scanning** with Trivy
- ✅ **Infrastructure as Code** validation
- ✅ **Security policy enforcement** (strict/permissive modes)
- ✅ **GitLab Security Dashboard** integration
- ✅ **Audit trail** with 7-day artifact retention

---

## 🧪 Testing Your Pipeline Locally

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

**Why Dagger?**
- ⚡ 10x faster than GitLab CI
- 🎯 Same containers as production
- 💰 No CI minutes consumed
- 🔄 Intelligent caching

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

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [TESTING.md](./TESTING.md) | Complete guide to local testing |
| [FIXES_APPLIED.md](./FIXES_APPLIED.md) | All fixes and improvements in v1.0.1 |
| [dagger/README.md](./dagger/README.md) | Dagger module usage guide |

---

## 🎯 Key Features

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
```yaml
variables:
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

### GitOps Deployment
```yaml
variables:
  GITOPS_REPO: "git@gitlab.com:gitops/myapp.git"
  GITOPS_PATH: "values.yaml"
  GITOPS_IMAGE_TAG_YQ_PATH: ".image.tag"
```

---

## 🏗️ Architecture

### Pipeline Stages

```
1. preflight          → Secrets detection, dependency scanning
2. build              → Application build with artifacts
3. test               → Unit tests with coverage
4. security-scan      → SAST with Semgrep
5. security-analysis  → Container & IaC scanning
6. deploy-staging     → GitOps deployment to staging
7. validate           → DAST against staging
8. deploy-production  → Manual deployment to production
9. monitor            → Health checks
10. report            → Aggregate security findings
```

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

## 🔧 Configuration Variables

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
ENABLE_SECRETS: "true"
ENABLE_DEPENDENCY_SCAN: "true"
ENABLE_SAST: "true"
ENABLE_CONTAINER_SCAN: "false"
ENABLE_IAC_SCAN: "false"
ENABLE_DAST: "false"
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

## 🐛 Troubleshooting

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

## 📊 Changelog

### v1.0.1 (2025-12-12) - Current

**Fixed:**
- ✅ Gitleaks report generation
- ✅ DAST policy enforcement
- ✅ PHP test dependencies
- ✅ Container registry authentication
- ✅ Mattermost webhook
- ✅ Missing GitOps templates

**Added:**
- ✅ Dagger testing module
- ✅ gitlab-ci-local configuration
- ✅ Comprehensive documentation
- ✅ GitLab Security Dashboard integration
- ✅ Pinned all tool versions

**See [FIXES_APPLIED.md](./FIXES_APPLIED.md) for complete details**

### v1.0.0 (Initial Release)
- Base templates
- Security scanning templates
- Language-specific builds
- Examples for Node.js, Python, PHP

---

## 🤝 Contributing

### Template Developers

1. Make changes to templates
2. Test locally: `dagger call test --source=./examples/node --language=node`
3. Validate YAML: `dagger call validate-yaml --yaml-file=.gitlab-ci.yml`
4. Run all examples: `./test-local.sh` (option 5)
5. Update documentation
6. Create merge request

### Security Researchers

Found a security issue? Please email security@company.com

---

## 📞 Support

- **Slack:** #devsecops-support
- **Email:** devsecops-team@company.com
- **Docs:** https://docs.company.com/devsecops
- **Issues:** https://gitlab.com/platform/devsecops-template/-/issues

---

## 📄 License

MIT License - See LICENSE file for details

---

## ⭐ Quick Links

- [Testing Guide](./TESTING.md)
- [Fix Report](./FIXES_APPLIED.md)
- [Dagger Module](./dagger/README.md)
- [Examples](./examples/)
- [Security Templates](./templates/security/)

---

**Built with ❤️ by the DevSecOps Team**
