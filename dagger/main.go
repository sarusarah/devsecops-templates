// DevSecOps CI/CD Template Testing with Dagger
//
// This Dagger module tests the GitLab DevSecOps templates locally before pushing to GitLab.
// It replicates the security scanning pipeline stages to catch issues early.

package main

import (
	"context"
	"dagger/devsecops/internal/dagger"
	"fmt"
)

type Devsecops struct{}

// Test runs all security scans on a project
func (m *Devsecops) Test(
	ctx context.Context,
	// +required
	source *dagger.Directory,
	// Language of the project (node, python, php)
	// +default="node"
	language string,
) error {
	fmt.Println("🔒 Running DevSecOps pipeline tests...")

	// Run all security scans in parallel (Dagger handles parallelization automatically)

	// 1. Secrets Detection
	secretsScan := m.SecretsDetection(ctx, source)

	// 2. Dependency Scanning
	depScan := m.DependencyScanning(ctx, source, language)

	// 3. SAST Scanning
	sastScan := m.SastScanning(ctx, source)

	// Wait for all scans
	if _, err := secretsScan.Sync(ctx); err != nil {
		return fmt.Errorf("secrets detection failed: %w", err)
	}

	if _, err := depScan.Sync(ctx); err != nil {
		return fmt.Errorf("dependency scanning failed: %w", err)
	}

	if _, err := sastScan.Sync(ctx); err != nil {
		return fmt.Errorf("SAST scanning failed: %w", err)
	}

	fmt.Println("✅ All security scans passed!")
	return nil
}

// SecretsDetection scans for secrets using Gitleaks
func (m *Devsecops) SecretsDetection(
	ctx context.Context,
	// +required
	source *dagger.Directory,
) *dagger.Container {
	fmt.Println("🔍 Running secrets detection with Gitleaks...")

	return dag.Container().
		From("zricethezav/gitleaks:v8.21.2").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{
			"gitleaks", "detect",
			"--redact",
			"--source", ".",
			"--report-path", "gitleaks-report.json",
			"--report-format", "json",
			"--no-git",
		})
}

// DependencyScanning scans dependencies for vulnerabilities
func (m *Devsecops) DependencyScanning(
	ctx context.Context,
	// +required
	source *dagger.Directory,
	// +required
	language string,
) *dagger.Container {
	fmt.Printf("📦 Running dependency scanning for %s...\n", language)

	switch language {
	case "node":
		return dag.Container().
			From("node:20-alpine").
			WithMountedDirectory("/src", source).
			WithWorkdir("/src").
			WithExec([]string{"sh", "-c", "npm audit --json > dependency-scan.json || true"})

	case "python":
		return dag.Container().
			From("python:3.12-slim").
			WithMountedDirectory("/src", source).
			WithWorkdir("/src").
			WithExec([]string{"pip", "install", "-U", "pip", "pip-audit"}).
			WithExec([]string{"sh", "-c", "pip-audit -r requirements.txt -f json > dependency-scan.json || true"})

	case "php":
		return dag.Container().
			From("php:8.3-cli").
			WithMountedDirectory("/src", source).
			WithWorkdir("/src").
			WithExec([]string{"sh", "-c", "curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer"}).
			WithExec([]string{"sh", "-c", "composer audit --format=json > dependency-scan.json || true"})

	default:
		return dag.Container().
			From("alpine:3.20").
			WithMountedDirectory("/src", source).
			WithWorkdir("/src").
			WithExec([]string{"sh", "-c", "echo '{}' > dependency-scan.json"})
	}
}

// SastScanning runs static application security testing with Semgrep
func (m *Devsecops) SastScanning(
	ctx context.Context,
	// +required
	source *dagger.Directory,
) *dagger.Container {
	fmt.Println("🔬 Running SAST with Semgrep...")

	return dag.Container().
		From("returntocorp/semgrep:1.97.0").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{
			"semgrep", "scan",
			"--config", "p/security-audit",
			"--json",
			"-o", "semgrep.json",
			".",
		})
}

// ContainerScanning scans a container image with Trivy
func (m *Devsecops) ContainerScanning(
	ctx context.Context,
	// +required
	imageName string,
	// +default="latest"
	imageTag string,
) *dagger.Container {
	fmt.Printf("🐳 Scanning container %s:%s with Trivy...\n", imageName, imageTag)

	imageRef := fmt.Sprintf("%s:%s", imageName, imageTag)

	return dag.Container().
		From("aquasec/trivy:0.58.1").
		WithExec([]string{
			"trivy", "image",
			"--severity", "CRITICAL,HIGH",
			"--exit-code", "0",
			"--format", "json",
			"--output", "trivy.json",
			imageRef,
		})
}

