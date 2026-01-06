# Testing the DevSecOps Templates

This document describes how to test the DevSecOps CI/CD templates locally before deploying to GitLab.

## Quick Start

We provide **two testing approaches**:

1. **Dagger** (Recommended) - Fast, portable, with intelligent caching
2. **gitlab-ci-local** - Tests actual `.gitlab-ci.yml` files

---

## Option 1: Testing with Dagger (Recommended)

### Prerequisites

```bash
# Install Dagger CLI
curl -L https://dl.dagger.io/dagger/install.sh | sh

# Verify installation
dagger version
```

### Run All Security Scans

Test a project with all security scans:

```bash
cd /DevSecOps/dagger

# Test Node.js example
dagger call test --source=./examples/node --language=node

# Test Python example
dagger call test --source=./examples/python --language=python

# Test PHP example
dagger call test --source=./examples/php-symfony --language=php
```

### Run Individual Scans

```bash
# Secrets detection only
dagger call secrets-detection --source=./examples/node

# Dependency scanning
dagger call dependency-scanning --source=./examples/node --language=node

# SAST (Semgrep)
dagger call sast-scanning --source=./examples/node

# Container scanning (requires built image)
dagger call container-scanning \
  --image-name=registry.gitlab.com/myproject/app \
  --image-tag=latest
```

### Build and Test

```bash
# Build Node.js app
dagger call build-node \
  --source=./examples/node \
  --node-version=20 \
  --package-manager=pnpm \
  export --path=./build-output

# Run tests
dagger call test-node \
  --source=./examples/node \
  --node-version=20 \
  --package-manager=pnpm
```

### Validate YAML Syntax

```bash
dagger call validate-yaml --yaml-file=./examples/node/.gitlab-ci.yml
```

---

## Option 2: Testing with gitlab-ci-local

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

This will present a menu:
```
Available tests:
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
| **Speed** | ⚡ Very fast (intelligent caching) | 🐢 Slower (no advanced caching) |
| **Accuracy** | 🎯 Uses same containers as CI | 🎯 Tests actual `.gitlab-ci.yml` |
| **Setup** | Single binary install | Node.js + npm required |
| **Parallel execution** | ✅ Yes | ❌ Limited |
| **Container registry auth** | ✅ Easy | ⚠️  Complex |
| **Learning curve** | 📚 Moderate (new tool) | 📖 Easy (GitLab CI knowledge) |
| **Debugging** | 🔍 Excellent (interactive mode) | ⚠️  Limited |
| **CI integration** | ✅ Works everywhere | ⚠️  GitLab-specific |

**Recommendation:**
- Use **Dagger** for daily development and pre-commit hooks
- Use **gitlab-ci-local** for final validation before pushing

---

## Pre-Commit Hooks

### Dagger Pre-Commit Hook

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash
echo "🔒 Running security scans..."

dagger call test --source=. --language=node

if [ $? -ne 0 ]; then
    echo "❌ Security scans failed. Fix issues before committing."
    exit 1
fi

echo "✅ Security scans passed!"
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
  allow_failure: false  # Block pipeline if Dagger tests fail
```

---

## Testing Workflow

### Daily Development

1. Make changes to templates
2. Run Dagger tests locally: `dagger call test --source=./examples/node --language=node`
3. Fix any issues
4. Commit changes

### Before Pushing

1. Run full test suite: `./test-local.sh` (option 5)
2. Validate all YAML files: `./test-local.sh` (option 6)
3. Push to GitLab

### In GitLab CI

1. GitLab CI runs the actual pipeline
2. Reviews results in GitLab Security Dashboard
3. Merge request shows security findings

---

## Troubleshooting

### Dagger: "failed to connect to engine"

```bash
# Check Docker is running
docker info

# Restart Docker
sudo systemctl restart docker
```

### gitlab-ci-local: "job not found"

Check that job names in `.gitlab-ci.yml` match the templates.

### Secrets detection failing locally

Dagger uses `--no-git` flag, so it scans all files. GitLab CI only scans committed files.

---

## Advanced Usage

### Custom Security Rules

Add custom Semgrep rules:

```bash
dagger call sast-scanning \
  --source=./examples/node \
  --config="p/security-audit"
```

### Test with Different Tool Versions

Modify `dagger/main.go` to use different container versions:

```go
From("returntocorp/semgrep:1.95.0")  // Use older version
```

---

## Getting Help

- **Dagger**: https://docs.dagger.io
- **gitlab-ci-local**: https://github.com/firecow/gitlab-ci-local
- **Issues**: Create an issue in this repository

---

## Examples

### Full Pipeline Test (Node.js)

```bash
# 1. Validate YAML
dagger call validate-yaml --yaml-file=./examples/node/.gitlab-ci.yml

# 2. Build
dagger call build-node \
  --source=./examples/node \
  --package-manager=pnpm \
  export --path=./dist

# 3. Test
dagger call test-node \
  --source=./examples/node \
  --package-manager=pnpm

# 4. Security scans
dagger call test --source=./examples/node --language=node
```

### Full Pipeline Test (Python)

```bash
# Validate + Security scans
dagger call validate-yaml --yaml-file=./examples/python/.gitlab-ci.yml
dagger call test --source=./examples/python --language=python
```

### Full Pipeline Test (PHP)

```bash
# Validate + Security scans
dagger call validate-yaml --yaml-file=./examples/php-symfony/.gitlab-ci.yml
dagger call test --source=./examples/php-symfony --language=php
```
