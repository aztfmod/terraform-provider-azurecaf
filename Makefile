default: go

go:
	go build -o $(HOME)/.terraform.d/plugins/linux_amd64/terraform-provider-caf

test:
	cd ./examples && terraform init && terraform plan && terraform apply -auto-approve