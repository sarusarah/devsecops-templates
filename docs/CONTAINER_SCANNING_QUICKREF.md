# Container Scanning Quick Reference

## Common Configurations

### 1. Single Image with Commit SHA (Default)
```yaml
# Your build pushes: registry/project:abc123
variables:
  DEVSECOPS_ENABLE_CONTAINER_SCAN: "true"
  # DEVSECOPS_IMAGE_TAG defaults to ${CI_COMMIT_SHORT_SHA}

container-security-scan:
  needs: [build-image]
```

### 2. Single Image with `latest` Tag
```yaml
# Your build pushes: registry/project:latest
variables:
  DEVSECOPS_ENABLE_CONTAINER_SCAN: "true"
  DEVSECOPS_IMAGE_TAG: "latest"

container-security-scan:
  needs: [build-image]
```

### 3. Multi-Environment (staging/prod)
```yaml
# Your build pushes: registry/project/staging:latest
#                    registry/project/prod:latest
variables:
  DEVSECOPS_ENABLE_CONTAINER_SCAN: "true"
  DEVSECOPS_IMAGE_TAG: "latest"

container-security-scan:staging:
  extends: container-security-scan
  variables:
    DEVSECOPS_CONTAINER_IMAGE_SUFFIX: "staging"
  rules:
    - if: '$DEVSECOPS_ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_BRANCH == "staging"'
  needs: [build-staging]

container-security-scan:prod:
  extends: container-security-scan
  variables:
    DEVSECOPS_CONTAINER_IMAGE_SUFFIX: "prod"
  rules:
    - if: '$DEVSECOPS_ENABLE_CONTAINER_SCAN == "true" && $CI_COMMIT_BRANCH == "prod"'
  needs: [build-prod]

container-security-scan:
  rules:
    - when: never
```

### 4. External Registry (Docker Hub, ACR, etc.)
```yaml
# Your build pushes: docker.io/myorg/app:v1.0.0
variables:
  DEVSECOPS_ENABLE_CONTAINER_SCAN: "true"
  DEVSECOPS_IMAGE_NAME: "docker.io/myorg/app"
  DEVSECOPS_IMAGE_TAG: "v${CI_COMMIT_TAG}"

container-security-scan:
  needs: [push-to-external]
```

## Configuration Variables

| Variable | Default | Use When |
|----------|---------|----------|
| `DEVSECOPS_IMAGE_NAME` | `${CI_REGISTRY_IMAGE}` | Using external registry |
| `DEVSECOPS_IMAGE_TAG` | `${CI_COMMIT_SHORT_SHA}` | Tag differs from commit SHA |
| `DEVSECOPS_CONTAINER_IMAGE_SUFFIX` | _(empty)_ | Using sub-images (e.g., `/staging`) |
| `DEVSECOPS_TRIVY_SEVERITY` | `UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL` | Want to filter by severity |
| `DEVSECOPS_TRIVY_EXIT_CODE` | `0` | Want to fail pipeline on findings (set to `1`) |

## Image Reference Resolution

The template builds image references as:
```
${DEVSECOPS_IMAGE_NAME}/${DEVSECOPS_CONTAINER_IMAGE_SUFFIX}:${DEVSECOPS_IMAGE_TAG}
```

### Examples

| DEVSECOPS_IMAGE_NAME | SUFFIX | TAG | Result |
|------------|--------|-----|--------|
| `registry/project` | - | `abc123` | `registry/project:abc123` |
| `registry/project` | `staging` | `latest` | `registry/project/staging:latest` |
| `docker.io/org/app` | - | `v1.0.0` | `docker.io/org/app:v1.0.0` |

## Troubleshooting Checklist

### ❌ Image Not Found Error

- [ ] Does the image exist with exact tag? (`docker pull <image>` to verify)
- [ ] Does `DEVSECOPS_IMAGE_TAG` match your build process?
- [ ] Do you need `DEVSECOPS_CONTAINER_IMAGE_SUFFIX`?
- [ ] Did you add `needs:` dependency?
- [ ] Is scan running after build stage?

### ❌ Authentication Failed Error

- [ ] For GitLab registry: Is token access enabled?
- [ ] For external registry: Are credentials set in CI/CD variables?
- [ ] For self-signed certs: Set `DEVSECOPS_TRIVY_NON_SSL: "true"`

### ❌ Too Many Vulnerabilities

- [ ] Filter by severity: `DEVSECOPS_TRIVY_SEVERITY: "HIGH,CRITICAL"`
- [ ] Create `.trivyignore` file for false positives
- [ ] Use permissive mode: `DEVSECOPS_SECURITY_POLICY: "permissive"`

## Quick Diagnostic Commands

### Check if image exists
```bash
docker pull ${CI_REGISTRY_IMAGE}/${DEVSECOPS_CONTAINER_IMAGE_SUFFIX}:${DEVSECOPS_IMAGE_TAG}
```

### Manual scan locally
```bash
export IMAGE_REF="registry/project/staging:latest"
trivy image --severity HIGH,CRITICAL ${IMAGE_REF}
```

### Test authentication
```bash
echo ${CI_JOB_TOKEN} | docker login ${CI_REGISTRY} -u gitlab-ci-token --password-stdin
```

## Full Documentation

- **Comprehensive Guide**: `../templates/security/CONTAINER_SCANNING.md`
- **Fix Summary**: `CONTAINER_SCANNING_FIX.md` (this directory)
- **Main README**: `../README.md` → "Container Scanning Configuration"
