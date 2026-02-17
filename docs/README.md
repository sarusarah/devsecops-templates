# Documentation

This directory contains detailed documentation for the DevSecOps CI/CD templates.

## Container Scanning Documentation

### Quick Reference (Start Here!)
**[CONTAINER_SCANNING_QUICKREF.md](CONTAINER_SCANNING_QUICKREF.md)** - One-page cheat sheet with common configurations

Use this when:
- You need a quick copy-paste solution
- You're troubleshooting an issue
- You want variable reference

### Comprehensive Guide
**[CONTAINER_SCANNING.md](CONTAINER_SCANNING.md)** - Complete reference

Use this when:
- Setting up container scanning for the first time
- Need detailed explanations
- Want to understand how it works
- Looking for advanced configurations

## Dependency-Track Integration

**[DEPENDENCY_TRACK.md](DEPENDENCY_TRACK.md)** - Complete guide for SBOM upload to Dependency-Track

Use this when:
- Setting up Dependency-Track integration
- Configuring monorepo SBOM uploads
- Troubleshooting authentication or upload issues
- Looking for examples and best practices

## Testing & Development

**[TESTING.md](TESTING.md)** - Complete guide for local testing

Use this when:
- Testing templates locally before pushing
- Setting up Dagger for development
- Configuring pre-commit hooks
- Running security scans locally

## Project Information

**[CHANGELOG.md](CHANGELOG.md)** - Version history and changes

Use this when:
- Reviewing what changed between versions
- Understanding feature additions and bug fixes
- Planning upgrades

## Other Documentation

- **[../README.md](../README.md)** - Main project README with overview and quick start
- **[../templates/gitlab/](../templates/gitlab/)** - GitLab CI templates with inline comments
- **[../templates/github/](../templates/github/)** - GitHub Actions reusable workflows
- **[../dagger/README.md](../dagger/README.md)** - Dagger module documentation

## Template Structure

The templates are organized by CI/CD system with a clear, flat structure:

```
templates/
├── gitlab/              # GitLab CI templates
│   ├── base.yml
│   ├── build.yml
│   ├── test.yml
│   ├── workflow.yml
│   ├── monitor.yml
│   ├── report.yml
│   ├── deploy-staging.yml
│   ├── deploy-production.yml
│   └── security/        # Security scanning templates
│       ├── secrets.yml
│       ├── dependency.yml
│       ├── sast.yml
│       ├── dast.yml
│       ├── container.yml
│       ├── iac.yml
│       └── dtrack.yml
└── github/              # GitHub Actions reusable workflows
    ├── build-node.yml
    ├── build-php.yml
    ├── build-python.yml
    ├── test-node.yml
    ├── test-php.yml
    ├── test-python.yml
    └── security/        # Security workflows
        ├── secrets.yml
        ├── dependency.yml
        ├── sast.yml
        ├── container.yml
        ├── iac.yml
        └── dtrack.yml
```

## Contributing

When adding new documentation:

1. Place comprehensive guides in `docs/` directory
2. Update this README with a link to your new documentation
3. Cross-reference between documents using relative paths
4. Keep template inline documentation up-to-date
