default: build

dev_container:
	go build -o $(HOME)/.terraform.d/plugins/linux_amd64/terraform-provider-azurecaf

build:
	go build -o ./terraform-provider-azurecaf

unittest:
	go test ./...
	
test:
	cd ./examples && terraform init && terraform plan && terraform apply -auto-approve