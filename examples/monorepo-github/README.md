# Monorepo Example - GitHub Actions

This example demonstrates how to use the DevSecOps CI/CD templates with a monorepo containing multiple projects in **GitHub Actions**.

> **Note:** For GitLab CI examples, see [examples/monorepo-gitlab/](../monorepo-gitlab/)

## Structure

```
monorepo-github/
├── frontend/              # Node.js frontend application
│   ├── src/
│   │   └── index.js
│   └── package.json
├── backend/               # Python backend application
│   ├── src/
│   │   └── main.py
│   └── requirements.txt
└── .github/
    └── workflows/
        ├── frontend.yml   # GitHub Actions for frontend
        └── backend.yml    # GitHub Actions for backend
```

## Projects

### Frontend (Node.js)
- **Language:** JavaScript/Node.js
- **Version:** Node.js 20
- **Package Manager:** npm
- **Location:** `frontend/`

### Backend (Python)
- **Language:** Python
- **Version:** Python 3.12
- **Dependencies:** requirements.txt
- **Location:** `backend/`

## GitHub Actions Configuration

The approach uses **separate workflow files per project**:
- `frontend.yml` - Runs when `frontend/**` changes
- `backend.yml` - Runs when `backend/**` changes

**Key features:**
- **Reusable workflows:** Call centralized workflow templates
- **Change detection:** Native `paths:` filter triggers workflows
- **project_path input:** Each workflow specifies its project directory
- **Independent pipelines:** Frontend and backend run independently

**Key patterns:**
```yaml
on:
  push:
    paths:
      - 'frontend/**'

jobs:
  build:
    uses: platform/devsecops-template/.github/workflows/build-node.yml@main
    with:
      project_path: frontend
      node_version: '20'
```

## CI Stages

Each project runs the following stages:

### 1. Secret Scanning
- Scans project directory for hardcoded secrets
- Runs on every commit

### 2. Build
- Installs dependencies
- Builds the project
- Uploads artifacts

### 3. Test
- Runs unit tests
- Generates coverage reports

### 4. SAST (Static Analysis)
- Scans code for security vulnerabilities
- Uploads results to GitHub Security

### 5. Dependency Scanning
- Checks dependencies for known vulnerabilities
- Language-specific scanning

## Change Detection

Uses `paths:` filter in workflow triggers to run only when project files change:

```yaml
on:
  push:
    paths:
      - 'frontend/**'
  pull_request:
    paths:
      - 'frontend/**'
```

## Testing Locally

### Using Dagger (Recommended)
```bash
# Test frontend
cd dagger
dagger call test --source=../examples/monorepo-github/frontend --language=node

# Test backend
dagger call test --source=../examples/monorepo-github/backend --language=python
```

### Using act (GitHub Actions locally)
```bash
# Install act
brew install act  # macOS
# or follow: https://github.com/nektos/act#installation

# Test frontend workflow
act -W .github/workflows/frontend.yml

# Test backend workflow
act -W .github/workflows/backend.yml
```

## Benefits of This Approach

### ✅ Efficient CI/CD
- Only affected projects run in CI
- Faster feedback for developers
- Reduced GitHub Actions minutes/costs

### ✅ Clear Separation
- Each project has its own workflow file
- Easy to understand what runs when
- Simple to add/remove stages per project

### ✅ Independent Deployment
- Frontend and backend can be deployed independently
- Different release cadences per project
- Reduced deployment risk

### ✅ Native GitHub Features
- Uses built-in `paths:` filter (no external actions)
- Integrates with GitHub Security tab
- Clean workflow organization

### ✅ Scalable
- Easy to add new projects (just add new workflow file)
- No complex matrix or dynamic generation
- Follows KISS principle

## Adding a New Project

1. Add project directory
2. Create new workflow file (e.g., `.github/workflows/api.yml`)
3. Set paths filter and project_path input

Example `.github/workflows/api.yml`:
```yaml
name: API CI

on:
  push:
    branches: [main, develop]
    paths:
      - 'api/**'
  pull_request:
    branches: [main, develop]
    paths:
      - 'api/**'

jobs:
  secrets:
    uses: platform/devsecops-template/.github/workflows/security-secrets.yml@main
    with:
      project_path: api

  build:
    uses: platform/devsecops-template/.github/workflows/build-node.yml@main
    with:
      project_path: api
      node_version: '20'

  test:
    needs: build
    uses: platform/devsecops-template/.github/workflows/test-node.yml@main
    with:
      project_path: api

  sast:
    uses: platform/devsecops-template/.github/workflows/security-sast.yml@main
    with:
      project_path: api

  dependency-scan:
    uses: platform/devsecops-template/.github/workflows/security-dependency.yml@main
    with:
      project_path: api
      language: node
```

## Common Patterns

### Root-level Scans
Some scans should run on the entire repository - create a separate workflow without path filters:

```yaml
# .github/workflows/root-security.yml
name: Root Security Scan

on:
  push:
    branches: [main]
  pull_request:

jobs:
  secrets:
    uses: platform/devsecops-template/.github/workflows/security-secrets.yml@main
    with:
      project_path: "."
```

### Shared Dependencies
If projects share dependencies (e.g., in `shared/` directory):

```yaml
on:
  push:
    paths:
      - 'frontend/**'
      - 'shared/**'  # Also trigger if shared code changes
```

### Container Builds
For projects that build containers:

```yaml
container-scan:
  uses: platform/devsecops-template/.github/workflows/security-container.yml@main
  with:
    image_name: ghcr.io/${{ github.repository }}/frontend
    image_tag: ${{ github.sha }}
```

### Matrix Builds (Multiple Versions)
Test against multiple Node.js versions:

```yaml
jobs:
  test:
    strategy:
      matrix:
        node-version: [18, 20, 22]
    uses: platform/devsecops-template/.github/workflows/test-node.yml@main
    with:
      project_path: frontend
      node_version: ${{ matrix.node-version }}
```

## Migration from Single Project

If you're migrating from a single-project repository:

1. **Create project directories:**
   ```bash
   mkdir -p frontend backend
   git mv src/ frontend/src/
   git mv package.json frontend/
   ```

2. **Create workflow files:**
   - Split single `.github/workflows/ci.yml` into per-project workflows
   - Add `paths:` filters to each workflow
   - Add `project_path` input to reusable workflow calls

3. **Update documentation:**
   - Update README with new structure
   - Document which project does what

4. **Validate:**
   ```bash
   # Test locally with act
   act -W .github/workflows/frontend.yml
   act -W .github/workflows/backend.yml
   ```

## Troubleshooting

### Workflows not running when code changes
- Check `paths:` filter matches your directory structure
- Verify paths don't have leading slashes (use `frontend/**` not `/frontend/**`)

### Security scans scanning entire repo
- Verify `project_path` input is set correctly in workflow calls
- Check that paths don't include parent directories

### Build artifacts not found
- Check `working-directory` is set correctly in workflow steps
- Ensure artifact paths are relative to project directory

### Reusable workflow not found
- Verify the repository path and ref/tag in `uses:`
- Ensure the template repository is accessible
- Check that the workflow file exists at the specified path

## Resources

- [DevSecOps Templates Documentation](../../README.md)
- [GitLab CI Example](../monorepo-gitlab/)
- [GitHub Actions Reusable Workflows](../../templates/github/workflows/README.md)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
