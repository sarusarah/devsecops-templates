
## Changelog

### v1.0.2 (2026-01-16)

**Fixed:**
- Container scanning: Fixed image registry authentication with proper `TRIVY_USERNAME` and `TRIVY_PASSWORD` environment variables
- Container scanning: Added support for remote image scanning (`--image-src remote`) to work without Docker daemon
- Container scanning: Fixed image reference resolution for multi-environment deployments (staging/prod)
- Container scanning: Added `CONTAINER_IMAGE_SUFFIX` variable for sub-image paths
- YAML linting: Removed trailing spaces in `container.yml`
- YAML linting: Added document start marker (`---`) to all YAML files for consistency
- GitLab CI: Fixed `lint:gitlab-ci` job by using Python image instead of Alpine without Python

**Added:**
- Container scanning: Comprehensive documentation (`templates/security/CONTAINER_SCANNING.md`) with configuration examples
- Container scanning: Quick reference guide (`docs/CONTAINER_SCANNING_QUICKREF.md`) with common patterns
- Documentation: YAML linter explanation (`docs/YAML_LINT_EXPLANATION.md`)
- Documentation: Documentation index (`docs/README.md`)
- Testing: Local linting script (`lint-local.sh`) for pre-commit validation
- Container scanning: Support for insecure registries via `TRIVY_NON_SSL` variable
- Container scanning: Enhanced logging and error messages for debugging
- README: Expanded container scanning configuration section with examples

**Changed:**
- Container scanning: Improved authentication flow with fallback variables
- Container scanning: Better support for GitLab CI and external registries
- All YAML files now follow consistent formatting with document start markers

**Documentation:**
- Added 5 new documentation files totaling 900+ lines
- Improved README with container scanning troubleshooting section
- Added inline comments and usage examples in templates

### v1.0.1 (2025-12-12)

**Fixed:**
- Gitleaks report generation
- DAST policy enforcement
- PHP test dependencies
- Container registry authentication
- Mattermost webhook
- Missing GitOps templates

**Added:**
- Dagger testing module
- gitlab-ci-local configuration
- Comprehensive documentation
- GitLab Security Dashboard integration
- Pinned all tool versions

### v1.0.0 (Initial Release)
- Base templates
- Security scanning templates
- Language-specific builds
- Examples for Node.js, Python, PHP

---

## Contributing

### Template Developers

1. Make changes to templates
2. Test locally: `dagger call test --source=./examples/node --language=node`
3. Validate YAML: `dagger call validate-yaml --yaml-file=.gitlab-ci.yml`
4. Run all examples: `./test-local.sh` (option 5)
5. Update documentation
6. Create merge request

### Security Researchers


---

## License

MIT License - See LICENSE file for details

---

## Quick Links

- [Testing Guide](./TESTING.md)
- [Dagger Module](./dagger/README.md)
- [Examples](./examples/)
- [Security Templates](./templates/security/)

---
