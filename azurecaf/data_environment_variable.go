package azurecaf

import (
	"context"
	"os"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

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
				Description: "Through an error if the environment variable is not set (default: false).",
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
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	value := os.Getenv(name)

	if d.Get("fails_if_empty").(bool) {
		return diag.Errorf("Value is not set for environment variable: %s", name)
	}

	d.SetId(name)
	_ = d.Set("value", value)

	return diags
}