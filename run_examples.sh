#!/bin/bash
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
LOCAL_PLUGIN_DIR=~/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/${GOOS}_${GOARCH}
echo "Using local plugin directory: $LOCAL_PLUGIN_DIR"
mkdir -p $LOCAL_PLUGIN_DIR
cp ./terraform-provider-azurecaf $LOCAL_PLUGIN_DIR/

# Create development override file for examples
cat > examples/terraform.rc << EOF
provider_installation {
  dev_overrides {
    "aztfmod.com/arnaudlh/azurecaf" = "${HOME}/.terraform.d/plugins/aztfmod.com/arnaudlh/azurecaf/1.0.0/${GOOS}_${GOARCH}"
  }
  direct {}
}
EOF

# Run terraform in examples directory using the local config
cd ./examples && TF_CLI_CONFIG_FILE=terraform.rc terraform init -upgrade && terraform plan && terraform apply -auto-approve
