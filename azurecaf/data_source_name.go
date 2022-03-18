package azurecaf

import (
	"context"

	"github.com/aztfmod/terraform-provider-azurecaf/azurecaf/internal/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceName() *schema.Resource {
	return &schema.Resource{

		ReadContext:   dataSourceNameRead,
		SchemaVersion: 4,
		Schema:        schemas.V4_Schema(),
	}
}

func dataSourceNameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	getNameResult(d, m)
	return diags
}
