#!/bin/bash
# Local CI/CD Testing Script
# Tests GitLab CI pipelines locally using gitlab-ci-local

set -e

echo "🔧 DevSecOps Template Local Testing"
echo "===================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if gitlab-ci-local is installed
if ! command -v gitlab-ci-local &> /dev/null; then
    echo -e "${YELLOW}gitlab-ci-local not found. Installing...${NC}"
    npm install -g gitlab-ci-local
fi

# Function to test a specific example
test_example() {
    local example=$1
    local language=$2

    echo ""
    echo -e "${GREEN}Testing: ${example}${NC}"
    echo "----------------------------------------"

    cd "examples/${example}"

    # Override variables for testing
    export LANGUAGE="${language}"
    export ENABLE_DAST="false"  # Disable DAST for local testing
    export GITOPS_REPO=""       # Disable GitOps for local testing

    # Run gitlab-ci-local
    if gitlab-ci-local --preview; then
        echo -e "${GREEN}✅ ${example}: Pipeline preview successful${NC}"
    else
        echo -e "${RED}❌ ${example}: Pipeline preview failed${NC}"
        return 1
    fi

    cd ../..
}

# Main testing flow
main() {
    echo ""
    echo "Available tests:"
    echo "  1) Node.js example"
    echo "  2) Python example"
    echo "  3) PHP Symfony example"
    echo "  4) PHP Drupal example"
    echo "  5) All examples"
    echo "  6) Validate template YAML"
    echo ""

    read -p "Choose test to run (1-6): " choice

    case $choice in
        1)
            test_example "node" "node"
            ;;
        2)
            test_example "python" "python"
            ;;
        3)
            test_example "php-symfony" "php"
            ;;
        4)
            test_example "php-drupal" "php"
            ;;
        5)
            echo -e "${GREEN}Running all examples...${NC}"
            test_example "node" "node"
            test_example "python" "python"
            test_example "php-symfony" "php"
            test_example "php-drupal" "php"
            ;;
        6)
            echo -e "${GREEN}Validating YAML files...${NC}"
            for file in templates/**/*.yml examples/**/.gitlab-ci.yml .gitlab-ci.yml; do
                if [ -f "$file" ]; then
                    echo "Checking: $file"
                    if yq eval "$file" > /dev/null 2>&1; then
                        echo -e "${GREEN}✅ $file${NC}"
                    else
                        echo -e "${RED}❌ $file - Invalid YAML${NC}"
                    fi
                fi
            done
            ;;
        *)
            echo -e "${RED}Invalid choice${NC}"
            exit 1
            ;;
    esac

    echo ""
    echo -e "${GREEN}✅ Testing complete!${NC}"
}

# Run main function
main
