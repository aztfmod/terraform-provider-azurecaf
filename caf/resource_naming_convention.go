package caf

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNamingConvention() *schema.Resource {
	return &schema.Resource{
		Create: resourceNamingConventionCreate,
		Read:   resourceNamingConventionRead,
		Update: resourceNamingConventionUpdate,
		Delete: resourceNamingConventionDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"convention": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ConventionCafRandom),
				ValidateFunc: validation.StringInSlice([]string{
					string(ConventionCafClassic),
					string(ConventionCafRandom),
					string(ConventionRandom),
					string(ConventionPassThrough),
				}, false),
			},
			"prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ResourceTypeRg),
					string(ResourceTypeSt),
				}, false),
			},
			"generated_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNamingConventionCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	d.SetId(name)

	generateName(d, m)

	return resourceNamingConventionRead(d, m)
}

func resourceNamingConventionRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceNamingConventionUpdate(d *schema.ResourceData, m interface{}) error {
	generateName(d, m)

	return resourceNamingConventionRead(d, m)
}

func resourceNamingConventionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func updateAddress(d *schema.ResourceData, m interface{}) error {
	return nil
}

func generateName(d *schema.ResourceData, m interface{}) error {
	prefix := d.Get("prefix").(string)
	name := d.Get("name").(string)

	d.Set("generated_name", prefix+name)

	return nil
}
