// Package azurecaf implements the Azure Cloud Adoption Framework (CAF) Terraform provider.
//
// This package provides resources and data sources for generating Azure resource names
// that comply with the Microsoft Cloud Adoption Framework naming conventions and Azure
// resource naming requirements.
//
// Key components:
//   - azurecaf_name resource: Creates names with full validation and customization options
//   - azurecaf_naming_convention resource: Legacy naming convention resource (deprecated)
//   - azurecaf_name data source: Generates names during plan phase for early validation
//   - azurecaf_environment_variable data source: Retrieves environment variables
//
// The provider supports multiple naming conventions including CAF classic, CAF random,
// passthrough, and fully random naming strategies.
package azurecaf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the configured Terraform provider schema with all supported
// resources and data sources.
//
// Resources:
//   - azurecaf_naming_convention: Legacy naming convention resource (use azurecaf_name instead)
//   - azurecaf_name: Primary resource for generating Azure-compliant resource names
//
// Data Sources:
//   - azurecaf_environment_variable: Retrieves environment variables with validation
//   - azurecaf_name: Generates names during plan phase for early validation
//
// The provider requires no configuration parameters and works out-of-the-box with
// the built-in Azure resource definitions.
func Provider() *schema.Provider {
	return &schema.Provider{
		// No provider-level configuration required
		Schema: map[string]*schema.Schema{},

		// Resources that can be created and managed
		ResourcesMap: map[string]*schema.Resource{
			"azurecaf_naming_convention": resourceNamingConvention(), // Legacy - use azurecaf_name instead
			"azurecaf_name":              resourceName(),             // Primary naming resource
		},

		// Data sources for retrieving information
		DataSourcesMap: map[string]*schema.Resource{
			"azurecaf_environment_variable": dataEnvironmentVariable(), // Environment variable lookup
			"azurecaf_name":                 dataName(),                // Name generation during plan
		},
	}
}
