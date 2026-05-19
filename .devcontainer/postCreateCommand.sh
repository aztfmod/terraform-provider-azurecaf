#!/usr/bin/env bash
set -euo pipefail

echo "🔧 Setting up development environment..."

# Download Go module dependencies
echo "📦 Downloading Go modules..."
go mod download

# Install development tools
echo "🛠️  Installing development tools..."
go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@v0.30.0
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

# Enable Terraform plugin cache
mkdir -p ~/.terraform.d
cat > ~/.terraform.d/plugin-cache.hcl <<'EOF'
plugin_cache_dir = "$HOME/.terraform.d/plugin-cache"
EOF
mkdir -p ~/.terraform.d/plugin-cache

echo ""
echo "✅ Development environment ready!"
echo ""
echo "Quick start:"
echo "  make build     - Build provider and run unit tests"
echo "  make unittest  - Run unit tests only"
echo "  make test      - Build and run Terraform examples"
echo "  make help      - Show all available targets"
