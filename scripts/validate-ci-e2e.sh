#!/bin/bash

# CI E2E Test Validation Script
# This script simulates the CI environment locally for testing

set -e

echo "🚀 Starting CI E2E Test Validation..."
echo

# Check dependencies
echo "📋 Checking dependencies..."
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    exit 1
fi

if ! command -v terraform &> /dev/null; then
    echo "❌ Terraform is not installed"
    exit 1
fi

echo "✅ Go version: $(go version)"
echo "✅ Terraform version: $(terraform version --json | grep version | head -1 | cut -d'"' -f4)"
echo

# Set CI environment variables
export CHECKPOINT_DISABLE=1
export TF_IN_AUTOMATION=1
export TF_CLI_ARGS_init="-upgrade=false"

echo "🔧 Environment variables set:"
echo "   CHECKPOINT_DISABLE=$CHECKPOINT_DISABLE"
echo "   TF_IN_AUTOMATION=$TF_IN_AUTOMATION"
echo "   TF_CLI_ARGS_init=$TF_CLI_ARGS_init"
echo

# Run CI-style tests
echo "🧪 Running CI-style E2E tests..."
echo

echo "1️⃣ Quick E2E Tests (Always run in CI):"
time make test_e2e_quick
echo

echo "2️⃣ Import E2E Tests (PR only):"
time make test_e2e_import
echo

echo "3️⃣ Full E2E Test Suite (PR only):"
time make test_e2e
echo

echo "4️⃣ Testing CI with Act (if available):"
if command -v act &> /dev/null && docker info > /dev/null 2>&1; then
    echo "   ✅ Act and Docker available - testing workflow structure"
    echo "   Running dry-run of E2E workflow..."
    act pull_request --job e2e-tests -n > /dev/null 2>&1
    echo "   ✅ E2E workflow structure validated with act"
else
    echo "   ⚠️ Act or Docker not available - skipping act validation"
    echo "   Install act with: brew install act"
fi
echo

echo "✅ All CI E2E tests completed successfully!"
echo "🎉 Provider is ready for CI/CD deployment!"
echo
echo "💡 Additional testing options:"
echo "   - Run ./scripts/quick-ci-test.sh for act validation"
echo "   - Run ./scripts/test-ci-with-act.sh for full CI simulation"
