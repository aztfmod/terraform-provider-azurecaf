#!/bin/bash

# Comprehensive E2E Test Validation
# This script runs complete E2E validation including act testing

set -e

echo "🎯 Starting Comprehensive E2E Test Validation..."
echo

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Check dependencies
print_status $BLUE "📋 Checking dependencies..."
dependencies_ok=true

if ! command -v go &> /dev/null; then
    print_status $RED "❌ Go is not installed"
    dependencies_ok=false
fi

if ! command -v terraform &> /dev/null; then
    print_status $RED "❌ Terraform is not installed"
    dependencies_ok=false
fi

if ! command -v act &> /dev/null; then
    print_status $YELLOW "⚠️ Act is not installed (CI simulation will be skipped)"
    print_status $YELLOW "   Install with: brew install act"
    act_available=false
else
    act_available=true
fi

if ! docker info > /dev/null 2>&1; then
    print_status $YELLOW "⚠️ Docker is not running (act tests will be skipped)"
    print_status $YELLOW "   Start Docker Desktop to enable act testing"
    act_available=false
fi

if [ "$dependencies_ok" = false ]; then
    print_status $RED "❌ Missing required dependencies"
    exit 1
fi

print_status $GREEN "✅ Go version: $(go version | cut -d' ' -f3)"
print_status $GREEN "✅ Terraform version: $(terraform version --json 2>/dev/null | grep version | head -1 | cut -d'"' -f4)"

# Set CI environment variables
export CHECKPOINT_DISABLE=1
export TF_IN_AUTOMATION=1
export TF_CLI_ARGS_init="-upgrade=false"

print_status $BLUE "🔧 Environment variables set for testing"
echo

# Phase 1: Local E2E Tests
print_status $BLUE "🧪 Phase 1: Local E2E Testing..."
echo

print_status $BLUE "1️⃣ Quick E2E Tests (Core functionality):"
if time make test_e2e_quick; then
    print_status $GREEN "   ✅ Quick E2E tests passed"
else
    print_status $RED "   ❌ Quick E2E tests failed"
    exit 1
fi
echo

print_status $BLUE "2️⃣ Import E2E Tests (Import functionality):"
if time make test_e2e_import; then
    print_status $GREEN "   ✅ Import E2E tests passed"
else
    print_status $RED "   ❌ Import E2E tests failed"
    exit 1
fi
echo

print_status $BLUE "3️⃣ Complete E2E Test Suite (All tests):"
if time make test_e2e; then
    print_status $GREEN "   ✅ Complete E2E test suite passed"
else
    print_status $RED "   ❌ Complete E2E test suite failed"
    exit 1
fi
echo

# Phase 2: Act CI Simulation
if [ "$act_available" = true ]; then
    print_status $BLUE "🎭 Phase 2: CI Simulation with Act..."
    echo

    print_status $BLUE "4️⃣ Act Workflow Validation (Dry-run):"
    if act pull_request --job e2e-tests -n > /dev/null 2>&1; then
        print_status $GREEN "   ✅ Workflow structure validation passed"
    else
        print_status $RED "   ❌ Workflow structure validation failed"
        exit 1
    fi

    print_status $BLUE "5️⃣ Act Full E2E Tests (CI environment):"
    print_status $YELLOW "   🔄 Running full E2E test suite in CI container..."
    print_status $YELLOW "   This may take several minutes..."
    
    if act workflow_dispatch --job e2e-tests --input test_type=all --env CHECKPOINT_DISABLE=1 --env TF_IN_AUTOMATION=1 > /tmp/act_output.log 2>&1; then
        print_status $GREEN "   ✅ Act CI simulation passed"
        print_status $GREEN "   🎉 All tests pass in CI environment!"
    else
        print_status $RED "   ❌ Act CI simulation failed"
        print_status $RED "   📄 Last 20 lines of output:"
        tail -20 /tmp/act_output.log
        exit 1
    fi
else
    print_status $YELLOW "⚠️ Phase 2: Skipping CI simulation (act/docker not available)"
fi

echo
print_status $GREEN "🎉 Comprehensive E2E Test Validation Completed Successfully!"
echo

# Summary
print_status $BLUE "📊 Test Summary:"
print_status $GREEN "   ✅ Local Quick E2E Tests"
print_status $GREEN "   ✅ Local Import E2E Tests"  
print_status $GREEN "   ✅ Local Complete E2E Suite"

if [ "$act_available" = true ]; then
    print_status $GREEN "   ✅ CI Workflow Validation"
    print_status $GREEN "   ✅ CI E2E Test Simulation"
    print_status $GREEN "   🚀 Provider ready for production deployment!"
else
    print_status $YELLOW "   ⚠️ CI simulation skipped (install act + docker for full validation)"
    print_status $GREEN "   🚀 Provider ready for deployment (local tests passed)!"
fi

echo
print_status $BLUE "💡 Next Steps:"
print_status $BLUE "   - Commit changes to trigger GitHub Actions CI"
print_status $BLUE "   - Create pull request for comprehensive CI testing"
print_status $BLUE "   - Monitor GitHub Actions for E2E test results"

if [ "$act_available" = true ]; then
    print_status $BLUE "   - Use './scripts/quick-ci-test.sh' for quick CI validation"
    print_status $BLUE "   - Use 'act pull_request --job e2e-tests' for on-demand CI testing"
fi
