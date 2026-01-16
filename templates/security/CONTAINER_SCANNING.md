# Container Security Scanning

This template provides automated container image vulnerability scanning using [Trivy](https://aquasecurity.github.io/trivy/) as part of the GitLab CI/CD pipeline.

## Overview

The `container.yml` template scans Docker/OCI container images for:
- **Vulnerabilities** in OS packages and application dependencies
- **Misconfigurations** in container images
- **Exposed secrets** in image layers

Results are integrated with GitLab's Security Dashboard for centralized vulnerability management.

## Quick Start

### Basic Usage

```yaml
# .gitlab-ci.yml
include:
  - project: platform/devsecops-template
    ref: main
    file:
      - /templates/security/container.yml

variables:
  ENABLE_CONTAINER_SCAN: "true"
  IMAGE_TAG: "latest"
```

This will scan `${CI_REGISTRY_IMAGE}:latest` in the `package` stage.

## Configuration

### Essential Variables

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `ENABLE_CONTAINER_SCAN` | `false` | Yes | Must be `"true"` to enable scanning |
| `IMAGE_NAME` | `${CI_REGISTRY_IMAGE}` | No | Full registry path to image repository |
| `IMAGE_TAG` | `${CI_COMMIT_SHORT_SHA}` | No | Tag of the image to scan |

### Security Policy Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TRIVY_SEVERITY` | `UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL` | Severity levels to report |
| `TRIVY_EXIT_CODE` | `0` | Exit code on vulnerabilities (set to `1` to fail pipeline) |
| `TRIVY_NON_SSL` | `false` | Set to `true` for insecure/self-signed registries |

### Advanced Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CONTAINER_IMAGE_SUFFIX` | _(empty)_ | Append sub-path to image name (e.g., `staging` → `IMAGE_NAME/staging:TAG`) |

## Common Use Cases

### 1. Simple Single-Image Project

Your build pushes to `${CI_REGISTRY_IMAGE}:${CI_COMMIT_SHORT_SHA}`:

```yaml
include:
  - project: platform/devsecops-template
    file: /templates/security/container.yml

stages:
  - build
  - package

build-image:
  stage: build
  script:
    - docker build -t ${CI_REGISTRY_IMAGE}:${CI_COMMIT_SHORT_SHA} .
    - docker push ${CI_REGISTRY_IMAGE}:${CI_COMMIT_SHORT_SHA}

variables:
  ENABLE_CONTAINER_SCAN: "true"
  # IMAGE_TAG defaults to ${CI_COMMIT_SHORT_SHA}, so no need to set
```

### 2. Using `latest` Tag

Your build always pushes as `latest`:

```yaml
build-image:
  script:
    - docker build -t ${CI_REGISTRY_IMAGE}:latest .
    - docker push ${CI_REGISTRY_IMAGE}:latest

variables:
  ENABLE_CONTAINER_SCAN: "true"
  IMAGE_TAG: "latest"  # Override default

container-security-scan:
  needs:
    - build-image  # Wait for image to be pushed
```

### 3. Multi-Environment Images (Staging/Prod)

Your project builds separate images like:
- `registry.gitlab.com/myorg/myapp/staging:latest`
- `registry.gitlab.com/myorg/myapp/prod:latest`

```yaml
stages:
  - build-deployment-image
  - package

# Build jobs
build-staging:
  stage: build-deployment-image
  script:
    - docker build -t ${CI_REGISTRY_IMAGE}/staging:latest .
    - docker push ${CI_REGISTRY_IMAGE}/staging:latest
  rules:
    - if: '$CI_COMMIT_BRANCH == "staging"'

build-prod:
  stage: build-deployment-image
  script:
    - docker build -t ${CI_REGISTRY_IMAGE}/prod:latest .
    - docker push ${CI_REGISTRY_IMAGE}/prod:latest
  rules:
    - if: '$CI_COMMIT_BRANCH == "prod"'

# Override container scan for staging
container-security-scan:staging:
  extends: container-security-scan
  variables:
    CONTAINER_IMAGE_SUFFIX: "staging"
    IMAGE_TAG: "latest"
  rules:
    - if: '$ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_BRANCH == "staging"'
  needs:
    - build-staging

# Override container scan for prod
container-security-scan:prod:
  extends: container-security-scan
  variables:
    CONTAINER_IMAGE_SUFFIX: "prod"
    IMAGE_TAG: "latest"
  rules:
    - if: '$ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_BRANCH == "prod"'
  needs:
    - build-prod

# Disable default scan (since we use branch-specific ones)
container-security-scan:
  rules:
    - when: never

variables:
  ENABLE_CONTAINER_SCAN: "true"
```

### 4. External Registry (Docker Hub, Azure ACR, etc.)

```yaml
variables:
  ENABLE_CONTAINER_SCAN: "true"
  IMAGE_NAME: "docker.io/myorg/myapp"  # External registry
  IMAGE_TAG: "${CI_COMMIT_SHORT_SHA}"
  
container-security-scan:
  needs:
    - push-to-external-registry
```

### 5. Scanning with Semantic Versioning

```yaml
variables:
  ENABLE_CONTAINER_SCAN: "true"
  IMAGE_TAG: "v${CI_COMMIT_TAG}"  # e.g., v1.2.3
  
container-security-scan:
  rules:
    - if: '$ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_TAG'
```

## How It Works

### Image Scanning Process

1. **Authentication**: Template authenticates with the container registry using `CI_REGISTRY_USER` and `CI_REGISTRY_PASSWORD` (or `CI_JOB_TOKEN`)
2. **Remote Scanning**: Trivy scans the image directly from the registry using `--image-src remote` (no Docker daemon required)
3. **Vulnerability Detection**: Trivy analyzes all layers for vulnerabilities
4. **Report Generation**: 
   - JSON report saved as artifact (`trivy.json`)
   - Table output in job logs for human review
   - Integrated with GitLab Security Dashboard

### Image Reference Resolution

The template constructs the full image reference as:

```
${IMAGE_NAME}/${CONTAINER_IMAGE_SUFFIX}:${IMAGE_TAG}
```

Where:
- `IMAGE_NAME` defaults to `${CI_REGISTRY_IMAGE}`
- `CONTAINER_IMAGE_SUFFIX` is optional (empty by default)
- `IMAGE_TAG` defaults to `${CI_COMMIT_SHORT_SHA}`

**Examples:**

| IMAGE_NAME | CONTAINER_IMAGE_SUFFIX | IMAGE_TAG | Final Reference |
|------------|------------------------|-----------|-----------------|
| `registry/project` | _(empty)_ | `abc123` | `registry/project:abc123` |
| `registry/project` | `staging` | `latest` | `registry/project/staging:latest` |
| `docker.io/org/app` | _(empty)_ | `v1.0.0` | `docker.io/org/app:v1.0.0` |

## Troubleshooting

### Issue: Image Not Found

**Error:**
```
unable to find the specified image "registry/project:tag"
5 errors occurred:
  * docker error: unable to inspect the image
  * containerd error: containerd socket not found
  * podman error: unable to initialize Podman client
  * remote error: GET https://registry/v2/project/manifests/tag: 404 NOT FOUND
```

**Causes & Solutions:**

1. **Image doesn't exist with that tag**
   - Verify the image was built and pushed with the exact tag
   - Check `IMAGE_TAG` matches your build process
   
2. **Wrong image path**
   - Verify `IMAGE_NAME` is correct
   - Check if you need `CONTAINER_IMAGE_SUFFIX`

3. **Image not ready yet**
   - Add `needs:` dependency to wait for build job:
     ```yaml
     container-security-scan:
       needs:
         - build-and-push-image
     ```

4. **Wrong stage order**
   - Ensure scan runs in a stage AFTER the build stage

### Issue: Authentication Failed

**Error:**
```
remote error: GET https://registry/v2/: 401 UNAUTHORIZED
```

**Solutions:**

1. **For GitLab Container Registry:**
   - Variables are set automatically, but ensure `CI_JOB_TOKEN` has registry access
   - Check project settings → CI/CD → Token Permissions

2. **For external registries:**
   - Set registry credentials in CI/CD variables:
     ```yaml
     variables:
       CI_REGISTRY_USER: "username"
       CI_REGISTRY_PASSWORD: "${EXTERNAL_REGISTRY_TOKEN}"
     ```

3. **For self-signed certificates:**
   ```yaml
   variables:
     TRIVY_NON_SSL: "true"
   ```

### Issue: Scan Taking Too Long

**Solutions:**

1. **Use remote scanning (default):**
   - The template already uses `--image-src remote`
   - No need to pull entire image

2. **Cache Trivy DB:**
   ```yaml
   container-security-scan:
     cache:
       key: trivy-db
       paths:
         - .trivycache/
     before_script:
       - export TRIVY_CACHE_DIR=.trivycache
   ```

### Issue: Too Many False Positives

**Solutions:**

1. **Scan only high/critical:**
   ```yaml
   variables:
     TRIVY_SEVERITY: "HIGH,CRITICAL"
   ```

2. **Create `.trivyignore`:**
   ```
   # .trivyignore
   CVE-2021-12345  # False positive, not exploitable in our use case
   ```

3. **Use permissive mode on non-default branches:**
   ```yaml
   variables:
     SECURITY_POLICY: "permissive"  # Allows failures on feature branches
   ```

## Best Practices

### 1. Fail on Critical/High Vulnerabilities

```yaml
variables:
  TRIVY_SEVERITY: "CRITICAL,HIGH"
  TRIVY_EXIT_CODE: "1"  # Fail pipeline if found
```

### 2. Always Add Dependencies

```yaml
container-security-scan:
  needs:
    - build-image  # Ensure image exists before scanning
```

### 3. Scan Close to Deployment

Place in `package` stage (after build, before deploy):

```yaml
stages:
  - build
  - package  # ← Container scan runs here
  - deploy
```

### 4. Monitor Scan Results

- Check GitLab Security Dashboard regularly
- Set up alerts for new vulnerabilities
- Review scan artifacts in merge requests

### 5. Keep Base Images Updated

```dockerfile
# Bad: outdated base image with known vulnerabilities
FROM ubuntu:18.04

# Good: recent base image with patches
FROM ubuntu:22.04
```

## Integration with GitLab Security Features

### Security Dashboard

Scan results automatically appear in:
- Project → Security & Compliance → Vulnerability Report
- Merge Request → Security widget

### Dependency Scanning vs Container Scanning

| Feature | Dependency Scanning | Container Scanning |
|---------|--------------------|--------------------|
| **Scope** | Application dependencies (package.json, requirements.txt) | OS packages + app dependencies in image |
| **Stage** | `source` (early) | `package` (after build) |
| **Use Case** | Catch issues before building | Final verification before deploy |
| **Recommendation** | Enable both for comprehensive coverage | |

## Example: Complete Azure Container App Setup

See `/home/sarah/Projects/Liip/azure/containerized-demo-app/.gitlab-ci.yml` for a real-world example of multi-environment container scanning.

## Further Reading

- [Trivy Documentation](https://aquasecurity.github.io/trivy/)
- [GitLab Container Scanning](https://docs.gitlab.com/ee/user/application_security/container_scanning/)
- [OWASP Container Security](https://owasp.org/www-project-docker-top-10/)
