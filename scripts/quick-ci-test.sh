#!/bin/bash

# Quick CI Test with Act
# This script runs a minimal E2E test using act to validate CI setup

set -e

echo "🎭 Quick CI Test with Act..."
echo

# Check dependencies
if ! command -v act &> /dev/null; then
    echo "❌ Act not installed. Install with: brew install act"
    exit 1
fi

if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker not running. Please start Docker Desktop."
    exit 1
fi

echo "✅ Act and Docker are ready"
echo

# Test E2E workflow in dry-run mode first
echo "1️⃣ Testing E2E workflow structure (dry-run)..."
act pull_request --job e2e-tests -n

echo
echo "2️⃣ Testing specific E2E job (actual run - just the setup)..."
echo "   This will only run the setup steps to validate the environment"

# Run just the first few steps to validate setup
act pull_request --job e2e-tests \
    --env CHECKPOINT_DISABLE=1 \
    --env TF_IN_AUTOMATION=1 \
    --env TF_CLI_ARGS_init="-upgrade=false" \
    --verbose

echo
echo "✅ Quick CI test completed!"
echo
echo "📋 Summary:"
echo "   ✅ Workflow structure is valid"
echo "   ✅ Dependencies can be installed in CI"
echo "   ✅ Environment setup works"
echo
echo "🚀 To run full E2E tests in CI environment:"
echo "   ./scripts/test-ci-with-act.sh"
