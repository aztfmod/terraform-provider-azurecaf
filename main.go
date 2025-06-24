// Package main provides the entry point for the Azure Cloud Adoption Framework (CAF) Terraform provider.
//
// This provider implements a set of methodologies for naming convention implementation
// including the default Microsoft Cloud Adoption Framework for Azure recommendations
// as per https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging.
//
// The provider allows users to:
//   - Generate compliant Azure resource names following CAF guidelines
//   - Clean inputs to ensure compliance with Azure naming restrictions
//   - Add prefixes, suffixes, and random characters to resource names
//   - Validate existing names against Azure resource naming rules
//   - Support multiple naming conventions (CAF classic, CAF random, passthrough, etc.)
package main

import (
	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// go:generate directive runs the code generation tool to create resource definitions
// from the resourceDefinition.json file. This ensures that all supported Azure
// resource types and their naming constraints are up-to-date.
//go:generate go run gen.go

// main initializes and serves the Terraform provider using the Terraform plugin SDK.
// The provider is configured through the azurecaf.Provider() function which defines
// the available resources and data sources.
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return azurecaf.Provider()
		},
	})
}
