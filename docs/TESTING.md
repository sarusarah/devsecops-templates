# Testing the DevSecOps Templates

Local testing helps you catch issues before pushing to CI/CD. We provide two testing approaches: **Dagger** (recommended) and **gitlab-ci-local**.

---

## Quick Start

**Most Common Commands:**

```bash
# Always run from the dagger/ directory
cd dagger

# Run all security scans
dagger call test --source=../examples/node --language=node

# Test your own project
dagger call test --source=/path/to/your/project --language=node

# Test Dependency-Track integration (no real upload)
dagger call dtrack-test --source=../examples/node
```

**Performance:** First run ~30-60s (downloading images), subsequent runs ~5-10s (cached!)

---

## Option 1: Dagger (Recommended)

### Why Dagger?

- **10x faster** than GitLab CI (intelligent caching)
- **Parallel execution** - All scans run simultaneously
- **Same containers** as production CI
- **Works everywhere** - Not tied to GitLab

### Prerequisites

```bash
# Install Dagger CLI
curl -L https://dl.dagger.io/dagger/install.sh | sh

# Verify installation
dagger version
```

### Available Test Functions

#### Full Security Scan Suite

Runs secrets detection, dependency scanning, and SAST:

```bash
cd dagger

# Node.js
dagger call test --source=../examples/node --language=node

# Python
dagger call test --source=../examples/python --language=python

# PHP
dagger call test --source=../examples/php-symfony --language=php
```

#### Individual Security Scans

```bash
# Secrets detection (Gitleaks)
dagger call secrets-detection --source=../examples/node

# Dependency scanning
dagger call dependency-scanning --source=../examples/node --language=node

# SAST (Semgrep)
dagger call sast-scanning --source=../examples/node

# Container scanning (requires built image)
dagger call container-scanning \
  --image-name=registry.gitlab.com/myproject/app \
  --image-tag=latest
```

#### Dependency-Track Testing

```bash
# Test SBOM generation and payload (no upload)
dagger call dtrack-test --source=../examples/node

# Test monorepo with project path
dagger call dtrack-test \
  --source=../examples/monorepo-gitlab \
  --project-path=frontend \
  --project-name=myorg/monorepo

# Upload to real Dependency-Track (requires credentials)
dagger call dtrack-upload \
  --source=../examples/node \
  --dtrack-url=https://api.dtrack.example.com \
  --dtrack-api-key=env:DTRACK_API_KEY \
  --project-name=myproject \
  --project-version=1.0.0
```

#### Build & Test

```bash
# Build Node.js app
dagger call build-node \
  --source=../examples/node \
  --node-version=20 \
  --package-manager=pnpm \
  export --path=../build-output

# Run tests
dagger call test-node \
  --source=../examples/node \
  --node-version=20 \
  --package-manager=pnpm
```

#### YAML Validation

```bash
dagger call validate-yaml --yaml-file=../examples/node/.gitlab-ci.yml
```

### Common Workflows

#### Before Committing

```bash
cd dagger
dagger call test --source=.. --language=node
```

#### Testing Monorepo Projects

```bash
cd dagger

# Test each component independently
dagger call test --source=../examples/monorepo-gitlab/frontend --language=node
dagger call test --source=../examples/monorepo-gitlab/backend --language=python

# Test with Dependency-Track monorepo support
dagger call dtrack-test \
  --source=../examples/monorepo-gitlab \
  --project-path=frontend \
  --project-name=myorg/monorepo
```

**Monorepo Benefits:**
- Each project tested independently
- Only test changed projects in CI
- Parallel execution for faster results

---

## Option 2: gitlab-ci-local

### Prerequisites

```bash
# Install gitlab-ci-local
npm install -g gitlab-ci-local

# Verify installation
gitlab-ci-local --version
```

### Interactive Testing

```bash
cd /home/sarah/Projects/Liip/DevSecOps

# Run interactive test script
./test-local.sh
```

**Menu options:**
```
1) Node.js example
2) Python example
3) PHP Symfony example
4) PHP Drupal example
5) All examples
6) Validate template YAML
```

### Manual Testing

```bash
cd examples/node

# Preview pipeline
gitlab-ci-local --preview

# Run specific job
gitlab-ci-local secrets-detection

# Run all jobs
gitlab-ci-local
```

