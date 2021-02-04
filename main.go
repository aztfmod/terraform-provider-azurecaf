package main

import (
	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate go run gen.go

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return azurecaf.Provider()
		},
	})
}
