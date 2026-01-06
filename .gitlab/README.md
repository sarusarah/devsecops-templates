# GitLab CI/CD

Fast validation pipeline for the DevSecOps Template Library.

## Pipeline Stages

### 1. Validate
- `lint:yaml` - Validates YAML syntax
- `lint:gitlab-ci` - Validates GitLab CI files
- `lint:templates` - Checks template structure
- `lint:stages` - Validates OWASP SPVS stage naming
- `lint:documentation` - Lints Markdown files
- `security:secrets` - Scans for secrets (Gitleaks)
- `security:trivy` - Scans for vulnerabilities (Trivy)

### 2. Test
- `test:composition` - Validates template composition

### 3. Release
- `release:package` - Creates release artifacts (on tags)

## Running Locally

### Test YAML
```bash
yamllint templates/
python3 -c "import yaml; yaml.safe_load(open('templates/base.yml'))"
```

### Test Templates
```bash
cd dagger
dagger call test --source=../examples/node --language=node
```

## Configuration

Edit `.gitlab-ci.yml` to customize the pipeline.

CI/CD files are in `.gitlab/ci/`:
- `lint.yml` - Linting jobs
- `contract.yml` - Contract tests (optional)
- `compliance.yml` - Compliance tests (optional)
