package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate go run gen.go

func main() {
	var debugMode bool
	debugCAF := strings.ToLower(os.Getenv("DEBUG_CAF_PROVIDER"))
	debugMode = len(debugCAF) > 0 && (debugCAF == "true" || debugCAF == "1" || debugCAF == "on" || debugCAF == "yes")
	opts := &plugin.ServeOpts{ProviderFunc: func() *schema.Provider {
		return azurecaf.Provider()
	}}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/aztfmod/azurecaf", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
