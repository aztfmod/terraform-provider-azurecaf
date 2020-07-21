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