// Build builds a Node.js application
func (m *Devsecops) BuildNode(
	ctx context.Context,
	// +required
	source *dagger.Directory,
	// +default="20"
	nodeVersion string,
	// +default="npm"
	packageManager string,
) (*dagger.Directory, error) {
	fmt.Printf("🔨 Building Node.js project with %s...\n", packageManager)

	container := dag.Container().
		From(fmt.Sprintf("node:%s-alpine", nodeVersion)).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"corepack", "enable"})

	switch packageManager {
	case "pnpm":
		container = container.
			WithExec([]string{"pnpm", "install", "--frozen-lockfile"}).
			WithExec([]string{"pnpm", "build"})
	case "yarn":
		container = container.
			WithExec([]string{"yarn", "install", "--frozen-lockfile"}).
			WithExec([]string{"yarn", "build"})
	default:
		container = container.
			WithExec([]string{"npm", "ci"}).
			WithExec([]string{"npm", "run", "build"})
	}

	return container.Directory("/src"), nil
}

// TestNode runs Node.js tests
func (m *Devsecops) TestNode(
	ctx context.Context,
	// +required
	source *dagger.Directory,
	// +default="20"
	nodeVersion string,
	// +default="npm"
	packageManager string,
) error {
	fmt.Printf("🧪 Running Node.js tests with %s...\n", packageManager)

	container := dag.Container().
		From(fmt.Sprintf("node:%s-alpine", nodeVersion)).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"corepack", "enable"})

	switch packageManager {
	case "pnpm":
		container = container.
			WithExec([]string{"pnpm", "install", "--frozen-lockfile"}).
			WithExec([]string{"pnpm", "test", "--", "--ci"})
	case "yarn":
		container = container.
			WithExec([]string{"yarn", "install", "--frozen-lockfile"}).
			WithExec([]string{"yarn", "test"})
	default:
		container = container.
			WithExec([]string{"npm", "ci"}).
			WithExec([]string{"npm", "test"})
	}

	_, err := container.Sync(ctx)
	return err
}

// ValidateYaml validates GitLab CI YAML syntax
func (m *Devsecops) ValidateYaml(
	ctx context.Context,
	// +required
	yamlFile *dagger.File,
) (string, error) {
	fmt.Println("✅ Validating GitLab CI YAML syntax...")

	output, err := dag.Container().
		From("alpine:3.20").
		WithExec([]string{"apk", "add", "--no-cache", "yq"}).
		WithMountedFile("/ci.yml", yamlFile).
		WithExec([]string{"yq", "eval", "/ci.yml"}).
		Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("YAML validation failed: %w", err)
	}

	return output, nil
}

