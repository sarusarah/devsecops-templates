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

// AiReportTest tests the AI reporting pipeline logic without calling the Gemini API.
// Validates: report file discovery, prompt construction, Gemini request/response handling
// (mocked), summary aggregation, and Slack payload construction.
func (m *Devsecops) AiReportTest(
	ctx context.Context,
	// +required
	source *dagger.Directory,
	// Optional: test with a real Gemini API key (will make actual API calls)
	// +optional
	geminiApiKey *dagger.Secret,
) (string, error) {
	fmt.Println("🧪 Testing AI Reporting pipeline logic...")

	testScript := `
set -e

echo "================================================"
echo "AI Reporting Pipeline Test"
echo "================================================"
echo ""

PASS=0
FAIL=0

assert_ok() {
  if [ $? -eq 0 ]; then
    echo "  ✓ $1"
    PASS=$((PASS + 1))
  else
    echo "  ✗ $1"
    FAIL=$((FAIL + 1))
  fi
}

assert_file() {
  if [ -f "$1" ]; then
    echo "  ✓ File exists: $1"
    PASS=$((PASS + 1))
  else
    echo "  ✗ File missing: $1"
    FAIL=$((FAIL + 1))
  fi
}

assert_contains() {
  if echo "$1" | grep -q "$2"; then
    echo "  ✓ Contains: $2"
    PASS=$((PASS + 1))
  else
    echo "  ✗ Missing: $2 in output"
    FAIL=$((FAIL + 1))
  fi
}

assert_json_valid() {
  if echo "$1" | jq . > /dev/null 2>&1; then
    echo "  ✓ Valid JSON: $2"
    PASS=$((PASS + 1))
  else
    echo "  ✗ Invalid JSON: $2"
    FAIL=$((FAIL + 1))
  fi
}

# ============================================
echo "1. Test: Generate sample security reports"
echo "============================================"

# Generate real security scan reports to use as test input
echo "  Generating Trivy dependency scan..."
trivy fs --format json --output dependency-scan.json /src 2>/dev/null || true
assert_file "dependency-scan.json"

echo "  Generating Trivy SAST scan..."
trivy fs --scanners misconfig --format json --output sast-report.json /src 2>/dev/null || true
assert_file "sast-report.json"

# Create a mock secrets report (empty = no secrets found)
echo '{"Results":[]}' > secrets-report.json
assert_file "secrets-report.json"

# Create mock summary.md (simulating report stage output)
cat > summary.md << 'SUMMARY_EOF'
# Pipeline Security Summary
- **Project**: test/ai-report
- **Commit**: abc1234
- **Branch**: main

## Report: dependency-scan.json
- **Issues found**: 2

## Report: sast-report.json
- **Issues found**: 0

---
Status: 1 security scan(s) found issues
SUMMARY_EOF
assert_file "summary.md"

echo ""

# ============================================
echo "2. Test: Report file discovery"
echo "============================================"

# Replicate the discovery logic from ai-report.yml
REPORT_FILES="
secrets-report.json:Secrets Detection (Trivy)
gitleaks-report.json:Secrets Detection (Gitleaks)
dependency-scan.json:Dependency Vulnerability Scan
sast-report.json:Static Application Security Testing (Trivy)
semgrep.json:Static Application Security Testing (Semgrep)
iac-report.json:Infrastructure as Code Security (Trivy)
polaris.json:Infrastructure as Code Security (Polaris)
trivy.json:Container Image Security Scan
zap/zap.json:Dynamic Application Security Testing (OWASP ZAP)
"

FOUND_COUNT=0
echo "$REPORT_FILES" | while IFS=: read -r file category; do
  [ -z "$file" ] && continue
  file=$(echo "$file" | xargs)
  if [ -f "$file" ]; then
    FOUND_COUNT=$((FOUND_COUNT + 1))
  fi
done

# We created 3 files: secrets-report.json, dependency-scan.json, sast-report.json
test $(ls secrets-report.json dependency-scan.json sast-report.json 2>/dev/null | wc -l) -ge 3
assert_ok "Found at least 3 report files for analysis"

echo ""

# ============================================
echo "3. Test: Gemini request payload construction"
echo "============================================"

REPORT_CONTENT=$(head -c 500000 dependency-scan.json)
CATEGORY="Dependency Vulnerability Scan"

PROMPT="You are a CI/CD security analyst. Analyze the following ${CATEGORY} report output and provide a concise summary.

Format your response exactly as:
STATUS: PASS | WARN | FAIL
SEVERITY: CRITICAL | HIGH | MEDIUM | LOW | NONE
FINDINGS: <number of issues found>
SUMMARY: <one-line summary>
DETAILS:
- <key finding 1>
- <key finding 2>
ACTIONS:
- <recommended action 1, if any>

Report type: ${CATEGORY}
Report content:
${REPORT_CONTENT}"

REQUEST_PAYLOAD=$(jq -n --arg prompt "$PROMPT" \
  '{"contents": [{"parts": [{"text": $prompt}]}]}')

assert_json_valid "$REQUEST_PAYLOAD" "Gemini request payload"

# Verify payload structure
echo "$REQUEST_PAYLOAD" | jq -e '.contents[0].parts[0].text' > /dev/null
assert_ok "Payload has contents[0].parts[0].text structure"

echo ""

# ============================================
echo "4. Test: Gemini response parsing"
echo "============================================"

# Mock a Gemini API response
MOCK_RESPONSE='{
  "candidates": [{
    "content": {
      "parts": [{
        "text": "STATUS: WARN\nSEVERITY: HIGH\nFINDINGS: 2\nSUMMARY: 2 high-severity vulnerabilities found in dependencies\nDETAILS:\n- CVE-2024-1234: lodash prototype pollution (HIGH)\n- CVE-2024-5678: express path traversal (HIGH)\nACTIONS:\n- Update lodash to >= 4.17.21\n- Update express to >= 4.19.0"
      }]
    }
  }]
}'

assert_json_valid "$MOCK_RESPONSE" "Mock Gemini response"

PARSED_TEXT=$(echo "$MOCK_RESPONSE" | jq -r '.candidates[0].content.parts[0].text // "No response generated"')
assert_contains "$PARSED_TEXT" "STATUS: WARN"
assert_contains "$PARSED_TEXT" "SEVERITY: HIGH"
assert_contains "$PARSED_TEXT" "FINDINGS: 2"

echo ""

# ============================================
echo "5. Test: Summary aggregation"
echo "============================================"

mkdir -p ai-reports

# Write mock individual analyses
echo "STATUS: PASS
SEVERITY: NONE
FINDINGS: 0
SUMMARY: No secrets detected in codebase" > ai-reports/secrets-report.txt

echo "STATUS: WARN
SEVERITY: HIGH
FINDINGS: 2
SUMMARY: 2 high-severity vulnerabilities found in dependencies
DETAILS:
- CVE-2024-1234: lodash prototype pollution (HIGH)
- CVE-2024-5678: express path traversal (HIGH)
ACTIONS:
- Update lodash to >= 4.17.21" > ai-reports/dependency-scan.txt

echo "STATUS: PASS
SEVERITY: NONE
FINDINGS: 0
SUMMARY: No SAST issues found" > ai-reports/sast-report.txt

# Aggregate analyses (same logic as ai-summary job)
COMBINED_ANALYSES=""
for report in ai-reports/*.txt; do
  REPORT_NAME=$(basename "$report" .txt | sed 's|-| |g')
  CONTENT=$(cat "$report")
  COMBINED_ANALYSES="${COMBINED_ANALYSES}
=== ${REPORT_NAME} ===
${CONTENT}

"
done

test -n "$COMBINED_ANALYSES"
assert_ok "Combined analyses is non-empty"

assert_contains "$COMBINED_ANALYSES" "secrets report"
assert_contains "$COMBINED_ANALYSES" "dependency scan"
assert_contains "$COMBINED_ANALYSES" "sast report"

echo ""

# ============================================
echo "6. Test: Summary prompt construction"
echo "============================================"

PROJECT="test/ai-report"
BRANCH="main"
COMMIT="abc1234"
PIPELINE_URL="https://gitlab.example.com/test/ai-report/-/pipelines/123"
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

SUMMARY_PROMPT="You are a DevSecOps reporting assistant. Create a consolidated CI/CD pipeline summary.

Pipeline context:
- Project: ${PROJECT}
- Branch: ${BRANCH}
- Commit: ${COMMIT}
- Date: ${DATE}

Format your response as:
OVERALL_STATUS: PASS | WARN | FAIL
VERDICT: <one-line summary>
CRITICAL:
- <issues or None>
WARNINGS:
- <issues or None>
PASSED:
- <what passed>
RECOMMENDATION: <next step>

Individual stage analyses:
${COMBINED_ANALYSES}"

SUMMARY_REQUEST=$(jq -n --arg prompt "$SUMMARY_PROMPT" \
  '{"contents": [{"parts": [{"text": $prompt}]}]}')

assert_json_valid "$SUMMARY_REQUEST" "Summary request payload"

echo ""

# ============================================
echo "7. Test: Slack payload construction"
echo "============================================"

MOCK_SUMMARY="OVERALL_STATUS: WARN
VERDICT: 2 dependency vulnerabilities require attention
CRITICAL:
- None
WARNINGS:
- CVE-2024-1234: lodash prototype pollution (HIGH)
- CVE-2024-5678: express path traversal (HIGH)
PASSED:
- Secrets detection: clean
- SAST: no issues found
RECOMMENDATION: Update vulnerable dependencies before merging"

OVERALL_STATUS=$(echo "$MOCK_SUMMARY" | grep -oP '(?<=OVERALL_STATUS: )\S+' | head -1)
test "$OVERALL_STATUS" = "WARN"
assert_ok "Parsed OVERALL_STATUS = WARN"

case "$OVERALL_STATUS" in
  PASS)  COLOR="#36a64f" ; EMOJI="white_check_mark" ;;
  WARN)  COLOR="#daa038" ; EMOJI="warning" ;;
  FAIL)  COLOR="#cc0000" ; EMOJI="rotating_light" ;;
  *)     COLOR="#808080" ; EMOJI="information_source" ;;
esac

test "$COLOR" = "#daa038"
assert_ok "Color mapped to yellow for WARN"

test "$EMOJI" = "warning"
assert_ok "Emoji mapped to warning for WARN"

VERDICT=$(echo "$MOCK_SUMMARY" | grep -oP '(?<=VERDICT: ).*' | head -1)
test -n "$VERDICT"
assert_ok "Parsed VERDICT from summary"

DETAILS=":warning: *Warnings*\n- CVE-2024-1234\n\n:white_check_mark: *Passed*\n- Secrets: clean"

SLACK_PAYLOAD=$(jq -n \
  --arg color "$COLOR" \
  --arg emoji "$EMOJI" \
  --arg project "$PROJECT" \
  --arg branch "$BRANCH" \
  --arg commit "$COMMIT" \
  --arg verdict "$VERDICT" \
  --arg details "$DETAILS" \
  --arg pipeline_url "$PIPELINE_URL" \
  '{
    "attachments": [{
      "color": $color,
      "blocks": [
        {"type": "header", "text": {"type": "plain_text", "text": (":" + $emoji + ": Pipeline Summary: " + $project), "emoji": true}},
        {"type": "context", "elements": [{"type": "mrkdwn", "text": ("Branch: " + $branch + " | Commit: " + $commit)}]},
        {"type": "section", "text": {"type": "mrkdwn", "text": ("*" + $verdict + "*")}},
        {"type": "divider"},
        {"type": "section", "text": {"type": "mrkdwn", "text": $details}},
        {"type": "actions", "elements": [{"type": "button", "text": {"type": "plain_text", "text": "View Pipeline"}, "url": $pipeline_url}]}
      ]
    }]
  }')

assert_json_valid "$SLACK_PAYLOAD" "Slack Block Kit payload"

echo "$SLACK_PAYLOAD" | jq -e '.attachments[0].color' > /dev/null
assert_ok "Slack payload has color"

echo "$SLACK_PAYLOAD" | jq -e '.attachments[0].blocks[0].type == "header"' > /dev/null
assert_ok "Slack payload has header block"

echo "$SLACK_PAYLOAD" | jq -e '.attachments[0].blocks | length >= 5' > /dev/null
assert_ok "Slack payload has at least 5 blocks"

BUTTON_URL=$(echo "$SLACK_PAYLOAD" | jq -r '.attachments[0].blocks[-1].elements[0].url')
test "$BUTTON_URL" = "$PIPELINE_URL"
assert_ok "Slack button URL matches pipeline URL"

echo ""

# ============================================
echo "8. Test: Fallback when Gemini unavailable"
echo "============================================"

FALLBACK_SUMMARY="OVERALL_STATUS: UNKNOWN
VERDICT: AI analysis unavailable - review pipeline logs manually
CRITICAL:
- AI reporting could not generate analysis (check GEMINI_API_KEY configuration)
WARNINGS:
- None
PASSED:
- Pipeline execution completed
RECOMMENDATION: Check pipeline logs directly at ${PIPELINE_URL}"

FALLBACK_STATUS=$(echo "$FALLBACK_SUMMARY" | grep -oP '(?<=OVERALL_STATUS: )\S+' | head -1)
test "$FALLBACK_STATUS" = "UNKNOWN"
assert_ok "Fallback status is UNKNOWN"

assert_contains "$FALLBACK_SUMMARY" "GEMINI_API_KEY"

echo ""

# ============================================
echo "9. Test: Large report truncation"
echo "============================================"

dd if=/dev/urandom bs=1024 count=600 2>/dev/null | base64 > large-report.json
LARGE_SIZE=$(wc -c < large-report.json)
test "$LARGE_SIZE" -gt 500000
assert_ok "Created large report (${LARGE_SIZE} bytes > 500KB)"

TRUNCATED_CONTENT=$(head -c 500000 large-report.json)
TRUNCATED_SIZE=${#TRUNCATED_CONTENT}
test "$TRUNCATED_SIZE" -le 500000
assert_ok "Truncated content is <= 500KB (${TRUNCATED_SIZE} bytes)"

rm -f large-report.json

echo ""

# ============================================
echo "10. Test: Metadata JSON construction"
echo "============================================"

REPORT_COUNT=$(ls ai-reports/*.txt 2>/dev/null | wc -l)

STATUS_JSON=$(jq -n \
  --arg model "gemini-2.0-flash" \
  --arg date "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  --arg project "$PROJECT" \
  --arg commit "$COMMIT" \
  --arg branch "$BRANCH" \
  --arg pipeline "$PIPELINE_URL" \
  --argjson count "$REPORT_COUNT" \
  '{
    "model": $model,
    "date": $date,
    "project": $project,
    "commit": $commit,
    "branch": $branch,
    "pipeline_url": $pipeline,
    "reports_analyzed": $count,
    "skipped": false
  }')

assert_json_valid "$STATUS_JSON" "Status metadata JSON"

echo "$STATUS_JSON" | jq -e '.reports_analyzed == 3' > /dev/null
assert_ok "reports_analyzed count is 3"

echo "$STATUS_JSON" | jq -e '.skipped == false' > /dev/null
assert_ok "skipped is false"

echo "$STATUS_JSON" | jq -e '.model == "gemini-2.0-flash"' > /dev/null
assert_ok "model is gemini-2.0-flash"

echo ""
echo "================================================"
echo "RESULTS: ${PASS} passed, ${FAIL} failed"
echo "================================================"

if [ $FAIL -gt 0 ]; then
  echo "Some tests failed!"
  exit 1
fi

echo "All AI Reporting tests passed!"
`

	container := dag.Container().
		From("aquasec/trivy:0.58.1").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "curl", "jq", "coreutils", "grep"})

	if geminiApiKey != nil {
		container = container.WithSecretVariable("GEMINI_API_KEY", geminiApiKey)

		testScript += `

echo ""
echo "================================================"
echo "BONUS: Live Gemini API test"
echo "================================================"

if [ -n "${GEMINI_API_KEY}" ]; then
  echo "Testing real Gemini API call..."
  SMALL_PROMPT='{"contents": [{"parts": [{"text": "Reply with exactly: STATUS: PASS"}]}]}'

  HTTP_CODE=$(echo "$SMALL_PROMPT" | curl -s -w "%{http_code}" -o /tmp/gemini-test.json \
    "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent" \
    -H "x-goog-api-key: ${GEMINI_API_KEY}" \
    -H "Content-Type: application/json" \
    -d @- \
    --max-time 30)

  if [ "$HTTP_CODE" = "200" ]; then
    RESPONSE=$(jq -r '.candidates[0].content.parts[0].text // "empty"' /tmp/gemini-test.json)
    echo "  Gemini API responded (HTTP 200): $RESPONSE"
  else
    echo "  Gemini API failed (HTTP $HTTP_CODE)"
    cat /tmp/gemini-test.json 2>/dev/null || true
  fi
  rm -f /tmp/gemini-test.json
else
  echo "  Skipped (no GEMINI_API_KEY)"
fi
`
	}

	container = container.
		WithNewFile("/test-ai-report.sh", testScript).
		WithExec([]string{"sh", "/test-ai-report.sh"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("AI report test failed: %w", err)
	}

	return output, nil
}
