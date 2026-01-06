# GitHub Actions

Fast validation workflow for the DevSecOps Template Library.

## Workflow Jobs

### Validation
- `lint-yaml` - Validates YAML syntax
- `lint-templates` - Checks template structure
- `lint-stages` - Validates OWASP SPVS stage naming
- `security-secrets` - Scans for secrets (Gitleaks)
- `security-trivy` - Scans for vulnerabilities (Trivy)

### Test
- `test-composition` - Validates template composition

### Release
- `release-package` - Creates release artifacts (on releases)

## Triggers

- Push to `main` or `develop` branches
- Pull requests
- Release events

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

- **Workflow**: `.github/workflows/ci.yml`
- **Dependabot**: `.github/dependabot.yml` (automatic dependency updates)
