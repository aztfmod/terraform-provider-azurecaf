package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func V4_Schema() map[string]*schema.Schema {
	resourceMapsKeys := getResourceMaps()
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: false,
			Default:  "",
		},
		"prefixes": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Optional: true,
			ForceNew: false,
		},
		"suffixes": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Optional: true,
			ForceNew: false,
		},
		"random_length": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     false,
			ValidateFunc: validation.IntAtLeast(0),
			Default:      0,
		},
		"result": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"results": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
		},
		"separator": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: false,
			Default:  "-",
		},
		"clean_input": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: false,
			Default:  true,
		},
		"passthrough": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: false,
			Default:  false,
		},
		"resource_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
			ForceNew:     false,
		},
		"resource_types": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
			},
			Optional: true,
			ForceNew: false,
		},
		"random_seed": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: false,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return new == "0" || old == new
			},
		},
		"random_string": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"use_slug": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: false,
			Default:  true,
		},
	}
}

func V4() *schema.Resource {
	return &schema.Resource{
		Schema: V4_Schema(),
	}
}
