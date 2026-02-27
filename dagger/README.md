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
cd /home/DevSecOps/dagger
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

#### Dependency-Track SBOM Testing

Test SBOM generation and payload construction (no real upload):

```bash
# Test single project
dagger call dtrack-test --source=../examples/node

# Test monorepo with project path
dagger call dtrack-test --source=../examples/monorepo-gitlab --project-path=frontend

# Test with UUID-based identification
dagger call dtrack-test --source=../examples/node --test-uuid=abc-123-def

# Test with auto-create (name + version)
dagger call dtrack-test --source=../examples/node \
  --project-name=myorg/myproject \
  --project-version=1.0.0
```

Upload to real Dependency-Track instance (requires credentials):

```bash
# Upload with auto-create
dagger call dtrack-upload \
  --source=../examples/node \
  --dtrack-url=https://api.dtrack.example.com \
  --dtrack-api-key=env:DTRACK_API_KEY \
  --project-name=myorg/myproject \
  --project-version=1.0.0

# Upload with explicit UUID
dagger call dtrack-upload \
  --source=../examples/node \
  --dtrack-url=https://api.dtrack.example.com \
  --dtrack-api-key=env:DTRACK_API_KEY \
  --project-uuid=abc-123-def-456
```

#### AI Reporting Testing

Test the AI reporting pipeline logic (no Gemini API key required):

```bash
# Run all AI reporting tests with mock data
dagger call ai-report-test --source=../examples/node

# Run with live Gemini API validation
dagger call ai-report-test \
  --source=../examples/node \
  --gemini-api-key=env:AI_API_KEY
```

Validates: report file discovery, Gemini request/response handling, summary aggregation, Slack Block Kit payload construction, fallback behavior, and large report truncation.

For full testing documentation see [docs/AI_REPORTING_TESTING.md](../docs/AI_REPORTING_TESTING.md).

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
| `dtrack-test` | Tests DTrack SBOM generation and payload (no upload) |
| `dtrack-upload` | Uploads SBOM to real Dependency-Track instance |
| `ai-report-test` | Tests AI reporting pipeline logic (mock + optional live API) |
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
  - local: /templates/gitlab/base.yml

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
