package main

import (
	"context"
	"flag"
	"log"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate go run gen.go

func main() {

	var debugMode bool

	flag.BoolVar(&debugMode, "debug", true, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: func() *schema.Provider {
		return azurecaf.Provider()
	}}

	if debugMode {
		// TODO: update this string with the full name of your provider as used in your configs
		err := plugin.Debug(context.Background(), "registry.terraform.io/aztfmod/azurecaf", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
