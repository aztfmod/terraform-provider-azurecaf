package azurecaf

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getResourceMaps() []string {
	resourceMapsKeys := make([]string, 0, len(ResourceDefinitions))
	for k := range ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}
	return resourceMapsKeys
}

func resourceNameV2() *schema.Resource {
	resourceMapsKeys := getResourceMaps()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
				ForceNew: true,
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
			},
		},
	}
}

func resourceNameV3() *schema.Resource {
	resourceMapsKeys := getResourceMaps()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: true,
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: true,
			},
			"random_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
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
				ForceNew: true,
				Default:  "-",
			},
			"clean_input": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"passthrough": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				ForceNew:     true,
			},
			"resource_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				},
				Optional: true,
				ForceNew: true,
			},
			"random_seed": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"use_slug": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
		},
	}
}

func resourceNameStateUpgradeV2(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState["use_slug"] = true
	return rawState, nil
}

func resourceNameStateUpgradeV3(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	results := rawState["results"]
	resourceType := rawState["resource_type"].(string)
	result := rawState["result"].(string)
	content := make(map[string]interface{})
	switch v := results.(type) {
	default:
		fmt.Printf("unexpected type %T", v)
	case map[string]interface{}:
		content = v
	}
	if _, ok := content[resourceType]; !ok {
		content[resourceType] = result
	}
	rawState["results"] = content
	ids := []string{}
	for k, v := range content {
		ids = append(ids, fmt.Sprintf("%s\t%s", k, v))
	}
	rawState["id"] = b64.StdEncoding.EncodeToString([]byte(strings.Join(ids, "\n")))

	return rawState, nil
}
