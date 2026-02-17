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

## Other Documentation

- **[../README.md](../README.md)** - Main project README with overview and quick start
- **[../templates/security/](../templates/security/)** - Security template files with inline comments

## Contributing

When adding new documentation:

1. Place comprehensive guides in `templates/<category>/` alongside the template
2. Place quick references and fix docs in `docs/`
3. Update this README with a link to your new documentation
4. Cross-reference between documents using relative paths