### Configuration

Edit `.gitlab-ci-local-config.yml` to customize variables for local testing.

---

## Comparison: Dagger vs gitlab-ci-local

| Feature | Dagger | gitlab-ci-local |
|---------|--------|-----------------|
| **Speed** | Very fast (intelligent caching) | Slower |
| **Accuracy** | Same containers as CI | Tests actual `.gitlab-ci.yml` |
| **Setup** | Single binary | Node.js + npm required |
| **Parallel execution** | Yes | Limited |
| **Learning curve** | Moderate | Easy |
| **Debugging** | Excellent | Limited |

**Recommendation:**
- Use **Dagger** for daily development and pre-commit hooks
- Use **gitlab-ci-local** for final validation before pushing

---

## Pre-Commit Hooks

### Dagger Pre-Commit Hook

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash
echo "Running security scans..."

cd dagger
dagger call test --source=.. --language=node

if [ $? -ne 0 ]; then
    echo "Security scans failed. Fix issues before committing."
    exit 1
fi

echo "Security scans passed!"
```

```bash
chmod +x .git/hooks/pre-commit
```

---

## CI/CD Integration

### Add Dagger to GitLab CI

```yaml
# .gitlab-ci.yml
dagger-validation:
  stage: preflight
  image: dagger/dagger:latest
  services:
    - docker:dind
  variables:
    DOCKER_HOST: tcp://docker:2375
  script:
    - cd dagger
    - dagger call test --source=.. --language=${LANGUAGE}
  allow_failure: false
```

---

## Testing Workflow

### Daily Development

1. Make changes to templates
2. Run: `cd dagger && dagger call test --source=../examples/node --language=node`
3. Fix any issues
4. Commit changes

### Before Pushing

1. Run full test suite: `./test-local.sh` (option 5)
2. Validate YAML: `./test-local.sh` (option 6)
3. Push to GitLab

### In GitLab CI

1. GitLab CI runs the actual pipeline
2. Review results in GitLab Security Dashboard
3. Merge request shows security findings

---

## Troubleshooting

### Dagger: "module not found"

**Solution:** Always run from `dagger/` directory:
```bash
cd dagger
```

### Dagger: "failed to connect to engine"

**Solution:** Ensure Docker is running:
```bash
docker info

# If not running:
sudo systemctl start docker
```

### gitlab-ci-local: "job not found"

Check that job names in `.gitlab-ci.yml` match the templates.

### Secrets detection failing locally

Dagger uses `--no-git` flag, so it scans all files. GitLab CI only scans committed files.

---

## Example Test Runs

### Full Pipeline Test (Node.js)

```bash
cd dagger

# 1. Validate YAML
dagger call validate-yaml --yaml-file=../examples/node/.gitlab-ci.yml

# 2. Build
dagger call build-node \
  --source=../examples/node \
  --package-manager=pnpm \
  export --path=../dist

# 3. Test
dagger call test-node \
  --source=../examples/node \
  --package-manager=pnpm

# 4. Security scans
dagger call test --source=../examples/node --language=node

# 5. Dependency-Track SBOM
dagger call dtrack-test --source=../examples/node
```

### Success Output

```bash
$ dagger call test --source=../examples/node --language=node

✔ connect 0.5s
✔ load module: . 0.8s
✔ devsecops: Devsecops! 0.0s
✔ .test(source: Directory!, language: "node"): Void 5.7s
 All security scans passed!
```

### Handling Failures

If a scan finds issues:

```bash
✘ .test(...): Void ERROR
✘ withExec gitleaks detect... ERROR
! Gitleaks found 3 potential secrets
! exit code: 1
```

**Action:**
1. Review the findings
2. Fix or whitelist false positives
3. Run again to verify

---

## Advanced Usage

### Custom Security Rules

```bash
dagger call sast-scanning \
  --source=../examples/node \
  --config="p/security-audit"
```

### Test with Different Tool Versions

Modify `dagger/main.go`:

```go
From("returntocorp/semgrep:1.95.0")  // Use older version
```

---

## Additional Resources

- **Dagger Documentation:** https://docs.dagger.io
- **gitlab-ci-local:** https://github.com/firecow/gitlab-ci-local
- **Dagger Module README:** [../dagger/README.md](../dagger/README.md)
