# Monorepo Example - GitLab CI

This example demonstrates how to use the DevSecOps CI/CD templates with a monorepo containing multiple projects in **GitLab CI**.

> **Note:** For GitHub Actions examples, see [examples/monorepo-github/](../monorepo-github/)

## Structure

```
monorepo-gitlab/
├── frontend/              # Node.js frontend application
│   ├── src/
│   │   └── index.js
│   └── package.json
├── backend/               # Python backend application
│   ├── src/
│   │   └── main.py
│   └── requirements.txt
└── .gitlab-ci.yml        # GitLab CI configuration
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

## GitLab CI Configuration

The `.gitlab-ci.yml` file demonstrates:
- **Per-project jobs:** Separate jobs for frontend and backend
- **Change detection:** Jobs only run when their project changes
- **PROJECT_PATH variable:** Each job specifies its project directory
- **Independent pipelines:** Frontend and backend can be tested/deployed independently

**Key patterns:**
```yaml
build:frontend:
  extends: .build:node
  variables:
    PROJECT_PATH: frontend
  rules:
    - changes:
        - frontend/**/*
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
- Checks for misconfigurations

### 5. Dependency Scanning
- Checks dependencies for known vulnerabilities
- Language-specific scanning

## Change Detection

Uses `rules:changes` to trigger jobs only when project files change:

```yaml
rules:
  - changes:
      - frontend/**/*
```

## Testing Locally

### Using Dagger (Recommended)
```bash
# Test frontend
cd dagger
dagger call test --source=../examples/monorepo-gitlab/frontend --language=node

# Test backend
dagger call test --source=../examples/monorepo-gitlab/backend --language=python
```

### Using gitlab-ci-local
```bash
# Install gitlab-ci-local
npm install -g gitlab-ci-local

# Test frontend jobs
gitlab-ci-local --cwd examples/monorepo-gitlab build:frontend

# Test backend jobs
gitlab-ci-local --cwd examples/monorepo-gitlab build:backend
```

## Benefits of This Approach

### ✅ Efficient CI/CD
- Only affected projects run in CI
- Faster feedback for developers
- Reduced CI minutes/costs

### ✅ Clear Separation
- Each project has its own CI configuration
- Easy to understand what runs when
- Simple to add/remove stages per project

### ✅ Independent Deployment
- Frontend and backend can be deployed independently
- Different release cadences per project
- Reduced deployment risk

### ✅ Scalable
- Easy to add new projects (just add more jobs)
- No complex matrix or dynamic pipeline generation
- Follows KISS principle

## Adding a New Project

1. Add project directory
2. Copy existing job definitions
3. Update PROJECT_PATH variable
4. Update rules:changes paths

Example:
```yaml
build:api:
  extends: .build:node
  variables:
    PROJECT_PATH: api
  rules:
    - changes:
        - api/**/*
```

## Common Patterns

### Root-level Scans
Some scans should run on the entire repository:

```yaml
secrets-detection:root:
  extends: secrets-detection
  variables:
    PROJECT_PATH: "."
```

### Shared Dependencies
If projects share dependencies (e.g., in `shared/` directory):

```yaml
rules:
  - changes:
      - frontend/**/*
      - shared/**/*  # Also run if shared code changes
```

### Container Builds
For projects that build containers:

```yaml
container-scan:frontend:
  extends: container-security-scan
  variables:
    IMAGE_NAME: "$CI_REGISTRY_IMAGE/frontend"
    IMAGE_TAG: "$CI_COMMIT_SHORT_SHA"
```

## Migration from Single Project

If you're migrating from a single-project repository:

1. **Create project directories:**
   ```bash
   mkdir -p frontend backend
   git mv src/ frontend/src/
   git mv package.json frontend/
   ```

2. **Update CI configuration:**
   - Add PROJECT_PATH to all jobs
   - Add change detection rules
   - Test locally before committing

3. **Update documentation:**
   - Update README with new structure
   - Document which project does what

4. **Validate:**
   ```bash
   # Test that jobs run correctly
   gitlab-ci-local build:frontend
   gitlab-ci-local build:backend
   ```

## Troubleshooting

### Jobs not running when code changes
- Check `rules:changes` paths include `/**/*` suffix
- Verify paths match your directory structure

### Security scans scanning entire repo
- Verify `PROJECT_PATH` is set correctly in job variables
- Check that paths don't include parent directories

### Build artifacts not found
- Ensure artifact paths include `${PROJECT_PATH:-.}/` prefix
- Check that PROJECT_PATH matches your directory structure

## Resources

- [DevSecOps Templates Documentation](../../README.md)
- [GitHub Actions Example](../monorepo-github/)
- [Template Reference](../../templates/)
