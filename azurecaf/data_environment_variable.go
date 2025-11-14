package azurecaf

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataEnvironmentVariable creates and returns the schema for the azurecaf_environment_variable data source.
//
// This data source provides a secure way to read environment variables from the system
// where Terraform is running. It includes validation and optional default values.
//
// Use cases:
//   - Reading configuration from environment variables
//   - Providing fallback values when environment variables are not set
//   - Integrating with CI/CD systems that inject configuration via environment
//
// Security note: Environment variables retrieved through this data source will be
// stored in Terraform state. Avoid using this for sensitive values that should not
// be persisted in state files.
func dataEnvironmentVariable() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceAction,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the environment variable.",
			},
			"fails_if_empty": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Throws an error if the environment variable is not set (default: false).",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the environment variable.",
				Sensitive:   true,
			},
		},
	}
}

func resourceAction(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name, ok := d.Get("name").(string)
	if !ok {
		return diag.Errorf("name must be a string")
	}
	value, ok := os.LookupEnv(name)

	if !ok {
		return diag.Errorf("Value is not set for environment variable: %s", name)
	}

	d.SetId(name)
	_ = d.Set("value", value)

	return diags
}
