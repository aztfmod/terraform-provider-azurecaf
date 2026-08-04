package azurecaf

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataEnvironmentVariable creates and returns the schema for the azurecaf_environment_variable data source.
//
// This data source provides a secure way to read environment variables from the system
// where Terraform is running.
//
// Use cases:
//   - Reading configuration from environment variables
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
				Description: "Throws an error if the environment variable is set to an empty value. Missing variables always produce an error (default: false).",
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

func resourceAction(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	failsIfEmpty := d.Get("fails_if_empty").(bool)
	value, ok := os.LookupEnv(name)

	if !ok {
		return diag.Errorf("Environment variable is not set: %s", name)
	}
	if failsIfEmpty && value == "" {
		return diag.Errorf("Environment variable is empty: %s", name)
	}

	d.SetId(name)
	if err := d.Set("value", value); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set environment variable value: %w", err))
	}

	return nil
}
