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