// DtrackTest tests Dependency-Track SBOM generation and payload construction (no real upload)
func (m *Devsecops) DtrackTest(
	ctx context.Context,
	// +required
	source *dagger.Directory,
	// Project path for monorepo support (e.g., "frontend")
	// +optional
	projectPath string,
	// Test with project UUID method
	// +optional
	testUuid string,
	// Test with project name (for auto-create method)
	// +default="test-project"
	projectName string,
	// Test with project version
	// +default="1.0.0-test"
	projectVersion string,
) (string, error) {
	fmt.Println("🧪 Testing Dependency-Track SBOM generation and payload construction...")

	scanPath := "."
	if projectPath != "" {
		scanPath = projectPath
		fmt.Printf("→ Testing monorepo path: %s\n", projectPath)
	}

	// Build test script that mimics the DTrack template logic
	testScript := `
set -e

echo "================================================"
echo "1. Generating CycloneDX SBOM"
echo "================================================"
trivy fs --format cyclonedx --output bom.json ` + scanPath + `

if [ ! -f "bom.json" ]; then
  echo "✗ SBOM generation failed"
  exit 1
fi

SBOM_SIZE=$(stat -c%s bom.json)
echo "✓ SBOM generated successfully (${SBOM_SIZE} bytes)"

echo ""
echo "================================================"
echo "2. Validating SBOM structure"
echo "================================================"
if ! jq -e '.bomFormat == "CycloneDX"' bom.json > /dev/null; then
  echo "✗ Invalid CycloneDX format"
  exit 1
fi
echo "✓ Valid CycloneDX format"

COMPONENT_COUNT=$(jq '.components | length' bom.json)
echo "✓ Found ${COMPONENT_COUNT} components"

echo ""
echo "================================================"
echo "3. Testing base64 encoding"
echo "================================================"
BOM_B64=$(base64 -w 0 < bom.json)
B64_SIZE=${#BOM_B64}
echo "✓ Base64 encoded (${B64_SIZE} bytes)"

echo ""
echo "================================================"
echo "4. Testing payload construction"
echo "================================================"
`

	if testUuid != "" {
		testScript += `
# Method 1: UUID-based identification
echo "→ Testing UUID-based project identification"
PAYLOAD=$(jq -n \
  --arg project "` + testUuid + `" \
  --arg bom "$BOM_B64" \
  '{project: $project, bom: $bom}')
echo "✓ UUID payload constructed"
`
	} else {
		finalProjectName := projectName
		if projectPath != "" {
			finalProjectName = projectName + "/" + projectPath
		}
		testScript += `
# Method 2: Auto-create with name+version
echo "→ Testing auto-create project identification"
PROJECT_NAME="` + finalProjectName + `"
PROJECT_VERSION="` + projectVersion + `"
echo "   Project Name: ${PROJECT_NAME}"
echo "   Project Version: ${PROJECT_VERSION}"
PAYLOAD=$(jq -n \
  --arg name "$PROJECT_NAME" \
  --arg version "$PROJECT_VERSION" \
  --arg bom "$BOM_B64" \
  '{projectName: $name, projectVersion: $version, bom: $bom}')
echo "✓ Auto-create payload constructed"
`
	}

	testScript += `
PAYLOAD_SIZE=$(echo "$PAYLOAD" | wc -c)
echo "✓ Payload size: ${PAYLOAD_SIZE} bytes"

echo ""
echo "================================================"
echo "5. Validating payload structure"
echo "================================================"
if ! echo "$PAYLOAD" | jq -e '.bom' > /dev/null; then
  echo "✗ Missing 'bom' field in payload"
  exit 1
fi
echo "✓ Payload has 'bom' field"

if echo "$PAYLOAD" | jq -e '.project' > /dev/null; then
  echo "✓ Payload has 'project' field (UUID method)"
elif echo "$PAYLOAD" | jq -e '.projectName' > /dev/null && echo "$PAYLOAD" | jq -e '.projectVersion' > /dev/null; then
  echo "✓ Payload has 'projectName' and 'projectVersion' fields (auto-create method)"
else
  echo "✗ Invalid payload structure"
  exit 1
fi

echo ""
echo "================================================"
echo "✅ All Dependency-Track tests passed!"
echo "================================================"
echo ""
echo "NOTE: This test validates SBOM generation and payload construction."
echo "      No actual upload to Dependency-Track was performed."
echo ""
echo "To test with a real Dependency-Track instance, use:"
echo "  dagger call dtrack-upload --source=. --dtrack-url=<URL> --dtrack-api-key=env:DTRACK_KEY"
`

	container := dag.Container().
		From("aquasec/trivy:0.58.1").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq", "coreutils"}).
		WithNewFile("/test.sh", testScript).
		WithExec([]string{"sh", "/test.sh"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("DTrack test failed: %w", err)
	}

	return output, nil
}

// DtrackUpload uploads SBOM to a real Dependency-Track instance (requires credentials)
// WARNING: This performs a real upload. Use DtrackTest for validation without uploading.
func (m *Devsecops) DtrackUpload(
	ctx context.Context,
	// +required
	source *dagger.Directory,
	// +required
	dtrackUrl string,
	// +required
	dtrackApiKey *dagger.Secret,
	// Optional: Explicit project UUID (takes precedence)
	// +optional
	projectUuid string,
	// Optional: Project name (defaults to "test-project")
	// +optional
	projectName string,
	// Optional: Project version (defaults to "test")
	// +default="test"
	projectVersion string,
	// Optional: Monorepo subproject path
	// +optional
	projectPath string,
) (string, error) {
	fmt.Println("⚠️  WARNING: Performing real upload to Dependency-Track")
	fmt.Printf("→ Target: %s\n", dtrackUrl)

	scanPath := "."
	if projectPath != "" {
		scanPath = projectPath
	}

	finalProjectName := projectName
	if finalProjectName == "" {
		finalProjectName = "test-project"
	}
	if projectPath != "" {
		finalProjectName = finalProjectName + "/" + projectPath
	}

	uploadScript := `
set -e
trivy fs --format cyclonedx --output bom.json ` + scanPath + `
BOM_B64=$(base64 -w 0 < bom.json)
`

	if projectUuid != "" {
		uploadScript += `PAYLOAD=$(printf '%s' "$BOM_B64" | jq -R -s --arg project "` + projectUuid + `" '{project: $project, bom: .}')`
	} else {
		uploadScript += `PAYLOAD=$(printf '%s' "$BOM_B64" | jq -R -s --arg name "` + finalProjectName + `" --arg version "` + projectVersion + `" '{projectName: $name, projectVersion: $version, bom: .}')`
	}

	uploadScript += `
HTTP_CODE=$(printf '%s' "$PAYLOAD" | curl -w "%{http_code}" -o response.json \
  --retry 1 --retry-delay 5 --max-time 60 \
  -X PUT "` + dtrackUrl + `/api/v1/bom" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: ${DTRACK_API_KEY}" \
  --data-binary @-)

echo "HTTP Status: ${HTTP_CODE}"
if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 201 ]; then
  echo "✓ Upload successful"
  cat response.json
else
  echo "✗ Upload failed"
  cat response.json
  exit 1
fi
`

	container := dag.Container().
		From("aquasec/trivy:0.58.1").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq", "coreutils"}).
		WithSecretVariable("DTRACK_API_KEY", dtrackApiKey).
		WithNewFile("/upload.sh", uploadScript).
		WithExec([]string{"sh", "/upload.sh"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("DTrack upload failed: %w", err)
	}

	return output, nil
}
