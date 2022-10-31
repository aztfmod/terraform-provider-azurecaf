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
			"sensitive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not display the value in the log is the value is sensitive (default: false).",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the environment variable.",
			},
			"value_sensitive": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Value (sensitive) of the environment variable.",
			},
		},
	}
}

func resourceAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	value := os.Getenv(name)

	if d.Get("fails_if_empty").(bool) {
		return diag.Errorf("Value is not set for environment variable: %s", name)
	}

	d.SetId(name)

	if d.Get("sensitive").(bool) {
		_ = d.Set("value_sensitive", value)
		_ = d.Set("value", "")
	} else {
		_ = d.Set("value_sensitive", "")
		_ = d.Set("value", value)
	}

	return diags
}