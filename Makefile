default: build

dev_container:
	go generate
	go fmt
	go build -o ~/.terraform.d/plugins/linux_amd64/terraform-provider-azurecaf

build:
	go generate
	go fmt ./...
	go build -o ./terraform-provider-azurecaf

unittest:
	go test ./...
	tfproviderlint ./...
	
test:
	cd ./examples && terraform init && terraform plan && terraform apply -auto-approve

generate_resource_table:
	cat resourceDefinition.json | jq -r '.[] | "| \(.name)| \(.slug)| \(.min_length)| \(.max_length)| \(.lowercase)| \(.validation_regex)|"'
	cat resourceDefinition_out_of_docs.json | jq -r '.[] | "| \(.name)| \(.slug)| \(.min_length)| \(.max_length)| \(.lowercase)| \(.validation_regex)|"'

