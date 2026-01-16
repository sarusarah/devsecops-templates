#!/bin/bash
# Local linting script for DevSecOps templates
# Run this before committing to catch linting errors early

set -e

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║           DevSecOps Templates - Local Lint Tests            ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker is required but not found${NC}"
    echo "Please install Docker to run these tests"
    exit 1
fi

echo -e "${YELLOW}→ Starting tests using Docker...${NC}"
echo ""

# Test 1: YAML Syntax Validation
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 1: YAML Syntax Validation"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker run --rm -v "$(pwd)":/workspace -w /workspace python:3.12-slim bash -c "
pip install -q yamllint pyyaml 2>&1 | tail -2
find templates/ examples/ -name '*.yml' -o -name '*.yaml' | while read f; do
  python3 -c \"import yaml; yaml.safe_load(open('\$f'))\" || exit 1
done
" && echo -e "${GREEN}✓ YAML syntax valid${NC}" || { echo -e "${RED}✗ YAML syntax errors${NC}"; exit 1; }
echo ""

# Test 2: YAML Linting (GitHub config)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 2: YAML Linting (GitHub config)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker run --rm -v "$(pwd)":/workspace -w /workspace python:3.12-slim bash -c "
pip install -q yamllint pyyaml 2>&1 | tail -2
find templates/ examples/ -name '*.yml' -o -name '*.yaml' | while read f; do
  yamllint -d '{extends: default, rules: {line-length: {max: 200}}}' \"\$f\" || exit 1
done
" && echo -e "${GREEN}✓ YAML linting passed (GitHub)${NC}" || { echo -e "${RED}✗ YAML linting failed (GitHub)${NC}"; exit 1; }
echo ""

# Test 3: YAML Linting (GitLab config)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 3: YAML Linting (GitLab config)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker run --rm -v "$(pwd)":/workspace -w /workspace python:3.12-slim bash -c "
pip install -q yamllint pyyaml 2>&1 | tail -2
find templates/ examples/ -name '*.yml' -o -name '*.yaml' | while read f; do
  yamllint -d '{extends: default, rules: {line-length: {max: 200}, comments: disable}}' \"\$f\" || exit 1
done
" && echo -e "${GREEN}✓ YAML linting passed (GitLab)${NC}" || { echo -e "${RED}✗ YAML linting failed (GitLab)${NC}"; exit 1; }
echo ""

# Test 4: Template Structure
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 4: Template Structure Validation"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker run --rm -v "$(pwd)":/workspace -w /workspace alpine:3.20 sh -c "
for template in templates/security/*.yml; do
  if ! grep -q 'OWASP SPVS' \"\$template\"; then
    echo \"Warning: \$template missing OWASP SPVS documentation\"
  fi
  if ! grep -q 'stage:' \"\$template\"; then
    echo \"Error: \$template missing stage definition\"
    exit 1
  fi
done
" && echo -e "${GREEN}✓ Template structure valid${NC}" || { echo -e "${RED}✗ Template structure invalid${NC}"; exit 1; }
echo ""

# Test 5: OWASP SPVS Stages
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 5: OWASP SPVS Stage Validation"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker run --rm -v "$(pwd)":/workspace -w /workspace alpine:3.20 sh -c "
VALID_STAGES='source build test package stage verify deploy operate report'
for template in templates/security/*.yml; do
  stages=\$(grep '^  stage:' \"\$template\" | awk '{print \$2}' || true)
  for stage in \$stages; do
    if ! echo \"\$VALID_STAGES\" | grep -wq \"\$stage\"; then
      echo \"Invalid stage '\$stage' in \$template\"
      exit 1
    fi
  done
done
" && echo -e "${GREEN}✓ OWASP SPVS stages valid${NC}" || { echo -e "${RED}✗ OWASP SPVS stage validation failed${NC}"; exit 1; }
echo ""

# Test 6: Template Composition
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 6: Template Composition Test"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker run --rm -v "$(pwd)":/workspace -w /workspace python:3.12-slim bash -c "
pip install -q pyyaml 2>&1 | tail -2
python3 << 'EOF'
import yaml, sys
templates = ['templates/base.yml', 'templates/security/secrets.yml',
             'templates/security/dependency.yml', 'templates/security/sast.yml']
jobs = {}
for t in templates:
  content = yaml.safe_load(open(t))
  if content:
    for k, v in content.items():
      if isinstance(v, dict) and 'stage' in v:
        if k in jobs:
          print(f'Conflict: {k}'); sys.exit(1)
        jobs[k] = v
print(f'Validated {len(templates)} templates, {len(jobs)} jobs')
EOF
" && echo -e "${GREEN}✓ Template composition valid${NC}" || { echo -e "${RED}✗ Template composition failed${NC}"; exit 1; }
echo ""

# Summary
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                   ✅ ALL TESTS PASSED                        ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "Your changes are ready to commit and push!"
echo ""
