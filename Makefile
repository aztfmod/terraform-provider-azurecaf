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

build:	## Build the project
	go generate
	go fmt ./...
	go build -o ./terraform-provider-azurecaf
	go test ./...

unittest: 	## Init go test
	go test ./...
	tfproviderlint ./...
	
test: # Start a terraform test / invisible help comment
	cd ./examples && terraform init && terraform plan && terraform apply -auto-approve

generate_resource_table:  	## Generate resource table (output only)
	cat resourceDefinition.json | jq -r '.[] | "| \(.name)| \(.slug)| \(.min_length)| \(.max_length)| \(.lowercase)| \(.validation_regex)|"'
	cat resourceDefinition_out_of_docs.json | jq -r '.[] | "| \(.name)| \(.slug)| \(.min_length)| \(.max_length)| \(.lowercase)| \(.validation_regex)|"'

