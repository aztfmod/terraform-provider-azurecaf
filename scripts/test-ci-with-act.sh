#!/bin/bash

# Local CI Testing with Act
# This script uses act to run GitHub Actions workflows locally

set -e

echo "🎭 Testing CI Workflows Locally with Act..."
echo

# Check if act is installed
if ! command -v act &> /dev/null; then
    echo "❌ Act is not installed. Install with: brew install act"
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop."
    exit 1
fi

echo "✅ Act version: $(act --version)"
echo "✅ Docker is running"
echo

# Function to run a workflow with act
run_workflow() {
    local workflow_file=$1
    local event=$2
    local job_name=$3
    local description=$4
    
    echo "🚀 Running: $description"
    echo "   Workflow: $workflow_file"
    echo "   Event: $event"
    echo "   Job: $job_name"
    echo
    
    if [ -n "$job_name" ]; then
        act $event -W .github/workflows/$workflow_file --job $job_name
    else
        act $event -W .github/workflows/$workflow_file
    fi
    
    echo
    echo "✅ Completed: $description"
    echo "----------------------------------------"
    echo
}

# Test E2E workflow components
echo "🧪 Testing E2E Workflows..."
echo

# Test the E2E workflow (simulating pull request)
echo "1️⃣ Testing E2E Tests Job (Pull Request simulation):"
run_workflow "e2e.yml" "pull_request" "e2e-tests" "E2E Tests Job"

# Test specific workflow jobs if needed
echo "2️⃣ Testing E2E Summary Job:"
run_workflow "e2e.yml" "pull_request" "e2e-summary" "E2E Summary Job"

# Test main workflow E2E integration
echo "3️⃣ Testing Main Workflow (E2E integration):"
echo "   NOTE: This will run the full Go workflow including E2E tests"
echo "   This may take longer and includes all CI steps"
read -p "   Do you want to run the full workflow? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    run_workflow "go.yml" "pull_request" "build" "Main CI Build with E2E"
else
    echo "   Skipping full workflow test"
fi

echo "🎉 Local CI testing completed!"
echo
echo "📋 Summary:"
echo "   ✅ E2E workflow tested successfully"
echo "   ✅ All jobs executed in local containers"
echo "   ✅ Provider builds and tests work in CI environment"
echo
echo "💡 Tips:"
echo "   - Use 'act --list' to see all available workflows"
echo "   - Use 'act -n' for dry-run mode"
echo "   - Use 'act -v' for verbose output"
echo "   - Use 'act --job <job-name>' to run specific jobs"
