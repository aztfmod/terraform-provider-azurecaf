default: help

# Add help text after each target name starting with '\#\#'
# Found here: https://gist.github.com/prwhite/8168133
.PHONY: help
help:  ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

dev_container:
	go generate
	go fmt
	go build -o ~/.terraform.d/plugins/linux_amd64/terraform-provider-azurecaf

build:	## Build the project and run unit tests
	go generate
	go fmt ./...
	go build -o ./terraform-provider-azurecaf
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -cover ./...

unittest: 	## Run unit tests without coverage
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test ./...
	tfproviderlint ./...

test_coverage: 	## Run tests with coverage reporting
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -cover ./...

test_coverage_html: 	## Generate HTML coverage report
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at: coverage.html"

test_coverage_specific: ## Run coverage-focused tests specifically
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v ./azurecaf/... -run="Test.*" -coverprofile=coverage.out

test_integration: 	## Run integration tests
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" TF_ACC=1 go test -v ./azurecaf/... -run="TestAcc"

test_data_sources: 	## Run data source integration tests
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" TF_ACC=1 go test -v ./azurecaf/... -run="TestAccDataSourcesIntegration"

test_error_handling: 	## Run error handling integration tests
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" TF_ACC=1 go test -v ./azurecaf/... -run="TestAccErrorHandling"

test_resource_naming: ## Run naming convention tests
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 TF_CLI_ARGS_init="-upgrade=false" go test -v ./azurecaf/... -run="TestAcc.*NamingConvention" -coverprofile=naming_coverage.out ./...
	go tool cover -html=naming_coverage.out -o naming_coverage.html
	@echo "Naming coverage report generated at: naming_coverage.html"

test_e2e: build	## Run end-to-end tests (requires built provider)
	go test -v ./e2e/...

test_e2e_ci: build	## Run E2E tests in CI mode
	CHECKPOINT_DISABLE=1 TF_IN_AUTOMATION=1 go test -v ./e2e/...

test_all: unittest test_integration test_e2e	## Run all tests (unit, integration, and e2e)

test_ci: unittest test_coverage	## Run CI tests (unit tests with coverage, no integration tests)

clean:	## Clean up build artifacts and test results
	rm -f coverage.out coverage.html terraform-provider-azurecaf
	go clean

test: ## Run terraform examples with local provider
	# First build the provider
	go build -o ./terraform-provider-azurecaf
	
	# Create script to set up and run the examples
	@echo '#!/bin/bash' > run_examples.sh
	@echo 'GOOS=$$(go env GOOS)' >> run_examples.sh
	@echo 'GOARCH=$$(go env GOARCH)' >> run_examples.sh
	@echo 'LOCAL_PLUGIN_DIR=~/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/$${GOOS}_$${GOARCH}' >> run_examples.sh
	@echo 'echo "Using local plugin directory: $$LOCAL_PLUGIN_DIR"' >> run_examples.sh
	@echo 'mkdir -p $$LOCAL_PLUGIN_DIR' >> run_examples.sh
	@echo 'cp ./terraform-provider-azurecaf $$LOCAL_PLUGIN_DIR/' >> run_examples.sh
	@echo '' >> run_examples.sh
	@echo '# Create development override file for examples' >> run_examples.sh
	@echo 'cat > examples/terraform.rc << EOF' >> run_examples.sh
	@echo 'provider_installation {' >> run_examples.sh
	@echo '  dev_overrides {' >> run_examples.sh
	@echo '    "aztfmod.com/arnaudlh/azurecaf" = "$${HOME}/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/$${GOOS}_$${GOARCH}"' >> run_examples.sh
	@echo '  }' >> run_examples.sh
	@echo '  direct {}' >> run_examples.sh
	@echo '}' >> run_examples.sh
	@echo 'EOF' >> run_examples.sh
	@echo '' >> run_examples.sh
	@echo '# Run terraform in examples directory using the local config' >> run_examples.sh
	@echo 'cd ./examples && TF_CLI_CONFIG_FILE=terraform.rc terraform init -upgrade && terraform plan && terraform apply -auto-approve' >> run_examples.sh
	
	# Make the script executable and run it
	@chmod +x run_examples.sh
	@./run_examples.sh
	@rm run_examples.sh

generate_resource_table:  	## Generate resource table (output only)
	cat resourceDefinition.json | jq -r '.[] | "| \(.name)| \(.slug)| \(.min_length)| \(.max_length)| \(.lowercase)| \(.validation_regex)|"'

