# Dependency-Track Integration Guide

Upload Software Bill of Materials (SBOM) to Dependency-Track for centralized dependency and vulnerability tracking across all projects.

---

## Table of Contents
- [Quick Start](#quick-start)
- [Configuration Modes](#configuration-modes)
- [GitLab CI Examples](#gitlab-ci-examples)
- [GitHub Actions Examples](#github-actions-examples)
- [Monorepo Setup](#monorepo-setup)
- [Variables Reference](#variables-reference)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

### GitLab CI

```yaml
# .gitlab-ci.yml
include:
  - project: platform/devsecops-template
    ref: v1.0.1
    file:
      - /templates/base.yml
      - /templates/security/dtrack.yml

variables:
  ENABLE_DTRACK: "true"
  DTRACK_URL: "https://api.dtrack.example.com"
  DTRACK_API_KEY: "${DTRACK_API_KEY}"  # Set as CI/CD secret
```

### GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI Pipeline

on:
  push:
    branches: [main]

jobs:
  dtrack:
    uses: platform/devsecops-template/.github/workflows/security-dtrack.yml@v1.0.1
    with:
      dtrack_url: "https://api.dtrack.example.com"
    secrets:
      dtrack_api_key: ${{ secrets.DTRACK_API_KEY }}
```

---

## Configuration Modes

### Mode 1: Auto-Create (Recommended)

Automatically creates Dependency-Track projects using repository path and version.

**GitLab CI:**
```yaml
variables:
  ENABLE_DTRACK: "true"
  DTRACK_URL: "https://api.dtrack.example.com"
  DTRACK_API_KEY: "${DTRACK_API_KEY}"
  # Project auto-named: "platform/my-repo"
  # Version auto-set: CI_COMMIT_TAG or CI_COMMIT_SHORT_SHA
```

**GitHub Actions:**
```yaml
with:
  dtrack_url: "https://api.dtrack.example.com"
  # Project auto-named: "org/my-repo"
  # Version auto-set: branch/tag name (GITHUB_REF_NAME)
secrets:
  dtrack_api_key: ${{ secrets.DTRACK_API_KEY }}
```

### Mode 2: Explicit Project UUID

Use existing Dependency-Track project by UUID (takes precedence over auto-create).

**GitLab CI:**
```yaml
variables:
  ENABLE_DTRACK: "true"
  DTRACK_URL: "https://api.dtrack.example.com"
  DTRACK_API_KEY: "${DTRACK_API_KEY}"
  DTRACK_PROJECT_UUID: "abc-123-def-456"  # Your project UUID from DTrack
```

**GitHub Actions:**
```yaml
with:
  dtrack_url: "https://api.dtrack.example.com"
  dtrack_project_uuid: "abc-123-def-456"
secrets:
  dtrack_api_key: ${{ secrets.DTRACK_API_KEY }}
```

### Mode 3: Custom Project Name and Version

Override auto-generated names with custom values.

**GitLab CI:**
```yaml
variables:
  ENABLE_DTRACK: "true"
  DTRACK_URL: "https://api.dtrack.example.com"
  DTRACK_API_KEY: "${DTRACK_API_KEY}"
  DTRACK_PROJECT_NAME: "my-custom-project-name"
  DTRACK_PROJECT_VERSION: "v1.2.3"
```

**GitHub Actions:**
```yaml
with:
  dtrack_url: "https://api.dtrack.example.com"
  dtrack_project_name: "my-custom-project-name"
  dtrack_project_version: "v1.2.3"
secrets:
  dtrack_api_key: ${{ secrets.DTRACK_API_KEY }}
```

---

## GitLab CI Examples

### Basic Single Project

```yaml
include:
  - project: platform/devsecops-template
    ref: v1.0.1
    file:
      - /templates/base.yml
      - /templates/security/dtrack.yml

variables:
  ENABLE_DTRACK: "true"
  DTRACK_URL: "https://api.dtrack.example.com"
  DTRACK_API_KEY: "${DTRACK_API_KEY}"
```

**Result:**
- Project name: `platform/my-repo`
- Version: Latest tag or commit SHA
- SBOM: Entire repository

### With Custom Version Tags

```yaml
variables:
  ENABLE_DTRACK: "true"
  DTRACK_URL: "https://api.dtrack.example.com"
  DTRACK_API_KEY: "${DTRACK_API_KEY}"
  DTRACK_PROJECT_VERSION: "${CI_COMMIT_TAG}"

# Only upload on tagged releases
dtrack-upload:
  rules:
    - if: '$ENABLE_DTRACK == "true" && $CI_COMMIT_TAG'
```

### Scheduled Dependency Audits

```yaml
# Run daily SBOM upload for main branch
dtrack-upload:
  rules:
    - if: '$ENABLE_DTRACK == "true" && $CI_PIPELINE_SOURCE == "schedule" && $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH'
```

---

## GitHub Actions Examples

### Basic Single Project

```yaml
name: Security Scanning

on:
  push:
    branches: [main, develop]

jobs:
  dtrack:
    uses: platform/devsecops-template/.github/workflows/security-dtrack.yml@v1.0.1
    with:
      dtrack_url: ${{ vars.DTRACK_URL }}
    secrets:
      dtrack_api_key: ${{ secrets.DTRACK_API_KEY }}
```

**Result:**
- Project name: `org/my-repo`
- Version: Branch/tag name
- SBOM: Entire repository

### Release Tags Only

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  dtrack:
    uses: platform/devsecops-template/.github/workflows/security-dtrack.yml@v1.0.1
    with:
      dtrack_url: ${{ vars.DTRACK_URL }}
      dtrack_project_version: ${{ github.ref_name }}
    secrets:
      dtrack_api_key: ${{ secrets.DTRACK_API_KEY }}
```

### Scheduled Dependency Audits

```yaml
name: Daily Dependency Audit

on:
  schedule:
    - cron: '0 2 * * *'  # 2 AM daily

jobs:
  dtrack:
    uses: platform/devsecops-template/.github/workflows/security-dtrack.yml@v1.0.1
    with:
      dtrack_url: ${{ vars.DTRACK_URL }}
    secrets:
      dtrack_api_key: ${{ secrets.DTRACK_API_KEY }}
```

---

## Monorepo Setup

### GitLab CI Monorepo

```yaml
include:
  - project: platform/devsecops-template
    ref: v1.0.1
    file:
      - /templates/base.yml
      - /templates/security/dtrack.yml

variables:
  DTRACK_URL: "https://api.dtrack.example.com"
  DTRACK_API_KEY: "${DTRACK_API_KEY}"

# Frontend component
dtrack-upload:frontend:
  extends: dtrack-upload
  variables:
    PROJECT_PATH: frontend
    ENABLE_DTRACK: "true"
  rules:
    - if: '$ENABLE_DTRACK == "true"'
      changes:
        - frontend/**/*

# Backend component
dtrack-upload:backend:
  extends: dtrack-upload
  variables:
    PROJECT_PATH: backend
    ENABLE_DTRACK: "true"
  rules:
    - if: '$ENABLE_DTRACK == "true"'
      changes:
        - backend/**/*

# Mobile app
dtrack-upload:mobile:
  extends: dtrack-upload
  variables:
    PROJECT_PATH: mobile-app
    ENABLE_DTRACK: "true"
  rules:
    - if: '$ENABLE_DTRACK == "true"'
      changes:
        - mobile-app/**/*
```

**Result:**
- Frontend: `platform/my-repo/frontend`
- Backend: `platform/my-repo/backend`
- Mobile: `platform/my-repo/mobile-app`

### GitHub Actions Monorepo

**Frontend workflow** (`.github/workflows/frontend.yml`):
```yaml
name: Frontend CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'frontend/**'

jobs:
  dtrack:
    uses: platform/devsecops-template/.github/workflows/security-dtrack.yml@v1.0.1
    with:
      project_path: frontend
      dtrack_url: ${{ vars.DTRACK_URL }}
    secrets:
      dtrack_api_key: ${{ secrets.DTRACK_API_KEY }}
```

**Backend workflow** (`.github/workflows/backend.yml`):
```yaml
name: Backend CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'backend/**'

jobs:
  dtrack:
    uses: platform/devsecops-template/.github/workflows/security-dtrack.yml@v1.0.1
    with:
      project_path: backend
      dtrack_url: ${{ vars.DTRACK_URL }}
    secrets:
      dtrack_api_key: ${{ secrets.DTRACK_API_KEY }}
```

**Result:**
- Frontend: `org/my-repo/frontend`
- Backend: `org/my-repo/backend`

---

## Variables Reference

### GitLab CI Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ENABLE_DTRACK` | No | `"false"` | Enable Dependency-Track upload |
| `DTRACK_URL` | **Yes** | - | DTrack base URL (e.g., `https://api.dtrack.example.com`) |
| `DTRACK_API_KEY` | **Yes** | - | DTrack API key (set as CI/CD secret) |
| `DTRACK_PROJECT_UUID` | No | - | Explicit project UUID (takes precedence) |
| `DTRACK_PROJECT_NAME` | No | `${CI_PROJECT_PATH}[/${PROJECT_PATH}]` | Override project name |
| `DTRACK_PROJECT_VERSION` | No | `${CI_COMMIT_TAG}` or `${CI_COMMIT_SHORT_SHA}` | Project version |
| `PROJECT_PATH` | No | `"."` | Subproject path for monorepo |

### GitHub Actions Inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `dtrack_url` | **Yes** | - | DTrack base URL |
| `project_path` | No | `"."` | Subproject path for monorepo |
| `dtrack_project_uuid` | No | - | Explicit project UUID (takes precedence) |
| `dtrack_project_name` | No | `${GITHUB_REPOSITORY}[/${project_path}]` | Override project name |
| `dtrack_project_version` | No | `${GITHUB_REF_NAME}` | Project version |

### GitHub Actions Secrets

| Secret | Required | Description |
|--------|----------|-------------|
| `dtrack_api_key` | **Yes** | DTrack API key |

---

## Troubleshooting

### Authentication Error (HTTP 401)

**Symptoms:**
```
✗ Failed to upload SBOM
HTTP Status Code: 401
```

**Solutions:**
1. **GitLab CI:** Set `DTRACK_API_KEY` in Settings → CI/CD → Variables
   - ✓ Check "Mask variable"
   - ✓ Uncheck "Protect variable" (unless limiting to protected branches)
2. **GitHub Actions:** Set `DTRACK_API_KEY` in Settings → Secrets → Actions
3. Verify API key in Dependency-Track:
   - Admin → Access Management → Teams → Automation → API Keys
   - Ensure key has `BOM_UPLOAD` permission
4. Test API key manually:
   ```bash
   curl -H "X-Api-Key: YOUR_KEY" https://api.dtrack.example.com/api/version
   ```

### Project Not Found (HTTP 404)

**Symptoms:**
```
✗ Failed to upload SBOM
HTTP Status Code: 404
{"error": "Project not found"}
```

**Solutions:**
1. Verify `DTRACK_PROJECT_UUID` is correct (copy from DTrack UI)
2. Check project exists: DTrack → Projects → Search
3. Use auto-create mode instead:
   - Remove `DTRACK_PROJECT_UUID` variable
   - Let DTrack auto-create project with name+version

### Network Timeout

**Symptoms:**
```
curl: (28) Connection timed out after 60000 milliseconds
```

**Solutions:**
1. Verify `DTRACK_URL` is correct
2. Check network/firewall rules allow outbound HTTPS
3. For self-hosted runners: Ensure connectivity to DTrack instance
4. Test connectivity:
   ```bash
   curl -v https://api.dtrack.example.com/api/version
   ```

### Empty or Minimal SBOM

**Symptoms:**
```
SBOM generated successfully (145 bytes)
```

**Solutions:**
1. Verify dependency files exist in scanned directory:
   - **Node.js:** `package.json`, `package-lock.json`, `yarn.lock`
   - **Python:** `requirements.txt`, `Pipfile.lock`, `poetry.lock`
   - **PHP:** `composer.json`, `composer.lock`
   - **Java:** `pom.xml`, `build.gradle`
   - **Go:** `go.mod`
   - **Ruby:** `Gemfile.lock`
2. Check `PROJECT_PATH` points to correct directory (monorepo)
3. Download and inspect SBOM artifact:
   - **GitLab CI:** Job → Browse → `bom.json`
   - **GitHub Actions:** Workflow run → Artifacts → `dtrack-sbom-*`

### Duplicate Projects in Monorepo

**Symptoms:**
- All components uploading to same DTrack project
- Project shows mixed dependencies from multiple components

**Solutions:**

**GitLab CI:** Set unique `PROJECT_PATH` for each component:
```yaml
dtrack-upload:frontend:
  extends: dtrack-upload
  variables:
    PROJECT_PATH: frontend  # Creates "platform/repo/frontend"

dtrack-upload:backend:
  extends: dtrack-upload
  variables:
    PROJECT_PATH: backend   # Creates "platform/repo/backend"
```

**GitHub Actions:** Set unique `project_path` in separate workflows:
```yaml
# frontend.yml
with:
  project_path: frontend  # Creates "org/repo/frontend"

# backend.yml
with:
  project_path: backend   # Creates "org/repo/backend"
```

### SBOM Upload Works But No Vulnerabilities Shown

**Symptoms:**
- Upload succeeds (HTTP 200/201)
- Project created in DTrack
- Dependency count shows correctly
- But no vulnerabilities displayed

**Solutions:**
1. **Vulnerability Analysis Not Complete:**
   - DTrack analyzes vulnerabilities asynchronously
   - Wait 1-5 minutes for initial analysis
   - Check: DTrack → System → Tasks (should show analysis jobs)

2. **Vulnerability Datasources Not Configured:**
   - DTrack → Admin → Analyzers
   - Ensure NVD, GitHub Advisories, OSS Index are enabled
   - Run manual mirror sync if needed

3. **No Known Vulnerabilities:**
   - Your dependencies may genuinely have no known CVEs
   - Check individual components in DTrack UI
   - Try uploading SBOM for project with known vulnerable dependencies (e.g., old packages) to verify DTrack is working

### Invalid SBOM Format Error

**Symptoms:**
```
✗ Failed to upload SBOM
HTTP Status Code: 400
{"error": "Invalid BOM"}
```

**Solutions:**
1. Verify Trivy version is up-to-date (uses latest CycloneDX spec)
2. Check Dependency-Track version compatibility
3. Manually validate SBOM:
   ```bash
   # Download bom.json from artifacts
   jq . bom.json  # Should be valid JSON
   ```
4. DTrack accepts CycloneDX 1.2, 1.3, 1.4, 1.5 - check Trivy output format

---

## Project Naming Patterns

### Single Project Repository

| Platform | Pattern | Example |
|----------|---------|---------|
| GitLab CI | `${CI_PROJECT_PATH}` | `platform/my-app` |
| GitHub Actions | `${GITHUB_REPOSITORY}` | `myorg/my-app` |

### Monorepo Components

| Platform | Pattern | Example |
|----------|---------|---------|
| GitLab CI | `${CI_PROJECT_PATH}/${PROJECT_PATH}` | `platform/my-app/frontend`<br>`platform/my-app/backend` |
| GitHub Actions | `${GITHUB_REPOSITORY}/${project_path}` | `myorg/my-app/frontend`<br>`myorg/my-app/backend` |

---

## Advanced Usage

### Conditional Upload Based on Branch

**GitLab CI:**
```yaml
dtrack-upload:
  rules:
    - if: '$ENABLE_DTRACK == "true" && $CI_COMMIT_BRANCH == "main"'
    - if: '$ENABLE_DTRACK == "true" && $CI_COMMIT_TAG'
```

**GitHub Actions:**
```yaml
on:
  push:
    branches:
      - main
      - release/*
    tags:
      - 'v*'
```

### Different Projects for Branches

**GitLab CI:**
```yaml
dtrack-upload:staging:
  extends: dtrack-upload
  variables:
    DTRACK_PROJECT_NAME: "platform/my-app-staging"
    DTRACK_PROJECT_VERSION: "staging-${CI_COMMIT_SHORT_SHA}"
  rules:
    - if: '$ENABLE_DTRACK == "true" && $CI_COMMIT_BRANCH == "staging"'

dtrack-upload:production:
  extends: dtrack-upload
  variables:
    DTRACK_PROJECT_NAME: "platform/my-app"
    DTRACK_PROJECT_VERSION: "${CI_COMMIT_TAG}"
  rules:
    - if: '$ENABLE_DTRACK == "true" && $CI_COMMIT_TAG'
```

### Fail Pipeline on High-Risk Vulnerabilities

DTrack upload always succeeds if API accepts SBOM. To fail based on vulnerabilities, use DTrack's Policy Compliance feature or separate vulnerability check job.

---

## Best Practices

1. **Use Auto-Create Mode:** Simpler setup, automatic project creation
2. **Set API Key as Secret:** Never commit API keys to repository
3. **Monorepo: Unique PROJECT_PATH:** Each component gets its own DTrack project
4. **Version Tagging:** Use semantic versioning for releases
5. **Scheduled Audits:** Run daily SBOM uploads on main branch to track dependency drift
6. **Artifact Retention:** Keep `bom.json` artifacts for debugging (7 days default)
7. **Monitor DTrack UI:** Review vulnerabilities regularly, set up alerts
8. **Policy Enforcement:** Configure DTrack policies for automatic violation detection

---

## Getting Help

- **Template Issues:** [GitHub Issues](https://github.com/platform/devsecops-template/issues)
- **Dependency-Track Docs:** [docs.dependencytrack.org](https://docs.dependencytrack.org)
- **Trivy SBOM Docs:** [aquasecurity.github.io/trivy](https://aquasecurity.github.io/trivy)
