# GitLab DevSecOps CI/CD Template Library

A comprehensive, enterprise-grade GitLab CI/CD template library implementing DevSecOps best practices with automated security scanning, testing, and deployment.


## Quick Start

### Using The Templates

Include the templates you need in your `.gitlab-ci.yml` and enable desired security scans:
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

**Detailed documentation for template usage can be found [here](../README.md).**

---
## .github directory content

Ccontains GitHub-specific configuration files for automation, CI/CD workflows, and dependency management.

### Workflows

#### [`workflows/ci.yml`](workflows/ci.yml)
Continuous Integration workflow that runs on every push and pull request.

**Jobs:**
- **Validation** - Fast linting and security checks
  - `lint-yaml` - Validates YAML syntax and structure
  - `lint-templates` - Checks template structure and OWASP SPVS compliance
  - `lint-stages` - Validates stage naming conventions
  - `security-secrets` - Scans for exposed secrets (Gitleaks)
  - `security-trivy` - Vulnerability scanning (Trivy)

- **Test** - Template composition validation
  - `test-composition` - Validates template composition and job conflicts

- **Release** - Automated release packaging
  - `release-package` - Creates release artifacts on new releases

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Release events (created, published)

**Concurrency:** Automatically cancels in-progress runs when new commits are pushed

### Dependency Management

#### [`dependabot.yml`](dependabot.yml)
Automated dependency updates configuration.

**Ecosystems Monitored:**
- GitHub Actions (weekly updates on Monday)
- Go modules in `/dagger` (weekly updates on Monday)
- Docker images in workflows (weekly updates on Monday)

**Configuration:**
- Opens up to 5 PRs per ecosystem
- Auto-labels PRs with `dependencies` and ecosystem-specific tags
- Commit messages prefixed with `chore(deps):`

## Running Checks Locally

### Quick Validation
```bash
# Run all linters locally (uses Docker)
./lint-local.sh

# Individual checks
yamllint templates/ examples/
python3 -c "import yaml; yaml.safe_load(open('templates/base.yml'))"
```

### Template Testing
```bash
# Using Dagger (requires Dagger CLI)
cd dagger
dagger call test --source=../examples/node --language=node
```

### GitLab CI Validation
```bash
# Validate GitLab CI syntax
docker run --rm -v "$PWD":/work -w /work python:3.12-slim sh -c \
  "pip install -q pyyaml && python3 .gitlab/ci/validate.py"
```

## Contributing

When adding new workflows or automation:

1. **Test locally first** - Use `./lint-local.sh` to validate changes
2. **Follow naming conventions** - Use descriptive job and workflow names
3. **Add documentation** - Update this README when adding new workflows
4. **Use concurrency groups** - Prevent wasted CI minutes
5. **Pin action versions** - Use `@v4` not `@latest` for stability

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Dependabot Configuration Reference](https://docs.github.com/en/code-security/dependabot/dependabot-version-updates/configuration-options-for-the-dependabot.yml-file)
- [Workflow Syntax Reference](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)
