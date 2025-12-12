# DevSecOps Dagger Module

This Dagger module allows you to test your GitLab DevSecOps CI/CD pipelines locally before pushing to GitLab.

## Prerequisites

- [Dagger CLI](https://docs.dagger.io/install) v0.16.1+
- Docker or Podman

## Installation

```bash
# Install Dagger CLI
curl -L https://dl.dagger.io/dagger/install.sh | sh

# Initialize the module (already done)
cd /home/sarah/Projects/Liip/DevSecOps/dagger
dagger develop
```

## Usage

### Run All Security Scans

Test a Node.js project with all security scans:

```bash
dagger call test --source=../examples/node --language=node
```

Test a Python project:

```bash
dagger call test --source=../examples/python --language=python
```

Test a PHP project:

```bash
dagger call test --source=../examples/php-symfony --language=php
```

### Run Individual Scans

#### Secrets Detection

```bash
dagger call secrets-detection --source=../examples/node
```

#### Dependency Scanning

```bash
dagger call dependency-scanning --source=../examples/node --language=node
```

#### SAST (Semgrep)

```bash
dagger call sast-scanning --source=../examples/node
```

#### Container Scanning

```bash
dagger call container-scanning --image-name=myapp --image-tag=latest
```

### Build & Test

#### Build Node.js Application

```bash
dagger call build-node \
  --source=../examples/node \
  --node-version=20 \
  --package-manager=pnpm \
  export --path=./dist
```

#### Run Node.js Tests

```bash
dagger call test-node \
  --source=../examples/node \
  --node-version=20 \
  --package-manager=pnpm
```

### Validate YAML

Validate GitLab CI YAML syntax:

```bash
dagger call validate-yaml --yaml-file=../examples/node/.gitlab-ci.yml
```

## Pipeline Functions

| Function | Description |
|----------|-------------|
| `test` | Runs all security scans (secrets, dependencies, SAST) |
| `secrets-detection` | Scans for secrets with Gitleaks |
| `dependency-scanning` | Scans dependencies for vulnerabilities |
| `sast-scanning` | Runs SAST with Semgrep |
| `container-scanning` | Scans container images with Trivy |
| `build-node` | Builds a Node.js application |
| `test-node` | Runs Node.js tests |
| `validate-yaml` | Validates GitLab CI YAML syntax |

## Integration with GitLab CI

You can use Dagger in your GitLab CI pipeline:

```yaml
test-with-dagger:
  stage: test
  image: dagger/dagger:latest
  services:
    - docker:dind
  script:
    - dagger call test --source=. --language=node
```

## Advantages

- **Local Testing**: Test pipelines before pushing to GitLab
- **Fast Feedback**: Catch security issues locally
- **Reproducible**: Same containers as GitLab CI
- **Parallel Execution**: Dagger runs scans in parallel
- **Caching**: Intelligent caching speeds up repeated runs

## Examples

### Pre-commit Hook

Create `.git/hooks/pre-push`:

```bash
#!/bin/bash
echo "Running security scans with Dagger..."
dagger call test --source=. --language=node
if [ $? -ne 0 ]; then
  echo "Security scans failed. Fix issues before pushing."
  exit 1
fi
```

### CI/CD Integration

Add to your `.gitlab-ci.yml`:

```yaml
include:
  - local: /templates/base.yml

variables:
  LANGUAGE: "node"

# Local validation with Dagger (optional)
dagger-test:
  stage: preflight
  image: dagger/dagger:latest
  services:
    - docker:dind
  script:
    - cd dagger
    - dagger call test --source=.. --language=${LANGUAGE}
  allow_failure: true
```

## Troubleshooting

### "dagger: command not found"

Install Dagger CLI:
```bash
curl -L https://dl.dagger.io/dagger/install.sh | sh
```

### Docker daemon not running

Ensure Docker or Podman is running:
```bash
docker info
```

### Module initialization errors

Reinitialize the module:
```bash
cd dagger
dagger develop
```
